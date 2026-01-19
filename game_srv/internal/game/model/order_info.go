package model

import (
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"kernel/dao"
	"math"
	"time"

	"github.com/zy/game_data/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OrderUnShipment = iota // 未发货
	OrderVerified          // 已验证未发货
	OrderShipment          // 已发货
)

const RedisOrderTime = 86400 * time.Second // 一天

type OrderInfo struct {
	Id        primitive.ObjectID `bson:"_id"`
	ChargeID  int                `bson:"charge_id"`
	UserId    uint64             `bson:"userid"`
	OrderId   string             `bson:"order_id"`   // 游戏订单号(服务端生成)
	ThirdNo   string             `bson:"third_no"`   // 第三方订单号（iOS渠道为雷霆订单号）
	ChannelNo string             `bson:"channel_no"` // 渠道编号
	Currency  string             `bson:"currency"`   // 币种
	Money     int                `bson:"money"`      // 第三方返回金额（分）
	ProductId string             `bson:"product_id"` // 价格档位/内购id 例：com.leiting.1.kulu
	AccountId string             `bson:"accountid"`  // 手游账号ID 雷霆OpenId
	Status    int                `bson:"status"`     // 发货状态
	CreateAt  time.Time          `bson:"create_at"`
	UpdateAt  time.Time          `bson:"update_at"`
}

func (order *OrderInfo) async_filed_to_db(filed_name string, value interface{}) {
	db_name := config.Conf.LocalMongo.DB
	collection := db.CollectionName_Order
	op := &dao.WriteOperation{
		Database:   db_name,
		Collection: collection,
		Type:       dao.Upsert,
		Filter:     bson.M{"_id": order.Id},
		Tms:        time.Now().UnixMilli(),
		Uuid:       order.UserId,
	}
	op.Update = bson.M{"$set": bson.M{filed_name: value}}
	db.LocalMongoWriter.AsyncWrite(op)
}

// 保存发货状态
func (order *OrderInfo) SaveStatus() {
	order.async_filed_to_db("status", order.Status)
}

func GetProductId(cfg *template.JCharge) string {
	return fmt.Sprintf(getDecimalPlaces(float64(cfg.CostRMB)), cfg.CostRMB, config.Conf.Leiting.Game)
}

// 获取格式化字符串
func getDecimalPlaces(f float64) string { //
	if math.Floor(f) == f {
		return "com.leiting.%.0f.%s"
	}
	if int(f*100)%10 == 0 {
		return "com.leiting.%.1f.%s"
	}
	return "com.leiting.%.2f.%s"
}
