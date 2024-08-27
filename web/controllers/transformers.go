package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func TransformersController(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/transformers/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/transformers/new", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET /web/transformers/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		h := templ.Handler(views.Transformers(
			views.TransformersProps{
				Transformers:       transformers,
				CurrentTransformer: "",
				Code:               ".",
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/transformers/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformerName := ""
		transformerCode := "."

		err = r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not parse form")
		} else {
			transformerName = r.FormValue("name")
			transformerCode = r.FormValue("code")

			if permissions.CanEditTransformers {
				if transformerName == "" {
					notifications = append(notifications, "❌ Transformer name is required")
				}

				if transformerCode == "" {
					notifications = append(notifications, "❌ Transformer code is required")
				}

				if transformerName != "" && transformerCode != "" {
					err = pipelinesManager.SaveTransformerScript(transformerName, transformerCode)
					if err != nil {
						slog.ErrorContext(
							r.Context(),
							"error saving transformer script",
							"channel", "web",
							"error", err.Error(),
						)

						notifications = append(notifications, "❌ Could not save transformer script")
					} else {
						notifications = append(notifications, "✅ Transformer script saved")
					}
				}
			}
		}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		h := templ.Handler(views.Transformers(
			views.TransformersProps{
				Transformers:       transformers,
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/transformers/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")
		transformerCode, err := pipelinesManager.GetTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
			return
		}

		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		h := templ.Handler(views.Transformers(
			views.TransformersProps{
				Transformers:       transformers,
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/transformers/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")
		transformerCode, err := pipelinesManager.GetTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
			return
		}

		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not parse form")
		} else {
			transformerName = r.FormValue("name")
			transformerCode = r.FormValue("code")

			if permissions.CanEditTransformers {
				if transformerName == "" {
					notifications = append(notifications, "❌ Transformer name is required")
				}

				if transformerCode == "" {
					notifications = append(notifications, "❌ Transformer code is required")
				}

				if transformerName != "" && transformerCode != "" {
					err = pipelinesManager.SaveTransformerScript(transformerName, transformerCode)
					if err != nil {
						slog.ErrorContext(
							r.Context(),
							"error saving transformer script",
							"channel", "web",
							"error", err.Error(),
						)

						notifications = append(notifications, "❌ Could not save transformer script")
					} else {
						notifications = append(notifications, "✅ Transformer script saved")
					}
				}
			}
		}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		h := templ.Handler(views.Transformers(
			views.TransformersProps{
				Transformers:       transformers,
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/transformers/delete/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")

		permissions := auth.Permissions{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !permissions.CanEditTransformers {
			http.Redirect(w, r, fmt.Sprintf("/web/transformers/edit/%s", transformerName), http.StatusSeeOther)
			return
		}

		err = pipelinesManager.DeleteTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error deleting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
		}

		http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
	})

	return mux
}
