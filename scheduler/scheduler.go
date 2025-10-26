// Package scheduler implements a dynamic, concurrent job scheduler in Go.
package scheduler

import (
	"context"
	"sync"
	"time"
)

// Job defines an interface that every scheduled task must implement.
type Job interface {
	Run(ctx context.Context) // Run is called at job interval with scheduler context
	Name() string           // Unique name for identification
}

type scheduledJob struct {
	job       Job
	interval  time.Duration
	stopCh    chan struct{}
	wg        *sync.WaitGroup
}

type Scheduler struct {
	mu       sync.Mutex
	jobs     map[string]*scheduledJob
	wg       sync.WaitGroup // Tracks all in-progress job-executions
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduler creates a new Scheduler instance.
func NewScheduler(ctx context.Context) *Scheduler {
	cctx, cancel := context.WithCancel(ctx)
	return &Scheduler{
		jobs:   make(map[string]*scheduledJob),
		ctx:    cctx,
		cancel: cancel,
	}
}

// AddJob registers and starts a new job at specified interval. If job with same name exists, removes it first.
func (s *Scheduler) AddJob(job Job, interval time.Duration) {
	s.mu.Lock()
	if existing, ok := s.jobs[job.Name()]; ok {
		s.mu.Unlock()
		s.RemoveJob(job.Name())
		s.mu.Lock()
	}
	j := &scheduledJob{
		job:      job,
		interval: interval,
		stopCh:   make(chan struct{}),
		wg:       &s.wg,
	}
	s.jobs[job.Name()] = j
	s.mu.Unlock()

	s.wg.Add(1)
	go s.runJob(j)
}

// RemoveJob stops and removes job with given name.
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	job, ok := s.jobs[name]
	if ok {
		delete(s.jobs, name)
	}
	s.mu.Unlock()
	if ok {
		close(job.stopCh) // Trigger the goroutine to exit and cleanup
	}
}

func (s *Scheduler) runJob(j *scheduledJob) {
	defer s.wg.Done()
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.wg.Add(1)
			go func() { // actual job may block, so run in separate goroutine
				defer j.wg.Done()
				// Run with root ctx, but allow job to choose further scoping inside
				j.job.Run(s.ctx)
			}()
		case <-j.stopCh:
			return
		case <-s.ctx.Done():
			return
		}
	}
}

// Shutdown stops all jobs, waits for all running jobs to finish, and cleans up.
// Returns when all in-flight jobs are done.
func (s *Scheduler) Shutdown() {
	s.cancel()

	s.mu.Lock()
	for _, job := range s.jobs {
		close(job.stopCh)
	}
	s.jobs = make(map[string]*scheduledJob)
	s.mu.Unlock()

	s.wg.Wait()
}
