package handle

// import (
// 	"github.com/v587-zyf/gc/log"
// 	"go.uber.org/zap"
// 	"kernel/tools"
// 	"msg"
// 	"server/internal/game/model"

// 	"server/internal/game/common"
// 	"server/internal/game/dao"
// )

// // RequenstEnterUnionBattleHandle 进入联盟战
// func RequenstEnterUnionBattleHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	retMsg := &msg.ResponseEnterUnionBattle{
// 		Result:  msg.ErrCode_SUCC,
// 		HurtMax: playerData.AccountInfo.AllianceBossHurtMax,
// 	}

// 	if member, err := dao.GetMember(playerData.AccountInfo.AccountId); err == nil {
// 		//retMsg.Rank, _ = rdb.GetUserRankingByTypes(uint32(member.AllianceID), int64(playerData.GetAccountId()), msg.AllianceRankType_Alliance_Rank_Boss_Single)

// 		ranking, err := model.GetAllianceRankModel().GetUserRanking(playerData.GetAccountId(), member.AllianceID, msg.AllianceRankType_Alliance_Rank_Boss_Single)
// 		if err != nil {
// 			log.Error("get alliance rank user ranking err", zap.Error(err),
// 				zap.Int64("accountId", playerData.GetAccountId()), zap.Uint32("allianceId", uint32(member.AllianceID)))
// 		} else {
// 			retMsg.Rank = ranking
// 		}
// 	}

// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }
