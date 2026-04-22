package calculate

import (
	"context"
	"errors"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCalculateBestFit(t *testing.T) {
	t.Run("returns exact match", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(250, []int{250, 500, 1000, 2000, 5000})

		require.True(t, ok)
		assert.Equal(t, 250, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 250, Quantity: 1},
		}, packs)
	})

	t.Run("returns minimal overfill", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(251, []int{250, 500, 1000, 2000, 5000})

		require.True(t, ok)
		assert.Equal(t, 500, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 500, Quantity: 1},
		}, packs)
	})

	t.Run("prioritizes minimal items before minimal packs", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(501, []int{250, 500, 1000, 2000, 5000})

		require.True(t, ok)
		assert.Equal(t, 750, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 500, Quantity: 1},
			{Size: 250, Quantity: 1},
		}, packs)
	})

	t.Run("returns challenge example result", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(12001, []int{250, 500, 1000, 2000, 5000})

		require.True(t, ok)
		assert.Equal(t, 12250, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 5000, Quantity: 2},
			{Size: 2000, Quantity: 1},
			{Size: 250, Quantity: 1},
		}, packs)
	})

	t.Run("returns tips edge case result", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(500000, []int{23, 31, 53})

		require.True(t, ok)
		assert.Equal(t, 500000, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 53, Quantity: 9429},
			{Size: 31, Quantity: 7},
			{Size: 23, Quantity: 2},
		}, packs)
	})

	t.Run("returns no solution when no pack sizes are configured", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(500000, []int{})

		require.False(t, ok)
		assert.Equal(t, 0, totalItems)
		assert.Nil(t, packs)
	})

	t.Run("returns no solution when invalid pack sizes are configured", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(500000, []int{0, 250, 500})

		require.False(t, ok)
		assert.Equal(t, 0, totalItems)
		assert.Nil(t, packs)
	})

	t.Run("works with one pack size only", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(1, []int{250})

		require.True(t, ok)
		assert.Equal(t, 250, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 250, Quantity: 1},
		}, packs)
	})

	t.Run("works with unsorted pack sizes input", func(t *testing.T) {
		totalItems, packs, ok := calculateBestFit(501, []int{1000, 250, 500, 5000, 2000})

		require.True(t, ok)
		assert.Equal(t, 750, totalItems)
		assert.Equal(t, []PackResult{
			{Size: 500, Quantity: 1},
			{Size: 250, Quantity: 1},
		}, packs)
	})
}

func TestCalculatePacks(t *testing.T) {
	t.Run("returns calculated packs", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return([]int{250, 500, 1000, 2000, 5000}, nil)

		controller := NewController(packSizesStore)
		input := &PacksInput{}
		input.Body.Items = 501

		output, err := controller.CalculatePacks(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, &PacksOutput{
			Body: PacksBody{
				RequestedItems: 501,
				TotalItems:     750,
				Packs: []PackResult{
					{Size: 500, Quantity: 1},
					{Size: 250, Quantity: 1},
				},
			},
		}, output)
	})

	t.Run("returns unprocessable entity when no pack sizes are configured", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return([]int{}, nil)

		controller := NewController(packSizesStore)
		input := &PacksInput{}
		input.Body.Items = 500000

		output, err := controller.CalculatePacks(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, output)
		assertStatusError(t, err, 422)
		assert.EqualError(t, err, "order cannot be fulfilled because no pack sizes are configured")
	})

	t.Run("returns unprocessable entity when invalid pack sizes are configured", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return([]int{0, 250, 500}, nil)

		controller := NewController(packSizesStore)
		input := &PacksInput{}
		input.Body.Items = 500000

		output, err := controller.CalculatePacks(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, output)
		assertStatusError(t, err, 422)
		assert.EqualError(t, err, "order cannot be fulfilled because invalid pack sizes are configured")
	})

	t.Run("returns store error", func(t *testing.T) {
		expectedErr := errors.New("store failed")
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return(nil, expectedErr)

		controller := NewController(packSizesStore)
		input := &PacksInput{}
		input.Body.Items = 501

		output, err := controller.CalculatePacks(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func assertStatusError(t *testing.T, err error, expectedStatus int) {
	t.Helper()

	var statusErr huma.StatusError
	require.ErrorAs(t, err, &statusErr)
	assert.Equal(t, expectedStatus, statusErr.GetStatus())
}
