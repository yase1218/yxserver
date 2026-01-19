package model

// import (
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"gopkg.in/mgo.v2/bson"

// 	"server/internal/config"
// 	"server/internal/db"
// 	"server/internal/publicconst"
// )

type (
	ArenaRankData struct {
		AccountId  int64   `bson:"_id"`
		ServerId   uint32  `bson:"server_id"`
		RankId     int32   `bson:"rank_id"`
		MonsterIds []int32 `bson:"monster_ids"`
		IsRobot    bool
	}

	ArenaEnemyData struct {
		AccountId int64  // 账号
		Nick      string // 昵称
		HeadImg   uint32 // 头像
		HeadFrame uint32 // 头像框
		Title     uint32 // 称号
		ShipId    uint32 // 机甲id
		Combat    uint32
	}
	ArenaPlayerPkRecordData struct {
		IsWin     bool            `bson:"isWin"`
		EnemyInfo *ArenaEnemyData `bson:"enemyinfo"`
		IsAttack  bool            `bson:"isattack"`
		OldRank   int32           `bson:"oldrank"`
		NewRank   int32           `bson:"newrank"`
		Stamp     int64           `bson:"stamp"`
	}

	ArenaPlayerData struct {
		AccountId        int64                      `bson:"_id"`
		ServerId         uint32                     `bson:"serverid"`
		TodayPkCnt       uint32                     `bson:"todayPkCnt"`
		TotalPkCnt       uint32                     `bson:"totalPkCnt"`
		TodayBuyPkCnt    uint32                     `bson:"todayBuyPkCnt"`
		TodayUseShips    []uint32                   `bson:"todayUseShips"`
		RewardBeginStamp int64                      `bson:"rewardBeginStamp"`
		LastResetStamp   int64                      `bson:"lastResetStamp"`
		UnlockMonsterIds []uint32                   `bson:"unlockMonsterIds"`
		DefendMonster    []int32                    `bson:"defendMonster"`
		Records          []*ArenaPlayerPkRecordData `bson:"records"`
	}

	// ArenaMongoModel struct{}
)

// var (
// 	ArenaModel = new(ArenaMongoModel)
// )

// func GetArenaModel() *ArenaMongoModel {
// 	return ArenaModel
// }

// func (m *ArenaMongoModel) GetDB() *mongo.Database {
// 	return db.GetLocalClient().Database(config.Conf.GetLocalDB())
// }

// func (m *ArenaMongoModel) GetCol() *mongo.Collection {
// 	return m.GetDB().Collection(publicconst.LOCAL_ARENA)
// }

// func (m *ArenaMongoModel) Create(data *ArenaPlayerData) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().InsertOne(ctx, data)
// 	return err
// }

// func (m *ArenaMongoModel) Update(data *ArenaPlayerData) error {
// 	ctx, cancel := GetDBCtx()
// 	defer cancel()

// 	_, err := m.GetCol().UpdateOne(ctx, bson.M{"_id": data.AccountId}, data)
// 	return err
// }
