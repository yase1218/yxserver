package service

import (
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tda"
	"kernel/tools"
	"msg"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

func GetGlobalAttrDetail(p *player.Player) (msg.ErrCode, []*msg.Attr) {
	var ret []*msg.Attr
	displayList := template.GetAttrListTemplate().GetDisplayAttr()
	for i := 0; i < len(displayList); i++ {
		attrConfig := displayList[i]
		if data, ok := p.UserData.BaseInfo.Attrs[attrConfig.Id]; ok {
			temp := &msg.Attr{
				Id:        attrConfig.Id,
				CalcValue: data.FinalValue,
				Value:     data.InitValue + data.LevelValue + data.Add,
			}

			if len(attrConfig.DisplayPara) > 0 {
				var finalValue float32
				for k := 0; k < len(attrConfig.DisplayPara); k++ {
					//	temp2 := template.GetAttrListTemplate().GetAttr(attrConfig.DisplayPara[k])
					if data2, ok2 := p.UserData.BaseInfo.Attrs[attrConfig.DisplayPara[k]]; ok2 {
						finalValue += data2.CalcFinalValue()
					}
				}
				temp.Value = finalValue
			}
			ret = append(ret, temp)
		}
	}
	return msg.ErrCode_SUCC, ret
}

func SetSupportId(p *player.Player, supportId []uint32) msg.ErrCode {
	if len(supportId) == 0 {
		return msg.ErrCode_INVALID_DATA
	}
	same := true
	for i := 0; i < len(supportId); i++ {
		if supportId[i] > 0 && supportId[i] == p.UserData.BaseInfo.ShipId {
			return msg.ErrCode_INVALID_DATA
		}
	}
	if len(p.UserData.BaseInfo.SupportId) > 0 {
		for i := 0; i < len(supportId); i++ {
			if supportId[i] != p.UserData.BaseInfo.SupportId[i] {
				same = false
				break
			}
		}
	} else {
		same = false
	}
	if same {
		return msg.ErrCode_INVALID_DATA
	}

	for i := 0; i < len(supportId); i++ {
		if supportId[i] > 0 {
			if ship := getShip(p, supportId[i]); ship == nil {
				return msg.ErrCode_SHIP_NOT_EXIST
			}
		}
	}

	if len(supportId) > int(template.GetSystemItemTemplate().TeamSupportShipNum) {
		return msg.ErrCode_INVALID_DATA
	}

	if p.UserData.BaseInfo.SupportId == nil {
		p.UserData.BaseInfo.SupportId = make([]uint32, len(supportId))
	}
	copy(p.UserData.BaseInfo.SupportId, supportId)
	p.SaveBaseInfo()
	GlobalAttrChange(p, true)

	return msg.ErrCode_SUCC
}

// SetShipId 设置出战机甲
func SetShipId(p *player.Player, shipId uint32) msg.ErrCode {
	ship := getShip(p, shipId)
	if ship == nil {
		return msg.ErrCode_SHIP_NOT_EXIST
	}

	if p.UserData.BaseInfo.ShipId == shipId {
		return msg.ErrCode_INVALID_DATA
	}

	p.UserData.BaseInfo.ShipId = shipId
	p.SaveBaseInfo()

	GlobalAttrChange(p, true)

	//common.PlayerMgr.UpdatePlayerBasic(&model.AccBasic{AccountId: p.GetAccountId(), ShipId: shipId})

	UpdateTask(p, true, publicconst.TASK_COND_SHIP_SWITCH, 1)
	return msg.ErrCode_SUCC
}

func GmDayResetAccout(p *player.Player) {
	if p.UserData.BaseInfo == nil {
		return
	}

	p.UserData.BaseInfo.ActiveDay = 0
	p.SaveBaseInfo()
}

// 升级机甲gm命令
func GmUpgrade(p *player.Player, addLv uint32) msg.ErrCode {

	curLevel := p.UserData.Level
	nextLevel := template.GetShipLevelTempalte().GetLevel(curLevel + addLv)
	// 满级了
	if nextLevel == nil {
		return msg.ErrCode_LEVEL_FULL
	}

	// 等级变化
	addAllShipLevelAttr(p, p.UserData.Level, nextLevel.Data.ShipLevel)

	// 计算全局属性
	GlobalAttrChange(p, true)

	p.UserData.Level = nextLevel.Data.ShipLevel
	p.SaveLevel()

	event.EventMgr.PublishEvent(event.NewLevelChangeEvent(p,
		curLevel, nextLevel.Data.ShipLevel, ListenLevelChangeEvent))

	return msg.ErrCode_SUCC
}

// Upgrade 升级
func Upgrade(p *player.Player) (msg.ErrCode, uint32) {
	curLevel := p.UserData.Level
	nextLevel := template.GetShipLevelTempalte().GetLevel(curLevel + 1)
	// 满级了
	if nextLevel == nil {
		return msg.ErrCode_LEVEL_FULL, 0
	}

	for i := 0; i < len(nextLevel.CostItem); i++ {
		if !EnoughItem(p.GetUserId(),
			nextLevel.CostItem[i].ItemId, nextLevel.CostItem[i].ItemNum) {
			return msg.ErrCode_NO_ENOUGH_ITEM, 0
		}
	}

	// 扣除道具
	var notifyClientItems []uint32
	tdaItems := make([]*tda.Item, 0, len(nextLevel.CostItem))
	for i := 0; i < len(nextLevel.CostItem); i++ {
		CostItem(p.GetUserId(), nextLevel.CostItem[i].ItemId, nextLevel.CostItem[i].ItemNum, publicconst.UpgradeCostItem,
			false)
		notifyClientItems = append(notifyClientItems, nextLevel.CostItem[i].ItemId)
		tdaItems = append(tdaItems, &tda.Item{ItemId: strconv.Itoa(int(nextLevel.CostItem[i].ItemId)), ItemNum: nextLevel.CostItem[i].ItemNum})
	}

	// 通知客户端
	updateClientItemsChange(p.GetUserId(), notifyClientItems)

	// 等级变化
	addAllShipLevelAttr(p, p.UserData.Level, nextLevel.Data.ShipLevel)

	// 计算全局属性
	GlobalAttrChange(p, true)

	p.UserData.Level = nextLevel.Data.ShipLevel
	p.SaveLevel()

	// // tda update level
	// tdaData := &tda.CommonUser{
	// 	Current_level: strconv.Itoa(int(p.AccountInfo.Level)),
	// }
	// tda.TdaUpdateCommonUser(p.TdaCommonAttr.AccountId, p.TdaCommonAttr.DistinctId, tdaData)

	event.EventMgr.PublishEvent(event.NewLevelChangeEvent(p,
		curLevel, nextLevel.Data.ShipLevel, ListenLevelChangeEvent))

	// // tda
	// tda.TdaKuluUpgrade(p.ChannelId, p.TdaCommonAttr, p.AccountInfo.Level, tdaItems)

	return msg.ErrCode_SUCC, p.UserData.Level
}

// UpdateNick 更新昵称
func UpdateNick(p *player.Player, nick string) msg.ErrCode {
	old := p.UserData.Nick
	p.UserData.Nick = nick
	p.UpdateNickTime = tools.GetCurTime()
	p.UserData.BaseInfo.IsEditName = true

	p.SaveNick()
	p.SaveBaseInfo()
	player.UpdateNick(old, nick)
	return msg.ErrCode_SUCC
}

// UpdateNickCheck 更新昵称检擦
func UpdateNickCheck(p *player.Player, nick string) msg.ErrCode {
	if len(nick) == 0 {
		return msg.ErrCode_NICK_EMPTY
	}

	nick = strings.Trim(nick, " ")
	if p.UserData.Nick == nick {
		return msg.ErrCode_NICK_SAME
	}

	if utf8.RuneCountInString(nick) > int(template.GetSystemItemTemplate().NickLen) {
		return msg.ErrCode_NICK_TO_LONG
	}

	if strings.HasPrefix(nick, "kulu_") {
		return msg.ErrCode_NICK_HAS_FORBIDDEN
	}
	//if template.GetForbiddenTemplate().HasForbidden(nick) {
	//	return msg.ErrCode_NICK_HAS_FORBIDDEN
	//}

	curTime := tools.GetCurTime()
	if p.UpdateNickTime > 0 && (curTime-p.UpdateNickTime) < uint32(publicconst.UPDATE_NICK_CD) {
		return msg.ErrCode_NICK_IN_CD
	}

	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
		key  = RdbUserNameKey()
	)

	old := p.UserData.Nick
	_, err := rc.SRem(rCtx, key, old).Result()
	if err != nil {
		log.Error("SRem remove nick error", zap.Error(err), zap.String("old nick", old))
		return msg.ErrCode_NICK_EXIST
	}

	_, err1 := rc.SAdd(rCtx, key, nick).Result()
	if err1 != nil {
		log.Error("SAdd player name err", zap.Error(err), zap.String("nick", nick))
		rc.SAdd(rCtx, key, old)
		return msg.ErrCode_Name_Not_Vaild
	}
	return msg.ErrCode_SUCC
}

