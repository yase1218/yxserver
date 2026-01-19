package service

import (
	"errors"
	"fmt"
	"gameserver/internal/async"
	"gameserver/internal/config"
	"gameserver/internal/game/builder"
	"gameserver/internal/game/common"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/game/uid"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"github.com/zy/game_data/template"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type CreateUserCost struct {
	TotalCost time.Duration

	NameCost   time.Duration
	AllocCost  time.Duration
	VerifyCost time.Duration

	UserCostData *UserCost
}

type UserCost struct {
	BaseCost     time.Duration
	GemsCost     time.Duration
	ItemCost     time.Duration
	ShipCost     time.Duration
	CardCost     time.Duration
	EquipCost    time.Duration
	WeaponCost   time.Duration
	TreasureCost time.Duration
	ContractCost time.Duration
	TeamCost     time.Duration
	PokerCost    time.Duration
	AppearCost   time.Duration
	FunctionCost time.Duration
	PeakCost     time.Duration
	MailCost     time.Duration
}

type CostRecord struct {
	LoginTotalCost    map[string]time.Duration
	LoadAccountCost   map[string]time.Duration
	CreateAccountCost map[string]time.Duration
	LoadUserCost      map[string]time.Duration
	CreateUserCost    map[string]*CreateUserCost
	InsertUserCost    map[string]time.Duration
	VerifyCost        map[string]time.Duration
	CorrectCost       map[string]time.Duration
	RefreshCost       map[string]time.Duration
	HodgeCost         map[string]time.Duration

	Logout1 map[string]time.Duration
	Logout2 map[string]time.Duration
	Logout3 map[string]time.Duration
}

var (
	CostRecordData = &CostRecord{
		LoginTotalCost:    make(map[string]time.Duration),
		LoadAccountCost:   make(map[string]time.Duration),
		CreateAccountCost: make(map[string]time.Duration),
		LoadUserCost:      make(map[string]time.Duration),
		CreateUserCost:    make(map[string]*CreateUserCost),
		InsertUserCost:    make(map[string]time.Duration),
		VerifyCost:        make(map[string]time.Duration),
		CorrectCost:       make(map[string]time.Duration),
		RefreshCost:       make(map[string]time.Duration),
		HodgeCost:         make(map[string]time.Duration),

		Logout1: make(map[string]time.Duration),
		Logout2: make(map[string]time.Duration),
		Logout3: make(map[string]time.Duration),
	}

	costLimit = time.Duration(30 * time.Millisecond)
)

//func OnLogin(p *player.Player, account_id string, packet_id uint32, ip string, req *msg.RequestLogin) error {
// now := time.Now()
// var (
// 	loadAccountCost   time.Duration
// 	createAccountCost time.Duration
// 	loadUserCost      time.Duration
// 	createUserCost    time.Duration
// 	verifyCost        time.Duration
// 	correctCost       time.Duration
// 	refreshCost       time.Duration
// 	hodgeCost         time.Duration
// )

// createUserCostData := &CreateUserCost{}
// // account_id := utils.StrToUInt64(req.UserId)
// // if account_id == 0 {
// // 	log.Error("handleLoginReq accountId error", zap.String("userId", req.UserId))
// // 	return nil
// // }

// defer func() {
// 	cost := time.Since(now)
// 	if cost > costLimit {
// 		CostRecordData.LoginTotalCost[account_id] = time.Duration(cost.Milliseconds())
// 	}
// 	if loadAccountCost > costLimit {
// 		CostRecordData.LoadAccountCost[account_id] = time.Duration(loadAccountCost.Milliseconds())
// 	}
// 	if createAccountCost > costLimit {
// 		CostRecordData.CreateAccountCost[account_id] = time.Duration(createAccountCost.Milliseconds())
// 	}
// 	if loadUserCost > costLimit {
// 		CostRecordData.LoadUserCost[account_id] = time.Duration(loadUserCost.Milliseconds())
// 	}
// 	if createUserCost > costLimit {
// 		CostRecordData.CreateUserCost[account_id] = createUserCostData
// 	}
// 	if verifyCost > costLimit {
// 		CostRecordData.VerifyCost[account_id] = time.Duration(verifyCost.Milliseconds())
// 	}
// 	if correctCost > costLimit {
// 		CostRecordData.CorrectCost[account_id] = time.Duration(correctCost.Milliseconds())
// 	}
// 	if refreshCost > costLimit {
// 		CostRecordData.RefreshCost[account_id] = time.Duration(refreshCost.Milliseconds())
// 	}
// 	if hodgeCost > costLimit {
// 		CostRecordData.HodgeCost[account_id] = time.Duration(hodgeCost.Milliseconds())
// 	}
// }()

// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// defer cancel()
// account_res := db.LocalMongoReader.FindOne(
// 	ctx,
// 	config.Conf.LocalMongo.DB,
// 	db.CollectionName_UserAccount,
// 	bson.M{"accountid": account_id},
// )
// loadAccountCost = time.Since(now)

// //var isNewAccount = 0
// // 账号查询
// user_account := &model.UserAccount{}
// err := account_res.Decode(user_account)
// if err != nil {
// 	if errors.Is(err, mongo.ErrNoDocuments) { // 新账号
// 		user_account.AccountId = account_id
// 		user_account.UserId = uid.GenUserId()
// 		user_account.CreateTime = uint32(now.Unix())
// 		//isNewAccount = 1
// 		// 创建账号
// 		if er := player.CreateAccount(user_account); er != nil {
// 			log.Error("handleLoginReq create_account failed",
// 				zap.String("accountId", account_id),
// 				zap.Uint64("userId", user_account.UserId),
// 				zap.Error(err),
// 			)
// 			return er
// 		}
// 	} else {
// 		log.Error("handleLoginReq load account failed", zap.String("accountId", account_id), zap.Error(err))
// 		return err
// 	}
// }

// // 角色查询
// p.UserData, err = player.LoadUser(ctx, user_account.UserId)
// channel_id := uint32(1)
// if err != nil {
// 	if errors.Is(err, mongo.ErrNoDocuments) {
// 		// 创建角色
// 		//channel_id := req.ChannelId
// 		e := CreateUser(p, account_id, user_account.UserId, channel_id, now, createUserCostData)
// 		if e != nil {
// 			log.Error("handleLoginReq create_user failed",
// 				zap.String("accountId", account_id),
// 				zap.Uint64("userId", user_account.UserId),
// 				zap.Error(e),
// 			)
// 			return e
// 		}
// 	} else {
// 		log.Error("handleLoginReq user not exist", zap.String("accountId", account_id), zap.Error(err))
// 		return err
// 	}
// } else {
// 	verifyTime := time.Now()
// 	e := VerifyLogin(p, now)
// 	verifyCost = time.Since(verifyTime)
// 	if e != nil {
// 		log.Error("verify login failed", zap.String("accountId", account_id), zap.Error(e))
// 		return e
// 	}
// }

// correctTime := time.Now()
// change := LoginCorrect(p, now)
// correctCost = time.Since(correctTime)
// if change {

// }

// refreshTime := time.Now()
// refresh_onlogin(p, now)
// refreshCost = time.Since(refreshTime)

// // TODO login check token、count_limit 、ban……

// hodgeTime := time.Now()
// p.ChannelId = channel_id
// p.UserData.BaseInfo.LoginTime = uint32(now.Unix())
// p.UserData.BaseInfo.ExtraInfo = req.ExtraInfo
// p.UserData.BaseInfo.Ip = ip
// // if p.UserData.BaseInfo.ShipId == 0 {
// // 	p.UserData.BaseInfo.ShipId = template.GetSystemItemTemplate().InitShip[0]
// // }

// // TODO rbi

// cur := uint32(now.Unix())
// if p.UserData.BaseInfo.ApData.BuyTimes > 0 &&
// 	cur > p.UserData.BaseInfo.ApData.NextBuyTime {
// 	p.UserData.BaseInfo.ApData.BuyTimes = 0
// 	p.UserData.BaseInfo.ApData.NextBuyTime = tools.GetDailyRefreshTime()
// }

// // 超过一天则清除连续活跃度
// if p.UserData.BaseInfo.ActiveDay > 0 {
// 	temp := p.UserData.BaseInfo.LastActiveTime + 2*24*3600 + template.GetSystemItemTemplate().RefreshHour*3600
// 	if p.UserData.BaseInfo.LoginTime > temp {
// 		p.UserData.BaseInfo.ActiveDay = 0
// 	}
// }

// // 重置首充
// for i := 0; i < len(p.UserData.BaseInfo.Charge); i++ {
// 	id := p.UserData.BaseInfo.Charge[i].Id
// 	if p.UserData.BaseInfo.Charge[i].Value == 0 {
// 		continue
// 	}

// 	if config := template.GetChargeTemplate().GetCharge(id); config != nil {
// 		if config.FirstPurchaseExtraResetTime > 0 &&
// 			int(cur) >= config.FirstPurchaseExtraResetTime &&
// 			p.UserData.BaseInfo.Charge[i].ResetTime < config.FirstPurchaseExtraResetTime {
// 			p.UserData.BaseInfo.Charge[i].ResetTime = int(cur)
// 			p.UserData.BaseInfo.Charge[i].Value = 0
// 		}
// 	}
// }

// ResetDailyAp(p, false)
// RefreshShopLogin(p) // 商店
// ResetAdData(p)

// LoadOfflineMail(p)
// //LoadOfflineOrder(p)
// UserOrdersShipment(p)

