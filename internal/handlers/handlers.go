package handlers

import (
	"encoding/json"
	"fileserverapi/internal/auth"
	"fileserverapi/internal/storage"
	"log"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files[]"]

	for _, fileHeader := range files {
		err = storage.UploadLocal(fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Write([]byte("Upload received"))

}

func List(w http.ResponseWriter, r *http.Request) {
	dir, err := storage.GetDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dir)
}

func Thumbnail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	data, err := storage.MakeThumbnail(name)
	if err != nil {
		log.Printf("GetThumbnail failed for %s: %v", name, err)
		http.Error(w, "could not create thumbnail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

func Download(w http.ResponseWriter, r *http.Request) {
	files := r.PathValue("files")
	downloadedFiles, err := storage.DownloadFiles(strings.Split(files, "/"))
	if err != nil {
		log.Println("could not download files: ", err)
		http.Error(w, "could not download files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Write(downloadedFiles)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
	}

	token, err := auth.Authenticate(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	cookie := http.Cookie{
		Name:     "auth",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		SameSite: http.SameSiteLaxMode,
		Domain:   "localhost",
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
