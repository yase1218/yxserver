package model

import (
	"msg"
)

type (
	Contract struct {
		AccountId uint64 `bson:"_id"`

		TaskIds      []uint32          `bson:"task_ids"`
		TaskId       uint32            `bson:"task_id"`
		TaskType     msg.ConditionType `bson:"task_type"`
		StageEventId uint32            `bson:"stage_event_id"`
		FinishCount  uint32            `bson:"finish_count"`
		Reward       bool              `bson:"reward"`
		Point        uint32            `bson:"point"`

		SignNum        uint32 `bson:"sign_num"`
		RandNum        uint32 `bson:"rand_num"`
		RandDiamondNum uint32 `bson:"rand_diamond_num"`
		ResetDate      uint32 `bson:"reset_date"`
	}

	ContractMongoModel struct{}
)

// var (
// 	ContractModel = new(ContractMongoModel)
// )

// func GetContractModel() *ContractMongoModel {
// 	return ContractModel
// }

// func (m *ContractMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *ContractMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_CONTRACT)
// }

// func (m *ContractMongoModel) Create(data *Contract) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *ContractMongoModel) Get(accountId int64) (*Contract, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data Contract
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *ContractMongoModel) Update(accountId int64, update bson.M) error {
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

// func (m *ContractMongoModel) UpdatePointsToZero(accountIds []int64, point uint32) error {
// 	if len(accountIds) == 0 {
// 		return nil
// 	}

// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().UpdateMany(ctx,
// 		bson.M{"_id": bson.M{"$in": accountIds}},
// 		bson.M{"$set": bson.M{"point": point}})
// 	return err
// }
