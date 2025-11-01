package specialty

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty"
	medicalService "github.com/shayesteh1hs/DrAppointment/internal/service/medical/specialty"
)

type SpecialtyHandler struct {
	service medicalService.Service
}

func NewSpecialtyHandler(service medicalService.Service) *SpecialtyHandler {
	return &SpecialtyHandler{
		service: service,
	}
}

func (h *SpecialtyHandler) ListSpecialties(c *gin.Context) {
	paginator := pagination.NewOffsetPaginator[ListItemDTO]()
	if err := paginator.BindQueryParam(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	specialties, totalCount, err := h.service.ListSpecialtiesOffset(c.Request.Context(), paginator.GetParams())
	if err != nil {
		log.Printf("failed to fetch specialties: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch specialties"})
		return
	}

	specialtyDTOs := NewListItemDTO(specialties)
	result, err := paginator.CreatePaginationResult(specialtyDTOs, totalCount)
	if err != nil {
		log.Printf("failed to create pagination result: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pagination result"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SpecialtyHandler) GetSpecialtyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialty ID"})
		return
	}

	spec, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, specialty.ErrSpecialtyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Specialty not found"})
			return
		}
		log.Printf("failed to fetch specialty by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch specialty"})
		return
	}

	response := NewDetailDTO(*spec)
	c.JSON(http.StatusOK, response)
}

func (h *SpecialtyHandler) RegisterRoutes(router *gin.RouterGroup) {
	specialtyRoutes := router.Group("/specialties")
	{
		specialtyRoutes.GET("", h.ListSpecialties)
		specialtyRoutes.GET("/:id", h.GetSpecialtyByID)
	}
}
