package packs

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, c *Controller) {
	huma.Register(api, huma.Operation{
		OperationID: "get-packs",
		Method:      http.MethodGet,
		Path:        "/api/v1/packs",
		Summary:     "Get available pack sizes",
		Tags:        []string{"Packs"},
	}, c.GetPackSizes)

	huma.Register(api, huma.Operation{
		OperationID:   "update-packs",
		Method:        http.MethodPost,
		Path:          "/api/v1/packs",
		Summary:       "Update available pack sizes",
		Tags:          []string{"Packs"},
		DefaultStatus: http.StatusOK,
	}, c.UpdatePackSizes)

	huma.Register(api, huma.Operation{
		OperationID: "delete-pack-size",
		Method:      http.MethodDelete,
		Path:        "/api/v1/packs/{size}",
		Summary:     "Delete one pack size",
		Tags:        []string{"Packs"},
	}, c.DeletePackSize)
}
