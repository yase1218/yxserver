package model

type (
	DesertFight struct {
		AccountId   uint64   `bson:"_id"`
		RetroSignin uint32   `bson:"retro_signin"` // 累计补签次数
		KillTimes   uint32   `bson:"kill_times"`   // 击杀数量
		RewardTimes []uint32 `bson:"reward_times"` // 已领取段位
		ResetDate   uint32   `bson:"reset_date"`   // 重置日期
	}

	// DesertFightMongoModel struct{}
)

func NewDesertFight(accountId uint64) *DesertFight {
	return &DesertFight{
		AccountId:   accountId,
		RewardTimes: make([]uint32, 0),
	}
}

// var (
// 	DesertFightModel = new(DesertFightMongoModel)
// )

// func GetDesertFightModel() *DesertFightMongoModel {
// 	return DesertFightModel
// }

// func (m *DesertFightMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *DesertFightMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_DESERT_FIGHT)
// }

// func (m *DesertFightMongoModel) CreateDesertFight(fight *DesertFight) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, fight)
// 	return err
// }

// func (m *DesertFightMongoModel) GetDesertFight(accountId int64) (*DesertFight, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data DesertFight
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *DesertFightMongoModel) UpdateDesertFight(accountId int64, update bson.M) error {
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
