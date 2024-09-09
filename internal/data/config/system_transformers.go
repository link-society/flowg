package config

import "path/filepath"

const TRANSFORMERS_STORAGE_TYPE = "transformers"

type TransformerSystem struct {
	backend *Storage
}

func NewTransformerSystem(backend *Storage) *TransformerSystem {
	backend.resolveStorageTypeDir(TRANSFORMERS_STORAGE_TYPE)
	return &TransformerSystem{backend: backend}
}

func (sys *TransformerSystem) List() ([]string, error) {
	items, err := sys.backend.listStorageTypeItems(TRANSFORMERS_STORAGE_TYPE)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, item := range items {
		if filepath.Ext(item) == ".vrl" {
			results = append(results, item[:len(item)-4])
		}
	}

	return results, nil
}

func (sys *TransformerSystem) Read(name string) (string, error) {
	return sys.backend.readStorageTypeItem(TRANSFORMERS_STORAGE_TYPE, name+".vrl")
}

func (sys *TransformerSystem) Write(name, script string) error {
	return sys.backend.writeStorageTypeItem(TRANSFORMERS_STORAGE_TYPE, name+".vrl", script)
}

func (sys *TransformerSystem) Delete(name string) error {
	return sys.backend.deleteStorageTypeItem(TRANSFORMERS_STORAGE_TYPE, name+".vrl")
}
