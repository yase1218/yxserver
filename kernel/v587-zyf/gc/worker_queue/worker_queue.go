package worker_queue

import (
	"context"
	"github.com/v587-zyf/gc/internal"
	"kernel/tools"
	"sync"
)

type (
	WorkerQueue struct {
		mu sync.Mutex

		ctx     context.Context
		cancel  context.CancelFunc
		options *WorkerQueueOption

		// double-ended queue to store asynchronous jobs
		q internal.Deque[asyncJob]

		// maximum concurrency
		maxConcurrency int

		// current concurrency
		curConcurrency int
	}

	// Asynchronous job
	asyncJob func()
)

func NewWorkerQueue() *WorkerQueue {
	c := &WorkerQueue{
		maxConcurrency: 256,
		curConcurrency: 0,
		options:        NewWorkerQueueOption(),
	}
	return c
}

func (c *WorkerQueue) Init(ctx context.Context, opts ...any) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(c.options)
		}
	}

	if c.options.maxCount != 0 {
		c.maxConcurrency = c.options.maxCount
	}

	return nil
}

// Retrieves a job from the worker queue
func (c *WorkerQueue) getJob(newJob asyncJob, delta int) asyncJob {
	c.mu.Lock()
	defer c.mu.Unlock()

	if newJob != nil {
		c.q.PushBack(newJob)
	}
	c.curConcurrency += delta
	if c.curConcurrency >= c.maxConcurrency {
		return nil
	}
	var job = c.q.PopFront()
	if job == nil {
		return nil
	}
	c.curConcurrency++
	return job
}

// Do continuously executes jobs in the worker queue
func (c *WorkerQueue) do(job asyncJob) {
	for job != nil {
		job()
		job = c.getJob(nil, -1)
	}
}

// Adds a job to the queue and executes it immediately if resources are available
func (c *WorkerQueue) Push(job asyncJob) {
	if nextJob := c.getJob(job, 0); nextJob != nil {
		go tools.GoSafe("work_queue do", func() {
			c.do(nextJob)
		})
	}
}

func (c *WorkerQueue) GetCtx() context.Context {
	return c.ctx
}
