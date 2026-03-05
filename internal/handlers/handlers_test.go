package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestHome_TableDriven tests the Home handler with multiple scenarios.
func TestHome_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		template       *template.Template
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success execution",
			path:           "/",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusOK,
			expectedBody:   "<h1>HOME</h1>",
		},
		{
			name:           "execution failure",
			path:           "/",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse(`{{call .}}`)),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
		{
			name:           "template missing",
			path:           "/",
			method:         http.MethodGet,
			template:       nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not Found",
		},
		{
			name:           "unknown path returns 404",
			path:           "/random-page",
			method:         http.MethodGet,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
		{
			name:           "wrong method returns 405",
			path:           "/",
			method:         http.MethodPost,
			template:       template.Must(template.New("index.html").Parse("<h1>HOME</h1>")),
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := make(map[string]*template.Template)
			if tt.template != nil {
				cache["index.html"] = tt.template
			}

			app := &Application{TemplateCache: cache}

			req := httptest.NewRequest(tt.method, tt.path, nil)
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

// TestNewTemplateCache tests the template cache initialization.
func TestNewTemplateCache(t *testing.T) {
	t.Run("returns error on missing files", func(t *testing.T) {
		_, err := NewTemplateCache()
		if err == nil {
			t.Fatal("expected error when template files are missing, got nil")
		}
	})
}

// TestGenerateASCII tests the core ASCII generation logic.
func TestGenerateASCII(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		banner     string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "valid request",
			text:       "Hello",
			banner:     "standard",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "default banner",
			text:       "Hello",
			banner:     "",
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "empty text",
			text:       "",
			banner:     "standard",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "text too long",
			text:       strings.Repeat("a", 1001),
			banner:     "standard",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "invalid banner",
			text:       "Hello",
			banner:     "invalid",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, status, err := GenerateASCII(tt.text, tt.banner)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateASCII() error = %v, wantErr %v", err, tt.wantErr)
			}

			if status != tt.wantStatus {
				t.Errorf("GenerateASCII() status = %v, want %v", status, tt.wantStatus)
			}

			if !tt.wantErr && result == "" {
				t.Errorf("GenerateASCII() returned empty result for valid input")
			}
		})
	}
}

// TestHandleASCIIArt tests the HandleASCIIArt handler with multiple scenarios.
func TestHandleASCIIArt(t *testing.T) {
	tmpl := template.Must(template.New("index.html").Parse(
		`{{if .Error}}ERROR:{{.Error}}{{end}}{{if .Result}}RESULT:{{.Result}}{{end}}`,
	))

	tests := []struct {
		name           string
		method         string
		formData       string
		template       *template.Template
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "wrong method returns 405",
			method:         http.MethodGet,
			formData:       "",
			template:       tmpl,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name:           "template missing returns 404",
			method:         http.MethodPost,
			formData:       "text=Hello&banner=standard",
			template:       nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Not Found",
		},
		{
			name:           "empty text returns 400 with error in body",
			method:         http.MethodPost,
			formData:       "text=&banner=standard",
			template:       tmpl,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "ERROR:text cannot be empty",
		},
		{
			name:           "invalid banner returns 404 with error in body",
			method:         http.MethodPost,
			formData:       "text=Hello&banner=invalid",
			template:       tmpl,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "ERROR:invalid banner name",
		},
		{
			name:           "valid request returns 200 with result in body",
			method:         http.MethodPost,
			formData:       "text=Hi&banner=standard",
			template:       tmpl,
			expectedStatus: http.StatusOK,
			expectedBody:   "RESULT:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := make(map[string]*template.Template)
			if tt.template != nil {
				cache["index.html"] = tt.template
			}

			app := &Application{TemplateCache: cache}

			req := httptest.NewRequest(tt.method, "/ascii-art", strings.NewReader(tt.formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			app.HandleASCIIArt(rr, req)

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
