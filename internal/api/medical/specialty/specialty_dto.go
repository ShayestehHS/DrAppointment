package specialty

import (
	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/api"
	"github.com/shayesteh1hs/DrAppointment/internal/entity/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/utils"
)

var _ api.PageEntityDTO = (*ListItemDTO)(nil)

type ListItemDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ImageURL *string   `json:"image_url"`
}

func (p ListItemDTO) IsPageEntityDTO() bool { return true }
func (p ListItemDTO) GetID() string         { return p.ID.String() }

func NewListItemDTO(specialties []medical.Specialty) []ListItemDTO {
	items := make([]ListItemDTO, 0, len(specialties))
	for _, specialty := range specialties {
		dto := ListItemDTO{
			ID:   specialty.ID,
			Name: specialty.Name,
		}
		if specialty.ImagePath != nil {
			s := utils.GetFullImageURL(specialty.ImagePath)
			dto.ImageURL = &s
		}
		items = append(items, dto)
	}
	return items
}

type DetailDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt string    `json:"created_at"`
}

func NewDetailDTO(specialty medical.Specialty) DetailDTO {
	dto := DetailDTO{
		ID:        specialty.ID,
		Name:      specialty.Name,
		CreatedAt: specialty.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if specialty.ImagePath != nil {
		s := utils.GetFullImageURL(specialty.ImagePath)
		dto.ImageURL = &s
	}

	return dto
}
