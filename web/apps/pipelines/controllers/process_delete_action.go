package controllers

import (
	"fmt"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"
	"link-society.com/flowg/internal/webutils"
)

func ProcessDeleteAction(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		pipelineName := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !webutils.Permissions(r.Context()).CanEditPipelines {
			url := fmt.Sprintf("/web/pipelines/edit/%s", pipelineName)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}

		err := pipelinesManager.DeletePipelineFlow(pipelineName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to delete pipeline flow", err)

			url := fmt.Sprintf("/web/pipelines/edit/%s", pipelineName)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
	}
}
