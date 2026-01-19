package dao

import (
	"context"
	"fmt"
	"kernel/tools"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mohae/deepcopy"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WriteOperation struct {
	Database   string
	Collection string
	Type       OperationType
	Filter     interface{}
	Update     interface{}
	Document   interface{}
	Uuid       uint64 // 唯一标识
	Tms        int64  // 时间戳
	Cover      bool
}

type OperationType int

const (
	Insert OperationType = iota
	Upsert
	Update
	Delete
	Replace
)

// 批量操作结果
type BatchResult struct {
	SuccessCount int
	ErrorCount   int
	Errors       []error
}

type OpInfo struct {
	DataBase string
	Tms      int64
}

type MongoWriter struct {
	client       *mongo.Client
	queues       []chan *WriteOperation // 多个任务队列
	syncMaps     []sync.Map
	workerWg     sync.WaitGroup // 用于等待所有worker退出
	isRunning    atomic.Bool    // 运行状态标志
	shutdownChan chan struct{}  // 关闭信号通道
	batchSize    int
	flushTimeout time.Duration
	workerCount  int       // 工作协程数量
	shutdownOnce sync.Once // 确保关闭操作只执行一次
	// chan BatchResult // 批量操作结果通道
	//stats        *WriterStats     // 统计信息

	//sync_map   sync.Map
	debug      bool
	collCache  *CollectionCache
	baseCtx    context.Context    // 基础上下文
	cancelFunc context.CancelFunc // 取消函数
}

// 写入器统计信息
type WriterStats struct {
	sync.RWMutex
	TotalProcessed int64
	TotalSucceeded int64
	TotalFailed    int64
	LastError      error
	LastErrorTime  time.Time
}

type CollectionCache struct {
	collections sync.Map // db.collection -> *mongo.Collection
}

func (c *CollectionCache) GetCollection(client *mongo.Client, dbName, collName string) *mongo.Collection {
	key := fmt.Sprintf("%s.%s", dbName, collName)
	if coll, ok := c.collections.Load(key); ok {
		return coll.(*mongo.Collection)
	}

	coll := client.Database(dbName).Collection(collName)
	c.collections.Store(key, coll)
	return coll
}

/*
queueSize:每个任务队列的大小 1000~5000 考虑并发流量和每个任务占用内存
batchSize:批量操作的大小 100~500 批量操作大小有16M限制，大了减少网络请求次数但会提升写入延迟
workerCount:工作协程数量 CPU核心数的1-2倍

如果队列经常满，增加 queueSize 或 workerCount
如果批量操作效率低，调整 batchSize 或 flushTimeout
如果MongoDB负载过高，减少 workerCount 或 batchSize
*/
func NewMongoWriter(client *mongo.Client, workerCount, queueSize, batchSize, flushSeconds int, debug bool, post func(string)) *MongoWriter {
	baseCtx, cancel := context.WithCancel(context.Background())
	// 创建多个队列
	queues := make([]chan *WriteOperation, workerCount)
	maps := make([]sync.Map, workerCount)
	for i := 0; i < workerCount; i++ {
		queues[i] = make(chan *WriteOperation, queueSize)
		maps[i] = sync.Map{}
	}

	writer := &MongoWriter{
		client:       client,
		queues:       queues,
		syncMaps:     maps,
		batchSize:    batchSize,
		flushTimeout: time.Duration(flushSeconds) * time.Second,
		workerCount:  workerCount,
		shutdownChan: make(chan struct{}),
		//batchResults: make(chan BatchResult, 1024), // 缓冲结果通道
		//stats:        &WriterStats{},
		debug:      debug,
		collCache:  &CollectionCache{},
		baseCtx:    baseCtx,
		cancelFunc: cancel,
	}

	writer.isRunning.Store(true)

	// 启动多个worker协程
	writer.workerWg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		//go writer.worker(i)
		go tools.GoSafeWithParam("Monogo Writer", writer.Worker, post, i)
	}

	// 启动统计协程
	//go tools.GoSafePost("Monogo CollectStats", writer.CollectStats, post)

	// 定期监控队列使用情况
	writer.workerWg.Add(1)
	go tools.GoSafePost("Monogo Monitor", func() {
		defer writer.workerWg.Done()
		defer func() {
			if writer.cancelFunc != nil {
				writer.cancelFunc()
			}
		}()
		for range time.Tick(3 * time.Second) {
			for i := 0; i < writer.workerCount; i++ {
				size := writer.WorkerQueueSize(i)
				//if size > queueSize*4/5 {
				log.Warn("mongo woker queue size", zap.Int("woker index", i), zap.Int("current lenth", size), zap.Int("queueSize", queueSize))
				//}
			}
		}
	}, post)

	log.Info("MongoWriter started ",
		zap.Int("workerCount", workerCount),
		zap.Int("batchSize", batchSize),
		zap.Int("queueSize", queueSize),
		zap.Int("flushSeconds", flushSeconds),
	)
	return writer
}

