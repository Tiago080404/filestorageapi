package storage

import (
	"archive/zip"
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

var mockDirPath = "/home/tiago/fileservertest/"
var thumbnailPath = "/home/tiago/fileservertest/thumbnails"

func UploadLocal(fileHeader *multipart.FileHeader) error {
	path := filepath.Join(mockDirPath, fileHeader.Filename)
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
	dir, err := os.ReadDir(mockDirPath)
	if err != nil {
		log.Printf("could not open dir: %s", err)
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
		log.Printf("could not marshal to json: %s", err)
		return nil, err
	}

	return byteFiles, nil
}

func MakeThumbnail(path string) ([]byte, error) { //refactoren func name passt nicht mehr
	if thumbnailExists(path) {
		log.Println("thumbnail already exists: ", thumbnailPath, path)
		file, _ := os.Open(filepath.Join(thumbnailPath, path))
		img, _, err := image.Decode(file)
		if err != nil {
			log.Printf("could not open file: %s", err)
			return nil, err
		}

		var buf bytes.Buffer
		jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
		return buf.Bytes(), nil
	} else {
		log.Println("thumbnail does not exists")

		file, err := os.Open(filepath.Join(mockDirPath, path))
		if err != nil {
			log.Printf("could not open file: %s", err)
			return nil, err
		}
		defer file.Close()

		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("could not decode image: %s", err)
			return nil, err
		}

		thumbnail := imaging.Resize(image, 150, 150, imaging.Lanczos)

		pathh, _ := os.Create(filepath.Join(thumbnailPath, path))

		var buffer bytes.Buffer
		imaging.Encode(&buffer, thumbnail, imaging.JPEG)

		_, err = io.Copy(pathh, &buffer)
		if err != nil {
			log.Printf("Could not copy thumbnail into folder: %s", err)
			return nil, err
		}

		buf, err := encodeJpeg(thumbnail)
		if err != nil {
			log.Printf("could not encode jpeg for %s: %v", path, err)
			return nil, err
		}

		return buf, nil
	}
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

func thumbnailExists(name string) bool {
	_, err := os.Open(filepath.Join(thumbnailPath, name))
	if err != nil {
		return false
	}

	return true
}

func DownloadFiles(files []string) ([]byte, error) { //todo: change later without zip
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	for _, fileName := range files {
		file, err := os.Open(filepath.Join(mockDirPath, fileName))
		if err != nil {
			log.Println("could not open file to download: ", err)
			continue
		}

		w, err := zw.Create(fileName)
		if err != nil {
			file.Close()
			log.Println("could not create zip for file: ", fileName)
			log.Println(err)
			continue
		}
		_, err = io.Copy(w, file)
		if err != nil {
			file.Close()
			log.Println("Could not copy the file to the zip: ", err)
			continue
		}
		file.Close()
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
