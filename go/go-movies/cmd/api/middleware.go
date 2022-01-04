package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

// All http middlewares in go commonly have this signature where it takes
// a next http.Handler and returns a http.Handler

// This is a middleware for enabling CORS - basically to respond to preflight requests from the clients
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                           // Allow all origins to connect to the backend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS") // Allow these methods to send pre-flight requests for cross-origin access to the backend
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") // Allow all origins to resend these headers when connecting to the backend
		next.ServeHTTP(w, r)
	})
}

// This is a middleware for authenticating requests based on a jwt token that is passed down in the headers from the frontend
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a response header that is set to indicate the caching servers not to respond with cached content when this header has a value present
		w.Header().Set("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			// You could set a dummy anonymous user if you want the the protected routes to be still available for some testing or demo purposes
		}

		// Here we are splitting the value of the Authorization header - "Authorization : Bearer <jwt_token>"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.badRequestErrorJSON(w, fmt.Errorf("Invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.badRequestErrorJSON(w, fmt.Errorf("Invalid auth header - not a bearer token"))
			return
		}

		// Get the token and check if the signature is correct
		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.authorizationErrorJSON(w, err)
			return
		}

		if !claims.Valid(time.Now()) {
			app.authorizationErrorJSON(w, fmt.Errorf("The Authorization token has expired"))
			return
		}

		if !claims.AcceptAudience("example.com") {
			app.authorizationErrorJSON(w, fmt.Errorf("The Authorization token has an invalid audience"))
			return
		}

		if claims.Issuer != "example.com" {
			app.authorizationErrorJSON(w, fmt.Errorf("The Authorization token has an invalid issuer"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.authorizationErrorJSON(w, fmt.Errorf("The Authorization token has an invalid user id"))
			return
		}

		log.Println("Logged in user id is : ", userID)

		next.ServeHTTP(w, r)
	})
}
