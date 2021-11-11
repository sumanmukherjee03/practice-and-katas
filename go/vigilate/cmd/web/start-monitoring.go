package main

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type job struct {
	HostServiceID int
}

func (j job) Run() {
}

func startMonitoring() {
	preferenceID, err := strconv.Atoi(preferenceMap["monitoring_live"])
	if err != nil {
		log.Error(err)
		return
	}
	if preferenceID == 1 {
		data := make(map[string]string)
		data["message"] = "starting"
		// trigger a message to broadcast to all clients letting them know that the app is starting to monitor
		// get all the services that we want to monitor
		// range through the services
		//   get the schedule unit and number
		//   create a job
		//   save the id of the job so that we can start and stop it
		//   broadcast over websockets the fact that the service is scheduled for monitoring
		// end range
	}
}
