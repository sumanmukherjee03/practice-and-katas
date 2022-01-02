package models

import (
	"context"
	"fmt"
	"time"
)

// InsertMovie is the func to create a single movie
func (m *DBModel) InsertMovie(movie Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into movies(title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	row := m.DB.QueryRowContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Rating,
		movie.Runtime,
		movie.MPAARating,
		time.Now(),
		time.Now(),
	)

	var newID int
	if err := row.Scan(&newID); err != nil {
		return newID, err
	}

	return newID, nil
}

// UpdateMovie is the func to create a single movie
func (m *DBModel) UpdateMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `update movies set
    title = $1,
    description = $2,
    year = $3,
    release_date = $4,
    rating = $5,
    runtime = $6,
    mpaa_rating = $7,
    updated_at = $8 where id = $9`
	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Rating,
		movie.Runtime,
		movie.MPAARating,
		time.Now(),
		movie.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteMovie deletes a movie from the database
func (m *DBModel) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM movies WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

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
