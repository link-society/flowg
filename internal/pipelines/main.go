package pipelines

import (
	"encoding/json"
	"os"
	"path/filepath"

	"link-society.com/flowg/internal/logstorage"
)

type Manager struct {
	transformersDir string
	pipelinesDir    string

	db *logstorage.Storage
}

func NewManager(db *logstorage.Storage, configDir string) *Manager {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	transformersDir := filepath.Join(configDir, "transformers")
	if _, err := os.Stat(transformersDir); os.IsNotExist(err) {
		os.Mkdir(transformersDir, os.ModePerm)
	}

	pipelinesDir := filepath.Join(configDir, "pipelines")
	if _, err := os.Stat(pipelinesDir); os.IsNotExist(err) {
		os.Mkdir(pipelinesDir, os.ModePerm)
	}

	return &Manager{
		transformersDir: transformersDir,
		pipelinesDir:    pipelinesDir,
		db:              db,
	}
}

func (m *Manager) ListTransformers() ([]string, error) {
	files, err := os.ReadDir(m.transformersDir)
	if err != nil {
		return nil, err
	}

	transformers := []string{}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".vrl" {
			transformer := file.Name()[0 : len(file.Name())-4]
			transformers = append(transformers, transformer)
		}
	}

	return transformers, nil
}

func (m *Manager) GetTransformerScript(name string) (string, error) {
	filePath := filepath.Join(m.transformersDir, name+".vrl")
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func (m *Manager) SaveTransformerScript(name, script string) error {
	filePath := filepath.Join(m.transformersDir, name+".vrl")
	return os.WriteFile(filePath, []byte(script), os.ModePerm)
}

func (m *Manager) DeleteTransformerScript(name string) error {
	filePath := filepath.Join(m.transformersDir, name+".vrl")
	return os.Remove(filePath)
}

func (m *Manager) ListPipelines() ([]string, error) {
	files, err := os.ReadDir(m.pipelinesDir)
	if err != nil {
		return nil, err
	}

	pipelines := []string{}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			pipeline := file.Name()[0 : len(file.Name())-5]
			pipelines = append(pipelines, pipeline)
		}
	}

	return pipelines, nil
}

func (m *Manager) GetPipelineFlow(name string) (string, error) {
	filePath := filepath.Join(m.pipelinesDir, name+".json")
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func (m *Manager) SavePipelineFlow(name, flow string) error {
	filePath := filepath.Join(m.pipelinesDir, name+".json")
	return os.WriteFile(filePath, []byte(flow), os.ModePerm)
}

func (m *Manager) DeletePipelineFlow(name string) error {
	filePath := filepath.Join(m.pipelinesDir, name+".json")
	return os.Remove(filePath)
}

func (m *Manager) GetPipeline(name string) (*Pipeline, error) {
	filePath := filepath.Join(m.pipelinesDir, name+".json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var flowGraph FlowGraph
	if err := json.NewDecoder(file).Decode(&flowGraph); err != nil {
		return nil, err
	}

	return flowGraph.BuildPipeline()
}
