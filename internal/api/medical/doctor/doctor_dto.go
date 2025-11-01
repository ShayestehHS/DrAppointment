package doctor

import (
	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/api"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
)

var _ api.PageEntityDTO = (*ListItemDTO)(nil)

type ListItemDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SpecialtyID uuid.UUID `json:"specialty_id"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Description string    `json:"description,omitempty"`
}

func newListItemDTO(doctors []medical.Doctor) []ListItemDTO {
	items := make([]ListItemDTO, 0, len(doctors))

	for _, doctor := range doctors {
		items = append(items, ListItemDTO{
			ID:          doctor.ID,
			Name:        doctor.Name,
			SpecialtyID: doctor.SpecialtyID,
			PhoneNumber: doctor.PhoneNumber,
			AvatarURL:   doctor.AvatarURL,
			Description: doctor.Description,
		})
	}

	return items
}

func (d ListItemDTO) IsPageEntityDTO() bool { return true }
func (d ListItemDTO) GetID() string         { return d.ID.String() }

type DetailDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SpecialtyID uuid.UUID `json:"specialty_id"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func NewDetailDTO(doctor medical.Doctor) DetailDTO {
	return DetailDTO{
		ID:          doctor.ID,
		Name:        doctor.Name,
		SpecialtyID: doctor.SpecialtyID,
		PhoneNumber: doctor.PhoneNumber,
		AvatarURL:   doctor.AvatarURL,
		Description: doctor.Description,
		CreatedAt:   doctor.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   doctor.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
