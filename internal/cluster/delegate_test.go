package cluster

import (
	"testing"

	"log/slog"
	"time"

	"encoding/json"
	"io"
	"net/url"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

func newSyncMailbox(t *testing.T) actor.Mailbox[*syncRequest] {
	t.Helper()
	m := actor.NewMailbox[*syncRequest]()
	m.Start()
	t.Cleanup(m.Stop)
	return m
}

func newMergeDelegate(t *testing.T, syncM actor.MailboxSender[*syncRequest]) *delegate {
	t.Helper()
	stateStorage := newClusterStateStorage(t)
	// Mark every namespace as freshly synced so the staleness gate treats this
	// node as a healthy replication source and exercises the normal push path.
	now := time.Now().UnixNano()
	for _, ns := range []string{"auth", "config", "log"} {
		if err := stateStorage.SetLiveness(t.Context(), ns, now); err != nil {
			t.Fatalf("SetLiveness(%q): %v", ns, err)
		}
	}
	d := &delegate{
		logger:              slog.New(slog.NewTextHandler(io.Discard, nil)),
		localNodeID:         "node-local",
		localEndpoint:       &url.URL{Scheme: "http", Host: "127.0.0.1:9113"},
		endpoints:           newEndpointCache(),
		bootstrapThreshold:  time.Hour,
		clusterStateStorage: stateStorage,
		syncRequestM:        syncM,
		storages: map[string]storage.Streamable{
			// A non-zero latest version means there is data to replicate, so the
			// caught-up short-circuit does not kick in unless a test asks for it.
			"auth":   &fakeStreamable{latest: 1000},
			"config": &fakeStreamable{latest: 1000},
			"log":    &fakeStreamable{latest: 1000},
		},
	}
	d.endpoints.Set("node-remote", &url.URL{Scheme: "http", Host: "127.0.0.1:9114"})
	return d
}

func setLatest(t *testing.T, d *delegate, namespace string, latest uint64) {
	t.Helper()
	store, ok := d.storages[namespace].(*fakeStreamable)
	if !ok {
		t.Fatalf("namespace %q is not a *fakeStreamable", namespace)
	}
	store.latest = latest
}

func expectSyncRequest(t *testing.T, m actor.Mailbox[*syncRequest]) *syncRequest {
	t.Helper()
	select {
	case req := <-m.ReceiveC():
		return req
	case <-time.After(time.Second):
		t.Fatal("expected a sync request, got none")
		return nil
	}
}

func expectNoSyncRequest(t *testing.T, m actor.Mailbox[*syncRequest]) {
	t.Helper()
	select {
	case req := <-m.ReceiveC():
		t.Fatalf("expected no sync request, got %+v", req)
	case <-time.After(100 * time.Millisecond):
	}
}

func sinceByNamespace(states []clusterstate.NamespaceSyncState) map[string]uint64 {
	out := make(map[string]uint64, len(states))
	for _, st := range states {
		out[st.Namespace] = st.Since
	}
	return out
}

// TestMergeRemoteStateFirstContact verifies that when a peer has never synced
// from us (no LastSync entry for the local node) we still trigger a sync, pushing
// every namespace from scratch instead of logging an error.
func TestMergeRemoteStateFirstContact(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)

	req := expectSyncRequest(t, m)
	if req.remoteNodeID != "node-remote" {
		t.Fatalf("remoteNodeID: got %q want node-remote", req.remoteNodeID)
	}

	since := sinceByNamespace(req.lastSync)
	for _, ns := range []string{"auth", "config", "log"} {
		v, ok := since[ns]
		if !ok {
			t.Fatalf("expected namespace %q in sync request, got %v", ns, since)
		}
		if v != 0 {
			t.Fatalf("expected since=0 for %q on first contact, got %d", ns, v)
		}
	}
}

// TestMergeRemoteStateUsesKnownWatermarks verifies that the watermarks the peer
// reports for OUR data drive incremental sync, that watermarks for other source
// nodes are ignored, and that namespaces the peer has never seen default to 0.
func TestMergeRemoteStateUsesKnownWatermarks(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	remote := clusterstate.NodeState{
		NodeID: "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{
			"node-local": {
				{Namespace: "auth", Since: 42},
				{Namespace: "config", Since: 7},
			},
			// Watermarks for data the peer received from a different node must
			// not leak into our push.
			"node-other": {
				{Namespace: "auth", Since: 999},
			},
		},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)

	req := expectSyncRequest(t, m)
	since := sinceByNamespace(req.lastSync)
	if since["auth"] != 42 {
		t.Fatalf("auth since: got %d want 42", since["auth"])
	}
	if since["config"] != 7 {
		t.Fatalf("config since: got %d want 7", since["config"])
	}
	if v, ok := since["log"]; !ok || v != 0 {
		t.Fatalf("log since: got %d (present=%v) want 0", v, ok)
	}
}

