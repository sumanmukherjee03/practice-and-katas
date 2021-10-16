package dbrepo

import (
	"context"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

// Insert method to add a new record to the hosts table.
func (m *postgresDBRepo) InsertHostService(h models.Host, s models.Service, activate int) (int, error) {
	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO host_services (host_id, service_id, active, schedule_number, schedule_unit, status, created_at, updated_at)
	  VALUES ($1, $2, $3, 3, 'm', 'pending', $4, $5) returning id`
	row := m.DB.QueryRowContext(ctx, stmt, h.ID, s.ID, activate, time.Now(), time.Now())
	err := row.Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}
