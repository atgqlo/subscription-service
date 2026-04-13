package service

import (
	"context"
	"io"
	"log"
	"subscriptons-service/internal/mocks"
	"subscriptons-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubService_TotalCost(t *testing.T) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	testUserID := uuid.New()
	mockSubs := []*models.Subscription{
		{ID: uuid.New(), Price: 200, ServiceName: "Netflix"},
		{ID: uuid.New(), Price: 300, ServiceName: "Spotify"},
	}
	mockRepo.On("FindSubscriptions", mock.Anything, testUserID, mock.Anything, mock.Anything, mock.Anything).Return(mockSubs, nil)

	silentLogger := log.New(io.Discard, "", 0)
	service := NewSubscriptionService(mockRepo, silentLogger)

	total, err := service.TotalCost(context.Background(), testUserID, nil, "03-2026", "04-2026")
	assert.NoError(t, err)

	assert.Equal(t, 500, total, "")
}
func TestSubscriptionService_List_Success(t *testing.T) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	mockSubs := []*models.Subscription{
		{ID: uuid.New(), ServiceName: "Netflix"},
		{ID: uuid.New(), ServiceName: "Spotify"},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(mockSubs, 2, nil)

	silentLogger := log.New(io.Discard, "", 0)
	service := NewSubscriptionService(mockRepo, silentLogger)

	subs, total, err := service.List(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, subs, 2)
}

func TestSubscriptionService_List_LimitError(t *testing.T) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	silentLogger := log.New(io.Discard, "", 0)
	service := NewSubscriptionService(mockRepo, silentLogger)

	subs, total, err := service.List(context.Background(), 150, 0)

	assert.Error(t, err)
	assert.EqualError(t, err, "limit must be 1-100")
	assert.Nil(t, subs)
	assert.Equal(t, 0, total)
}

func TestSubscriptionService_Update_Success(t *testing.T) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	subToUpdate := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix Premium",
	}

	mockRepo.On("Update", mock.Anything, subToUpdate).Return(nil)

	silentLogger := log.New(io.Discard, "", 0)
	service := NewSubscriptionService(mockRepo, silentLogger)

	err := service.Update(context.Background(), subToUpdate)

	assert.NoError(t, err)
}

func TestSubscriptionService_Delete_Success(t *testing.T) {
	mockRepo := mocks.NewSubscriptionRepository(t)
	testID := uuid.New()

	mockRepo.On("Delete", mock.Anything, testID).Return(nil)

	silentLogger := log.New(io.Discard, "", 0)
	service := NewSubscriptionService(mockRepo, silentLogger)

	err := service.Delete(context.Background(), testID)

	assert.NoError(t, err)
}

func TestSubscriptionService_GetByID_Success(t *testing.T) {

	mockRepo := mocks.NewSubscriptionRepository(t)

	testID := uuid.New()

	expectedSub := &models.Subscription{
		ID:          testID,
		ServiceName: "Spotify",
	}

	mockRepo.On("GetByID", mock.Anything, testID).Return(expectedSub, nil)

	silentLogger := log.New(io.Discard, "", 0)

	service := NewSubscriptionService(mockRepo, silentLogger)

	result, err := service.GetByID(context.Background(), testID)

	assert.NoError(t, err, "ошибка должна быть nil при успешном поиске")

	assert.NotNil(t, result, "результат не должен быть nil")

	assert.Equal(t, expectedSub.ServiceName, result.ServiceName, "имена сервисов должны совпадать")
	assert.Equal(t, expectedSub.ID, result.ID, "ID должны совпадать")

	mockRepo.AssertExpectations(t)
}
