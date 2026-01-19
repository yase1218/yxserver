package handle

// import (
// 	"github.com/v587-zyf/gc/log"
// 	"github.com/zy/game_data/template"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.uber.org/zap"
// 	"kernel/tools"
// 	"msg"
// 	"server/internal/config"
// 	"server/internal/game/common"
// 	"server/internal/game/model"
// 	"server/internal/game/service"
// 	"time"
// )

// // AllianceRankReq 请求联盟排行榜数据
// func AllianceRankReq(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.AllianceRankReq)
// 	//err, data, me := service.ServMgr.GetAllianceRankService().GetRank(playerData, req)
// 	err, data, me := service.ServMgr.GetRankService().GetAllianceRank(playerData, req)
// 	ack := &msg.AllianceRankAck{
// 		Result:   err,
// 		RankType: req.GetRankType(),
// 		RankList: data,
// 		Me:       me,
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, ack.Result)
// }

// // 玩家退出、被踢出联盟，从联盟排行榜删除
// func DelAllianceRank(allianceId uint32, accountId int64) {
// 	{
// 		//if err := rdb.DelAllianceRankByType(allianceId, accountId, msg.AllianceRankType_Alliance_Rank_Single_Power); err != nil {
// 		//	log.Error("del alliance rank single power err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		//}
// 		singlePowerFilter := bson.M{
// 			"alliance_id": allianceId,
// 			"account_id":  accountId,
// 			"server_id":   config.Conf.ServerId,
// 			"type":        msg.AllianceRankType_Alliance_Rank_Single_Power,
// 		}
// 		if err := model.GetAllianceRankModel().Delete(singlePowerFilter); err != nil {
// 			log.Error("del alliance rank single power err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		}

// 		date := tools.GetYearWeekByOffset(time.Now(), int(template.GetSystemItemTemplate().RefreshHour))

// 		//if err := rdb.DelAllianceRankByType(allianceId, accountId, msg.AllianceRankType_Alliance_Rank_Single_Active); err != nil {
// 		//	log.Error("del alliance rank single active err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		//}
// 		singleActiveFilter := bson.M{
// 			"alliance_id": allianceId,
// 			"account_id":  accountId,
// 			"server_id":   config.Conf.ServerId,
// 			"type":        msg.AllianceRankType_Alliance_Rank_Single_Active,
// 			"date":        date,
// 		}
// 		if err := model.GetAllianceRankModel().Delete(singleActiveFilter); err != nil {
// 			log.Error("del alliance rank single active err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		}

// 		//if err := rdb.DelAllianceRankByType(allianceId, accountId, msg.AllianceRankType_Alliance_Rank_Boss_Single); err != nil {
// 		//	log.Error("del alliance rank boss single err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		//}
// 		bossSingleFilter := bson.M{
// 			"alliance_id": allianceId,
// 			"account_id":  accountId,
// 			"server_id":   config.Conf.ServerId,
// 			"type":        msg.AllianceRankType_Alliance_Rank_Boss_Single,
// 			"date":        date,
// 		}
// 		if err := model.GetAllianceRankModel().Delete(bossSingleFilter); err != nil {
// 			log.Error("del alliance rank boss single err", zap.Error(err), zap.Int64("accountId", accountId), zap.Uint32("allianceId", allianceId))
// 		}

// 	}
// }

// // 玩家新进入联盟，上战力榜
// func AddAllianceRank(allianceId uint32, accountId int64, combat uint32) {
// 	if err := service.ServMgr.GetRankService().AddAllianceRank(allianceId, msg.AllianceRankType_Alliance_Rank_Single_Power, accountId, uint64(combat)); err != nil {
// 		log.Error("add alliance rank single power err", zap.Error(err),
// 			zap.Int64("accountId", accountId),
// 			zap.Uint32("allianceId", allianceId))
// 	}

// 	//if err := rdb.AddAllianceRankByType(allianceId, msg.AllianceRankType_Alliance_Rank_Single_Power, accountId, float64(combat)); err != nil {
// 	//	log.Error("add alliance rank single power err", zap.Error(err), zap.Int64("accountId", accountId))
// 	//}
// }
