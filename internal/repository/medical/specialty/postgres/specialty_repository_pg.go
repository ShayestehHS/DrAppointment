package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"

	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty"
)

type specialtyRepository struct {
	db *sql.DB
}

func (r *specialtyRepository) ListOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "image_path", "created_at", "updated_at")
	sb.From("specialties")
	sb.Limit(params.Limit)
	sb.Offset(params.GetOffset())

	query, args := sb.Build()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []medical.Specialty{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	specialties, err := r.scanSpecialties(rows)
	if err != nil {
		return []medical.Specialty{}, err
	}

	return specialties, nil
}

func (r *specialtyRepository) Count(ctx context.Context) (int, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("count(*)")
	sb.From("specialties")

	query, args := sb.Build()
	var totalCount int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to scan total count: %w", err)
	}
	return totalCount, nil
}

func (r *specialtyRepository) GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "image_path", "created_at", "updated_at")
	sb.From("specialties")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	row := r.db.QueryRowContext(ctx, query, args...)

	var spec medical.Specialty
	var imagePath sql.NullString
	err := row.Scan(
		&spec.ID,
		&spec.Name,
		&imagePath,
		&spec.CreatedAt,
		&spec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, specialty.ErrSpecialtyNotFound
		}
		return nil, fmt.Errorf("failed to scan specialty: %w", err)
	}

	if imagePath.Valid {
		spec.ImagePath = medical.NewSpecialtyImage(imagePath.String)
	}

	return &spec, nil
}

func (r *specialtyRepository) scanSpecialties(rows *sql.Rows) ([]medical.Specialty, error) {
	var specialties []medical.Specialty
	for rows.Next() {
		var spec medical.Specialty
		var imagePath sql.NullString
		err := rows.Scan(
			&spec.ID,
			&spec.Name,
			&imagePath,
			&spec.CreatedAt,
			&spec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if imagePath.Valid {
			spec.ImagePath = medical.NewSpecialtyImage(imagePath.String)
		}

		specialties = append(specialties, spec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return specialties, nil
}

func NewSpecialtyRepository(db *sql.DB) specialty.Repository {
	return &specialtyRepository{db: db}
}
