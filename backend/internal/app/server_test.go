package app_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"pack_wise/internal/app"
)

type testPackSizesStore struct {
	sizes []int
}

func (s *testPackSizesStore) GetPackSizes(_ context.Context) ([]int, error) {
	return append([]int(nil), s.sizes...), nil
}

func (s *testPackSizesStore) ReplacePackSizes(_ context.Context, sizes []int) error {
	s.sizes = append([]int(nil), sizes...)
	return nil
}

func (s *testPackSizesStore) DeletePackSize(_ context.Context, size int) error {
	filtered := make([]int, 0, len(s.sizes))
	for _, currentSize := range s.sizes {
		if currentSize != size {
			filtered = append(filtered, currentSize)
		}
	}

	s.sizes = filtered
	return nil
}

func TestGetPackSizesEndpoint(t *testing.T) {
	tests := []struct {
		name         string
		sizes        []int
		expectedBody string
	}{
		{
			name:         "returns configured pack sizes",
			sizes:        []int{250, 500, 1000, 2000, 5000},
			expectedBody: "[250,500,1000,2000,5000]\n",
		},
		{
			name:         "returns empty list when no pack sizes are configured",
			sizes:        []int{},
			expectedBody: "[]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/packs", nil)
			recorder := httptest.NewRecorder()

			server := app.NewServer(":8080", &testPackSizesStore{sizes: tt.sizes})
			server.Handler.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			assert.Equal(t, tt.expectedBody, recorder.Body.String())
		})
	}
}

func TestCalculatePacksEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"items":501}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server := app.NewServer(":8080", &testPackSizesStore{sizes: []int{250, 500, 1000, 2000, 5000}})
	server.Handler.ServeHTTP(recorder, req)

	expectedBody := "{\"requestedItems\":501,\"totalItems\":750,\"packs\":[{\"size\":500,\"quantity\":1},{\"size\":250,\"quantity\":1}]}\n"
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.Equal(t, expectedBody, recorder.Body.String())
}

func TestUpdatePackSizesEndpoint(t *testing.T) {
	packSizesStore := &testPackSizesStore{sizes: []int{250, 500, 1000, 2000, 5000}}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/packs", bytes.NewBufferString(`{"sizes":[23,31,53]}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server := app.NewServer(":8080", packSizesStore)
	server.Handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, []int{23, 31, 53}, packSizesStore.sizes)
}

func TestDeletePackSizeEndpoint(t *testing.T) {
	packSizesStore := &testPackSizesStore{sizes: []int{23, 31, 53}}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/packs/31", nil)
	recorder := httptest.NewRecorder()

	server := app.NewServer(":8080", packSizesStore)
	server.Handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, []int{23, 53}, packSizesStore.sizes)
}

func TestCalculatePacksEndpointReturns422WhenNoPackSizesConfigured(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"items":500000}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	server := app.NewServer(":8080", &testPackSizesStore{sizes: []int{}})
	server.Handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "no pack sizes are configured")
}
