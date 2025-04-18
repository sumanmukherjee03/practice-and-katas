package repository

import "github.com/tsawler/vigilate/internal/models"

// DatabaseRepo is the database repository
type DatabaseRepo interface {
	// preferences
	AllPreferences() ([]models.Preference, error)
	SetSystemPref(name, value string) error
	InsertOrUpdateSitePreferences(pm map[string]string) error
	UpdatePreference(string, string) error

	// users and authentication
	GetUserById(id int) (models.User, error)
	InsertUser(u models.User) (int, error)
	UpdateUser(u models.User) error
	DeleteUser(id int) error
	UpdatePassword(id int, newPassword string) error
	Authenticate(email, testPassword string) (int, string, error)
	AllUsers() ([]*models.User, error)
	InsertRememberMeToken(id int, token string) error
	DeleteToken(token string) error
	CheckForToken(id int, token string) bool

	// hosts
	AllHosts() ([]*models.Host, error)
	GetHostById(int) (models.Host, error)
	InsertHost(models.Host) (int, error)
	UpdateHost(models.Host) error
	DeleteHost(id int) error

	// services
	AllServices() ([]*models.Service, error)
	GetServiceById(int) (models.Service, error)

	// host-services
	GetHostServiceById(int) (models.HostService, error)
	GetHostServiceByHostAndService(int, int) (models.HostService, error)
	InsertHostService(models.HostService) (int, error)
	UpdateHostService(models.HostService) error
	GetAllHostServiceStatusCount() (int, int, int, int, error)
	GetAllHostServicesWithStatus(string) ([]*models.HostService, error)
	GetAllHostServicesToMonitor() ([]*models.HostService, error)

	// events
	InsertEvent(models.Event) (int, error)
	GetAllEvents() ([]*models.Event, error)
}
