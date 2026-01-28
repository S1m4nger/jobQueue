package service

import (
	"context"
	"job-queue/internal/domain"
	"job-queue/internal/queue"
	"job-queue/internal/repository"
	"time"

	"github.com/google/uuid"
)

type JobService struct {
	repo  repository.JobRepository
	queue queue.Queue
}

func NewJobService(repo repository.JobRepository, queue queue.Queue) *JobService {
	return &JobService{
		repo:  repo,
		queue: queue,
	}
}

func (s *JobService) CreateJob(ctx context.Context, jobType string, payload []byte) (domain.Job, error) {
	job := domain.Job{
		ID:        uuid.New().String(),
		Type:      jobType,
		Payload:   payload,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return domain.Job{}, err
	}

	s.queue.Enqueue(job)

	return job, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (domain.Job, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *JobService) StartWorkers(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			jobs := s.queue.Dequeue()
			for {
				select {
				case job := <-jobs:
					s.executeJob(ctx, job)
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

func (s *JobService) executeJob(ctx context.Context, job domain.Job) {
	// 1. Update status to running
	_ = s.repo.UpdateStatus(ctx, job.ID, domain.StatusRunning)

	// 2. Perform mock processing
	result, err := s.process(job)

	// 3. Save result and update status
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, job.ID, domain.StatusFailed)
		return
	}

	_ = s.repo.SaveResult(ctx, job.ID, result)
}

func (s *JobService) process(job domain.Job) ([]byte, error) {
	// Mock processing time
	time.Sleep(2 * time.Second)

	// Simple mock result
	return []byte("Processed: " + job.Type), nil
}
