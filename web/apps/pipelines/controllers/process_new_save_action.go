package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/pipelines/templates/views"
)

func ProcessNewSaveAction(
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

		pipelineName := ""
		pipelineFlow := "null"

		if !webutils.Permissions(r.Context()).CanEditPipelines {
			webutils.NotifyError(r.Context(), "You do not have permission to edit pipelines")
			goto response
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		pipelineName = r.FormValue("name")
		pipelineFlow = r.FormValue("flow")

		if pipelineName == "" {
			webutils.NotifyError(r.Context(), "Pipeline name is required")
		}

		if pipelineFlow == "" {
			webutils.NotifyError(r.Context(), "Pipeline flow is required")
		}

		if pipelineName == "" || pipelineFlow == "" {
			goto response
		}

		if err := pipelineSys.Write(pipelineName, pipelineFlow); err != nil {
			webutils.LogError(r.Context(), "Failed to save pipeline flow", err)
			webutils.NotifyError(r.Context(), "Could not save pipeline")
			goto response
		}

		webutils.NotifyInfo(r.Context(), "Pipeline saved")

	response:
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
