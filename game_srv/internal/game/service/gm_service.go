package service

import (
	"encoding/json"
	"fmt"
	"gameserver/internal/game/player"
	"msg"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"

	"gameserver/internal/config"
	"gameserver/internal/game/event"
	"gameserver/internal/game/model"
	"gameserver/internal/publicconst"
)

type CommandItem struct {
	ItemId uint32 `json:"itemId"`
	Num    uint32 `json:"num"`
}

type CommandMail struct {
	TemplateId uint32              `json:"templateId"`
	Title      string              `json:"title"`
	Content    string              `json:"content"`
	Attachment []*model.SimpleItem `json:"attachment"`
	Para       string              `json:"para"`
}

type Equip struct {
	EquipId  uint32 `json:"equipId"`
	EquipNum uint32 `json:"num"`
}

type ShopItem struct {
	ShopItemId uint32 `json:"shopItemId"`
	Num        uint32 `json:"num"`
}

type ChargeItem struct {
	Id  uint32 `json:"id"`
	Num uint32 `json:"num"`
}

type CheckBattle struct {
	Check uint32 `json:"check"`
}

type Novice struct {
	Id    uint32 `json:"id"`
	Value uint32 `json:"value"`
}

type UpgradeStarShipStruct struct {
	Id    uint32 `json:"id"`
	AddLv uint32 `json:"addLv"`
}

type WeaponUpgrade struct {
	WeaponId uint32 `json:"weaponId"`
	Lv       uint32 `json:"lv"`
}

