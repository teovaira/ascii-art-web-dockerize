// Package main provides the ASCII art web server entry point.
//
// It reads the PORT environment variable (defaulting to 8080), initializes
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
	"time"

	"ascii-art-web-dockerize/internal/handlers"
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
	http.HandleFunc("/ascii-art", app.HandleASCIIArt)

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}
