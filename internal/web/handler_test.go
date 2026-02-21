package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome_TemplateMissing_Returns404(t *testing.T) {
	app := &Application{
		TemplateCache: map[string]*template.Template{}, // empty cache
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Home(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.StatusCode)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "The template does not exist") {
		t.Fatalf("expected body to contain %q, got %q", "The template does not exist", body)
	}
}

func TestHome_TemplateExists_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		template       *template.Template
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success execution",
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusOK,
			expectedBody:   "<h1>HOME</h1>",
		},
		{
			name:           "execution failure",
			template:       template.Must(template.New("index.html").Parse(`{{call .}}`)),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to connect to the internal service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				TemplateCache: map[string]*template.Template{
					"index.html": tt.template,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			app.Home(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Fatalf("expected body to contain %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
