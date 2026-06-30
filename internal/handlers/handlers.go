package handlers

import (
	"fileserverapi/internal/storage"
	"log"
	"net/http"
)

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
