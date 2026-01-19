package worker_queue

type WorkerQueueOption struct {
	maxCount int

	errHandler func(args ...any)
}

type Option func(o *WorkerQueueOption)

func NewWorkerQueueOption() *WorkerQueueOption {
	return &WorkerQueueOption{}
}

func WithMaxCount(maxCount int) Option {
	return func(o *WorkerQueueOption) {
		o.maxCount = maxCount
	}
}
func WithErrHandler(errHandler func(args ...any)) Option {
	return func(o *WorkerQueueOption) {
		o.errHandler = errHandler
	}
}
