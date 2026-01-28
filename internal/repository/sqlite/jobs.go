package sqlite

import (
	"context"
	"database/sql"
	"job-queue/internal/domain"
	"time"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job domain.Job) error {
	query := `INSERT INTO jobs (id, type, payload, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, job.ID, job.Type, job.Payload, job.Status, job.CreatedAt, job.UpdatedAt)
	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (domain.Job, error) {
	query := `SELECT id, type, payload, status, result, created_at, updated_at FROM jobs WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var job domain.Job
	err := row.Scan(&job.ID, &job.Type, &job.Payload, &job.Status, &job.Result, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *JobRepository) SaveResult(ctx context.Context, id string, result []byte) error {
	query := `UPDATE jobs SET result = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, result, domain.StatusDone, time.Now(), id)
	return err
}
