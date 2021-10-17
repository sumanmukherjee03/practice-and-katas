package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

// AllServices returns all services
func (m *postgresDBRepo) AllServices() ([]*models.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, service_name, icon, active, created_at, updated_at FROM services`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service

	for rows.Next() {
		s := &models.Service{}
		err = rows.Scan(&s.ID, &s.ServiceName, &s.Icon, &s.Active, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		// Append it to the slice
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return services, nil
}

// GetServiceById returns a service by id
func (m *postgresDBRepo) GetServiceById(id int) (models.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, service_name, icon, active,
			created_at, updated_at
			FROM services where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var s models.Service

	err := row.Scan(
		&s.ID,
		&s.ServiceName,
		&s.Icon,
		&s.Active,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		log.Println(err)
		return s, err
	}

	return s, nil
}
