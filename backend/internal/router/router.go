package router

import (
	"fileserverapi/internal/handlers"
	"net/http"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/upload", handlers.Upload)
	mux.HandleFunc("GET /api/list", handlers.List)
	mux.HandleFunc("GET /api/list/{path...}", handlers.ListFolder)
	mux.HandleFunc("GET /api/download/{files...}", handlers.Download)
	mux.HandleFunc("POST /api/login", handlers.Authenticate)
	mux.HandleFunc("POST /api/dir/create", handlers.CreateDir)
	mux.HandleFunc("POST /api/register", handlers.Register)
	mux.HandleFunc("PATCH /api/move/{file}/{dest}", handlers.MoveData)
	mux.HandleFunc("PATCH /api/rename", handlers.Rename)
	mux.HandleFunc("PATCH /api/remove", handlers.Remove)
	mux.HandleFunc("GET /api/open/{file...}", handlers.OpenFile)
	return mux
}

//auth.ProtectedRoutes()