func OnHeart(p *player.Player) {
	//a.ResetData(playerData)
	p.UpdateTime = tools.GetCurTime()

	//recoveryAp(p)
}

func PlayerLogout(p *player.Player) {
	accountId := p.GetOpenId()
	p.EndFightType = msg.EndFightType_End_Fight_Disconnect
	LeaveFight(p)

	p.State = publicconst.Offline
	OnLogout(p, time.Now())
	log.Info("PlayerLogout", zap.String("accountId", accountId))
}

// NotifyAccountChange 通知客户端账号变化
func NotifyAccountChange(p *player.Player) {
	p.SendNotify(&msg.NotifyAccountChange{
		Info: ToProtocolAccountInfo(p.UserData),
	})
}

// AddComboSkill 添加组合技能
func AddComboSkill(p *player.Player, skills []uint32) {
	if len(skills) == 0 {
		return
	}
	for i := 0; i < len(skills); i++ {
		if !tools.ListContain(p.UserData.BaseInfo.ComboSkill, skills[i]) {
			p.UserData.BaseInfo.ComboSkill = append(p.UserData.BaseInfo.ComboSkill, skills[i])
		}
	}
	//dao.AccountDao.UpdateComboSkill(p.GetAccountId(), p.UserData.BaseInfo.ComboSkill)
	p.SaveBaseInfo()
}

