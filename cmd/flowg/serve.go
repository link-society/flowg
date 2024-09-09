package main

import (
	"log/slog"

	"context"
	"time"

	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/spf13/cobra"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/app/bootstrap"
	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/app/metrics"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"
)

type serveCommandOpts struct {
	bindAddress string
	authDir     string
	logDir      string
	configDir   string
	verbose     bool
}

func NewServeCommand() *cobra.Command {
	opts := &serveCommandOpts{}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start FlowG standalone server",
		PreRun: func(cmd *cobra.Command, args []string) {
			logging.Setup(opts.verbose)
			metrics.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			authDb, err := auth.NewDatabase(opts.authDir)
			if err != nil {
				slog.Error(
					"Failed to open auth database",
					"channel", "main",
					"path", opts.authDir,
					"error", err,
				)
				exitCode = 1
				return
			}
			defer func() {
				err := authDb.Close()
				if err != nil {
					slog.Error(
						"Failed to close auth database",
						"channel", "main",
						"path", opts.authDir,
						"error", err,
					)
					exitCode = 1
				}
			}()

			logDb, err := logstorage.NewStorage(opts.logDir)
			if err != nil {
				slog.Error(
					"Failed to open logs database",
					"channel", "main",
					"path", opts.logDir,
					"error", err,
				)
				exitCode = 1
				return
			}
			defer func() {
				err := logDb.Close()
				if err != nil {
					slog.Error(
						"Failed to close logs database",
						"channel", "main",
						"path", opts.logDir,
						"error", err,
					)
					exitCode = 1
				}
			}()

			logNotifier := lognotify.NewLogNotifier()
			logNotifier.Start()
			defer logNotifier.Stop()

			configStorage := config.NewStorage(opts.configDir)

			if err := bootstrap.DefaultRolesAndUsers(authDb); err != nil {
				slog.Error(
					"Failed to bootstrap default roles and users",
					"channel", "main",
					"error", err,
				)
				exitCode = 1
				return
			}

			if err := bootstrap.DefaultPipeline(configStorage); err != nil {
				slog.Error(
					"Failed to bootstrap default pipeline",
					"channel", "main",
					"error", err,
				)
				exitCode = 1
				return
			}

			apiHandler := api.NewHandler(authDb, logDb, configStorage, logNotifier)
			webHandler := web.NewHandler(authDb, logDb, configStorage)
			metricsHandler := promhttp.Handler()

			rootHandler := http.NewServeMux()
			rootHandler.HandleFunc(
				"/health",
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK\r\n"))
				},
			)
			rootHandler.Handle("/metrics", metricsHandler)

			rootHandler.Handle("/api/", apiHandler)
			rootHandler.Handle("/auth/", webHandler)
			rootHandler.Handle("/web/", webHandler)
			rootHandler.Handle("/static/", webHandler)

			rootHandler.HandleFunc(
				"GET /{$}",
				func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, "/web/", http.StatusPermanentRedirect)
				},
			)

			server := &http.Server{
				Addr:    opts.bindAddress,
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
						"bind", opts.bindAddress,
						"error", err,
					)
				}
			}()

			slog.Info(
				"Starting server",
				"channel", "main",
				"bind", opts.bindAddress,
			)
			err = server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				slog.Error(
					"Failed to start server",
					"channel", "main",
					"bind", opts.bindAddress,
					"error", err,
				)
				exitCode = 1
				return
			}
		},
	}

	cmd.Flags().StringVar(
		&opts.bindAddress,
		"bind",
		defaultBindAddress,
		"Address to bind the server to",
	)

	cmd.Flags().StringVar(
		&opts.authDir,
		"auth-dir",
		defaultAuthDir,
		"Path to the auth database directory",
	)

	cmd.Flags().StringVar(
		&opts.logDir,
		"log-dir",
		defaultLogDir,
		"Path to the log database directory",
	)
	cmd.MarkFlagDirname("log-dir")

	cmd.Flags().StringVar(
		&opts.configDir,
		"config-dir",
		defaultConfigDir,
		"Path to the config directory",
	)
	cmd.MarkFlagDirname("config-dir")

	cmd.Flags().BoolVar(
		&opts.verbose,
		"verbose",
		false,
		"Enable verbose logging",
	)

	return cmd
}
