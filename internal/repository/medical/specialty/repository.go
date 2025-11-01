package specialty

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
)

var ErrSpecialtyNotFound = errors.New("specialty not found")

type Repository interface {
	ListOffset(ctx context.Context, params pagination.LimitOffsetParams) ([]medical.Specialty, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*medical.Specialty, error)
}
