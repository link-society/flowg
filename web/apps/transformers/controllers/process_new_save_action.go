package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/transformers/templates/views"
)

func ProcessNewSaveAction(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformerName := ""
		transformerCode := "."

		if !webutils.Permissions(r.Context()).CanEditTransformers {
			webutils.NotifyError(r.Context(), "You do not have permission to edit transformers")
			goto response
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		transformerName = r.FormValue("name")
		transformerCode = r.FormValue("code")

		if transformerName == "" {
			webutils.NotifyError(r.Context(), "Transformer name is required")
		}

		if transformerCode == "" {
			webutils.NotifyError(r.Context(), "Transformer code is required")
		}

		if transformerName == "" || transformerCode == "" {
			goto response
		}

		if err := pipelinesManager.SaveTransformerScript(transformerName, transformerCode); err != nil {
			webutils.LogError(r.Context(), "Failed to save transformer script", err)
			webutils.NotifyError(r.Context(), "Could not save transformer")
			goto response
		}

		webutils.NotifyInfo(r.Context(), "Transformer script saved")

	response:
		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch transformers", err)
			webutils.NotifyError(r.Context(), "Could not fetch transformers")
			transformers = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Transformers:       transformers,
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
		))
		h.ServeHTTP(w, r)
	}
}
