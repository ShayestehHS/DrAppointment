package specialty

import (
	"context"

	"github.com/google/uuid"

	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty"
)

type Service interface {
	ListSpecialtiesOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error)
}

type specialtyService struct {
	repo specialty.Repository
}

func NewSpecialtyService(repo specialty.Repository) Service {
	return &specialtyService{
		repo: repo,
	}
}

func (s *specialtyService) ListSpecialtiesOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, int, error) {
	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		return []medical.Specialty{}, 0, err
	}

	specialties, err := s.repo.ListOffset(ctx, params)
	return specialties, totalCount, err
}

func (s *specialtyService) GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error) {
	return s.repo.GetByID(ctx, id)
}
