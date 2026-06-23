package main

import (
	"fmt"
	"os"

	"go.uber.org/fx"

	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"

	"link-society.com/flowg/api"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

func main() {
	var handler http.Handler

	app := fx.New(
		fx.NopLogger,

		// The schema generation only inspects the route table, so nil backends
		// are enough to build the handler.
		fx.Provide(func() auth.Storage { return nil }),
		fx.Provide(func() config.Storage { return nil }),
		fx.Provide(func() log.Storage { return nil }),
		fx.Provide(func() lognotify.LogNotifier { return nil }),
		fx.Provide(func() pipelines.Runner { return nil }),

		api.Module("openapi-handler"),

		fx.Populate(fx.Annotate(&handler, fx.ParamTags(`name:"openapi-handler"`))),
	)
	if err := app.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not build OpenAPI handler: %v\n", err)
		os.Exit(1)
	}

	apiService := handler.(*web.Service)
	reflector := apiService.OpenAPIReflector().(*openapi31.Reflector)
	schema, err := reflector.Spec.MarshalJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not marshal OpenAPI schema to JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("./website/src/openapi.json", schema, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not write OpenAPI schema to file: %v\n", err)
		os.Exit(1)
	}
}
