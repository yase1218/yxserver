package db

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"time"

	"kernel/dao"

	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	global_mongo_conn *mongo.Client
	GlobalMongoWriter *dao.MongoWriter
	GlobalMongoReader *dao.MongoReader
)

func InitGlobal(post func(string)) {
	init_global_mongo(post)
}

func init_global_mongo(post func(string)) {
	// 配置连接池参数
	clientOptions := options.Client().
		ApplyURI(config.Conf.GlobalMongo.Url).
		SetMaxConnecting(100).                                         // 允许快速建立连接
		SetMaxPoolSize(uint64(config.Conf.GlobalMongo.WorkCount * 2)). // 连接池大小为worker数量的2倍
		SetMinPoolSize(uint64(config.Conf.GlobalMongo.WorkCount)).     // 最小连接数等于worker数量
		// SetMaxPoolSize(200).                 // 连接池大小为worker数量的2倍
		// SetMinPoolSize(50).                  // 最小连接数等于worker数量
		SetMaxConnIdleTime(5 * time.Minute). // 连接最大空闲时间
		SetConnectTimeout(10 * time.Second). // 连接超时时间
		SetSocketTimeout(30 * time.Second)   // Socket操作超时时间

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MongoDB url: %s err: %v", config.Conf.GlobalMongo.Url, err))
	}

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to ping MongoDB url: %s err: %v", config.Conf.GlobalMongo.Url, err))
	}
	global_mongo_conn = client
	GlobalMongoWriter = dao.NewMongoWriter(
		global_mongo_conn,
		config.Conf.GlobalMongo.WorkCount,
		config.Conf.GlobalMongo.QueueSize,
		config.Conf.GlobalMongo.BatchSize,
		config.Conf.GlobalMongo.FlushSeconds,
		config.Conf.Debug,
		post,
	)

	GlobalMongoReader = dao.NewMongoReader(global_mongo_conn)
}

func StopGlobal() {
	log.Info("Global DB close start")
	GlobalMongoWriter.Stop()
	global_mongo_conn.Disconnect(context.Background())
	log.Info("Global DB close succ")
}
