package models

import (
	"database/sql"
	"time"
)

// DBModel is the wrapper type for sql.DB
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

// Movie is the type for movie
type Movie struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Year        int            `json:"year"`
	ReleaseDate time.Time      `json:"release_date"`
	Runtime     int            `json:"runtime"`
	Rating      int            `json:"rating"`
	MPAARating  string         `json:"mpaa_rating"`
	CreatedAt   time.Time      `json:"created_at"` // Use `json:"-"` to ignore this field in json
	UpdatedAt   time.Time      `json:"updated_at"` // Use `json:"-"` to ignore this field in json
	Poster      string         `json:"string"`
	MovieGenres map[int]string `json:"movie_genres"`
}

// Genre is the type for genre
type Genre struct {
	ID          int           `json:"id"`
	GenreName   string        `json:"genre_name"`
	CreatedAt   time.Time     `json:"-"`
	UpdatedAt   time.Time     `json:"-"`
	MovieGenres []*MovieGenre `json:"movie_genres"`
}

// MovieGenre is the type for movie-genre
type MovieGenre struct {
	ID        int       `json:"-"`
	MovieID   int       `json:"-"`
	GenreID   int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Genre     Genre     `json:"genre"`
	Movie     Movie     `json:"movie"`
}

// User is the type for an user of the website
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"` // The hashed password value
}
