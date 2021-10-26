package dbrepo

import (
	"context"
	"fmt"
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

// GetHostServiceStatusCount returns the active hosts with services that have status pending, healthy, warning and problem
func (m *postgresDBRepo) GetAllHostServiceStatusCount() (int, int, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT
		(SELECT COUNT(id) FROM host_services WHERE active = 1 AND status = 'pending') AS pending,
		(SELECT COUNT(id) FROM host_services WHERE active = 1 AND status = 'healthy') AS healthy,
		(SELECT COUNT(id) FROM host_services WHERE active = 1 AND status = 'warning') AS warning,
		(SELECT COUNT(id) FROM host_services WHERE active = 1 AND status = 'problem') AS problem`
	row := m.DB.QueryRowContext(ctx, stmt)

	var pending int
	var healthy int
	var warning int
	var problem int

	err := row.Scan(
		&pending,
		&healthy,
		&warning,
		&problem,
	)

	if err != nil {
		log.Println(err)
		return 0, 0, 0, 0, err
	}

	return pending, healthy, warning, problem, nil
}

// GetAllHostServicesWithStatus returns the host services that have status pending
func (m *postgresDBRepo) GetAllHostServicesWithStatus(status string) ([]*models.HostService, error) {
	var hss []*models.HostService
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
    s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at,
    h.id, h.host_name, h.canonical_name, h.url, h.ip, h.ipv6, h.location, h.os, h.active, h.created_at, h.updated_at
    FROM host_services hs
    LEFT JOIN services s ON (s.id = hs.service_id)
    LEFT JOIN hosts h ON (h.id = hs.host_id)
    WHERE hs.status = $1`
	rows, err := m.DB.QueryContext(ctx, stmt, status)
	if err != nil {
		return hss, fmt.Errorf("ERROR - Could not fetch rows for host services : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		hs := &models.HostService{}
		err = rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.LastCheck,
			&hs.Status,
			&hs.CreatedAt,
			&hs.UpdatedAt,
			&hs.Service.ID,
			&hs.Service.ServiceName,
			&hs.Service.Active,
			&hs.Service.Icon,
			&hs.Service.CreatedAt,
			&hs.Service.UpdatedAt,
			&hs.Host.ID,
			&hs.Host.HostName,
			&hs.Host.CanonicalName,
			&hs.Host.URL,
			&hs.Host.IP,
			&hs.Host.IPV6,
			&hs.Host.Location,
			&hs.Host.OS,
			&hs.Host.Active,
			&hs.Host.CreatedAt,
			&hs.Host.UpdatedAt,
		)
		if err != nil {
			return hss, err
		}
		hss = append(hss, hs)
	}

	return hss, nil
}
