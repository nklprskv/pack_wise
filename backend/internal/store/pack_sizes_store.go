package store

import "context"

type PackSizesStore interface {
	GetPackSizes(ctx context.Context) ([]int, error)
	ReplacePackSizes(ctx context.Context, sizes []int) error
	DeletePackSize(ctx context.Context, size int) error
}
