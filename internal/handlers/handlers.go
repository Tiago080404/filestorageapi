package handlers

import (
	"fmt"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println(header.Filename)
	fmt.Println(header.Size)
	w.Write([]byte("Upload received"))
}
