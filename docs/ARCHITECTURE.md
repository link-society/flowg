# FlowG Architecture

FlowG is a log management platform that lets you ingest, transform, and query
logs using a visual pipeline builder. This document describes the overall
technical architecture of the project: the technologies it is built on, how it
is continuously integrated and released, and how the source tree is organized.

It is meant as an entry point for contributors. For deeper dives into specific
subsystems, see [Related Documentation](#related-documentation).

## Technical Stack

FlowG is a polyglot project. Each language is used where it fits best.

### Backend — Go

The core of FlowG (server, CLI client, health probe, HTTP/Syslog/management
services, storage and processing engines) is written in **Go**.

Notable dependencies:

| Package | Description |
| --- | --- |
| [go.uber.org/fx](https://github.com/uber-go/fx) | dependency injection and application lifecycle management |
| [go-actor](https://github.com/vladopajic/go-actor) | Actor model |
| [BadgerDB](https://github.com/dgraph-io/badger) | embedded key/value store |
| [FoundationDB](https://www.foundationdb.org) | distributed key-value store (optional, replaces BadgerDB for clustered deployments) |
| [swaggest](https://github.com/swaggest) | OpenAPI framework |

### Log Transformation Engine — Rust

Log transformations use the
**[Vector Remap Language (VRL)](https://vector.dev/docs/reference/vrl/)**. Since
VRL is implemented in Rust, FlowG embeds a small Rust crate (`flowg-vrl`) that
wraps the [`vrl`](https://crates.io/crates/vrl) crate and exposes it as a static
library. Go calls into it through **cgo**.

The crate lives in [internal/utils/langs/vrl/rust-crate](../internal/utils/langs/vrl/rust-crate)
and is compiled to a `staticlib` that is linked into the Go binary at build
time.

### Frontend — TypeScript / React

The Web UI is a single-page application built with **TypeScript** and
**React**, bundled with **[Vite](https://vite.dev)**.

Notable dependencies:


| Package | Description |
| --- | --- |
| [MUI](https://mui.com) | UI component library |
| [@xyflow/react](https://reactflow.dev) | visual pipeline builder |
| [@monaco-editor/react](https://github.com/suren-atoyan/monaco-react) | code editor |

The built assets are embedded into the Go binary (see [web/](../web)) so FlowG
ships as a single self-contained executable.

### Documentation Website — Docusaurus

The public documentation site is built with
**[Docusaurus](https://docusaurus.io)** (TypeScript), with
**[redocusaurus](https://github.com/rohit-gohri/redocusaurus)** for rendering
the OpenAPI spec and **[Mermaid](https://mermaid.js.org)** for diagrams. It
lives in [website/](../website).

### End-to-End Test Suite — Python

The end-to-end test runner is written in **Python** and managed with
**[PDM](https://pdm-project.org)**, using **[pytest](https://pytest.org)** as the
orchestrator. Each suite is driven by the tool best suited to it:

- **API** — **[Hurl](https://hurl.dev)** (`.hurl` specs) exercises the REST API.
- **Web** — **[Robot Framework](https://robotframework.org)** (with
  SeleniumLibrary) drives the Web UI.
- **Benchmark** — **[oha](https://github.com/hatoo/oha)** load-tests the
  ingestion endpoints.

Tests spin up FlowG via its Docker image. See [tests/](../tests).

### Build Orchestration & Packaging

- **[Task](https://taskfile.dev)** (`Taskfile.yml` + `scripts/*.taskfile.yml`) —
  the task runner used to build, test, run and release every component.
- **Docker** — container image build (`docker/flowg.dockerfile`).
- **Helm** — Kubernetes deployment chart under [k8s/charts](../k8s/charts).

## Continuous Integration

CI/CD runs on **GitHub Actions**. Workflows live in `.github/workflows/` and are
composed of reusable building blocks (`workflow_call`) wired together by a few
top-level entry points.

### Entry points

| Workflow | Trigger | Purpose |
| --- | --- | --- |
| `ci-main.yml` | push to `main` | Full pipeline: build, test, build binaries, build & deploy website. |
| `ci-pr-app.yml` | PR touching app code | Build Docker image and run end-to-end tests. |
| `ci-pr-website.yml` | PR touching `website/**` | Build binaries and the website. |
| `ci-release.yml` | release | Build and publish release artifacts. |
| `trigger-release.yml` | manual / tag | Kicks off the release process. |

### Reusable workflows

- **`build-docker.yml`** — builds the FlowG Docker image.
- **`build-binary.yml`** — cross-compiles binaries (parameterized by
  `goos`/`goarch`/`rust-target`) and caches artifacts in the GitHub Container
  Registry via [ORAS](https://oras.land).
- **`test-e2e.yml`** — runs the API, Web and Kubernetes end-to-end suites using
  the composite actions under `.github/actions/`.
- **`build-website.yml`** / **`deploy-website.yml`** — build and publish the
  documentation site.
- **`release-binary.yml`** / **`release-docker.yml`** — publish release
  artifacts.

### Quality gates

- **Unit tests** — Go (`go test`) and Rust (`cargo test`), runnable locally with
  `task test:unit`.
- **End-to-end tests** — API, Web and benchmark suites (`task test:e2e:*`).
- **Static analysis** — [SonarCloud](https://sonarcloud.io/summary/new_code?id=link-society_flowg)
  quality gate, surfaced on the README badge.
- **Linting** — `oxlint` + `prettier` for the frontend.

## Directory Structure

```text
flowg/
├── api/                  # REST API handlers (one file per use case) + middlewares
├── cmd/                  # Entry points for the binaries
│   ├── flowg-server/     # Main server (HTTP, Syslog, management APIs)
│   ├── flowg-client/     # CLI client
│   └── flowg-health/     # Health-check probe
├── internal/             # Private application code (not importable externally)
│   ├── app/              # Application wiring (fx modules, bootstrap, logging, metrics, server)
│   ├── engines/          # Processing engines (lognotify, pipelines)
│   ├── models/           # Domain models (auth, flows, forwarders, log records, ...)
│   ├── services/         # Long-running services (http, mgmt, syslog)
│   ├── storage/          # Storage backends (auth, config, log) over BadgerDB or FoundationDB
│   └── utils/            # Shared utilities (api, auth, client, kvstore, langs, otlp, ...)
│       └── langs/vrl/    # VRL bindings, including the Rust crate (rust-crate/)
├── web/                  # Embedded Web UI
│   ├── app/              # React/TypeScript source (Vite project)
│   └── public/           # Built assets embedded into the Go binary
├── website/              # Public documentation site (Docusaurus)
├── tests/                # End-to-end test suite (Python / pytest / Robot Framework)
│   └── specs/            # API, web and benchmark specifications
├── k8s/charts/           # Helm chart for Kubernetes deployment
├── docker/               # Dockerfile and entrypoint
├── scripts/              # Task definitions and code generators (gen/)
├── docs/                 # Project documentation (this file, slides, screenshots)
└── Taskfile.yml          # Root task runner entry point
```

### Versioning conventions in models

The `internal/models` package versions some on-disk/serialized structures with
explicit `_v1` / `_v2` suffixes (e.g. `flow_v1.go`, `flow_v2.go`,
`forwarder_v1.go`, `forwarder_v2.go`) and provides `_convert.go` helpers to
migrate between versions. New work should target the latest version and add a
conversion path when the schema changes.

## Related Documentation

- [README.md](../README.md) — project overview, features, build instructions.
- [CONTRIBUTING.md](../CONTRIBUTING.md) — contribution guidelines.
- [SECURITY.md](../SECURITY.md) — security policy.
- [docs/README.md](./README.md) — documentation index.
- [Public documentation](https://docs.flowg.cloud) — user-facing docs.
- **Go API reference** — generated from in-code doc comments and rendered with
  [`pkgsite`](https://pkg.go.dev/golang.org/x/pkgsite).
- **REST API reference** — OpenAPI spec served at `/api/docs` and generated via
  `task doc:gen:openapi`.
- **CLI reference** — generated via `task doc:gen:cli`.
