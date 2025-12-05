package medical

import (
	"time"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity"
)

var _ entity.ModelEntity = (*Doctor)(nil)

type Doctor struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	SpecialtyID uuid.UUID `json:"specialty_id" db:"specialty_id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (d Doctor) GetPK() string {
	return d.ID.String()
}
