package db_global

import (
	"context"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgoptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type AccountLog struct {
	ID        primitive.ObjectID `bson:"_id"`
	OpenId    string             `bson:"open_id"`
	Platform  uint32             `bson:"platform"`
	ChannelId uint32             `bson:"channel_id"`
	UserId    uint64             `bson:"userid"`
	CreateAt  time.Time          `bson:"create_at"`
	UpdateAt  time.Time          `bson:"update_at"` // 最后更新时间
}

type AccountLogDBModel struct{}

var (
	AccountLogModel = &AccountLogDBModel{}
)

func AccountLogCreateIndex() {
	collecttion := GetAccountLogCol()
	err := collecttion.CreateOneIndex(context.Background(), options.IndexModel{
		Key: []string{
			"open_id",
			"platform",
			"channel_id",
		},
		IndexOptions: mgoptions.Index().SetUnique(true),
	})
	if err != nil {
		log.Panic("create account_log index err", zap.Error(err))
	}
}

func GetAccountLogModel() *AccountLogDBModel {
	return AccountLogModel
}

func GetAccountLogCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_ACCOUNT_LOG)
}

func (m *AccountLogDBModel) Upsert(data *AccountLog) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	return GetAccountLogCol().Upsert(context.Background(), filter, data)
}

func (m *AccountLogDBModel) GetOneByKey(openId string, channel, platform uint32) (*AccountLog, error) {
	var data *AccountLog
	var err error
	filter := bson.M{"open_id": openId, "channel_id": channel, "platform": platform}
	err = GetAccountLogCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *AccountLogDBModel) GetOneByPlatform(platform string) (*AccountLog, error) {
	var data *AccountLog
	var err error
	filter := bson.M{"platform": platform}
	err = GetAccountLogCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *AccountLogDBModel) GetAll(filter any) ([]*AccountLog, error) {
	var dataSlice []*AccountLog
	var err error
	err = GetAccountLogCol().Find(context.Background(), filter).All(&dataSlice)
	return dataSlice, err
}

// func (m *AccountLogDBModel) GetUserPlatform(openId, platform string) (*AccountLog, error) {
// 	filter := bson.M{"open_id": openId, "platform": platform}
// 	data := new(AccountLog)
// 	err := GetAccountCol().Find(context.Background(), filter).One(&data)
// 	return data, err
// }
