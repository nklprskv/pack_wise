package calculate

import (
	"context"
	"sort"

	"github.com/danielgtaylor/huma/v2"
)

type PacksInput struct {
	Body struct {
		Items int `json:"items" minimum:"1" doc:"Requested number of items"`
	}
}

type PackResult struct {
	Size     int `json:"size"`
	Quantity int `json:"quantity"`
}

type PacksBody struct {
	RequestedItems int          `json:"requestedItems"`
	TotalItems     int          `json:"totalItems"`
	Packs          []PackResult `json:"packs"`
}

type PacksOutput struct {
	Body PacksBody
}

// CalculatePacks returns the smallest shippable quantity that fulfills the
// requested item count and, for that quantity, the fewest possible packs.
func (c *Controller) CalculatePacks(ctx context.Context, input *PacksInput) (*PacksOutput, error) {
	packSizes, err := c.packSizesStore.GetPackSizes(ctx)
	if err != nil {
		return nil, err
	}

	if len(packSizes) == 0 {
		return nil, huma.Error422UnprocessableEntity("order cannot be fulfilled because no pack sizes are configured")
	}

	if !hasOnlyPositivePackSizes(packSizes) {
		return nil, huma.Error422UnprocessableEntity("order cannot be fulfilled because invalid pack sizes are configured")
	}

	totalItems, packs, ok := calculateBestFit(input.Body.Items, packSizes)
	if !ok {
		return nil, huma.Error422UnprocessableEntity("order cannot be fulfilled with the current pack sizes")
	}

	return &PacksOutput{
		Body: PacksBody{
			RequestedItems: input.Body.Items,
			TotalItems:     totalItems,
			Packs:          packs,
		},
	}, nil
}

type planStep struct {
	packCount int
	prevTotal int
	packSize  int
	reachable bool
}

// calculateBestFit uses dynamic programming to find the minimal fulfillable
// total for the requested quantity and the corresponding minimal-pack
// breakdown for that total.
func calculateBestFit(requestedItems int, packSizes []int) (int, []PackResult, bool) {
	if len(packSizes) == 0 {
		return 0, nil, false
	}

	if !hasOnlyPositivePackSizes(packSizes) {
		return 0, nil, false
	}

	limit := requestedItems + maxPackSize(packSizes)
	plan := make([]planStep, limit+1)
	plan[0] = planStep{reachable: true}

	for total := 1; total <= limit; total++ {
		for _, packSize := range packSizes {
			if total < packSize || !plan[total-packSize].reachable {
				continue
			}

			candidatePackCount := plan[total-packSize].packCount + 1
			if !plan[total].reachable || candidatePackCount < plan[total].packCount {
				plan[total] = planStep{
					packCount: candidatePackCount,
					prevTotal: total - packSize,
					packSize:  packSize,
					reachable: true,
				}
			}
		}
	}

	bestTotal := 0
	for total := requestedItems; total <= limit; total++ {
		if plan[total].reachable {
			bestTotal = total
			break
		}
	}

	if bestTotal == 0 {
		return 0, nil, false
	}

	packQuantities := make(map[int]int)
	for total := bestTotal; total > 0; total = plan[total].prevTotal {
		packQuantities[plan[total].packSize]++
	}

	packs := make([]PackResult, 0, len(packQuantities))
	for _, packSize := range packSizes {
		if quantity := packQuantities[packSize]; quantity > 0 {
			packs = append(packs, PackResult{
				Size:     packSize,
				Quantity: quantity,
			})
		}
	}

	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Size > packs[j].Size
	})

	return bestTotal, packs, true
}

// maxPackSize returns the largest configured pack size.
func maxPackSize(packSizes []int) int {
	maxPackSize := 0

	for _, packSize := range packSizes {
		if packSize > maxPackSize {
			maxPackSize = packSize
		}
	}

	return maxPackSize
}

// hasOnlyPositivePackSizes reports whether all configured pack sizes are valid.
func hasOnlyPositivePackSizes(packSizes []int) bool {
	for _, packSize := range packSizes {
		if packSize <= 0 {
			return false
		}
	}

	return true
}
