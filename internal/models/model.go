package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ServiceName string    `db:"service_name" json:"service_name" binding:"required"`
	Price       int       `db:"price" json:"price" binding:"required"`
	UserID      uuid.UUID `db:"user_id" json:"user_id" binding:"required"`
	StartDate   string    `db:"start_date" json:"start_date" binding:"required"`
	EndDate     *string   `db:"end_date" json:"end_date"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
