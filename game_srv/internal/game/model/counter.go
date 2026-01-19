package model

// import (
// 	"context"
// 	"gameserver/internal/config"
// 	"gameserver/internal/db"
// 	"gameserver/internal/publicconst"

// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	"go.mongodb.org/mongo-driver/bson"
// )

type Counter struct {
	Name string `bson:"_id"`
	Seq  uint64
}

// const (
// 	RankID = "rank_id"
// )

// func GenRankIdSeq() (uint64, error) {
// 	return genSeqByName(RankID, 0)
// }

// func genSeqByName(name string, initValue uint64) (uint64, error) {
// 	col := db.GetLocalClient().Database(config.Conf.GetLocalDB()).Collection(publicconst.LOCAL_CONTRACT)

// 	filter := bson.M{"_id": name}
// 	update := bson.M{
// 		"$setOnInsert": bson.M{"seq": initValue},
// 		"$inc":         bson.M{"seq": uint64(1)},
// 	}
// 	opts := options.FindOneAndUpdate().
// 		SetUpsert(true).
// 		SetReturnDocument(options.After)

// 	var ret Counter
// 	err := col.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&ret)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ret.Seq, nil
// }
