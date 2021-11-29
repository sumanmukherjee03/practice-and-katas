package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	// ErrNoRecord no record found in database error
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials invalid username/password error
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail duplicate email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrInactiveAccount inactive account error
	ErrInactiveAccount = errors.New("models: Inactive Account")
)

// User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	UserActive  int
	AccessLevel int
	Email       string
	Password    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Preferences map[string]string
}

// Preference model
type Preference struct {
	ID         int
	Name       string
	Preference []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Host model
type Host struct {
	ID            int
	HostName      string
	CanonicalName string
	URL           string
	IP            string
	IPV6          string
	Location      string
	OS            string
	Active        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	HostServices  []HostService
}

func (h *Host) IsValid() error {
	if len(h.HostName) == 0 {
		return fmt.Errorf("HostName is empty")
	}
	if len(h.CanonicalName) == 0 {
		return fmt.Errorf("CanonicalName is empty")
	}
	if len(h.URL) == 0 {
		return fmt.Errorf("URL is empty")
	}
	if h.Active > 1 {
		return fmt.Errorf("Active can only be 0 or 1")
	}
	return nil
}

// Service model
type Service struct {
	ID          int
	ServiceName string
	Icon        string
	Active      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// HostService model
type HostService struct {
	ID             int
	HostID         int
	ServiceID      int
	ScheduleNumber int
	ScheduleUnit   string
	Status         string
	LastCheck      time.Time
	Active         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Service        Service
	Host           Host
}

func (hs *HostService) ScheduleText() (string, error) {
	switch hs.ScheduleUnit {
	case "d":
		return fmt.Sprintf("@every %d%s", hs.ScheduleNumber*24, "h"), nil
	case "h":
		return fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit), nil
	case "m":
		return fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit), nil
	case "s":
		return fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit), nil
	default:
		return "", fmt.Errorf("Invalid schedule unit - %s", hs.ScheduleUnit)
	}
}

// Schedule model
type Schedule struct {
	ID                     int
	EntryID                cron.EntryID
	Entry                  cron.Entry
	Host                   string
	Service                string
	LastRunFromHostService time.Time
	HostServiceID          int
	ScheduleText           string
}

// Event model
type Event struct {
	ID            int
	EventType     string
	HostID        int
	ServiceID     int
	HostServiceID int
	HostName      string
	ServiceName   string
	Message       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Host          Host
	Service       Service
	HostService   HostService
}
