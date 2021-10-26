package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
	log "github.com/sirupsen/logrus"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/models"
)

// AllHealthyServices lists all healthy services
func (repo *DBRepo) AllHealthyServices(w http.ResponseWriter, r *http.Request) {
	repo.allServicesOfStatus(w, r, "warning")
}

// AllWarningServices lists all warning services
func (repo *DBRepo) AllWarningServices(w http.ResponseWriter, r *http.Request) {
	repo.allServicesOfStatus(w, r, "warning")
}

// AllProblemServices lists all problem services
func (repo *DBRepo) AllProblemServices(w http.ResponseWriter, r *http.Request) {
	repo.allServicesOfStatus(w, r, "problem")
}

// AllPendingServices lists all pending services
func (repo *DBRepo) AllPendingServices(w http.ResponseWriter, r *http.Request) {
	repo.allServicesOfStatus(w, r, "pending")
}

func (repo *DBRepo) allServicesOfStatus(w http.ResponseWriter, r *http.Request, status string) {
	var hostServices []*models.HostService
	hostServices, err := repo.DB.GetAllHostServicesWithStatus(status)
	if err != nil {
		log.Error(err)
		ServerError(w, r, err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("hostServices", hostServices)
	// The template name being rendered is the same as the status
	err = helpers.RenderPage(w, r, status, vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
