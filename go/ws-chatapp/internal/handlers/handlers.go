package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	// Create a set of views from the dirs that are passed so that jet templates can be looked up by name
	views = jet.NewHTMLSet("./html")

	// Create a variable of websocket.Upgrader type to upgrade a normal http connection to a websocket connection
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

/////////////////////////////////////////////////////////////////////
///////////////////////// HANDLER FUNCS /////////////////////////////
/////////////////////////////////////////////////////////////////////
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Error("ERROR - could not fetch template for rendering", err)
		http.Error(w, "Could not render home page", http.StatusInternalServerError)
		return
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Error("ERROR - could not upgrade http connection to websocket connection", err)
		http.Error(w, "Could not upgrade http connection to websocket connection", http.StatusInternalServerError)
		return
	}
	log.Info("Client connected to endpoint")
	var resp WsJsonResponse
	resp.Message = `<em><small>Connected to server</small></em>`
	err = ws.WriteJSON(resp)
	if err != nil {
		log.Error("ERROR - could not send back json response to client", err)
		http.Error(w, "Could not return a json response for upgraded websocket conection", http.StatusInternalServerError)
		return
	}
}

/////////////////////////////////////////////////////////////////////
////////////////////////// HELPER FUNCS /////////////////////////////
/////////////////////////////////////////////////////////////////////
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Error("ERROR - could not fetch template for rendering", err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Error("ERROR - could not render template", err)
		return err
	}
	return nil
}
