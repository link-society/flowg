package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/pipelines/templates/views"
)

func PageEdit(
	userSys *auth.UserSystem,
	pipelineSys *config.PipelineSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		pipelineName := r.PathValue("name")
		pipelineFlow, err := pipelineSys.Read(pipelineName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch pipeline flow", err)
			http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
			return
		}

		pipelines, err := pipelineSys.List()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch pipelines", err)
			webutils.NotifyError(r.Context(), "Could not fetch pipelines")
			pipelines = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
		))
		h.ServeHTTP(w, r)
	}
}
