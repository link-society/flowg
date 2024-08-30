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
		transformerName := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !webutils.Permissions(r.Context()).CanEditTransformers {
			url := fmt.Sprintf("/web/transformers/edit/%s", transformerName)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}

		if err := pipelinesManager.DeleteTransformerScript(transformerName); err != nil {
			webutils.LogError(r.Context(), "Failed to delete transformer script", err)
		}

		http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
	}
}
