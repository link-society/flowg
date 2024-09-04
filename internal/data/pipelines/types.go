package pipelines

import (
	"context"
	"errors"
	"sync"

	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/ffi/vrl"
)

type Pipeline struct {
	Name string
	Root Node
}

func (p *Pipeline) Run(
	ctx context.Context,
	manager *Manager,
	entry *logstorage.LogEntry,
) error {
	metrics.IncPipelineLogCounter(p.Name)
	return p.Root.Process(ctx, manager, entry)
}

type Node interface {
	Process(ctx context.Context, manager *Manager, entry *logstorage.LogEntry) error
}

type TransformNode struct {
	TransformerName string
	Next            []Node
}

type SwitchNode struct {
	Condition logstorage.Filter
	Next      []Node
}

type RouterNode struct {
	Stream string
}

func (n *TransformNode) Process(
	ctx context.Context,
	manager *Manager,
	entry *logstorage.LogEntry,
) error {
	vrlScript, err := manager.GetTransformerScript(n.TransformerName)
	if err != nil {
		return err
	}

	output, err := vrl.ProcessRecord(entry.Fields, vrlScript)
	if err != nil {
		return err
	}

	errC := make(chan error, len(n.Next))
	wg := sync.WaitGroup{}

	for _, next := range n.Next {
		wg.Add(1)
		go func(next Node) {
			defer wg.Done()

			newEntry := &logstorage.LogEntry{
				Timestamp: entry.Timestamp,
				Fields:    output,
			}
			err := next.Process(ctx, manager, newEntry)
			if err != nil {
				errC <- err
			}
		}(next)
	}

	wg.Wait()
	close(errC)

	var errs []error
	for err := range errC {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (n *SwitchNode) Process(
	ctx context.Context,
	manager *Manager,
	entry *logstorage.LogEntry,
) error {
	if n.Condition.Evaluate(entry) {
		errC := make(chan error, len(n.Next))
		wg := sync.WaitGroup{}

		for _, next := range n.Next {
			wg.Add(1)
			go func(next Node) {
				defer wg.Done()
				err := next.Process(ctx, manager, entry)
				if err != nil {
					errC <- err
				}
			}(next)
		}

		wg.Wait()
		close(errC)

		var errs []error
		for err := range errC {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (n *RouterNode) Process(
	ctx context.Context,
	manager *Manager,
	entry *logstorage.LogEntry,
) error {
	key, err := manager.collectorSys.Ingest(ctx, n.Stream, entry)
	if err == nil {
		manager.notifier.Notify(n.Stream, string(key), *entry)
		metrics.IncStreamLogCounter(n.Stream)
	}

	return err
}
