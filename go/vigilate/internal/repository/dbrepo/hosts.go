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
		// Append it to the slice
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
	    VALUES ($1, $2, 1, 3, 'm', 'pending', $3, $4)`
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
