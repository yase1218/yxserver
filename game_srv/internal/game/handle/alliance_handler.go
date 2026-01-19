package handle

// import (
// 	"kernel/tda"
// 	"kernel/tools"
// 	"math/rand"
// 	"msg"
// 	"strconv"
// 	"time"
// 	"unicode/utf8"

// 	"github.com/v587-zyf/gc/log"
// 	"github.com/zy/game_data/template"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.uber.org/zap"

// 	"gameserver/internal/game/common"
// 	"server/internal/game/dao"
// 	"server/internal/game/model"
// 	"server/internal/game/service"
// 	"server/internal/publicconst"
// )

// // HandleSearchAlliance 处理搜索联盟请求
// func HandleSearchAlliance(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.SearchAllianceReq)

// 	alliances := make([]model.Alliance, 0)
// 	var err error

// 	// 检查玩家是否已有联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err == nil && member != nil {
// 		// 获取目标联盟
// 		alliance, err := dao.GetAlliance(member.AllianceID)
// 		if err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.SearchAllianceRsp{}, packetId, msg.ErrCode_SYSTEM_ERROR)
// 			return
// 		}
// 		if alliance != nil {
// 			// 广播联盟基础信息
// 			bMsg := &msg.NotifyAllianceInfo{
// 				Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 			}
// 			tools.SendMsg(playerData.PlayerAgent, bMsg, packetId, msg.ErrCode_SUCC)
// 			return
// 		}
// 	}

// 	// 构建搜索条件
// 	filter := bson.M{}
// 	if req.Id != "" || req.Name != "" {
// 		if req.Id != "" {
// 			// 支持按名称或ID搜索
// 			if id, err := strconv.ParseInt(req.Id, 10, 64); err == nil {
// 				filter["_id"] = id
// 			}

// 		} else if req.Name != "" {
// 			filter["name"] = bson.M{"$regex": req.Name}
// 		}
// 		// 搜索联盟
// 		// 人数不满的前50联盟
// 		alliances, err = dao.SearchAlliance(filter, 50)
// 		if err != nil {
// 			alliances = nil
// 		}
// 		if alliances == nil || len(alliances) == 0 {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.SearchAllianceRsp{}, packetId, msg.ErrCode_ALLIANCE_NOT_FOUND_LINK)
// 			return
// 		}
// 	} else {
// 		//
// 		//alls := make(map[uint32]uint32, 0)
// 		//// 搜索联盟
// 		//applys, _ := dao.GetApplicationByAccountIDByCD(playerData.AccountInfo.AccountId)
// 		//if applys != nil {
// 		//	for _, i2 := range applys {
// 		//		alls[i2.AllianceID] = i2.AllianceID
// 		//	}
// 		//}
// 		//alls2 := make([]int64, 0)
// 		//for _, u2 := range alls {
// 		//	alls2 = append(alls2, int64(u2))
// 		//}

// 		// 人数不满的前50联盟
// 		alliances, err = dao.SearchAllianceCanJoin(filter, 50, nil)
// 		if err != nil {
// 			//tools.SendMsg(playerData.PlayerAgent, &msg.SearchAllianceRsp{}, packetId, msg.ErrCode_ALLIANCE_NOT_FOUND_LINK)
// 			//return
// 			alliances = nil
// 		}
// 	}

// 	if alliances == nil {
// 		alliances = make([]model.Alliance, 0)
// 	}

// 	// 转换为响应消息
// 	rsp := &msg.SearchAllianceRsp{
// 		Alliances: make([]*msg.AllianceInfoA, 0, len(alliances)),
// 	}

// 	for _, a := range alliances {
// 		_, leaderInfo := service.ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(a.LeaderID))
// 		if leaderInfo != nil {
// 			rsp.Alliances = append(rsp.Alliances, &msg.AllianceInfoA{
// 				AllianceId:    a.ID,
// 				Name:          a.Name,
// 				Banner:        a.Banner,
// 				Level:         a.Level,
// 				MemberCount:   a.MemberCount,
// 				MemberLimit:   a.MaxMemberCount,
// 				PowerRequired: a.PowerRequired,
// 				TotalPower:    uint32(a.TotalPower),
// 				ChestCount:    a.TotalTreasure,
// 				LeaderName:    leaderInfo.Name,
// 			})
// 		}
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// }

// //// HandleSearchAlliance2 处理搜索联盟请求
// //func HandleSearchAlliance2(packetId uint32, args interface{}, playerData *common.PlayerData) {
// //	req := args.(*msg.SearchAllianceReq)
// //
// //	_, alls := service.ServMgr.GetAllianceService().Search(playerData, req.Id, req.Name)
// //
// //	// 转换为响应消息
// //	rsp := &msg.SearchAllianceRsp{
// //		Alliances: make([]*msg.AllianceInfoA, 0, len(alls)),
// //	}
// //
// //	for _, a := range alls {
// //		_, leaderInfo := service.ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(a.LeaderID))
// //		if leaderInfo != nil {
// //			rsp.Alliances = append(rsp.Alliances, &msg.AllianceInfoA{
// //				AllianceId:    a.ID,
// //				Name:          a.Name,
// //				Banner:        a.Banner,
// //				Level:         a.Level,
// //				MemberCount:   a.MemberCount,
// //				MemberLimit:   a.MaxMemberCount,
// //				PowerRequired: a.PowerRequired,
// //				TotalPower:    uint32(a.TotalPower),
// //				ChestCount:    a.TotalTreasure,
// //				LeaderName:    leaderInfo.Name,
// //			})
// //		}
// //	}
// //	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// //}

// // HandleCreateAlliance 处理创建联盟请求
// func HandleCreateAlliance(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.CreateAllianceReq)

// 	// 检查玩家是否已有联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err == nil && member != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_HAS_LINK,
// 		}, packetId, msg.ErrCode_ALLIANCE_HAS_LINK)
// 		return
// 	}

// 	// 检查消耗道具
// 	for _, item := range template.GetSystemItemTemplate().AllianceCreateCostItems {
// 		if !service.ServMgr.GetItemService().EnoughItem(playerData.AccountInfo.AccountId, item.ItemId, item.ItemNum) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 				Result: msg.ErrCode_NO_ENOUGH_ITEM,
// 			}, packetId, msg.ErrCode_NO_ENOUGH_ITEM)
// 			return
// 		}
// 	}

// 	if uint32(utf8.RuneCountInString(req.Declaration)) > template.GetSystemItemTemplate().AllianceDeclareLen {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_DECLARATION_LEN,
// 		}, packetId, msg.ErrCode_ALLIANCE_DECLARATION_LEN)
// 		return
// 	}

