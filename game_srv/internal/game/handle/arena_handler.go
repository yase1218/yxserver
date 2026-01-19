package handle

// import (
// 	"kernel/tools"
// 	"msg"

// 	"server/internal/game/common"
// 	"server/internal/game/service"
// )

// func RequestArenaInfo(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := service.ServMgr.GetArenaService().GetPlayerInfo(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestRanks(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := &msg.ArenaRankAck{}
// 	ack.Ranks = service.ServMgr.GetArenaService().GetRanks(playerData)
// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestMonsterInfo(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := service.ServMgr.GetArenaService().GetPlayerMonsterInfo(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenUnlockPos(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.ArenaUnlockPosReq)
// 	ack := &msg.ArenaUnlockPosAck{}
// 	ack.Result = service.ServMgr.GetArenaService().UnlockPos(playerData, int(req.Pos))
// 	ack.Pos = req.Pos
// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenUnlockMonster(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.ArenaUnlockMonsterReq)

// 	ack := &msg.ArenaUnlockMonsterAck{}
// 	ack.Result = service.ServMgr.GetArenaService().UnlockMonster(playerData, req.MonsterId)
// 	ack.MonsterId = req.MonsterId
// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenBuyPkCnt(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := &msg.ArenaBuyPkCntAck{}
// 	ack.Result = service.ServMgr.GetArenaService().BuyPkCnt(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenGetRewardList(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := &msg.ArenaGetRewardAck{}
// 	ack.RewardBeginStamp, ack.RewardItem = service.ServMgr.GetArenaService().GetAFKReward(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenPkList(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	ack := &msg.ArenaPkListAck{}
// 	ack.PkList = service.ServMgr.GetArenaService().RefreshPkList(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenaPk(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.ArenaPkReq)

// 	ack := &msg.ArenaPkAck{}
// 	ack.Result = service.ServMgr.GetArenaService().Pk(playerData, req.TargetRank)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenaPkRecord(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	ack := &msg.ArenaRecordAck{}
// 	ack.Records = service.ServMgr.GetArenaService().GetAllRecord(playerData)

// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func RequestArenaSetDefend(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.ArenaSetDefendReq)

// 	ack := &msg.ArenaSetDefendAck{}
// 	ack.Result = service.ServMgr.GetArenaService().SetDefend(playerData, req.MonsterIds)
// 	if ack.Result == msg.ErrCode_SUCC {
// 		ack.MonsterIds = req.MonsterIds
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }
