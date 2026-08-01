package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameData(t *testing.T) {
	tmpDir := t.TempDir()
	thumbDir := t.TempDir()

	origThumbnailPath := thumbnailPath
	thumbnailPath = thumbDir
	defer func() { thumbnailPath = origThumbnailPath }()

	oldPath := filepath.Join(tmpDir, "old.JPG")
	newPath := filepath.Join(tmpDir, "new.JPG")

	err := os.WriteFile(oldPath, []byte("data"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(thumbDir, "old.JPG"), []byte("thumb"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = RenameData(newPath, oldPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMovefile(t *testing.T) {
	tmpDir := t.TempDir()

	origMockDirPath := mockDirPath
	mockDirPath = tmpDir
	defer func() { mockDirPath = origMockDirPath }()

	err := os.Mkdir(filepath.Join(tmpDir, "destFolder"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "old.txt"), []byte("data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = MoveFile("old.txt", "destFolder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
