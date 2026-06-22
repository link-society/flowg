package config

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/utils/hlc"
)

// newConfigNode spins up a standalone config storage backed by an in-memory
// badger instance, with an HLC clock whose NodeID is the given node id. Each
// node is fully independent — there is no cross-node communication other than
// the explicit anti-entropy exchanges performed by the tests.
func newConfigNode(t *testing.T, nodeID string) Storage {
	t.Helper()

	opts := DefaultOptions()
	opts.InMemory = true

	var s Storage

	app := fxtest.New(
		t,
		fx.Provide(func() *hlc.Clock { return hlc.NewClock(nodeID) }),
		changefeed.NewNotifier(),
		NewStorage(opts),
		fx.Populate(&s),
		fx.NopLogger,
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return s
}

func mustWriteTransformer(t *testing.T, ctx context.Context, s Storage, name, content string) {
	t.Helper()
	if err := s.WriteTransformer(ctx, name, content); err != nil {
		t.Fatalf("write transformer %q: %v", name, err)
	}
}

func mustDeleteTransformer(t *testing.T, ctx context.Context, s Storage, name string) {
	t.Helper()
	if err := s.DeleteTransformer(ctx, name); err != nil {
		t.Fatalf("delete transformer %q: %v", name, err)
	}
}

func dumpNode(t *testing.T, ctx context.Context, s Storage) []byte {
	t.Helper()
	var buf bytes.Buffer
	if _, err := s.Dump(ctx, &buf, 0); err != nil {
		t.Fatalf("dump: %v", err)
	}
	return buf.Bytes()
}

func snapshotTransformers(t *testing.T, ctx context.Context, s Storage) map[string]string {
	t.Helper()
	names, err := s.ListTransformers(ctx)
	if err != nil {
		t.Fatalf("list transformers: %v", err)
	}
	out := make(map[string]string, len(names))
	for _, name := range names {
		content, err := s.ReadTransformer(ctx, name)
		if err != nil {
			t.Fatalf("read transformer %q: %v", name, err)
		}
		out[name] = content
	}
	return out
}

// captureDumps snapshots every node's full state at a single instant — the
// moment a healed network partition begins anti-entropy. Dumps are captured
// up-front so a later round can replay the exact same payloads (badger's
// streaming backup is expensive, so callers reuse these snapshots rather than
// re-dumping).
func captureDumps(t *testing.T, ctx context.Context, nodes []Storage) [][]byte {
	t.Helper()

	dumps := make([][]byte, len(nodes))
	for i, n := range nodes {
		dumps[i] = dumpNode(t, ctx, n)
	}
	return dumps
}

// mergeRound models a healed network partition: every node merges every other
// node's snapshot. Because the LWW merge is commutative, associative and
// idempotent, a single round over full snapshots is sufficient for every node
// to converge on the global winner set, and replaying the same snapshots is a
// no-op.
func mergeRound(t *testing.T, ctx context.Context, nodes []Storage, dumps [][]byte) {
	t.Helper()

	for i, n := range nodes {
		for j := range nodes {
			if i == j {
				continue
			}
			if err := n.Merge(ctx, bytes.NewReader(dumps[j])); err != nil {
				t.Fatalf("node %d merge node %d: %v", i, j, err)
			}
		}
	}
}

// TestMultiNodePartitionRejoinConvergence simulates three nodes accepting
// independent, conflicting writes while partitioned, then rejoining and running
// anti-entropy. After rejoin, all nodes must hold byte-identical state, LWW must
// pick the latest writer for conflicting keys, node-local writes must propagate
// everywhere, and a delete on the latest writer must win as a tombstone.
func TestMultiNodePartitionRejoinConvergence(t *testing.T) {
	ctx := t.Context()

	node0 := newConfigNode(t, "node0")
	node1 := newConfigNode(t, "node1")
	node2 := newConfigNode(t, "node2")
	nodes := []Storage{node0, node1, node2}

	// --- Partition phase: independent, conflicting writes ---

	// Same key written on all three. node2 writes last (later HLC wall time, and
	// a winning NodeID tiebreak if wall times collide) so node2 must win.
	mustWriteTransformer(t, ctx, node0, "shared", "from-node0")
	mustWriteTransformer(t, ctx, node1, "shared", "from-node1")
	mustWriteTransformer(t, ctx, node2, "shared", "from-node2")

	// Node-local keys that must propagate to every node.
	mustWriteTransformer(t, ctx, node0, "only0", "owned-by-0")
	mustWriteTransformer(t, ctx, node1, "only1", "owned-by-1")

	// Key created everywhere, then deleted on the latest writer: the tombstone
	// must win over every prior create.
	mustWriteTransformer(t, ctx, node0, "victim", "v0")
	mustWriteTransformer(t, ctx, node1, "victim", "v1")
	mustWriteTransformer(t, ctx, node2, "victim", "v2")
	mustDeleteTransformer(t, ctx, node2, "victim")

	// --- Rejoin phase: anti-entropy across all nodes ---
	// Snapshot the partition-era state once, then merge those snapshots into
	// every node. The same snapshots are replayed below to assert idempotency.
	dumps := captureDumps(t, ctx, nodes)
	mergeRound(t, ctx, nodes, dumps)

	// --- Convergence: every node must hold identical state ---
	want := snapshotTransformers(t, ctx, node0)
	for i := 1; i < len(nodes); i++ {
		got := snapshotTransformers(t, ctx, nodes[i])
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("node %d diverged from node 0:\n node0=%#v\n node%d=%#v", i, want, i, got)
		}
	}

	// --- LWW correctness ---
	if want["shared"] != "from-node2" {
		t.Fatalf("expected shared=from-node2 (latest writer wins), got %q", want["shared"])
	}
	if want["only0"] != "owned-by-0" {
		t.Fatalf("expected only0=owned-by-0 to propagate, got %q", want["only0"])
	}
	if want["only1"] != "owned-by-1" {
		t.Fatalf("expected only1=owned-by-1 to propagate, got %q", want["only1"])
	}
	if content, ok := want["victim"]; ok {
		t.Fatalf("expected victim to be deleted everywhere, but it is present: %q", content)
	}

	// --- Stability: replaying the same anti-entropy payloads must be a no-op ---
	// Re-merging the original (now stale) partition snapshots must not regress
	// any node: LWW keeps the winners already converged on above.
	mergeRound(t, ctx, nodes, dumps)
	for i := range nodes {
		got := snapshotTransformers(t, ctx, nodes[i])
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("node %d changed after an idempotent round:\n want=%#v\n got=%#v", i, want, got)
		}
	}
}
