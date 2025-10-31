//go:build test
// +build test

package memory

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty"
)

type specialtyRepository struct {
	specialties []medical.Specialty
}

func (r *specialtyRepository) ListOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, error) {
	// Sort by created_at desc for consistent ordering
	sortedSpecialties := make([]medical.Specialty, len(r.specialties))
	copy(sortedSpecialties, r.specialties)
	sort.Slice(sortedSpecialties, func(i, j int) bool {
		return sortedSpecialties[i].CreatedAt.After(sortedSpecialties[j].CreatedAt)
	})

	total := len(sortedSpecialties)

	// Apply pagination
	offset := params.GetOffset()
	start := offset
	if start < 0 {
		start = 0
	}
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

	items := make([]medical.Specialty, 0)
	if start < total && end > start {
		items = sortedSpecialties[start:end]
	}

	return items, nil
}

func (r *specialtyRepository) GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error) {
	for _, spec := range r.specialties {
		if spec.ID == id {
			return &spec, nil
		}
	}
	return nil, fmt.Errorf("specialty not found: %s", id)
}

func (r *specialtyRepository) Count(ctx context.Context) (int, error) {
	return len(r.specialties), nil
}

// AddSpecialty adds a specialty to the in-memory store (for testing)
func (r *specialtyRepository) AddSpecialty(spec medical.Specialty) {
	r.specialties = append(r.specialties, spec)
}

// Clear removes all specialties from the in-memory store (for testing)
func (r *specialtyRepository) Clear() {
	r.specialties = []medical.Specialty{}
}

func NewSpecialtyRepository() specialty.Repository {
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
		},
	}
}