func ProcessCommand(command msg.GMCommandId, content string, p *player.Player) msg.ErrCode {
	if command == msg.GMCommandId_AdCallBack {
		AdCallBack(p)
		return msg.ErrCode_SUCC
	}

	if !config.Conf.Gm {
		return msg.ErrCode_SYSTEM_ERROR
	}

	log.Debug("ProcessCommand: ", zap.Any("cmd", command), zap.Any("conten", content), zap.Any("accountId", p.GetUserId()))
	params := strings.Split(content, ",")

	//if GameInfo.Gm == 0 {
	//	return msg.ErrCode_SUCC
	//}

	//	ServMgr.GetAttrService().GlobalAttrChange(playerData, true)
	switch command {
	case msg.GMCommandId_AddItem:
		return gm_addItem(content, p)
	case msg.GMCommandId_AddEquip:
		return gm_addEquip(content, p)
	case msg.GMCommandId_Charge:
		return gm_charge(content, p)
	case msg.GMCommandId_SetNovice:
		return gm_setNovice(content, p)
	case msg.GMCommandId_CheckBattle:
		return gm_checkBattle(content, p)
	case msg.GMCommandId_PassAllMainMission:
		return gm_passAllMainMission(content, p)
	//case msg.GMCommandId_AdCallBack:
	//	ServMgr.GetAdService().AdCallBack(playerData)
	//	return msg.ErrCode_SUCC
	case msg.GMCommandId_RefreshRank:
		// TODO 排行榜
		//RefreshRank(p)
		return msg.ErrCode_SUCC
	case msg.GMCommandId_AddAllianceExp:
		// TODO 公会联盟
		// // 检查玩家是否已有联盟
		// member, err := dao.GetMember(p.AccountInfo.AccountId)
		// if err == nil && member != nil {
		// 	ServMgr.GetAllianceService().AddAllianceExp(member.AllianceID, 1000)
		// }
		return msg.ErrCode_SUCC

	case msg.GMCommandId_DayReset:
		gm_day_reset(p)

		return msg.ErrCode_SUCC
	case msg.GMCommandId_LevelUpgrade:
		value, _ := strconv.ParseUint(content, 10, 64)
		addLv := uint32(value)
		if addLv <= 0 {
			return msg.ErrCode_INVALID_DATA
		}
		return GmUpgrade(p, addLv)
	case msg.GMCommandId_ShipStarUpgrade:
		if len(params) != 2 {
			return msg.ErrCode_INVALID_DATA
		}

		id, _ := strconv.ParseUint(params[0], 10, 64)
		addlv, _ := strconv.ParseUint(params[1], 10, 64)

		return GmUpgradeStarShip(p, uint32(id), uint32(addlv))
	case msg.GMCommandId_EquipPosUpgrade:

		if len(params) != 2 {
			return msg.ErrCode_INVALID_DATA
		}

		pos, _ := strconv.ParseUint(params[0], 10, 64)
		lv, _ := strconv.ParseUint(params[1], 10, 64)

		ec := GmUpgradeEquipPos(p, uint32(pos), uint32(lv))
		return ec
	case msg.GMCommandId_WeaponUpgrade:
		if len(params) != 2 {
			return msg.ErrCode_INVALID_DATA
		}

		weaponId, _ := strconv.ParseUint(params[0], 10, 64)
		lv, _ := strconv.ParseUint(params[1], 10, 64)

		ec := GmUpgradeWeapon(p, uint32(weaponId), uint32(lv))
		return ec
	case msg.GMCommandId_PassChallenge:

		return passChallenge(content, p)

	case msg.GMCommandId_ResetSystemTime:
		if !config.Conf.GmTime {
			return msg.ErrCode_INVALID_DATA
		}
		// ntpServers := []string{
		// 	"ntp.ntsc.ac.cn",   // 中国国家授时中心主服务器
		// 	"cn.pool.ntp.org",  // 公共NTP服务器池
		// 	"time1.aliyun.com", // 阿里云NTP服务器
		// }
		// ntpTime, err := getTimeFromNTP(ntpServers)
		// if err != nil {
		// 	log.Error("无法从任何NTP服务器获取时间: ", zap.Error(err))
		// }

		// // 将时间转换为北京时间(UTC+8)
		// beijingTime := ntpTime.In(time.FixedZone("CST", 8*3600))
		// return g.setSystemTime(beijingTime)
	case msg.GMCommandId_NextDaySystemTime:
		// if !config.Conf.GmTime {
		// 	return msg.ErrCode_INVALID_DATA
		// }
		// return g.setSystemTime(time.Now().AddDate(0, 0, 1))
	case msg.GMCommandId_FinishOpenActivityTask:
		return GmFinishTodayOpenTask(p)
	case msg.GMCommandId_SetPeakFightCup:
		// TODO 重构
		value, _ := strconv.ParseUint(content, 10, 64)
		return GmSetPeakFightCup(p, uint32(value))
	case msg.GMCommandId_CopyAccount:
		// if len(params) != 3 {
		// 	return msg.ErrCode_INVALID_DATA
		// }

		templateId, _ := strconv.ParseUint(params[0], 10, 64)
		// prefix := params[1]
		// num, _ := strconv.ParseUint(params[2], 10, 64)
		// return copyAccount(int64(templateId), prefix, int(num))
		ChargeCallBack(p, int(templateId), 0)
	}

	return msg.ErrCode_SYSTEM_ERROR
}

func gm_day_reset(playerData *player.Player) {
	//ServMgr.GetExploreService().GmResetFreeCard(playerData)
	GmRefreshPlayMethod(playerData)
	GmDayResetAccout(playerData)
	GmDayRefreshShopItems(playerData)
	GmDayResetActivity(playerData)
	GmResetDataCardPool(playerData)
	//ServMgr.GetArenaService().GmResetData(playerData)
	ResetContractAll(playerData, true, true)
	//ServMgr.GetPeakFightService().Reset(playerData, true)
	DayResetDesert(playerData, true)
}

func getTimeFromNTP(servers []string) (time.Time, error) {
	var lastError error

	for _, server := range servers {
		fmt.Printf("尝试连接到NTP服务器: %s...\n", server)

		// 设置选项，包括超时和详细模式
		options := ntp.QueryOptions{
			Timeout: 5 * time.Second,
			// 可以添加更多选项，如Version、Tolerance等
		}

		// 获取带统计信息的时间
		timeInfo, err := ntp.QueryWithOptions(server, options)
		if err != nil {
			lastError = err
			log.Error("获取时间失败 ")
			continue
		}

		// 计算校正后的时间
		ntpTime := timeInfo.Time

		// 验证时间精度
		if timeInfo.ClockOffset > 100*time.Millisecond {
			log.Error("警告: 时钟偏移较大 ")
		}

		return ntpTime, nil
	}

	return time.Time{}, fmt.Errorf("所有NTP服务器均失败: %v", lastError)
}

