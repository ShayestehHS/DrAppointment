//go:build test
// +build test

package specialty

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	medicalService "github.com/shayesteh1hs/DrAppointment/internal/service/medical/specialty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/shayesteh1hs/DrAppointment/internal/pagination"
	"github.com/shayesteh1hs/DrAppointment/internal/repository/medical/specialty/memory"
)

type SpecialtyOffsetPageDTO = pagination.Result[ListItemDTO]

func setupSpecialtyHandler() *SpecialtyHandler {
	repo := memory.NewSpecialtyRepository()
	service := medicalService.NewSpecialtyService(repo)
	return NewSpecialtyHandler(service)
}

func TestSpecialtyHandler_ListSpecialties_Success(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	req, err := http.NewRequest("GET", "/specialties?page=1&limit=10", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response SpecialtyOffsetPageDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 3, response.TotalCount)
	assert.Len(t, response.Items, 3)

	// Verify the first specialty (should be sorted by created_at desc)
	firstSpecialty := response.Items[0]
	assert.Equal(t, "Dermatology", firstSpecialty.Name)
	assert.NotEmpty(t, firstSpecialty.ID)
	assert.NotEmpty(t, firstSpecialty.ImageURL)
}

func TestSpecialtyHandler_ListSpecialties_WithPagination(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Test with page=2, limit=2 to get second and third specialties
	req, err := http.NewRequest("GET", "/specialties?page=2&limit=2", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response SpecialtyOffsetPageDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 3, response.TotalCount)
	assert.Len(t, response.Items, 1)

	// Should get Cardiology (only 1 item on page 2 with limit 2)
	assert.Equal(t, "Cardiology", response.Items[0].Name)
}

func TestSpecialtyHandler_ListSpecialties_EmptyResult(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Test with page beyond available data
	req, err := http.NewRequest("GET", "/specialties?page=100&limit=10", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response SpecialtyOffsetPageDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 3, response.TotalCount)
	assert.Len(t, response.Items, 0)
}

func TestSpecialtyHandler_ListSpecialties_InvalidPagination(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Test with invalid page
	req, err := http.NewRequest("GET", "/specialties?page=invalid&limit=10", nil)
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

func TestSpecialtyHandler_GetSpecialtyByID_Success(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Use the known specialty ID from memory repository
	specialtyID := "223e4567-e89b-12d3-a456-426614174000"
	req, err := http.NewRequest("GET", "/specialties/"+specialtyID, nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DetailDTO
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Cardiology", response.Name)
	assert.Equal(t, specialtyID, response.ID.String())
	assert.NotEmpty(t, response.ImageURL)
	assert.NotEmpty(t, response.CreatedAt)
}

func TestSpecialtyHandler_GetSpecialtyByID_NotFound(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	// Use a non-existent specialty ID
	nonExistentID := uuid.New().String()
	req, err := http.NewRequest("GET", "/specialties/"+nonExistentID, nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Specialty not found")
}

func TestSpecialtyHandler_GetSpecialtyByID_InvalidID(t *testing.T) {
	handler := setupSpecialtyHandler()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router.Group("/"))

	req, err := http.NewRequest("GET", "/specialties/invalid-uuid", nil)
	require.NoError(t, err)
	req.Host = "localhost:8080"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["error"], "Invalid specialty ID")
}
