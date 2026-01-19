package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"msg"
)

// RequestLoadCardPoolHandle 请求加载卡池
func RequestLoadCardPoolHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseLoadCardPool{
		Result: msg.ErrCode_SUCC,
	}

	res.Data = service.ToProtocolCardPools(p.UserData.CardPool.CardPools)
	p.SendResponse(packetId, res, res.Result)
}

// RequestLotteryHandle 请求抽奖
func RequestLotteryHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLottery)
	log.Debug("RequestLotteryHandle", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req))
	err, items := service.Lottery(p, req.CardId, req.LotteryTimes)
	res := &msg.ResponseLottery{
		Result: err,
	}
	if err == msg.ErrCode_SUCC {
		res.CardId = req.CardId
		res.LotteryTimes = req.LotteryTimes
		res.RewardItems = service.TemplateItemToProtocolItems(items)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestLotteryTimesRewardHandle 领取抽奖次数奖励
func RequestLotteryTimesRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLotteryTimesReward)
	err, items := service.GetLotteryTimesReward(p, req.CardId)
	res := &msg.ResponseLotteryTimesReward{
		Result: err,
		CardId: req.CardId,
	}
	if err == msg.ErrCode_SUCC {
		res.RewardItems = service.ToProtocolSimpleItems(items)
	}
	p.SendResponse(packetId, res, res.Result)
}
