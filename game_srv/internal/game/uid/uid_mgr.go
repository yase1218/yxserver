package uid

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SMaxUserId struct {
	Id uint64 `bson:"maxUserId"`
}

var (
	MaxUserId uint64
)

const UidEra = uint64(1000000)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 构建聚合管道
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "maxUserId", Value: bson.D{
					{Key: "$max", Value: "$userid"},
				}}, 
			}},
		},
	}

	// 执行聚合查询
	cursor, err := db.LocalMongoReader.Aggregate(ctx, config.Conf.LocalMongo.DB, db.CollectionName_User, pipeline)
	if err != nil {
		panic(fmt.Sprintf("aggregate failed: %v", err))
	}
	defer cursor.Close(ctx)

	// 解码结果
	var result []SMaxUserId
	if err = cursor.All(ctx, &result); err != nil {
		panic(fmt.Sprintf("decode failed: %v", err))
	}

	// 处理结果：如果集合为空，结果可能为空切片
	if len(result) > 0 {
		MaxUserId = result[0].Id
	} else {
		MaxUserId = uint64(config.Conf.ServerId) * UidEra
	}
}

func GenUserId() uint64 {
	MaxUserId++
	return MaxUserId
}
