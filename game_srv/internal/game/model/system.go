package model

import (
	"time"

	"gameserver/internal/enum"
)

type System struct {
	Key     enum.SystemKey `bson:"_id"`
	Str     string         `bson:"str,omitempty"`
	Int64   int64          `bson:"int64,omitempty"`
	Float64 float64        `bson:"float64,omitempty"`
	Time    time.Time      `bson:"time,omitempty"`
}

type SystemMongoModel struct{}

var (
	SystemModel = &SystemMongoModel{}
)

func GetSystemModel() *SystemMongoModel {
	return SystemModel
}

// func GetSystemCol() *mongo.Collection {
// 	return db.GetGlobalClient().Database(config.Conf.LocalMongo.DB).Collection(enum.DB_LOCAL_COL_SYSTEM)
// }

// func (m *SystemMongoModel) UpdateOne(data *System) (*mongo.UpdateResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), publicconst.DB_OP_TIME_OUT)
// 	defer cancel()

// 	filter := bson.M{"_id": data.Key}
// 	update := bson.M{
// 		"$set": data,
// 	}
// 	return GetSystemCol().UpdateOne(ctx, filter, update)
// }

// func (m *SystemMongoModel) InsertOne(data *System) (*mongo.InsertOneResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), publicconst.DB_OP_TIME_OUT)
// 	defer cancel()

// 	return GetSystemCol().InsertOne(ctx, data)
// }

// func (m *SystemMongoModel) GetOneByKey(key enum.SystemKey) (*System, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), publicconst.DB_OP_TIME_OUT)
// 	defer cancel()

// 	var data *System
// 	filter := bson.M{"_id": key}

// 	ret := GetSystemCol().FindOne(ctx, filter)
// 	if err := ret.Decode(&data); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
// 		log.Error("get db system err", zap.Error(err))
// 		return nil, err
// 	}

// 	return data, nil
// }

// func (m *SystemMongoModel) GetAll(filter any) ([]*System, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), publicconst.DB_OP_TIME_OUT)
// 	defer cancel()

// 	var dataSlice []*System
// 	var err error
// 	cur, err := GetSystemCol().Find(ctx, filter)
// 	if err != nil {
// 		log.Error("find all db system err", zap.Error(err))
// 		return nil, err
// 	}

// 	for {
// 		if !cur.Next(ctx) {
// 			break
// 		}
// 		data := new(System)
// 		if err = cur.Decode(&data); err == nil {
// 			dataSlice = append(dataSlice, data)
// 		}
// 	}

// 	return dataSlice, err
// }
