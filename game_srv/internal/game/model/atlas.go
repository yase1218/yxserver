package model

type (
	AtlasUnit struct {
		Id     uint32 `bson:"id"`     // 图鉴id
		Reward bool   `bson:"reward"` // 是否已领取奖励
	}
	Atlas struct {
		AccountId int64 `bson:"_id"`

		Data map[uint32]*AtlasUnit `bson:"data"` // id:data
	}
)

// type AtlasMongoModel struct{}

// var (
// 	AtlasModel = &AtlasMongoModel{}
// )

// func GetAtlasModel() *AtlasMongoModel {
// 	return AtlasModel
// }

// func (m *AtlasMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *AtlasMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_ATLAS)
// }

// func (m *AtlasMongoModel) Create(data *Atlas) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *AtlasMongoModel) Get(accountId int64) (*Atlas, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data Atlas
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *AtlasMongoModel) Update(accountId int64, update bson.M) error {
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

// func (m *AtlasMongoModel) UpdateOnlyNewData(accountId int64, newData map[uint32]*AtlasUnit) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	update := bson.M{
// 		"$set": bson.M{},
// 	}

// 	for id, unit := range newData {
// 		key := "data." + strconv.Itoa(int(id))
// 		update["$setOnInsert"] = bson.M{
// 			"_id": accountId,
// 		}
// 		update["$set"].(bson.M)[key] = unit
// 	}

// 	opts := options.Update().SetUpsert(true)
// 	_, err := m.GetCol().UpdateOne(
// 		ctx,
// 		bson.M{"_id": accountId},
// 		update,
// 		opts,
// 	)
// 	return err
// }

// func (m *AtlasMongoModel) UpdateReward(accountID int64, atlasID uint32, reward bool) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	key := "data." + strconv.Itoa(int(atlasID)) + ".reward"
// 	update := bson.M{
// 		"$set": bson.M{
// 			key: reward,
// 		},
// 	}

// 	_, err := m.GetCol().UpdateOne(ctx, bson.M{"_id": accountID}, update)
// 	return err
// }
