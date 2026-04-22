package packs

import "context"

type DeletePackSizeInput struct {
	Size int `path:"size" minimum:"1"`
}

// DeletePackSize removes one configured pack size identified by its value.
func (c *Controller) DeletePackSize(ctx context.Context, input *DeletePackSizeInput) (*struct{}, error) {
	if err := c.packSizesStore.DeletePackSize(ctx, input.Size); err != nil {
		return nil, err
	}

	return nil, nil
}
