package mongo

import (
	"context"
	"errors"
	"github.com/v587-zyf/gc/db/mongo"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Counter struct {
	Name string `bson:"_id"`
	Seq  uint64
}

const (
	COUNTER_TEST = "test_id"
)

// GenBattleTokenSeq 获取自增id
func GetTestIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_TEST, 1)
}

func genSeqByName(name string, initValue uint64) (uint64, error) {
	col := mongo.DB(DB_TEST).Collection(COL_COUNTER)

	ret := Counter{}

	err := col.Find(context.Background(), bson.M{"_id": name}).Apply(qmgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    false,
		ReturnNew: true,
	}, &ret)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		_, err = col.InsertOne(context.Background(), &Counter{Name: name, Seq: initValue})
		if err != nil && !qmgo.IsDup(err) {
			return 0, err
		}

		err = col.Find(context.Background(), bson.M{"_id": name}).Apply(qmgo.Change{
			Update:    bson.M{"$inc": bson.M{"seq": 1}},
			Upsert:    false,
			ReturnNew: true,
		}, &ret)
	}

	return ret.Seq, err
}
