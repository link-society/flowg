# web

The package at `web` serves FlowG's single-page web UI.

It exists to keep the front-end concern self-contained behind a single HTTP
handler. The UI is built ahead of time into static assets that are embedded
into the binary, so the server ships as one self-contained executable with no
external file dependencies, and the rest of the application only ever sees an
`http.Handler`.

## Responsibilities

- **Asset serving** — `NewHandler` serves the pre-built UI from the embedded
  filesystem. Requests under `assets/` return immutable static files with
  long-lived cache headers; every other request returns the SPA entry point.
- **Runtime injection** — `index.html` is rendered as a Go template so values
  only known at runtime — the enabled feature flags and the mount path the UI
  is served from — can be injected into the page on each request.
- **Wiring** — `Module` provides the handler to the application's `fx`
  dependency graph under a named tag so the HTTP server can mount it.

## Compression

The UI is built and stored gzip-compressed (`public/**/*.gz`), and assets are
served compressed as-is. Clients must therefore advertise gzip support through
the `Accept-Encoding` header; requests that do not are rejected with
`406 Not Acceptable`.

## Layout

- **[app](app)** — the front-end source and its build tooling. Its build output
  lands in `public/`, which is what gets embedded.
- **public/** — the generated, gzip-compressed assets embedded into the binary.
