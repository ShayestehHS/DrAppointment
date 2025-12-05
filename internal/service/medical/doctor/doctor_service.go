package doctor

import (
	"context"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"

	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor"
)

type Service interface {
	ListDoctorsOffset(ctx context.Context, filters filter.DoctorQueryParam, params pagination.LimitOffsetParams) ([]medical.Doctor, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*medical.Doctor, error)
}

type doctorService struct {
	repo doctor.Repository
}

func NewDoctorService(repo doctor.Repository) Service {
	return &doctorService{
		repo: repo,
	}
}

func (s *doctorService) ListDoctorsOffset(ctx context.Context, filters filter.DoctorQueryParam, params pagination.LimitOffsetParams) ([]medical.Doctor, int, error) {
	totalCount, err := s.repo.Count(ctx, filters)
	if err != nil {
		return []medical.Doctor{}, 0, err
	}

	doctors, err := s.repo.ListOffset(ctx, filters, params)
	return doctors, totalCount, err
}

func (s *doctorService) GetByID(ctx context.Context, id uuid.UUID) (*medical.Doctor, error) {
	return s.repo.GetByID(ctx, id)
}