// 	//name = strings.Trim(name, " ")
// 	//if len(name) == 0 || pet.Name == name {
// 	//	return msg.ErrCode_INVALID_DATA
// 	//}

// 	//if uint32(len(req.Declaration)) > template.GetSystemItemTemplate().AllianceDeclareLen {
// 	//	tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 	//		Result: msg.ErrCode_ALLIANCE_PARAM_ERROR,
// 	//	}, packetId, msg.ErrCode_SUCC)
// 	//	return
// 	//}

// 	nameLen := uint32(utf8.RuneCountInString(req.Name))
// 	if nameLen < 1 || nameLen > template.GetSystemItemTemplate().AllianceNameLen {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NAME_LEN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NAME_LEN)
// 		return
// 	}

// 	if template.GetForbiddenTemplate().HasForbidden(req.Name) {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NAME_HAS_FORBIDDEN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NAME_HAS_FORBIDDEN)
// 		return
// 	}

// 	if uint32(len(req.Declaration)) > 0 {
// 		if template.GetForbiddenTemplate().HasForbidden(req.Declaration) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_DECLARATION_HAS_FORBIDDEN,
// 			}, packetId, msg.ErrCode_ALLIANCE_DECLARATION_HAS_FORBIDDEN)
// 			return
// 		}
// 	}

// 	// 创建新联盟
// 	alliance := &model.Alliance{
// 		ID:             uint32(time.Now().Unix()),
// 		Name:           req.Name,
// 		Banner:         req.Banner,
// 		Level:          1,
// 		Declaration:    req.Declaration,
// 		LeaderID:       playerData.AccountInfo.AccountId,
// 		PowerRequired:  0,
// 		MemberCount:    0,
// 		MaxMemberCount: uint32(template.GetGuildTemplate().GetAlliance(1).Num),
// 		TotalPower:     0,
// 		TotalTreasure:  0,
// 		CreateTime:     time.Now(),
// 		TreasureRule:   0,
// 		Exp:            0,
// 		AutoJoin:       true,
// 	}

// 	if err := dao.CreateAlliance(alliance); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NAME_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_NAME_EXIST)
// 		return
// 	}

// 	// 添加创建者为成员
// 	member = &model.AllianceMember{
// 		AllianceID:   alliance.ID,
// 		PlayerID:     playerData.AccountInfo.AccountId,
// 		Name:         playerData.AccountInfo.Nick,
// 		HeadImg:      playerData.AccountInfo.HeadImg,
// 		HeadFrame:    playerData.AccountInfo.HeadFrame,
// 		Position:     1, // 指挥官
// 		Power:        playerData.AccountInfo.Combat,
// 		WeeklyActive: 0,
// 		IsMuted:      false,
// 		JoinTime:     time.Now(),
// 		LastOnline:   time.Now(),
// 	}
// 	if err := dao.AddMember(member); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 		}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 		return
// 	}

// 	alliance.MemberCount = 1
// 	alliance.TotalPower = uint64(member.Power)

// 	if alliance != nil {
// 		// 广播联盟基础信息
// 		bMsg := &msg.NotifyAllianceInfo{
// 			Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 		}
// 		tools.SendMsg(playerData.PlayerAgent, bMsg, packetId, msg.ErrCode_SUCC)
// 	}

// 	// 扣除消耗道具
// 	for _, item := range template.GetSystemItemTemplate().AllianceCreateCostItems {
// 		service.ServMgr.GetItemService().CostItem(playerData.AccountInfo.AccountId, item.ItemId, item.ItemNum, publicconst.CreateAllianceCostItem, true)
// 	}

// 	AddAllianceRank(alliance.ID, playerData.AccountInfo.AccountId, playerData.AccountInfo.Combat)

// 	// 发送响应
// 	tools.SendMsg(playerData.PlayerAgent, &msg.CreateAllianceRsp{
// 		Result: msg.ErrCode_SUCC,
// 	}, packetId, msg.ErrCode_SUCC)

// 	// tda
// 	tda.TdaGuildCreate(playerData.ChannelId, playerData.TdaCommonAttr, alliance.Level, alliance.ID, alliance.Name)
// }

// // HandleJoinAlliance 处理加入联盟请求
// func HandleJoinAlliance(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.JoinAllianceReq)

// 	// AllianceLeaveTime 是最近一小时以内, 就不允许加入
// 	if playerData.AccountInfo.AllianceLeaveTime > tools.GetCurTime()-template.GetSystemItemTemplate().AllianceLeaveCD {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			//Result: msg.ErrCode_ALLIANCE_EXIT_COOLDOWN,
// 			CDTime: playerData.AccountInfo.AllianceLeaveTime + template.GetSystemItemTemplate().AllianceLeaveCD,
// 			//}, packetId, msg.ErrCode_ALLIANCE_EXIT_COOLDOWN)
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 检查玩家是否已有联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err == nil && member != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_HAS_LINK, // 已有联盟
// 		}, packetId, msg.ErrCode_ALLIANCE_HAS_LINK)
// 		return
// 	}

// 	// 获取目标联盟
// 	alliance, err := dao.GetAlliance(req.AllianceId)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_EXIST)
// 		return
// 	}

// 	// 检查联盟是否满员
// 	if alliance.MemberCount >= alliance.MaxMemberCount {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 		}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 		return
// 	}

// 	// 检查战力要求
// 	playerPower := playerData.AccountInfo.Combat
// 	if alliance.PowerRequired > 0 {
// 		if playerPower < alliance.PowerRequired {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 			return
// 		}
// 	}

// 	// 获取联盟成员列表
// 	userids := make([]uint32, 0)
// 	userids = append(userids, uint32(playerData.AccountInfo.AccountId))
// 	_, players := service.ServMgr.GetSocialService().GetPlayerSimpleInfos(userids)
// 	if players == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		return
// 	}

