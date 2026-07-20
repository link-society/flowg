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

// staticfiles holds the pre-built, gzip-compressed assets of the web UI,
// embedded into the binary at build time.
//
// The whole public/ tree is embedded recursively: Go's embed patterns do not
// support "**", so a suffix glob like "public/**/*.gz" only reaches a fixed
// depth and would silently drop nested assets (e.g. the i18n locale files under
// public/assets/locales/<lng>/). The build emits only gzip-compressed files, so
// embedding the directory captures exactly the assets that are served.
//
//go:embed all:public
var staticfiles embed.FS

// NewHandler builds the HTTP handler that serves FlowG's single-page web UI.
//
// The handler is mounted under "/web/" and serves two kinds of requests:
//
//   - Requests under "assets/" are served as static, immutable files straight
//     from the embedded filesystem, with long-lived cache headers.
//   - Every other request returns the SPA entry point ("index.html"), rendered
//     as a Go template so runtime values — the enabled feature flags and the
//     mountPath the UI is served from — can be injected into the page.
//
// All assets are stored and served gzip-compressed, so clients must advertise
// gzip support through the Accept-Encoding header; requests that do not are
// rejected with http.StatusNotAcceptable.
func NewHandler(mountPath string) http.Handler {
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
					"MountPath": mountPath,
				}

				w.WriteHeader(http.StatusOK)
				htmlTemplate.Execute(w, data)
			}
		}),
	)
}
