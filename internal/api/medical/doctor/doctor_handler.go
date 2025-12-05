package doctor

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	medicalFilter "github.com/shayesteh1hs/DrAppointment/internal/filter/medical"
	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor"
	medicalService "github.com/shayesteh1hs/DrAppointment/internal/service/medical/doctor"
)

type Handler struct {
	service medicalService.Service
}

func NewHandler(service medicalService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListDoctors(c *gin.Context) {
	paginator := pagination.NewOffsetPaginator[ListItemDTO]()
	if err := paginator.BindQueryParam(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filterParams medicalFilter.DoctorQueryParam
	if err := c.ShouldBindQuery(&filterParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}
	if err := filterParams.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctors, totalCount, err := h.service.ListDoctorsOffset(c.Request.Context(), filterParams, paginator.GetParams())
	if err != nil {
		log.Printf("failed to fetch doctors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	doctorsDTO := newListItemDTO(doctors)
	result, err := paginator.CreatePaginationResult(doctorsDTO, totalCount)
	if err != nil {
		log.Printf("failed to create pagination result: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pagination result"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetDoctorByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doc, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, doctor.ErrDoctorNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
			return
		}
		log.Printf("failed to fetch doctor by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctor"})
		return
	}

	response := NewDetailDTO(*doc)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	doctorRoutes := router.Group("/doctors")
	{
		doctorRoutes.GET("", h.ListDoctors)
		doctorRoutes.GET("/:id", h.GetDoctorByID)
	}
}
