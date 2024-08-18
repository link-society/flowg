package main

import (
	"log/slog"

	"context"
	"time"

	"os"
	"os/signal"
	"syscall"

	"net/http"

	"link-society.com/flowg/internal/logging"
	"link-society.com/flowg/internal/pipelines"
	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"
)

func main() {
	logging.Setup(*verbose)
	os.Exit(run())
}

func run() int {
	db, err := storage.NewStorage(*dbPath)
	if err != nil {
		slog.Error(
			"Failed to open database",
			"channel", "main",
			"dbpath", *dbPath,
			"error", err,
		)
		return 1
	}
	defer db.Close()

	pipelinesManager := pipelines.NewManager(db, *configDir)

	apiHandler := api.NewHandler(db, pipelinesManager)
	webHandler := web.NewHandler(db, pipelinesManager)

	rootHandler := http.NewServeMux()
	rootHandler.Handle("/api/", apiHandler)
	rootHandler.Handle("/web/", webHandler)
	rootHandler.Handle("/static/", webHandler)

	rootHandler.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusPermanentRedirect)
	})

	server := &http.Server{
		Addr:    *bindAddress,
		Handler: logging.NewMiddleware(rootHandler),
	}

	go func() {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		<-sigC

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error(
				"Failed to shutdown server",
				"channel", "main",
				"bind", *bindAddress,
				"error", err,
			)
		}
	}()

	slog.Info(
		"Starting server",
		"channel", "main",
		"bind", *bindAddress,
	)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error(
			"Failed to start server",
			"channel", "main",
			"bind", *bindAddress,
			"error", err,
		)
		return 1
	}

	return 0
}