// TestMergeRemoteStateIgnoresOwnState verifies we never sync against our own
// gossiped state.
func TestMergeRemoteStateIgnoresOwnState(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	remote := clusterstate.NodeState{NodeID: "node-local"}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)
	expectNoSyncRequest(t, m)
}

// TestMergeRemoteStateUnknownEndpoint verifies that a peer whose endpoint we have
// not learned yet does not trigger a sync.
func TestMergeRemoteStateUnknownEndpoint(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)
	d.endpoints.Delete("node-remote")

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)
	expectNoSyncRequest(t, m)
}

// TestMergeRemoteStateSkipsCaughtUpNamespaces verifies that namespaces for which
// the peer's watermark already covers our latest version are not included in the
// sync request, while lagging namespaces still are.
func TestMergeRemoteStateSkipsCaughtUpNamespaces(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	setLatest(t, d, "auth", 50)
	setLatest(t, d, "config", 50)
	setLatest(t, d, "log", 50)

	remote := clusterstate.NodeState{
		NodeID: "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{
			"node-local": {
				{Namespace: "auth", Since: 50},   // caught up exactly -> skip
				{Namespace: "config", Since: 80}, // ahead of us -> skip
				{Namespace: "log", Since: 10},    // behind -> sync
			},
		},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)

	req := expectSyncRequest(t, m)
	since := sinceByNamespace(req.lastSync)
	if len(since) != 1 {
		t.Fatalf("expected only the lagging namespace, got %v", since)
	}
	if v, ok := since["log"]; !ok || v != 10 {
		t.Fatalf("log since: got %d (present=%v) want 10", v, ok)
	}
}

// TestMergeRemoteStateNoRequestWhenAllCaughtUp verifies that no sync request is
// emitted at all when the peer already holds everything we have.
func TestMergeRemoteStateNoRequestWhenAllCaughtUp(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	setLatest(t, d, "auth", 50)
	setLatest(t, d, "config", 50)
	setLatest(t, d, "log", 50)

	remote := clusterstate.NodeState{
		NodeID: "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{
			"node-local": {
				{Namespace: "auth", Since: 50},
				{Namespace: "config", Since: 50},
				{Namespace: "log", Since: 50},
			},
		},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)
	expectNoSyncRequest(t, m)
}

// TestMergeRemoteStateEmptyStoreNoRequest verifies that a node with no data yet
// (latest version 0) does not generate any traffic on first contact.
func TestMergeRemoteStateEmptyStoreNoRequest(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	setLatest(t, d, "auth", 0)
	setLatest(t, d, "config", 0)
	setLatest(t, d, "log", 0)

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)
	expectNoSyncRequest(t, m)
}

// setLiveness rewrites the persisted liveness marker for a namespace so a test
// can place the local node anywhere on the staleness timeline.
func setLiveness(t *testing.T, d *delegate, namespace string, at time.Time) {
	t.Helper()
	if err := d.clusterStateStorage.SetLiveness(t.Context(), namespace, at.UnixNano()); err != nil {
		t.Fatalf("SetLiveness(%q): %v", namespace, err)
	}
}

func bootstrapSet(req *syncRequest) map[string]struct{} {
	out := make(map[string]struct{}, len(req.bootstrap))
	for _, ns := range req.bootstrap {
		out[ns] = struct{}{}
	}
	return out
}

// TestMergeRemoteStateStaleBootstrapsFromHealthyPeer verifies the push gate: a
// namespace we have not synced within the threshold is never pushed (which could
// resurrect GC'd tombstones) and is instead scheduled for a destructive
// bootstrap, but only because the remote advertises it as a healthy source.
func TestMergeRemoteStateStaleBootstrapsFromHealthyPeer(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	// auth is stale (out of contact > threshold) but config/log stay fresh.
	setLiveness(t, d, "auth", time.Now().Add(-2*time.Hour))

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
		Healthy:  map[string]bool{"auth": true, "config": true, "log": true},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)

	req := expectSyncRequest(t, m)

	boot := bootstrapSet(req)
	if _, ok := boot["auth"]; !ok {
		t.Fatalf("expected auth to be bootstrapped, got %v", req.bootstrap)
	}

	since := sinceByNamespace(req.lastSync)
	if _, ok := since["auth"]; ok {
		t.Fatalf("stale auth must not be pushed, got lastSync %v", since)
	}
	for _, ns := range []string{"config", "log"} {
		if _, ok := since[ns]; !ok {
			t.Fatalf("fresh namespace %q should still push, got %v", ns, since)
		}
	}
}

