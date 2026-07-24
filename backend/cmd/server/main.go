package main

import (
	"fileserverapi/internal/auth"
	"fileserverapi/internal/database"
	"fileserverapi/internal/router"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := auth.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Println("couild not read .env")
	}

	err = database.Conn()
	if err != nil {
		log.Fatal(err)
	}
	r := router.Setup()

	log.Println("Server läuft auf :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
