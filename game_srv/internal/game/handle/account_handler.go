package handle

import (
	"gameserver/internal/content"
	"msg"
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"

	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RpcRegist 注册
func RpcRegist(args []interface{}) {

	//uidServer := service.ServMgr.GetInterService().GetUidServer()
	//if uidServer == 0 {
	//	res.Result = msg.ErrCode_SYSTEM_ERROR
	//	agent.WriteMsg(res)
	//	return
	//}

	// 测试消息
	//subject := fmt.Sprintf("%v", uidServer)
	//interMsg := &msg.RequestCommonInterMsg{
	//	ServerId: int64(conf.Server.ServerId),
	//	MsgId:    uint32(msg.MsgId_ID_RequestInterRegistUid),
	//}
	//registMsg := &msg.RequestInterRegistUid{
	//	AccountId: accountId,
	//}
	//temp, _ := proto.Marshal(registMsg)
	//interMsg.Data = temp
	//data, _ := proto.Marshal(interMsg)
	//service.ServMgr.GetNatsService().Publish(subject, data)
}

// RpcLoginRegist 登陆账号不存在则注册
func RpcLoginRegist(args []interface{}) {
}

// // RpcLogin 登录
// func RpcLogin(args []interface{}) {
// 	userId := args[0].(string)
// 	accountId := args[1].(int64)
// 	agent := args[2].(gate.Agent)
// 	packetId := args[3].(uint32)
// 	channelId := args[4].(uint32)
// 	nick := args[5].(string)
// 	extraInfo := args[6].(string)
// 	ip := args[7].(string)
// 	//tdaCommAtr := args[8].(string)

// 	// var tdaCommonAttr *tda.CommonAttr
// 	// if err := json.Unmarshal([]byte(tdaCommAtr), &tdaCommonAttr); err != nil {
// 	// 	log.Error("unmarshal err", zap.Error(err))
// 	// 	return
// 	// }

// 	// TODO 暂时屏蔽 需要客户端登录game时上报
// 	tdaCommonAttr := &tda.CommonAttr{AccountId: fmt.Sprintf("%v", accountId)}

// 	//if len(nick) == 0 {
// 	//	nick = userId
// 	//}
// 	//log.Infof("RpcLogin accountId:%v", accountId)

// 	if common.PlayerMgr.GetOnlineNum() >= config.Conf.MaxRegisterNum {
// 		resMsg := &msg.ResponseLogin{Result: msg.ErrCode_SYSTEM_ERROR}
// 		tools.SendMsg(agent, resMsg, packetId, resMsg.Result)
// 		agent.Close()
// 		return
// 	}

// 	account := dao.AccountDao.GetAccount(accountId)
// 	if account != nil && account.Forbidden {
// 		resMsg := &msg.ResponseLogin{Result: msg.ErrCode_FORBIDDEN_USER}
// 		tools.SendMsg(agent, resMsg, packetId, resMsg.Result)
// 		agent.Close()
// 		return
// 	}

// 	req := &DoLoginReq{
// 		UserId:        userId,
// 		AccountId:     accountId,
// 		Agent:         agent,
// 		PackId:        packetId,
// 		Account:       account,
// 		ExtraInfo:     extraInfo,
// 		Ip:            ip,
// 		ChanelId:      channelId,
// 		TdaCommonAttr: tdaCommonAttr,
// 	}
// 	if account == nil || account.UserId == "" {
// 		var err msg.ErrCode
// 		//log.Infof("AddAccount:%v", accountId)
// 		nick = fmt.Sprintf("kulu_%d", accountId)
// 		err, account = dao.AccountDao.AddAccount(userId, accountId, channelId, nick)

// 		// 注册成功的处理
// 		if err == msg.ErrCode_SUCC {
// 			// 添加dnu 统计
// 			dao.UserStaticDao.AddDNUStatics(model.NewDNUStatics(tools.GetCurDateStart(), channelId, uint32(accountId), tools.GetCurTime(), agent.RemoteAddr().String(), extraInfo, userId))
// 			req.Account = account
// 			//doLogin(userId, accountId, agent, packetId, account, extraInfo, ip)
// 			doLogin(req)

// 			// tda user update rolename

// 		} else {
// 			resMsg := &msg.ResponseLogin{Result: err}
// 			tools.SendMsg(agent, resMsg, packetId, resMsg.Result)
// 		}
// 	} else {
// 		//doLogin(userId, accountId, agent, packetId, account, extraInfo, ip)
// 		doLogin(req)
// 	}
// }

// type (
// 	DoLoginReq struct {
// 		UserId        string
// 		AccountId     int64
// 		Agent         gate.Agent
// 		PackId        uint32
// 		Account       *model.Account
// 		ExtraInfo     string
// 		Ip            string
// 		ChanelId      uint32
// 		TdaCommonAttr *tda.CommonAttr
// 	}
// )

// // func doLogin(userId string, accountId int64, agent gate.Agent, packetId uint32, account *model.Account, extraInfo, ip string) {
// func doLogin(req *DoLoginReq) {
// 	// 登录中
// 	if checkLogining(req.Agent) {
// 		res := &msg.ResponseLogin{
// 			Result: msg.ErrCode_ISLOGINING,
// 		}
// 		tools.SendMsg(req.Agent, res, req.PackId, res.Result)
// 		return
// 	}

// 	dao.AccountDao.UpdateExtraInfoIp(req.AccountId, req.ExtraInfo, req.Ip)

// 	// 设置userdata
// 	var userData = common.PlayerMgr.FindPlayerData(req.AccountId)
// 	if userData != nil {
// 		if userData.IsOnline() {
// 			//service.ServMgr.GetAccountService().OnClose(userData)
// 			service.ServMgr.GetAccountService().KnickOut(userData, msg.ErrCode_OTHER_LOGIN)
// 		}
// 		userData.UpdateTime = uint32(time.Now().Unix())
// 		userData.PlayerAgent = req.Agent
// 	} else {
// 		userData = common.NewPlayerData(req.UserId, req.Agent)
// 		userData.AccountInfo = req.Account
// 		//common.PlayerMgr.AddPlayerData(userData)
// 	}
// 	if userData.AccountInfo.AccountId == 0 {
// 		userData.AccountInfo.AccountId = req.AccountId
// 	}

// 	userData.ChannelId = req.ChanelId
// 	userData.TdaCommonAttr = req.TdaCommonAttr
// 	userData.Ip = req.Ip

// 	timeNow := time.Now()
// 	// rbi login
// 	rbiData := &rbi.PlayerLogin{
// 		GameSvrId:      req.TdaCommonAttr.SeverId,
// 		DtEventTime:    timeNow,
// 		VGameAppid:     "",
// 		PlatID:         rbi.PlatMap[strings.ToLower(req.TdaCommonAttr.Os)],
// 		IZoneAreaID:    utils.StrToInt(req.TdaCommonAttr.SeverId),
// 		VOpenID:        req.TdaCommonAttr.GopOpenId,
// 		VRoleID:        strconv.FormatInt(userData.AccountInfo.AccountId, 10),
// 		VRoleName:      strconv.FormatInt(userData.AccountInfo.AccountId, 10),
// 		VClientIP:      req.Ip,
// 		Region:         req.TdaCommonAttr.URegion,
// 		Country:        req.TdaCommonAttr.Country,
// 		GarenaOpenID:   req.TdaCommonAttr.GopOpenId,
// 		Timekey:        timeNow.Unix(),
// 		ClientVersion:  "",
// 		SystemSoftware: req.TdaCommonAttr.Os_version,
// 		SystemHardware: req.TdaCommonAttr.Device_model,
// 		TelecomOper:    req.TdaCommonAttr.Manufacturer,
// 		Network:        req.TdaCommonAttr.Network_type,
// 		ScreenWidth:    utils.StrToInt(req.TdaCommonAttr.Screen_width),
// 		ScreenHight:    utils.StrToInt(req.TdaCommonAttr.Screen_height),
// 		Density:        0,
// 		CpuHardware:    "",
// 		Memory:         0,
// 		GLRender:       "",
// 		GLVersion:      "",
// 		DeviceId:       req.TdaCommonAttr.Device_id,
// 		GenderType:     0,
// 		ILevel:         int(userData.AccountInfo.Level),
// 		RegisterTime:   time.Unix(int64(userData.AccountInfo.CreateTime), 0),
// 		RoleTotalCash:  0,
// 	}
// 	rbi.RbiWrite(rbiData)

// 	common.PlayerMgr.AddPlayerData(userData)
// 	common.PlayerMgr.AddPlayerBasic(userData)
// 	dao.AccountDao.UpdateAccountLogin(req.AccountId)
// 	userData.State = publicconst.Online
// 	lastLoginTime := userData.AccountInfo.LoginTime
// 	userData.AccountInfo.LoginTime = tools.GetCurTime()
// 	req.Agent.SetUserData(userData)
// 	req.Agent.SetAgentId(req.AccountId)
// 	req.Agent.SetUserId(userData.UserId)

// 	// 登录成功做一些处理
// 	service.ServMgr.GetAccountService().LoginSucc(userData)

// 	res := &msg.ResponseLogin{
// 		Result: msg.ErrCode_SUCC,
// 	}
// 	res.ServerTime = time.Now().UnixMilli()
// 	res.Info = service.ToProtocolAccountInfo(userData.AccountInfo)
// 	res.ApInfo = service.ToProtocolApInfo(userData.AccountInfo)
// 	res.MissionId = uint32(userData.AccountInfo.MissionId)
// 	for i := 0; i < len(userData.AccountInfo.GuideData); i++ {
// 		res.GuideData = append(res.GuideData, &msg.GuideInfo{Id: userData.AccountInfo.GuideData[i].Id,
// 			Value: userData.AccountInfo.GuideData[i].Value})
// 	}
// 	for i := 0; i < len(userData.AccountInfo.PopUps); i++ {
// 		res.PopUps = append(res.PopUps, &msg.PopUpInfo{Id: userData.AccountInfo.PopUps[i].Id, PopType: userData.AccountInfo.PopUps[i].PopUpType})
// 	}
// 	toProtocolChargeInfo(userData, res)
// 	toProtocolFundInfo(userData, res)
// 	toProtocolAdInfo(userData, res)
// 	res.DailyAp = service.ToProtocolDailApInfo(userData.AccountInfo.DailyApData)
// 	res.MonthCardDailyRewardTime = userData.AccountInfo.MonthCardDailyRewardTime
// 	res.TalentData = service.ToProtocolTalentData(userData.AccountInfo.TalentData)
// 	res.OpenServerDays = common.GetOpenServerDays()

// 	tools.SendMsg(req.Agent, res, req.PackId, res.Result)

// 	service.ServMgr.GetPokerService().SendPokerNtf(userData)

// 	service.ServMgr.GetCommonService().AddStaticsData(userData, publicconst.Statics_Login_Id, fmt.Sprintf("loginTime:%v,lastLoginTime:%v|", userData.AccountInfo.LoginTime, lastLoginTime))
// 	log.Info("Login succ", zap.String("userId", req.UserId), zap.Int64("accountId", req.AccountId))
// 	//log.Infof("Login userId:%v id:%v login succ", userId, userData.AccountInfo.AccountId)
// }

// RequestClientHeartHandle 处理客户端心跳
func RequestClientHeartHandle(packetId uint32, args interface{}, p *player.Player) {
	service.OnHeart(p)
	nowAt := time.Now().UnixMilli()
	p.SendResponse(packetId, &msg.ResponseClientHert{
		Result:     msg.ErrCode_SUCC,
		ServerTime: nowAt,
	}, msg.ErrCode_SUCC)
}

// RequestLogoutHandle 客户端退出
func RequestLogoutHandle(packetId uint32, args interface{}, p *player.Player) {
	// p.SendMsg(0, &msg.ResponseLogout{
	// 	Result: msg.ErrCode_SUCC,
	// })
	//charge.GetOrderManager().DelUserOrder(p.GetUserId()) // 退出删除玩家订单
	service.PlayerLogout(p)
}

// RequestUpdateNickHandle 请求更新昵称
func RequestUpdateNickHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUpdateNick)
	res := &msg.ResponseUpdateNick{}
	//res.Result = service.UpdateNick(p, req.Nick)
	//if res.Result == msg.ErrCode_SUCC {
	//	if err := content.PushContent(&common.ContentData{
	//		Content:     req.Nick,
	//		UserId:      p.UserData.AccountId,
	//		ChannelNo:   p.SdkChannelNo,
	//		Os:          p.Os,
	//		RoleId:      strconv.Itoa(int(p.GetUserId())),
	//		RoleName:    p.UserData.Nick,
	//		RoleLevel:   strconv.Itoa(int(p.UserData.Level)),
	//		ServerId:    strconv.Itoa(int(p.UserData.ServerId)),
	//		Ip:          p.Ip,
	//		ContentType: common.ContentTypeNick,
	//		Cb: func(ec msg.ErrCode) {
	//			if ec == msg.ErrCode_SUCC {
	//				res.Nick = p.UserData.Nick
	//				p.SendResponse(packetId, res, res.Result)
	//			}
	//
	//			res.Result = ec
	//			log.Debug("RequestWorldChatHandle", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res))
	//			p.SendResponse(packetId, res, res.Result)
	//		},
	//	}); err != nil {
	//		p.SendResponse(packetId, res, res.Result)
	//		log.Error("RequestWorldChatHandle failed", zap.Error(err))
	//	}
	//} else {
	//	res.Nick = p.UserData.Nick
	//	p.SendResponse(packetId, res, res.Result)
	//}
	if code := service.UpdateNickCheck(p, req.Nick); code != msg.ErrCode_SUCC {
		res.Result = code
		p.SendResponse(packetId, res, res.Result)
		return
	}
	userId := p.GetUserId()
	if err := content.PushContent(&content.ContentData{
		Content:     []string{req.Nick},
		UserId:      p.UserData.AccountId,
		ChannelNo:   p.SdkChannelNo,
		Os:          p.Os,
		RoleId:      strconv.Itoa(int(p.GetUserId())),
		RoleName:    p.UserData.Nick,
		RoleLevel:   strconv.Itoa(int(p.UserData.Level)),
		ServerId:    strconv.Itoa(int(p.UserData.ServerId)),
		Ip:          p.UserData.BaseInfo.Ip,
		ContentType: content.ContentTypeNick,
		Cb: func(ec msg.ErrCode) {
			nickPlayer := player.FindByUserId(userId)
			if nickPlayer != nil {
				res.Result = ec
				if ec == msg.ErrCode_SUCC {
					res.Result = service.UpdateNick(nickPlayer, req.Nick)
				}
				res.Nick = nickPlayer.UserData.Nick
				//log.Debug("RequestUpdateNickHandle PushContent ", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res))
				nickPlayer.SendResponse(packetId, res, res.Result)
			}
		},
	}); err != nil {
		log.Error("RequestUpdateNickHandle failed", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		res.Nick = p.UserData.Nick
		res.Result = msg.ErrCode_SYSTEM_ERROR
		p.SendResponse(packetId, res, res.Result)
	}
}

