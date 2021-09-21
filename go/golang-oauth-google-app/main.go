package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Regarding how to setup a google oauth2 app follow the instructions here : https://support.google.com/cloud/answer/6158849

const (
	googleOAuthAppClientIDEnvVar     = "GOOGLE_OAUTH2_APP_CLIENT_ID"
	googleOAuthAppClientSecretEnvVar = "GOOGLE_OAUTH2_APP_CLIENT_SECRET"
	googleEmailScope                 = "https://www.googleapis.com/auth/userinfo.email"
	googleAPIEndpoint                = "https://www.googleapis.com/oauth2/v2/userinfo"
	googleOAuthStateCookieName       = "googleOAuthState"
	jwtSigningSecretKeyEnvVar        = "JWT_SECRET"
	sessionCookieName                = "session"
	sessionExpiry                    = 5 * time.Minute
)

var (
	googleOAuthConfig   *oauth2.Config
	googleOAuthScopes   []string
	emailToDBUserMap    = make(map[string]User)   // Key is email id and value is a User type
	sessions            = make(map[string]string) // Key is a session ID and value is a user ID
	jwtSigningSecretKey []byte
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	GoogleID  string `json:"google_id"`
}

type googleResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func init() {
	googleOAuthAppClientID, ok := os.LookupEnv(googleOAuthAppClientIDEnvVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "Google oauth2 app client id not provided in env var %s", googleOAuthAppClientIDEnvVar)
		os.Exit(2)
	}
	googleOAuthAppClientSecret, ok := os.LookupEnv(googleOAuthAppClientSecretEnvVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "Google oauth2 app client secret not provided in env var %s", googleOAuthAppClientSecretEnvVar)
		os.Exit(2)
	}
	jwtSigningSecretKeyStr, ok := os.LookupEnv(jwtSigningSecretKeyEnvVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "Jwt signing secret not provided in env var %s", jwtSigningSecretKeyEnvVar)
		os.Exit(2)
	}
	if len(jwtSigningSecretKeyStr) != 32 {
		fmt.Fprintf(os.Stderr, "Jwt signing secret should be a alphanumeric string of 32 characters")
		os.Exit(2)
	}
	jwtSigningSecretKey = []byte(jwtSigningSecretKeyStr)
	// The config for google oauth2 client setup requires a redirect url to be present.
	// This is the same callback url that you put into the google oauth2 app setup.
	// The scopes mentioned here are also used during the setting up of the authorization console for the oauth2 app.
	// The scopes are basically permissions that would be granted by google during the authorization step of the login process.
	// These scopes are also setup while setting up the oauth2 app authorization console.
	googleOAuthScopes = []string{
		googleEmailScope,
	}
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8002/oauth2/receive",
		ClientID:     googleOAuthAppClientID,
		ClientSecret: googleOAuthAppClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       googleOAuthScopes,
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/oauth2/google", googleOAuthLoginHandler)
	http.HandleFunc("/oauth2/receive", googleOAuthReceiveHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/register/submit", registerSubmitHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.ListenAndServe(":8002", nil)
}

///////////////////////////////////// HTTP HANDLER FUNCS /////////////////////////////////////

func indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	html := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>index</title>
  </head>
  <body>`
	form := `
    <form action="/oauth2/google" method="post" accept-charset="utf-8">
      <input type="submit" value="Login with Google" name="login_with_google" id="login_with_google"/>
    </form>`
	footer := `
    <form action="/logout" method="post" accept-charset="utf-8">
      <input type="submit" value="Logout" name="logout" id="logout"/>
    </form>
  </body>
