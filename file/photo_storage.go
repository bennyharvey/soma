package file

import (
	"fmt"
	"os"
	"path"
)

type PhotoStorage struct {
	path string
}

func NewPhotoStorage(path string) *PhotoStorage {
	return &PhotoStorage{path: path}
}

func (ps *PhotoStorage) AddPhoto(photoID string, photo []byte) error {
	dirPath, filePath := formPhotoPaths(ps.path, photoID)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("make dir path: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, err = f.Write(photo)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("close file: %w", err)
	}

	return nil
}

func (ps *PhotoStorage) PhotoPath(photoID string) string {
	_, photoPath := formPhotoPaths(ps.path, photoID)
	return photoPath
}

func formPhotoPaths(basePath, photoID string) (string, string) {
	if len(photoID) < 4 {
		return basePath, path.Join(basePath, photoID)
	}
	basePath = path.Join(basePath, photoID[:2], photoID[2:4])
	return basePath, path.Join(basePath, photoID)
}
