package models

import (
	"context"
	"time"
)

// GetGenreByID is the func to get a single genre by name
func (m *DBModel) GetGenreByID(id int) (*Genre, error) {
	var genre Genre
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `select id, genre_name, created_at, updated_at
    from genres
    where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&genre.ID,
		&genre.GenreName,
		&genre.CreatedAt,
		&genre.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	genreMovies, err := m.getMovieGenresByGenreID(ctx, genre.ID)
	if err != nil {
		return nil, err
	}
	genre.MovieGenres = genreMovies

	return &genre, nil
}

// GetAllGenres is the func to get all genres
func (m *DBModel) GetAllGenres() ([]*Genre, error) {
	var genres []*Genre

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, genre_name, created_at, updated_at
    from genres
    order by genre_name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return genres, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre Genre

		err := rows.Scan(
			&genre.ID,
			&genre.GenreName,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return genres, err
		}

		movieGenres, err := m.getMovieGenresByGenreID(ctx, genre.ID)
		if err != nil {
			return genres, err
		}
		genre.MovieGenres = movieGenres

		genres = append(genres, &genre)
	}

	return genres, nil
}

// GetMovieGenresByGenreID is the func to get all movie-genres by a given genre id
func (m *DBModel) getMovieGenresByGenreID(ctx context.Context, genreID int) ([]*MovieGenre, error) {
	var movieGenres []*MovieGenre

	query := `select mg.id, mg.movie_id, mg.genre_id, mg.created_at, mg.updated_at,
    m.id, m.title, m.description, m.year, m.release_date, m.rating, m.runtime, m.mpaa_rating, m.created_at, m.updated_at
    from movies_genres mg
    left join movies m on (m.id = mg.movie_id)
    where mg.genre_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, genreID)
	if err != nil {
		return movieGenres, err
	}
	defer rows.Close()

	for rows.Next() {
		var mg MovieGenre
		mgErr := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.CreatedAt,
			&mg.UpdatedAt,
			&mg.Movie.ID,
			&mg.Movie.Title,
			&mg.Movie.Description,
			&mg.Movie.Year,
			&mg.Movie.ReleaseDate,
			&mg.Movie.Rating,
			&mg.Movie.Runtime,
			&mg.Movie.MPAARating,
			&mg.Movie.CreatedAt,
			&mg.Movie.UpdatedAt,
		)
		if mgErr != nil {
			return movieGenres, err
		}
		movieGenres = append(movieGenres, &mg)
	}

	return movieGenres, nil
}
