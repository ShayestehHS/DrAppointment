//go:build test
// +build test

package doctor

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor/memory"
)

func setupDoctorService() Service {
	repo := memory.NewDoctorRepositoryWithTestData()
	return NewDoctorService(repo)
}

func TestDoctorService_ListDoctorsOffset_Success(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Equal(t, 2, totalCount)
	assert.Len(t, doctors, 2)

	// Verify doctors are sorted by created_at desc (newer first)
	assert.Equal(t, "Dr. Jane Doe", doctors[0].Name)
	assert.Equal(t, "Dr. John Smith", doctors[1].Name)
}

func TestDoctorService_ListDoctorsOffset_WithPagination(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  2,
		Limit: 1,
	}
	filters := filter.DoctorQueryParam{}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Equal(t, 2, totalCount)
	assert.Len(t, doctors, 1)

	// Should get the second doctor (Dr. John Smith)
	assert.Equal(t, "Dr. John Smith", doctors[0].Name)
}

func TestDoctorService_ListDoctorsOffset_ValidatesLimitTooLow(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 0, // Invalid limit
	}
	filters := filter.DoctorQueryParam{}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	// Service should still return results (limit validation happens in service)
	assert.Equal(t, 2, totalCount)
	assert.Len(t, doctors, 0)
}

func TestDoctorService_ListDoctorsOffset_ValidatesLimitTooHigh(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 200, // Too high limit
	}
	filters := filter.DoctorQueryParam{}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	// Service should still return results (limit validation happens in service)
	assert.Equal(t, 2, totalCount)
	assert.Len(t, doctors, 2)
}

func TestDoctorService_ListDoctorsOffset_Success_FirstPage(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	// Service should still return results
	assert.Equal(t, 2, totalCount)
	assert.Len(t, doctors, 2)
}

func TestDoctorService_ListDoctorsOffset_WithNameFilter(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		Name: "Jane",
	}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Equal(t, 1, totalCount)
	assert.Len(t, doctors, 1)
	assert.Equal(t, "Dr. Jane Doe", doctors[0].Name)
}

func TestDoctorService_ListDoctorsOffset_WithSpecialtyFilter(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		SpecialtyID: uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
	}

	doctors, totalCount, err := service.ListDoctorsOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Equal(t, 1, totalCount)
	assert.Len(t, doctors, 1)
	assert.Equal(t, "Dr. John Smith", doctors[0].Name)
}

func TestDoctorService_GetByID_Success(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	doctor, err := service.GetByID(ctx, doctorID)
	require.NoError(t, err)
	require.NotNil(t, doctor)

	assert.Equal(t, doctorID, doctor.ID)
	assert.Equal(t, "Dr. John Smith", doctor.Name)
}

func TestDoctorService_GetByID_NotFound(t *testing.T) {
	service := setupDoctorService()
	ctx := context.Background()

	nonExistentID := uuid.New()

	doc, err := service.GetByID(ctx, nonExistentID)
	require.Error(t, err)
	assert.Nil(t, doc)
	assert.True(t, errors.Is(err, doctor.ErrDoctorNotFound))
}
