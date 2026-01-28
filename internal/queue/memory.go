package queue

import "job-queue/internal/domain"

type MemoryQueue struct {
	jobs chan domain.Job
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		jobs: make(chan domain.Job, bufferSize),
	}
}

func (q *MemoryQueue) Enqueue(job domain.Job) {
	q.jobs <- job
}

func (q *MemoryQueue) Dequeue() <-chan domain.Job {
	return q.jobs
}
