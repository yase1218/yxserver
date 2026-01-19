package worker_pool

type WorkerPoolOption struct {
	maxCount int

	errHandler func(args ...any)
}

type Option func(o *WorkerPoolOption)

func NewWorkerPoolOption() *WorkerPoolOption {
	return &WorkerPoolOption{}
}

func WithMaxCount(maxCount int) Option {
	return func(o *WorkerPoolOption) {
		o.maxCount = maxCount
	}
}
func WithErrHandler(errHandler func(args ...any)) Option {
	return func(o *WorkerPoolOption) {
		o.errHandler = errHandler
	}
}
