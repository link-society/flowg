package config

import (
	"os"
	"path/filepath"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"
)

type migratorProcH struct {
	baseDir string
}

var _ proctree.ProcessHandler = (*migratorProcH)(nil)

func (p *migratorProcH) Init(ctx actor.Context) proctree.ProcessResult {
	if err := p.migrateAlerts(ctx); err != nil {
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (p *migratorProcH) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (p *migratorProcH) Terminate(ctx actor.Context, err error) error {
	return err
}

func (p *migratorProcH) migrateAlerts(ctx actor.Context) error {
	alertsDir := filepath.Join(p.baseDir, "alerts")
	forwardersDir := filepath.Join(p.baseDir, "forwarders")

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
