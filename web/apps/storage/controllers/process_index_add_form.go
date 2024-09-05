package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/storage/templates/components"
)

func ProcessIndexAddForm(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanEditStreams {
			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "delete")
			w.WriteHeader(http.StatusOK)
			return
		}

		streamName := r.PathValue("name")

		streamConfig, err := metaSys.GetStreamConfig(streamName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch stream config", err)
			webutils.NotifyError(r.Context(), "Failed to fetch stream config")

			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "delete")
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Failed to parse form")

			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "delete")
			w.WriteHeader(http.StatusOK)
			return
		}

		fieldName := r.FormValue("data_field_name")
		for _, field := range streamConfig.IndexedFields {
			if field == fieldName {
				webutils.NotifyError(r.Context(), "Field already indexed")

				trigger := htmx.Trigger{
					ToastEvent: &htmx.ToastEvent{
						Messages: webutils.Notifications(r.Context()),
					},
				}

				trigger.Write(r.Context(), w)
				htmx.Reswap(w, "delete")
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		streamConfig.IndexedFields = append(streamConfig.IndexedFields, fieldName)
		if err := metaSys.ConfigureStream(streamName, streamConfig); err != nil {
			webutils.LogError(r.Context(), "Failed to configure stream", err)
			webutils.NotifyError(r.Context(), "Failed to configure stream")

			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "delete")
			w.WriteHeader(http.StatusOK)
			return
		}

		webutils.NotifyInfo(r.Context(), "Field index added")

		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}
		trigger.Write(r.Context(), w)

		h := templ.Handler(components.IndexFieldForm(
			components.IndexFieldFormProps{
				StreamName: streamName,
				Field:      fieldName,
			},
		))
		h.ServeHTTP(w, r)
	}
}