// 	// 同时进行申请冷却时间检查
// 	if applyInfo, err := dao.GetApplicationByAccountIDByCD2(alliance.ID, playerData.AccountInfo.AccountId); err == nil && applyInfo != nil {
// 		if applyInfo.PlayerID == playerData.AccountInfo.AccountId {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_IN_APPLY,
// 			}, packetId, msg.ErrCode_ALLIANCE_IN_APPLY)
// 			return
// 		}
// 	}
// 	apply := &model.AllianceApplication{
// 		ID:         time.Now().UnixNano(),
// 		AllianceID: alliance.ID,
// 		PlayerID:   playerData.AccountInfo.AccountId,
// 		Name:       playerData.AccountInfo.Nick,
// 		HeadImg:    playerData.AccountInfo.HeadImg,
// 		HeadFrame:  playerData.AccountInfo.HeadFrame,
// 		Power:      playerPower,
// 		ApplyTime:  time.Now(),
// 		Status:     0, // 待处理
// 	}
// 	// 如果允许自动加入
// 	if alliance.AutoJoin {
// 		// 直接批准加入
// 		if !doAddMemberToAlliance(packetId, alliance, apply.PlayerID, apply.Name, apply.HeadImg, apply.HeadFrame, apply.Power) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 			}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 			return
// 		} else {
// 			// 广播联盟基础信息
// 			bMsg := &msg.NotifyAllianceInfo{
// 				Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 			}
// 			service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 			// 我在线
// 			service.InterSendUserMsg(func(arg interface{}, targerPlayer *common.PlayerData) {
// 				// 新成员提示
// 				notifyMsg := &msg.NotifyJoinAlliance{
// 					Name: alliance.Name,
// 				}
// 				tools.SendNotifyMsg(playerData.PlayerAgent, notifyMsg)
// 			}, playerData, playerData)
// 		}
// 	} else {
// 		// 添加申请记录
// 		dao.DeleteApplication(alliance.ID, playerData.AccountInfo.AccountId)
// 		if err := dao.AddApplication(apply); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 			}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 			return
// 		}

// 		// 获取申请列表
// 		if applies, err := dao.GetApplicationList(alliance.ID); err == nil {
// 			service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
// 				Count: int32(len(applies)),
// 			}, 0, msg.ErrCode_SUCC, 0, 0)
// 		}
// 	}

// 	AddAllianceRank(alliance.ID, playerData.AccountInfo.AccountId, playerData.AccountInfo.Combat)

// 	service.ServMgr.GetTaskService().UpdateTask(
// 		playerData, true, publicconst.TASK_COND_ALLIANCE_JOIN, 1)

// 	tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 		Result: msg.ErrCode_SUCC,
// 	}, packetId, msg.ErrCode_SUCC)

// 	// tda
// 	tda.TdaGuildJoinOrLeave(playerData.ChannelId, playerData.TdaCommonAttr, alliance.Level, alliance.ID, alliance.Name, "1")
// }

// // HandleQuickJoinAlliance 处理快速加入联盟请求
// func HandleQuickJoinAlliance(packetId uint32, args interface{}, playerData *common.PlayerData) {

// 	// AllianceLeaveTime 是最近一小时以内, 就不允许加入
// 	if playerData.AccountInfo.AllianceLeaveTime > tools.GetCurTime()-template.GetSystemItemTemplate().AllianceLeaveCD {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 			//Result: msg.ErrCode_ALLIANCE_EXIT_COOLDOWN,
// 			CDTime: playerData.AccountInfo.AllianceLeaveTime + template.GetSystemItemTemplate().AllianceLeaveCD,
// 			//}, packetId, msg.ErrCode_ALLIANCE_EXIT_COOLDOWN)
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 检查玩家是否已有联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err == nil && member != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_HAS_LINK, // 已有联盟
// 		}, packetId, msg.ErrCode_ALLIANCE_HAS_LINK)
// 		return
// 	}
// 	//
// 	alls := make(map[uint32]uint32, 0)
// 	// 搜索联盟
// 	applys, _ := dao.GetApplicationByAccountIDByCD(playerData.AccountInfo.AccountId)
// 	if applys != nil {
// 		for _, i2 := range applys {
// 			alls[i2.AllianceID] = i2.AllianceID
// 		}
// 	}
// 	alls2 := make([]int64, 0)
// 	for _, u2 := range alls {
// 		alls2 = append(alls2, int64(u2))
// 	}
// 	// 获取推荐联盟列表
// 	playerPower := dao.GetPlayerPower(playerData.AccountInfo.AccountId)
// 	alliances, err := dao.SearchAllianceCanJoin(bson.M{
// 		//"mem_count": bson.M{"$lt": publicconst.ALLIANCE_MAX_MEMBER},
// 		//"mem_count": bson.M{"$lt": "$max_count"},
// 		"power_req": bson.M{"$lte": playerPower},
// 	}, 0, alls2)
// 	if err != nil || len(alliances) == 0 {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_FOUND_LINK, // 找不到合适联盟
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_FOUND_LINK)
// 		return
// 	}

// 	var availableAlliances []model.Alliance
// 	for _, alliance := range alliances {
// 		if alliance.MemberCount < alliance.MaxMemberCount {
// 			availableAlliances = append(availableAlliances, alliance)
// 		}
// 	}
// 	if len(availableAlliances) == 0 {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_FOUND_LINK,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_FOUND_LINK)
// 		return
// 	}

// 	// 随机选择一个联盟
// 	alliance := availableAlliances[rand.Intn(len(availableAlliances))]

// 	if applyInfo, err := dao.GetApplicationByAccountID(alliance.ID, playerData.AccountInfo.AccountId); err == nil && applyInfo != nil {
// 		if applyInfo.PlayerID == playerData.AccountInfo.AccountId {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_IN_APPLY,
// 			}, packetId, msg.ErrCode_ALLIANCE_IN_APPLY)
// 			return
// 		}
// 	}

// 	// 添加申请记录
// 	apply := &model.AllianceApplication{
// 		ID:         time.Now().UnixNano(),
// 		AllianceID: alliance.ID,
// 		PlayerID:   playerData.AccountInfo.AccountId,
// 		Name:       playerData.AccountInfo.Nick,
// 		HeadImg:    playerData.AccountInfo.HeadImg,
// 		HeadFrame:  playerData.AccountInfo.HeadFrame,
// 		Power:      playerPower,
// 		ApplyTime:  time.Now(),
// 		Status:     0, // 待处理
// 	}
// 	if alliance.AutoJoin {
// 		// 直接批准加入
// 		if !doAddMemberToAlliance(packetId, &alliance, apply.PlayerID, apply.Name, apply.HeadImg, apply.HeadFrame, apply.Power) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 			}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 			return
// 		} else {
// 			// 广播联盟基础信息
// 			bMsg := &msg.NotifyAllianceInfo{
// 				Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(&alliance),
// 			}
// 			service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 			// 我在线
// 			service.InterSendUserMsg(func(arg interface{}, targerPlayer *common.PlayerData) {
// 				// 新成员提示
// 				notifyMsg := &msg.NotifyJoinAlliance{
// 					Name: alliance.Name,
// 				}
// 				tools.SendNotifyMsg(targerPlayer.PlayerAgent, notifyMsg)
// 			}, playerData, playerData)
// 		}
// 	} else {
// 		dao.DeleteApplication(alliance.ID, playerData.AccountInfo.AccountId)
// 		if err := dao.AddApplication(apply); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 				Result: msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		// 获取申请列表
// 		if applies, err := dao.GetApplicationList(alliance.ID); err == nil {
// 			service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
// 				Count: int32(len(applies)),
// 			}, 0, msg.ErrCode_SUCC, 0, 0)
// 		}
// 	}

