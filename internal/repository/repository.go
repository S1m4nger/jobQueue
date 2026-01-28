package repository

import (
	"context"
	"job-queue/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job domain.Job) error
	GetByID(ctx context.Context, id string) (domain.Job, error)
	UpdateStatus(ctx context.Context, id, status string) error
	SaveResult(ctx context.Context, id string, result []byte) error
}
