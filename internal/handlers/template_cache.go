package handlers

import (
	"html/template"
)

// NewTemplateCache parses HTML template files and stores them in a cache.
//
// It loads the base layout and the index page template at startup and returns
// a map keyed by template filename for efficient reuse across HTTP requests.
// The function must be called before any HTTP handler executes.
//
// Returns:
//   - A map of template name to parsed *template.Template.
//   - An error if any template file is missing or contains a parse error.
func NewTemplateCache() (map[string]*template.Template, error) {
	templateMap := make(map[string]*template.Template)
	ts, err := template.ParseFiles("templates/index.html", "templates/base.html")
	if err != nil {
		return nil, err
	}
	templateMap["index.html"] = ts
	return templateMap, nil
}
