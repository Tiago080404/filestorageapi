package main

import (
	"fileserverapi/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.Setup()

	log.Println("Server läuft auf :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
