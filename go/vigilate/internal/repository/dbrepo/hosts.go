package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

// AllHosts returns all hosts
func (m *postgresDBRepo) AllHosts() ([]*models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at FROM hosts`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []*models.Host

	for rows.Next() {
		h := &models.Host{}
		err = rows.Scan(&h.ID, &h.HostName, &h.CanonicalName, &h.URL, &h.IP, &h.IPV6, &h.Location, &h.OS, &h.Active, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, err
		}
		hostServices, err := m.getAllHostServicesForHost(ctx, h.ID)
		if err != nil {
			return nil, err
		}
		h.HostServices = hostServices
		hosts = append(hosts, h)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return hosts, nil
}

// GetHostById returns a host by id
func (m *postgresDBRepo) GetHostById(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active,
			created_at, updated_at
			FROM hosts where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var h models.Host

	err := row.Scan(
		&h.ID,
		&h.HostName,
		&h.CanonicalName,
		&h.URL,
		&h.IP,
		&h.IPV6,
		&h.Location,
		&h.OS,
		&h.Active,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return h, err
	}

	hostServices, err := m.getAllHostServicesForHost(ctx, h.ID)
	if err != nil {
		log.Println(err)
		return h, err
	}
	h.HostServices = hostServices

	return h, nil
}

// Insert method to add a new record to the hosts table.
func (m *postgresDBRepo) InsertHost(h models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id `

	row := m.DB.QueryRowContext(ctx, stmt,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		&h.Active,
		time.Now(),
		time.Now(),
	)

	var newId int
	err := row.Scan(&newId)
	if err != nil {
		return 0, err
	}

	services, err := m.AllServices()
	if err != nil {
		return newId, fmt.Errorf("Encountered error in fetching all services - %v", err)
	}

	for _, service := range services {
		hostServicesStmt := `INSERT INTO host_services (host_id, service_id, active, schedule_number, schedule_unit, status, created_at, updated_at)
	    VALUES ($1, $2, 0, 3, 'm', 'pending', $3, $4)`
		_, err = m.DB.ExecContext(ctx, hostServicesStmt, newId, service.ID, time.Now(), time.Now())
		if err != nil {
			return newId, fmt.Errorf("Encountered error in associating service with id %d with host - %v", service.ID, err)
		}
	}

	return newId, nil
}

// UpdateHost updates a host by id
func (m *postgresDBRepo) UpdateHost(h models.Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE hosts SET host_name = $1, canonical_name = $2, url = $3, ip = $4, ipv6 = $5, location = $6, os = $7, active = $8, updated_at = $9 WHERE id = $10`

	_, err := m.DB.ExecContext(ctx, stmt,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		h.ID,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteHost deletes a host
func (m *postgresDBRepo) DeleteHost(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM hosts WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ************************** HELPER FNS ***************************

func (m *postgresDBRepo) getAllHostServicesForHost(ctx context.Context, hostID int) ([]models.HostService, error) {
	var hostServices []models.HostService

	stmt := `SELECT hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.last_check, hs.status, hs.created_at, hs.updated_at,
    s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
    FROM host_services hs
    LEFT JOIN services s ON (s.id = hs.service_id)
    WHERE host_id = $1
    ORDER BY s.service_name`

	rows, err := m.DB.QueryContext(ctx, stmt, hostID)
	if err != nil {
		log.Println(err)
		return hostServices, fmt.Errorf("Encountered error in populating HostServices for host - %d : %v", hostID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var hs models.HostService
		err := rows.Scan(
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
		)
		if err != nil {
			log.Println(err)
			return hostServices, fmt.Errorf("Encountered error in populating HostServices for host - %d, host_service - %d : %v", hostID, hs.ID, err)
		}
		hostServices = append(hostServices, hs)
	}

	return hostServices, nil
}
