package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/driver"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/models"
	"github.com/tsawler/vigilate/internal/repository"
	"github.com/tsawler/vigilate/internal/repository/dbrepo"
)

//Repo is the repository
var Repo *DBRepo
var app *config.AppConfig

// DBRepo is the db repo
type DBRepo struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewHandlers creates the handlers
func NewHandlers(repo *DBRepo, a *config.AppConfig) {
	Repo = repo
	app = a
}

// NewPostgresqlHandlers creates db repo for postgres
func NewPostgresqlHandlers(db *driver.DB, a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// AdminDashboard displays the dashboard
func (repo *DBRepo) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	pending, healthy, warning, problem, err := repo.DB.GetAllHostServiceStatusCount()
	if err != nil {
		log.Error(err)
		ServerError(w, r, err)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("no_healthy", healthy)
	vars.Set("no_problem", problem)
	vars.Set("no_pending", pending)
	vars.Set("no_warning", warning)

	err = helpers.RenderPage(w, r, "dashboard", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Events displays the events page
func (repo *DBRepo) Events(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "events", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Settings displays the settings page
func (repo *DBRepo) Settings(w http.ResponseWriter, r *http.Request) {
	err := helpers.RenderPage(w, r, "settings", nil, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostSettings saves site settings
func (repo *DBRepo) PostSettings(w http.ResponseWriter, r *http.Request) {
	prefMap := make(map[string]string)

	prefMap["site_url"] = r.Form.Get("site_url")
	prefMap["notify_name"] = r.Form.Get("notify_name")
	prefMap["notify_email"] = r.Form.Get("notify_email")
	prefMap["smtp_server"] = r.Form.Get("smtp_server")
	prefMap["smtp_port"] = r.Form.Get("smtp_port")
	prefMap["smtp_user"] = r.Form.Get("smtp_user")
	prefMap["smtp_password"] = r.Form.Get("smtp_password")
	prefMap["sms_enabled"] = r.Form.Get("sms_enabled")
	prefMap["sms_provider"] = r.Form.Get("sms_provider")
	prefMap["twilio_phone_number"] = r.Form.Get("twilio_phone_number")
	prefMap["twilio_sid"] = r.Form.Get("twilio_sid")
	prefMap["twilio_auth_token"] = r.Form.Get("twilio_auth_token")
	prefMap["smtp_from_email"] = r.Form.Get("smtp_from_email")
	prefMap["smtp_from_name"] = r.Form.Get("smtp_from_name")
	prefMap["notify_via_sms"] = r.Form.Get("notify_via_sms")
	prefMap["notify_via_email"] = r.Form.Get("notify_via_email")
	prefMap["sms_notify_number"] = r.Form.Get("sms_notify_number")

	if r.Form.Get("sms_enabled") == "0" {
		prefMap["notify_via_sms"] = "0"
	}

	err := repo.DB.InsertOrUpdateSitePreferences(prefMap)
	if err != nil {
		log.Error(err)
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	// update app config
	for k, v := range prefMap {
		app.PreferenceMap[k] = v
	}

	app.Session.Put(r.Context(), "flash", "Changes saved")

	if r.Form.Get("action") == "1" {
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
	}
}

// AllHosts displays list of all hosts
func (repo *DBRepo) AllHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := repo.DB.AllHosts()
	if err != nil {
		ServerError(w, r, err)
		return
	}
	vars := make(jet.VarMap)
	vars.Set("hosts", hosts)
	err = helpers.RenderPage(w, r, "hosts", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// Host shows the host add/edit form
func (repo *DBRepo) Host(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not read url param id - %v", err))
	}

	vars := make(jet.VarMap)
	var h models.Host
	if id > 0 {
		h, err = repo.DB.GetHostById(id)
		if err != nil {
			ClientError(w, r, http.StatusNotFound)
			return
		}
	}
	vars.Set("host", h)

	// Make sure you pass the vars to RenderPage
	err = helpers.RenderPage(w, r, "host", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostHost handles the create/update of a host
func (repo *DBRepo) PostHost(w http.ResponseWriter, r *http.Request) {
	var h models.Host
	var hostID int

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not read url param id - %v", err))
	}

	// If there is an existing host, retrieve that from the DB
	if id > 0 {
		h, err = repo.DB.GetHostById(id)
		if err != nil {
			ClientError(w, r, http.StatusNotFound)
			return
		}
	}

	// get values from form and populate the fields of a new host or existing host
	h.HostName = r.Form.Get("host_name")
	h.CanonicalName = r.Form.Get("canonical_name")
	h.URL = r.Form.Get("url")
	h.IP = r.Form.Get("ip")
	h.IPV6 = r.Form.Get("ipv6")
	h.Location = r.Form.Get("location")
	h.OS = r.Form.Get("os")
	active, err := strconv.Atoi(r.Form.Get("active"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not read the form value for host field active - %v", err))
	}
	h.Active = active

	err = h.IsValid()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	if id > 0 {
		err = repo.DB.UpdateHost(h)
		if err != nil {
			helpers.ServerError(w, r, err)
			return
		}
		hostID = id
	} else {
		newHostID, err := repo.DB.InsertHost(h)
		if err != nil {
			helpers.ServerError(w, r, err)
			return
		}
		hostID = newHostID
	}

	repo.App.Session.Put(r.Context(), "flash", "Successfully created or updated host")
	http.Redirect(w, r, fmt.Sprintf("/admin/host/%d", hostID), http.StatusSeeOther)
}

type toggleServiceForHostResp struct {
	OK bool `json:"ok"`
}

// ToggleServiceForHost handles the association or dissociation of a host with a service
func (repo *DBRepo) ToggleServiceForHost(w http.ResponseWriter, r *http.Request) {
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

	// toggle service on or off for host
	activate, err := strconv.Atoi(r.Form.Get("activate"))
	if err != nil {
		log.Error(fmt.Errorf("ERROR - Could not get an integer value for activate from form - %v", err))
		ClientErrorJSON(w, r, http.StatusBadRequest)
		return
	}

	var hs models.HostService
	hs, err = repo.DB.GetHostServiceByHostAndService(h.ID, s.ID)
	if err != nil {
		hs.HostID = h.ID
		hs.ServiceID = s.ID
		hs.Active = activate
		_, err = repo.DB.InsertHostService(hs)
		if err != nil {
			log.Error(fmt.Errorf("ERROR - Could not associate/dissociate host with/from service - %v", err))
			ServerErrorJSON(w, r, err)
			return
		}
	} else {
		hs.Active = activate
		err = repo.DB.UpdateHostService(hs)
		if err != nil {
			log.Error(fmt.Errorf("ERROR - Could not associate/dissociate host with/from service - %v", err))
			ServerErrorJSON(w, r, err)
			return
		}
	}

	var resp toggleServiceForHostResp
	resp.OK = true
	out, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// AllUsers lists all admin users
func (repo *DBRepo) AllUsers(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	u, err := repo.DB.AllUsers()
	if err != nil {
		ClientError(w, r, http.StatusBadRequest)
		return
	}

	vars.Set("users", u)

	err = helpers.RenderPage(w, r, "users", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// OneUser displays the add/edit user page
func (repo *DBRepo) OneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
	}

	vars := make(jet.VarMap)

	if id > 0 {

		u, err := repo.DB.GetUserById(id)
		if err != nil {
			ClientError(w, r, http.StatusBadRequest)
			return
		}

		vars.Set("user", u)
	} else {
		var u models.User
		vars.Set("user", u)
	}

	err = helpers.RenderPage(w, r, "user", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}

// PostOneUser adds/edits a user
func (repo *DBRepo) PostOneUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error(err)
	}

	var u models.User

	if id > 0 {
		u, _ = repo.DB.GetUserById(id)
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		err := repo.DB.UpdateUser(u)
		if err != nil {
			log.Error(err)
			ClientError(w, r, http.StatusBadRequest)
			return
		}

		if len(r.Form.Get("password")) > 0 {
			// changing password
			err := repo.DB.UpdatePassword(id, r.Form.Get("password"))
			if err != nil {
				log.Error(err)
				ClientError(w, r, http.StatusBadRequest)
				return
			}
		}
	} else {
		u.FirstName = r.Form.Get("first_name")
		u.LastName = r.Form.Get("last_name")
		u.Email = r.Form.Get("email")
		u.UserActive, _ = strconv.Atoi(r.Form.Get("user_active"))
		u.Password = []byte(r.Form.Get("password"))
		u.AccessLevel = 3

		_, err := repo.DB.InsertUser(u)
		if err != nil {
			log.Error(err)
			ClientError(w, r, http.StatusBadRequest)
			return
		}
	}

	repo.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// DeleteUser soft deletes a user
func (repo *DBRepo) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = repo.DB.DeleteUser(id)
	repo.App.Session.Put(r.Context(), "flash", "User deleted")
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
