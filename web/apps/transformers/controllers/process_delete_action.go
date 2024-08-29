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
		transformerName := r.PathValue("name")

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
	}
}
