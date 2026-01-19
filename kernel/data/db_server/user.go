package db_server

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type (
	User struct {
		ID   uint64 `bson:"_id"`
		Nick string `bson:"nick"`
		Bag  *Bag   `bson:"bag"`

		CreateAt    time.Time `bson:"create_at"`     // 创建时间
		UpdateAt    time.Time `bson:"update_at"`     // 最后更新时间
		LastLoginAt time.Time `bson:"last_login_at"` // 最后登录时间
	}
)

func NewUser(id uint64, nick string) *User {
	return &User{
		ID:       id,
		Nick:     nick,
		Bag:      NewBag(),
		CreateAt: time.Now(),
	}
}

type UserDBModel struct{}

var (
	UserModel = &UserDBModel{}
)

func GetUserModel() *UserDBModel {
	return UserModel
}

func GetUserCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_USER)
}

func (m *UserDBModel) InsertOne(data *User) (*qmgo.InsertOneResult, error) {
	return GetUserCol().InsertOne(context.Background(), data)
}

func (m *UserDBModel) Upsert(data *User) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	return GetUserCol().Upsert(context.Background(), filter, data)
}

func (m *UserDBModel) GetOne(userID uint64) (*User, error) {
	var data *User
	var err error
	filter := bson.M{"_id": userID}
	err = GetUserCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *UserDBModel) GetAll(filter any) ([]*User, error) {
	var dataSlice []*User
	var err error
	err = GetUserCol().Find(context.Background(), filter).All(&dataSlice)
	return dataSlice, err
}
