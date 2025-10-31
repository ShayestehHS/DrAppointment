//go:build test
// +build test

package memory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

func setupSpecialtyMemoryRepo() *specialtyRepository {
	return &specialtyRepository{
		specialties: []medical.Specialty{
			{
				ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
				Name:      "Cardiology",
				ImagePath: medical.NewSpecialtyImage("cardiology.jpg"),
				CreatedAt: time.Now().Add(-24 * time.Hour),
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			{
				ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Name:      "Neurology",
				ImagePath: medical.NewSpecialtyImage("neurology.jpg"),
				CreatedAt: time.Now().Add(-12 * time.Hour),
				UpdatedAt: time.Now().Add(-12 * time.Hour),
			},
			{
				ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174002"),
				Name:      "Dermatology",
				ImagePath: medical.NewSpecialtyImage("dermatology.jpg"),
				CreatedAt: time.Now().Add(-6 * time.Hour),
				UpdatedAt: time.Now().Add(-6 * time.Hour),
			},
			{
				ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174003"),
				Name:      "Orthopedics",
				ImagePath: medical.NewSpecialtyImage("orthopedics.jpg"),
				CreatedAt: time.Now().Add(-3 * time.Hour),
				UpdatedAt: time.Now().Add(-3 * time.Hour),
			},
		},
	}
}

func TestSpecialtyMemoryRepository_ListOffset_Success(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 4)

	// Verify specialties are sorted by created_at desc (newer first)
	assert.Equal(t, "Orthopedics", result[0].Name)
	assert.Equal(t, "Dermatology", result[1].Name)
	assert.Equal(t, "Neurology", result[2].Name)
	assert.Equal(t, "Cardiology", result[3].Name)
}

func TestSpecialtyMemoryRepository_ListOffset_WithPagination(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  2,
		Limit: 2,
	}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 2)

	// Page 2 with limit 2 = offset 2, should get the 3rd and 4th specialties
	assert.Equal(t, "Neurology", result[0].Name)
	assert.Equal(t, "Cardiology", result[1].Name)
}

func TestSpecialtyMemoryRepository_ListOffset_LimitExceedsAvailable(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 100, // Limit exceeds available items
	}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 4) // Should return all 4 items even though limit is 100

	// Should get all specialties
	assert.Equal(t, "Orthopedics", result[0].Name)
	assert.Equal(t, "Dermatology", result[1].Name)
	assert.Equal(t, "Neurology", result[2].Name)
	assert.Equal(t, "Cardiology", result[3].Name)
}

func TestSpecialtyMemoryRepository_ListOffset_EmptyResult(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  100,
		Limit: 10,
	}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)
}

func TestSpecialtyMemoryRepository_ListOffset_ZeroLimit(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 0,
	}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 0) // No items returned with limit 0
}

func TestSpecialtyMemoryRepository_GetByID_Success(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")

	specialty, err := repo.GetByID(ctx, specialtyID)
	require.NoError(t, err)
	require.NotNil(t, specialty)

	assert.Equal(t, specialtyID, specialty.ID)
	assert.Equal(t, "Cardiology", specialty.Name)
	assert.NotNil(t, specialty.ImagePath)
	assert.Contains(t, specialty.ImagePath.GetPath().String(), "cardiology.jpg")
}

func TestSpecialtyMemoryRepository_GetByID_NotFound(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	nonExistentID := uuid.New()

	specialty, err := repo.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	assert.Nil(t, specialty)
	assert.Contains(t, err.Error(), "specialty not found")
}

func TestSpecialtyMemoryRepository_AddSpecialty(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()

	newSpecialty := medical.Specialty{
		ID:        uuid.New(),
		Name:      "Psychiatry",
		ImagePath: medical.NewSpecialtyImage("psychiatry.jpg"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.AddSpecialty(newSpecialty)

	// Verify the specialty was added
	ctx := context.Background()
	params := pagination.LimitOffsetParams{Page: 1, Limit: 10}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 5)

	// The new specialty should be first (most recent)
	assert.Equal(t, "Psychiatry", result[0].Name)
}

func TestSpecialtyMemoryRepository_Clear(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()

	repo.Clear()

	// Verify all specialties were removed
	ctx := context.Background()
	params := pagination.LimitOffsetParams{Page: 1, Limit: 10}

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)
}

func TestSpecialtyMemoryRepository_GetByID_AllSpecialties(t *testing.T) {
	repo := setupSpecialtyMemoryRepo()
	ctx := context.Background()

	// Test getting all specialties by ID
	expectedSpecialties := map[string]string{
		"223e4567-e89b-12d3-a456-426614174000": "Cardiology",
		"223e4567-e89b-12d3-a456-426614174001": "Neurology",
		"223e4567-e89b-12d3-a456-426614174002": "Dermatology",
		"223e4567-e89b-12d3-a456-426614174003": "Orthopedics",
	}

	for idStr, expectedName := range expectedSpecialties {
		id := uuid.MustParse(idStr)
		specialty, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, specialty)
		assert.Equal(t, expectedName, specialty.Name)
		assert.Equal(t, id, specialty.ID)
	}
}
