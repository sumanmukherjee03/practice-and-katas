package handlers

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type job struct {
	HostServiceID int
}

func (j job) Run() {
	Repo.ScheduledCheck(j.HostServiceID)
}

func (repo *DBRepo) StartMonitoring() {
	preferenceID, err := strconv.Atoi(repo.App.PreferenceMap["monitoring_live"])
	if err != nil {
		log.Error(err)
		return
	}
	if preferenceID == 1 {
		// data is the payload that is sent via websockets to all the clients
		data := make(map[string]string)
		data["message"] = "Monitoring is starting"
		// trigger a message to broadcast to all clients letting them know that the app is starting to monitor
		// make sure that the event name used here is the same as the one used in the listener
		err := repo.App.WsClient.Trigger("public-channel", "AppStarting", data)
		if err != nil {
			log.Error(err)
			return
		}

		servicesToMonitor, err := repo.DB.GetAllHostServicesToMonitor()
		if err != nil {
			log.Error(err)
			return
		}

		for _, hs := range servicesToMonitor {
			//   get the schedule unit and number and form the cron job string
			sch, err := hs.ScheduleText()
			if err != nil {
				log.Error(err)
				return
			}

			//   create a job
			var j job
			j.HostServiceID = hs.ID
			scheduledID, err := repo.App.Scheduler.AddJob(sch, j)
			if err != nil {
				log.Error(err)
				return
			}

			//   save the id of the job in app config MonitorMap so that we can start and stop it
			repo.App.MonitorMap[hs.ID] = scheduledID

			// Generate a payload and broadcast the message that monitoring for a host service has started.
			// This is necessary to broadcast because we can start/stop monitoring with a toggle in the UI.
			// So, this code doesnt just run at startup. It can also run at other times.
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			payload["host_service_id"] = strconv.Itoa(hs.ID)

			// year1 is a dummy value, set to an old date at year 1
			year1 := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			// If the MonitorMap already had a job scheduled for this host service at some point then the next scheduled event will be after year 1
			// If the app is first starting and monitoring is enabled then the next run will be pending...
			// But if it was stopped and started somewhere in the middle then there will be a next date/time to run the check
			if repo.App.Scheduler.Entry(repo.App.MonitorMap[hs.ID]).Next.After(year1) {
				payload["next_run"] = repo.App.Scheduler.Entry(repo.App.MonitorMap[hs.ID]).Next.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["next_run"] = "pending..."
			}

			payload["host_name"] = hs.Host.HostName
			payload["service_name"] = hs.Service.ServiceName

			if hs.LastCheck.After(year1) {
				payload["last_run"] = hs.LastCheck.Format("2006-01-02 3:04:05 PM")
			} else {
				payload["last_run"] = "pending..."
			}

			payload["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)

			err = repo.App.WsClient.Trigger("public-channel", "NextRunEvent", payload)
			if err != nil {
				log.Error(err)
				return
			}

			err = repo.App.WsClient.Trigger("public-channel", "ScheduleChangedEvent", payload)
			if err != nil {
				log.Error(err)
				return
			}
		}

		repo.App.Scheduler.Start()
	}
}
