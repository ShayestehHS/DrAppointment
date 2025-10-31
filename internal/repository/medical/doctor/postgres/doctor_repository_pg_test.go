package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestDoctorPostgresRepository_ListOffset_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	// Mock select query
	doctorID1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	doctorID2 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
	specialtyID1 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	specialtyID2 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}).
			AddRow(doctorID1, "Dr. John Smith", specialtyID1, "+1234567890", "https://example.com/avatar1.jpg", "Experienced cardiologist", now, now).
			AddRow(doctorID2, "Dr. Jane Doe", specialtyID2, "+1234567891", "https://example.com/avatar2.jpg", "Skilled neurologist", now, now))

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 2)
	assert.Equal(t, "Dr. John Smith", result[0].Name)
	assert.Equal(t, "Dr. Jane Doe", result[1].Name)

	// Verify all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_ListOffset_WithNameFilter(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		Name: "John",
	}

	// Mock select query with name filter
	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors WHERE name LIKE $1 LIMIT $2 OFFSET $3")).
		WithArgs("%John%", 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}).
			AddRow(doctorID, "Dr. John Smith", specialtyID, "+1234567890", "https://example.com/avatar1.jpg", "Experienced cardiologist", now, now))

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, "Dr. John Smith", result[0].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_ListOffset_WithSpecialtyFilter(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{
		SpecialtyID: specialtyID,
	}

	// Mock select query with specialty filter
	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors WHERE specialty_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(specialtyID.String(), 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}).
			AddRow(doctorID, "Dr. John Smith", specialtyID, "+1234567890", "https://example.com/avatar1.jpg", "Experienced cardiologist", now, now))

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, "Dr. John Smith", result[0].Name)
	assert.Equal(t, specialtyID, result[0].SpecialtyID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_ListOffset_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	// Mock select query returning empty result
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}))

	result, err := repo.ListOffset(ctx, filters, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_ListOffset_CountError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}
	filters := filter.DoctorQueryParam{}

	// Mock select query with error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.ListOffset(ctx, filters, params)
	require.Error(t, err)
	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_GetByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors WHERE id = $1")).
		WithArgs(doctorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}).
			AddRow(doctorID, "Dr. John Smith", specialtyID, "+1234567890", "https://example.com/avatar1.jpg", "Experienced cardiologist", now, now))

	doctor, err := repo.GetByID(ctx, doctorID)
	require.NoError(t, err)
	require.NotNil(t, doctor)

	assert.Equal(t, doctorID, doctor.ID)
	assert.Equal(t, "Dr. John Smith", doctor.Name)
	assert.Equal(t, specialtyID, doctor.SpecialtyID)
	assert.Equal(t, "+1234567890", doctor.PhoneNumber)
	assert.Equal(t, "Experienced cardiologist", doctor.Description)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors WHERE id = $1")).
		WithArgs(doctorID).
		WillReturnError(sql.ErrNoRows)

	doctor, err := repo.GetByID(ctx, doctorID)
	require.Error(t, err)
	assert.Nil(t, doctor)
	assert.Contains(t, err.Error(), "doctor not found")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDoctorPostgresRepository_GetByID_ScanError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewDoctorRepository(db)
	ctx := context.Background()

	doctorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	// Return invalid data that will cause scan error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, specialty_id, phone_number, avatar_url, description, created_at, updated_at FROM doctors WHERE id = $1")).
		WithArgs(doctorID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at"}).
			AddRow("invalid-uuid", "Dr. John Smith", "invalid-uuid", "+1234567890", "https://example.com/avatar1.jpg", "Experienced cardiologist", "invalid-time", "invalid-time"))

	doctor, err := repo.GetByID(ctx, doctorID)
	require.Error(t, err)
	assert.Nil(t, doctor)

	require.NoError(t, mock.ExpectationsWereMet())
}