// 	AddAllianceRank(alliance.ID, playerData.AccountInfo.AccountId, playerData.AccountInfo.Combat)

// 	tools.SendMsg(playerData.PlayerAgent, &msg.QuickJoinAllianceRsp{
// 		Result: msg.ErrCode_SUCC,
// 	}, packetId, msg.ErrCode_SUCC)

// 	// tda
// 	tda.TdaGuildJoinOrLeave(playerData.ChannelId, playerData.TdaCommonAttr, alliance.Level, alliance.ID, alliance.Name, "1")
// }

// // HandleGetAllianceMembers 处理获取联盟成员列表请求
// func HandleGetAllianceMembers(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.GetAllianceMembersRsp{
// 			Members: make([]*msg.AllianceMemberInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 获取联盟成员列表
// 	members, err := dao.GetMemberList(member.AllianceID)
// 	if err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.GetAllianceMembersRsp{
// 			Members: make([]*msg.AllianceMemberInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 转换为响应消息
// 	rsp := &msg.GetAllianceMembersRsp{
// 		Members: make([]*msg.AllianceMemberInfo, 0, len(members)),
// 	}
// 	// 获取联盟成员列表
// 	userids := make([]uint32, len(members))
// 	for _, member := range members {
// 		userids = append(userids, uint32(member.PlayerID))
// 	}
// 	_, players := service.ServMgr.GetSocialService().GetPlayerSimpleInfos(userids)
// 	if players == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.GetAllianceMembersRsp{
// 			Members: make([]*msg.AllianceMemberInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}
// 	for _, m := range members {
// 		rsp.Members = append(rsp.Members, &msg.AllianceMemberInfo{
// 			Data:           service.ToPlayerSimpleInfo(players[uint32(m.PlayerID)]),
// 			Position:       msg.AlliancePosition(m.Position),
// 			WeeklyActivity: m.WeeklyActive,
// 			LastOnlineTime: players[uint32(m.PlayerID)].LastOnlineTime,
// 		})
// 	}

// 	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// }

// // HandleQuitAlliance 处理退出联盟请求
// func HandleQuitAlliance(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuitAllianceRsp{
// 			AccountId: uint32(playerData.AccountInfo.AccountId),
// 			Result:    msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		log.Error("alliance nil", zap.Int64("accountId", playerData.GetAccountId()),
// 			zap.Uint32("allianceId", member.AllianceID), zap.Error(err))

// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuitAllianceRsp{
// 			AccountId: uint32(playerData.AccountInfo.AccountId),
// 			Result:    msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}

// 	// 指挥官 检查人数，最后一人解散联盟 否则不能退出
// 	if member.Position == 1 {
// 		if alliance.MemberCount > 1 {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.QuitAllianceRsp{
// 				AccountId: uint32(playerData.AccountInfo.AccountId),
// 				Result:    msg.ErrCode_ALLIANCE_LEADER_CANNOT_EXIT,
// 			}, packetId, msg.ErrCode_ALLIANCE_LEADER_CANNOT_EXIT)
// 			return
// 		}

// 		// del alliance
// 		if err = dao.DeleteAlliance(int64(member.AllianceID)); err != nil {
// 			log.Error("del alliance err", zap.Error(err),
// 				zap.Int64("accountId", playerData.GetAccountId()), zap.Uint32("allianceId", member.AllianceID))
// 		}
// 	}

// 	// 删除成员
// 	if err := dao.RemoveMember(member.PlayerID, member.AllianceID); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.QuitAllianceRsp{
// 			AccountId: uint32(playerData.AccountInfo.AccountId),
// 			Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	currTime := tools.GetCurTime()
// 	playerData.AccountInfo.AllianceLeaveTime = currTime
// 	dao.AccountDao.UpdateAllianceLeaveTime(playerData.AccountInfo.AccountId, currTime)

// 	DelAllianceRank(member.AllianceID, playerData.AccountInfo.AccountId)

// 	bMsg := &msg.QuitAllianceRsp{
// 		AccountId: uint32(playerData.AccountInfo.AccountId),
// 		Result:    msg.ErrCode_SUCC,
// 	}

// 	tools.SendMsg(playerData.PlayerAgent, bMsg, packetId, msg.ErrCode_SUCC)

// 	service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(member.AllianceID,
// 		bMsg, packetId, msg.ErrCode_SUCC,
// 		playerData.GetAccountId(), 0)

// 	// tda
// 	tda.TdaGuildJoinOrLeave(playerData.ChannelId, playerData.TdaCommonAttr, alliance.Level, alliance.ID, alliance.Name, "2")
// }

// // HandleAllianceManage 处理联盟管理请求
// func HandleAllianceManage(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.AllianceManageReq)
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}

// 	// 检查权限
// 	if member.Position > 2 { // 只有指挥官和副指挥官可以管理
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	// 获取目标成员
// 	targetMember, err := dao.GetMember(int64(req.TargetAccountId))
// 	if err != nil || targetMember == nil || targetMember.AllianceID != member.AllianceID {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_MEMBER_NOT_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_MEMBER_NOT_EXIST)
// 		return
// 	}

// 	// 更新联盟成员数
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_NOT_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_EXIST)
// 	}

// 	// 目标不能是自己
// 	if member.PlayerID == targetMember.PlayerID {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	// 执行管理操作
// 	switch req.OperateType {
// 	case 1: // 设置职位
// 		if member.Position != 1 { // 只有指挥官可以任命
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		if err := dao.UpdateMember(targetMember.PlayerID, bson.M{"position": req.NewPosition}); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}

// 		bMsg := &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_SUCC,
// 		}
// 		// 广播联盟基础信息
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(member.AllianceID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)
// 	case 2: // 转让指挥官
// 		if member.Position != 1 { // 只有指挥官可以转让
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		// 更新目标成员为指挥官
// 		if err := dao.UpdateMember(targetMember.PlayerID, bson.M{"position": msg.AlliancePosition_LEADER}); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		// 更新自己为副指挥官
// 		if err := dao.UpdateMember(member.PlayerID, bson.M{"position": msg.AlliancePosition_MEMBER}); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}

