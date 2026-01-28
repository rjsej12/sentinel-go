package worker

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type JobType string

const (
	JobTypeProcess JobType = "process"
	JobTypeDelay   JobType = "delay"
)

type Job struct {
	ID        string
	Type      JobType
	Data      []byte
	CreatedAt time.Time
}

type Queue struct {
	jobs chan Job
	ctx  context.Context
}

func NewQueue(ctx context.Context, size int) *Queue {
	return &Queue{
		jobs: make(chan Job, size),
		ctx:  ctx,
	}
}

func (q *Queue) Enqueue(job Job) error {
	select {
	case q.jobs <- job:
		return nil
	case <-q.ctx.Done():
		return q.ctx.Err()
	}
}

func (q *Queue) Dequeue() (Job, error) {
	select {
	case job := <-q.jobs:
		return job, nil
	case <-q.ctx.Done():
		return Job{}, q.ctx.Err()
	}
}

func (q *Queue) Length() int {
	return len(q.jobs)
}

func (q *Queue) Capacity() int {
	return cap(q.jobs)
}

func (q *Queue) Close() {
	close(q.jobs)
}

func NewJob(jobType JobType, data []byte) Job {
	return Job{
		ID:        generateID(),
		Type:      jobType,
		Data:      data,
		CreatedAt: time.Now(),
	}
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
