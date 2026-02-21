package web

import (
	"html/template"
	"net/http"
)

type Application struct {
	TemplateCache map[string]*template.Template
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	ts, found := app.TemplateCache["index.html"]
	if !found {
		http.Error(w, "The template does not exist", http.StatusNotFound)
		return
	}

	err := ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to connect to the internal service", http.StatusInternalServerError)
	}

}
