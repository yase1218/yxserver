package model

import (
	"msg"
	"time"
)

type (
	Fight struct {
		FightId      uint32          `bson:"fightId"`      // 当前所在战斗服id
		FightStageId int             `bson:"fightStageId"` // 当前stageId
		FightStartAt time.Time       `bson:"fightStartAt"` // 战斗开始时间
		Weapons      []uint32        `bson:"weapons"`      // 选择的武器
		Faction      msg.FactionType `bson:"faction"`      // 流派
	}

	FightMongoModel struct{}
)

func (f *Fight) Clear() {
	f.FightId = 0
	f.FightStageId = 0
	f.FightStartAt = time.Time{}
}

// var (
// 	FightModel = new(FightMongoModel)
// )

// func GetFightModel() *FightMongoModel {
// 	return FightModel
// }

// func (m *FightMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.LocalMongo.DB)
// }

// func (m *FightMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_FIGHT)
// }

// func (m *FightMongoModel) CreateFight(fight *Fight) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, fight)
// 	return err
// }

// func (m *FightMongoModel) GetFight(accountId int64) (*Fight, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data Fight
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *FightMongoModel) UpdateFight(accountId int64, update bson.M) error {
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
