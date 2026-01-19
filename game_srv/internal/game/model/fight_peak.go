package model

import "time"

type (
	PeakFight struct {
		BattleMatchId uint32 `bson:"battle_match_id"` // battle_match 表id
		RoomId        uint64 `bson:"room_id"`         // 当前所在房间id
		Cup           uint32 `bson:"cup"`             // 奖杯数量
		Streak        uint32 `bson:"streak"`          // 连胜
		FreeTimes     uint32 `bson:"free_times"`      // 已使用免费次数
		ResetDate     uint32 `bson:"reset_date"`      // 重置日期
		Season        uint32 `bson:"season"`          // 赛季
	}

	PeakFightMongoModel struct{}
)

func (f *PeakFight) Reset(now time.Time) {
	f.BattleMatchId = 1
	f.RoomId = 0
	f.Cup = 0
	f.Streak = 0
	f.FreeTimes = 0
}

// func NewPeakFight(accountId int64, season uint32) *PeakFight {
// 	return &PeakFight{
// 		BattleMatchId: 1,
// 		Season:        season,
// 	}
// }

// var (
// 	PeakFightModel = new(PeakFightMongoModel)
// )

// func GetPeakFightModel() *PeakFightMongoModel {
// 	return PeakFightModel
// }

// func (m *PeakFightMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *PeakFightMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_PEAK_FIGHT)
// }

// func (m *PeakFightMongoModel) CreatePeakFight(fight *PeakFight) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, fight)
// 	return err
// }

// func (m *PeakFightMongoModel) GetPeakFight(accountId int64) (*PeakFight, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data PeakFight
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *PeakFightMongoModel) UpdatePeakFight(accountId int64, update bson.M) error {
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
