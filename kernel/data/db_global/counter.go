package db_global

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
	COUNTER_ACCOUNT_ID     = "account_id"
	COUNTER_ACCOUNT_LOG_ID = "account_log_id"
	COUNTER_SERVER_INFO_ID = "server_info_id"
	COUNTER_WHITE_ID       = "white_id"
)

func GenAccountIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_ACCOUNT_ID, 9000000)
}
func GenServerInfoIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_SERVER_INFO_ID, 0)
}
func GenWhiteIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_WHITE_ID, 0)
}
func GenAccountLogIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_ACCOUNT_LOG_ID, 0)
}

func genSeqByName(name string, initValue uint64) (uint64, error) {
	col := mongo.DB(GetDB()).Collection(COL_COUNTER)

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
