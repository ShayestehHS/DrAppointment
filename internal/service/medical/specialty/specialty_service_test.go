//go:build test
// +build test

package specialty

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty/memory"
)

func setupSpecialtyService() Service {
	repo := memory.NewSpecialtyRepository()
	return NewSpecialtyService(repo)
}

func TestSpecialtyService_ListSpecialtiesOffset_Success(t *testing.T) {
	service := setupSpecialtyService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	specialties, totalCount, err := service.ListSpecialtiesOffset(ctx, params)
	require.NoError(t, err)

	assert.Equal(t, 3, totalCount)
	assert.Len(t, specialties, 3)

	// Verify specialties are sorted by created_at desc (newer first)
	assert.Equal(t, "Dermatology", specialties[0].Name)
	assert.Equal(t, "Neurology", specialties[1].Name)
	assert.Equal(t, "Cardiology", specialties[2].Name)
}

func TestSpecialtyService_ListSpecialtiesOffset_WithPagination(t *testing.T) {
	service := setupSpecialtyService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  2,
		Limit: 2,
	}

	specialties, totalCount, err := service.ListSpecialtiesOffset(ctx, params)
	require.NoError(t, err)

	assert.Equal(t, 3, totalCount)
	assert.Len(t, specialties, 1)

	// Page 2 with limit 2 = offset 2, with only 3 items, should get only the last one
	assert.Equal(t, "Cardiology", specialties[0].Name)
}

func TestSpecialtyService_ListSpecialtiesOffset_EmptyResult(t *testing.T) {
	service := setupSpecialtyService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  100,
		Limit: 10,
	}

	specialties, totalCount, err := service.ListSpecialtiesOffset(ctx, params)
	require.NoError(t, err)

	assert.Equal(t, 3, totalCount)
	assert.Len(t, specialties, 0)
}

func TestSpecialtyService_GetByID_Success(t *testing.T) {
	service := setupSpecialtyService()
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")

	specialty, err := service.GetByID(ctx, specialtyID)
	require.NoError(t, err)
	require.NotNil(t, specialty)

	assert.Equal(t, specialtyID, specialty.ID)
	assert.Equal(t, "Cardiology", specialty.Name)
	assert.NotNil(t, specialty.ImagePath)
}

func TestSpecialtyService_GetByID_NotFound(t *testing.T) {
	service := setupSpecialtyService()
	ctx := context.Background()

	nonExistentID := uuid.New()

	specialty, err := service.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	assert.Nil(t, specialty)
	assert.True(t, errors.Is(err, specialty.ErrSpecialtyNotFound))
}
