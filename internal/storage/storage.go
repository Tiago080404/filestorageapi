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

	if thumbnailExists(fileHeader.Filename) {
		log.Println("Thumbnail already exists")
	} else {
		createThumbnail(fileHeader.Filename)
	}

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
func createThumbnail(fileName string) error {
	log.Println("thumbnail does not exists")

	file, err := os.Open(filepath.Join(mockDirPath, fileName))
	if err != nil {
		log.Printf("could not open file: %s", err)
		return err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		log.Printf("could not decode image: %s", err)
		return err
	}

	thumbnail := imaging.Resize(image, 150, 150, imaging.Lanczos)

	pathh, _ := os.Create(filepath.Join(thumbnailPath, fileName))

	var buffer bytes.Buffer
	imaging.Encode(&buffer, thumbnail, imaging.JPEG)

	_, err = io.Copy(pathh, &buffer)
	if err != nil {
		log.Printf("Could not copy thumbnail into folder: %s", err)
		return err
	}

	_, err = encodeJpeg(thumbnail)
	if err != nil {
		log.Printf("could not encode jpeg for %s: %v", fileName, err)
		return err
	}

	return nil
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

func DownloadFiles(file string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(mockDirPath, file))
	if err != nil {
		return nil, err
	}

	return data, nil

}

func DownloadFolder(folder string) ([]byte, error) {
	dir, err := os.ReadDir(filepath.Join(mockDirPath, folder))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, entry := range dir {
		file, err := w.Create(entry.Name())
		if err != nil {
			log.Println("Coulkd not crate file zip")
			return nil, err
		}
		data, err := os.ReadFile(filepath.Join(mockDirPath, folder, entry.Name()))
		if err != nil {
			log.Println("could not read file")
			return nil, err

		}

		_, err = file.Write(data)
		if err != nil {
			log.Println("could not write in zip file")
			return nil, err
		}

	}
	err = w.Close()
	if err != nil {
		log.Println("could not close w: ", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func MakeNewDir(dir string) error {
	err := os.Mkdir(filepath.Join(mockDirPath, dir), 0750)
	if err != nil {
		log.Println("Could not create dir: ", err)
		return err
	}
	return err
}
