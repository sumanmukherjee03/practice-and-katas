package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

////////////////////////////////////// HELPER FUNCS ////////////////////////////////////////

// The expires attribute is only sent with the Set-Cookie response header, not with the Cookie request header.
// ie, the expiry is only set on the response cookie and this attribute is not present on the request cookie.
// Our server sends a response cookie to google oauth when redirecting to google's login page.
// Google sends the cookie back to us as a request cookie when redirecting us to the receive endpoint.
// The Cookie request header when recieved back on the /oauth2/receive endpoint contains only the names and values of the cookies.
// It does not contain any other metadata like expiry.
// So, to check for expiry on the receive endpoint simply check for presence of the cookie.
// Had the cookie expired google wouldnt have sent it back to us on the request header.
func genStateInCookie(w http.ResponseWriter) (string, error) {
	// Here's an example of reading some random bytes and base64 encoding it to send across the wire for the state variable :
	// bytesBuffer := make([]byte, 16)
	// rand.Read(bytesBuffer)
	// state := base64.URLEncoding.EncodeToString(bytesBuffer)
	expiry := time.Now().UTC().Add(10 * time.Minute)
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("ERROR - Could not generate a uuid to be used for representing the login page session for preventing CSRF attacks - %v", err)
	}
	state := id.String()

	// If you want to set a session cookie, dont add the Expires. Session cookies are deleted when the session ends.
	// That is determined by the browser. Where as cookies with Expires are permanent cookies and are deleted at the specified date+time.
	// A cookie with the HttpOnly attribute is inaccessible to the JavaScript Document.cookie API
	// IRL we would also set the cookie to be `Secure: true` so that it is only used in HTTPS, but not for this example since we are using http and localhost as callback.
	// For setting domains and paths on cookies, this discussion on StackOverflow is very relevant
	//   - https://stackoverflow.com/questions/1062963/how-do-browser-cookie-domains-work
	// When setting a cookie remember that the Path optional parameter defaults to whatever path the request is being made to when setting the cookie.
	// ie the default path in this case is going to be "/oauth2/google". However, we'd want the cookie to be available
	// on all the paths. Otherwise, we would not get the cookie in other requests made from the browser.
	// Hence, the Path optional paramater is set to "/".
	cookie := &http.Cookie{
		Name:     googleOAuthStateCookieName,
		Value:    state,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,
	}
	// This essentially does the same thing as - w.Header().Add("Set-Cookie", cookie.String())
	http.SetCookie(w, cookie)
	return state, nil
}

// The purpose of this function is to generate a session id and store the session id to user id mapping in the database
// Simultaneously, also generating a signed jwt token with the claim as the session id.
// And subsequently setting a response cookie in the client browser with the signed jwt toen containing the session id as the value.
// That way for subsequent requests the cookie is sent back to the server in the http request and the server can validate the jwt token,
// extract the session id from the cookie and find the user id from the session id from the DB.
func createSession(userID string, w http.ResponseWriter) error {
	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("ERROR - Could not generate a uuid to be used for representing the logged in session of the user - %v", err)
	}
	sessionID := id.String()
	// The sessions map being used here is a representation of how you would ideally be storing session info in a DB.
	sessions[sessionID] = userID
	token, err := createToken(sessionID)
	if err != nil {
		return fmt.Errorf("ERROR - Could not generate a token to be used for representing the logged in session of the user - %v", err)
	}
	// When setting a cookie remember that the Path optional parameter defaults to whatever path the request is being made to when setting the cookie.
	// ie the default path in this case is going to be "/oauth2/receive". However, we'd want the cookie to be available
	// on all the paths. Otherwise, we would not get the cookie in other requests made from the browser.
	// Hence, the Path optional paramater is set to "/".
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	// This essentially does the same thing as - w.Header().Add("Set-Cookie", cookie.String())
	http.SetCookie(w, cookie)
	return nil
}

// Generate a hashed string from a password which is in plain text.
// This is usually used for basic auth.
// Hashing is a one way function and you cant regenerate password from the hash.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ERROR - Could not generate hashed string from password - %v", err)
	}
	return string(bytes), nil
}

// Compare a hashed password and a user provided password in plain text to validate a login.
// This mechanism is used in basic auth. Hashing is one way. So, the only way to compare is by comparing the hashes.
func comparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("ERROR - Passwords did not match - %v", err)
	}
	return nil
}
