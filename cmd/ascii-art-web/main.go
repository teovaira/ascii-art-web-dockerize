package main

import (
	"fmt"
	"log"
	"os"

	"ascii-art-web/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
