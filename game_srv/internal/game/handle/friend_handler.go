package handle

import (
	"gameserver/internal/content"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
	"strconv"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

// import (
// 	"gameserver/internal/game/player"
// 	"gameserver/internal/game/service"
// 	"gameserver/internal/publicconst"
// 	"msg"

// 	"github.com/zy/game_data/template"
// )

// RequestFriListHandle 请求好友列表
func RequestFriListHandle(pid uint32, args interface{}, p *player.Player) {
	service.GetFriendList(pid, p)
}

func RequestFriApplyListHandle(pid uint32, args interface{}, p *player.Player) {
	service.GetFriApplyList(pid, p)
}

func RequestAddFriendHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestAddFriend)
	service.AddFriend(pid, p, req)
}

func RequestFriApplyOpHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestFriApplyOp)
	service.FriendApplyOp(pid, p, req)
}

func RequestDelFriendHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestDelFriend)
	service.DelFriend(pid, p, req)
}

func RequestRecommandFriendHandle(pid uint32, args interface{}, p *player.Player) {
	service.RecommandPlayer(pid, p)
}

func RequestSearchPlayerHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSearchPlayer)
	service.SearchPlayer(pid, p, req)
}

func RequestFriBlackListHandle(pid uint32, args interface{}, p *player.Player) {
	service.GetBlackList(pid, p)
}

func RequestBlackListOpHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestBlackListOp)
	service.BlackOp(pid, p, req)
}

func RequestPrivateChatHandle(pid uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPrivateChat)
	res := &msg.ResponsePrivateChat{
		Result:  msg.ErrCode_SUCC,
		Content: req.Content,
	}
	if code := service.PrivateChatCheck(p, req); code != msg.ErrCode_SUCC {
		res.Result = code
		p.SendResponse(pid, res, res.Result)
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
				res.Result = ec
				if ec == msg.ErrCode_SUCC {
					service.PrivateChat(pid, chatPlayer, req)
				} else {
					chatPlayer.SendResponse(pid, res, res.Result)
				}
			}
		},
	}); err != nil {
		log.Error("RequestPrivateChatHandle failed", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		p.SendResponse(pid, res, res.Result)
	}
}

// // RequestFriApplyListHandle 好友申请列表
// func RequestFriApplyListHandle(packetId uint32, args interface{}, p *player.Player) {
// 	retMsg := &msg.ResponseFriApplyList{
// 		Result: msg.ErrCode_SUCC,
// 	}

// 	if err := service.FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}

// 	retMsg.Data = service.ToPlayerSimpleInfos(service.GetFriApplyList(p))
// 	p.SendMsg(packetId, retMsg)
// }

// // RequestAddFriendHandle 添加好友
// func RequestAddFriendHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestAddFriend)
// 	retMsg := &msg.ResponseAddFriend{}
// 	if err := service.FunctionOpen(p, publicconst.Add_Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}

// 	retMsg.Result = service.AddFriend(p, req.AccountId)
// 	p.SendMsg(packetId, retMsg)
// }

// // RequestFriApplyOpHandle 好友申请操作
// func RequestFriApplyOpHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestFriApplyOp)
// 	retMsg := &msg.ResponseFriApplyOp{
// 		Result:    msg.ErrCode_SUCC,
// 		AccountId: req.AccountId,
// 		Op:        req.Op,
// 	}

// 	if err := service.FunctionOpen(p, publicconst.Add_Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}

// 	retMsg.Result = service.FriendApplyOp(p, req.AccountId, req.Op)
// 	p.SendMsg(packetId, retMsg)
// }

// // RequestDelFriendHandle 删除好友
// func RequestDelFriendHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestDelFriend)
// 	retMsg := &msg.ResponseDelFriend{
// 		Result:    service.DelFriend(p, req.AccountId),
// 		AccountId: req.AccountId,
// 	}
// 	p.SendMsg(packetId, retMsg)
// }

// func RequestRecommandFriendHandle(packetId uint32, args interface{}, p *player.Player) {
// 	retMsg := &msg.ResponseRecommandFriend{
// 		Result: msg.ErrCode_SUCC,
// 	}
// 	if err := service.FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}
// 	retMsg.Data = service.ToPlayerSimpleInfos(service.RecommandPlayer(p))
// 	p.SendMsg(packetId, retMsg)
// }

// func RequestSearchPlayerHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestSearchPlayer)
// 	retMsg := &msg.ResponseSearchPlayer{}
// 	if err := service.FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}

// 	data := service.SearchPlayer(p, req.Content)
// 	if data == nil {
// 		retMsg.Result = msg.ErrCode_PLAYER_NOT_EXIST
// 	} else {
// 		retMsg.Data = service.ToPlayerSimpleInfo(data)
// 	}
// 	p.SendMsg(packetId, retMsg)
// }

// func RequestFriBlackListHandle(packetId uint32, args interface{}, p *player.Player) {
// 	retMsg := &msg.ResponseFriBlackList{
// 		Result: msg.ErrCode_SUCC,
// 		Data:   service.ToPlayerSimpleInfos(service.GetBlackList(p)),
// 	}
// 	p.SendMsg(packetId, retMsg)
// }

// func RequestBlackListOpHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestBlackListOp)
// 	retMsg := &msg.ResponseBlackListOp{
// 		Result: msg.ErrCode_SUCC,
// 		Op:     req.Op,
// 	}
// 	retMsg.Result = service.BlackOp(p, req.AccountId, req.Op)
// 	if retMsg.Result == msg.ErrCode_SUCC {
// 		_, players := service.ServMgr.GetSocialService().GetPlayerSimpleInfos(req.AccountId)
// 		for _, player := range players {
// 			retMsg.Data = append(retMsg.Data, service.ToPlayerSimpleInfo(player))
// 		}
// 	}
// 	p.SendMsg(packetId, retMsg)
// }

// func RequestPrivateChatHandle(packetId uint32, args interface{}, p *player.Player) {
// 	req := args.(*msg.RequestPrivateChat)
// 	retMsg := &msg.ResponsePrivateChat{}
// 	if err := service.FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
// 		retMsg.Result = err
// 		p.SendMsg(packetId, retMsg)
// 		return
// 	}

// 	err, target := service.PrivateChat(p, req.AccountId, req.Content, req.Para)
// 	retMsg.Result = err
// 	retMsg.Content = template.GetForbiddenTemplate().Filter(req.Content)
// 	retMsg.Para = req.Para
// 	if retMsg.Result == msg.ErrCode_SUCC {
// 		retMsg.TargetPlayer = service.ToPlayerSimpleInfo(target)
// 	}
// 	p.SendMsg(packetId, retMsg)
// }
