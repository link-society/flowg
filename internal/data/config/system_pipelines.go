package config

import (
	"encoding/json"
	"path/filepath"
)

const PIPELINES_STORAGE_TYPE = "pipelines"

type PipelineSystem struct {
	backend *Storage
}

func NewPipelineSystem(backend *Storage) *PipelineSystem {
	return &PipelineSystem{backend: backend}
}

func (sys *PipelineSystem) List() ([]string, error) {
	items, err := sys.backend.listStorageTypeItems(PIPELINES_STORAGE_TYPE)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, item := range items {
		if filepath.Ext(item) == ".json" {
			results = append(results, item[:len(item)-5])
		}
	}

	return results, nil
}

func (sys *PipelineSystem) Read(name string) (string, error) {
	return sys.backend.readStorageTypeItem(PIPELINES_STORAGE_TYPE, name+".json")
}

func (sys *PipelineSystem) Write(name, flow string) error {
	return sys.backend.writeStorageTypeItem(PIPELINES_STORAGE_TYPE, name+".json", flow)
}

func (sys *PipelineSystem) Delete(name string) error {
	return sys.backend.deleteStorageTypeItem(PIPELINES_STORAGE_TYPE, name+".json")
}

func (sys *PipelineSystem) Parse(name string) (*FlowGraph, error) {
	file, err := sys.backend.openStorageTypeItem(PIPELINES_STORAGE_TYPE, name+".json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var flowGraph *FlowGraph
	if err := json.NewDecoder(file).Decode(&flowGraph); err != nil {
		return nil, err
	}

	return flowGraph, nil
}
