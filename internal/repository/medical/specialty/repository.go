package specialty

import (
	"context"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

type Repository interface {
	ListOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error)
}
