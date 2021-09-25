package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet"
	log "github.com/sirupsen/logrus"
)

var views = jet.NewHTMLSet("./html")

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Error("ERROR - could not fetch template for rendering", err)
		http.Error(w, "Could not render home page", http.StatusInternalServerError)
		return
	}
}

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
