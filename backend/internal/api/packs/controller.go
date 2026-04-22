package packs

import "pack_wise/internal/store"

//go:generate mockgen -source=../../store/pack_sizes_store.go -destination=./mocks.go -package=packs

type Controller struct {
	packSizesStore store.PackSizesStore
}

// NewController builds a packs controller backed by the configured store.
func NewController(packSizesStore store.PackSizesStore) *Controller {
	return &Controller{packSizesStore: packSizesStore}
}
