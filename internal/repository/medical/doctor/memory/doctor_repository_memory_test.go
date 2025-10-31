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
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

func setupDoctorMemoryRepo() *doctorRepository {
	return &doctorRepository{
		doctors: []medical.Doctor{
			{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name:        "Dr. John Smith",
				SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
				PhoneNumber: "+1234567890",
				AvatarURL:   "https://example.com/avatar1.jpg",
				Description: "Experienced cardiologist",
				CreatedAt:   time.Now().Add(-24 * time.Hour),
				UpdatedAt:   time.Now().Add(-24 * time.Hour),
			},
			{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
				Name:        "Dr. Jane Doe",
				SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				PhoneNumber: "+1234567891",
				AvatarURL:   "https://example.com/avatar2.jpg",
				Description: "Skilled neurologist",
				CreatedAt:   time.Now().Add(-12 * time.Hour),
				UpdatedAt:   time.Now().Add(-12 * time.Hour),
			},
			{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
				Name:        "Dr. Alice Johnson",
				SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
				PhoneNumber: "+1234567892",
				AvatarURL:   "https://example.com/avatar3.jpg",
				Description: "Expert cardiologist",
				CreatedAt:   time.Now().Add(-6 * time.Hour),
				UpdatedAt:   time.Now().Add(-6 * time.Hour),
			},
		},
	}
}

func TestDoctorMemoryRepository_ListOffset_Success(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 3)

	// Verify doctors are sorted by created_at desc (newer first)
	assert.Equal(t, "Dr. Alice Johnson", result[0].Name)
	assert.Equal(t, "Dr. Jane Doe", result[1].Name)
	assert.Equal(t, "Dr. John Smith", result[2].Name)
}

func TestDoctorMemoryRepository_ListOffset_WithPagination(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  2,
		Limit: 2,
	}
	filters := filter.DoctorQueryParam{}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 1)

	// Page 2 with limit 2 = offset 2, should get only the last doctor
	assert.Equal(t, "Dr. John Smith", result[0].Name)
}

func TestDoctorMemoryRepository_ListOffset_WithNameFilter(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		Name: "Jane",
	}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, "Dr. Jane Doe", result[0].Name)
}

func TestDoctorMemoryRepository_ListOffset_WithSpecialtyFilter(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
	}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 2)

	// Should get both cardiologists (Alice and John)
	assert.Equal(t, "Dr. Alice Johnson", result[0].Name)
	assert.Equal(t, "Dr. John Smith", result[1].Name)
}

func TestDoctorMemoryRepository_ListOffset_EmptyResult(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  100,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)
}

func TestDoctorMemoryRepository_ListOffset_NoMatchingFilter(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		Name: "NonExistent",
	}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)
}

func TestDoctorMemoryRepository_GetByID_Success(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	doctor, err := repo.GetByID(ctx, doctorID)
	require.NoError(t, err)
	require.NotNil(t, doctor)

	assert.Equal(t, doctorID, doctor.ID)
	assert.Equal(t, "Dr. John Smith", doctor.Name)
	assert.Equal(t, "Experienced cardiologist", doctor.Description)
}

func TestDoctorMemoryRepository_GetByID_NotFound(t *testing.T) {
	repo := setupDoctorMemoryRepo()
	ctx := context.Background()

	nonExistentID := uuid.New()

	doctor, err := repo.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	assert.Nil(t, doctor)
	assert.Contains(t, err.Error(), "doctor not found")
}

func TestDoctorMemoryRepository_AddDoctor(t *testing.T) {
	repo := setupDoctorMemoryRepo()

	newDoctor := medical.Doctor{
		ID:          uuid.New(),
		Name:        "Dr. Bob Wilson",
		SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174002"),
		PhoneNumber: "+1234567893",
		AvatarURL:   "https://example.com/avatar4.jpg",
		Description: "Dermatologist",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.AddDoctor(newDoctor)

	// Verify the doctor was added
	ctx := context.Background()
	params := pagination.LimitOffsetParams{Page: 1, Limit: 10}
	filters := filter.DoctorQueryParam{}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 4)

	// The new doctor should be first (most recent)
	assert.Equal(t, "Dr. Bob Wilson", result[0].Name)
}

func TestDoctorMemoryRepository_Clear(t *testing.T) {
	repo := setupDoctorMemoryRepo()

	repo.Clear()

	// Verify all doctors were removed
	ctx := context.Background()
	params := pagination.LimitOffsetParams{Page: 1, Limit: 10}
	filters := filter.DoctorQueryParam{}

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)
}