// 		upV := bson.M{}
// 		upV["leader_id"] = req.TargetAccountId
// 		err = dao.UpdateAlliance(member.AllianceID, upV)
// 		if err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		alliance.LeaderID = int64(req.TargetAccountId)

// 		bMsg := &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     msg.AlliancePosition_LEADER,
// 			Result:          msg.ErrCode_SUCC,
// 		}
// 		// 广播联盟基础信息
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(member.AllianceID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 		bMsg4 := &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: uint32(playerData.GetAccountId()),
// 			NewPosition:     msg.AlliancePosition_MEMBER,
// 			Result:          msg.ErrCode_SUCC,
// 		}
// 		// 广播联盟基础信息
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(member.AllianceID, bMsg4, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 		// 广播联盟基础信息
// 		bMsg2 := &msg.NotifyAllianceInfo{
// 			Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 		}
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg2, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 	case 3: // 踢出成员
// 		if targetMember.Position <= member.Position { // 不能踢出职位相同或更高的成员
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		if err := dao.RemoveMember(targetMember.PlayerID, member.AllianceID); err != nil {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 				OperateType:     req.OperateType,
// 				TargetAccountId: req.TargetAccountId,
// 				NewPosition:     req.NewPosition,
// 				Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 		bMsg := &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_SUCC,
// 		}
// 		// 广播联盟基础信息
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(member.AllianceID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), int64(req.TargetAccountId))

// 		// 广播联盟基础信息
// 		alliance.MemberCount--
// 		bMsg2 := &msg.NotifyAllianceInfo{
// 			Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 		}
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg2, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)

// 		currTime := tools.GetCurTime()
// 		if playerTarget := common.PlayerMgr.FindPlayerData(int64(req.TargetAccountId)); playerTarget != nil {
// 			service.InterSendUserMsg(func(msg interface{}, playerData2 *common.PlayerData) {
// 				playerData2.AccountInfo.AllianceLeaveTime = currTime
// 			}, nil, playerTarget)
// 		}

// 		DelAllianceRank(targetMember.AllianceID, targetMember.PlayerID)

// 		dao.AccountDao.UpdateAllianceLeaveTime(int64(req.TargetAccountId), currTime)
// 	default:
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceManageRsp{
// 			OperateType:     req.OperateType,
// 			TargetAccountId: req.TargetAccountId,
// 			NewPosition:     req.NewPosition,
// 			Result:          msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}
// }

// // HandleAllianceApplyList 处理获取联盟申请列表请求
// func HandleAllianceApplyList(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil || member.Position > 2 { // 只有指挥官和副指挥官可以查看
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceApplyListRsp{
// 			Applies: make([]*msg.AllianceApplyInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 获取申请列表
// 	applies, err := dao.GetApplicationList(member.AllianceID)
// 	if err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceApplyListRsp{
// 			Applies: make([]*msg.AllianceApplyInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 转换为响应消息
// 	rsp := &msg.AllianceApplyListRsp{
// 		Applies: make([]*msg.AllianceApplyInfo, 0, len(applies)),
// 	}
// 	for _, a := range applies {
// 		if a.Status == 0 { // 只返回待处理的申请
// 			rsp.Applies = append(rsp.Applies, &msg.AllianceApplyInfo{
// 				AccountId: uint32(a.PlayerID),
// 				Name:      a.Name,
// 				HeadImg:   a.HeadImg,
// 				HeadFrame: a.HeadFrame,
// 				Power:     uint32(a.Power),
// 				ApplyTime: uint32(a.ApplyTime.Unix()),
// 			})
// 		}
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// }

// // HandleAllianceApplyList2 处理获取联盟申请列表请求
// func HandleAllianceApplyList2(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	// 获取玩家所在联盟
// 	errCode, applies := service.ServMgr.GetAllianceService().GetAllianceApplyList(playerData)
// 	if errCode != msg.ErrCode_SUCC {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceApplyListRsp{
// 			Applies: make([]*msg.AllianceApplyInfo, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 转换为响应消息
// 	rsp := &msg.AllianceApplyListRsp{
// 		Applies: make([]*msg.AllianceApplyInfo, 0, len(applies)),
// 	}
// 	for _, a := range applies {
// 		if a.Status == 0 { // 只返回待处理的申请
// 			rsp.Applies = append(rsp.Applies, &msg.AllianceApplyInfo{
// 				AccountId: uint32(a.PlayerID),
// 				Name:      a.Name,
// 				HeadImg:   a.HeadImg,
// 				HeadFrame: a.HeadFrame,
// 				Power:     a.Power,
// 				ApplyTime: uint32(a.ApplyTime.Unix()),
// 			})
// 		}
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// }

// // HandleAllianceApply 处理联盟申请处理请求
// func HandleAllianceApply(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.HandleAllianceApplyReq)

// 	var targetPlayer *service.PlayerSimpleInfo

// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil || member.Position > 2 { // 只有指挥官和副指挥官可以处理
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}

// 	// 获取联盟信息
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_NOT_IN,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}

// 	// 获取申请记录
// 	apply, err := dao.GetApplication(int64(req.AccountId), member.AllianceID)
// 	if err != nil || apply == nil || apply.AllianceID != member.AllianceID {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_APPLY_NOT_EXIST,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_APPLY_NOT_EXIST)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}

// 	_, targetPlayer = service.ServMgr.GetSocialService().GetPlayerSimpleInfo(req.AccountId)
// 	if targetPlayer == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}

// 	if req.IsAccept {
// 		// 检查联盟是否满员
// 		if alliance.MemberCount >= alliance.MaxMemberCount {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 				AccountId: req.AccountId,
// 				Result:    msg.ErrCode_ALLIANCE_MEMBER_FULL,
// 				IsAccept:  req.IsAccept,
// 			}, packetId, msg.ErrCode_ALLIANCE_MEMBER_FULL)
// 			service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 			return
// 		}

// 		if !doAddMemberToAlliance(packetId, alliance, apply.PlayerID, apply.Name, apply.HeadImg, apply.HeadFrame, apply.Power) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 				AccountId: req.AccountId,
// 				Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 				IsAccept:  req.IsAccept,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		} else {
// 			// 广播联盟基础信息
// 			bMsg := &msg.NotifyAllianceInfo{
// 				Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 			}
// 			service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)
// 		}
// 	}

// 	// 更新申请状态
// 	status := uint8(2) // 拒绝
// 	if req.IsAccept {
// 		status = 1 // 接受
// 	}
// 	if err := dao.UpdateApplication(apply.ID, status); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}
// 	// 删除这个入盟申请 DeleteApplication
// 	if err := dao.DeleteApplication(apply.AllianceID, apply.PlayerID); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 			AccountId: req.AccountId,
// 			Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			IsAccept:  req.IsAccept,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 		return
// 	}

// 	if req.IsAccept {
// 		// 更新成员列表
// 		rsp := &msg.GetAllianceMembersRsp{
// 			Members: make([]*msg.AllianceMemberInfo, 0, 1),
// 			IsNew:   true,
// 		}
// 		rsp.Members = append(rsp.Members, &msg.AllianceMemberInfo{
// 			Data:           service.ToPlayerSimpleInfo(targetPlayer),
// 			Position:       msg.AlliancePosition_MEMBER,
// 			WeeklyActivity: 0,
// 			LastOnlineTime: targetPlayer.LastOnlineTime,
// 		})
// 		tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)

// 		// 对方在线
// 		if other := common.PlayerMgr.FindPlayerData(int64(req.AccountId)); other != nil {
// 			service.InterSendUserMsg(func(arg interface{}, targerPlayer *common.PlayerData) {
// 				// 新成员提示
// 				notifyMsg := &msg.NotifyJoinAlliance{
// 					Name: alliance.Name,
// 				}
// 				tools.SendNotifyMsg(targerPlayer.PlayerAgent, notifyMsg)
// 			}, playerData, other)
// 		}
// 	}

// 	// 获取申请列表
// 	if applies, err := dao.GetApplicationList(alliance.ID); err == nil {
// 		service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
// 			Count: int32(len(applies)),
// 		}, 0, msg.ErrCode_SUCC, 0, 0)
// 	}

// 	service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.HandleAllianceApplyRsp{
// 		AccountId: req.AccountId,
// 		Result:    msg.ErrCode_SUCC,
// 		IsAccept:  req.IsAccept,
// 	}, packetId, msg.ErrCode_SUCC, playerData.AccountInfo.AccountId, 0)

// 	// 发送申请结果
// 	service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceApplyMailId, nil)
// }

// func doAddMemberToAlliance(packetId uint32, alliance *model.Alliance, playerID int64, name string, headImg uint32, headFrame uint32, power uint32) bool {
// 	// 添加新成员
// 	newMember := &model.AllianceMember{
// 		AllianceID: alliance.ID,
// 		PlayerID:   playerID,
// 		Name:       name,
// 		HeadImg:    headImg,
// 		HeadFrame:  headFrame,
// 		Position:   uint8(msg.AlliancePosition_MEMBER), // 普通成员
// 		Power:      power,
// 		JoinTime:   time.Now(),
// 		LastOnline: time.Now(),
// 	}
// 	if err := dao.AddMember(newMember); err != nil {
// 		return false
// 	}

// 	// 更新联盟成员数
// 	alliance.MemberCount++
// 	if err := dao.UpdateAlliance(alliance.ID, bson.M{
// 		"mem_count": alliance.MemberCount,
// 	}); err != nil {
// 		return false
// 	}
// 	return true
// }

// //
// //// HandleAllianceRank 处理联盟排行榜请求
// //func HandleAllianceRank(packetId uint32, args interface{}, playerData *common.PlayerData) {
// //	req := args.(*msg.AllianceRankReq)
// //	// 获取玩家所在联盟
// //	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// //	if err != nil || member == nil {
// //		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceRankRsp{}, packetId, msg.ErrCode_SUCC)
// //		return
// //	}
// //
// //	var rankList interface{}
// //	var err2 error
// //
// //	switch req.RankType {
// //	case publicconst.ALLIANCE_RANK_POWER: // 成员战力榜
// //		rankList, err2 = dao.GetMemberPowerRank(member.AllianceID)
// //	case publicconst.ALLIANCE_RANK_ACTIVE: // 成员活跃榜
// //		rankList, err2 = dao.GetMemberActiveRank(member.AllianceID)
// //	case publicconst.ALLIANCE_RANK_BOSS: // BOSS伤害榜
// //		rankList, err2 = dao.GetBossDamageRank(member.AllianceID)
// //	case publicconst.ALLIANCE_RANK_DAMAGE: // 联盟总伤害榜
// //		rankList, err2 = dao.GetAllianceDamageRank()
// //	}
// //
// //	if err2 != nil {
// //		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceRankRsp{}, packetId, msg.ErrCode_SUCC)
// //		return
// //	}
// //
// //	// 转换为响应消息
// //	rsp := &msg.AllianceRankRsp{
// //		RankType: req.RankType,
// //		// RankList: rankList,
// //	}
// //	fmt.Println("rankList", rankList)
// //	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// //}

// // HandleGetRedPacketList 处理获取红包列表请求
// func HandleGetRedPacketList(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.GetRedPacketListRsp{
// 			RedPackets: make([]*msg.AllianceRedPacket, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 获取联盟红包列表
// 	redPackets, err := dao.GetAllianceRedPackets(playerData.AccountInfo.AccountId)
// 	if err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.GetRedPacketListRsp{
// 			RedPackets: make([]*msg.AllianceRedPacket, 0),
// 		}, packetId, msg.ErrCode_SUCC)
// 		return
// 	}

// 	// 转换为响应消息
// 	rsp := &msg.GetRedPacketListRsp{
// 		RedPackets: make([]*msg.AllianceRedPacket, 0, len(redPackets)),
// 	}
// 	for _, rp := range redPackets {
// 		rsp.RedPackets = append(rsp.RedPackets, &msg.AllianceRedPacket{
// 			Id:          rp.ID,
// 			RedPacketId: uint32(rp.ID),
// 			SenderName:  rp.SenderName,
// 			SenderHead:  rp.HeadImg,
// 			SenderFrame: rp.HeadFrame,
// 		})
// 	}
// 	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// }

// // HandleClaimRedPacket 处理领取红包请求
// func HandleClaimRedPacket(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.ClaimRedPacketReq)

// 	ec, aditems := service.ServMgr.GetAllianceService().ClaimRedPacket(playerData, req.Id)
// 	res := &msg.ClaimRedPacketRsp{
// 		Id:     req.Id,
// 		Result: ec,
// 	}

// 	for _, v := range aditems {
// 		res.RewardItem = append(res.RewardItem, &msg.SimpleItem{
// 			ItemId:  v.Id,
// 			ItemNum: v.Num,
// 		})
// 	}

// 	tools.SendMsg(playerData.PlayerAgent, res, packetId, msg.ErrCode_SUCC)
// }

// // HandleAllianceInfoUpdate 处理修改联盟信息请求
// func HandleAllianceInfoUpdate(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	req := args.(*msg.AllianceInfoUpdateReq)

// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 			UpdateType: req.UpdateType,
// 			Result:     msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}

// 	// 获取目标联盟
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.JoinAllianceRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_EXIST)
// 		return
// 	}

