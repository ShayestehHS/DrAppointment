//go:build test
// +build test

package memory

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor"
)

type doctorRepository struct {
	doctors []medical.Doctor
}

func (r *doctorRepository) ListOffset(ctx context.Context, filters filter.DoctorQueryParam, params pagination.LimitOffsetParams) ([]medical.Doctor, error) {
	filteredDoctors := r.applyFilters(r.doctors, filters)

	sort.Slice(filteredDoctors, func(i, j int) bool {
		return filteredDoctors[i].CreatedAt.After(filteredDoctors[j].CreatedAt)
	})

	offset := params.GetOffset()
	start := offset
	if start < 0 {
		start = 0
	}
	total := len(filteredDoctors)
	if start > total {
		start = total
	}

	end := start + params.Limit
	if params.Limit < 0 {
		end = start
	}
	if end > total {
		end = total
	}

	items := make([]medical.Doctor, 0)
	if start < total && end > start {
		items = filteredDoctors[start:end]
	}

	return items, nil
}

func (r *doctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*medical.Doctor, error) {
	for i := range r.doctors {
		if r.doctors[i].ID == id {
			return &r.doctors[i], nil
		}
	}
	return nil, doctor.ErrDoctorNotFound
}

func (r *doctorRepository) Count(ctx context.Context, filters filter.DoctorQueryParam) (int, error) {
	filteredDoctors := r.applyFilters(r.doctors, filters)
	return len(filteredDoctors), nil
}

func (r *doctorRepository) applyFilters(doctors []medical.Doctor, filters filter.DoctorQueryParam) []medical.Doctor {
	var filtered []medical.Doctor

	for _, doc := range doctors {
		match := true

		trimmedName := strings.TrimSpace(filters.Name)
		if trimmedName != "" {
			if !strings.Contains(strings.ToLower(doc.Name), strings.ToLower(trimmedName)) {
				match = false
			}
		}

		if filters.SpecialtyID != uuid.Nil {
			if doc.SpecialtyID != filters.SpecialtyID {
				match = false
			}
		}

		if match {
			filtered = append(filtered, doc)
		}
	}

	return filtered
}

func (r *doctorRepository) AddDoctor(doc medical.Doctor) {
	r.doctors = append(r.doctors, doc)
}

func (r *doctorRepository) Clear() {
	r.doctors = []medical.Doctor{}
}

func NewDoctorRepository() doctor.Repository {
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
		},
	}
}
