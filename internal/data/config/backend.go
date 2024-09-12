package config

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type backend interface {
	PathExist(dir string) bool
	MakePath(dir string) error

	ListFiles(dir string) ([]string, error)
	OpenFile(path string) (io.ReadCloser, error)
	WriteFile(path string, content []byte) error
	DeleteFile(path string) error
}

type fsBackend struct {
	baseDir string
}

func newFsBackend(baseDir string) *fsBackend {
	return &fsBackend{baseDir: baseDir}
}

func (b *fsBackend) PathExist(dir string) bool {
	path := filepath.Join(b.baseDir, dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (b *fsBackend) MakePath(dir string) error {
	path := filepath.Join(b.baseDir, dir)
	return os.MkdirAll(path, os.ModePerm)
}

func (b *fsBackend) ListFiles(dir string) ([]string, error) {
	path := filepath.Join(b.baseDir, dir)
	files, err := os.ReadDir(path)
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

func (b *fsBackend) OpenFile(path string) (io.ReadCloser, error) {
	filePath := filepath.Join(b.baseDir, path)
	return os.Open(filePath)
}

func (b *fsBackend) WriteFile(path string, content []byte) error {
	filePath := filepath.Join(b.baseDir, path)
	return os.WriteFile(filePath, content, os.ModePerm)
}

func (b *fsBackend) DeleteFile(path string) error {
	filePath := filepath.Join(b.baseDir, path)
	return os.Remove(filePath)
}

type memBackend struct {
	files map[string][]byte
}

func newMemBackend() *memBackend {
	return &memBackend{
		files: make(map[string][]byte),
	}
}

func (b *memBackend) PathExist(dir string) bool {
	dir = filepath.Clean(dir) + string(filepath.Separator)

	for path := range b.files {
		path = filepath.Clean(path)
		if strings.HasPrefix(path, dir) {
			return true
		}
	}

	return false
}

func (b *memBackend) MakePath(dir string) error {
	return nil
}

func (b *memBackend) ListFiles(dir string) ([]string, error) {
	dir = filepath.Clean(dir)

	items := []string{}

	for path := range b.files {
		path = filepath.Clean(path)
		if filepath.Dir(path) == dir {
			items = append(items, path)
		}
	}

	return items, nil
}

func (b *memBackend) OpenFile(path string) (io.ReadCloser, error) {
	content, ok := b.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	return io.NopCloser(strings.NewReader(string(content))), nil
}

func (b *memBackend) WriteFile(path string, content []byte) error {
	b.files[path] = content
	return nil
}

func (b *memBackend) DeleteFile(path string) error {
	delete(b.files, path)
	return nil
}
