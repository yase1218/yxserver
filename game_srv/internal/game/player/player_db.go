package player

import (
	"context"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"gameserver/internal/game/model"
	"kernel/dao"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

const (
	dirty_weight_max = 100
)

// TODO: 优化 一个更新间隔内一个字段只更新一次(最新值)

func LoadUser(ctx context.Context, uid uint64) (*model.UserData, error) {
	user_data := &model.UserData{}
	db_res := db.LocalMongoReader.FindOne(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_User,
		bson.M{"userid": uid},
	)
	err := db_res.Decode(user_data)
	if err != nil {
		return nil, err
	}

	return user_data, nil
}

func LoadByNickname(ctx context.Context, nickname string) (*model.UserData, error) {
	user_data := &model.UserData{}
	db_res := db.LocalMongoReader.FindOne(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_User,
		bson.M{"nick": nickname},
	)
	err := db_res.Decode(user_data)
	if err != nil {
		return nil, err
	}

	return user_data, nil
}

// func CreateAccount(account *model.UserAccount) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	db_name := config.Conf.LocalMongo.DB
// 	collection := db.CollectionName_UserAccount
// 	op := &dao.WriteOperation{
// 		Database:   db_name,
// 		Collection: collection,
// 		Type:       dao.Upsert,
// 		Filter:     bson.M{"userid": account.UserId},
// 		Document:   account,
// 		Uuid:       account.UserId,
// 		Tms:        time.Now().Unix(),
// 		Cover:      true,
// 	}
// 	_, _, _, err := db.LocalMongoWriter.SyncWrite(ctx, op)
// 	return err
// }

// 异步写入所有脏数据
func (p *Player) saveDirtyAsync() {
	dirtyMap := make(map[string]interface{}, len(p.dirty))
	for k := range p.dirty {
		dirtyMap[k] = p.get_filed_by_name(k)
	}
	p.asyncFiledsToDb(dirtyMap)
	p.clear_dirty()
}

// 异步写入数据
func (p *Player) asyncFiledsToDb(values map[string]interface{}) {
	if len(values) == 0 {
		return
	}
	db_name := config.Conf.LocalMongo.DB
	collection := db.CollectionName_User
	op := &dao.WriteOperation{
		Database:   db_name,
		Collection: collection,
		Type:       dao.Upsert,
		Filter:     bson.M{"userid": p.GetUserId()},
		Uuid:       p.GetUserId(),
		Tms:        time.Now().UnixMilli(),
	}
	copy_values := make(map[string]interface{})
	for k, v := range values {
		copy_values[k] = dao.DeepCopy(v)
	}
	op.Update = bson.M{"$set": copy_values}

	db.LocalMongoWriter.AsyncWrite(op)
}

// 同步写入所有脏数据 仅在下线(或关服)时调用
func (p *Player) SaveDirtySync() {
	dirtyMap := make(map[string]interface{}, len(p.dirty))
	for k := range p.dirty {
		dirtyMap[k] = p.get_filed_by_name(k)
	}
	p.syncFiledsToDb(dirtyMap)
	p.clear_dirty()
}

// 同步写入数据
func (p *Player) syncFiledsToDb(values map[string]interface{}) {
	if len(values) == 0 {
		return
	}
	opt := dao.Upsert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db_name := config.Conf.LocalMongo.DB
	collection := db.CollectionName_User
	op := &dao.WriteOperation{
		Database:   db_name,
		Collection: collection,
		Type:       opt,
		Filter:     bson.M{"userid": p.GetUserId()},
		Document:   p.UserData,
		Uuid:       p.GetUserId(),
		Tms:        time.Now().UnixMilli(),
		Cover:      true,
	}

	_, _, _, err := db.LocalMongoWriter.SyncWrite(ctx, op)
	if err != nil {
		log.Error("sync_all_to_db error",
			zap.Uint64("uid", p.GetUserId()),
			zap.Uint64("userId", p.GetUserId()),
			zap.Error(err))
	}
}

// 同步写入所有数据(全量)  创角使用
func (p *Player) InsertUserToDb() {
	p.syncAllToDb()
}

// 同步全量更新
func (p *Player) syncAllToDb() {
	opt := dao.Insert
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db_name := config.Conf.LocalMongo.DB
	collection := db.CollectionName_User

	op := &dao.WriteOperation{
		Database:   db_name,
		Collection: collection,
		Type:       opt,
		Filter:     bson.M{"userid": p.GetUserId()},
		Document:   p.UserData,
		Uuid:       p.GetUserId(),
		Tms:        time.Now().UnixMilli(),
		Cover:      true,
	}

	_, _, _, err := db.LocalMongoWriter.SyncWrite(ctx, op)
	if err != nil {
		log.Error("sync_all_to_db error",
			zap.String("accountId", p.GetOpenId()),
			zap.Uint64("userId", p.GetUserId()),
			zap.Error(err))
	}
}

func (p *Player) get_filed_by_name(filedName string) interface{} {
	switch filedName {
	case "nick":
		return p.UserData.Nick
	case "level":
		return p.UserData.Level
	case "headimage":
		return p.UserData.HeadImg
	case "headframe":
		return p.UserData.HeadFrame
	case "title":
		return p.UserData.Title
	case "baseinfo":
		return p.UserData.BaseInfo
	case "stageinfo":
		return p.UserData.StageInfo
	case "items":
		return p.UserData.Items
	case "task":
		return p.UserData.Task
	case "mission":
		return p.UserData.Mission
	case "ships":
		return p.UserData.Ships
	case "team":
		return p.UserData.Team
	case "equip":
		return p.UserData.Equip
	case "shop":
		return p.UserData.Shop
	case "playmethod":
		return p.UserData.PlayMethod
	case "weapon":
		return p.UserData.Weapon
	case "treasure":
		return p.UserData.Treasure
	case "poker":
		return p.UserData.Poker
	case "accountactivity":
		return p.UserData.AccountActivity
	case "cardpool":
		return p.UserData.CardPool
	case "appearance":
		return p.UserData.Appearance
	case "petdata":
		return p.UserData.PetData
	case "frienddata":
		return p.UserData.FriendData
	case "fight":
		return p.UserData.Fight
	case "peakfight":
		return p.UserData.PeakFight
	case "contract":
		return p.UserData.Contract
	case "desert":
		return p.UserData.Desert
	case "arena":
		return p.UserData.Arena
	case "atlas":
		return p.UserData.Atlas
	case "lucksale":
		return p.UserData.LuckSale
	case "functionpreview":
		return p.UserData.FunctionPreview
	case "mail":
		return p.UserData.MailData
	case "isregister":
		return p.UserData.IsRegister
	case "equipstage":
		return p.UserData.EquipStage
	case "resources_pass":
		return p.UserData.ResourcesPass
	case "likes":
		return p.UserData.Likes
	case "weekpass":
		return p.UserData.WeekPass
	case "personalized":
		return p.UserData.Personalized
	}
	return nil
}

func (p *Player) set_dirty_async(filed string) {
	p.dirty[filed] = struct{}{}
	if len(p.dirty) > dirty_weight_max {
		p.saveDirtyAsync()
	}
}

func (p *Player) is_dirty() bool {
	return len(p.dirty) > 0
}

func (p *Player) clear_dirty() {
	p.dirty = make(map[string]struct{})
}

// // TODO:乐观锁版本号
// func (p *Player) async_filed_to_db(filed_name string, value interface{}, dc bool) {
// 	db_name := config.Conf.LocalMongo.DB
// 	collection := db.CollectionName_User
// 	op := &dao.WriteOperation{
// 		Database:   db_name,
// 		Collection: collection,
// 		Type:       dao.Upsert,
// 		Filter:     bson.M{"userid": p.GetUserId()},
// 		//Update:     bson.M{"$set": bson.M{filed_name: copy_value}},
// 		Uuid: p.GetUserId(),
// 		Tms:  time.Now().Unix(),
// 	}
// 	if dc {
// 		copy_value := dao.DeepCopy(value) // use deep copy
// 		op.Update = bson.M{"$set": bson.M{filed_name: copy_value}}
// 	} else {
// 		op.Update = bson.M{"$set": bson.M{filed_name: value}}
// 	}

// 	db.LocalMongoWriter.AsyncWrite(op)
// }

// // TODO:乐观锁版本号
// func (p *Player) async_all_to_db() {
// 	db_name := config.Conf.LocalMongo.DB
// 	collection := db.CollectionName_User
// 	copy_value := dao.DeepCopy(p.UserData)
// 	op := &dao.WriteOperation{
// 		Database:   db_name,
// 		Collection: collection,
// 		Type:       dao.Upsert,
// 		Filter:     bson.M{"userid": p.GetUserId()},
// 		Update:     bson.M{"$set": copy_value},
// 		Uuid:       p.GetUserId(),
// 		Tms:        time.Now().Unix(),
// 	}
// 	db.LocalMongoWriter.AsyncWrite(op)
// }