func gm_addItem(content string, p *player.Player) msg.ErrCode {
	var items []*CommandItem
	// json.Unmarshal([]byte(content), &items)

	itemStrs := strings.Split(content, "|")
	for _, v := range itemStrs {
		params := strings.Split(v, ",")
		if len(params) != 2 {
			continue
		}

		items = append(items, &CommandItem{
			ItemId: utils.StrToUInt32(params[0]),
			Num:    utils.StrToUInt32(params[1]),
		})
	}

	if len(items) == 0 {
		return msg.ErrCode_SUCC
	}

	var notifyItems []uint32
	for i := 0; i < len(items); i++ {
		if items[i].Num > 0 {
			addItems := AddItem(p.GetUserId(), items[i].ItemId, int32(items[i].Num), publicconst.GMAddItem, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		} else {
			notifyItems = append(notifyItems, items[i].ItemId)
			temp := uint32(-(items[i].Num))
			totalNum := GetItemNum(p.GetUserId(), items[i].ItemId)
			if uint64(temp) > totalNum {
				temp = uint32(totalNum)
			}
			CostItem(p.GetUserId(), items[i].ItemId, temp, publicconst.GMCostItem, true)
		}
	}

	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(p, notifyItems, ListenNotifyClientItemEventEvent))
	return msg.ErrCode_SUCC
}

func gm_addEquip(content string, p *player.Player) msg.ErrCode {
	var equips []*Equip
	json.Unmarshal([]byte(content), &equips)

	for i := 0; i < len(equips); i++ {
		addEquip(p, equips[i].EquipId, equips[i].EquipNum, false)
	}
	if len(equips) > 0 {
		p.SaveEquip()
	}
	return msg.ErrCode_SUCC
}

func gm_charge(content string, playerData *player.Player) msg.ErrCode {
	var item ChargeItem
	json.Unmarshal([]byte(content), &item)

	// TODO 充值
	// ChargeCallBack(playerData, int(item.Id), int(item.Num))
	return msg.ErrCode_SUCC
}

func gm_checkBattle(content string, p *player.Player) msg.ErrCode {
	var battle CheckBattle
	json.Unmarshal([]byte(content), &battle)

	p.CheckBattle = battle.Check
	return msg.ErrCode_SUCC
}

// passAllMainMission 通关所有主线
func gm_passAllMainMission(content string, playerData *player.Player) msg.ErrCode {
	var (
		stageCfg *template.JMission

		playerStageId  = playerData.UserData.StageInfo.MissionId
		passEndStageId = utils.StrToInt(content)
	)

	if playerStageId == 0 {
		stageCfg = template.GetMissionTemplate().FirstMission
	} else {
		stageCfg = template.GetMissionTemplate().GetMission(int(playerStageId))
		if stageCfg == nil || stageCfg.NextId == 0 {
			//log.Error("stage cfg nil", zap.Uint32("playerStageId", playerStageId))
			return msg.ErrCode_SUCC
		}
	}

	// 先加体力
	{
		var notifyItems []uint32
		addItems := AddItem(playerData.GetUserId(), uint32(publicconst.ITEM_CODE_AP), 1000, publicconst.GMAddItem, false)
		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(playerData, notifyItems, ListenNotifyClientItemEventEvent))
	}

	ntf := &msg.FsFightResultNtf{
		AccountId:    playerData.GetUserId(),
		Victory:      true,
		KillCnt:      200,
		FightTime:    350000,
		Hp:           900000,
		Chip:         100,
		PokerTypes:   nil,
		KillMonsters: nil,
	}

	for {
		ntf.StageId = uint32(stageCfg.Id)

		MainFightBeforeCheck(playerData, stageCfg)
		MainFightAfter(playerData, ntf)

		//if stageCfg.Id == 10007 ||
		if (passEndStageId != 0 && stageCfg.Id >= passEndStageId) ||
			stageCfg.NextId == 0 {
			break
		}

		nextId := stageCfg.NextId
		stageCfg = template.GetMissionTemplate().GetMission(nextId)
		if stageCfg == nil {
			log.Error("stage cfg nil", zap.Int("stageId", nextId))
			break
		}
	}

	return msg.ErrCode_SUCC
}

