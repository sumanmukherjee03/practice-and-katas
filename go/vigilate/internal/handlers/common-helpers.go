package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// ClientError will display error page for client error i.e. bad request
func ClientError(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		show404(w, r)
	case http.StatusInternalServerError:
		show500(w, r)
	default:
		http.Error(w, http.StatusText(status), status)
	}
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Trace(trace)
	show500(w, r)
}

func show404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/404.html")
}

func show500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./ui/static/500.html")
}

type errRespJSON struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// ClientError will display error page for client error i.e. bad request
func ClientErrorJSON(w http.ResponseWriter, r *http.Request, status int) {
	switch status {
	case http.StatusNotFound:
		returnErrorJSON(w, r, status, "Entity could not be found")
	case http.StatusInternalServerError:
		returnErrorJSON(w, r, status, "Internal server error")
	default:
		returnErrorJSON(w, r, status, "Encountered an error")
	}
}

// ServerError will display error page for internal server error
func ServerErrorJSON(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Trace(trace)
	returnErrorJSON(w, r, http.StatusInternalServerError, "Internal server error")
}

func returnErrorJSON(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	var resp errRespJSON
	resp.OK = false
	resp.Message = msg
	out, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func printTemplateError(w http.ResponseWriter, err error) {
	_, _ = fmt.Fprint(w, fmt.Sprintf(`<small><span class='text-danger'>Error executing template: %s</span></small>`, err))
}

func broadcastMessage(channel, eventType string, payload map[string]string) {
	err := app.WsClient.Trigger(channel, eventType, payload)
	if err != nil {
		log.Error(err)
		return
	}
}