// RequestUpgradeHandle 升级
func RequestUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseUpgrade{}
	res.Result, res.Level = service.Upgrade(p)
	p.SendResponse(packetId, res, res.Result)
}

// RequestSetShipHandle 设置出战机甲
func RequestSetShipHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetShip)
	res := &msg.ResponseSetShip{
		ShipId: req.ShipId,
	}
	res.Result = service.SetShipId(p, req.ShipId)
	p.SendResponse(packetId, res, res.Result)
}

// RequestSetSupportShipHandle 设置支援机甲
func RequestSetSupportShipHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetSupportShip)
	res := &msg.ResponseSetSupportShip{
		ShipIds: req.ShipIds,
	}
	res.Result = service.SetSupportId(p, req.ShipIds)
	p.SendResponse(packetId, res, res.Result)
}

// RequestSetPlayerNameAndShipReq 初始化玩家名和机甲
func RequestSetPlayerNameAndShipReq(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.InitPlayerNameAndShipReq)
	res := &msg.InitPlayerNameAndShipRsp{}
	if code := service.OnInitPlayerNameAndShipCheck(p, req.Name, req.ShipId); code != msg.ErrCode_SUCC {
		res.Result = code
		p.SendResponse(packetId, res, res.Result)
		return
	}

	userId := p.GetUserId()

	if err := content.PushContent(&content.ContentData{
		Content:     []string{req.Name},
		UserId:      p.UserData.AccountId,
		ChannelNo:   p.SdkChannelNo,
		Os:          p.Os,
		RoleId:      strconv.Itoa(int(p.GetUserId())),
		RoleName:    p.UserData.Nick,
		RoleLevel:   strconv.Itoa(int(p.UserData.Level)),
		ServerId:    strconv.Itoa(int(p.UserData.ServerId)),
		Ip:          p.UserData.BaseInfo.Ip,
		ContentType: content.ContentTypeNick,
		Cb: func(ec msg.ErrCode) {
			initPlayer := player.FindByUserId(userId)
			if initPlayer != nil {
				res.Result = ec
				if ec == msg.ErrCode_SUCC {
					errCode := service.OnInitPlayerNameAndShip(initPlayer, req.Name, req.ShipId)
					res.Result = errCode
					if errCode == msg.ErrCode_SUCC {
						res.Name = req.Name
						res.ShipId = req.ShipId
					}
				}
				//log.Debug("RequestSetPlayerNameAndShipReq PushContent ", zap.Uint64("uid", userId), zap.Any("req", req), zap.Any("res", res))
				initPlayer.SendResponse(packetId, res, res.Result)
			}

		},
	}); err != nil {
		log.Error("RequestSetPlayerNameAndShipReq failed", zap.Uint64("uid", p.GetUserId()), zap.Any("req", req), zap.Any("res", res), zap.Error(err))
		res.Result = msg.ErrCode_SYSTEM_ERROR
		p.SendResponse(packetId, res, res.Result)
	}
}

