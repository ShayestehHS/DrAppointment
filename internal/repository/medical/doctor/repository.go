package doctor

import (
	"context"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	filter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

type Repository interface {
	ListOffset(ctx context.Context, filters filter.DoctorQueryParam, params pagination.LimitOffsetParams) ([]medical.Doctor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*medical.Doctor, error)
	Count(ctx context.Context, filters filter.DoctorQueryParam) (int, error)
}
