package model

import (
	"gameserver/internal/publicconst"

	"github.com/zy/game_data/template"
)

const (
	LuckSaleDefJackpotId = 1
)

type (
	LuckSaleUnit struct {
		Ids   []uint32 `bson:"ids"`
		Times uint32   `bson:"times"`
	}

	LuckSaleTaskUnit struct {
		TaskId     uint32 `bson:"taskId"`
		Value      uint32 `bson:"value"`
		Extra      int64  `bson:"extra"`
		State      uint32 `bson:"state"`
		UpdateTime uint32 `bson:"updateTime"`
		Reward     bool   `bson:"reward"`
	}

	LuckSale struct {
		AccountId uint64 `bson:"_id"`

		Jackpot int                   `bson:"jackpot"`
		Data    map[int]*LuckSaleUnit `bson:"data"`
		Task    []*LuckSaleTaskUnit   `bson:"task"`
	}
)

func NewLuckSale(accountId uint64) *LuckSale {
	taskCfgs := template.GetTaskTemplate().GetTaskByType(publicconst.LUCK_SALE_TASK)

	dbData := &LuckSale{
		AccountId: accountId,
		Jackpot:   LuckSaleDefJackpotId,
		Data:      make(map[int]*LuckSaleUnit),
		Task:      make([]*LuckSaleTaskUnit, 0, len(taskCfgs)),
	}
	dbData.Data[LuckSaleDefJackpotId] = NewLuckSaleUnit()
	for _, v := range taskCfgs {
		if v.Data.PreTask != 0 {
			continue
		}
		dbData.Task = append(dbData.Task, &LuckSaleTaskUnit{TaskId: v.Data.Id})
	}

	return dbData
}

func NewLuckSaleUnit() *LuckSaleUnit {
	return &LuckSaleUnit{
		Ids: make([]uint32, 0),
	}
}

func (a *LuckSaleTaskUnit) GetTaskId() uint32 {
	return a.TaskId
}

func (a *LuckSaleTaskUnit) GetTaskValue() uint32 {
	return a.Value
}

func (a *LuckSaleTaskUnit) SetTaskValue(value uint32) {
	a.Value = value
}

func (a *LuckSaleTaskUnit) AddTaskValue(add uint32) {
	a.Value += add
}

func (a *LuckSaleTaskUnit) SetTaskCompleteTime(time uint32) {
	a.UpdateTime = time
}

func (a *LuckSaleTaskUnit) GetTaskCompleteTime() uint32 {
	return a.UpdateTime
}

func (a *LuckSaleTaskUnit) GetTaskState() publicconst.TaskState {
	return publicconst.TaskState(a.State)
}

func (a *LuckSaleTaskUnit) SetTaskState(state publicconst.TaskState) {
	a.State = uint32(state)
}

func (a *LuckSaleTaskUnit) SetExtraPara(value uint32) {
	a.Extra = int64(value)
}

func (a *LuckSaleTaskUnit) GetExtraPara() uint32 {
	return uint32(a.Extra)
}

// type LuckSaleMongoModel struct{}

// var (
// 	LuckSaleModel = &LuckSaleMongoModel{}
// )

// func GetLuckSaleModel() *LuckSaleMongoModel {
// 	return LuckSaleModel
// }

// func (m *LuckSaleMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *LuckSaleMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_LUCK_SALE)
// }

// func (m *LuckSaleMongoModel) Create(data *LuckSale) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *LuckSaleMongoModel) Get(accountId int64) (*LuckSale, error) {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	var data LuckSale
// 	if err := m.GetCol().FindOne(ctx, bson.M{"_id": accountId}).Decode(&data); err != nil || data.AccountId == 0 {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (m *LuckSaleMongoModel) Update(accountId int64, update bson.M) error {
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

// func (m *LuckSaleMongoModel) UpdateOnlyNewData(accountId int64, newData map[uint32]*LuckSaleUnit) error {
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

// func (m *LuckSaleMongoModel) UpdateTask(accountId int64, data *LuckSaleTaskUnit) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	update := bson.M{
// 		"$set": bson.M{
// 			"task.$[task].value":      data.Value,
// 			"task.$[task].extra":      data.Extra,
// 			"task.$[task].state":      data.State,
// 			"task.$[task].updateTime": data.UpdateTime,
// 			"task.$[task].reward":     data.Reward,
// 		},
// 	}

// 	arrayFilters := options.ArrayFilters{
// 		Filters: []interface{}{
// 			bson.M{"task.taskId": data.TaskId},
// 		},
// 	}

// 	res := m.GetCol().FindOneAndUpdate(
// 		ctx,
// 		bson.M{"_id": accountId},
// 		update,
// 		options.FindOneAndUpdate().SetArrayFilters(arrayFilters),
// 	)

// 	if res.Err() != nil {
// 		log.Error("UpdateTask", zap.Int64("accountId", accountId), zap.Error(res.Err()))
// 		return res.Err()
// 	}

// 	return nil
// }

// func (m *LuckSaleMongoModel) BatchUpdateTasks(accountId int64, tasks []*LuckSaleTaskUnit) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	if len(tasks) == 0 {
// 		return nil
// 	}

// 	update := bson.M{"$set": bson.M{}}

// 	arrayFilters := make([]interface{}, 0, len(tasks))

// 	for i, task := range tasks {
// 		fieldPrefix := "task.$[task" + strconv.Itoa(i) + "]"
// 		update["$set"].(bson.M)[fieldPrefix+".value"] = task.Value
// 		update["$set"].(bson.M)[fieldPrefix+".extra"] = task.Extra
// 		update["$set"].(bson.M)[fieldPrefix+".state"] = task.State
// 		update["$set"].(bson.M)[fieldPrefix+".updateTime"] = task.UpdateTime
// 		update["$set"].(bson.M)[fieldPrefix+".reward"] = task.Reward

// 		arrayFilters = append(arrayFilters, bson.M{"task" + strconv.Itoa(i) + ".taskId": task.TaskId})
// 	}

// 	res := m.GetCol().FindOneAndUpdate(
// 		ctx,
// 		bson.M{"_id": accountId},
// 		update,
// 		options.FindOneAndUpdate().SetArrayFilters(
// 			options.ArrayFilters{Filters: arrayFilters},
// 		),
// 	)

// 	if res.Err() != nil {
// 		log.Error("BatchUpdateTasks", zap.Int64("accountId", accountId), zap.Error(res.Err()))
// 		return res.Err()
// 	}

// 	return nil
// }

// func (m *LuckSaleMongoModel) AddNewTask(accountId uint64, task *LuckSaleTaskUnit) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	filter := bson.M{"_id": accountId}
// 	updateBson := bson.D{{"$addToSet", bson.D{{"task", task}}}}
// 	if _, err := m.GetCol().UpdateOne(ctx, filter, updateBson); err != nil {
// 		log.Error("AddTask err", zap.Uint64("accountId", accountId), zap.Error(err))
// 		return err
// 	}

// 	return nil
// }