// 	// 检查权限
// 	if member.Position > uint8(msg.AlliancePosition_VICE_LEADER) { // 只有指挥官和副指挥官可以修改
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 			UpdateType: req.UpdateType,
// 			Result:     msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	// 副指挥官只能修改宣言、加入要求、自动加入（位掩码 2/8/16）
// 	if member.Position == uint8(msg.AlliancePosition_VICE_LEADER) {
// 		// 检查是否有非授权修改项（使用异或取反）
// 		if req.UpdateType & ^uint32(2|8|16) != 0 {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 				UpdateType: req.UpdateType,
// 				Result:     msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 			}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 			return
// 		}
// 	}

// 	upV := bson.M{}

// 	// 按位判断各个更新标志位
// 	if req.UpdateType&1 != 0 { // 1: 修改名称
// 		upV["name"] = req.Name
// 		alliance.Name = req.Name
// 		nameLen := uint32(utf8.RuneCountInString(req.Name))
// 		if nameLen < 1 || nameLen > template.GetSystemItemTemplate().AllianceNameLen {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 				UpdateType: req.UpdateType,
// 				Result:     msg.ErrCode_ALLIANCE_NAME_LEN, // 名称
// 			}, packetId, msg.ErrCode_ALLIANCE_NAME_LEN)
// 			return
// 		}
// 		if template.GetForbiddenTemplate().HasForbidden(req.Name) {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 				Result: msg.ErrCode_ALLIANCE_NAME_HAS_FORBIDDEN,
// 			}, packetId, msg.ErrCode_ALLIANCE_NAME_HAS_FORBIDDEN)
// 			return
// 		}
// 		// 检查消耗道具
// 		for _, item := range template.GetSystemItemTemplate().AllianceRenameCostItems {
// 			if !service.ServMgr.GetItemService().EnoughItem(playerData.AccountInfo.AccountId, item.ItemId, item.ItemNum) {
// 				tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 					UpdateType: req.UpdateType,
// 					Result:     msg.ErrCode_NO_ENOUGH_ITEM,
// 				}, packetId, msg.ErrCode_NO_ENOUGH_ITEM)
// 				return
// 			}
// 		}
// 	}
// 	if req.UpdateType&2 != 0 { // 2: 修改宣言
// 		upV["declaration"] = req.Declaration
// 		alliance.Declaration = req.Declaration
// 		nameDeclaration := uint32(utf8.RuneCountInString(req.Declaration))
// 		if nameDeclaration > 0 {
// 			if nameDeclaration > template.GetSystemItemTemplate().AllianceDeclareLen {
// 				tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 					UpdateType: req.UpdateType,
// 					Result:     msg.ErrCode_ALLIANCE_DECLARATION_LEN,
// 				}, packetId, msg.ErrCode_ALLIANCE_DECLARATION_LEN)
// 				return
// 			}
// 			if template.GetForbiddenTemplate().HasForbidden(req.Declaration) {
// 				tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 					Result: msg.ErrCode_ALLIANCE_DECLARATION_HAS_FORBIDDEN,
// 				}, packetId, msg.ErrCode_ALLIANCE_DECLARATION_HAS_FORBIDDEN)
// 				return
// 			}
// 		}
// 	}
// 	if req.UpdateType&4 != 0 { // 4: 修改旗帜
// 		upV["banner"] = req.Banner
// 		alliance.Banner = req.Banner
// 	}
// 	if req.UpdateType&8 != 0 { // 8: 修改加入要求
// 		upV["power_req"] = req.PowerRequired
// 		alliance.PowerRequired = req.PowerRequired
// 		if req.PowerRequired < 0 {
// 			tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 				UpdateType: req.UpdateType,
// 				Result:     msg.ErrCode_ALLIANCE_POWER_TOO_SMALL,
// 			}, packetId, msg.ErrCode_ALLIANCE_POWER_TOO_SMALL)
// 			return
// 		}
// 	}
// 	if req.UpdateType&16 != 0 { // 16: 修改自动加入
// 		upV["auto_join"] = req.AutoJoin
// 		alliance.AutoJoin = req.AutoJoin
// 	}

