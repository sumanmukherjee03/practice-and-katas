package models

import (
	"context"
	"fmt"
	"time"
)

// GetMovieByID is the func to get a single movie
func (m *DBModel) GetMovieByID(id int) (*Movie, error) {
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

	movieGenres, err := m.getMovieGenresByMovieID(ctx, id)
	if err != nil {
		return nil, err
	}
	movie.MovieGenres = movieGenres

	return &movie, nil
}

// GetAllMovies is the func to get all movies.
// We made an improvement to the initial version of this function to make it variadic
// This helps us in finding movies for a genre if a genre is given or no genre if no genre is given
func (m *DBModel) GetAllMovies(genre ...int) ([]*Movie, error) {
	var movies []*Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at
    from movies %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return movies, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
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
			return movies, err
		}

		movieGenres, err := m.getMovieGenresByMovieID(ctx, movie.ID)
		if err != nil {
			return movies, err
		}
		movie.MovieGenres = movieGenres

		movies = append(movies, &movie)
	}

	return movies, nil
}

// GetMovieGenresByMovieID is the func to get all movie-genres by a given movie id
func (m *DBModel) getMovieGenresByMovieID(ctx context.Context, movieID int) (map[int]string, error) {
	movieGenres := make(map[int]string)

	query := `select mg.id, mg.movie_id, mg.genre_id, mg.created_at, mg.updated_at,
    g.id, g.genre_name, g.created_at, g.updated_at
    from movies_genres mg
    left join genres g on (g.id = mg.genre_id)
    where mg.movie_id = $1`
	rows, err := m.DB.QueryContext(ctx, query, movieID)
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
			&mg.Genre.ID,
			&mg.Genre.GenreName,
			&mg.Genre.CreatedAt,
			&mg.Genre.UpdatedAt,
		)
		if mgErr != nil {
			return movieGenres, err
		}
		movieGenres[mg.ID] = mg.Genre.GenreName
	}

	return movieGenres, nil
}
