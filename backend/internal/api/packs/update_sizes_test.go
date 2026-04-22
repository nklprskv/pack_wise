package packs

import (
	"context"
	"errors"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"pack_wise/internal/store"
)

func TestUpdatePackSizes(t *testing.T) {
	t.Run("updates pack sizes", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().ReplacePackSizes(gomock.Any(), []int{23, 31, 53}).Return(nil)

		controller := NewController(packSizesStore)
		input := &UpdatePackSizesInput{}
		input.Body.Sizes = []int{23, 31, 53}

		output, err := controller.UpdatePackSizes(context.Background(), input)

		require.NoError(t, err)
		require.NotNil(t, output)
		assert.IsType(t, &UpdatePackSizesOutput{}, output)
	})

	t.Run("returns bad request when sizes are empty", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)

		controller := NewController(packSizesStore)
		input := &UpdatePackSizesInput{}

		output, err := controller.UpdatePackSizes(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, output)
		assertStatusError(t, err, 400)
		assert.EqualError(t, err, "sizes must not be empty")
	})

	t.Run("returns bad request when store rejects invalid sizes", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().ReplacePackSizes(gomock.Any(), []int{23, 23, 53}).Return(store.ErrInvalidPackSizes)

		controller := NewController(packSizesStore)
		input := &UpdatePackSizesInput{}
		input.Body.Sizes = []int{23, 23, 53}

		output, err := controller.UpdatePackSizes(context.Background(), input)

		require.Error(t, err)
		assert.Nil(t, output)
		assertStatusError(t, err, 400)
		assert.EqualError(t, err, "sizes must be unique positive integers")
	})

	t.Run("returns store error", func(t *testing.T) {
		expectedErr := errors.New("replace failed")
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().ReplacePackSizes(gomock.Any(), []int{23, 31, 53}).Return(expectedErr)

		controller := NewController(packSizesStore)
		input := &UpdatePackSizesInput{}
		input.Body.Sizes = []int{23, 31, 53}

		output, err := controller.UpdatePackSizes(context.Background(), input)

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
