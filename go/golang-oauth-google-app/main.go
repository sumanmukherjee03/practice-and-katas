package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
)

var (
	googleOAuthConfig *oauth2.Config
	googleOAuthScopes []string

	// Key is google user ID and value is our own user ID
	googleUserToDBUserMap = make(map[string]string)
)

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
	http.ListenAndServe(":8002", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>index</title>
  </head>
  <body>
    <form action="/oauth2/google" method="post" accept-charset="utf-8">
      <input type="submit" value="Login with Google" name="submit" id="submit"/>
    </form>
  </body>
</html>
	`
	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func googleOAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	// The state variable being passed to AuthCodeURL is generally a uuid representing a login attempt.
	// It usually is an id maintained in a DB. The id represents a login attempt and has an expiration time associated with it,
	// Usually the login attempt is not valid after that expiration time.
	googleLoginRedirectURL := googleOAuthConfig.AuthCodeURL("0000")
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

	// If the login session has expired or the session id is not a valid one, return a http error
	if state != "0000" {
		log.Error("ERROR - The state returned back from the google oauth login is incorrect or has expired")
		http.Error(w, "Login with google either expired or is in an invalid state", http.StatusBadRequest)
		return
	}

	// Validate that the scope of the login contains the scope we need to query the google API
	for _, s := range googleOAuthScopes {
		if !strings.Contains(scope, s) {
			log.Error("ERROR - The scope returned back from the google oauth login does not contain the required scopes")
			http.Error(w, "Login with google has incorrect scope", http.StatusBadRequest)
			return
		}
	}

	// At this point use the code from the query params to exchange it for a auth token with google.
	// For the exchange http call use the same context that was provided with the request, so that if the request times out
	// or is cancelled the token exchange call is also cancelled
	token, err := googleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Error("ERROR - Could not exchange oauth login code with google for an auth token", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// TokenSource is an interface, basically anything that can respond to a Token() method to return a token
	tokenSource := googleOAuthConfig.TokenSource(r.Context(), token)
	// Use this token source to get an authenticated google http client
	// You can use this http client to make calls to the google api
	client := oauth2.NewClient(r.Context(), tokenSource)
	resp, err := client.Get(googleAPIEndpoint)
	if err != nil {
		log.Error("ERROR - Could not fetch user info using the access token provided after exchange", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

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
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	googleID := gr.ID
	userID, ok := googleUserToDBUserMap[googleID]
	if !ok {
		// Create a new user account
		id, err := uuid.NewV4()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed to create new user from authenticated google user id", http.StatusInternalServerError)
			return
		}
		userID = id.String()
		googleUserToDBUserMap[googleID] = userID
	}
	fmt.Println("googleUserToDBUserMap : ", googleUserToDBUserMap)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
