package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"encoding/base64"
)

func migrateAlerts(baseDir string) error {
	alertsDir := filepath.Join(baseDir, "alerts")
	forwardersDir := filepath.Join(baseDir, "forwarders")

	if _, err := os.Stat(alertsDir); os.IsNotExist(err) {
		return nil
	}

	if _, err := os.Stat(forwardersDir); os.IsNotExist(err) {
		if err := os.MkdirAll(forwardersDir, 0755); err != nil {
			return err
		}
	}

	alerts, err := os.ReadDir(alertsDir)
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		if !alert.IsDir() {
			err := os.Rename(
				filepath.Join(alertsDir, alert.Name()),
				filepath.Join(forwardersDir, alert.Name()),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func migrateToBadger(ctx context.Context, baseDir string, storage *storageImpl) error {
	fileStorages := []struct {
		dir       string
		extension string
		converter func(string, []byte) error
	}{
		{
			dir:       "transformers",
			extension: ".vrl",
			converter: func(name string, content []byte) error {
				return storage.writeItem(ctx, transformerItemType, name, content)
			},
		},
		{
			dir:       "pipelines",
			extension: ".json",
			converter: func(name string, content []byte) error {
				return storage.writeItem(ctx, pipelineItemType, name, content)
			},
		},
		{
			dir:       "forwarders",
			extension: ".json.b64",
			converter: func(name string, b64content []byte) error {
				content := make([]byte, base64.StdEncoding.DecodedLen(len(b64content)))
				n, err := base64.StdEncoding.Decode(content, b64content)
				if err != nil {
					return err
				}

				return storage.writeItem(ctx, forwarderItemType, name, content[:n])
			},
		},
	}

	for _, storage := range fileStorages {
		dir := filepath.Join(baseDir, storage.dir)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), storage.extension) {
				path := filepath.Join(dir, file.Name())
				itemName := strings.TrimSuffix(file.Name(), storage.extension)

				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if err := storage.converter(itemName, content); err != nil {
					return err
				}

				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
