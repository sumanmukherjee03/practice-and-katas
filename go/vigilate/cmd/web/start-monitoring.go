package main

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type job struct {
	HostServiceID int
}

func (j job) Run() {
	repo.ScheduledCheck(j.HostServiceID)
}

func startMonitoring() {
	preferenceID, err := strconv.Atoi(preferenceMap["monitoring_live"])
	if err != nil {
		log.Error(err)
		return
	}
	if preferenceID == 1 {
		// data is the payload that is sent via websockets to all the clients
		data := make(map[string]string)
		data["message"] = "Monitoring is starting"
		// trigger a message to broadcast to all clients letting them know that the app is starting to monitor
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
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
			fmt.Println(hs.ID, hs.Host.HostName, hs.Service.ServiceName)
		}
		// get all the services that we want to monitor
		// range through the services
		//   get the schedule unit and number
		//   create a job
		//   save the id of the job so that we can start and stop it
		//   broadcast over websockets the fact that the service is scheduled for monitoring
		// end range
	}
}
