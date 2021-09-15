package main

import (
	"fmt"
	"net/http"

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

var github

var githubOAuthConfig = &oauth2.Config{
	ClientID:     "052b8521298530d2d897",
	ClientSecret: "186975de0888546e28f8f5f0d9a573aff582ed13",
	Endpoint:     github.Endpoint,
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandlerFunc("/oauth/github", githubOAuthHandler)
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
	`)
	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func githubOAuthHandler(w http.ResponseWriter, r *http.Request) {
	// The state variable being passed to AuthCodeURL is generally a uuid representing a login attempt.
	// It usually is an id maintained in a DB. The id represents a login attempt and has an expiration time associated with it,
	// Usually the login attempt is not valid after that expiration time.
	redirectURL := githubOAuthConfig.AuthCodeURL("0000")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