// getQueueIndex 根据玩家ID或Key计算应该使用哪个队列
func (w *MongoWriter) getQueueIndex(op *WriteOperation) int {
	// 根据玩家ID或Key计算队列索引
	return int(op.Uuid % uint64(w.workerCount))
}

// SyncWrite 同步写入操作
func (w *MongoWriter) SyncWrite(ctx context.Context, op *WriteOperation) (*mongo.InsertOneResult, *mongo.UpdateResult, *mongo.DeleteResult, error) {
	if !w.isRunning.Load() {
		return nil, nil, nil, fmt.Errorf("writer is shutting down, cannot accept new requests")
	}

	// 获取对应的集合
	collection := w.client.Database(op.Database).Collection(op.Collection)

	// 设置操作超时（如果调用方没有设置超时）
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	queueIndex := w.getQueueIndex(op)

	// 记录同步更新时间
	if op.Cover && op.Uuid != 0 {
		w.syncMaps[queueIndex].Store(op.Uuid, &OpInfo{
			DataBase: op.Database,
			Tms:      op.Tms,
		})
	}

	if w.debug {
		//log.Debug("mongo sync write operation", zap.Any("op", op))
	}

	opCostFn := func(op *WriteOperation) {
		nowTms := time.Now().UnixMilli()
		if nowTms-op.Tms >= 1000 {
			log.Warn("mongo sync cost over 1 second", zap.Any("op", op))
		}
	}
	// 根据操作类型执行同步写入
	switch op.Type {
	case Insert:
		result, err := collection.InsertOne(ctx, op.Document)
		opCostFn(op)
		return result, nil, nil, err
	case Upsert:
		opts := options.Replace().SetUpsert(true)
		result, err := collection.ReplaceOne(ctx, op.Filter, op.Document, opts)
		opCostFn(op)
		return nil, result, nil, err
	case Update:
		result, err := collection.UpdateOne(ctx, op.Filter, op.Update)
		opCostFn(op)
		return nil, result, nil, err
	case Delete:
		result, err := collection.DeleteOne(ctx, op.Filter)
		opCostFn(op)
		return nil, nil, result, err
	case Replace:
		result, err := collection.ReplaceOne(ctx, op.Filter, op.Document)
		opCostFn(op)
		return nil, result, nil, err
	default:
		return nil, nil, nil, fmt.Errorf("unknown operation type: %d", op.Type)
	}
}

// AsyncWrite 尝试异步写入一个操作
func (w *MongoWriter) AsyncWrite(op *WriteOperation) error {
	if !w.isRunning.Load() {
		return fmt.Errorf("writer is shutting down, cannot accept new requests")
	}

	queueIndex := w.getQueueIndex(op)

	select {
	case w.queues[queueIndex] <- op:
		return nil
	case <-w.shutdownChan:
		return fmt.Errorf("writer is shutting down, cannot accept new requests")
	default:
		return fmt.Errorf("write queue %d is full", queueIndex)
	}
}

