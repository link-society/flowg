package web

import (
	"embed"

	"io"
	"net/http"
	"strings"
)

//go:embed public/**/*.css
//go:embed public/**/*.js
//go:embed public/index.html
var staticfiles embed.FS

func NewHandler() http.Handler {
	return http.StripPrefix(
		"/web/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "assets/") {
				r.URL.Path = "/public/" + r.URL.Path
				http.FileServer(http.FS(staticfiles)).ServeHTTP(w, r)
			} else {
				html, err := staticfiles.Open("public/index.html")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				io.Copy(w, html)
			}
		}),
	)
}
