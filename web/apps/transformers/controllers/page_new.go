package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/transformers/templates/views"
)

func PageNew(
	userSys *auth.UserSystem,
	transformerSys *config.TransformerSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformers, err := transformerSys.List()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch transformers", err)
			webutils.NotifyError(r.Context(), "Could not fetch transformers")
			transformers = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Transformers:       transformers,
				CurrentTransformer: "",
				Code:               ".",
			},
		))
		h.ServeHTTP(w, r)
	}
}
