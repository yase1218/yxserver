package model

import (
	"msg"
	"time"
)

type (
	AllianceRank struct {
		//ID         uint64               `bson:"_id"`
		AccountId  int64                `bson:"account_id"`  // 账号
		Data       uint64               `bson:"data"`        // 榜单分数
		ServerId   uint32               `bson:"server_id"`   // 区服id
		AllianceId uint32               `bson:"alliance_id"` // 工会id
		Date       int                  `bson:"date"`        // 榜单周期时间
		Type       msg.AllianceRankType `bson:"type"`        // 榜单类型
		CreateAt   time.Time            `bson:"create_at"`
	}

	// AllianceRankMongoModel struct{}
)

// var (
// 	AllianceRankModel = new(AllianceRankMongoModel)
// )

// func GetAllianceRankModel() *AllianceRankMongoModel {
// 	return AllianceRankModel
// }

// func (m *AllianceRankMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *AllianceRankMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_RankAlliance)
// }

// func (m *AllianceRankMongoModel) Create(data *AllianceRank) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *AllianceRankMongoModel) Get(accountId int64) (*AllianceRank, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data AllianceRank
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *AllianceRankMongoModel) Update(filter any, update bson.M) error {
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

// 	_, err := m.GetCol().UpdateOne(ctx, filter, update)
// 	return err
// }

// func (m *AllianceRankMongoModel) UpsertData(filter bson.M, update bson.M) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	opts := options.FindOneAndUpdate().
// 		SetUpsert(true).
// 		SetReturnDocument(options.After)

// 	var result bson.M
// 	err := m.GetCol().FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (m *AllianceRankMongoModel) GetSettlementList(serverId uint32, tp msg.AllianceRankType, date int64) ([]*AllianceRank, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	opts := options.Find().
// 		SetSort(bson.D{
// 			{"data", -1},
// 			{"create_at", 1},
// 		})
// 	opts.SetLimit(100)

// 	filter := bson.M{
// 		"server_id": serverId,
// 		"date":      date,
// 	}

// 	if tp != msg.AllianceRankType_Alliance_Rank_Empty {
// 		filter["type"] = tp
// 	}

// 	cur, err := m.GetCol().Find(ctx, filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ret []*AllianceRank
// 	for {
// 		if !cur.Next(ctx) {
// 			break
// 		}
// 		var rank = AllianceRank{}
// 		if err = cur.Decode(&rank); err == nil {
// 			ret = append(ret, &rank)
// 		}
// 	}

// 	return ret, nil
// }

// func (m *AllianceRankMongoModel) GetList(serverId, allianceId uint32, tp msg.AllianceRankType) ([]*AllianceRank, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	opts := options.Find().
// 		SetSort(bson.D{
// 			{"data", -1},
// 			{"create_at", 1},
// 		})
// 	opts.SetLimit(100)

// 	filter := bson.M{
// 		"server_id": serverId,
// 		"type":      tp,
// 	}
// 	if tp != msg.AllianceRankType_Alliance_Rank_Boss {
// 		filter["alliance_id"] = allianceId
// 	}
// 	cur, err := m.GetCol().Find(ctx, filter, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ret []*AllianceRank
// 	for {
// 		if !cur.Next(ctx) {
// 			break
// 		}
// 		var rank = AllianceRank{}
// 		if err = cur.Decode(&rank); err == nil {
// 			ret = append(ret, &rank)
// 		}
// 	}

// 	return ret, nil
// }

// func (m *AllianceRankMongoModel) Delete(filter any) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().DeleteOne(ctx, filter)
// 	return err
// }

// func (m *AllianceRankMongoModel) GetUserRanking(accountId int64, allianceId uint32, tp msg.AllianceRankType) (int64, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	filter := bson.M{
// 		"account_id":  accountId,
// 		"alliance_id": allianceId,
// 		"type":        tp,
// 		"server_id":   config.Conf.ServerId,
// 		"date":        tools.GetYearWeekByOffset(time.Now(), int(template.GetSystemItemTemplate().RefreshHour)),
// 	}
// 	var rankInfo AllianceRank
// 	err := m.GetCol().FindOne(ctx, filter).Decode(&rankInfo)
// 	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
// 		return 0, err
// 	}
// 	if rankInfo.AccountId == 0 {
// 		return -1, nil
// 	}

// 	count, err := m.GetCol().CountDocuments(ctx, bson.M{
// 		"data":      bson.M{"$gt": rankInfo.Data},
// 		"type":      tp,
// 		"server_id": config.Conf.ServerId,
// 		"date":      rankInfo.Date,
// 	})
// 	if err != nil {
// 		return 0, err
// 	}

// 	return int64(count) + 1, nil
// }
