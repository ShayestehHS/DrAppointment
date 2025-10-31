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

	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func TestSpecialtyPostgresRepository_ListOffset_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	// Mock select query
	specialtyID1 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	specialtyID2 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	specialtyID3 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174002")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
			AddRow(specialtyID1, "Cardiology", "cardiology.jpg", now, now).
			AddRow(specialtyID2, "Neurology", "neurology.jpg", now, now).
			AddRow(specialtyID3, "Dermatology", sql.NullString{}, now, now)) // Test null image path

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 3)
	assert.Equal(t, "Cardiology", result[0].Name)
	assert.Equal(t, "Neurology", result[1].Name)
	assert.Equal(t, "Dermatology", result[2].Name)

	// Check image paths
	assert.NotNil(t, result[0].ImagePath)
	assert.NotNil(t, result[1].ImagePath)
	assert.Nil(t, result[2].ImagePath) // Should be nil for null image path

	// Verify all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_ListOffset_WithPagination(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  2,
		Limit: 2,
	}

	// Mock select query with pagination
	specialtyID1 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	specialtyID2 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174002")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(2, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
			AddRow(specialtyID1, "Neurology", "neurology.jpg", now, now).
			AddRow(specialtyID2, "Dermatology", "dermatology.jpg", now, now))

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 2)
	assert.Equal(t, "Neurology", result[0].Name)
	assert.Equal(t, "Dermatology", result[1].Name)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_ListOffset_EmptyResult(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	// Mock select query returning empty result
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}))

	result, err := repo.ListOffset(ctx, params)
	require.NoError(t, err)

	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_ListOffset_CountError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	// Mock select query with error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.ListOffset(ctx, params)
	require.Error(t, err)
	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_ListOffset_SelectError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	// Mock select query with error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnError(sql.ErrConnDone)

	result, err := repo.ListOffset(ctx, params)
	require.Error(t, err)
	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_GetByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties WHERE id = $1")).
		WithArgs(specialtyID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
			AddRow(specialtyID, "Cardiology", "cardiology.jpg", now, now))

	specialty, err := repo.GetByID(ctx, specialtyID)
	require.NoError(t, err)
	require.NotNil(t, specialty)

	assert.Equal(t, specialtyID, specialty.ID)
	assert.Equal(t, "Cardiology", specialty.Name)
	assert.NotNil(t, specialty.ImagePath)
	assert.Contains(t, specialty.ImagePath.GetPath().String(), "cardiology.jpg")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_GetByID_WithNullImagePath(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties WHERE id = $1")).
		WithArgs(specialtyID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
			AddRow(specialtyID, "Cardiology", sql.NullString{}, now, now))

	specialty, err := repo.GetByID(ctx, specialtyID)
	require.NoError(t, err)
	require.NotNil(t, specialty)

	assert.Equal(t, specialtyID, specialty.ID)
	assert.Equal(t, "Cardiology", specialty.Name)
	assert.Nil(t, specialty.ImagePath) // Should be nil for null image path

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties WHERE id = $1")).
		WithArgs(specialtyID).
		WillReturnError(sql.ErrNoRows)

	specialty, err := repo.GetByID(ctx, specialtyID)
	require.Error(t, err)
	assert.Nil(t, specialty)
	assert.Contains(t, err.Error(), "specialty not found")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_GetByID_ScanError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")

	// Return invalid data that will cause scan error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties WHERE id = $1")).
		WithArgs(specialtyID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
			AddRow("invalid-uuid", "Cardiology", "cardiology.jpg", "invalid-time", "invalid-time"))

	specialty, err := repo.GetByID(ctx, specialtyID)
	require.Error(t, err)
	assert.Nil(t, specialty)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSpecialtyPostgresRepository_ScanSpecialties_RowsError(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewSpecialtyRepository(db)
	ctx := context.Background()

	params := pagination.LimitOffsetParams{
		Page:  1,
		Limit: 10,
	}

	// Mock select query with rows that will have an error during iteration
	specialtyID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "image_path", "created_at", "updated_at"}).
		AddRow(specialtyID, "Cardiology", "cardiology.jpg", now, now).
		RowError(0, sql.ErrConnDone) // Add error to first row

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image_path, created_at, updated_at FROM specialties LIMIT $1 OFFSET $2")).
		WithArgs(10, 0).
		WillReturnRows(rows)

	result, err := repo.ListOffset(ctx, params)
	require.Error(t, err)
	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}
