package syslog

import (
	"fmt"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
	gosyslog "gopkg.in/mcuadros/go-syslog.v2"
)

func NewServer(bindAddress string) actor.Actor {
	worker := &worker{
		bindAddress: bindAddress,
	}

	return actor.New(
		worker,
		actor.OptOnStart(worker.onStart),
		actor.OptOnStop(worker.onStop),
	)
}

type worker struct {
	bindAddress string

	channel gosyslog.LogPartsChannel
	handler *gosyslog.ChannelHandler
	server  *gosyslog.Server
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case logParts, ok := <-w.channel:
		if !ok {
			return actor.WorkerEnd
		}

		slog.DebugContext(
			ctx,
			"Received syslog message",
			"channel", "syslog",
			"logrecord", fmt.Sprintf("%v", logParts),
		)
		return actor.WorkerContinue
	}
}

func (w *worker) onStart(ctx actor.Context) {
	w.channel = make(gosyslog.LogPartsChannel)
	w.handler = gosyslog.NewChannelHandler(w.channel)

	w.server = gosyslog.NewServer()
	w.server.SetFormat(gosyslog.Automatic)
	w.server.SetHandler(w.handler)
	w.server.ListenUDP(w.bindAddress)

	w.server.Boot()
}

func (w *worker) onStop() {
	w.server.Kill()
	w.server.Wait()
}
