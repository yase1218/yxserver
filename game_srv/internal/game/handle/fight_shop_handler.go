package handle

// import (
// 	"github.com/v587-zyf/gc/log"
// 	"go.uber.org/zap"
// 	"msg"
// 	"gameserver/internal/game/common"
// 	"gameserver/internal/game/service"
// )

// func FightShop(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.FightShopReq)
// 	//ack := new(msg.MoveAck)
// 	if err := service.ServMgr.GetFightService().SendToFight(playerData, req); err != nil {
// 		log.Error("fight shop err", zap.Error(err))
// 	}
// 	//tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func FightShopLock(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.FightShopLockReq)
// 	//ack := new(msg.MoveAck)
// 	if err := service.ServMgr.GetFightService().SendToFight(playerData, req); err != nil {
// 		log.Error("fight shop lock err", zap.Error(err))
// 	}
// 	//tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func FightShopBuy(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.FightShopBuyReq)
// 	//ack := new(msg.MoveAck)
// 	if err := service.ServMgr.GetFightService().SendToFight(playerData, req); err != nil {
// 		log.Error("fight shop buy err", zap.Error(err))
// 	}
// 	//tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }

// func FightShopRefresh(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.FightShopRefreshReq)
// 	//ack := new(msg.MoveAck)
// 	if err := service.ServMgr.GetFightService().SendToFight(playerData, req); err != nil {
// 		log.Error("fight shop refresh err", zap.Error(err))
// 	}
// 	//tools.SendMsg(playerData.PlayerAgent, ack, packetId, msg.ErrCode_SUCC)
// }
