package health

import (
	"github.com/danielgtaylor/huma/v2"
	"net/http"
)

func RegisterRoutes(api huma.API, c *Controller) {
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Check health",
		Tags:        []string{"Health"},
	}, c.HealthCheck)
}