// ResetData 通知客户端重置数据
func ResetData(p *player.Player, now time.Time) {
	curTime := uint32(now.Unix())
	var resetDatas []msg.ResetDataType
	if curTime >= p.UserData.BaseInfo.NextDailyRefreshTime {
		p.UserData.BaseInfo.NextDailyRefreshTime = tools.GetDailyRefreshTime()
		GetTasksByType(p, publicconst.DAILY_TASK)
		resetDatas = append(resetDatas, msg.ResetDataType_Reset_Daily_Data)
		resetDatas = append(resetDatas, msg.ResetDataType_Reset_Explore_GetFreeCard)
	}
	if curTime >= p.UserData.BaseInfo.NextWeeklyRefreshTime {
		p.UserData.BaseInfo.NextWeeklyRefreshTime = tools.GetWeeklyRefreshTime(template.GetSystemItemTemplate().RefreshHour)
		GetTasksByType(p, publicconst.WEEKLY_TASK)
		resetDatas = append(resetDatas, msg.ResetDataType_Reset_Weekly_Data)
	}

	for i := 0; i < len(resetDatas); i++ {
		// 重置日常
		if resetDatas[i] == msg.ResetDataType_Reset_Daily_Data {
			//ServMgr.GetMailService().DeleteOverTimeMailWithLock(p)
			RefreshPlayMethod(p, true, false)
			ResetCardPool(p, true)
			ResetAdData(p)
			ResetDailyAp(p, true)
			//ResetPetData(p, true)
		} else if resetDatas[i] == msg.ResetDataType_Reset_Weekly_Data {

		} else if resetDatas[i] == msg.ResetDataType_Reset_Explore_GetFreeCard {
			// 	ServMgr.GetExploreService().ResetData(p)
		}
	}

	if len(resetDatas) > 0 {
		p.SaveBaseInfo()
		ntf := &msg.NotifyResetData{}
		ntf.DataType = resetDatas
		ntf.Time = ToProtocolServerTime(now)
		p.SendNotify(ntf)
	}

	if curTime >= p.UserData.BaseInfo.NextRefreshActivityTime {
		p.UserData.BaseInfo.NextRefreshActivityTime = tools.GetHourRefreshTime(0)
		RefreshActivity(p, true)
		updateActivitiesData(p)
	}
}

func NotifyApRecovery(p *player.Player) {
	p.SendNotify(&msg.NotifyApRecover{
		RecoverStartTime: p.UserData.BaseInfo.ApData.RecoverStartTime,
	})
}

func ToProtocolAccountInfo(u *model.UserData) *msg.AccountInfo {
	return &msg.AccountInfo{
		AccountId:         int64(u.UserId),
		Nick:              u.Nick,
		Level:             u.Level,
		HeadImg:           u.HeadImg,
		ShipId:            u.BaseInfo.ShipId,
		SupportId:         u.BaseInfo.SupportId,
		PokerSlotCount:    u.BaseInfo.PokerSlotCount,
		DisplayId:         uint32(u.UserId),
		VideoFlag:         u.BaseInfo.VideoFlag,
		QuestionIds:       u.BaseInfo.QuestionIds,
		RewardQuestionIds: u.BaseInfo.RewardQuestionIds,
		Combat:            u.BaseInfo.Combat,
		Attrs:             ToProtocolGlobalAttr(u.BaseInfo.Attrs),
		HeadFrame:         u.HeadFrame,
		Title:             u.Title,
		ForbiddenChat:     u.BaseInfo.ForbiddenChat,
		IsEditName:        u.BaseInfo.IsEditName,
		Uid:               u.UserId,
	}
}

// ToProtocolApInfo 体力相关数据
func ToProtocolApInfo(accountInfo *model.Account) *msg.ApData {
	return &msg.ApData{
		BuyTimes:         accountInfo.ApData.BuyTimes,
		RecoverStartTime: accountInfo.ApData.RecoverStartTime,
	}
}

// 踢出玩家
func KickOutPlayer(p *player.Player) {
	// todo 通知fight暂停pve战斗

	// if p.FightType != msg.BattleType_Battle_None {
	// 	if p.FightType != msg.BattleType_Battle_EquipStage && p.FightType != msg.BattleType_Battle_Peak {
	// 		PauseFight(p)
	// 	} else {
	// 		// if !EquipStageLeave(p) {
	// 		// 	EquipStageTeamLeave(0, p, true)
	// 		// }
	// 		// OnPFLogOut(p)
	// 	}
	// }
	kickout := &msg.NotifyKickOut{
		Result: msg.ErrCode_SUCC,
	}
	p.SendNotify(kickout)

	p.State = publicconst.Offline
	accountId := p.GetOpenId()
	p.EndFightType = msg.EndFightType_End_Fight_Kick
	LeaveFight(p)

	log.Info("KickOutPlayer", zap.String("accountId", accountId))

	OnLogout(p, time.Now())

}
