package cluster

import (
	"log/slog"

	"fmt"
	"strconv"

	"io"
	"net/http"
	"net/url"

	"sync"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage"
)

type syncPool struct {
	logger *slog.Logger

	nodeID string
	cookie string

	authStorage   storage.Streamable
	configStorage storage.Streamable
	logStorage    storage.Streamable

	workers map[string]*syncActor
}

type syncActor struct {
	actor.Actor
	mbox actor.MailboxSender[nodeSyncState]
}

type syncWorker struct {
	pool     *syncPool
	nodeID   string
	endpoint *url.URL
	mbox     actor.MailboxReceiver[nodeSyncState]
}

var _ actor.Actor = (*syncActor)(nil)
var _ actor.Worker = (*syncWorker)(nil)

func (p *syncPool) AddWorker(nodeID string, endpoint *url.URL) {
	if _, exists := p.workers[nodeID]; exists {
		return
	}

	mbox := actor.NewMailbox[nodeSyncState]()

	worker := &syncWorker{
		pool:     p,
		nodeID:   nodeID,
		endpoint: endpoint,
		mbox:     mbox,
	}

	syncer := &syncActor{
		Actor: actor.Combine(mbox, actor.New(worker)).Build(),
		mbox:  mbox,
	}
	syncer.Start()

	p.workers[nodeID] = syncer
}

func (p *syncPool) RemoveWorker(nodeID string) {
	if worker, exists := p.workers[nodeID]; exists {
		worker.Stop()
		delete(p.workers, nodeID)
	}
}

func (p *syncPool) RemoveAll() {
	for nodeID, worker := range p.workers {
		worker.Stop()
		delete(p.workers, nodeID)
	}
}

func (w *syncWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.mbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		wg := &sync.WaitGroup{}
		wg.Add(3)

		go func() {
			defer wg.Done()
			w.syncStorage(ctx, w.pool.authStorage, msg.Auth, "auth")
		}()

		go func() {
			defer wg.Done()
			w.syncStorage(ctx, w.pool.configStorage, msg.Config, "config")
		}()

		go func() {
			defer wg.Done()
			w.syncStorage(ctx, w.pool.logStorage, msg.Log, "log")
		}()

		wg.Wait()

		return actor.WorkerContinue
	}
}

func (w *syncWorker) syncStorage(
	ctx actor.Context,
	s storage.Streamable,
	lastSync uint64,
	dbType string,
) {
	logger := w.pool.logger.With(
		slog.String("cluster.remote.node", w.nodeID),
		slog.String("cluster.remote.endpoint", w.endpoint.String()),
		slog.String("cluster.replication.storage", dbType),
	)

	url, err := url.JoinPath(w.endpoint.String(), "/cluster/sync", dbType)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to build sync url",
			slog.String("error", err.Error()),
		)
		return
	}

	reader, writer := io.Pipe()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		reader,
	)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to create sync request",
			slog.String("error", err.Error()),
		)
		return
	}

	req.Header.Set(COOKIE_HEADER_NAME, w.pool.cookie)
	req.Header.Set(NODEID_HEADER_NAME, w.pool.nodeID)
	req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set("Trailer", SINCE_HEADER_NAME)
	req.Trailer = http.Header{
		SINCE_HEADER_NAME: nil,
	}

	go func() {
		defer writer.Close()

		newSyncTs, err := s.Dump(ctx, writer, lastSync)
		if err != nil {
			logger.ErrorContext(
				ctx,
				"failed to dump storage",
				slog.String("error", err.Error()),
			)
			return
		}

		req.Trailer.Set(SINCE_HEADER_NAME, strconv.FormatUint(newSyncTs, 10))
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to send sync request",
			slog.String("error", err.Error()),
		)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			message = fmt.Appendf(nil, "%d", resp.StatusCode)
		}

		logger.ErrorContext(
			ctx,
			"failed to sync storage",
			slog.String("error", string(message)),
		)
		return
	}
}