// // 默认检查
// p.CheckBattle = 1

// p.InWorldChannel = false

// UpdateFunctionPreview(p, msg.ConditionType_Condition_Open_Server_Days)
// UpdateFunctionPreview(p, msg.ConditionType_Condition_Account_Days)
// UpdateFunctionPreview(p, msg.ConditionType_Condition_Pass_Mission)

// GlobalAttrChange(p, false)

// timeNow := time.Now()

// if !utils.IsSameDay(timeNow.Unix(), p.UserData.BaseInfo.LastLoginAt.Unix()) {
// 	OnCrossDay(p)
// }

// p.LastTick = int64(timeNow.Unix())
// if utils.IsPreviousDay(p.UserData.BaseInfo.LastLoginAt) {
// 	p.UserData.BaseInfo.LoginCnt++
// 	p.UserData.BaseInfo.LastLoginAt = timeNow
// 	p.SaveBaseInfo()
// }
// p.SdkChannelNo = req.SdkChannelNo
// p.Os = int(req.Os)
// // p.TdaCommonAttr.UserId = p.GetOpenId()
// // p.TdaCommonAttr.AccountId = fmt.Sprintf("%d", p.GetUserId())
// // tapData := &tapping.Login{
// // 	CommonAttr:    p.TdaCommonAttr,
// // 	RoleName:      p.GetNick(),
// // 	Action:        1,
// // 	IsFirst:       isNewAccount,
// // 	LastLoginTime:  p.UserData.BaseInfo.LastLoginAt,
// // 	AddDay:        p.UserData.BaseInfo.LoginCnt,
// // 	LastOutTime:   p.UserData.BaseInfo.LastLogoutAt,
// // 	AccountId:     p.GetOpenId(),
// // }
// // p.TappingBoth(tapData, "login")
// hodgeCost = time.Since(hodgeTime)
// 	return nil
// }

func OnLogout(p *player.Player, now time.Time) {
	var (
		cost1 time.Duration
		cost2 time.Duration
		cost3 time.Duration
	)

	defer func() {
		if cost1 > costLimit {
			CostRecordData.Logout1[p.GetOpenId()] = time.Duration(cost1.Milliseconds())
		}
		if cost2 > costLimit {
			CostRecordData.Logout2[p.GetOpenId()] = time.Duration(cost2.Milliseconds())
		}
		if cost3 > costLimit {
			CostRecordData.Logout3[p.GetOpenId()] = time.Duration(cost3.Milliseconds())
		}
	}()

	time1 := time.Now()
	if !EquipStageLeave(p) {
		EquipStageTeamLeave(0, p, true)
	}
	cost1 = time.Since(time1)

	time2 := time.Now()

	OnPFLogOut(p)
	cost2 = time.Since(time2)

	// p.TappingBoth(&tapping.Logout{
	// 	CommonAttr: p.TdaCommonAttr,
	// }, "logout")
	time3 := time.Now()
	p.UserData.BaseInfo.LastLogoutAt = now
	log.Info("player logout", ZapUser(p))
	p.SaveDirtySync()
	player.DelPlayer(p)
	cost3 = time.Since(time3)
}

func ResLogin(packet_id uint32, now time.Time, p *player.Player) {
	res := &msg.ResponseLogin{
		Result:                   msg.ErrCode_SUCC,
		Info:                     ToProtocolAccountInfo(p.UserData),
		ServerTime:               now.UnixMilli(),
		ApInfo:                   builder.BuildApData(p.UserData),
		MissionId:                uint32(p.UserData.StageInfo.MissionId),
		GuideData:                builder.BuildGuideInfo(p.UserData),
		PopUps:                   builder.BuildPopUps(p.UserData),
		ReInfo:                   builder.BuildChargeInfo(p.UserData),
		FundInfo:                 builder.BuildFundInfo(p.UserData),
		AdData:                   builder.BuildAdData(p.UserData),
		MonthCardDailyRewardTime: p.UserData.BaseInfo.MonthCardDailyRewardTime,
		DailyAp:                  builder.BuildDailyApInfo(p.UserData),
		TalentData:               builder.BuildTalentData(p.UserData),
		OpenServerDays:           1, // TODO get serverinfo from center
		IsNewAccount:             p.UserData.IsRegister,
		Time:                     ToProtocolServerTime(now),
	}

	if bc := GetBattleCacheByUid(p.GetUserId()); bc != nil {
		if now.Before(bc.DeadLine) {
			res.FsId = bc.FsId
			res.BattleId = bc.BattleId
			res.StageId = bc.StageId
		} else {
			RemoveBattleCacheByUid(p.GetUserId())
		}
	}

	p.State = publicconst.Online

	p.SendResponse(packet_id, res, res.Result)

	p.NtfPoker()
	// TODO service.ServMgr.GetCommonService().AddStaticsData(userData, publicconst.Statics_Login_Id, fmt.Sprintf("loginTime:%v,lastLoginTime:%v|", userData.AccountInfo.LoginTime, lastLoginTime))
	log.Info("Login succ", ZapUser(p))
}

