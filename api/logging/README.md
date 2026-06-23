# logging

The package at `api/logging` provides the logger shared by every endpoint and
middleware of FlowG's REST API.

It exists so API code never reaches for the global logger directly. Routing all
API logs through a single helper guarantees they carry a common `api` channel
tag, which lets operators tell API activity apart from the rest of FlowG's
output when filtering logs.

The logger is resolved lazily on each call, so it always reflects the logging
configuration in force when a request is handled rather than the one present at
startup.
