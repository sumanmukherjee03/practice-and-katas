package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

// InsertEvent creates a new event and retuns it's id or an error
func (m *postgresDBRepo) InsertEvent(ev models.Event) (int, error) {
	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO events (event_type, host_id, service_id, host_service_id, host_name, service_name, message, created_at, updated_at)
	  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	row := m.DB.QueryRowContext(ctx, stmt, ev.EventType, ev.HostID, ev.ServiceID, ev.HostServiceID, ev.HostName, ev.ServiceName, ev.Message, time.Now(), time.Now())
	err := row.Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// GetAllEvents returns a list of all events that occurred
func (m *postgresDBRepo) GetAllEvents() ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var events []*models.Event
	stmt := `SELECT id, event_type, host_id, service_id, host_service_id, host_name, service_name, message, created_at, updated_at FROM events ev ORDER BY created_at`
	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		log.Println(err)
		return events, fmt.Errorf("Encountered error in fetching all the events - %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ev models.Event
		err := rows.Scan(
			&ev.ID,
			&ev.EventType,
			&ev.HostID,
			&ev.ServiceID,
			&ev.HostServiceID,
			&ev.HostName,
			&ev.ServiceName,
			&ev.Message,
			&ev.CreatedAt,
			&ev.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return events, fmt.Errorf("Encountered error in fetching event - %v", err)
		}
		events = append(events, &ev)
	}

	return events, nil
}
