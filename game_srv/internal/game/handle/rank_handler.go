package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"

	"github.com/zy/game_data/template"
)

func RequestRankData(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GetTargetRankInfoReq)
	data, selfInfo, nextRewardTime := service.GetRankData(p, template.RankType(req.Type))
	resp := &msg.GetTargetRankInfoResp{}
	resp.RankData = data
	resp.Type = req.Type
	resp.SelfInfo = selfInfo
	resp.NextRewardTime = uint64(nextRewardTime)
	p.SendResponse(packetId, resp, msg.ErrCode_SUCC)
}

func RequestAddLikesForRank(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.AddLikesForRankPlayerReq)
	code, likeNum := service.HandleAddLikesForRank(p, template.RankType(req.Type), req.AccountId)
	resp := &msg.AddLikesForRankPlayerResp{}
	resp.Code = code
	resp.AccountId = req.AccountId
	resp.Type = req.Type
	resp.LikeNum = uint32(likeNum)
	p.SendResponse(packetId, resp, resp.Code)
}

func RequestGetMaxFirstPassReward(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GetMaxFirstPassRankRewardReq)
	code, items := service.GetFirstPassMaxReward(p, template.RankType(req.Type))
	resp := &msg.GetMaxFirstPassRankRewardResp{}
	resp.Code = code
	resp.Items = items
	resp.Type = req.Type
	p.SendResponse(packetId, resp, resp.Code)
}

func RequestFirstPassRecordData(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GetFirstPassRecordDataReq)
	code, recordMap, selfInfo := service.GetFirstPassRecordData(p, template.RankType(req.Type))
	resp := &msg.GetFirstPassRecordDataResp{}
	resp.Code = code
	resp.TotalRecord = recordMap
	resp.SelfPassRecord = selfInfo
	resp.Type = req.Type
	p.SendResponse(packetId, resp, resp.Code)
}

func HandleUpdateRankInfo(p *player.Player, args interface{}, rankType template.RankType) {
	service.UpdateCommonRankInfo(p, args, rankType)
}

// // RequestLikesHandle 点赞排行榜
// func RequestLikesHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.RequestLikes)
// 	err, items, likes := service.ServMgr.GetRankService().LikesMissionRank(playerData, req.Tp)
// 	retMsg := &msg.ResponseLikes{
// 		Result: err,
// 	}
// 	if err == msg.ErrCode_SUCC {
// 		retMsg.GetItems = service.TemplateItemToProtocolItems(items)
// 		retMsg.Likes = likes
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

// func RequestRankRewardInfoHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.RequestRankRewardInfo)
// 	err, data := service.ServMgr.GetRankService().GetRankRewardInfo(playerData, req.Tp)
// 	retMsg := &msg.ResponseRankRewardInfo{
// 		Result: err,
// 		Tp:     req.Tp,
// 	}
// 	if err == msg.ErrCode_SUCC {
// 		retMsg.Data = service.ToProtocolRankMissionRewards(data)
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

// func RequestRankMissionRewardHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.RequestRankMissionReward)
// 	err, ids, items := service.ServMgr.GetRankService().GetRankMissionReward(playerData, req.Tp)
// 	retMsg := &msg.ResponseRankMissionReward{
// 		Result: err,
// 		Tp:     req.Tp,
// 	}
// 	if err == msg.ErrCode_SUCC {
// 		for _, v := range ids {
// 			retMsg.MissionIds = append(retMsg.MissionIds, uint32(v))
// 		}
// 		retMsg.GetItems = service.ToProtocolSimpleItems(items)
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }

func RequestSpecialMissionRankHandle(packetId uint32, args interface{}, playerData *player.Player) {
	// 	req := args.(*msg.RequestSpecialMissionRank)
	// 	data, passTime := service.ServMgr.GetRankService().GetSpecialMissionRank(playerData, int(req.MissionId))
	// 	retMsg := &msg.ResponseSpecialMissionRank{
	// 		Result:    msg.ErrCode_SUCC,
	// 		MissionId: req.MissionId,
	// 		Data:      service.ToProtocolSpecialMissionRank(data),
	// 		PassTime:  passTime,
	// 	}
	// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
	retMsg := &msg.ResponseSpecialMissionRank{
		Result: msg.ErrCode_SUCC,
	}
	playerData.SendResponse(packetId, retMsg, retMsg.Result)
}

// func RequestDesertRankHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	data := service.ServMgr.GetRankService().GetDesertRankData(playerData)
// 	retMsg := &msg.ResponseDesertRank{
// 		Result: msg.ErrCode_SUCC,
// 		Data:   service.ToProtocolDesertRank(data),
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, retMsg, packetId, retMsg.Result)
// }