func AfterLogin(p *player.Player) {
	RefreshActivity(p, false)
	updateActivitiesData(p)
	CheckNewMethod(p)
	// TODO 老宠物移除暂留
	// ResetPetData(p, false)
	UpdateFriend(p)
}

func CreateUser(p *player.Player, account_id string, uid uint64, channel uint32, now time.Time, cost *CreateUserCost) error {
	// nameTime := time.Now()
	// nick := fmt.Sprintf("kulu_%v", account_id)

	// var (
	// 	rc   = rdb_single.Get()
	// 	rCtx = rdb_single.GetCtx()
	// 	key  = RdbUserNameKey()
	// )
	// _, err1 := rc.SAdd(rCtx, key, nick).Result()
	// if err1 != nil {
	// 	log.Error("SAdd player name err", zap.Error(err1), zap.String("nick", nick))
	// 	return err1
	// }
	// cost.NameCost = time.Duration(time.Since(nameTime).Milliseconds())

	// allocTime := time.Now()
	// p.UserData = &model.UserData{
	// 	ServerId:   config.Conf.ServerId,
	// 	ServerName: config.Conf.ServerName,
	// 	UserId:     uid,
	// 	ChannelId:  channel,
	// 	AccountId:  account_id,
	// 	Nick:       nick,
	// 	Level:      1,
	// 	IsNew:      true,
	// 	IsRegister: true,
	// 	HeadImg:    template.GetAppearanceTemplate().GetDefaultByType(1),
	// 	HeadFrame:  template.GetAppearanceTemplate().GetDefaultByType(2),
	// 	Title:      template.GetAppearanceTemplate().GetDefaultByType(3),
	// }
	// cost.AllocCost = time.Duration(time.Since(allocTime).Milliseconds())
	// verifyTime := time.Now()
	// if err := VerifyLogin(p, now); err != nil {
	// 	return err
	// }
	// cost.VerifyCost = time.Duration(time.Since(verifyTime).Milliseconds())

	// cost.UserCostData = CreateUserData(p, now)
	// InsertUserTime := time.Now()
	// p.InsertUserToDb()
	// insertUserCost := time.Since(InsertUserTime)
	// if insertUserCost > costLimit {
	// 	CostRecordData.InsertUserCost[account_id] = time.Duration(insertUserCost.Milliseconds())
	// }
	return nil
}

// func CreateUserData(p *player.Player, now time.Time) *UserCost {
// 	baseTime := time.Now()
// 	create_baseinfo(p, now)
// 	baseCost := time.Since(baseTime)

// 	gemTime := time.Now()
// 	create_gems(p)
// 	gemCost := time.Since(gemTime)

// 	itemTime := time.Now()
// 	create_items(p)
// 	itemCost := time.Since(itemTime)

// 	shipTime := time.Now()
// 	create_ship(p)
// 	shipCost := time.Since(shipTime)

// 	cardTime := time.Now()
// 	create_card_pool(p)
// 	cardCost := time.Since(cardTime)

// 	equipTime := time.Now()
// 	create_equip(p)
// 	equipCost := time.Since(equipTime)

// 	weaponTime := time.Now()
// 	create_weapon(p)
// 	weaponCost := time.Since(weaponTime)

// 	treasureTime := time.Now()
// 	create_treasure(p)
// 	treasureCost := time.Since(treasureTime)

// 	contractTime := time.Now()
// 	create_contract(p)
// 	contractCost := time.Since(contractTime)

// 	teamTime := time.Now()
// 	create_team(p)
// 	teamCost := time.Since(teamTime)

// 	pokerTime := time.Now()
// 	create_poker(p)
// 	pokerCost := time.Since(pokerTime)

// 	appearTime := time.Now()
// 	create_appearance(p)
// 	appearCost := time.Since(appearTime)

// 	functionTime := time.Now()
// 	create_function_preview(p)
// 	functionCost := time.Since(functionTime)

// 	peakTime := time.Now()
// 	create_peak_fight(p, now)
// 	peakCost := time.Since(peakTime)

// 	mailTime := time.Now()
// 	create_mail(p)
// 	mailCost := time.Since(mailTime)

