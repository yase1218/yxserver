package mongo

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

type Test struct {
	ID   uint64    `bson:"_id"`
	Time time.Time `bson:"time"`
}

type TestMongoModel struct{}

var (
	TestMongo = &TestMongoModel{}
	TestOnce  = sync.Once{}
)

func GetTestMongo() *TestMongoModel {
	TestOnce.Do(CWWTestInitIndex)
	return TestMongo
}

func CWWTestInitIndex() {
	//unique := true
	//err := GetTestCol().CreateIndexes(context.Background(), []options.IndexModel{
	//	{
	//		Key:          []string{"userId", "version"},
	//		IndexOptions: &options2.IndexOptions{Unique: &unique},
	//	},
	//})
	//if err != nil {
	//	log.Error("cw_big_high_ladder_rank CreateIndexes err", zap.String("err", err.Error()))
	//}
}

func GetTestCol() *qmgo.Collection {
	return mongo.DB(DB_TEST).Collection(COL_TEST)
}

func (m *TestMongoModel) InsertOne(data *Test) error {
	_, err := GetTestCol().Upsert(context.Background(), bson.M{}, data)
	if err != nil {
		return err
	}
	return nil
}

func (m *TestMongoModel) LoadOne(id uint64) (*Test, error) {
	var data *Test
	var err error
	err = GetTestCol().Find(context.Background(), bson.M{"_id": id}).One(&data)
	return data, err
}
