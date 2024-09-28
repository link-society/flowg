package pipelines

import (
	"context"
	"errors"
	"sync"

	"link-society.com/flowg/internal/app/metrics"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/ffi/filterdsl"
	"link-society.com/flowg/internal/utils/ffi/vrl"
)

type Node interface {
	Process(ctx context.Context, record *models.LogRecord) error
}

type SourceNode struct {
	Next []Node
}

type TransformNode struct {
	TransformerName string
	Next            []Node
}

type SwitchNode struct {
	Condition filterdsl.Filter
	Next      []Node
}

type PipelineNode struct {
	Pipeline string
}

type AlertNode struct {
	Alert string
}

type RouterNode struct {
	Stream string
}

func (n *SourceNode) Process(ctx context.Context, record *models.LogRecord) error {
	errC := make(chan error, len(n.Next))
	wg := sync.WaitGroup{}

	for _, next := range n.Next {
		wg.Add(1)
		go func(next Node) {
			defer wg.Done()
			err := next.Process(ctx, record)
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

func (n *TransformNode) Process(ctx context.Context, record *models.LogRecord) error {
	configStorage := getConfigStorage(ctx)
	vrlScript, err := configStorage.ReadTransformer(ctx, n.TransformerName)
	if err != nil {
		return err
	}

	output, err := vrl.ProcessRecord(record.Fields, vrlScript)
	if err != nil {
		return err
	}

	errC := make(chan error, len(n.Next))
	wg := sync.WaitGroup{}

	for _, next := range n.Next {
		wg.Add(1)
		go func(next Node) {
			defer wg.Done()

			newEntry := &models.LogRecord{
				Timestamp: record.Timestamp,
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

func (n *SwitchNode) Process(ctx context.Context, record *models.LogRecord) error {
	if n.Condition.Evaluate(record) {
		errC := make(chan error, len(n.Next))
		wg := sync.WaitGroup{}

		for _, next := range n.Next {
			wg.Add(1)
			go func(next Node) {
				defer wg.Done()
				err := next.Process(ctx, record)
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

func (n *PipelineNode) Process(ctx context.Context, record *models.LogRecord) error {
	configStorage := getConfigStorage(ctx)
	pipeline, err := Build(ctx, configStorage, n.Pipeline)
	if err != nil {
		return err
	}

	return pipeline.Process(ctx, DIRECT_ENTRYPOINT, record)
}

func (n *AlertNode) Process(ctx context.Context, record *models.LogRecord) error {
	configStorage := getConfigStorage(ctx)
	alert, err := configStorage.ReadAlert(ctx, n.Alert)
	if err != nil {
		return err
	}

	return alert.Call(ctx, record)
}

func (n *RouterNode) Process(ctx context.Context, record *models.LogRecord) error {
	logStorage := getLogStorage(ctx)
	logNotifier := getLogNotifier(ctx)

	key, err := logStorage.Ingest(ctx, n.Stream, record)
	if err == nil {
		err = logNotifier.Notify(ctx, n.Stream, string(key), *record)
	}
	if err == nil {
		metrics.IncStreamLogCounter(n.Stream)
	}

	return err
}
