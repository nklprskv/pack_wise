package packs

import "context"

type GetPackSizesOutput struct {
	Body []int
}

// GetPackSizes returns the currently configured pack sizes as a JSON array.
func (c *Controller) GetPackSizes(ctx context.Context, _ *struct{}) (*GetPackSizesOutput, error) {
	packSizes, err := c.packSizesStore.GetPackSizes(ctx)
	if err != nil {
		return nil, err
	}

	if packSizes == nil {
		packSizes = []int{}
	}

	return &GetPackSizesOutput{
		Body: packSizes,
	}, nil
}
