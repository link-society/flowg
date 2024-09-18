package web

import (
	"embed"

	"encoding/base64"
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

				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Cache-Control", "public, max-age=86400")
				w.Header().Set("ETag", base64.StdEncoding.EncodeToString([]byte(r.URL.Path)))

				http.FileServer(http.FS(staticfiles)).ServeHTTP(w, r)
			} else {
				html, err := staticfiles.Open("public/index.html")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Encoding", "gzip")
				w.WriteHeader(http.StatusOK)
				io.Copy(w, html)
			}
		}),
	)
}
