package main

import (
	"log"
	"os"

	"go.uber.org/fx"

	"net/http"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"

	"link-society.com/flowg/api"
	_ "link-society.com/flowg/api/operations"
	"link-society.com/flowg/api/routing"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	logStorage "link-society.com/flowg/internal/storage/log"
)

func main() {
	var handler http.Handler

	app := fx.New(
		fx.NopLogger,

		// The schema generation only inspects the route table, so nil backends
		// are enough to build the handler.
		fx.Provide(func() authStorage.Storage { return nil }),
		fx.Provide(func() config.Storage { return nil }),
		fx.Provide(func() logStorage.Storage { return nil }),
		fx.Provide(func() lognotify.LogNotifier { return nil }),
		fx.Provide(func() pipelines.Runner { return nil }),

		routing.Module("api.operations"),
		fx.Provide(api.NewHandler),

		fx.Populate(&handler),
	)
	if err := app.Err(); err != nil {
		log.Fatalf("ERROR: Could not build OpenAPI handler: %v", err)
	}

	apiService := handler.(*web.Service)
	reflector := apiService.OpenAPIReflector().(*openapi31.Reflector)
	schema, err := reflector.Spec.MarshalJSON()
	if err != nil {
		log.Fatalf("ERROR: Could not marshal OpenAPI schema to JSON: %v", err)
	}

	if err := os.WriteFile("./website/src/openapi.json", schema, 0644); err != nil {
		log.Fatalf("ERROR: Could not write OpenAPI schema to file: %v", err)
	}
}
