package config

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	Dir string
}

func NewStorage(dir string) *Storage {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	return &Storage{Dir: dir}
}

func (s *Storage) resolveStorageTypeDir(storageType string) string {
	storageDir := filepath.Join(s.Dir, storageType)
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		os.Mkdir(storageDir, os.ModePerm)
	}

	return storageDir
}

func (s *Storage) listStorageTypeItems(storageType string) ([]string, error) {
	storageDir := s.resolveStorageTypeDir(storageType)

	files, err := os.ReadDir(storageDir)
	if err != nil {
		return nil, err
	}

	items := []string{}

	for _, file := range files {
		if !file.IsDir() {
			item := file.Name()
			items = append(items, item)
		}
	}

	return items, nil
}

func (s *Storage) openStorageTypeItem(storageType, name string) (*os.File, error) {
	storageDir := s.resolveStorageTypeDir(storageType)
	filePath := filepath.Join(storageDir, name)
	return os.Open(filePath)
}

func (s *Storage) readStorageTypeItem(storageType, name string) (string, error) {
	file, err := s.openStorageTypeItem(storageType, name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *Storage) writeStorageTypeItem(storageType, name, content string) error {
	storageDir := s.resolveStorageTypeDir(storageType)
	filePath := filepath.Join(storageDir, name)
	return os.WriteFile(filePath, []byte(content), os.ModePerm)
}

func (s *Storage) deleteStorageTypeItem(storageType, name string) error {
	storageDir := s.resolveStorageTypeDir(storageType)
	filePath := filepath.Join(storageDir, name)
	return os.Remove(filePath)
}
