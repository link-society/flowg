package clusterstate

import (
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func newTestStorage(t *testing.T) Storage {
	t.Helper()

	opts := DefaultOptions()
	opts.Directory = t.TempDir()

	var s Storage
	app := fxtest.New(t, NewStorage(opts), fx.Populate(&s), fx.NopLogger)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return s
}

func mustUpdate(t *testing.T, s Storage, nodeID, namespace string, since uint64) {
	t.Helper()
	if err := s.UpdateLocalState(t.Context(), nodeID, namespace, since); err != nil {
		t.Fatalf("UpdateLocalState(%s, %s, %d): %v", nodeID, namespace, since, err)
	}
}

func collect(states []NamespaceSyncState) map[string]uint64 {
	out := make(map[string]uint64, len(states))
	for _, st := range states {
		out[st.Namespace] = st.Since
	}
	return out
}

// TestFetchLocalStateGroupsBySourceNode verifies that watermarks recorded for
// several source nodes/namespaces are reported grouped by SOURCE node id (not by
// namespace), which is what MergeRemoteState relies on.
func TestFetchLocalStateGroupsBySourceNode(t *testing.T) {
	ctx := t.Context()
	s := newTestStorage(t)

	mustUpdate(t, s, "node-a", "auth", 10)
	mustUpdate(t, s, "node-a", "config", 20)
	mustUpdate(t, s, "node-b", "log", 30)
	// Node ids may contain colons; the namespace never does, so parsing must
	// split on the LAST colon.
	mustUpdate(t, s, "node:weird", "auth", 5)

	state, err := s.FetchLocalState(ctx, "node-local", nil)
	if err != nil {
		t.Fatalf("FetchLocalState: %v", err)
	}

	if state.NodeID != "node-local" {
		t.Fatalf("NodeID: got %q want node-local", state.NodeID)
	}

	nodeA := collect(state.LastSync["node-a"])
	if len(nodeA) != 2 || nodeA["auth"] != 10 || nodeA["config"] != 20 {
		t.Fatalf("node-a watermarks: got %v want {auth:10, config:20}", nodeA)
	}

	nodeB := collect(state.LastSync["node-b"])
	if len(nodeB) != 1 || nodeB["log"] != 30 {
		t.Fatalf("node-b watermarks: got %v want {log:30}", nodeB)
	}

	weird := collect(state.LastSync["node:weird"])
	if len(weird) != 1 || weird["auth"] != 5 {
		t.Fatalf("node:weird watermarks: got %v want {auth:5}", weird)
	}
}

// TestUpdateLocalStateOverwrites verifies the watermark for a given
// (source node, namespace) advances in place rather than accumulating.
func TestUpdateLocalStateOverwrites(t *testing.T) {
	ctx := t.Context()
	s := newTestStorage(t)

	mustUpdate(t, s, "node-a", "auth", 10)
	mustUpdate(t, s, "node-a", "auth", 25)

	state, err := s.FetchLocalState(ctx, "node-local", nil)
	if err != nil {
		t.Fatalf("FetchLocalState: %v", err)
	}

	nodeA := collect(state.LastSync["node-a"])
	if len(nodeA) != 1 || nodeA["auth"] != 25 {
		t.Fatalf("node-a watermarks: got %v want {auth:25}", nodeA)
	}
}
