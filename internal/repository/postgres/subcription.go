package postgres

import (
	"context"
	"fmt"
	"subscriptons-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	conn *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{conn: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
              VALUES($1, $2, $3, $4, $5) 
              RETURNING id`

	var id uuid.UUID
	err := r.conn.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)

	if err != nil {
		return fmt.Errorf("Create: %w", err)
	}

	return nil
}
func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at
	FROM subscriptions
	WHERE id = $1`
	result := &models.Subscription{}
	err := r.conn.QueryRow(ctx, query, id).Scan(&result.ID, &result.ServiceName, &result.Price, &result.UserID, &result.StartDate, &result.EndDate, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get:%w", err)
	}
	return result, nil
}
func (r *SubscriptionRepository) List(ctx context.Context) ([]*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at
	FROM subscriptions`
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("rows:%w", err)
	}
	defer rows.Close()
	result := []*models.Subscription{}
	for rows.Next() {
		row := &models.Subscription{}
		err := rows.Scan(&row.ID, &row.ServiceName, &row.Price, &row.UserID, &row.StartDate, &row.EndDate, &row.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan:%w", err)
		}
		result = append(result, row)
	}
	return result, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions 
	SET service_name=COALESCE($1, service_name),
	price=COALESCE($2, price),
	start_date=COALESCE($3, start_date),
	end_date=COALESCE($4, end_date)
	WHERE id=$5
	RETURNING id, service_name, price, user_id, start_date, end_date, created_at`

	err := r.conn.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate, &sub.CreatedAt)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions
	WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete:%w", err)
	}
	return nil
}

func (r *SubscriptionRepository) TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, startDate, endDate string) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0)
	FROM subscriptions
	WHERE user_id = $1
	AND (service_name = $2 OR $2 IS NULL)
	AND start_date =$3
	`
	var total int
	err := r.conn.QueryRow(ctx, query, userID, serviceName, startDate).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("total:%w", err)
	}
	return total, nil

}
