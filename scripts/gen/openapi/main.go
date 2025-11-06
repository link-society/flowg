package main

import (
	"log"
	"os"

	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"

	"link-society.com/flowg/api"
)

func main() {
	apiService := api.NewHandler(api.Dependencies{}).(*web.Service)
	reflector := apiService.OpenAPIReflector().(*openapi31.Reflector)
	schema, err := reflector.Spec.MarshalJSON()
	if err != nil {
		log.Fatalf("ERROR: Could not marshal OpenAPI schema to JSON: %v", err)
	}

	if err := os.WriteFile("./website/src/openapi.json", schema, 0644); err != nil {
		log.Fatalf("ERROR: Could not write OpenAPI schema to file: %v", err)
	}
}
