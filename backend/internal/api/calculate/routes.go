package calculate

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, c *Controller) {
	huma.Register(api, huma.Operation{
		OperationID: "calculate-packs",
		Method:      http.MethodPost,
		Path:        "/api/v1/calculate",
		Summary:     "Calculate pack breakdown",
		Tags:        []string{"Calculate"},
	}, c.CalculatePacks)
}
