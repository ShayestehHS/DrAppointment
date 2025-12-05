//go:build test
// +build test

package doctor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	medicalService "github.com/shayesteh1hs/DrAppointment/internal/service/medical/doctor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/doctor/memory"
)

type DoctorOffsetPageDTO = pagination.Result[ListItemDTO]

func setupDoctorHandler() *Handler {
	repo := memory.NewDoctorRepositoryWithTestData()
	service := medicalService.NewDoctorService(repo)
	return NewHandler(service)
}

func TestDoctorHandler_ListDoctors_Success(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	req, err := http.NewRequest("GET", "/doctors?page=1&limit=10", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DoctorOffsetPageDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 2, response.TotalCount)
	assert.Len(t, response.Items, 2)

	// Verify the first doctor (should be sorted by created_at desc)
	firstDoctor := response.Items[0]
	assert.Equal(t, "Dr. Jane Doe", firstDoctor.Name)
	assert.NotEmpty(t, firstDoctor.ID)
	assert.NotEmpty(t, firstDoctor.SpecialtyID)
}

func TestDoctorHandler_ListDoctors_WithPagination(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Test with page=2, limit=1 to get second doctor
	req, err := http.NewRequest("GET", "/doctors?page=2&limit=1", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DoctorOffsetPageDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 2, response.TotalCount)
	assert.Len(t, response.Items, 1)

	// Should get the second doctor (Dr. John Smith)
	doctor := response.Items[0]
	assert.Equal(t, "Dr. John Smith", doctor.Name)
}

func TestDoctorHandler_ListDoctors_InvalidPagination(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Test with invalid limit
	req, err := http.NewRequest("GET", "/doctors?page=1&limit=invalid", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "invalid pagination parameters")
}

func TestDoctorHandler_GetDoctorByID_Success(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Use the known doctor ID from memory repository
	doctorID := "123e4567-e89b-12d3-a456-426614174000"
	req, err := http.NewRequest("GET", "/doctors/"+doctorID, nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DetailDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Dr. John Smith", response.Name)
	assert.Equal(t, doctorID, response.ID.String())
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
}

func TestDoctorHandler_GetDoctorByID_NotFound(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Use a non-existent doctor ID
	nonExistentID := uuid.New().String()
	req, err := http.NewRequest("GET", "/doctors/"+nonExistentID, nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Doctor not found")
}

func TestDoctorHandler_GetDoctorByID_InvalidID(t *testing.T) {
	handler := setupDoctorHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	req, err := http.NewRequest("GET", "/doctors/invalid-uuid", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Invalid doctor ID")
}
