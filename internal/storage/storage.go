package storage

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type FileInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

var mockDirPath = "/home/tiago/testneu"

func UploadLocal(fileHeader *multipart.FileHeader) error {
	path := filepath.Join("/home/tiago/Downloads/", fileHeader.Filename)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, file)
	return err
}

func GetDir() ([]byte, error) {
	var infos []FileInfo
	dir, err := os.ReadDir("/home/tiago/testneu")
	if err != nil {
		log.Fatal("could not open dir", err)
		return nil, err
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		infos = append(infos, FileInfo{Name: file.Name(), Url: "/thumbnail/" + file.Name()})
	}

	byteFiles, err := json.Marshal(infos)
	if err != nil {
		log.Fatal("could not marshal to json", err)
		return nil, err
	}

	return byteFiles, nil
}

func MakeThumbnail(path string) ([]byte, error) {
	file, err := os.Open(filepath.Join(mockDirPath, path))
	if err != nil {
		log.Fatal("could not open file", err)
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("could not decode image", err)
		return nil, err
	}

	thumbnail := imaging.Resize(image, 150, 150, imaging.Lanczos)

	buf, err := encodeJpeg(thumbnail)
	if err != nil {
		log.Printf("could not encode jpeg for %s: %v", path, err)
		return nil, err
	}

	return buf, nil
}

func encodeJpeg(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		log.Fatal("could not encode", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
