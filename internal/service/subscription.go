package service

import (
	"context"
	"fmt"
	"log"
	"subscriptons-service/internal/models"
	"subscriptons-service/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo repository.SubscriptionRepository
	log  *log.Logger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, log *log.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		log:  log,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, sub *models.Subscription) error {
	err := s.repo.Create(ctx, sub)
	if err != nil {
		s.log.Printf("create repo error: %v", err)
		return err
	}
	s.log.Printf("created subscription")
	return nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Printf("error get subscription id=%s: %v", id, err)
		return nil, err
	}
	s.log.Printf("found subscription id=%s service=%s", id, sub.ServiceName)
	return sub, nil
}

func (s *SubscriptionService) List(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error) {
	if limit > 100 || limit <= 0 {
		return nil, 0, fmt.Errorf("limit must be 1-100")
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *models.Subscription) error {
	if err := s.repo.Update(ctx, sub); err != nil {
		s.log.Printf("error update subscription id=%s: %v", sub.ID, err)
		return err
	}
	s.log.Printf("updated subscription id=%s", sub.ID)
	return nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Printf("error delete subscription id=%s: %v", id, err)
		return err
	}
	s.log.Printf("deleted subscription id=%s", id)
	return nil
}

func (s *SubscriptionService) TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, startDate, endDate string) (int, error) {
	subs, err := s.repo.FindSubscriptions(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		s.log.Printf("error finding subs: %v", err)
		return 0, err
	}
	var total int

	for _, sub := range subs {
		total += sub.Price
	}
	s.log.Printf("total cost user: %v", err)
	return total, nil
}