// passAllMainMission_bak 旧通关所有主线
func passAllMainMission_bak(content string, p *player.Player) msg.ErrCode {
	id := utils.StrToInt(content)

	tempConfig := template.GetMissionTemplate().GetMission(int(p.UserData.StageInfo.MissionId))
	var startMission *template.JMission
	if tempConfig == nil {
		startMission = template.GetMissionTemplate().FirstMission
	} else {
		if tempConfig.NextId == 0 {
			return msg.ErrCode_SUCC
		}
	}

	var notifyItems []uint32
	addItems := AddItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), 1000, publicconst.GMAddItem, false)
	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(p, notifyItems, ListenNotifyClientItemEventEvent))

	// 开始战斗
	p.CheckBattle = 0
	for {
		if tempConfig != nil {
			startMission = template.GetMissionTemplate().GetMission(tempConfig.NextId)
		}

		StartBattle(p, startMission.Id)
		EndBattle(p, &msg.BattleResult{})

		time.Sleep(time.Millisecond * 100)
		tempConfig = startMission
		if tempConfig.Id == 10007 {
			break
		}

		if id != 0 && tempConfig.Id >= id {
			break
		}

		if tempConfig.NextId == 0 {
			break
		}
	}
	p.CheckBattle = 1
	return msg.ErrCode_SUCC
}

// passAllMainMission 通关精英关卡
func passChallenge(content string, p *player.Player) msg.ErrCode {
	if content == "" {
		return msg.ErrCode_SUCC
	}
	targetMissionId := utils.StrToInt(content)

	if template.GetMissionTemplate().FirstMission == nil {
		return msg.ErrCode_SUCC
	}

	var notifyItems []uint32
	addItems := AddItem(p.GetUserId(), uint32(publicconst.ITEM_CODE_AP), 1000, publicconst.GMAddItem, false)
	notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	event.EventMgr.PublishEvent(event.NewNotifyClientItemEvent(p, notifyItems, ListenNotifyClientItemEventEvent))

	// 开始战斗
	p.CheckBattle = 0

	missionId := template.GetMissionTemplate().FirstJMission.Id
	safeCount := 100
	for {
		safeCount--
		if safeCount <= 0 {
			log.Error("passChallenge dead loop")
			break
		}
		tempConfig := template.GetMissionTemplate().GetMission(missionId)

		if tempConfig == nil {
			break
		}

		if findMission(p, tempConfig.Id, false) != nil {
			continue
		}

		StartBattle(p, missionId)
		ChallengeEndBattle(p, &msg.BattleResult{})
		time.Sleep(time.Millisecond * 100)

		if missionId >= targetMissionId {
			break
		}

		missionId++
	}
	p.CheckBattle = 1
	return msg.ErrCode_SUCC
}

func addMail(content string, playerData *player.Player) msg.ErrCode {
	var mail CommandMail
	json.Unmarshal([]byte(content), &mail)

	//	ServMgr.GetMailService().AddMail(playerData, mail.TemplateId, mail.Title, mail.Content, mail.Para, mail.Attachment)
	return msg.ErrCode_SUCC
}

func gm_setNovice(content string, p *player.Player) msg.ErrCode {
	var novice Novice
	json.Unmarshal([]byte(content), &novice)

	SetGuideInfo(p, novice.Id, novice.Value)
	return msg.ErrCode_SUCC
}

