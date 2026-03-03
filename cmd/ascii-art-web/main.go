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
