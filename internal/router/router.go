package router

import (
	"fileserverapi/internal/auth"
	"fileserverapi/internal/handlers"
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload", auth.ProtectedRoutes(handlers.Upload))
	mux.HandleFunc("/api/list", auth.ProtectedRoutes(handlers.List))
	mux.HandleFunc("/api/thumbnail/{name}", auth.ProtectedRoutes(handlers.Thumbnail))
	mux.HandleFunc("/api/download/{files...}", auth.ProtectedRoutes(handlers.Download))
	mux.HandleFunc("/api/login", handlers.Authenticate)
	return mux
}
