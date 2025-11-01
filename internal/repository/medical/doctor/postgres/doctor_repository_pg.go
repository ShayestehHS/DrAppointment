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
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor"
)

type doctorRepository struct {
	db *sql.DB
}

func (r *doctorRepository) ListOffset(ctx context.Context, filters filter.DoctorQueryParam, params pagination.LimitOffsetParams) ([]medical.Doctor, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at")
	sb.From("doctors")
	sb = filters.Apply(sb)
	sb.Limit(params.Limit)
	sb.Offset(params.GetOffset())

	query, args := sb.Build()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []medical.Doctor{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	doctors, err := r.scanDoctors(rows)
	if err != nil {
		return []medical.Doctor{}, err
	}

	return doctors, nil
}

func (r *doctorRepository) Count(ctx context.Context, filters filter.DoctorQueryParam) (int, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("count(*)")
	sb.From("doctors")
	filters.Apply(sb)

	query, args := sb.Build()
	var totalCount int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to scan total count: %w", err)
	}
	return totalCount, nil
}

func (r *doctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*medical.Doctor, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "name", "specialty_id", "phone_number", "avatar_url", "description", "created_at", "updated_at")
	sb.From("doctors")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	row := r.db.QueryRowContext(ctx, query, args...)

	var doc medical.Doctor
	err := row.Scan(
		&doc.ID,
		&doc.Name,
		&doc.SpecialtyID,
		&doc.PhoneNumber,
		&doc.AvatarURL,
		&doc.Description,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, doctor.ErrDoctorNotFound
		}
		return nil, fmt.Errorf("failed to scan doctor: %w", err)
	}

	return &doc, nil
}

func (r *doctorRepository) scanDoctors(rows *sql.Rows) ([]medical.Doctor, error) {
	var doctors []medical.Doctor
	for rows.Next() {
		var doc medical.Doctor
		err := rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.SpecialtyID,
			&doc.PhoneNumber,
			&doc.AvatarURL,
			&doc.Description,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func NewDoctorRepository(db *sql.DB) doctor.Repository {
	return &doctorRepository{db: db}
}
