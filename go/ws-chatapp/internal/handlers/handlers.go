package handlers

import (
	"fmt"
	"net/http"
	"sort"

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

	wsChan  = make(chan WsJsonPayload)
	clients = make(map[WebSocketConn]string)
)

type WebSocketConn struct {
	*websocket.Conn
}

type WsJsonPayload struct {
	Username string `json:"username"`
	Action   string `json:"action"`
	Message  string `json:"message"`

	// This is to leave it out of the json, ie not show up in json
	// Rather this is an internal field that will be populated on the server side
	Conn WebSocketConn `json:"-"`
}

type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
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

	conn := WebSocketConn{
		Conn: ws,
	}

	go listenForWs(&conn)

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

// This function runs as a goroutine for each websocket connection we establish with a client.
// The purpose of this goroutine is to run forever.
// If there is a panic while running this goroutine, we should recover and continue.
func listenForWs(conn *WebSocketConn) {
	defer func() {
		// Recover is a built-in function that regains control of a panicking goroutine.
		// Recover is only useful inside deferred functions.
		// During normal execution, a call to recover will return nil and have no other effect.
		// For more on defer and recover behavior : https://go.dev/blog/defer-panic-and-recover
		if r := recover(); r != nil {
			log.Error("ERROR - Encountered an error in listening to websocket connections. Recovered from panic in goroutine.", r)
		}
	}()

	var payload WsJsonPayload
	for {
		err := conn.ReadJSON(&payload)
		// If there is an error reading the data as json from the client, dont do anything
		if err == nil {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenToWsChan() {
	var resp WsJsonResponse
	for {
		ev := <-wsChan
		switch ev.Action {
		case "addUser":
			// get a list of users and broadcast it
			if len(ev.Username) > 0 {
				clients[ev.Conn] = ev.Username
			}
			broadcastUserListToAll(resp)
		case "userLeft":
			delete(clients, ev.Conn)
			broadcastUserListToAll(resp)
		}
	}
}

func getUserList() []string {
	var users []string
	for _, user := range clients {
		if len(user) > 0 {
			users = append(users, user)
		}
	}
	sort.Strings(users)
	return users
}

func broadcastUserListToAll(resp WsJsonResponse) {
	users := getUserList()
	resp.Action = "listUsers"
	resp.ConnectedUsers = users
	resp.Message = "List of users"
	broadcastToAll(resp)
}

func broadcastToAll(resp WsJsonResponse) {
	for clientConn, username := range clients {
		err := clientConn.WriteJSON(resp)
		if err != nil {
			log.Error(fmt.Sprintf("ERROR - Encountered an error in broadcasting to client connection for user %s", username), err)
			closeErr := clientConn.Close()
			if closeErr != nil {
				log.Error(fmt.Sprintf("ERROR - Could not close client connection for user %s, possibly because we lost the client. Ignoring it, and moving on to removing it from our list of conns.", username), err)
			}
			delete(clients, clientConn)
		}
	}
}
