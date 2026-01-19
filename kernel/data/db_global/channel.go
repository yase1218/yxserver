package db_global

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type EnvParam struct {
	Name     string `uri:"name" comment:"环境名"`
	LoginUrl string `uri:"loginUrl" comment:"登录地址"`
	ResUrl   string `uri:"resUrl" comment:"资源地址"`
	ResVer   string `uri:"resVer" comment:"资源版本号"`
}

type Channel struct {
	ChannelId  uint32 `json:"_id"`
	Name       string `json:"name" `
	Remark     string `json:"remark"`
	UpdateTime uint32 `json:"updateTime"`
	//StrUpdateTime string      `json:"strUpdateTime"`
	Envs []*EnvParam `json:"envs"`
}

type ChannelDBModel struct{}

var (
	ChannelModel = &ChannelDBModel{}
)

func GetChannelModel() *ChannelDBModel {
	return ChannelModel
}

func GetChannelCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_CHANNEL)
}

func (m *ChannelDBModel) Upsert(data *Channel) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ChannelId}
	return GetChannelCol().Upsert(context.Background(), filter, data)
}

func (m *ChannelDBModel) GetOne(channelId uint32) (*Channel, error) {
	var data *Channel
	var err error
	filter := bson.M{"_id": channelId}
	err = GetChannelCol().Find(context.Background(), filter).One(&data)
	return data, err
}
