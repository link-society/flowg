package pipelines

import (
	"context"
	"errors"
	"strings"
	"sync"

	"link-society.com/flowg/internal/app/metrics"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/langs/filtering"
	"link-society.com/flowg/internal/utils/langs/vrl"
)

// Node is one vertex of a compiled pipeline. Records flow through nodes from a
// source to its successors; each node transforms, filters, routes or forwards
// the record before passing it on.
type Node interface {
	// Init prepares the node (compile scripts, open connections) before any
	// record is processed.
	Init(ctx context.Context) error
	// Close releases whatever Init acquired.
	Close(ctx context.Context) error
	// Process handles one record and forwards it to the node's successors.
	Process(ctx context.Context, record *models.LogRecord) error
}

// SourceNode is a pipeline entrypoint; it simply forwards records to its
// successors.
type SourceNode struct {
	ID   string
	Next []Node
}

// TransformNode runs a VRL transformer, which may emit zero, one or many records
// for each input, forwarding each to its successors.
type TransformNode struct {
	ID          string
	Transformer string
	Next        []Node

	runner *vrl.ScriptRunner
}

// SwitchNode forwards a record to its successors only when the record matches
// its filtering condition.
type SwitchNode struct {
	ID        string
	Condition string
	Next      []Node

	runner filtering.Filter
}

// PipelineNode delegates processing to another named pipeline.
type PipelineNode struct {
	ID       string
	Pipeline string
}

// ForwardNode sends the record to an external destination through a forwarder.
type ForwardNode struct {
	ID        string
	Forwarder *models.ForwarderV2
}

// RouterNode persists the record into a log stream and notifies live
// subscribers; it is a terminal node.
type RouterNode struct {
	ID     string
	Stream string
}

var _ Node = (*SourceNode)(nil)
var _ Node = (*TransformNode)(nil)
var _ Node = (*SwitchNode)(nil)
var _ Node = (*PipelineNode)(nil)
var _ Node = (*ForwardNode)(nil)
var _ Node = (*RouterNode)(nil)

// sendRecordToNextNodes processes a record through every successor concurrently
// and joins their errors.
func sendRecordToNextNodes(ctx context.Context, next []Node, record *models.LogRecord) error {
	errC := make(chan error, len(next))
	wg := sync.WaitGroup{}

	for _, nextNode := range next {
		wg.Add(1)
		go func(nextNode Node) {
			defer wg.Done()
			err := nextNode.Process(ctx, record)
			if err != nil {
				errC <- err
			}
		}(nextNode)
	}

	wg.Wait()
	close(errC)

	var errs []error
	for err := range errC {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// traceNode appends a per-node trace entry when a tracer is active (dry runs);
// it is a no-op during normal processing.
func traceNode(
	ctx context.Context,
	nodeID string,
	err error,
	input map[string]string,
	output []map[string]string,
) {
	tracer := GetTracer(ctx)
	if tracer != nil {
		tracer.Trace = append(tracer.Trace, NodeTrace{
			NodeID: nodeID,
			Input:  input,
			Output: output,
			Error:  TraceError(err),
		})
	}
}

// isDryRun reports whether the context carries a tracer, in which case nodes
// skip their side effects (forwarding, ingestion) and only record traces.
func isDryRun(ctx context.Context) bool {
	return GetTracer(ctx) != nil
}

// MARK: source
func (n *SourceNode) Init(context.Context) error {
	return nil
}

func (n *SourceNode) Close(context.Context) error {
	return nil
}

func (n *SourceNode) Process(ctx context.Context, record *models.LogRecord) error {
	err := sendRecordToNextNodes(ctx, n.Next, record)
	traceNode(ctx, n.ID, err, record.Fields, []map[string]string{record.Fields})
	return err
}

// MARK: transform
func (n *TransformNode) Init(ctx context.Context) error {
	var err error
	n.runner, err = vrl.NewScriptRunner(n.Transformer)
	return err
}

func (n *TransformNode) Close(ctx context.Context) error {
	n.runner.Close()
	return nil
}

func (n *TransformNode) Process(ctx context.Context, record *models.LogRecord) error {
	output, err := n.runner.TransformLog(record.Fields)
	if err != nil {
		traceNode(ctx, n.ID, err, record.Fields, nil)
		return err
	}

	for _, event := range output {
		err := sendRecordToNextNodes(ctx, n.Next, &models.LogRecord{
			Timestamp: record.Timestamp,
			Fields:    event,
		})
		if err != nil {
			traceNode(ctx, n.ID, err, record.Fields, output)
			return err
		}
	}

	traceNode(ctx, n.ID, nil, record.Fields, output)
	return nil
}

// MARK: switch
func (n *SwitchNode) Init(ctx context.Context) error {
	var err error
	n.runner, err = filtering.Compile(n.Condition)
	return err
}

func (n *SwitchNode) Close(ctx context.Context) error {
	return nil
}

func (n *SwitchNode) Process(ctx context.Context, record *models.LogRecord) error {
	matches, err := n.runner.Evaluate(record)
	if err != nil {
		traceNode(ctx, n.ID, err, record.Fields, nil)
		return err
	}

	if matches {
		err = sendRecordToNextNodes(ctx, n.Next, record)
		traceNode(ctx, n.ID, err, record.Fields, []map[string]string{record.Fields})
		return err
	}

	traceNode(ctx, n.ID, nil, record.Fields, nil)
	return err
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
	traceNode(ctx, n.ID, err, record.Fields, nil)
	if err != nil {
		return err
	}

	if isDryRun(ctx) {
		return nil
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
	traceNode(ctx, n.ID, nil, record.Fields, nil)
	if isDryRun(ctx) {
		return nil
	}

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
	traceNode(ctx, n.ID, nil, record.Fields, nil)
	if isDryRun(ctx) {
		return nil
	}

	w := getWorker(ctx)

	key, err := w.logStorage.Ingest(ctx, n.Stream, record)
	if err == nil {
		err = w.logNotifier.Notify(ctx, n.Stream, strings.Join(key, ":"), *record)
	}
	if err == nil {
		metrics.IncStreamLogCounter(n.Stream)
	}

	return err
}
