package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tsawler/vigilate/internal/models"
)

const (
	HTTP           = 2 // id of the HTTP service in the database so that we can use it for pattern matching
	HTTPS          = 3 // id of the HTTPS service in the database so that we can use it for pattern matching
	SSLCertificate = 4 // id of the SSLCertificate service in the database so that we can use it for pattern matching
)

type performCheckOnServiceForHostResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostID        int       `json:"host_id"`
	HostServiceID int       `json:"host_service_id"`
	LastCheck     time.Time `json:"last_check"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
}

// ScheduledCheck performs a scheduled check on a host service id
func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	var hs models.HostService
	var err error

	hs, err = repo.DB.GetHostServiceById(hostServiceID)
	if err != nil {
		if err != nil {
			log.Error(fmt.Errorf("ERROR - Could not find host-service with provided id - %v", err))
			return
		}
	}

	msg, newStatus := repo.testServiceForHost(hs)
	if newStatus != hs.Status {
		repo.updateHostServiceStatusCount(hs, newStatus, msg)
	}
}

// PerformCheckOnServiceForHost handles the request for a manual check on a service in a host
func (repo *DBRepo) PerformCheckOnServiceForHost(w http.ResponseWriter, r *http.Request) {
	var h models.Host
	var s models.Service

	hostID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not read url param id to get host id - %v", err))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	serviceID, err := strconv.Atoi(chi.URLParam(r, "service_id"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not read url param service_id to get service id - %v", err))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	// If there is an existing host, retrieve that from the DB
	if hostID == 0 || serviceID == 0 {
		log.Error(fmt.Errorf("ERROR - Either host id or service id value is not valid - host id : %d, service id : %d", hostID, serviceID))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	h, err = repo.DB.GetHostById(hostID)
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not find host by id %d - %v", hostID, err))
		ClientErrorJSON(w, r, http.StatusNotFound)
		return
	}

	s, err = repo.DB.GetServiceById(serviceID)
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not find service by id %d - %v", serviceID, err))
		ClientErrorJSON(w, r, http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not parse form data : %v", err))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	oldStatus := r.Form.Get("old_status")
	if len(oldStatus) == 0 {
		log.Error(fmt.Errorf("ERROR - old status value is not valid - old status : %s", oldStatus))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	var hs models.HostService
	hs, err = repo.DB.GetHostServiceByHostAndService(h.ID, s.ID)
	if err != nil {
		if err != nil {
			log.Error(fmt.Errorf("ERROR - Could not find host-service with provided host id and service id - %v", err))
			ServerErrorJSON(w, r, err)
			return
		}
	}

	msg, newStatus := repo.testServiceForHost(hs)
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not perform check and update DB for service on host - %v", err))
		ServerErrorJSON(w, r, err)
		return
	}

	var resp performCheckOnServiceForHostResp
	resp.OK = true
	resp.Message = msg
	resp.HostID = hs.HostID
	resp.ServiceID = hs.ServiceID
	resp.HostServiceID = hs.ID
	resp.LastCheck = hs.LastCheck
	resp.OldStatus = oldStatus
	resp.NewStatus = hs.Status

	out, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

/////////////////////////////////////////////////////////////////////
///////////////////////////// HELPERS ///////////////////////////////
/////////////////////////////////////////////////////////////////////

func (repo *DBRepo) updateHostServiceStatusCount(hs models.HostService, newStatus string, msg string) {
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	if err := repo.DB.UpdateHostService(hs); err != nil {
		log.Error(fmt.Errorf("ERROR - Could not update DB for service on host - %v", err))
		return
	}

	pending, healthy, warning, problem, err := repo.DB.GetAllHostServiceStatusCount()
	if err != nil {
		log.Error(err)
		return
	}
	data := make(map[string]string)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["pending_count"] = strconv.Itoa(pending)
	data["warning_count"] = strconv.Itoa(warning)
	data["problem_count"] = strconv.Itoa(problem)
	broadcastMessage("public-channel", "HostServiceCountChanged", data)

	log.Info(msg)
}

func (repo *DBRepo) testServiceForHost(hs models.HostService) (string, string) {
	staleStatus := hs.Status
	var msg, newStatus string
	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = repo.testHTTPServiceForHost(hs.Host.URL)
		break
	}
	if newStatus != staleStatus {
		repo.pushStatusChangedEvent(hs, newStatus)
		event := models.Event{
			EventType:     newStatus,
			HostID:        hs.HostID,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostName:      hs.Host.HostName,
			ServiceName:   hs.Service.ServiceName,
			Message:       msg,
		}
		_, err := repo.DB.InsertEvent(event)
		if err != nil {
			log.Error(err)
			return "", ""
		}
		if staleStatus == "healthy" && newStatus != "healthy" {
			log.Info("Send an email or sms indicating that a service is misbehaving")
		}
	}
	repo.pushScheduleChangedEvent(hs, newStatus)
	// TODO : Send an email or sms notification if this needs to be notified as an alert
	return msg, newStatus
}

func (repo *DBRepo) pushStatusChangedEvent(hs models.HostService, newStatus string) {
	staleStatus := hs.Status
	payload := make(map[string]string)
	payload["host_service_id"] = strconv.Itoa(hs.ID)
	payload["host_id"] = strconv.Itoa(hs.HostID)
	payload["service_id"] = strconv.Itoa(hs.ServiceID)
	payload["host_name"] = hs.Host.HostName
	payload["service_name"] = hs.Service.ServiceName
	payload["new_status"] = newStatus
	payload["icon"] = hs.Service.Icon
	payload["message"] = fmt.Sprintf("%s on %s status changed from %s to %s", hs.Service.ServiceName, hs.Host.HostName, staleStatus, newStatus)
	payload["last_check"] = time.Now().Format("2006-01-02 3:04:06 PM")
	payload["stale_status"] = staleStatus
	broadcastMessage("public-channel", "HostServiceStatusChanged", payload)
}

func (repo *DBRepo) pushScheduleChangedEvent(hs models.HostService, newStatus string) {
	yearOne := time.Date(0001, 1, 1, 0, 0, 0, 1, time.UTC)
	data := make(map[string]string)
	data["schedule_id"] = strconv.Itoa(hs.ID)
	data["host_service_id"] = strconv.Itoa(hs.ID)
	data["host_id"] = strconv.Itoa(hs.HostID)
	data["service_id"] = strconv.Itoa(hs.ServiceID)
	data["host"] = hs.Host.HostName
	data["service"] = hs.Service.ServiceName
	data["last_run"] = time.Now().Format("2006-01-02 3:04:05 PM")
	nextScheduledEv := repo.App.Scheduler.Entry(repo.App.MonitorMap[hs.ID]).Next
	if nextScheduledEv.After(yearOne) {
		data["next_run"] = nextScheduledEv.Format("2006-01-02 3:04:05 PM")
	} else {
		data["next_run"] = "pending..."
	}
	data["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
	data["status"] = newStatus
	data["icon"] = hs.Service.Icon
	data["message"] = fmt.Sprintf("%s on %s schedule has updated", hs.Service.ServiceName, hs.Host.HostName)
	broadcastMessage("public-channel", "HostServiceScheduleChanged", data)
}

func (repo *DBRepo) addToMonitorMap(hs models.HostService) {
	if repo.App.PreferenceMap["monitoring_live"] == "1" {
		sch, err := hs.ScheduleText()
		if err != nil {
			log.Error(err)
			return
		}

		// This struct type is declared in handlers/start-monitoring.go
		var j job
		j.HostServiceID = hs.ID
		scheduledJobID, err := repo.App.Scheduler.AddJob(sch, j)
		if err != nil {
			log.Error(err)
			return
		}
		repo.App.MonitorMap[hs.ID] = scheduledJobID

		data := make(map[string]string)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["host"] = hs.Host.HostName
		data["service"] = hs.Service.ServiceName
		data["last_run"] = hs.LastCheck.Format("2006-01-02 3:04:05 PM")
		data["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
		data["next_run"] = "pending..."
		data["message"] = "scheduling"
		broadcastMessage("public-channel", "HostServiceScheduleChanged", data)
	}
}

func (repo *DBRepo) removeFromMonitorMap(hs models.HostService) {
	if repo.App.PreferenceMap["monitoring_live"] == "1" {
		repo.App.Scheduler.Remove(repo.App.MonitorMap[hs.ID])
		data := make(map[string]string)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		broadcastMessage("public-channel", "ScheduleItemRemovedEvent", data)
	}
}

func (repo *DBRepo) testHTTPServiceForHost(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	url = strings.Replace(url, "https", "http", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}
