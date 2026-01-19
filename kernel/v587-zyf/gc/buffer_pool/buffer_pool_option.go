package buffer_pool

import "time"

// 缓存池选项
type BufferPoolOption func(*BufferPoolOptions)

// 缓存池选项结构体
type BufferPoolOptions struct {
	size          int64
	bufferSize    int64
	maxSize       int64
	autoCleanup   bool
	cleanupPeriod time.Duration
}

func DefaultBufferPoolOptions() *BufferPoolOptions {
	return &BufferPoolOptions{
		size:          100,
		bufferSize:    1024,
		maxSize:       100,
		autoCleanup:   true,
		cleanupPeriod: 1 * time.Second,
	}
}

func WithSize(size int64) BufferPoolOption {
	return func(options *BufferPoolOptions) {
		options.size = size
	}
}

func WithBufferSize(bufferSize int64) BufferPoolOption {
	return func(options *BufferPoolOptions) {
		options.bufferSize = bufferSize
	}
}

func WithMaxSize(maxSize int64) BufferPoolOption {
	return func(options *BufferPoolOptions) {
		options.maxSize = maxSize
	}
}

func WithAutoCleanup(autoCleanup bool) BufferPoolOption {
	return func(options *BufferPoolOptions) {
		options.autoCleanup = autoCleanup
	}
}

func WithCleanupPeriod(cleanupPeriod time.Duration) BufferPoolOption {
	return func(options *BufferPoolOptions) {
		options.cleanupPeriod = cleanupPeriod
	}
}
