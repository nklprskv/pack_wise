package packs

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeletePackSize(t *testing.T) {
	t.Run("deletes one pack size", func(t *testing.T) {
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().DeletePackSize(gomock.Any(), 31).Return(nil)

		controller := NewController(packSizesStore)

		output, err := controller.DeletePackSize(context.Background(), &DeletePackSizeInput{Size: 31})

		require.NoError(t, err)
		assert.Nil(t, output)
	})

	t.Run("returns store error", func(t *testing.T) {
		expectedErr := errors.New("delete failed")
		mockController := gomock.NewController(t)
		packSizesStore := NewMockPackSizesStore(mockController)
		packSizesStore.EXPECT().DeletePackSize(gomock.Any(), 31).Return(expectedErr)

		controller := NewController(packSizesStore)

		output, err := controller.DeletePackSize(context.Background(), &DeletePackSizeInput{Size: 31})

		require.Error(t, err)
		assert.Nil(t, output)
		assert.ErrorIs(t, err, expectedErr)
	})
}
