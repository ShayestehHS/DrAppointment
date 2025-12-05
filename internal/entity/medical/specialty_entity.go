package medical

import (
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/entity"
)

var _ entity.Image = &SpecialtyImage{}
var _ entity.ModelEntity = &Specialty{}

type Specialty struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	ImagePath *SpecialtyImage `json:"image_path" db:"-"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

func (s Specialty) GetPK() string {
	return s.ID.String()
}

type SpecialtyImage struct {
	Path url.URL `json:"path" db:"image_path"`
}

func (si *SpecialtyImage) GetPath() *url.URL {
	return &si.Path
}

func NewSpecialtyImage(fileName string) *SpecialtyImage {
	if fileName == "" {
		return nil
	}

	path := "specialties/" + filepath.Base(fileName)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return &SpecialtyImage{Path: url.URL{
		Path: path,
	}}
}
