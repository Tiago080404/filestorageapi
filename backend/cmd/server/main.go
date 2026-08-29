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
	handler := corsMiddleware(r)

	log.Println("Server läuft auf :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PATCH, GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
