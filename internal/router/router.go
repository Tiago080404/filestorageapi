package router

import (
	"fileserverapi/internal/handlers"
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload", handlers.Upload)
	mux.HandleFunc("/api/list", handlers.List)
	mux.HandleFunc("/api/thumbnail/{name}", handlers.Thumbnail)
	mux.HandleFunc("/api/download/{files...}", handlers.Download)
	return mux
}
