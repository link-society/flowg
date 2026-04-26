package web

import (
	"embed"

	"mime"
	"path"

	"html/template"
	"strings"

	"compress/gzip"
	"encoding/base64"

	"io"
	"net/http"

	"link-society.com/flowg/internal/app/featureflags"
)

//go:embed public/**/*.gz
//go:embed public/*.gz
var staticfiles embed.FS

func NewHandler() http.Handler {
	return http.StripPrefix(
		"/web/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				http.Error(w, "gzip required", http.StatusNotAcceptable)
				return
			}

			if strings.HasPrefix(r.URL.Path, "assets/") {
				reqpath := "public/" + r.URL.Path
				realpath := reqpath + ".gz"
				r.URL.Path = realpath

				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(reqpath)))
				w.Header().Set("Cache-Control", "public, max-age=86400")
				w.Header().Set("ETag", base64.StdEncoding.EncodeToString([]byte(r.URL.Path)))

				http.FileServer(http.FS(staticfiles)).ServeHTTP(w, r)
			} else {
				htmlTemplateFile, err := staticfiles.Open("public/index.html.gz")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer htmlTemplateFile.Close()

				htmlTemplateReader, err := gzip.NewReader(htmlTemplateFile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer htmlTemplateReader.Close()

				htmlTemplateSource, err := io.ReadAll(htmlTemplateReader)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				htmlTemplate, err := template.New("index").Parse(string(htmlTemplateSource))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				data := map[string]any{
					"FeatureFlags": map[string]bool{
						"DemoMode": featureflags.GetDemoMode(),
					},
				}

				w.WriteHeader(http.StatusOK)
				htmlTemplate.Execute(w, data)
			}
		}),
	)
}