// Worker 工作协程
func (w *MongoWriter) Worker(param interface{}) {
	id := param.(int)
	defer w.workerWg.Done()
	log.Info("mongo worker started", zap.Int("id", id))

	queue := w.queues[id]
	var batch []*WriteOperation
	timer := time.NewTimer(w.flushTimeout)
	defer timer.Stop()
	localCache := make(map[uint64]*OpInfo)
	for {
		select {
		case op, ok := <-queue:
			// if w.debug {
			// 	log.Debug("pop sync write operation", zap.Any("op", op))
			// }
			if !ok {
				// 队列通道被关闭，处理最后一批数据后退出
				w.flushBatch(batch)
				log.Info("mongo worker exited, queue closed and flushed", zap.Int("id", id))
				return
			}

			// 有同步更新记录 同集合 且操作时间晚于最后一次全量更新时间才写入
			if info, exists := localCache[op.Uuid]; exists {
				if op.Database == info.DataBase && op.Tms > info.Tms {
					batch = append(batch, op)
				}
			} else if info, ok := w.syncMaps[id].Load(op.Uuid); ok {
				op_info := info.(*OpInfo)
				localCache[op.Uuid] = op_info
				if op.Database == op_info.DataBase && op.Tms > op_info.Tms {
					batch = append(batch, op)
				}
			} else {
				batch = append(batch, op)
			}

			if len(batch) >= w.batchSize {
				w.flushBatch(batch)
				//w.batchResults <- result
				batch = nil
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(w.flushTimeout)
			}

		case <-timer.C:
			if len(batch) > 0 {
				w.flushBatch(batch)
				//w.batchResults <- result
				batch = nil
			}
			timer.Reset(w.flushTimeout)
		case <-w.shutdownChan:
			// 收到关闭信号，处理剩余任务
			log.Info("mongo worker exited, shutdown signal received", zap.Int("id", id))

			// 处理当前批次
			if len(batch) > 0 {
				w.flushBatch(batch)
				//w.batchResults <- result
				batch = nil
			}

			// 处理队列中剩余的任务
		drainLoop:
			for {
				select {
				case op, ok := <-queue:
					if !ok {
						break drainLoop
					}
					w.flushBatch([]*WriteOperation{op})
					//w.batchResults <- result
				default:
					break drainLoop
				}
			}
			log.Info("mongo worker exited, shutdown processed", zap.Int("id", id))
			return
			// default:
			// 	log.Debug("test ", zap.Int("i", id))
		}
	}
}

// flushBatch 处理批量写入并返回结果
func (w *MongoWriter) flushBatch(ops []*WriteOperation) BatchResult {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			log.Warn("Slow batch operation",
				zap.Duration("duration", duration),
				zap.Int("batchSize", len(ops)))
		}
	}()
	if len(ops) == 0 {
		return BatchResult{}
	}

	// 使用更高效的分组方式
	type opGroup struct {
		dbName   string
		collName string
		models   []mongo.WriteModel
	}

	groups := make(map[string]*opGroup) // key: "db.coll"

	// 单次遍历分组
	for _, op := range ops {
		key := fmt.Sprintf("%s.%s", op.Database, op.Collection)

		group, exists := groups[key]
		if !exists {
			group = &opGroup{
				dbName:   op.Database,
				collName: op.Collection,
				models:   make([]mongo.WriteModel, 0, len(ops)),
			}
			groups[key] = group
		}

		// 根据操作类型创建模型
		switch op.Type {
		case Insert:
			group.models = append(group.models, mongo.NewInsertOneModel().SetDocument(op.Document))
		case Upsert, Update:
			model := mongo.NewUpdateOneModel().
				SetFilter(op.Filter).
				SetUpdate(op.Update).
				SetUpsert(op.Type == Upsert)
			group.models = append(group.models, model)
		case Delete:
			group.models = append(group.models, mongo.NewDeleteOneModel().SetFilter(op.Filter))
		}
	}

	// 执行批量操作
	var result BatchResult
	for key, group := range groups {
		if len(group.models) == 0 {
			continue
		}

		// 使用预缓存的集合引用
		collection := w.collCache.GetCollection(w.client, group.dbName, group.collName)

		// 使用带超时的上下文
		ctx, cancel := context.WithTimeout(w.baseCtx, 30*time.Second)

		_, err := collection.BulkWrite(ctx, group.models)
		if err != nil {
			log.Error("BulkWrite failed",
				zap.String("collection", key),
				zap.Error(err),
				zap.Int("batchSize", len(group.models)))
			result.ErrorCount += len(group.models)
		} else {
			result.SuccessCount += len(group.models)
			if w.debug {
				//log.Debug("mongo flush batch", zap.Any("models", group.models))
			}

		}
		cancel()
	}

	return result
}

