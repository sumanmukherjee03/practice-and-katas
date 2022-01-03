package main

import (
	"encoding/json"
	"fmt"
	"go-movies/models"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

// This is a dummy user that would ideally reside in the DB. But we are avoiding doing that
// for development purposes and just hardcoding a single user here so that we can use this for experimenting.
// This is the URL for the go-playground link to generate hashed password : https://go.dev/play/p/uKMMCzJWGsW
// The password for which the hash is pasted in the user is : "password"
var validUser = models.User{
	ID:        11,
	Email:     "john.doe@example.com",
	FirstName: "John",
	LastName:  "Doe",
	Password:  "$2a$12$p09srt6y2w.uXuqpsa39yeHL1esS4ntPTUFU0RinmXnaMHSCZmEmi",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.badRequestErrorJSON(w, err)
		return
	}

	hashedPassword := validUser.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.authorizationErrorJSON(w, err)
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "example.com"
	claims.Audiences = []string{"example.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorJSON(w, err)
		return
	}
	app.writeJSON(w, http.StatusOK, jwtBytes, "response")
}
