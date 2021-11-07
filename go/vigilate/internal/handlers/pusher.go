package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/pusher/pusher-http-go"
	log "github.com/sirupsen/logrus"
)

// This handler handles an authenticated request.
// Hence things like the user are available to it in the session.
func (repo *DBRepo) PusherAuth(w http.ResponseWriter, r *http.Request) {
	userID := repo.App.Session.GetInt(r.Context(), "userID")

	user, err := repo.DB.GetUserById(userID)
	if err != nil {
		log.Error("ERROR - Could not fetch user with id retrieved from session", err)
		http.Error(w, "Could not fetch user with id retrieved from session", http.StatusInternalServerError)
		return
	}

	params, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("ERROR - Could not read request body", err)
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(userID),
		UserInfo: map[string]string{
			"name": user.FirstName,
			"id":   strconv.Itoa(userID),
		},
	}

	resp, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
	if err != nil {
		log.Error("ERROR - Could not authenticate user with pusher server via websocket client", err)
		http.Error(w, "Could not authenticate user with pusher server via websocket client", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		log.Error("ERROR - Could not write json response", err)
		http.Error(w, "Could not write json response", http.StatusInternalServerError)
		return
	}
}

func (repo *DBRepo) TestPusher(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["message"] = "hello testing"
	// We only have one channel at the moment called "public-channel".
	// It is being subscribed to on the frontend as well. Look in partials/js.jet.
	err := repo.App.WsClient.Trigger("public-channel", "testEvent", data)
	if err != nil {
		log.Error("ERROR - Could not push event to pusher via the websocket client", err)
		return
	}
}
