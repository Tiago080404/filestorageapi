package router

import (
	"fileserverapi/internal/handlers"
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload", handlers.Upload)

	return mux
}
