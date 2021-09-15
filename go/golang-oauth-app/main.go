package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
)

var githubOAuthConfig *oauth2.Config

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
    <form action="/oauth/github" method="post" accept-charset="utf-8">
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
	redirectURL := githubOAuthConfig.AuthCodeURL("0000")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// After login and granting permissions the github oauth login page will redirect us to
// http://localhost:8001/oauth2/receive?code=<a_unique_token>&state=0000
// Note here that the state containing the uuid for the session is passed back to us as a query param in the callback url
// This callback url was set in the github oauth application setup
func githubOAuthReceiveHandler(w http.ResponseWriter, r *http.Request) {
}
