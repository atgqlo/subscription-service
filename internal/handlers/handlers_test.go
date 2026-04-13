package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"subscriptons-service/internal/mocks"
	"subscriptons-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testID := uuid.New()
		expectedSub := &models.Subscription{ID: testID, ServiceName: "Netflix"}

		mockService.On("GetByID", mock.Anything, testID).Return(expectedSub, nil)

		r := gin.Default()
		r.GET("/subscriptions/:id", h.GetByIDHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/subscriptions/"+testID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Netflix")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testID := uuid.New()
		mockService.On("GetByID", mock.Anything, testID).Return(nil, errors.New("not found"))

		r := gin.Default()
		r.GET("/subscriptions/:id", h.GetByIDHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/subscriptions/"+testID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSubscriptionHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Valid Request", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Subscription")).Return(nil)

		r := gin.Default()
		r.POST("/subscriptions", h.CreateHandler())

		// ИСПРАВЛЕНИЕ: Теперь мы передаем JSON со всеми полями, которые требует binding:"required"
		validJSON := `{"service_name": "Spotify", "price": 300, "user_id": "123e4567-e89b-12d3-a456-426614174000", "start_date": "2024-01"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBufferString(validJSON))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		h := NewSubscriptionHandler(nil, log.New(io.Discard, "", 0))

		r := gin.Default()
		r.POST("/subscriptions", h.CreateHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBufferString("{invalid json}"))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubscriptionHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success with defaults", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		mockList := []*models.Subscription{
			{ID: uuid.New(), ServiceName: "Netflix", Price: 500},
		}
		mockService.On("List", mock.Anything, 20, 0).Return(mockList, 1, nil)

		r := gin.Default()
		r.GET("/subscriptions", h.ListHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/subscriptions", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Netflix")
		assert.Contains(t, w.Body.String(), `"total":1`)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		mockService.On("List", mock.Anything, 20, 0).Return(nil, 0, errors.New("db error"))

		r := gin.Default()
		r.GET("/subscriptions", h.ListHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/subscriptions", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSubscriptionHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testID := uuid.New()
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*models.Subscription")).Return(nil)

		r := gin.Default()
		r.PUT("/subscriptions/:id", h.UpdateHandler())

		updateJSON := `{"service_name": "Netflix Premium", "price": 1000, "start_date": "2024-01"}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/subscriptions/"+testID.String(), bytes.NewBufferString(updateJSON))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Netflix Premium")
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		h := NewSubscriptionHandler(nil, log.New(io.Discard, "", 0))

		r := gin.Default()
		r.PUT("/subscriptions/:id", h.UpdateHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/subscriptions/bad-uuid", bytes.NewBufferString(`{}`))

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubscriptionHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testID := uuid.New()
		mockService.On("Delete", mock.Anything, testID).Return(nil)

		r := gin.Default()
		r.DELETE("/subscriptions/:id", h.DeleteHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/subscriptions/"+testID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testID := uuid.New()
		mockService.On("Delete", mock.Anything, testID).Return(errors.New("db delete error"))

		r := gin.Default()
		r.DELETE("/subscriptions/:id", h.DeleteHandler())

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/subscriptions/"+testID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSubscriptionHandler_TotalCost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := mocks.NewSubscriptionServiceInterface(t)
		h := NewSubscriptionHandler(mockService, log.New(io.Discard, "", 0))

		testUserID := uuid.New()
		expectedTotal := 1500

		mockService.On("TotalCost", mock.Anything, testUserID, mock.Anything, "2024-01", "").Return(expectedTotal, nil)

		r := gin.Default()
		r.GET("/total-cost", h.TotalCostHandler())

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/total-cost?user_id=%s&start_date=2024-01", testUserID.String())
		req, _ := http.NewRequest(http.MethodGet, url, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"total_cost": 1500}`, w.Body.String())
	})

	t.Run("Missing Required Params", func(t *testing.T) {
		h := NewSubscriptionHandler(nil, log.New(io.Discard, "", 0))

		r := gin.Default()
		r.GET("/total-cost", h.TotalCostHandler())

		w := httptest.NewRecorder()

		url := fmt.Sprintf("/total-cost?user_id=%s", uuid.New().String())
		req, _ := http.NewRequest(http.MethodGet, url, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
