package storage

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadLocal(file *multipart.File, fileName string) error {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, *file)
	if err != nil {
		return err
	}

	path := filepath.Join("/home/tiago/Downloads/", fileName)
	return os.WriteFile(path, buf.Bytes(), 0644)
}
