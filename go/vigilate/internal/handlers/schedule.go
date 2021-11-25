package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	log "github.com/sirupsen/logrus"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/models"
)

// We need to sort the schedules by host and for that we are creating a new type
// and then satisfying the interface expected by sort so that we can sort the schedules by host
type ByHost []models.Schedule

func (a ByHost) Len() int           { return len(a) }
func (a ByHost) Less(i, j int) bool { return a[i].Host < a[j].Host }
func (a ByHost) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ListEntries lists schedule entries
func (repo *DBRepo) ListEntries(w http.ResponseWriter, r *http.Request) {
	var items []models.Schedule
	for k, v := range repo.App.MonitorMap {
		var item models.Schedule
		item.ID = k
		item.EntryID = v
		item.Entry = repo.App.Scheduler.Entry(v)
		hs, err := repo.DB.GetHostServiceById(k)
		if err != nil {
			log.Error("ERROR - Could not fetch host-service by id : %v", err)
			return
		}
		item.ScheduleText = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
		item.LastRunFromHostService = hs.LastCheck
		item.Host = hs.Host.HostName
		item.Service = hs.Service.ServiceName
		items = append(items, item)
	}

	// sort the slice of items by converting it to the ByHost type which satisfies the sort.Interface
	sort.Sort(ByHost(items))

	data := make(jet.VarMap)
	data.Set("items", items)

	err := helpers.RenderPage(w, r, "schedule", data, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
