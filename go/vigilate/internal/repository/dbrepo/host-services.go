package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

// GetHostServiceByHostAndService returns a host by id
func (m *postgresDBRepo) GetHostServiceByHostAndService(hostID, serviceID int) (models.HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, host_id, service_id, active, schedule_number, schedule_unit, status, created_at, updated_at
		FROM host_services WHERE host_id = $1 AND service_id = $2`
	row := m.DB.QueryRowContext(ctx, stmt, hostID, serviceID)

	var hs models.HostService

	err := row.Scan(
		&hs.ID,
		&hs.HostID,
		&hs.ServiceID,
		&hs.Active,
		&hs.ScheduleNumber,
		&hs.ScheduleUnit,
		&hs.Status,
		&hs.CreatedAt,
		&hs.UpdatedAt,
	)

	if err != nil {
		log.Println(err)
		return hs, err
	}

	return hs, nil
}

// Insert method to add a new record to the hosts table.
func (m *postgresDBRepo) InsertHostService(hs models.HostService) (int, error) {
	var newId int
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `INSERT INTO host_services (host_id, service_id, active, schedule_number, schedule_unit, status, created_at, updated_at)
	  VALUES ($1, $2, $3, 3, 'm', 'pending', $4, $5) returning id`
	row := m.DB.QueryRowContext(ctx, stmt, hs.HostID, hs.ServiceID, hs.Active, time.Now(), time.Now())
	err := row.Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// Insert method to add a new record to the hosts table.
func (m *postgresDBRepo) UpdateHostService(hs models.HostService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `UPDATE host_services SET host_id = $1, service_id = $2, active = $3, updated_at = $4 WHERE id = $5`
	_, err := m.DB.ExecContext(ctx, stmt, hs.HostID, hs.ServiceID, hs.Active, time.Now(), hs.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
