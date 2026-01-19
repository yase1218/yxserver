package model

// import (
// 	"gameserver/internal/config"
// 	"gameserver/internal/db"
// 	"gameserver/internal/publicconst"
// 	"strings"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

type (
	Cumulative struct {
		AccountId int64         `bson:"_id"`
		Data      map[int][]int // 类型:参数
	}

	//CumulativeMongoModel struct{}
)

// var (
// 	CumulativeModel = new(CumulativeMongoModel)
// )

// func GetCumulativeModel() *CumulativeMongoModel {
// 	return CumulativeModel
// }

// func (m *CumulativeMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *CumulativeMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_CUMULATIVE)
// }

// func (m *CumulativeMongoModel) Create(data *Cumulative) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *CumulativeMongoModel) Get(accountId int64) (*Cumulative, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data Cumulative
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *CumulativeMongoModel) Update(accountId int64, update bson.M) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	hasOperator := false
// 	for k := range update {
// 		if strings.HasPrefix(k, "$") {
// 			hasOperator = true
// 			break
// 		}
// 	}
// 	if !hasOperator {
// 		update = bson.M{"$set": update}
// 	}

// 	_, err := m.GetCol().UpdateOne(ctx, bson.M{"_id": accountId}, update)
// 	return err
// }