// 	return &UserCost{
// 		BaseCost:     time.Duration(baseCost.Milliseconds()),
// 		GemsCost:     time.Duration(gemCost.Milliseconds()),
// 		ItemCost:     time.Duration(itemCost.Milliseconds()),
// 		ShipCost:     time.Duration(shipCost.Milliseconds()),
// 		CardCost:     time.Duration(cardCost.Milliseconds()),
// 		EquipCost:    time.Duration(equipCost.Milliseconds()),
// 		WeaponCost:   time.Duration(weaponCost.Milliseconds()),
// 		TreasureCost: time.Duration(treasureCost.Milliseconds()),
// 		ContractCost: time.Duration(contractCost.Milliseconds()),
// 		TeamCost:     time.Duration(teamCost.Milliseconds()),
// 		PokerCost:    time.Duration(pokerCost.Milliseconds()),
// 		AppearCost:   time.Duration(appearCost.Milliseconds()),
// 		FunctionCost: time.Duration(functionCost.Milliseconds()),
// 		PeakCost:     time.Duration(peakCost.Milliseconds()),
// 		MailCost:     time.Duration(mailCost.Milliseconds()),
// 	}
// }

func create_mail(p *player.Player) {
	p.UserData.MailData = &model.UserMail{
		Mails: make(map[int64]*model.Mail),
	}
}

func create_peak_fight(p *player.Player, now time.Time) {
	season := uint32(0)
	season_t := template.GetBattlePassTemplate().GetCurSeason(uint32(now.Unix()))
	if season_t != nil {
		season = season_t.Season
	} else {
		season = 1
	}
	p.UserData.PeakFight = &model.PeakFight{
		BattleMatchId: 1,
		Season:        season,
	}
}

func create_baseinfo(p *player.Player, now time.Time) {
	p.UserData.BaseInfo.CreateTime = uint32(now.Unix())
}

func create_function_preview(p *player.Player) {
	for k, v := range fp_ids {
		if v == nil {
			if _, ok := p.UserData.FunctionPreview.Data[k]; !ok {
				p.UserData.FunctionPreview.Data[k] = msg.TaskState_Task_Complete
			}
		}
	}
}

func create_appearance(p *player.Player) {
	defaultIds := template.GetAppearanceTemplate().GetDefault()
	var addAppearances []*model.Appearance
	for _, id := range defaultIds {
		addAppearances = append(addAppearances, model.NewAppearance(id))
	}
	AddAppearances(p, addAppearances, false)
}

func create_poker(p *player.Player) {
	p.UserData.Poker.CommData = template.GetSystemItemTemplate().InitPoker
}

func create_team(p *player.Player) {
	shipId := template.GetSystemItemTemplate().InitShipsList[0]
	roleId := template.GetSystemItemTemplate().InitRole[0]
	for i := 0; i < int(template.GetSystemItemTemplate().TeamNum); i++ {
		team := model.NewTeam(uint32(i), shipId, roleId, nil, nil)
		p.UserData.Team.TeamData = append(p.UserData.Team.TeamData, team)
	}
}

func create_contract(p *player.Player) {
	p.UserData.Contract = &model.Contract{
		AccountId: p.GetUserId(),
		TaskIds:   RandCrontractTaskIds(),
		ResetDate: tools.GetDailyRefreshTime(),
	}
}

func create_treasure(p *player.Player) {
	initCommonTreasure := template.GetSystemItemTemplate().InitTreasure
	if !tools.ListUint32Equal(p.UserData.Treasure.CommData, initCommonTreasure) {
		p.UserData.Treasure.CommData = initCommonTreasure
	}
}

func create_card_pool(p *player.Player) {
	var initCardPools []*model.CardPool
	cardList := template.GetLotteryShipTemplate().GetInitPool()
	for i := 0; i < len(cardList); i++ {
		var freeTimes uint32 = 0
		var nextResetTime uint32 = 0
		if cardList[i].XDayFreeTime > 0 {
			freeTimes = 1
		}

		if cardList[i].CardType != 2 {
			nextResetTime = tools.GetDailyRefreshTime()
		}
		initCardPools = append(initCardPools,
			model.NewCardPool(cardList[i].Id, freeTimes,
				tools.GetDailyXRefreshTime(cardList[i].XDayFreeTime, template.GetSystemItemTemplate().RefreshHour),
				nextResetTime, 0, 0, 0))
	}
	p.UserData.CardPool.CardPools = initCardPools
}

func create_equip(p *player.Player) {
	initEquip := template.GetSystemItemTemplate().InitEquip
	for i := 0; i < len(initEquip); i++ {
		addEquip(p, initEquip[i].ItemId, initEquip[i].ItemNum, false)
	}
}

func create_gems(p *player.Player) {
	for i := 0; i < len(p.UserData.Equip.GemPos); i++ {
		p.UserData.Equip.GemPos[i] = make([]uint64, template.GetSystemItemTemplate().GemSlotMax)
		for j := 0; j < len(p.UserData.Equip.GemPos[i]); j++ {
			p.UserData.Equip.GemPos[i][j] = 0
		}
	}
}

