package pipelines

import (
	"context"
	"errors"
	"sync"

	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/ffi/vrl"
)

type Node interface {
	Process(ctx context.Context, entry *logstorage.LogEntry) error
}

type TransformNode struct {
	TransformerName string
	Next            []Node
}

type SwitchNode struct {
	Condition logstorage.Filter
	Next      []Node
}

type PipelineNode struct {
	Pipeline string
}

type RouterNode struct {
	Stream string
}

func (n *TransformNode) Process(ctx context.Context, entry *logstorage.LogEntry) error {
	transformerSys := getTransformerSystem(ctx)
	vrlScript, err := transformerSys.Read(n.TransformerName)
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
			err := next.Process(ctx, newEntry)
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

func (n *SwitchNode) Process(ctx context.Context, entry *logstorage.LogEntry) error {
	if n.Condition.Evaluate(entry) {
		errC := make(chan error, len(n.Next))
		wg := sync.WaitGroup{}

		for _, next := range n.Next {
			wg.Add(1)
			go func(next Node) {
				defer wg.Done()
				err := next.Process(ctx, entry)
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

func (n *PipelineNode) Process(ctx context.Context, entry *logstorage.LogEntry) error {
	pipelineSys := getPipelineSystem(ctx)
	pipeline, err := Build(pipelineSys, n.Pipeline)
	if err != nil {
		return err
	}

	return pipeline.Process(ctx, entry)
}

func (n *RouterNode) Process(ctx context.Context, entry *logstorage.LogEntry) error {
	collectorSys := getCollectorSystem(ctx)
	logNotifier := getLogNotifier(ctx)

	key, err := collectorSys.Ingest(ctx, n.Stream, entry)
	if err == nil {
		logNotifier.Notify(n.Stream, string(key), *entry)
		metrics.IncStreamLogCounter(n.Stream)
	}

	return err
}
