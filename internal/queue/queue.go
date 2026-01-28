package queue

import "job-queue/internal/domain"

type Queue interface {
	Enqueue(job domain.Job)
	Dequeue() <-chan domain.Job
}