func create_items(p *player.Player) {
	initItems := template.GetSystemItemTemplate().InitItem
	var notifyItems []uint32
	for k := 0; k < len(initItems); k++ {
		addItems := AddPlayerItem(p,
			initItems[k].ItemId, int32(initItems[k].ItemNum), publicconst.InitAddItem, false)

		notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
	}

}

func create_ship(p *player.Player) {
	// ships := template.GetSystemItemTemplate().InitShip
	// for i := 0; i < len(ships); i++ {
	// 	AddShip(p, ships[i], publicconst.InitAddItem, true)
	// }
}

func create_weapon(p *player.Player) {
	var newWeapons []uint32
	initWeaponIds := template.GetSystemItemTemplate().InitWeapons
	for i := 0; i < len(initWeaponIds); i++ {
		if temp := getWeapon(p, initWeaponIds[i]); temp == nil {
			newWeapons = append(newWeapons, initWeaponIds[i])
		}
	}

	if len(newWeapons) > 0 {
		for i := 0; i < len(newWeapons); i++ {
			AddWeapon(p, newWeapons[i], publicconst.InitAddItem)
		}
	}
}

func OnCrossDay(p *player.Player) {
	DayResetDesert(p, true)
	HandleUpdateChargePackageState(p)
	RefreshEquipStage(p)
	HandleRefreshLikesInfo(p)
	OnCrossDayFreshWeekPass(p)
}

func OnInitPlayerNameAndShip(p *player.Player, playerName string, shipId uint32) msg.ErrCode {
	p.UserData.Nick = playerName
	AddShip(p, shipId, publicconst.InitAddItem, false)
	UpdateClientShipChange(p, []uint32{shipId})
	p.UserData.BaseInfo.ShipId = shipId
	p.UserData.IsRegister = false

	GlobalAttrChange(p, true)

	p.SaveRegister()
	p.SaveNick()
	p.SaveBaseInfo()
	return msg.ErrCode_SUCC
}

func OnInitPlayerNameAndShipCheck(p *player.Player, playerName string, shipId uint32) msg.ErrCode {
	if !p.UserData.IsRegister {
		return msg.ErrCode_Name_Is_Setting
	}

	if playerName == "" {
		return msg.ErrCode_Name_Not_Vaild
	}

	playerName = strings.Trim(playerName, " ")
	if utf8.RuneCountInString(playerName) <= 0 {
		return msg.ErrCode_Name_Not_Vaild
	}

	if utils.ContainsEmoji(playerName) || utils.HasInterpretedEscapeChars(playerName) {
		return msg.ErrCode_Name_Not_Vaild
	}

	if utf8.RuneCountInString(playerName) > int(template.GetSystemItemTemplate().NickLen) {
		return msg.ErrCode_NICK_TO_LONG
	}

	//if template.GetForbiddenTemplate().HasForbidden(playerName) {
	//	return msg.ErrCode_NICK_HAS_FORBIDDEN
	//}

	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
		key  = fmt.Sprintf("%v:UserNames", config.Conf.ServerId)
	)

	old := p.UserData.Nick
	_, err1 := rc.SRem(rCtx, key, old).Result()
	if err1 != nil {
		log.Error("SRem remove nick error", zap.Error(err1), zap.String("old nick", old))
		return msg.ErrCode_NICK_EXIST
	}

	res, err := rc.SAdd(rCtx, key, playerName).Result()
	if err != nil {
		log.Error("SAdd player name err", zap.Error(err))
		return msg.ErrCode_Name_Not_Vaild
	}

	if res == 0 {
		return msg.ERRCODE_NAME_EXIST
	}

	defaultShips := template.GetSystemItemTemplate().InitShipsList
	isFind := false
	for i := 0; i < len(defaultShips); i++ {
		if defaultShips[i] == shipId {
			isFind = true
			break
		}
	}

	if !isFind {
		return msg.ERRCODE_SHIP_NOT_EXIST
	}
	return msg.ErrCode_SUCC
}

func OnRandomGenPlayerName(p *player.Player) string {
	var name string
	//for i := 0; i < 100; i++ {
	randomNameId := template.GetRandomNameTemplate().RandOne()

	num, err := strconv.Atoi(randomNameId)
	if err != nil {
		log.Error("random player name error", zap.String("random cfg id", randomNameId))
	} else {
		name = template.GetLanguageTemplate().GetContent(uint32(num))
		//break
	}
	//}
	return name
}

func ZapUser(p *player.Player) zap.Field {
	if p == nil {
		return zap.String("info", "nil player")
	}
	return zap.String("info", p.Info())
}

func refresh_onlogin(p *player.Player, now time.Time) {
	refresh_global_mail_onlogin(p, now)
}

