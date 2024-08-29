package controllers

import (
	"log/slog"

	"fmt"

	"net/http"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"
)

func ProcessDeleteAction(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("name")

		permissions := auth.Permissions{}
		user := auth.GetContextUser(r.Context())
		scopes, err := userSys.ListUserScopes(user.Name)
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

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !permissions.CanEditPipelines {
			http.Redirect(w, r, fmt.Sprintf("/web/pipelines/edit/%s", pipelineName), http.StatusSeeOther)
			return
		}

		err = pipelinesManager.DeletePipelineFlow(pipelineName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error deleting pipeline",
				"channel", "web",
				"error", err.Error(),
			)
		}

		http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
	}
}
