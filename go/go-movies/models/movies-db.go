package models

import (
	"context"
	"database/sql"
	"time"
)

// Movie is the wrapper type for sql.DB
type DBModel struct {
	DB *sql.DB
}

// Get is the func to get a single movie
func (m *DBModel) Get(id int) (*Movie, error) {
	var movie Movie
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at
    from movies
    where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `select mg.id, mg.movie_id, mg.genre_id, mg.created_at, mg.updated_at,
    g.id, g.genre_name, g.created_at, g.updated_at
    from movies_genres mg
    left join genres g on (g.id = mg.genre_id)
    where mg.movie_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieGenres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		mgErr := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.CreatedAt,
			&mg.UpdatedAt,
			&mg.Genre.ID,
			&mg.Genre.GenreName,
			&mg.Genre.CreatedAt,
			&mg.Genre.UpdatedAt,
		)
		if mgErr != nil {
			return nil, err
		}
		movieGenres[mg.ID] = mg.Genre.GenreName
	}
	movie.MovieGenres = movieGenres

	return &movie, nil
}

// All is the func to get all movies
func (m *DBModel) All(id int) ([]*Movie, error) {
	return nil, nil
}