func setSystemTime(t time.Time) msg.ErrCode {
	// 使用 date 命令设置时间（格式："2023-01-01 12:00:00"）

	cmd := exec.Command("date", "-s", t.Format("2006-01-02 15:04:05"))

	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("GmService setSystemTime err, %s", zap.String("err", string(output)))
		return msg.ErrCode_SYSTEM_ERROR
	}

	// 同步硬件时钟（可选，但建议执行）
	// syncCmd := exec.Command("hwclock", "--systohc")
	// if out, err := syncCmd.CombinedOutput(); err != nil {
	// log.Error("GmService setSystemTime err, %s", zap.String("err", string(out)))
	// 	return msg.ErrCode_SYSTEM_ERROR
	// }

	return msg.ErrCode_SUCC
}

func copyAccount(templateId int64, prefix string, num int) msg.ErrCode {
	// globalDb := db.GetGlobalClient().Database(config.Conf.GlobalMongo.DB)

	// dbAccount := db_global.GetAccountModel()
	// templateAccountInfo, err := dbAccount.GetOne(uint64(templateId))
	// if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	// 	log.Error("get account info err", zap.Error(err), zap.Any("templateId", templateId))
	// 	return msg.ErrCode_ERR_NONE
	// }

	// msdkChannel := templateAccountInfo.Account.Msdk.Platform

	// ec, templateItem := dao.ItemDao.LoadItems(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadItems failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateTask := dao.TaskDao.LoadTasks(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadTasks failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateMission := dao.MissionDao.LoadMissions(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadMissions failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateShip := dao.ShipDao.LoadShips(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadShips failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateRole := dao.RoleDao.LoadRoles(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadRoles failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateTeams := dao.TeamDao.LoadTeams(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadTeams failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateEquips := dao.EquipDao.LoadEquips(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadEquips failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// ec, templateWeapon := dao.WeaponDao.LoadWeapon(templateId)
	// if ec != msg.ErrCode_SUCC {
	// 	log.Error("LoadWeapon failed", zap.Any("ec", ec))
	// 	return ec
	// }

	// for i := 1; i <= num; i++ {

	// 	newAccId, err := db_global.GenAccountIdSeq()
	// 	if err != nil {
	// 		log.Error("gen account id err", zap.Error(err))
	// 		continue
	// 	}

	// 	openId := fmt.Sprintf("%s_%d", prefix, i)
	// 	accountInfo := &db_global.Account{
	// 		ID:       newAccId,
	// 		Account:  db_global.NewAccMsdk(openId, templateAccountInfo.Account.Msdk.Platform),
	// 		SerId:    templateAccountInfo.SerId,
	// 		Ip:       templateAccountInfo.Ip,
	// 		CreateAt: time.Now(),

	// 		DeviceId:       templateAccountInfo.DeviceId,
	// 		Os:             templateAccountInfo.Os,
	// 		OsVersion:      templateAccountInfo.OsVersion,
	// 		AppVersion:     templateAccountInfo.AppVersion,
	// 		Manufacturer:   templateAccountInfo.Manufacturer,
	// 		DeviceModel:    templateAccountInfo.DeviceModel,
	// 		ScreenHeight:   templateAccountInfo.ScreenHeight,
	// 		ScreenWidth:    templateAccountInfo.ScreenWidth,
	// 		Ram:            templateAccountInfo.Ram,
	// 		Disk:           templateAccountInfo.Disk,
	// 		NetworkType:    templateAccountInfo.NetworkType,
	// 		Carrier:        templateAccountInfo.Carrier,
	// 		Country:        templateAccountInfo.Country,
	// 		CountryCode:    templateAccountInfo.CountryCode,
	// 		SystemLanguage: templateAccountInfo.SystemLanguage,
	// 		ChannelId:      templateAccountInfo.ChannelId,
	// 	}

	// 	channelInfo := &db_global.AccountChannelInfo{
	// 		Channel: kenum.Account_Type_Msdk,
	// 		AccountInfo: &db_global.AccMsdk{
	// 			OpenId:   openId,
	// 			Platform: msdkChannel,
	// 		},
	// 	}

	// 	if err = dbAccount.NewUserUnique(channelInfo, accountInfo); err != nil {
	// 		log.Error("gm accountRegisterUnique err", zap.Error(err), zap.Any("channelInfo", channelInfo), zap.Any("accountInfo", accountInfo))
	// 		continue
	// 	}

	// 	acctIdStr := fmt.Sprintf("%v", newAccId)
	// 	dao.AccountDao.AddAccount(acctIdStr, int64(newAccId), templateAccountInfo.ChannelId, acctIdStr)

	// 	dao.ItemDao.LoadItems(int64(newAccId))
	// 	dao.ItemDao.AddItems(int64(newAccId), templateItem.Items)

	// 	dao.TaskDao.LoadTasks(int64(newAccId))
	// 	dao.TaskDao.AddTasks(int64(newAccId), publicconst.MAIN_TASK, templateTask.MainTasks)
	// 	dao.TaskDao.AddTasks(int64(newAccId), publicconst.ACHIEVE_TASK, templateTask.AchieveTasks)
	// 	dao.TaskDao.AddTasks(int64(newAccId), publicconst.DAILY_TASK, templateTask.DailyTasks)
	// 	dao.TaskDao.AddTasks(int64(newAccId), publicconst.WEEKLY_TASK, templateTask.WeeklyTasks)

	// 	dao.MissionDao.LoadMissions(int64(newAccId))
	// 	missionIdMax := 0
	// 	for _, v := range templateMission.Missions {
	// 		dao.MissionDao.AddMission(int64(newAccId), v, true)

	// 		if missionIdMax < v.MissionId {
	// 			missionIdMax = v.MissionId
	// 		}
	// 	}
	// 	dao.AccountDao.UpdateMissionId(int64(newAccId), missionIdMax)
	// 	for _, v := range templateMission.Challenges {
	// 		dao.MissionDao.AddMission(int64(newAccId), v, false)
	// 	}

	// 	dao.ShipDao.LoadShips(int64(newAccId))
	// 	for _, v := range templateShip.Ships {
	// 		dao.ShipDao.AddShip(int64(newAccId), v)
	// 	}

	// 	dao.RoleDao.LoadRoles(int64(newAccId))
	// 	for _, v := range templateRole.Roles {
	// 		dao.RoleDao.AddRole(int64(newAccId), v)
	// 	}

	// 	dao.TeamDao.LoadTeams(int64(newAccId))
	// 	for _, v := range templateTeams.TeamData {
	// 		dao.TeamDao.AddTeam(int64(newAccId), v)
	// 	}
	// 	for _, v := range templateTeams.BattleData {
	// 		dao.TeamDao.AddBattleTeam(int64(newAccId), v)
	// 	}

	// 	dao.EquipDao.LoadEquips(int64(newAccId))
	// 	for _, v := range templateEquips.EquipData {
	// 		dao.EquipDao.AddEquip(int64(newAccId), v)
	// 	}
	// 	for _, v := range templateEquips.EquipPosData {
	// 		dao.EquipDao.AddEquipPos(int64(newAccId), v)
	// 	}
	// 	for _, v := range templateEquips.EquipSuits {
	// 		dao.EquipDao.AddEquipSuit(int64(newAccId), v)
	// 	}
	// 	for _, v := range templateEquips.SuitReward {
	// 		dao.EquipDao.AddSuitInfo(int64(newAccId), v)
	// 	}
	// 	dao.EquipDao.UpdateUseEquipSuit(int64(newAccId), templateEquips.UseEquipSuit, false)

	// 	dao.WeaponDao.LoadWeapon(int64(newAccId))
	// 	dao.WeaponDao.AddWeapons(int64(newAccId), templateWeapon.Weapons)
	// 	dao.WeaponDao.AddSecondaryWeapons(int64(newAccId), templateWeapon.SecondaryWeapons)
	// }

	return msg.ErrCode_SUCC
}
