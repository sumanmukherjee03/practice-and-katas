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
	"golang.org/x/oauth2/github"
)

// The 3 stage OAuth flow :
// login to mywebsite.com with a form on mywebsite.com
//   -> on hitting submit, the form data is sent to google for authentication by the OAuth2 library
//        - user is asked to grant permissions
//        - especially what to share with mywebsite.com from google
//     -> google authenticates the clients credentials
//       -> google sends back a redirect url (usually referred to as callback url when setting up oauth), a url which points to some location on mywebsite.com
//         -> mywebsite.com redirects the client to mywebsite.com oauth callback url with the query parameters in the url
//           -> mywebsite.com exchanges that token + a secret with google to get back a jwt token, commonly referred to as access token
//                 - this jwt access token usually contains a google user id which maps directly to an internal user id in the database of mywebsite.com
//             -> from that point onwards, mywebsite.com uses the jwt token to ask google if the user is authenticated or not
//                OR mywebsite.com can use the user id to issue it's own jwt token and store that in cookies for logged in information

// There is also a 2 stage OAuth flow which is implicit, because that does not exchange the token from redirect url to get back an access token

// In this example we are using a sample github oauth2 app registered for playing with oauth.

const (
	githubOAuthAppClientIDEnvVar     = "GITHUB_OAUTH2_APP_CLIENT_ID"
	githubOAuthAppClientSecretEnvVar = "GITHUB_OAUTH2_APP_CLIENT_SECRET"
	githubGraphqlApiEndpoint         = "https://api.github.com/graphql"
)

var (
	githubOAuthConfig *oauth2.Config

	// Key is github ID and value is our user ID
	githubUserToDBUserMap = make(map[string]string)
)

// The response body of the github graphql api would be like this - {"data":{"viewer":{"id":"..."}}}
// The type here is to be able to unmarshal that response json into a concrete type
type githubResponse struct {
	Data struct {
		Viewer struct {
			ID string `json:"id"`
		} `json:"viewer"`
	} `json:"data"`
}

func init() {
	githubOAuthAppClientID, ok := os.LookupEnv(githubOAuthAppClientIDEnvVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "Github oauth2 app client id not provided in env var %s", githubOAuthAppClientIDEnvVar)
		os.Exit(2)
	}
	githubOAuthAppClientSecret, ok := os.LookupEnv(githubOAuthAppClientSecretEnvVar)
	if !ok {
		fmt.Fprintf(os.Stderr, "Github oauth2 app client secret not provided in env var %s", githubOAuthAppClientSecretEnvVar)
		os.Exit(2)
	}
	githubOAuthConfig = &oauth2.Config{
		ClientID:     githubOAuthAppClientID,
		ClientSecret: githubOAuthAppClientSecret,
		Endpoint:     github.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/oauth2/github", githubOAuthLoginHandler)
	http.HandleFunc("/oauth2/receive", githubOAuthReceiveHandler)
	http.ListenAndServe(":8001", nil)
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
    <form action="/oauth2/github" method="post" accept-charset="utf-8">
      <input type="submit" value="Login with Github" name="submit" id="submit"/>
    </form>
  </body>
</html>
	`
	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func githubOAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	// The state variable being passed to AuthCodeURL is generally a uuid representing a login attempt.
	// It usually is an id maintained in a DB. The id represents a login attempt and has an expiration time associated with it,
	// Usually the login attempt is not valid after that expiration time.
	githubLoginRedirectURL := githubOAuthConfig.AuthCodeURL("0000")
	http.Redirect(w, r, githubLoginRedirectURL, http.StatusSeeOther)
}

// After login and granting permissions the github oauth login page will redirect us to
// http://localhost:8001/oauth2/receive?code=<a_unique_token>&state=0000
// Note here that the state containing the uuid for the session is passed back to us as a query param in the callback url
// This callback url was set in the github oauth application setup
// This redirect url can also be optionally provided in the config for the oauth2 client, but the host and port must match what is configured in the app on github.
// Usually though, it is not necessary to be mentioned in the config and can be skipped.
func githubOAuthReceiveHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	state := r.FormValue("state")
	// If the login session has expired or the session id is not a valid one, return a http error
	if state != "0000" {
		http.Error(w, "The state returned back from the github oauth login is either incorrect or has expired", http.StatusBadRequest)
		return
	}
	// At this point use the code from the query params to exchange it for a auth token with github.
	// For the exchange http call use the same context that was provided with the request, so that if the request times out
	// or is cancelled the token exchange call is also cancelled
	token, err := githubOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Could not exchange oauth login code with github for an auth token", http.StatusInternalServerError)
		return
	}
	// TokenSource is an interface, basically anything that can respond to a Token() method to return a token
	tokenSource := githubOAuthConfig.TokenSource(r.Context(), token)
	// Now use this token source to get an authenticated github http client
	// You can use this http client to make calls to the github api
	client := oauth2.NewClient(r.Context(), tokenSource)

	// You need to pass a reader to the client.Post method for the request body
	// Also, this request body is in json format, but represents a graphql query
	// Simple curl examples of graphql api queries here : https://docs.github.com/en/graphql/guides/forming-calls-with-graphql
	requestBody := strings.NewReader(`{"query": "query {viewer {id}}"}`)
	resp, err := client.Post(githubGraphqlApiEndpoint, "application/json", requestBody)
	if err != nil {
		log.Error(err)
		http.Error(w, "Could not fetch logged in user details from github using an authenticated client", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// bytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// log.Error(err)
	// http.Error(w, "Could not read response body for request from github using an authenticated client to fetch user id", http.StatusInternalServerError)
	// return
	// }

	var gr githubResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Error(err)
		http.Error(w, "Could not unmarshal response body for request from github using an authenticated client to fetch user id", http.StatusInternalServerError)
		return
	}
	githubID := gr.Data.Viewer.ID
	userID, ok := githubUserToDBUserMap[githubID]
	if !ok {
		// Create a new user account
		id, err := uuid.NewV4()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed to create new user from authenticated github user id", http.StatusInternalServerError)
			return
		}
		userID = id.String()
		githubUserToDBUserMap[githubID] = userID
	}
	fmt.Println("githubUserToDBUserMap : ", githubUserToDBUserMap)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