// CollectStats 收集统计信息
// func (w *MongoWriter) CollectStats() {
// 	for result := range w.batchResults {
// 		w.stats.Lock()
// 		w.stats.TotalProcessed += int64(result.SuccessCount + result.ErrorCount)
// 		w.stats.TotalSucceeded += int64(result.SuccessCount)
// 		w.stats.TotalFailed += int64(result.ErrorCount)
// 		if len(result.Errors) > 0 {
// 			w.stats.LastError = result.Errors[0]
// 			w.stats.LastErrorTime = time.Now()
// 		}
// 		w.stats.Unlock()
// 	}
// }

// Stop 优雅关闭
func (w *MongoWriter) Stop() {
	w.shutdownOnce.Do(func() {
		// 1. 设置运行状态为false
		w.isRunning.Store(false)
		log.Info("Initiating async MongoDB writer shutdown...")

		// 2. 发送关闭信号给所有worker
		close(w.shutdownChan)

		// 3. 关闭所有队列通道
		for i := 0; i < w.workerCount; i++ {
			close(w.queues[i])
		}

		// 4. 等待所有worker协程处理完毕
		w.workerWg.Wait()

		// 5. 关闭批量结果通道
		//close(w.batchResults)

		// 6. 关闭MongoDB连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := w.client.Disconnect(ctx); err != nil {
			log.Error("disconnecting from mongo failed.", zap.Error(err))
		} else {
			log.Info("mongo writer stopped gracefully.")
		}
	})
}

// QueueSize 获取所有队列的总大小
func (w *MongoWriter) QueueSize() int {
	total := 0
	for i := 0; i < w.workerCount; i++ {
		total += len(w.queues[i])
	}
	return total
}

// WorkerQueueSize 获取特定worker队列的大小
func (w *MongoWriter) WorkerQueueSize(workerIndex int) int {
	if workerIndex < 0 || workerIndex >= w.workerCount {
		return 0
	}
	return len(w.queues[workerIndex])
}

// GetStats 获取写入器统计信息
// func (w *MongoWriter) GetStats() (processed, succeeded, failed int64, lastErrorTime time.Time, lastError error) {
// 	w.stats.RLock()
// 	defer w.stats.RUnlock()
// 	return w.stats.TotalProcessed, w.stats.TotalSucceeded,
// 		w.stats.TotalFailed, w.stats.LastErrorTime, w.stats.LastError
// }

// IsRunning 检查写入器是否正在运行
func (w *MongoWriter) IsRunning() bool {
	return w.isRunning.Load()
}

func DeepCopy(data interface{}) interface{} {
	// b, err := bson.Marshal(&data)
	// if err != nil {
	// 	log.Error("", zap.Error(err))
	// 	return nil
	// }
	// var result interface{}
	// err = bson.Unmarshal(b, &result)
	// if err != nil {
	// 	return nil
	// }
	// return result
	res := deepcopy.Copy(data)
	return res
}
