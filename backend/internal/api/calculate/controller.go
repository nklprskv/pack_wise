package calculate

import "pack_wise/internal/store"

//go:generate mockgen -source=../../store/pack_sizes_store.go -destination=./mocks.go -package=calculate

type Controller struct {
	packSizesStore store.PackSizesStore
}

func NewController(packSizesStore store.PackSizesStore) *Controller {
	return &Controller{packSizesStore: packSizesStore}
}