func ToProtocolServerTime(now time.Time) *msg.ServerTimeSt {
	info := common.GetServerInfo()
	weeks := tools.GetWeekCount(info.OpenTime)
	var serverTime = &msg.ServerTimeSt{
		TimeStamp: uint64(now.UnixMilli()),
		Week:      uint32(weeks),
		Day:       uint32(now.Day()),
		Zone:      float64(tools.GetCurrentTimezoneOffset()),
	}

	now.Local().Location()
	return serverTime
}

func OnLogin(accountId string,
	gateId, fsId, packetId uint32,
	sessionId uint64,
	ip string,
	req *msg.RequestLogin) error {
	// 异步加载角色数据
	log.Info("trace login push async load req",
		zap.String("accountId", accountId),
		zap.Uint64("sessionId", sessionId),
	)
	if err := async.Push(&async.AyncReadUser{
		AccountId:    accountId,
		GateId:       gateId,
		FsId:         fsId,
		PacketId:     packetId,
		SessionId:    sessionId,
		Extra:        req.ExtraInfo,
		SdkChannelNo: req.SdkChannelNo,
		Os:           int(req.Os),
		Ip:           ip,
		Cb: func(err error, account string, userData *model.UserData) {
			log.Info("trace login async callback",
				zap.String("accountId", accountId),
				zap.Uint64("sessionId", sessionId),
			)
			now := time.Now()
			if err != nil {
				if !errors.Is(err, mongo.ErrNoDocuments) { //错误
					log.Error("login failed",
						zap.String("account", account),
						zap.Uint32("gateId", gateId),
						zap.Uint64("sessionId", sessionId),
						zap.Error(err))
					return
				} else { // 没有角色数据 创建
					log.Info("trace login no user then create",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
					p := player.CreatePlayer(now)
					p.BindGateId(gateId)
					p.BindSessionId(sessionId)

					log.Info("trace login no user start create user data",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
					e := CreateUserData(p, account, now)
					if e != nil {
						log.Error("login create user failed",
							zap.String("account", account),
							zap.Uint32("gateId", gateId),
							zap.Uint64("sessionId", sessionId),
							zap.Error(e))
						return
					}
					ResLogin(packetId, now, p)
					log.Info("trace login ResLogin",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
					AfterLogin(p)
					player.AddPlayer(p)
				}
			} else {
				log.Info("trace login has user",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				p := player.CreatePlayer(now)
				p.BindGateId(gateId)
				p.BindSessionId(sessionId)
				p.UserData = userData

				if fsId != 0 {
					p.SetFsId(fsId)
				}
				e := VerifyLogin(p, now)
				if e != nil {
					log.Error("verify login failed", zap.String("accountId", account), zap.Error(e))
					return
				}
				log.Info("trace login has user VerifyLogin",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				LoginCorrect(p, now)
				log.Info("trace login has user LoginCorrect",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				refresh_onlogin(p, now)
				log.Info("trace login has user refresh_onlogin",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)

				p.UserData.BaseInfo.LoginTime = uint32(now.Unix())
				p.UserData.BaseInfo.ExtraInfo = req.ExtraInfo
				p.UserData.BaseInfo.Ip = ip

				cur := uint32(now.Unix())
				if p.UserData.BaseInfo.ApData.BuyTimes > 0 &&
					cur > p.UserData.BaseInfo.ApData.NextBuyTime {
					p.UserData.BaseInfo.ApData.BuyTimes = 0
					p.UserData.BaseInfo.ApData.NextBuyTime = tools.GetDailyRefreshTime()
					log.Info("trace login has user ap data buy times",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
				}

				// 超过一天则清除连续活跃度
				if p.UserData.BaseInfo.ActiveDay > 0 {
					temp := p.UserData.BaseInfo.LastActiveTime + 2*24*3600 + template.GetSystemItemTemplate().RefreshHour*3600
					if p.UserData.BaseInfo.LoginTime > temp {
						p.UserData.BaseInfo.ActiveDay = 0
					}
					log.Info("trace login has user ActiveDays",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
				}

				// 重置首充
				for i := 0; i < len(p.UserData.BaseInfo.Charge); i++ {
					id := p.UserData.BaseInfo.Charge[i].Id
					if p.UserData.BaseInfo.Charge[i].Value == 0 {
						continue
					}

					if config := template.GetChargeTemplate().GetCharge(id); config != nil {
						if config.FirstPurchaseExtraResetTime > 0 &&
							int(cur) >= config.FirstPurchaseExtraResetTime &&
							p.UserData.BaseInfo.Charge[i].ResetTime < config.FirstPurchaseExtraResetTime {
							p.UserData.BaseInfo.Charge[i].ResetTime = int(cur)
							p.UserData.BaseInfo.Charge[i].Value = 0
						}
					}
				}
				log.Info("trace login has user check first charge",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)

				ResetDailyAp(p, false)
				log.Info("trace login has user ResetDailyAp",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				RefreshShopLogin(p) // 商店
				log.Info("trace login has user RefreshShopLogin",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				ResetAdData(p)
				log.Info("trace login has user ResetAdData",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)

				LoadOfflineMail(p)
				log.Info("trace login has user LoadOfflineMail",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				//LoadOfflineOrder(p)
				UserOrdersShipment(p)
				log.Info("trace login has user UserOrdersShipment",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)

				// 默认检查
				p.CheckBattle = 1

				p.InWorldChannel = false

				// UpdateFunctionPreview(p, msg.ConditionType_Condition_Open_Server_Days)
				// UpdateFunctionPreview(p, msg.ConditionType_Condition_Account_Days)
				// UpdateFunctionPreview(p, msg.ConditionType_Condition_Pass_Mission)
				log.Info("trace login has user UpdateFunctionPreview",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)

				GlobalAttrChange(p, false)
				log.Info("trace login has user GlobalAttrChange",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				if !utils.IsSameDay(now.Unix(), p.UserData.BaseInfo.LastLoginAt.Unix()) {
					OnCrossDay(p)
					log.Info("trace login has user is OnCrossDay",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
				} else {
					log.Info("trace login has user not OnCrossDay",
						zap.String("accountId", accountId),
						zap.Uint64("sessionId", sessionId),
					)
				}

				p.LastTick = int64(now.Unix())
				if utils.IsPreviousDay(p.UserData.BaseInfo.LastLoginAt) {
					p.UserData.BaseInfo.LoginCnt++
					p.UserData.BaseInfo.LastLoginAt = now
					p.SaveBaseInfo()
				}
				p.SdkChannelNo = req.SdkChannelNo
				p.Os = int(req.Os)
				ResLogin(packetId, now, p)
				log.Info("trace login has user ResLogin",
					zap.String("accountId", accountId),
					zap.Uint64("sessionId", sessionId),
				)
				AfterLogin(p)
				player.AddPlayer(p)
			}
		},
	}); err != nil {
		log.Error("trace login push async failed", zap.Error(err))
	}
	return nil
}

func CreateUserData(p *player.Player, accountId string, now time.Time) error {
	userId := uid.GenUserId()
	nick := "kulu_" + accountId

	p.UserData = &model.UserData{
		ServerId:   config.Conf.ServerId,
		ServerName: config.Conf.ServerName,
		UserId:     userId,
		AccountId:  accountId,
		Nick:       nick,
		Level:      1,
		IsNew:      true,
		IsRegister: true,
		HeadImg:    template.GetAppearanceTemplate().GetDefaultByType(1),
		HeadFrame:  template.GetAppearanceTemplate().GetDefaultByType(2),
		Title:      template.GetAppearanceTemplate().GetDefaultByType(3),
	}

	if err := VerifyLogin(p, now); err != nil {
		return err
	}
	log.Info("trace login no user verify struct",
		zap.String("accountId", accountId))

	create_baseinfo(p, now)
	log.Info("trace login no user create_baseinfo",
		zap.String("accountId", accountId))

	create_gems(p)
	log.Info("trace login no user create_gems",
		zap.String("accountId", accountId))

	create_items(p)
	log.Info("trace login no user create_items",
		zap.String("accountId", accountId))

	create_ship(p)
	log.Info("trace login no user create_ship",
		zap.String("accountId", accountId))

	create_card_pool(p)
	log.Info("trace login no user create_card_pool",
		zap.String("accountId", accountId))

	create_equip(p)
	log.Info("trace login no user create_equip",
		zap.String("accountId", accountId))

	create_weapon(p)
	log.Info("trace login no user create_weapon",
		zap.String("accountId", accountId))

	create_treasure(p)
	log.Info("trace login no user create_treasure",
		zap.String("accountId", accountId))

	create_contract(p)
	log.Info("trace login no user create_contract",
		zap.String("accountId", accountId))

	create_team(p)
	log.Info("trace login no user create_team",
		zap.String("accountId", accountId))

	create_poker(p)
	log.Info("trace login no user create_poker",
		zap.String("accountId", accountId))

	create_appearance(p)
	log.Info("trace login no user create_appearance",
		zap.String("accountId", accountId))

	create_function_preview(p)
	log.Info("trace login no user create_function_preview",
		zap.String("accountId", accountId))

	create_peak_fight(p, now)
	log.Info("trace login no user create_peak_fight",
		zap.String("accountId", accountId))

	create_mail(p)
	log.Info("trace login no user create_mail",
		zap.String("accountId", accountId))

	p.InsertUserToDb()
	log.Info("trace login no user sync insert",
		zap.String("accountId", accountId))
	return nil
}
