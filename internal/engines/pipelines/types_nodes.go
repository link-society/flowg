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
	Init(ctx context.Context) error
	Close(ctx context.Context) error
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

type ForwardNode struct {
	Forwarder *models.ForwarderV2
}

type RouterNode struct {
	Stream string
}

var _ Node = (*SourceNode)(nil)
var _ Node = (*TransformNode)(nil)
var _ Node = (*SwitchNode)(nil)
var _ Node = (*PipelineNode)(nil)
var _ Node = (*ForwardNode)(nil)
var _ Node = (*RouterNode)(nil)

// MARK: source
func (n *SourceNode) Init(context.Context) error {
	return nil
}

func (n *SourceNode) Close(context.Context) error {
	return nil
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

// MARK: transform
func (n *TransformNode) Init(ctx context.Context) error {
	return nil
}

func (n *TransformNode) Close(ctx context.Context) error {
	return nil
}

func (n *TransformNode) Process(ctx context.Context, record *models.LogRecord) error {
	w := getWorker(ctx)
	vrlScript, err := w.configStorage.ReadTransformer(ctx, n.TransformerName)
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

// MARK: switch
func (n *SwitchNode) Init(ctx context.Context) error {
	return nil
}

func (n *SwitchNode) Close(ctx context.Context) error {
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

// MARK: pipeline
func (n *PipelineNode) Init(ctx context.Context) error {
	return nil
}

func (n *PipelineNode) Close(ctx context.Context) error {
	return nil
}

func (n *PipelineNode) Process(ctx context.Context, record *models.LogRecord) error {
	w := getWorker(ctx)
	pipeline, err := w.getOrBuildPipeline(ctx, n.Pipeline)
	if err != nil {
		return err
	}

	return pipeline.Process(ctx, DIRECT_ENTRYPOINT, record)
}

// MARK: forward
func (n *ForwardNode) Init(ctx context.Context) error {
	return n.Forwarder.Init(ctx)
}

func (n *ForwardNode) Close(ctx context.Context) error {
	return n.Forwarder.Close(ctx)
}

func (n *ForwardNode) Process(ctx context.Context, record *models.LogRecord) error {
	return n.Forwarder.Call(ctx, record)
}

// MARK: router
func (n *RouterNode) Init(ctx context.Context) error {
	return nil
}

func (n *RouterNode) Close(ctx context.Context) error {
	return nil
}

func (n *RouterNode) Process(ctx context.Context, record *models.LogRecord) error {
	w := getWorker(ctx)

	key, err := w.logStorage.Ingest(ctx, n.Stream, record)
	if err == nil {
		err = w.logNotifier.Notify(ctx, n.Stream, string(key), *record)
	}
	if err == nil {
		metrics.IncStreamLogCounter(n.Stream)
	}

	return err
}
