package packs

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"pack_wise/internal/store"
)

type UpdatePackSizesInput struct {
	Body struct {
		Sizes []int `json:"sizes"`
	}
}

type UpdatePackSizesOutput struct{}

// UpdatePackSizes replaces the configured pack sizes after validating that the
// request contains unique positive integers.
func (c *Controller) UpdatePackSizes(ctx context.Context, input *UpdatePackSizesInput) (*UpdatePackSizesOutput, error) {
	if len(input.Body.Sizes) == 0 {
		return nil, huma.Error400BadRequest("sizes must not be empty")
	}

	if err := c.packSizesStore.ReplacePackSizes(ctx, input.Body.Sizes); err != nil {
		if errors.Is(err, store.ErrInvalidPackSizes) {
			return nil, huma.Error400BadRequest("sizes must be unique positive integers")
		}

		return nil, err
	}

	return &UpdatePackSizesOutput{}, nil
}
