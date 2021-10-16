package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/tsawler/vigilate/internal/models"
)

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
