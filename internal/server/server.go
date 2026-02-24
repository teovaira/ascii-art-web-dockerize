// Package server provides HTTP server functionality for ASCII art web application.
package server

import (
	"net/http"
	"os"

	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
	"ascii-art-web/internal/validation"
)

func Start(addr string) error {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handleAsciiArt)
	return http.ListenAndServe(addr, nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Home"))
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	if banner == "" {
		banner = "standard"
	}

	if err := validation.ValidateText(text); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateBanner(banner); err != nil {
		http.Error(w, "Not Found: invalid banner", http.StatusNotFound)
		return
	}

	bannerPath := "cmd/ascii-art/testdata/" + banner + ".txt"
	bannerData, err := parser.LoadBanner(os.DirFS("."), bannerPath)
	if err != nil {
		http.Error(w, "Not Found: banner file not found", http.StatusNotFound)
		return
	}

	result, err := renderer.ASCII(text, bannerData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
