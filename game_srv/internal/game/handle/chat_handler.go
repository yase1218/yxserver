package handle

import (
	"msg"
	"strconv"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"gameserver/internal/content"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RequestWorldChatHandle 请求世界聊天
func RequestWorldChatHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestWorldChat)
	res := &msg.ResponseWorldChat{
		Result: msg.ErrCode_SUCC,
		Tp:     req.Tp,
	}

	if ok, data := service.SendGmHandle(p, req.Content); ok {
		res.Data = service.ToProtocolWorldChatInfo(data)
		p.SendResponse(packetId, res, res.Result)
		return
	}

	if req.IsShare {
		ec, data := service.SendWorldChat(p, req.Content, req.Para)
		res.Result = ec
		res.Data = service.ToProtocolWorldChatInfo(data)
		log.Debug("RequestWorldChatHandle share", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res))
		p.SendResponse(packetId, res, res.Result)
		return
	}

	if code := service.SendWorldCheck(p, req.Content, req.Para); code != msg.ErrCode_SUCC {
		res.Result = code
		p.SendResponse(packetId, res, res.Result)
		return
	}

	checkContent := []string{req.Content}
	if len(req.Para) > 0 {
		checkContent = append(checkContent, req.Para)
	}
	userId := p.GetUserId()
	if err := content.PushContent(&content.ContentData{
		Content:     checkContent,
		UserId:      p.UserData.AccountId,
		ChannelNo:   p.SdkChannelNo,
		Os:          p.Os,
		RoleId:      strconv.Itoa(int(p.GetUserId())),
		RoleName:    p.UserData.Nick,
		RoleLevel:   strconv.Itoa(int(p.UserData.Level)),
		ServerId:    strconv.Itoa(int(p.UserData.ServerId)),
		Ip:          p.UserData.BaseInfo.Ip,
		ContentType: content.ContentTypeChat,
		Cb: func(ec msg.ErrCode) {
			chatPlayer := player.FindByUserId(userId)
			if chatPlayer != nil {
				data := &service.WorldMsg{}
				if ec == msg.ErrCode_SUCC {
					if req.Tp == msg.ChatChannelType_Chat_Channel_World {
						ec, data = service.SendWorldChat(chatPlayer, req.Content, req.Para)
					}
				}

				res.Result = ec
				if ec == msg.ErrCode_SUCC {
					res.Data = service.ToProtocolWorldChatInfo(data)
				}
				chatPlayer.SendResponse(packetId, res, res.Result)
			}

			// TODO 公会聊天
			// else if req.Tp == msg.ChatChannelType_Chat_Channel_Alliance {
			// 	ec, data = service.ServMgr.GetAllianceService().SendChat(p, req.Content, req.Para)
			// }

		},
	}); err != nil {
		log.Error("RequestWorldChatHandle failed", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		p.SendResponse(packetId, res, res.Result)
	}
}

// RequestEnterWorldChatHandle 进入世界聊天
func RequestEnterWorldChatHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestEnterWorldChat)

	res := &msg.ResponseEnterWorldChat{
		Result: msg.ErrCode_SUCC,
		Tp:     req.Tp,
	}

	var data []*service.WorldMsg
	if req.Tp == msg.ChatChannelType_Chat_Channel_World {
		data = service.EnterWorldChat(p)
	}
	// TODO 公会聊天
	// else if req.Tp == msg.ChatChannelType_Chat_Channel_Alliance {
	// 	data = service.ServMgr.GetAllianceService().EnterChat(p)
	// }

	res.Data = service.ToProtocolWorldChatInfos(data)

	p.SendResponse(packetId, res, res.Result)
}

// RequestLeaveWorldChatHandle 离开世界聊天
func RequestLeaveWorldChatHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestLeaveWorldChat)

	service.LeaveWorldChat(p)
	res := &msg.ResponseLeaveWorldChat{
		Result: msg.ErrCode_SUCC,
		Tp:     req.Tp,
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestWorldChatDataHandle 请求世界聊天数据
func RequestWorldChatDataHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestWorldChatData)
	res := &msg.ResponseWorldChatData{
		Result: msg.ErrCode_SUCC,
		Tp:     req.Tp,
	}

	if req.Tp == msg.ChatChannelType_Chat_Channel_World {
		res.Data = service.ToProtocolWorldChatInfos(service.GetWorldChatData(req.StartId))
	}
	// TODO 公会聊天
	// else if req.Tp == msg.ChatChannelType_Chat_Channel_Alliance {
	// 	retMsg.Data = service.ToProtocolWorldChatInfos(service.ServMgr.GetAllianceService().GetChatData(p, req.StartId))
	// }

	p.SendResponse(packetId, res, res.Result)
}

// 发送表情
func EmoteSend(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.EmoteSendReq)
	//ack := &msg.SelectAccessoryAck{Result: msg.ErrCode_SUCC}
	if err := service.EmoteSend(p, req); err != nil {
		log.Error("send emote err", zap.Error(err))
		//ack.Result = msg.ErrCode_SYSTEM_ERROR
	}
	//tools.SendMsg(playerData.PlayerAgent, ack, packetId, ack.Result)
}
