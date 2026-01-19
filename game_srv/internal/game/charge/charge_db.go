package charge

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"gameserver/internal/game/model"
	"kernel/dao"
	"kernel/tools"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// 加载30天以内的订单
func LoadAllOrder(ctx context.Context) {
	t := tools.GetOffsetDays(-30)
	cursor, err := db.LocalMongoReader.Find(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_Order,
		bson.M{
			"create_at": bson.M{"$gte": t},
		},
	)
	tools.GetCurTime()

	if err != nil {
		panic(fmt.Errorf("load all order failed,%s", err.Error()))
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		order := &model.OrderInfo{}
		err = cursor.Decode(order)
		if err != nil {
			panic(fmt.Errorf("load all order decode failed,%s", err.Error()))
		}
		GetOrderManager().AddOrder(order)
	}
}

func CreateOrder(order *model.OrderInfo) error { // 异步落地
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	dbName := config.Conf.LocalMongo.DB
	collection := db.CollectionName_Order
	op := &dao.WriteOperation{
		Database:   dbName,
		Collection: collection,
		Type:       dao.Insert,
		//Filter:     bson.M{"uuid": order.OrderId},
		Document: order,
		Tms:      time.Now().UnixMilli(),
		Cover:    true,
	}
	err := db.LocalMongoWriter.AsyncWrite(op)
	return err
}

// func GetOrderRedisKey(orderID string, serverID uint32) string {
// 	return fmt.Sprintf("charge_order:%s:%d", orderID, serverID)
// }

// func LoadVerifiedOrders(ctx context.Context, uid uint64) []*model.OrderInfo {
// 	var data []*model.OrderInfo
// 	cursor, err := db.LocalMongoReader.Find(
// 		ctx,
// 		config.Conf.LocalMongo.DB,
// 		"order_info",
// 		bson.M{ /*"status": model.OrderVerified, */ "userid": uid},
// 	)
// 	if err != nil {
// 		panic(fmt.Errorf("load order failed,%s", err.Error()))
// 	}
// 	defer cursor.Close(ctx)
// 	for cursor.Next(ctx) {
// 		order := &model.OrderInfo{}
// 		err = cursor.Decode(order)
// 		if err != nil {
// 			panic(fmt.Errorf("load order decode failed,%s", err.Error()))
// 		}
// 		data = append(data, order)
// 	}
// 	return data
// }
