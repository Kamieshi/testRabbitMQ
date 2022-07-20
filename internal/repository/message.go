package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"testRabitMQ/internal/models"
)

type MessageRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*models.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Write(ctx context.Context, message *models.Message) error
	BatchQuery(ctx context.Context, bt *pgx.Batch) error
}
