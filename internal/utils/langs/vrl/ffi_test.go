package vrl_test

import (
	"testing"

	"fmt"
	"reflect"
	"sync"

	"go.link-society.com/flowg/internal/utils/langs/vrl"
)

// TransformLog is called concurrently by the pipeline engine, one goroutine
// per record; each caller must get back the result of its own event.
func TestTransformLog_Concurrent(t *testing.T) {
	runner, err := vrl.NewScriptRunner(`.echo = .id`)
	if err != nil {
		t.Fatalf("NewScriptRunner() failed: %v", err)
	}
	defer runner.Close()

	const workers, iterations = 8, 500

	var wg sync.WaitGroup
	errs := make(chan error, workers)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()

			for i := range iterations {
				id := fmt.Sprintf("%d-%d", w, i)

				output, err := runner.TransformLog(map[string]string{"id": id})
				if err != nil {
					errs <- err
					return
				}

				if len(output) != 1 || output[0]["echo"] != id {
					errs <- fmt.Errorf("TransformLog(%q) = %v", id, output)
					return
				}
			}
		}(w)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

func TestTransformLog(t *testing.T) {
	runner, err := vrl.NewScriptRunner(`
		.foo = "bar"
		.bar.baz = [1, 2, 3, "a"]
	`)
	if err != nil {
		t.Fatalf("NewScriptRunner() failed: %v", err)
	}
	defer runner.Close()

	input := map[string]string{}
	output, err := runner.TransformLog(input)
	if err != nil {
		t.Errorf("TransformLog() failed: %v", err)
	}

	expected := []map[string]string{{
		"foo":       "bar",
		"bar.baz.0": "1",
		"bar.baz.1": "2",
		"bar.baz.2": "3",
		"bar.baz.3": "a",
	}}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("TransformLog() = %v, want %v", output, expected)
	}
}

func TestTransformLog_EmptyScript(t *testing.T) {
	runner, err := vrl.NewScriptRunner(``)
	if err != nil {
		t.Fatalf("NewScriptRunner() failed: %v", err)
	}
	defer runner.Close()

	input := map[string]string{
		"foo": "bar",
	}

	output, err := runner.TransformLog(input)
	if err != nil {
		t.Errorf("TransformLog() failed: %v", err)
	}

	if !reflect.DeepEqual(output, []map[string]string{input}) {
		t.Errorf("TransformLog() = %v, want %v", output, input)
	}
}
