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

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type syncRequest struct {
	remoteNodeID   string
	remoteEndpoint *url.URL
	lastSync       []clusterstate.NamespaceSyncState
}

type syncActor struct {
	actor.Actor
}

type syncWorker struct {
	logger *slog.Logger

	requestM actor.MailboxReceiver[*syncRequest]
	cookie   string

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

	url, err := url.JoinPath(remoteEndpoint.String(), syncState.Namespace)
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
	req.Header.Set(NODEID_HEADER_NAME, remoteNodeID)
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
}
