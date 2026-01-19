package db

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"time"

	"kernel/dao"

	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	CollectionName_UserAccount = "useraccount"
	CollectionName_User        = "user"
	CollectionName_GlobalMail  = "global_mail"
	CollectionName_Order       = "order_info"
	CollectionName_FriendOp    = "friendop"

	CollectionName_Global_Account    = "account"
	CollectionName_Global_AccountLog = "account_log"
)

var (
	local_mongo_conn *mongo.Client
	LocalMongoWriter *dao.MongoWriter
	LocalMongoReader *dao.MongoReader
)

func Init(post func(string)) {
	init_mongo(post)
}

func init_mongo(post func(string)) {
	// 配置连接池参数
	clientOptions := options.Client().
		ApplyURI(config.Conf.LocalMongo.Url).
		SetMaxConnecting(2). // 允许快速建立连接
		SetMaxPoolSize(200).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(10 * time.Minute). // 连接最大空闲时间
		SetConnectTimeout(5 * time.Second).   // 连接超时时间
		SetSocketTimeout(15 * time.Second)    // Socket操作超时时间

	// clientOptions.SetPoolMonitor(&event.PoolMonitor{
	// 	Event: func(evt *event.PoolEvent) {
	// 		switch evt.Type {
	// 		case event.ConnectionCreated:
	// 			log.Debug("MongoDB connection created")
	// 		case event.GetFailed:
	// 			log.Warn("MongoDB get connection failed", zap.String("reason", evt.Reason))
	// 		case event.PoolCleared:
	// 			log.Warn("MongoDB pool cleared")
	// 		case event.ConnectionReturned:
	// 			log.Warn("MongoDB Returned connection:", zap.Uint64("connectionId", evt.ConnectionID))
	// 		case event.ConnectionClosed:
	// 			log.Warn("MongoDB Closed connection:", zap.Uint64("connectionId", evt.ConnectionID), zap.String("reason", evt.Reason))
	// 		}
	// 	},
	// })

	// clientOptions.SetMonitor(&event.CommandMonitor{
	// 	Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
	// 		log.Warnf("命令开始: 请求ID:%d, 数据库:%s, 命令名:%s, 命令内容: %s",
	// 			startedEvent.RequestID, startedEvent.DatabaseName, startedEvent.CommandName, startedEvent.Command)
	// 	},
	// 	Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
	// 		log.Warnf("命令成功: 请求ID:%d, 命令名:%s, 执行时间:%d 纳秒",
	// 			succeededEvent.RequestID, succeededEvent.CommandName, succeededEvent.Duration)
	// 	},
	// 	Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
	// 		log.Warnf("命令失败: 请求ID:%d, 命令名:%s, 执行时间:%d 纳秒, 错误: %v",
	// 			failedEvent.RequestID, failedEvent.CommandName, failedEvent.Duration, failedEvent.Failure)
	// 	},
	// })

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MongoDB url: %s err: %v", config.Conf.LocalMongo.Url, err))
	}

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to ping MongoDB url: %s err: %v", config.Conf.LocalMongo.Url, err))
	}
	local_mongo_conn = client
	LocalMongoWriter = dao.NewMongoWriter(
		local_mongo_conn,
		config.Conf.LocalMongo.WorkCount,
		config.Conf.LocalMongo.QueueSize,
		config.Conf.LocalMongo.BatchSize,
		config.Conf.LocalMongo.FlushSeconds,
		config.Conf.Debug,
		post,
	)

	LocalMongoReader = dao.NewMongoReader(local_mongo_conn)

	create_local_index()
}

func create_index(dataBase *mongo.Database, collectionName, filedName string, unique bool) {
	collecttion := dataBase.Collection(collectionName)
	if index, err := collecttion.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: filedName, Value: 1}},
		Options: options.Index().SetUnique(unique),
	}); err != nil {
		log.Panic("create index err", zap.String("name", collectionName), zap.Error(err))
	} else {
		log.Info("create index  succ", zap.String("name", collectionName), zap.String("index", index), zap.Bool("unique", unique))
	}
}

func create_compound_index(dataBase *mongo.Database, collectionName string, indexFields bson.D, unique bool) {
	collection := dataBase.Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys:    indexFields, // 联合索引字段[7](@ref)
		Options: options.Index().SetUnique(unique),
	}

	// 执行创建索引操作
	if index, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
		log.Panic("create compound index failed",
			zap.String("collect", collectionName),
			zap.Error(err))
	} else {
		log.Info("create compound index succ",
			zap.String("collect", collectionName),
			zap.String("index", index),
			zap.Bool("unique", unique))
	}
}

func create_local_index() {
	dataBase := local_mongo_conn.Database(config.Conf.LocalMongo.DB)
	// create_index(dataBase, CollectionName_UserAccount, "accountid", true)
	// create_index(dataBase, CollectionName_UserAccount, "userid", true)

	create_index(dataBase, CollectionName_User, "userid", true)
	create_index(dataBase, CollectionName_User, "accountid", true)
	//create_index(dataBase, CollectionName_User, "nick", true)
	create_index(dataBase, CollectionName_User, "serverid", false)

	create_index(dataBase, CollectionName_Order, "userid", false)

	create_compound_index(dataBase, CollectionName_FriendOp, bson.D{
		{Key: "tar_id", Value: 1},
		{Key: "op_id", Value: 1},
		{Key: "op_type", Value: 1}}, false)

	create_index(dataBase, CollectionName_GlobalMail, "mail_id", true)
}

func Stop() {
	log.Info("Local DB close start")
	LocalMongoWriter.Stop()
	local_mongo_conn.Disconnect(context.Background())
	log.Info("Local DB close succ")
}
