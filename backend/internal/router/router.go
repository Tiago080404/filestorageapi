package router

import (
	"fileserverapi/internal/auth"
	"fileserverapi/internal/handlers"
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/upload", auth.ProtectedRoutes(handlers.Upload))
	mux.HandleFunc("GET /api/list", auth.ProtectedRoutes(handlers.List))
	mux.HandleFunc("GET /api/download/{files...}", auth.ProtectedRoutes(handlers.Download))
	mux.HandleFunc("POST /api/login", handlers.Authenticate)
	mux.HandleFunc("POST /api/dir/create", handlers.CreateDir)
	mux.HandleFunc("POST /api/register", handlers.Register)
	mux.HandleFunc("PATCH /api/move/{file}/{dest}", handlers.MoveData)
	mux.HandleFunc("PATCH /api/rename", handlers.Rename)
	mux.HandleFunc("PATCH /api/remove", handlers.Remove)
	return mux
}
