package buffer_pool

import (
	"context"
	"kernel/tools"
	"sync"
	"sync/atomic"
	"time"
)

// 缓存区结构体
type Buffer struct {
	Data []byte
}

// 缓存池结构体
type BufferPool struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	options *BufferPoolOptions
	Pool    chan *Buffer

	MaxSize      int64
	currentSize  int64
	releaseCount int64

	stats  map[string]int
	ticker *time.Ticker
}

func NewBufferPool() *BufferPool {
	bp := &BufferPool{
		Pool:    make(chan *Buffer, 512),
		MaxSize: 512,
		stats:   make(map[string]int),
		options: DefaultBufferPoolOptions(),
	}
	return bp
}

func (bp *BufferPool) Init(ctx context.Context, post func(string), opts ...any) error {
	bp.ctx, bp.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(BufferPoolOption)(bp.options)
		}
		if bp.options.maxSize != 0 {
			bp.MaxSize = bp.options.maxSize
		}
	}

	for i := 0; i < int(bp.options.size); i++ {
		bp.Pool <- &Buffer{Data: make([]byte, bp.options.bufferSize)}
	}

	if bp.options.autoCleanup {
		bp.ticker = time.NewTicker(bp.options.cleanupPeriod)
		if post == nil {
			go tools.GoSafe("buffer pool clean loop", func() {
				bp.cleanupLoop()
			})
		} else {
			go tools.GoSafePost("buffer pool clean loop", func() {
				bp.cleanupLoop()
			}, post)
		}
	}

	return nil
}

func (bp *BufferPool) GetCtx() context.Context {
	return bp.ctx
}

// 获取缓冲区
func (bp *BufferPool) Get() *Buffer {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.stats["get"]++

	select {
	case <-bp.ctx.Done():
		bp.mu.Unlock()
		return nil
	case buf := <-bp.Pool:
		bp.stats["hit"]++
		return buf
	default:
		bp.stats["miss"]++
		//fmt.Println("缓存池已空，创建新缓冲区")
		atomic.AddInt64(&bp.currentSize, 1)
		if bp.currentSize > bp.MaxSize {
			atomic.AddInt64(&bp.currentSize, -1)
			return &Buffer{Data: make([]byte, bp.options.bufferSize)}
		}
		return &Buffer{Data: make([]byte, bp.options.bufferSize)}
	}
}

// 释放缓冲区
func (bp *BufferPool) Put(buf *Buffer) {
	buf.Data = buf.Data[:0]

	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.stats["put"]++
	bp.releaseCount++

	select {
	case <-bp.ctx.Done():
		bp.mu.Unlock()
		return
	case bp.Pool <- buf:
		//fmt.Printf("释放缓冲区，当前池大小：%d\n", len(bp.Pool))
	default:
		//fmt.Println("缓存池已满，丢弃缓冲区")
		atomic.AddInt64(&bp.currentSize, -1)
	}

	if bp.releaseCount%10 == 0 && len(bp.Pool) > int(bp.MaxSize)/2 {
		close(bp.Pool)
		bp.Pool = make(chan *Buffer, bp.MaxSize)
		for i := 0; i < int(bp.MaxSize); i++ {
			bp.Pool <- &Buffer{Data: make([]byte, bp.options.bufferSize)}
		}
		bp.currentSize = bp.MaxSize
	}
}

// 获取统计信息
func (bp *BufferPool) GetStats() map[string]int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return bp.stats
}

// 清理循环
func (bp *BufferPool) cleanupLoop() {
	for {
		select {
		case <-bp.ticker.C:
			//fmt.Println("执行自动清理")
			bp.cleanup()
		case <-bp.ctx.Done():
			return
		}
	}
}

// 执行清理
func (bp *BufferPool) cleanup() {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	if len(bp.Pool) > int(bp.MaxSize)/2 {
		//fmt.Println("清理部分缓冲区")
		close(bp.Pool)
		bp.Pool = make(chan *Buffer, bp.MaxSize)
		for i := 0; i < int(bp.MaxSize); i++ {
			bp.Pool <- &Buffer{Data: make([]byte, bp.options.bufferSize)}
		}
		bp.currentSize = bp.MaxSize
	}
}
