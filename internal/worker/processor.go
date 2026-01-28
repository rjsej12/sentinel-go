package worker

import (
	"context"
	"log"
	"time"

	"github.com/rjsej12/sentinel-go/internal/metrics"
)

type Processor struct {
	queue      *Queue
	workerPool int
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewProcessor(ctx context.Context, queue *Queue, workerPool int) *Processor {
	procCtx, cancel := context.WithCancel(ctx)
	return &Processor{
		queue:      queue,
		workerPool: workerPool,
		ctx:        procCtx,
		cancel:     cancel,
	}
}

func (p *Processor) Start() {
	for i := 0; i < p.workerPool; i++ {
		go p.worker(i)
	}
	log.Printf("Processor started with %d workers", p.workerPool)

	metrics.WorkerQueueCapacity.Set(float64(p.queue.Capacity()))
	go p.updateQueueMetrics()
}

func (p *Processor) updateQueueMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			metrics.WorkerQueueLength.Set(float64(p.queue.Length()))
		}
	}
}

func (p *Processor) worker(id int) {
	log.Printf("Worker %d started", id)
	defer log.Printf("Worker %d stopped", id)

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			job, err := p.queue.Dequeue()
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("Worker %d: error dequeing: %v", id, err)
				metrics.WorkerErrorsTotal.WithLabelValues("dequeue").Inc()
				continue
			}

			p.processJob(id, job)
		}
	}
}

func (p *Processor) processJob(workerID int, job Job) {
	start := time.Now()

	metrics.WorkerJobsInFlight.Inc()
	defer metrics.WorkerJobsInFlight.Dec()

	log.Printf("Worker: %d: processing job %s (type: %s)", workerID, job.ID, job.Type)

	var err error
	jobTypeStr := string(job.Type)

	switch job.Type {
	case JobTypeProcess:
		err = p.handleProcessJob(job)
	case JobTypeDelay:
		err = p.handleDelayJob(job)
	default:
		log.Printf("Worker: %d: unknown job type: %s", workerID, job.Type)
		metrics.WorkerErrorsTotal.WithLabelValues("unknown_type").Inc()
		metrics.WorkerJobsProcessedTotal.WithLabelValues(jobTypeStr, "error").Inc()
		return
	}

	duration := time.Since(start)

	metrics.WorkerJobDurationSeconds.WithLabelValues(jobTypeStr).Observe(duration.Seconds())

	if err != nil {
		metrics.WorkerJobsProcessedTotal.WithLabelValues(jobTypeStr, "error").Inc()
		metrics.WorkerErrorsTotal.WithLabelValues("process").Inc()
		log.Printf("Worker: %d: error processing job %s: %v", workerID, job.ID, err)
	} else {
		metrics.WorkerJobsProcessedTotal.WithLabelValues(jobTypeStr, "success").Inc()
		log.Printf("Worker: %d: completed job %s in %v", workerID, job.ID, duration)
	}
}

func (p *Processor) handleProcessJob(job Job) error {
	processTime := 100 + (len(job.Data) % 400)
	time.Sleep(time.Duration(processTime) * time.Millisecond)
	return nil
}

func (p *Processor) handleDelayJob(job Job) error {
	delayTime := 1000 + (len(job.Data) % 2000)
	time.Sleep(time.Duration(delayTime) * time.Millisecond)
	return nil
}

func (p *Processor) Stop() {
	log.Println("Stopping processor...")
	p.cancel()

	done := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		log.Println("Processor stopped")
	case <-time.After(5 * time.Second):
		log.Println("Processor stop timeout")
	}
}

func (p *Processor) Stats() ProcessorStats {
	return ProcessorStats{
		QueueLength:   p.queue.Length(),
		QueueCapacity: p.queue.Capacity(),
		WorkerPool:    p.workerPool,
	}
}

type ProcessorStats struct {
	QueueLength   int
	QueueCapacity int
	WorkerPool    int
}
