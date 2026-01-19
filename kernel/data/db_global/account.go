package db_global

import (
	"context"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	mgoptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type (
	AccMsdk struct {
		Platform string `bson:"platform"`
		OpenId   string `bson:"open_id"`
	}
	AccNormal struct {
		Account string `bson:"account"` // 账号
	}
	Accounts struct {
		Normal *AccNormal `bson:"normal,omitempty"`
		Msdk   *AccMsdk   `bson:"msdk,omitempty"`
	}
	AccountChannelInfo struct {
		Channel     string
		AccountInfo any
	}

	//Account struct {
	//	ID        uint64    `bson:"_id"`
	//	Account   *Accounts `bson:"accounts"` // 账号系统
	//	SerId     uint64    `bson:"ser_id"`   // 服务器id
	//	Ip        string    `bson:"ip"`       // ip
	//	Nick      string    `bson:"nick"`     // 昵称
	//	ChannelId uint32    `bson:"channelId"`
	//
	//	DeviceId       string `bson:"device_id"`       // 设备id
	//	Os             string `bson:"os"`              // 操作系统
	//	OsVersion      string `bson:"os_version"`      // 操作系统版本
	//	AppVersion     string `bson:"app_version"`     // app版本号
	//	Manufacturer   string `bson:"manufacturer"`    // 设备制造商
	//	DeviceModel    string `bson:"device_model"`    // 手机型号
	//	ScreenHeight   string `bson:"screen_height"`   // 屏幕高度
	//	ScreenWidth    string `bson:"screen_width"`    // 屏幕宽度
	//	Ram            string `bson:"ram"`             // 设备运行内存状态
	//	Disk           string `bson:"disk"`            // 设备存储空间状态
	//	NetworkType    string `bson:"network_type"`    // 网络状态
	//	Carrier        string `bson:"carrier"`         // 网络运营商
	//	Country        string `bson:"country"`         // 国家地区
	//	CountryCode    string `bson:"country_code"`    // 国家地区代码
	//	SystemLanguage string `bson:"system_language"` // 系统语言
	//
	//	//LoginContinue uint32    `bson:"login_continue"` // 连续登录天数
	//	LoginCnt    uint32    `bson:"login_cnt"`     // 累计登录天数
	//	CreateAt    time.Time `bson:"create_at"`     // 创建时间
	//	UpdateAt    time.Time `bson:"update_at"`     // 最后更新时间
	//	LastLoginAt time.Time `bson:"last_login_at"` // 最后登录时间
	//	LastOutAt   time.Time `bson:"last_out_at"`   // 最后登出时间
	//}

	Account struct {
		ID       uint64    `bson:"_id"` // userid
		OpenId   string    `bson:"open_id"`
		PlatForm uint32    `bson:"platform"`
		ServerId uint32    `bson:"server_id"` // 服务器id
		CreateAt time.Time `bson:"create_at"` // 创建时间
		UpdateAt time.Time `bson:"update_at"` // 最后更新时间
	}
)

func NewAccNormal(acc string) *Accounts {
	return &Accounts{
		Normal: &AccNormal{Account: acc},
	}
}

func NewAccMsdk(openId, platform string) *Accounts {
	return &Accounts{
		Msdk: &AccMsdk{OpenId: openId, Platform: platform},
	}
}

// func MakeAccountFilter(accountInfo *AccountChannelInfo) bson.M {
// 	filter := bson.M{}
// 	switch accountInfo.Channel {
// 	case kenum.Account_Type_Normal:
// 		accInfo := accountInfo.AccountInfo.(*AccNormal)
// 		filter = bson.M{"accounts.normal.account": accInfo.Account}
// 	case kenum.Account_Type_Msdk:
// 		accInfo := accountInfo.AccountInfo.(*AccMsdk)
// 		filter = bson.M{"accounts.msdk.open_id": accInfo.OpenId, "accounts.msdk.platform": accInfo.Platform}
// 	}
// 	return filter
// }

type AccountDBModel struct{}

var (
	AccountModel = &AccountDBModel{}
)

func AccountCreateIndex() {
	collecttion := GetAccountCol()

	err := collecttion.CreateOneIndex(context.Background(), options.IndexModel{
		Key: []string{
			"open_id",
			"platform",
			"server_id",
		},
		IndexOptions: mgoptions.Index().SetUnique(true),
	})
	if err != nil {
		log.Panic("create account index err", zap.Error(err))
	}

	err = collecttion.CreateOneIndex(context.Background(), options.IndexModel{
		Key: []string{
			"open_id",
		},
	})
	if err != nil {
		log.Panic("create account index err", zap.Error(err))
	}
}

func GetAccountModel() *AccountDBModel {
	return AccountModel
}

func GetAccountCol() *qmgo.Collection {
	return mongo.DB(GetDB()).Collection(COL_ACCOUNT)
}

func (m *AccountDBModel) Upsert(data *Account) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	return GetAccountCol().Upsert(context.Background(), filter, data)
}

func (m *AccountDBModel) GetAllByOpenId(openId string) ([]*Account, error) {
	var data []*Account
	var err error
	filter := bson.M{"open_id": openId}
	err = GetAccountCol().Find(context.Background(), filter).All(&data)
	return data, err
}

func (m *AccountDBModel) GetOneByKey(openId string, platform, serverId uint32) (*Account, error) {
	var data *Account
	var err error
	filter := bson.M{"open_id": openId, "platform": platform, "server_id": serverId}
	err = GetAccountCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *AccountDBModel) GetOne(userID uint64) (*Account, error) {
	var data *Account
	var err error
	filter := bson.M{"_id": userID}
	err = GetAccountCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *AccountDBModel) GetAll(filter any) ([]*Account, error) {
	var dataSlice []*Account
	var err error
	err = GetAccountCol().Find(context.Background(), filter).All(&dataSlice)
	return dataSlice, err
}

// func (m *AccountDBModel) NewUserUnique(accountInfo *AccountChannelInfo, data *Account) error {
// 	filter := MakeAccountFilter(accountInfo)

// 	count, err := GetAccountCol().Find(context.Background(), filter).Count()
// 	if err != nil {
// 		log.Error("user get account err", zap.Reflect("filter", filter), zap.String("err", err.Error()))
// 		return err
// 	}

// 	if count != 0 {
// 		return fmt.Errorf("accountRegisterAlready")
// 	}

// 	if _, err = GetAccountCol().InsertOne(context.Background(), data); err != nil {
// 		if qmgo.IsDup(err) {
// 			log.Error("account insert duplicate err", zap.Uint64("id", data.ID), zap.Error(err))
// 		} else {
// 			log.Error("account get account err", zap.Reflect("filter", filter), zap.Error(err))
// 		}

// 		return err
// 	}

// 	return nil
// }