</html>`

	if len(msg) > 0 {
		html += `<p>Notice : ` + msg + `</p>`
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		cookie = &http.Cookie{
			Name:  sessionCookieName,
			Value: "",
		}
	}

	sessionID, err := parseToken(cookie.Value)
	if err != nil {
		log.Info(fmt.Sprintf("Could not get a valid session id from the session cookie - %v", err))
		html += form
	} else {
		userID, ok := sessions[sessionID]
		if !ok {
			html += form
		} else {
			html += `<p>Logged in user id : ` + userID + `</p>`
		}
	}
	html += footer

	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func googleOAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	// The state variable being passed to AuthCodeURL is generally a uuid representing a login attempt.
	// It usually is an id maintained in a DB or a cookie on the client side.
	// This unique id is usually associated with an expiry time.
	// The login attempt is not valid after that expiration time.
	// When the google login redirects you back to the server, the server should match this id and the expiry time.
	// This is necessary to protect the client from CSRF attacks.
	// ie, this expiry essentially represents how to set a hard boundary on how long it takes for the oauth provider
	// to redirect users back to our website during a login attempt

	// Validate that this is a POST request
	if r.Method != http.MethodPost {
		log.Error("ERROR - This path only handles a POST request")
		http.Error(w, "This needs to be a POST request to accept the form submission for login with google", http.StatusBadRequest)
		return
	}

	state, err := genStateInCookie(w)
	if err != nil {
		log.Error("ERROR - Could not handle login attempt and redirect to google", err)
		http.Error(w, "Failed to handle login attempt and redirect to google", http.StatusInternalServerError)
		return
	}
	googleLoginRedirectURL := googleOAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, googleLoginRedirectURL, http.StatusSeeOther)
}

// After login and granting permissions the github oauth login page will redirect us to
// http://localhost:8001/oauth2/receive?code=<a_unique_token>&state=0000
// Note here that the state containing the uuid for the session is passed back to us as a query param in the callback url
// This callback url was set in the github oauth application setup
// This redirect url can also be optionally provided in the config for the oauth2 client, but the host and port must match what is configured in the app on github.
// Usually though, it is not necessary to be mentioned in the config and can be skipped.
func googleOAuthReceiveHandler(w http.ResponseWriter, r *http.Request) {
	// If you want to know what's in the query params, then use ParseForm and marshal the data so that you can print it
	// if err := r.ParseForm(); err != nil {
	// log.Error(err)
	// }
	// data, err := json.Marshal(r.Form)
	// if err != nil {
	// log.Error(err)
	// }
	// fmt.Println(string(data))

	code := r.FormValue("code")
	state := r.FormValue("state")
	scope := r.FormValue("scope")
	cookieState, err := r.Cookie(googleOAuthStateCookieName)
	if err != nil {
		log.Error("ERROR - The state cookie could not be found in the request which indicates that the cookie must have expired")
		msg := url.QueryEscape("Cookie to protect against CSRF attack not found or must have expired")
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	// Check if the session is a valid one or not. This is required to protect against CSRF attacks
	if state != cookieState.Value {
		log.Error("ERROR - The state returned back from the google oauth login doesnt match what is in the cookie")
		msg := "Login with google either has an invalid state"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	// Validate that the scope of the login contains the scope we need to query the google API
	for _, s := range googleOAuthScopes {
		if !strings.Contains(scope, s) {
			log.Error("ERROR - The scope returned back from the google oauth login does not contain the required scopes")
			msg := "Login with google has incorrect scope"
			http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
			return
		}
	}

	// At this point use the code from the query params to exchange it for a auth token with google.
	// This also uses the client id and secret because it is in the oauth config.
	// And the url that is called is the TokenURL from the oauth endpoints config.
	// For the exchange http call use the same context that was provided with the request, so that if the request times out
	// or is cancelled the token exchange call is also cancelled
	// The token that you get back here is a Bearer token that comes with it's own expiry.
	// This token is not necessarily always jwt, although in some cases it can be. For google's oauth it is not a jwt token.
	token, err := googleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Error("ERROR - Could not exchange oauth login code with google for an auth token", err)
		msg := "Login failed because we could not exchange code for oauth token"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	// TokenSource is an interface, basically anything that can respond to a Token() method to return a token
	// The difference between Token concrete type and this interface is that if your token is expiring,
	// calling the Token() method in the the TokenSource interface will essentially give you back a new token based on a refresh token if necessary.
	tokenSource := googleOAuthConfig.TokenSource(r.Context(), token)
	// Use this token source to get an authenticated google http client
	// You can use this http client to make calls to the google api
	client := oauth2.NewClient(r.Context(), tokenSource)
	resp, err := client.Get(googleAPIEndpoint)
	if err != nil {
		log.Error("ERROR - Could not fetch user info using the access token provided after exchange", err)
		msg := "Failed to fetch user details on login"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Error("ERROR - Google returned an unsuccessful response when fetching user details")
		msg := "Failed to fetch user details on login"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	// Use this code chunk if you want to print out the json response body so that you can create the struct to unmarshal it
	// bytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// log.Error("ERROR - Could not read the response from google while fetching user info", err)
	// http.Error(w, "Login failed", http.StatusInternalServerError)
	// return
	// }
	// fmt.Println(string(bytes))

	var gr googleResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Error("ERROR - Could not unmarshal user info fetched from google", err)
		msg := "Login failed because we couldnt unmarshal user details from google"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	// If user already exists in our DB
	//   - find the user from the DB
	//     create a session id
	//     create a token with the session id in it's claims
	//     persist the session id in DB to extract user id from session id for subsequent requests
	//     stick the signed token containing the session id in it's claims in a cookie to set it in the client side
	// If user does not exist in our DB
	//   - find the user from the DB
	//     create a token with the email instead of a session id
	//     get the name, email and other details that you could fetch from the google response
	//     and redirect the user to a partial register page so that some of these details can be prefilled
	email := gr.Email
	user, ok := emailToDBUserMap[email]
	if !ok {
		// Redirect the user to a registration page
		email := gr.Email
		vals := url.Values{}
		vals.Add("email", email)
		vals.Add("google_id", gr.ID)
		http.Redirect(w, r, "/register?"+vals.Encode(), http.StatusSeeOther)
		return
	}

	userID := user.ID
	createSessionErr := createSession(userID, w)
	if createSessionErr != nil {
		log.Error("Failed to create session for authenticated user", createSessionErr)
		msg := "Failed to create session for authenticated user"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	msg := fmt.Sprintf("Successfully logged in user %s %s - %s", user.FirstName, user.LastName, user.Email)
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Validate that this is a POST request
	if r.Method != http.MethodPost {
		log.Error("ERROR - This path only handles a POST request")
		http.Error(w, "This needs to be a POST request to accept the form submission for logout", http.StatusBadRequest)
		return
	}

	sessionCookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		sessionCookie = &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
		}
	}

	sessionID, err := parseToken(sessionCookie.Value)
	if err != nil {
		log.Info(fmt.Sprintf("Could not get a valid session id from the session cookie - %v", err))
	}
	delete(sessions, sessionID)
	// Setting the max age to -1 on a cookie tells the browser to expire the cookie
	sessionCookie.MaxAge = -1
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	googleID := r.FormValue("google_id")
	msg := r.FormValue("msg")
	html := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>index</title>
  </head>
  <body>`
	form := `
    <form action="/register/submit" method="post" accept-charset="utf-8">
      <label for="email">Email:</label>`
	footer := `
      <input type="submit" value="Register" name="register" id="register"/>
    </form>
  </body>
</html>`

	if len(msg) > 0 {
		html += `<p>Notice : ` + msg + `</p>`
	}

	if len(email) > 0 {
		form += `
      <input type="text" name="email" id="email" value="` + email + `" />`
	} else {
		form += `
      <input type="text" name="email" id="email" />`
		form += `
      <label for="password">Password:</label>
      <input type="password" id="password" name="password" minlength="16" required>`
	}

	form += `
      <label for="first_name">First Name:</label>
      <input type="text" name="first_name" id="first_name" />
      <label for="last_name">Last Name:</label>
      <input type="text" name="last_name" id="last_name" />`

	if len(googleID) > 0 {
		form += `
      <input type="hidden" id="google_id" name="google_id" value="` + googleID + `">`
	}

	html += form
	html += footer

	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func registerSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Validate that this is a POST request
	if r.Method != http.MethodPost {
		log.Error("ERROR - This path only handles a POST request")
		http.Error(w, "This needs to be a POST request to accept the form submission for registration", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	password := r.FormValue("password")
	googleID := r.FormValue("google_id")
	if len(password) > 0 {
		log.Error("Currently we havent implemented a registration page without google oauth")
		http.Error(w, "Normal registration remains unimplemented at this point", http.StatusNotImplemented)
		return
	}

	var userID string
	if len(email) > 0 {
		id, err := uuid.NewV4()
		if err != nil {
			log.Error("ERROR - Failed to generate user id", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		userID = id.String()
		emailToDBUserMap[email] = User{
			ID:        userID,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			GoogleID:  googleID,
		}
	} else {
		log.Error("ERROR - An email id is required to be registered")
		http.Error(w, "Missing email id", http.StatusBadRequest)
		return
	}

	createSessionErr := createSession(userID, w)
	if createSessionErr != nil {
		log.Error("Failed to create session for registered user", createSessionErr)
		msg := "Failed to create session for registered user"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	msg := fmt.Sprintf("Successfully registered and logged in user with email %s", email)
	http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)

}

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
