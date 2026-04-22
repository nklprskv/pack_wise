package app

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"pack_wise/internal/api/calculate"
	"pack_wise/internal/api/health"
	"pack_wise/internal/api/packs"
	"pack_wise/internal/store"
)

// NewServer builds the HTTP server and wires the application routes and
// middleware around a shared pack sizes store.
func NewServer(addr string, packSizesStore store.PackSizesStore) *http.Server {
	mux := http.NewServeMux()
	config := huma.DefaultConfig("pack_wise API", "0.1.0")
	config.CreateHooks = nil
	config.DocsPath = "/swagger"
	config.DocsRenderer = huma.DocsRendererSwaggerUI

	api := humago.New(mux, config)
	registerRoutes(api, packSizesStore)

	return &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(mux),
	}
}

// registerRoutes attaches all API groups to the shared Huma instance.
func registerRoutes(api huma.API, packSizesStore store.PackSizesStore) {
	healthController := health.NewController()
	health.RegisterRoutes(api, healthController)

	packsController := packs.NewController(packSizesStore)
	packs.RegisterRoutes(api, packsController)

	calculateController := calculate.NewController(packSizesStore)
	calculate.RegisterRoutes(api, calculateController)
}
