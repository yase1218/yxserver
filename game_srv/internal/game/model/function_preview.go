package model

import (
	"msg"
)

type (
	FunctionPreview struct {
		AccountId uint64 `bson:"_id"`

		Data map[uint32]msg.TaskState `bson:"data"`
	}

	//FunctionPreviewMongoModel struct{}
)

// var (
// 	FunctionPreviewModel = new(FunctionPreviewMongoModel)
// )

// func GetFunctionPreviewModel() *FunctionPreviewMongoModel {
// 	return FunctionPreviewModel
// }

// func (m *FunctionPreviewMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *FunctionPreviewMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_FUNCTION_PREVIEW)
// }

// func (m *FunctionPreviewMongoModel) Create(data *FunctionPreview) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *FunctionPreviewMongoModel) Get(accountId int64) (*FunctionPreview, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data FunctionPreview
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *FunctionPreviewMongoModel) Update(accountId int64, update bson.M) error {
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

// func (m *FunctionPreviewMongoModel) UpdateData(accountId int64, data map[uint32]msg.TaskState) error {
// 	if len(data) == 0 {
// 		return nil
// 	}

// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	update := bson.M{}
// 	for k, v := range data {
// 		key := "data." + string(rune(k))
// 		update[key] = v
// 	}

// 	_, err := m.GetCol().UpdateOne(ctx,
// 		bson.M{"_id": accountId},
// 		bson.M{"$set": update})
// 	return err
// }
