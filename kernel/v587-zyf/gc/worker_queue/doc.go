package worker_queue

import (
	"context"
)

var defaultWorkQueue *WorkerQueue

func Init(ctx context.Context, opts ...any) (err error) {
	defaultWorkQueue = NewWorkerQueue()
	if err = defaultWorkQueue.Init(ctx, opts...); err != nil {
		return err
	}

	return
}
func GetCtx() context.Context {
	return defaultWorkQueue.GetCtx()
}

func Push(job asyncJob) {
	defaultWorkQueue.Push(job)
}
