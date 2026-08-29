package handlers

import (
	"encoding/json"
	"fileserverapi/internal/database"
	"fileserverapi/internal/storage"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewDir struct {
	Name string `json:"name"`
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
	dir, err := storage.GetHomeDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dir)
}

func Download(w http.ResponseWriter, r *http.Request) {
	var mockDirPath = "/home/tiago/fileservertest/"

	files := r.PathValue("files")
	stat, err := os.Stat(filepath.Join(mockDirPath, files))
	if err != nil {
		log.Println("could not read stat")
	}

	if stat.IsDir() {
		folder, err := storage.DownloadFolder(files)
		if err != nil {
			log.Println("Could not download folder")
			http.Error(w, "could not download folder", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Write(folder)
	} else {
		downloadedFile, err := storage.DownloadFiles(files)
		if err != nil {
			log.Println("could not download files: ", err)
			http.Error(w, "could not download files", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(downloadedFile)
	}

}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := database.Authenticate(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
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

func CreateDir(w http.ResponseWriter, r *http.Request) {
	var req NewDir
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.MakeNewDir(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var registerReq LoginRequest //type is the same

	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		return
	}

	database.CreateUser(registerReq.Username, registerReq.Password)
}

func MoveData(w http.ResponseWriter, r *http.Request) { //kann das auch in body tun ist ja kein get
	/* 	data := r.PathValue("file")
	   	dest := r.PathValue("dest") */

	var dataMove struct {
		Data string `json:"file"`
		Dest string `json:"dest"`
	}

	json.NewDecoder(r.Body).Decode(&dataMove)
	err := storage.MoveFile(dataMove.Data, dataMove.Dest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Moved data"))
}

func Rename(w http.ResponseWriter, r *http.Request) {
	var newName struct {
		NewName string `json:"newname"`
		OldPath string `json:"oldpath"`
	}

	json.NewDecoder(r.Body).Decode(&newName)
	log.Println(newName.OldPath, newName.NewName)

	err := storage.RenameData(newName.NewName, newName.OldPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Remove(w http.ResponseWriter, r *http.Request) {
	var dataPath struct {
		Path string `json:"path"`
	}
	json.NewDecoder(r.Body).Decode(&dataPath)

	err := storage.RemoveData(dataPath.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListFolder(w http.ResponseWriter, r *http.Request) {
	getDir := r.PathValue("path")

	dir, err := storage.ListDir(getDir)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dir)
}

func OpenFile(w http.ResponseWriter, r *http.Request) {
	fileToOpen := r.PathValue("file")

	openedFile, err := storage.OpenFile(fileToOpen)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(openedFile)
}

func PreviewFile(w http.ResponseWriter, r *http.Request) {
	fileToPreview := r.PathValue("path")

	preview, err := storage.PreviewFile(fileToPreview)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "max-age=3600")
	w.WriteHeader(http.StatusOK)
	w.Write(preview)
}
