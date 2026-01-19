package db_global

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type White struct {
	ID  uint64 `bson:"_id"`
	Uid string `bson:"uid"` // 账号
}

type WhiteDBModel struct{}

var (
	WhiteModel = &WhiteDBModel{}
)

func GetWhiteModel() *WhiteDBModel {
	return WhiteModel
}

func GetWhiteCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_WHITE)
}

func (m *WhiteDBModel) Upsert(data *White) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	return GetWhiteCol().Upsert(context.Background(), filter, data)
}

func (m *WhiteDBModel) GetOne(id uint64) (*White, error) {
	var data *White
	var err error
	filter := bson.M{"_id": id}
	err = GetWhiteCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *WhiteDBModel) GetAll(filter any) ([]*White, error) {
	var dataSlice []*White
	var err error
	err = GetWhiteCol().Find(context.Background(), filter).All(&dataSlice)
	return dataSlice, err
}

func (m *WhiteDBModel) InsertMany(data []*White) ([]uint64, error) {
	result, insertErr := GetWhiteCol().InsertMany(context.Background(), data)
	if insertErr != nil {
		return []uint64{}, insertErr
	}
	ids := make([]uint64, len(result.InsertedIDs))
	for i := 0; i < len(ids); i++ {
		ids[i] = uint64(result.InsertedIDs[i].(int64))
	}
	return ids, nil
}

func (m *WhiteDBModel) LoadAllUid() ([]string, error) {
	var dataSlice []*White
	if err := GetWhiteCol().Find(context.Background(), bson.M{}).All(&dataSlice); err != nil {
		return nil, err
	}
	uidSlice := make([]string, 0, len(dataSlice))
	for _, data := range dataSlice {
		uidSlice = append(uidSlice, data.Uid)
	}
	return uidSlice, nil
}
