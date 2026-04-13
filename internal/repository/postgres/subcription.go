package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"subscriptons-service/internal/models"
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
              RETURNING id, created_at`

	err := r.conn.QueryRow(ctx, query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID, &sub.CreatedAt)

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
func (r *SubscriptionRepository) List(ctx context.Context, limit, offset int) ([]*models.Subscription, int, error) {
	var total int

	if err := r.conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM subscriptions`).
		Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at
        FROM subscriptions 
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
	rows, err := r.conn.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("rows: %w", err)
	}
	defer rows.Close()

	result := []*models.Subscription{}
	for rows.Next() {
		row := &models.Subscription{}
		if err := rows.Scan(&row.ID, &row.ServiceName, &row.Price, &row.UserID, &row.StartDate, &row.EndDate, &row.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan: %w", err)
		}
		result = append(result, row)
	}
	return result, total, nil
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

func (r *SubscriptionRepository) FindSubscriptions(ctx context.Context, userID uuid.UUID, serviceName *string, startDate, endDate string) ([]*models.Subscription, error) {
	if endDate == "" {
		endDate = startDate
	}
	query := `
		SELECT id, price, service_name
		FROM subscriptions
		WHERE user_id = $1
		  AND ($2::text IS NULL OR service_name ILIKE '%' || $2::text || '%')
		  AND start_date <= $4
		  AND COALESCE(end_date, '99-9999') >= $3
	`
	rows, err := r.conn.Query(ctx, query, userID, serviceName, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("find subs: %w", err)
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.Price, &sub.ServiceName); err != nil {
			return nil, err
		}
		subs = append(subs, &sub)
	}
	return subs, nil
}
