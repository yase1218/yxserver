package buffer_pool

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 基本功能测试
func TestBasicFunctionality(t *testing.T) {
	err := Init(
		context.Background(),
		WithSize(100),
		WithBufferSize(1024),
		WithMaxSize(150),
		WithAutoCleanup(true),
		WithCleanupPeriod(1*time.Second),
	)
	if err != nil {
		t.Errorf("Failed to initialize buffer pool: %v", err)
		return
	}

	// 获取缓冲区
	buf := GetBuffer()
	if buf == nil {
		t.Errorf("Expected non-nil buffer, got nil")
		return
	}

	// 使用缓冲区
	buf.Data = append(buf.Data, "Hello, world!"...)

	// 释放缓冲区
	Put(buf)

	// 验证统计信息
	stats := GetStats()
	expectedStats := map[string]int{
		"get":  1,
		"hit":  1,
		"miss": 0,
		"put":  1,
	}

	for k, v := range expectedStats {
		if stats[k] != v {
			t.Errorf("Expected %s count to be %d, got %d", k, v, stats[k])
		}
	}
}

// 并发测试
func TestConcurrency(t *testing.T) {
	err := Init(
		context.Background(),
		WithSize(100),
		WithBufferSize(1024),
		WithMaxSize(150),
		WithAutoCleanup(true),
		WithCleanupPeriod(1*time.Second),
	)
	if err != nil {
		t.Errorf("Failed to initialize buffer pool: %v", err)
		return
	}

	numWorkers := 10
	numMessages := 1000
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	startTime := time.Now()

	bufferPool := GetBufferPool()
	if bufferPool == nil {
		t.Errorf("Expected non-nil buffer, got nil")
		return
	}
	for i := 0; i < numWorkers; i++ {
		go func() {
			testConcurrent(bufferPool, numMessages)
			wg.Done()
		}()
	}

	// 等待所有并发操作完成
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("所有缓冲区已处理，耗时: %s\n", duration)

	// 验证统计信息
	stats := GetStats()
	expectedStats := map[string]int{
		"get":  numWorkers * numMessages,
		"hit":  1000, // 假设一半命中
		"miss": 0,    // 假设一半未命中
		"put":  numWorkers * numMessages,
	}

	for k, v := range expectedStats {
		if stats[k] < v {
			t.Errorf("Expected %s count to be at least %d, got %d", k, v, stats[k])
		}
	}
}

// 测试并发获取和释放
func testConcurrent(bufferPool *BufferPool, numMessages int) {
	for i := 0; i < numMessages; i++ {
		buf := bufferPool.Get()
		buf.Data = append(buf.Data, fmt.Sprintf("Message %d", i)...)
		bufferPool.Put(buf)
	}
}

// 统计信息测试
func TestStatistics(t *testing.T) {
	err := Init(
		context.Background(),
		WithSize(100),
		WithBufferSize(1024),
		WithMaxSize(150),
		WithAutoCleanup(true),
		WithCleanupPeriod(1*time.Second),
	)
	if err != nil {
		t.Errorf("Failed to initialize buffer pool: %v", err)
		return
	}

	bufferPool := GetBufferPool()
	if bufferPool == nil {
		t.Errorf("Expected non-nil buffer, got nil")
		return
	}
	// 获取和释放缓冲区
	for i := 0; i < 1000; i++ {
		buf := bufferPool.Get()
		buf.Data = append(buf.Data, fmt.Sprintf("Message %d", i)...)
		bufferPool.Put(buf)
	}

	// 验证统计信息
	stats := bufferPool.GetStats()
	expectedStats := map[string]int{
		"get":  1000,
		"hit":  1000, // 假设一半命中
		"miss": 0,    // 假设一半未命中
		"put":  1000,
	}

	for k, v := range expectedStats {
		if stats[k] < v {
			t.Errorf("Expected %s count to be at least %d, got %d", k, v, stats[k])
		}
	}
}

// 清理机制测试
func TestCleanup(t *testing.T) {
	err := Init(
		context.Background(),
		WithSize(100),
		WithBufferSize(1024),
		WithMaxSize(150),
		WithAutoCleanup(true),
		WithCleanupPeriod(1*time.Second),
	)
	if err != nil {
		t.Errorf("Failed to initialize buffer pool: %v", err)
		return
	}

	bufferPool := GetBufferPool()
	if bufferPool == nil {
		t.Errorf("Expected non-nil buffer, got nil")
		return
	}

	// 获取和释放缓冲区
	for i := 0; i < 1000; i++ {
		buf := bufferPool.Get()
		buf.Data = append(buf.Data, fmt.Sprintf("Message %d", i)...)
		bufferPool.Put(buf)
	}

	time.Sleep(2 * time.Second) // 等待自动清理

	// 验证清理后的池大小
	if len(bufferPool.Pool) > int(bufferPool.MaxSize)/2 {
		t.Errorf("Expected pool size to be less than or equal to %d, got %d", bufferPool.MaxSize/2, len(bufferPool.Pool))
	}
}
