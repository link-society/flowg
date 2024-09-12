package config

import (
	"io"
	"path/filepath"
)

type StorageOpts struct {
	dir      string
	inMemory bool
}

func DefaultStorageOpts() StorageOpts {
	return StorageOpts{
		dir:      "./data/config",
		inMemory: false,
	}
}

func (s StorageOpts) WithDir(dir string) StorageOpts {
	s.dir = dir
	return s
}

func (s StorageOpts) WithInMemory(inMemory bool) StorageOpts {
	s.inMemory = inMemory
	return s
}

type Storage struct {
	backend backend
}

func NewStorage(opts StorageOpts) *Storage {
	var b backend

	if opts.inMemory {
		b = newMemBackend()
	} else {
		b = newFsBackend(opts.dir)
	}

	if !b.PathExist("") {
		b.MakePath("")
	}

	return &Storage{backend: b}
}

func (s *Storage) listStorageTypeItems(storageType string) ([]string, error) {
	s.backend.MakePath(storageType)
	return s.backend.ListFiles(storageType)
}

func (s *Storage) openStorageTypeItem(storageType, name string) (io.ReadCloser, error) {
	s.backend.MakePath(storageType)
	return s.backend.OpenFile(filepath.Join(storageType, name))
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
	s.backend.MakePath(storageType)
	return s.backend.WriteFile(filepath.Join(storageType, name), []byte(content))
}

func (s *Storage) deleteStorageTypeItem(storageType, name string) error {
	s.backend.MakePath(storageType)
	return s.backend.DeleteFile(filepath.Join(storageType, name))
}
