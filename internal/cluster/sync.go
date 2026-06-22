package cluster

import (
	"context"
	"log/slog"

	"fmt"
	"io"
	"strconv"

	"net/http"
	"net/url"

	"sync"
	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type syncRequest struct {
	remoteNodeID   string
	remoteEndpoint *url.URL
	lastSync       []clusterstate.NamespaceSyncState
	bootstrap      []string
}

type syncActor struct {
	actor.Actor
}

type syncWorker struct {
	logger *slog.Logger

	localNodeID        string
	requestM           actor.MailboxReceiver[*syncRequest]
	cookie             string
	bootstrapThreshold time.Duration

	storages            map[string]storage.Streamable
	clusterStateStorage clusterstate.Storage
}

var _ actor.Worker = (*syncWorker)(nil)

func (w *syncWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case req, ok := <-w.requestM.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		wg := sync.WaitGroup{}

		for _, syncState := range req.lastSync {
			wg.Add(1)

			go func(syncState clusterstate.NamespaceSyncState) {
				defer wg.Done()
				w.syncStorage(ctx, req.remoteNodeID, req.remoteEndpoint, syncState)
			}(syncState)
		}

		for _, namespace := range req.bootstrap {
			wg.Add(1)

			go func(namespace string) {
				defer wg.Done()
				w.bootstrapStorage(ctx, req.remoteNodeID, req.remoteEndpoint, namespace)
			}(namespace)
		}

		wg.Wait()

		return actor.WorkerContinue
	}
}

func (w *syncWorker) syncStorage(ctx context.Context, remoteNodeID string, remoteEndpoint *url.URL, syncState clusterstate.NamespaceSyncState) {
	logger := w.logger.With(
		slog.String("cluster.remote.node", remoteNodeID),
		slog.String("cluster.remote.endpoint", remoteEndpoint.String()),
		slog.String("cluster.replication.namespace", syncState.Namespace),
	)

	storage, ok := w.storages[syncState.Namespace]
	if !ok {
		logger.ErrorContext(ctx, "unknown namespace")
		return
	}

	url, err := url.JoinPath(remoteEndpoint.String(), "cluster", "sync", syncState.Namespace)
	if err != nil {
		logger.ErrorContext(ctx, "failed to build sync URL", slog.String("error", err.Error()))
		return
	}

	reader, writer := io.Pipe()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reader)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create sync request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set(COOKIE_HEADER_NAME, w.cookie)
	req.Header.Set(NODEID_HEADER_NAME, w.localNodeID)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Trailer", SINCE_HEADER_NAME)
	req.Trailer = http.Header{SINCE_HEADER_NAME: nil}

	go func() {
		defer writer.Close()

		newSyncTs, err := storage.Dump(ctx, writer, syncState.Since)
		if err != nil {
			logger.ErrorContext(ctx, "failed to dump storage", slog.String("error", err.Error()))
			return
		}

		req.Trailer.Set(SINCE_HEADER_NAME, strconv.FormatUint(newSyncTs, 10))
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "failed to send sync request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			message = fmt.Appendf(nil, "%d", resp.StatusCode)
		}

		logger.ErrorContext(ctx, "failed to sync storage", slog.String("status", string(message)))
		return
	}

	if err := w.clusterStateStorage.SetLiveness(ctx, syncState.Namespace, time.Now().UnixNano()); err != nil {
		logger.ErrorContext(ctx, "failed to update liveness", slog.String("error", err.Error()))
	}
}

func (w *syncWorker) bootstrapStorage(ctx context.Context, remoteNodeID string, remoteEndpoint *url.URL, namespace string) {
	logger := w.logger.With(
		slog.String("cluster.remote.node", remoteNodeID),
		slog.String("cluster.remote.endpoint", remoteEndpoint.String()),
		slog.String("cluster.replication.namespace", namespace),
	)

	store, ok := w.storages[namespace]
	if !ok {
		logger.ErrorContext(ctx, "unknown namespace")
		return
	}

	endpoint, err := url.JoinPath(remoteEndpoint.String(), "cluster", "sync", namespace)
	if err != nil {
		logger.ErrorContext(ctx, "failed to build sync URL", slog.String("error", err.Error()))
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create bootstrap request", slog.String("error", err.Error()))
		return
	}

	req.Header.Set(COOKIE_HEADER_NAME, w.cookie)
	req.Header.Set(NODEID_HEADER_NAME, w.localNodeID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "failed to send bootstrap request", slog.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			message = fmt.Appendf(nil, "%d", resp.StatusCode)
		}
		logger.WarnContext(ctx, "bootstrap source unavailable", slog.String("status", string(message)))
		return
	}

	if err := store.DropAll(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to drop namespace", slog.String("error", err.Error()))
		return
	}

	if err := store.Merge(ctx, resp.Body); err != nil {
		logger.ErrorContext(ctx, "failed to merge bootstrap snapshot", slog.String("error", err.Error()))
		return
	}

	if err := w.clusterStateStorage.ResetLocalState(ctx, namespace); err != nil {
		logger.ErrorContext(ctx, "failed to reset local state", slog.String("error", err.Error()))
		return
	}

	if err := w.clusterStateStorage.SetLiveness(ctx, namespace, time.Now().UnixNano()); err != nil {
		logger.ErrorContext(ctx, "failed to update liveness", slog.String("error", err.Error()))
		return
	}

	logger.InfoContext(ctx, "bootstrapped namespace from remote peer")
}
