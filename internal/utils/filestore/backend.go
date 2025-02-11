package filestore

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type backend interface {
	ListFiles() ([]string, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, content []byte) error
	DeleteFile(path string) error
	StatFile(path string) (fs.FileInfo, error)
}

type fsBackend struct {
	baseDir string
}

func newFsBackend(baseDir string) (*fsBackend, error) {
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &fsBackend{baseDir: baseDir}, nil
}

func (b *fsBackend) ListFiles() ([]string, error) {
	files, err := os.ReadDir(b.baseDir)
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

func (b *fsBackend) ReadFile(path string) ([]byte, error) {
	filePath := filepath.Join(b.baseDir, path)
	return os.ReadFile(filePath)
}

func (b *fsBackend) WriteFile(path string, content []byte) error {
	filePath := filepath.Join(b.baseDir, path)
	return os.WriteFile(filePath, content, os.ModePerm)
}

func (b *fsBackend) DeleteFile(path string) error {
	filePath := filepath.Join(b.baseDir, path)
	return os.Remove(filePath)
}

func (b *fsBackend) StatFile(path string) (fs.FileInfo, error) {
	filePath := filepath.Join(b.baseDir, path)
	return os.Stat(filePath)
}

type memBackend struct {
	files map[string][]byte
}

type memFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func newMemBackend() *memBackend {
	return &memBackend{
		files: make(map[string][]byte),
	}
}

func (b *memBackend) ListFiles() ([]string, error) {
	items := []string{}

	for path := range b.files {
		items = append(items, path)
	}

	return items, nil
}

func (b *memBackend) ReadFile(path string) ([]byte, error) {
	content, ok := b.files[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	return content, nil
}

func (b *memBackend) WriteFile(path string, content []byte) error {
	b.files[path] = content
	return nil
}

func (b *memBackend) DeleteFile(path string) error {
	delete(b.files, path)
	return nil
}

func (b *memBackend) StatFile(path string) (fs.FileInfo, error) {
	content, err := b.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &memFileInfo{
		name:    path,
		size:    int64(len(content)),
		mode:    os.ModePerm,
		modTime: time.Now(),
	}, nil
}

func (f *memFileInfo) Name() string {
	return f.name
}

func (f *memFileInfo) Size() int64 {
	return f.size
}

func (f *memFileInfo) Mode() os.FileMode {
	return f.mode
}

func (f *memFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *memFileInfo) IsDir() bool {
	return false
}

func (f *memFileInfo) Sys() any {
	return nil
}