// 	// 检查是否有有效更新项
// 	if len(upV) == 0 {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 			UpdateType: req.UpdateType,
// 			Result:     msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	err = dao.UpdateAlliance(member.AllianceID, upV)
// 	if err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 			UpdateType: req.UpdateType,
// 			Result:     msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 		return
// 	}

// 	tools.SendMsg(playerData.PlayerAgent, &msg.AllianceInfoUpdateRsp{
// 		UpdateType: req.UpdateType,
// 		Result:     msg.ErrCode_SUCC,
// 	}, packetId, msg.ErrCode_SUCC)

// 	// 消耗盟主道具
// 	if req.UpdateType&1 != 0 { // 1: 修改名称
// 		// 扣除消耗道具
// 		for _, item := range template.GetSystemItemTemplate().AllianceRenameCostItems {
// 			service.ServMgr.GetItemService().CostItem(playerData.AccountInfo.AccountId, item.ItemId, item.ItemNum, publicconst.CreateAllianceCostItem, true)
// 		}
// 	}

// 	// 广播联盟基础信息
// 	bMsg := &msg.NotifyAllianceInfo{
// 		Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 	}
// 	service.ServMgr.GetAllianceService().BroadcastMsgToAlliance(alliance.ID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)
// }

// func HandleAcceptAllianceApply(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	//req := args.(*msg.AcceptAllianceApplyReq)
// 	// 一条条同意, 一共多少条?

// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil || member.Position > 2 { // 只有指挥官和副指挥官可以处理
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AcceptAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		return
// 	}

// 	// 获取联盟信息
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AcceptAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}

// 	// 获取玩家所在联盟
// 	errCode, applies := service.ServMgr.GetAllianceService().GetAllianceApplyList(playerData)
// 	if errCode != msg.ErrCode_SUCC {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.AcceptAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		return
// 	}

// 	for _, apply := range applies {
// 		if acceptMember(playerData, &apply, alliance) != msg.ErrCode_SUCC {
// 		} else {
// 		}
// 	}
// }

