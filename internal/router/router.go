package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/shayesteh1hs/DrAppointment/internal/api/medical/doctor"
	"github.com/shayesteh1hs/DrAppointment/internal/api/medical/specialty"
	doctor2 "github.com/shayesteh1hs/DrAppointment/internal/service/medical/doctor"
	medicalService "github.com/shayesteh1hs/DrAppointment/internal/service/medical/specialty"

	"github.com/shayesteh1hs/DrAppointment/internal/middleware"
	doctorPostgres "github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor/postgres"
	specialtyPostgres "github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty/postgres"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to DrGo API",
		})
	})

	api.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Setup medical group routes
	medicalGroup := api.Group("/medical")
	setupMedicalRoutes(medicalGroup, db)

	return r
}

func setupMedicalRoutes(rg *gin.RouterGroup, db *sql.DB) {
	// Setup doctor routes
	doctorRepo := doctorPostgres.NewDoctorRepository(db)
	doctorService := doctor2.NewDoctorService(doctorRepo)
	doctorHandler := doctor.NewHandler(doctorService)
	doctorHandler.RegisterRoutes(rg)

	// Setup specialty routes
	specialtyRepo := specialtyPostgres.NewSpecialtyRepository(db)
	specialtyService := medicalService.NewSpecialtyService(specialtyRepo)
	specialtyHandler := specialty.NewSpecialtyHandler(specialtyService)
	specialtyHandler.RegisterRoutes(rg)
}
