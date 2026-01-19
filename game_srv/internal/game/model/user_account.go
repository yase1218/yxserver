package model

type UserAccount struct { // TODO 设置索引
	AccountId  string `bson:"accountid"`
	UserId     uint64 `bson:"userid"`
	CreateTime uint32
}
