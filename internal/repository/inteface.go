package repository

import (
	"context"
	"subscriptons-service/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	List(ctx context.Context) ([]*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, startDate, endDate string) (int, error)
}
