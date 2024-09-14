package web

import (
	"embed"

	"strings"

	"io"
	"net/http"
)

//go:embed public/**/*.css
//go:embed public/**/*.js
//go:embed public/**/*.js.map
//go:embed public/index.html
var staticfiles embed.FS

//go:generate templ generate

func NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"/web/",
		http.StripPrefix("/web/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		})),
	)

	return mux
}
