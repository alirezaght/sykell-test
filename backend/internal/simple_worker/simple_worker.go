package simpleworker

import (
	"context"
	"errors"
	"sykell-backend/internal/crawl"
	"sync"
	"time"
)

// SimpleWorker is a basic implementation of a worker that processes jobs
type SimpleWorker struct {
	jobChan chan crawl.WorlFlowInput
	mux sync.Mutex
	runningJobs map[string]context.CancelFunc
	timeout time.Duration
}

// NewSimpleWorker creates a new instance of SimpleWorker
func NewSimpleWorker(buffer int, timeout time.Duration) *SimpleWorker {
	return &SimpleWorker{
		jobChan: make(chan crawl.WorlFlowInput, buffer),
		runningJobs: make(map[string]context.CancelFunc),
		timeout: timeout,
	}
}

// StartWorkerLoop starts the worker loop to process incoming jobs
func (sw *SimpleWorker) StartWorkerLoop(ctx context.Context) {
	
	go func() {
		for {
			select {
				case <-ctx.Done():
					return
				case job := <-sw.jobChan:					
					sw.mux.Lock()
					if _, exists := sw.runningJobs[job.WorkflowID]; exists { 
						sw.mux.Unlock()
						continue
					}
					ctx, cancel := context.WithTimeout(ctx, sw.timeout)
					sw.runningJobs[job.WorkflowID] = cancel
					sw.mux.Unlock()
					go func(job crawl.WorlFlowInput) {
						defer func() {
							sw.mux.Lock()
							delete(sw.runningJobs, job.WorkflowID)
							sw.mux.Unlock()
						}()
						crawl.CrawlURLActivity(ctx, job, func(ctx context.Context, details string) {})
						}(job)							
			}
		}
	}()
}


// Start begins processing jobs from the job channel
func (sw *SimpleWorker) Start(ctx context.Context, input crawl.WorlFlowInput) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case sw.jobChan <- input:
		return nil
	}
}


// Stop stops processing jobs for the given workflow ID
func (sw *SimpleWorker) Stop(ctx context.Context, workflowID string) error {
	sw.mux.Lock()
	cancel, exists := sw.runningJobs[workflowID]
	sw.mux.Unlock()
	if !exists {
		return errors.New("no running job found for the given workflow ID")
	}
	cancel()
	return nil
}

// Shutdown gracefully shuts down the worker, waiting for running jobs to complete
func (sw *SimpleWorker) Shutdown() {
	sw.mux.Lock()
	for _, c := range sw.runningJobs {
		c()
	}
	sw.mux.Unlock()
	close(sw.jobChan)
}