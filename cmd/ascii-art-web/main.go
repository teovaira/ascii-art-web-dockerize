// Package main provides the ASCII art web server entry point.
//
// It reads the PORT environment variable (defaulting to 8080), initialises
// the HTML template cache, registers HTTP routes, and starts the server.
//
// Usage:
//
//	go run ./cmd/ascii-art-web
//	PORT=9090 go run ./cmd/ascii-art-web
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ascii-art-web/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cache, err := handlers.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app := &handlers.Application{TemplateCache: cache}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/ascii-art", app.HandleAsciiArt)

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