// TestMergeRemoteStateStaleNoHealthySourceWaits verifies that a stale namespace
// is neither pushed nor bootstrapped when no reachable peer advertises it as
// healthy and the node has not yet crossed the self-heal horizon.
func TestMergeRemoteStateStaleNoHealthySourceWaits(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	// Stale (> threshold) but not past the 2*threshold self-heal horizon.
	setLiveness(t, d, "auth", time.Now().Add(-90*time.Minute))

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
		Healthy:  map[string]bool{"auth": false, "config": true, "log": true},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)

	// config/log still push, so a request is emitted, but auth appears nowhere.
	req := expectSyncRequest(t, m)
	if len(req.bootstrap) != 0 {
		t.Fatalf("expected no bootstrap, got %v", req.bootstrap)
	}
	if _, ok := sinceByNamespace(req.lastSync)["auth"]; ok {
		t.Fatalf("stale auth without healthy source must not push")
	}

	// Liveness stays old: the node has not been declared healthy again.
	got, err := d.clusterStateStorage.GetLiveness(t.Context(), "auth")
	if err != nil {
		t.Fatalf("GetLiveness: %v", err)
	}
	if time.Since(time.Unix(0, got)) < time.Hour {
		t.Fatalf("liveness should not have been refreshed, got %v ago", time.Since(time.Unix(0, got)))
	}
}

// TestMergeRemoteStateStaleSelfHeals verifies that a node out of contact beyond
// the self-heal horizon (2*threshold) resumes on its own when it can reach peers
// but none advertise a healthy source — the cold-restart escape hatch.
func TestMergeRemoteStateStaleSelfHeals(t *testing.T) {
	t.Parallel()
	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)

	setLiveness(t, d, "auth", time.Now().Add(-3*time.Hour)) // > 2*threshold (2h)

	remote := clusterstate.NodeState{
		NodeID:   "node-remote",
		LastSync: map[string][]clusterstate.NamespaceSyncState{},
		Healthy:  map[string]bool{"auth": false, "config": true, "log": true},
	}
	buf, err := json.Marshal(remote)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	d.MergeRemoteState(buf, false)
	// Drain whatever request config/log produce; auth must not appear in it.
	req := expectSyncRequest(t, m)
	if _, ok := sinceByNamespace(req.lastSync)["auth"]; ok {
		t.Fatalf("auth must not push during the self-heal round")
	}
	if len(req.bootstrap) != 0 {
		t.Fatalf("auth must not bootstrap without a healthy source, got %v", req.bootstrap)
	}

	// Liveness has been refreshed: the node considers itself healthy again.
	got, err := d.clusterStateStorage.GetLiveness(t.Context(), "auth")
	if err != nil {
		t.Fatalf("GetLiveness: %v", err)
	}
	if time.Since(time.Unix(0, got)) > time.Minute {
		t.Fatalf("expected liveness refreshed to ~now, got %v ago", time.Since(time.Unix(0, got)))
	}
}

func newClusterStateStorage(t *testing.T) clusterstate.Storage {
	t.Helper()

	opts := clusterstate.DefaultOptions()
	opts.InMemory = true

	var s clusterstate.Storage
	app := fxtest.New(t, clusterstate.NewStorage(opts), fx.Populate(&s), fx.NopLogger)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	return s
}

// TestWatermarkRoundTripDrivesSync ties the cluster-state storage and the
// delegate together to validate the full keying contract: a watermark recorded
// on node-b for data sourced from node-a (exactly what handleSync stores from the
// NODEID header) is reported by node-b's gossiped state under "node-a", and
// node-a resumes pushing that namespace from the recorded version.
func TestWatermarkRoundTripDrivesSync(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	stateStorage := newClusterStateStorage(t)
	if err := stateStorage.UpdateLocalState(ctx, "node-a", "config", 100); err != nil {
		t.Fatalf("UpdateLocalState: %v", err)
	}

	nodeBState, err := stateStorage.FetchLocalState(ctx, "node-b", nil)
	if err != nil {
		t.Fatalf("FetchLocalState: %v", err)
	}
	buf, err := json.Marshal(nodeBState)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	m := newSyncMailbox(t)
	d := newMergeDelegate(t, m)
	d.localNodeID = "node-a"
	d.endpoints.Set("node-b", &url.URL{Scheme: "http", Host: "127.0.0.1:9114"})

	d.MergeRemoteState(buf, false)

	req := expectSyncRequest(t, m)
	if req.remoteNodeID != "node-b" {
		t.Fatalf("remoteNodeID: got %q want node-b", req.remoteNodeID)
	}

	since := sinceByNamespace(req.lastSync)
	if since["config"] != 100 {
		t.Fatalf("config since: got %d want 100", since["config"])
	}
}