// RequestRandomPlayerName 随机生成玩家名
func RequestRandomPlayerName(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.RandomPlayerNameResp{}
	name := service.OnRandomGenPlayerName(p)
	res.Name = name
	p.SendResponse(packetId, res, msg.ErrCode_SUCC)
}

// RequestGlobalAttrDetailHandle 请求全局属性详细
func RequestGlobalAttrDetailHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseGlobalAttrDetail{}
	res.Result, res.Data = service.GetGlobalAttrDetail(p)
	p.SendResponse(packetId, res, res.Result)
}

// func checkLogining(agent gate.Agent) bool {
// 	userData := agent.UserData()
// 	if userData != nil {
// 		if playerData := userData.(*player.Player); playerData != nil {
// 			if playerData.State == publicconst.Logining {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func toProtocolChargeInfo(p *player.Player, res *msg.ResponseLogin) {
	// 充值信息
	res.ReInfo = &msg.RechargeInfo{}
	for _, v := range p.UserData.BaseInfo.Charge {
		if v.Value > 0 {
			res.ReInfo.RechargeIds = append(res.ReInfo.RechargeIds, uint32(v.Id))
		}
	}

	for _, v := range p.UserData.BaseInfo.MonthCard {
		pbData := &msg.MonthcardInfo{
			MonthCardId:       uint32(v.Id),
			EndTime:           uint32(v.EndTime),
			NextGetRewardTime: uint32(v.NextGetRewardTime),
		}
		res.ReInfo.McInfo = append(res.ReInfo.McInfo, pbData)
	}
}

func toProtocolFundInfo(p *player.Player, res *msg.ResponseLogin) {
	// 充值信息
	res.FundInfo = make([]*msg.MainFundInfo, 0, len(p.UserData.BaseInfo.MainFund))
	for _, v := range p.UserData.BaseInfo.MainFund {
		pbData := &msg.MainFundInfo{
			FundId:      uint32(v.Id),
			RewardMaxId: uint32(v.FreeId),
			BuyFlag:     uint32(v.BuyFlag),
		}
		res.FundInfo = append(res.FundInfo, pbData)
	}
}

func toProtocolAdInfo(p *player.Player, res *msg.ResponseLogin) {
	// 充值信息
	for i := 0; i < len(p.UserData.BaseInfo.Ad); i++ {
		res.AdData = append(res.AdData, &msg.AdInfo{
			AdId:  p.UserData.BaseInfo.Ad[i].AdId,
			Times: p.UserData.BaseInfo.Ad[i].Times,
		})
	}
}
