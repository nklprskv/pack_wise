package packs

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetPackSizesReturnsConfiguredSizes(t *testing.T) {
	expectedErr := errors.New("store failed")

	t.Run("returns configured sizes", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return([]int{250, 500, 1000}, nil)

		controller := NewController(packSizesStore)

		output, err := controller.GetPackSizes(context.Background(), &struct{}{})
		require.NoError(t, err)
		assert.Equal(t, []int{250, 500, 1000}, output.Body)
	})

	t.Run("returns empty slice when store returns nil", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return(nil, nil)

		controller := NewController(packSizesStore)

		output, err := controller.GetPackSizes(context.Background(), &struct{}{})
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, []int{}, output.Body)
	})

	t.Run("returns store error", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().GetPackSizes(gomock.Any()).Return(nil, expectedErr)

		controller := NewController(packSizesStore)

		_, err := controller.GetPackSizes(context.Background(), &struct{}{})
		require.Error(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
}
