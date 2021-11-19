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

	var h models.Host
	h, err = repo.DB.GetHostById(hs.HostID)
	if err != nil {
		if err != nil {
			log.Error(fmt.Errorf("ERROR - Could not find host with host id from host-services - %v", err))
			return
		}
	}
	hs.Host = h

	var hasHostServiceStatusChanged = false
	oldStatus := hs.Status
	msg, newStatus := repo.testServiceForHost(hs)
	if newStatus != oldStatus {
		hasHostServiceStatusChanged = true
	}

	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not perform check and update DB for service on host - %v", err))
		return
	}

	if hasHostServiceStatusChanged {
		payload := make(map[string]string)
		payload["message"] = fmt.Sprintf("HostService status changed from %s to %s", oldStatus, newStatus)
		payload["old_status"] = oldStatus
		payload["new_status"] = newStatus
		broadcastMessage("public-channel", "HostServiceStatusChanged", payload)
		if oldStatus == "healthy" && newStatus != "healthy" {
			// Also, if appropriate, send an email or sms
			log.Info("Send an email or sms indicating that a service is misbehaving")
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
	}

	log.Info(msg)
}

// ToggleServiceForHost handles the association or dissociation of a host with a service
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

func (repo *DBRepo) testServiceForHost(hs models.HostService) (string, string) {
	var msg, newStatus string
	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = repo.testHTTPServiceForHost(hs.Host.URL)
		break
	}
	return msg, newStatus
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
