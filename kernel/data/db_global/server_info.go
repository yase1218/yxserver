package db_global

import (
	"context"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type ServerInfo struct {
	SerID             uint64    `bson:"_id"`
	SerName           string    `bson:"ser_name"`       // 服务器名称
	SerAddr           string    `bson:"ser_addr"`       // 服务器地址
	RegisterNum       int       `bson:"register_num"`   // 已注册人数
	RegisterLimit     int       `bson:"register_limit"` // 注册人数上限
	OpenTime          int64     `bson:"open_time"`      // 开服时间
	DisplayTime       int64     `bson:"display_time"`   // 显示时间
	Whitelist         bool      `bson:"whitelist"`      // 是否开启白名单
	IsMaintain        bool      `bson:"is_maintain"`    // 是否维护中
	MaintainDesc      string    `bson:"maintain_desc"`  // 维护描述
	MaintainBeginTime int64     `bson:"maintain_begin_time"`
	MaintainEndTime   int64     `bson:"maintain_end_time"`
	InterceptArea     string    `bson:"intercept_area"` // 拦截区域
	OnlineNum         int       `bson:"online_num"`     // 在线人数
	CreateAt          time.Time `bson:"create_at"`
	UpdateAt          time.Time `bson:"update_at"`
}

// 0->关闭 1->正常 2->爆满
func (s *ServerInfo) Status() uint32 {
	status := uint32(1)
	if s.SerAddr == "" || s.IsMaintain {
		status = 0
	}
	return status
}

type ServerInfoDBModel struct{}

var (
	ServerInfoModel = &ServerInfoDBModel{}
)

func GetServerInfoModel() *ServerInfoDBModel {
	return ServerInfoModel
}

func GetServerInfoCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_SERVER_INFO)
}

func (m *ServerInfoDBModel) Upsert(data *ServerInfo) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.SerID}
	return GetServerInfoCol().Upsert(context.Background(), filter, data)
}

func (m *ServerInfoDBModel) UpdateServerInfo(data *ServerInfo) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.SerID}

	err := GetServerInfoCol().UpdateOne(context.Background(), filter,
		bson.M{"$set": bson.M{"register_num": data.RegisterNum, "update_at": time.Now()}})
	return nil, err
}

func (m *ServerInfoDBModel) GetOne(serverId uint64) (*ServerInfo, error) {
	var data *ServerInfo
	var err error
	filter := bson.M{"_id": serverId}
	err = GetServerInfoCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *ServerInfoDBModel) GetAll(filter any) ([]*ServerInfo, error) {
	var dataSlice []*ServerInfo
	err := GetServerInfoCol().Find(context.Background(), filter).Sort("-_id").All(&dataSlice)
	return dataSlice, err
}
