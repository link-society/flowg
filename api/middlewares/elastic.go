package middlewares

import (
	"log/slog"

	"encoding/base64"
	"encoding/json"

	"fmt"
	"slices"
	"strconv"
	"strings"

	"net/http"

	"go.uber.org/fx"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ElasticDeps lists the dependencies of [newElasticHandler]: the backends it
// uses to authenticate callers and turn the documents they index into pipeline
// runs.
type ElasticDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	ConfigStorage  storage.ConfigStorage
	PipelineRunner pipelines.Runner
}

func init() {
	routing.RegisterMiddleware(
		NewElasticHandler,
		"/api/v1/middlewares/elastic/",
	)
}

// elasticProduct advertises the handler as Elasticsearch, which the official
// clients probe for before they agree to talk to a server.
func elasticProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		next.ServeHTTP(w, r)
	})
}

// elasticAuth resolves the HTTP Basic credentials Elasticsearch clients send
// into a FlowG user, so downstream handlers can enforce the user's permissions.
func elasticAuth(deps ElasticDeps, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		basicAuthData := authHeader[len("Basic "):]
		decoded, err := base64.StdEncoding.DecodeString(basicAuthData)
		if err != nil {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}

		username, password := parts[0], parts[1]
		ok, err := deps.AuthStorage.VerifyUserPassword(r.Context(), username, password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := deps.AuthStorage.FetchUser(r.Context(), username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ctx := auth.ContextWithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewElasticHandler serves the subset of the Elasticsearch API that log
// shippers exercise: checking that an index (a FlowG pipeline) exists, and
// indexing a document (running a log record through that pipeline).
func NewElasticHandler(deps ElasticDeps) http.Handler {
	logger := slog.Default().With(slog.String("channel", "input.middleware.elastic"))
	mux := http.NewServeMux()

	mux.HandleFunc(
		"HEAD /api/v1/middlewares/elastic/{index}",
		func(w http.ResponseWriter, r *http.Request) {
			index := r.PathValue("index")

			user := auth.GetContextUser(r.Context())
			authorized, err := deps.AuthStorage.VerifyUserPermission(
				r.Context(),
				user.Name,
				models.SCOPE_READ_PIPELINES,
			)
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"Failed to verify user permission",
					slog.String("error", err.Error()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !authorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			pipelines, err := deps.ConfigStorage.ListPipelines(r.Context())
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"Failed to list pipelines",
					slog.String("error", err.Error()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if slices.Contains(pipelines, index) {
				w.WriteHeader(http.StatusOK)
				return
			}

			http.Error(w, "Not Found", http.StatusNotFound)
		},
	)

	mux.HandleFunc(
		"POST /api/v1/middlewares/elastic/{index}/_doc",
		func(w http.ResponseWriter, r *http.Request) {
			index := r.PathValue("index")

			user := auth.GetContextUser(r.Context())
			authorized, err := deps.AuthStorage.VerifyUserPermission(
				r.Context(),
				user.Name,
				models.SCOPE_SEND_LOGS,
			)
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"Failed to verify user permission",
					slog.String("error", err.Error()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !authorized {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			defer r.Body.Close()
			dec := json.NewDecoder(r.Body)
			dec.UseNumber()

			var doc map[string]any
			if err := dec.Decode(&doc); err != nil {
				logger.ErrorContext(
					r.Context(),
					"Failed to decode request body",
					slog.String("error", err.Error()),
				)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			fields := map[string]string{}

			joinKey := func(prefix, key string) string {
				if prefix == "" {
					return key
				}
				return prefix + "." + key
			}

			var flatten func(prefix string, v any)
			flatten = func(prefix string, v any) {
				switch t := v.(type) {
				case map[string]any:
					for k, val := range t {
						flatten(joinKey(prefix, k), val)
					}

				case []any:
					for i, val := range t {
						flatten(joinKey(prefix, strconv.Itoa(i)), val)
					}

				case nil:
					fields[prefix] = ""

				case string:
					fields[prefix] = t

				case json.Number:
					fields[prefix] = t.String()

				case bool:
					fields[prefix] = strconv.FormatBool(t)

				default:
					fields[prefix] = fmt.Sprint(t)
				}
			}

			flatten("", doc)

			record := models.NewLogRecord(fields)
			err = deps.PipelineRunner.Run(
				r.Context(),
				index,
				pipelines.DIRECT_ENTRYPOINT,
				record,
			)
			if err != nil {
				logger.ErrorContext(
					r.Context(),
					"Failed to process log entry",
					slog.String("pipeline", index),
					slog.String("error", err.Error()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		},
	)

	return elasticProduct(elasticAuth(deps, mux))
}