// // TODO 完成接受全部联盟申请
// func acceptMember(playerData *common.PlayerData, apply *model.AllianceApplication, alliance *model.Alliance) msg.ErrCode {
// 	//_, targetPlayer := service.ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(apply.PlayerID))
// 	//if targetPlayer == nil {
// 	//	return msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH
// 	//}
// 	//
// 	//// 检查联盟是否满员
// 	//if alliance.MemberCount >= alliance.MaxMemberCount {
// 	//	return msg.ErrCode_ALLIANCE_MEMBER_FULL
// 	//}
// 	//
// 	//if !doAddMemberToAlliance(packetId, alliance, apply.PlayerID, apply.Name, apply.HeadImg, apply.HeadFrame, apply.Power) {
// 	//	tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 	//		AccountId: req.AccountId,
// 	//		Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 	//		IsAccept:  req.IsAccept,
// 	//	}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 	//	return
// 	//} else {
// 	//	// 广播联盟基础信息
// 	//	bMsg := &msg.NotifyAllianceInfo{
// 	//		Alliance: service.ServMgr.GetAllianceService().ToAllianceInfo(alliance),
// 	//	}
// 	//	BroadcastMsgToAlliance(alliance.ID, bMsg, packetId, msg.ErrCode_SUCC, playerData.GetAccountId(), 0)
// 	//}
// 	//
// 	//// 更新申请状态
// 	//status := uint8(2) // 拒绝
// 	//if req.IsAccept {
// 	//	status = 1 // 接受
// 	//}
// 	//if err := dao.UpdateApplication(apply.ID, status); err != nil {
// 	//	tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 	//		AccountId: req.AccountId,
// 	//		Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 	//		IsAccept:  req.IsAccept,
// 	//	}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 	//	service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 	//	return
// 	//}
// 	//// 删除这个入盟申请 DeleteApplication
// 	//if err := dao.DeleteApplication(apply.AllianceID, apply.PlayerID); err != nil {
// 	//	tools.SendMsg(playerData.PlayerAgent, &msg.HandleAllianceApplyRsp{
// 	//		AccountId: req.AccountId,
// 	//		Result:    msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH,
// 	//		IsAccept:  req.IsAccept,
// 	//	}, packetId, msg.ErrCode_ALLIANCE_PERMISSION_NOT_ENOUGH)
// 	//	service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceRefuseMailId, nil)
// 	//	return
// 	//}
// 	//
// 	//if req.IsAccept {
// 	//	// 更新成员列表
// 	//	rsp := &msg.GetAllianceMembersRsp{
// 	//		Members: make([]*msg.AllianceMemberInfo, 0, 1),
// 	//		IsNew:   true,
// 	//	}
// 	//	rsp.Members = append(rsp.Members, &msg.AllianceMemberInfo{
// 	//		Data:           service.ToPlayerSimpleInfo(targetPlayer),
// 	//		Position:       msg.AlliancePosition_MEMBER,
// 	//		WeeklyActivity: 0,
// 	//		LastOnlineTime: targetPlayer.LastOnlineTime,
// 	//	})
// 	//	tools.SendMsg(playerData.PlayerAgent, rsp, packetId, msg.ErrCode_SUCC)
// 	//
// 	//	// 对方在线
// 	//	if other := common.PlayerMgr.FindPlayerData(int64(req.AccountId)); other != nil {
// 	//		service.InterSendUserMsg(func(arg interface{}, targerPlayer *common.PlayerData) {
// 	//			// 新成员提示
// 	//			notifyMsg := &msg.NotifyJoinAlliance{
// 	//				Name: alliance.Name,
// 	//			}
// 	//			tools.SendNotifyMsg(targerPlayer.PlayerAgent, notifyMsg)
// 	//		}, playerData, other)
// 	//	}
// 	//}
// 	//
// 	//// 获取申请列表
// 	//if applies, err := dao.GetApplicationList(alliance.ID); err == nil {
// 	//	BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
// 	//		Count: int32(len(applies)),
// 	//	}, 0, msg.ErrCode_SUCC, 0, 0)
// 	//}
// 	//
// 	//BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.HandleAllianceApplyRsp{
// 	//	AccountId: req.AccountId,
// 	//	Result:    msg.ErrCode_SUCC,
// 	//	IsAccept:  req.IsAccept,
// 	//}, packetId, msg.ErrCode_SUCC, playerData.AccountInfo.AccountId, 0)
// 	//
// 	//// 发送申请结果
// 	//service.ServMgr.GetAllianceService().SendTempMail([]int64{int64(req.AccountId)}, playerData.GenId(), template.GetSystemItemTemplate().AllianceApplyMailId, nil)
// 	return msg.ErrCode_SUCC
// }

// func ClearAllianceApplyReqHandle(packetId uint32, args interface{}, playerData *common.PlayerData) {
// 	//_ := args.(*msg.ClearAllianceApplyReq)
// 	// 获取玩家所在联盟
// 	member, err := dao.GetMember(playerData.AccountInfo.AccountId)
// 	if err != nil || member == nil || member.Position > 2 { // 只有指挥官和副指挥官可以处理
// 		tools.SendMsg(playerData.PlayerAgent, &msg.ClearAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH,
// 		}, packetId, msg.ErrCode_ALLIANCE_POWER_NOT_ENOUGH)
// 		return
// 	}

// 	// 获取联盟信息
// 	alliance, err := dao.GetAlliance(member.AllianceID)
// 	if err != nil || alliance == nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.ClearAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_NOT_IN,
// 		}, packetId, msg.ErrCode_ALLIANCE_NOT_IN)
// 		return
// 	}

// 	// 获取申请记录
// 	if err := dao.RemoveApplication(member.AllianceID); err != nil {
// 		tools.SendMsg(playerData.PlayerAgent, &msg.ClearAllianceApplyRsp{
// 			Result: msg.ErrCode_ALLIANCE_APPLY_NOT_EXIST,
// 		}, packetId, msg.ErrCode_ALLIANCE_APPLY_NOT_EXIST)
// 		return
// 	}

// 	// 获取申请列表
// 	service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.NotifyAllianceApply{
// 		Count: 0,
// 	}, 0, msg.ErrCode_SUCC, 0, 0)

// 	service.ServMgr.GetAllianceService().BroadcastMsgToAllianceForOfficer(alliance.ID, &msg.ClearAllianceApplyRsp{
// 		Result: msg.ErrCode_SUCC,
// 	}, packetId, msg.ErrCode_SUCC, playerData.AccountInfo.AccountId, 0)
// }
