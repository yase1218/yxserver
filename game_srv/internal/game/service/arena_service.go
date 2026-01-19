package service

// import (
// 	"errors"
// 	"fmt"
// 	"kernel/protocol"
// 	"kernel/tools"
// 	"math/rand"
// 	"msg"
// 	"slices"
// 	"sort"
// 	"time"

// 	"github.com/v587-zyf/gc/log"
// 	"github.com/zy/game_data/template"
// 	"go.uber.org/zap"

// 	common2 "server/internal/common"
// 	"server/internal/config"
// 	"server/internal/game/common"
// 	"server/internal/game/condition"
// 	"server/internal/game/dao"
// 	"server/internal/game/model"
// 	"server/internal/publicconst"
// )

// type ArenaFightSession struct {
// 	fightId    uint32
// 	attackerId int64
// 	defenserId int64
// 	timeout    int64
// 	rankId     int32
// }

// var (
// 	arena_fightSession   map[int32]*ArenaFightSession //<rankId, data>
// 	arena_playerSessions map[int64]*ArenaFightSession
// 	arena_ranks          []*model.ArenaRankData
// 	arena_playerRank     map[int64]int32
// )

// func NewArenaService() *ArenaService {
// 	return &ArenaService{
// 		arena_fightSession:   make(map[int32]*ArenaFightSession),
// 		arena_playerSessions: make(map[int64]*ArenaFightSession),
// 		arena_playerRank:     make(map[int64]int32),
// 	}
// }

// func (a *ArenaService) OnInit() {
// 	err, rankDatas := dao.ArenaRankDao.LoadAllRank(config.Conf.ServerId)
// 	if err != msg.ErrCode_SUCC {
// 		log.Error("ArenaService LoadRank err", zap.String("err", err.String()))
// 		//log.Errorf("DesertActRank loadMissionRank err:%v", err)
// 		return
// 	}

// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(-1)
// 	if cfg == nil {
// 		return
// 	}

// 	for i := 0; i < int(cfg.RankEnd); i++ {
// 		rankId := int32(i + 1)

// 		robot := randRobot(rankId)
// 		d := &model.ArenaRankData{
// 			RankId:     rankId,
// 			IsRobot:    true,
// 			AccountId:  int64(robot.AccountId),
// 			MonsterIds: template.GetSystemItemTemplate().ArenaInitDefendMonster,
// 		}
// 		a.arena_ranks = append(a.arena_ranks, d)
// 	}

// 	for i := 0; i < len(rankDatas); i++ {
// 		rank := rankDatas[i]
// 		if int(rank.RankId) >= len(a.arena_ranks) {
// 			continue
// 		}

// 		rank.IsRobot = false

// 		if needUseDefultMonsters(rank.MonsterIds) {
// 			rank.MonsterIds = template.GetSystemItemTemplate().ArenaInitDefendMonster
// 		}

// 		a.arena_playerRank[rank.AccountId] = rank.RankId
// 		a.arena_ranks[rank.RankId-1] = rank
// 	}
// }

// func (a *ArenaService) OnDestroy() {
// }

// func (a *ArenaService) LoadPlayerData(playerData *common.PlayerData) error {
// 	if playerData.Arena == nil {
// 		if _, data := dao.ArenaDao.LoadPlayerData(playerData.GetAccountId()); data != nil {
// 			playerData.Arena = data
// 		} else {
// 			playerData.Arena = &model.ArenaPlayerData{
// 				AccountId: playerData.GetAccountId(),
// 				ServerId:  config.Conf.ServerId,
// 			}

// 			playerData.Arena.Records = make([]*model.ArenaPlayerPkRecordData, 0)
// 			playerData.Arena.UnlockMonsterIds = make([]uint32, 0)
// 			playerData.Arena.DefendMonster = []int32{-1, -1, -1, -1, -1, -1}
// 			for _, v := range template.GetArenaMonsterPosTemplate().GetAll() {
// 				if v.Id > len(playerData.Arena.DefendMonster) {
// 					continue
// 				}

// 				pos := v.Id - 1
// 				if len(v.TaskType) == 0 && len(v.UnlockCost) == 0 {
// 					playerData.Arena.DefendMonster[pos] = 0
// 				}
// 			}

// 			dbModel := model.GetArenaModel()
// 			dbModel.Create(playerData.Arena)
// 		}
// 	}

// 	return nil
// }

// func (a *ArenaService) GetRanks(playerData *common.PlayerData) []*msg.ArenaRankInfo {
// 	var ranks []*msg.ArenaRankInfo
// 	for i := 1; i <= 100; i++ {
// 		ranks = append(ranks, rankData2Proto(int32(i), a.arena_ranks[i-1]))
// 	}

// 	return ranks
// }

// func (a *ArenaService) DayReset(playerData *common.PlayerData) {
// 	playerData.Arena.TodayBuyPkCnt = 0
// 	playerData.Arena.TodayPkCnt = 0
// 	playerData.Arena.TodayUseShips = make([]uint32, 0)
// 	playerData.Arena.LastResetStamp = int64(tools.GetDailyRefreshTime())
// }

// func (a *ArenaService) GetPlayerInfo(playerData *common.PlayerData) *msg.ArenaInfoAck {
// 	needSave := false
// 	cur := time.Now().Unix()

// 	if playerData.Arena.RewardBeginStamp == 0 && msg.ErrCode_SUCC == ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.Arena) {
// 		playerData.Arena.RewardBeginStamp = cur
// 		needSave = true
// 	}

// 	log.Debug("ArenaService GetPlayerInfo,", zap.Any("LastResetStamp", playerData.Arena.LastResetStamp), zap.Any("cur", cur))
// 	if playerData.Arena.LastResetStamp <= cur {
// 		a.DayReset(playerData)
// 		needSave = true
// 	}

// 	ack := &msg.ArenaInfoAck{}
// 	ack.TodayPkCnt = playerData.Arena.TodayPkCnt
// 	ack.TodayBuyCnt = playerData.Arena.TodayBuyPkCnt
// 	ack.RewardBeginStamp = playerData.Arena.RewardBeginStamp
// 	ack.MyRank = a.FindPlayerRank(playerData.GetAccountId())
// 	ack.RankSettlementStamp = int64(tools.GetWeeklyRefreshTime(0))
// 	ack.TotalPkCnt = playerData.Arena.TotalPkCnt
// 	ack.TodayUseShipIds = playerData.Arena.TodayUseShips

// 	for i := 1; i <= 5; i++ {
// 		ack.TopRank = append(ack.TopRank, rankData2Proto(int32(i), a.arena_ranks[i-1]))
// 	}

// 	// showList := []*msg.ArenaRankInfo{}
// 	// cfg := template.GetZombieColiseumTemplate().GetCfgByRank(myrank)
// 	// if cfg != nil {
// 	// 	showList = a.RefreshRankList(myrank)
// 	// }

// 	if needSave {
// 		dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	}

// 	return ack
// }

// func (a *ArenaService) GetPlayerMonsterInfo(playerData *common.PlayerData) *msg.ArenaMonsterInfoAck {

// 	for pos, v := range playerData.Arena.DefendMonster {
// 		if v >= 0 {
// 			continue
// 		}

// 		id := pos + 1
// 		cfg := template.GetArenaMonsterPosTemplate().Get(id)
// 		if cfg == nil {
// 			continue
// 		}

// 		canUnlock := false
// 		if len(cfg.UnlockCost) == 0 {
// 			if len(cfg.TaskType) == 0 {
// 				canUnlock = true
// 			} else if _, ok := condition.GetCondition().Check(playerData, cfg.TaskType); ok {
// 				canUnlock = true
// 			}
// 		}

// 		if canUnlock {
// 			playerData.Arena.DefendMonster[pos] = 0
// 		}
// 	}

// 	for _, v := range template.GetArenaMonsterTemplate().GetAll() {
// 		slices.Sort(playerData.Arena.UnlockMonsterIds)
// 		if _, b := slices.BinarySearch(playerData.Arena.UnlockMonsterIds, v.MonsterId); b {
// 			continue
// 		}

// 		canUnlock := false
// 		if len(v.UnlockCost) == 0 {
// 			if len(v.TaskType) == 0 {
// 				canUnlock = true
// 			} else if _, ok := condition.GetCondition().Check(playerData, v.TaskType); ok {
// 				canUnlock = true
// 			}

// 		}

// 		if canUnlock {
// 			playerData.Arena.UnlockMonsterIds = append(playerData.Arena.UnlockMonsterIds, v.MonsterId)
// 		}
// 	}

// 	info := &msg.ArenaMonsterInfoAck{}
// 	info.DefendMonster = playerData.Arena.DefendMonster
// 	info.UnlockMonster = playerData.Arena.UnlockMonsterIds
// 	return info
// }

// func (a *ArenaService) UnlockPos(playerData *common.PlayerData, pos int) msg.ErrCode {

// 	if !(1 <= pos && pos <= len(playerData.Arena.DefendMonster)) {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	cfg := template.GetArenaMonsterPosTemplate().Get(pos)
// 	if cfg == nil {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	// info := &msg.ArenaMonsterInfoAck{}
// 	if playerData.Arena.DefendMonster[pos-1] >= 0 {
// 		return msg.ErrCode_ARENA_POS_IS_UNLCOK
// 	}

// 	for _, v := range cfg.UnlockCost {
// 		if !ServMgr.GetItemService().EnoughItem(playerData.GetAccountId(),
// 			v.ItemId, v.ItemNum) {
// 			return msg.ErrCode_NO_ENOUGH_ITEM
// 		}
// 	}

// 	var notifyItems []uint32
// 	// 消耗道具
// 	for _, v := range cfg.UnlockCost {
// 		ServMgr.GetItemService().CostItem(playerData.GetAccountId(), v.ItemId, v.ItemNum, publicconst.ArenaBuyPkCnt, false)
// 		notifyItems = append(notifyItems, v.ItemId)
// 	}
// 	ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)

// 	playerData.Arena.DefendMonster[pos-1] = 0
// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) UnlockMonster(playerData *common.PlayerData, monsterId uint32) msg.ErrCode {

// 	// info := &msg.ArenaMonsterInfoAck{}
// 	cfg := template.GetArenaMonsterTemplate().Get(monsterId)
// 	if cfg == nil {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	slices.Sort(playerData.Arena.UnlockMonsterIds)
// 	if _, b := slices.BinarySearch(playerData.Arena.UnlockMonsterIds, monsterId); b {
// 		return msg.ErrCode_ARENA_MONSTER_IS_UNLCOK
// 	}

// 	for _, v := range cfg.UnlockCost {
// 		if !ServMgr.GetItemService().EnoughItem(playerData.GetAccountId(),
// 			v.ItemId, v.ItemNum) {
// 			return msg.ErrCode_NO_ENOUGH_ITEM
// 		}
// 	}

// 	var notifyItems []uint32
// 	// 消耗道具
// 	for _, v := range cfg.UnlockCost {
// 		ServMgr.GetItemService().CostItem(playerData.GetAccountId(), v.ItemId, v.ItemNum, publicconst.ArenaBuyPkCnt, false)
// 		notifyItems = append(notifyItems, v.ItemId)
// 	}
// 	ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)

// 	playerData.Arena.UnlockMonsterIds = append(playerData.Arena.UnlockMonsterIds, monsterId)

// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) BuyPkCnt(playerData *common.PlayerData) msg.ErrCode {
// 	if playerData.Arena.TodayBuyPkCnt >= template.GetSystemItemTemplate().ArenaDayBuyCntLimit {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	for _, v := range template.GetSystemItemTemplate().ArenaDayBuyCost {
// 		if !ServMgr.GetItemService().EnoughItem(playerData.GetAccountId(),
// 			v.ItemId, v.ItemNum) {
// 			return msg.ErrCode_NO_ENOUGH_ITEM
// 		}
// 	}

// 	var notifyItems []uint32
// 	// 消耗道具
// 	for _, v := range template.GetSystemItemTemplate().ArenaDayBuyCost {
// 		ServMgr.GetItemService().CostItem(playerData.GetAccountId(), v.ItemId, v.ItemNum, publicconst.ArenaBuyPkCnt, false)
// 		notifyItems = append(notifyItems, v.ItemId)
// 	}
// 	ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)

// 	playerData.Arena.TodayBuyPkCnt++

// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) GetAFKReward(playerData *common.PlayerData) (int64, []*msg.SimpleItem) {
// 	var rewardRet []*msg.SimpleItem
// 	myrank := a.FindPlayerRank(playerData.GetAccountId())

// 	nextBeginStamp := playerData.Arena.RewardBeginStamp
// 	cur := time.Now().Unix()
// 	if playerData.Arena.RewardBeginStamp == 0 && msg.ErrCode_SUCC == ServMgr.GetCommonService().FunctionOpen(playerData, publicconst.Arena) {
// 		playerData.Arena.RewardBeginStamp = cur
// 	}

// 	if playerData.Arena.RewardBeginStamp == 0 {
// 		return nextBeginStamp, rewardRet
// 	}

// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(myrank)
// 	if cfg == nil {
// 		return nextBeginStamp, rewardRet
// 	}

// 	diff := cur - playerData.Arena.RewardBeginStamp
// 	max := int64(8 * 60 * 60)
// 	if diff > max {
// 		diff = max
// 	}
// 	cnt := diff / (60 * 60)

// 	if cnt <= 0 {
// 		return nextBeginStamp, rewardRet
// 	}

// 	var notifyItems []uint32
// 	for _, item := range cfg.PatrolRewards {
// 		itemId := item.ItemId
// 		itemNum := item.ItemNum * uint32(cnt)
// 		addItems := ServMgr.GetItemService().AddItem(playerData.GetAccountId(), itemId, itemNum, publicconst.PeakFight, false)
// 		notifyItems = append(notifyItems, ServMgr.GetItemService().GetSimpleItemIds(addItems)...)
// 		rewardRet = append(rewardRet, &msg.SimpleItem{ItemId: itemId, ItemNum: itemNum})
// 	}

// 	nextBeginStamp = cur - diff%(60*60)
// 	if diff > max {
// 		nextBeginStamp = cur
// 	}

// 	playerData.Arena.RewardBeginStamp = nextBeginStamp

// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)

// 	return nextBeginStamp, rewardRet
// }

// func (a *ArenaService) RefreshPkList(playerData *common.PlayerData) []*msg.ArenaRankInfo {

// 	myrank := a.FindPlayerRank(playerData.GetAccountId())

// 	showList := []*msg.ArenaRankInfo{}
// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(myrank)
// 	if cfg == nil {
// 		return showList
// 	}

// 	if myrank < 0 {
// 		myrank = cfg.ShowRankEnd
// 	}

// 	var rankIds []int
// 	if cfg.ShowRankEnd <= 10 {
// 		for i := 1; i <= int(cfg.ShowRankEnd); i++ {
// 			if i == int(myrank) {
// 				continue
// 			}

// 			rankIds = append(rankIds, i)
// 		}
// 	} else {

// 		frontCnt := int32(7)

// 		if myrank-cfg.ShowRankStart <= 8 {
// 			frontCnt = myrank - 1
// 		}

// 		if cfg.ShowRankEnd-myrank <= 3 {
// 			frontCnt = 10 - (cfg.ShowRankEnd - myrank)
// 		}

// 		maxRandCnt := 100

// 		if frontCnt <= 7 && myrank == cfg.ShowRankStart+7 {
// 			for i := int32(1); i <= frontCnt; i++ {
// 				rankIds = append(rankIds, int(cfg.ShowRankStart+i))
// 			}
// 		} else {
// 			for i := 0; i < maxRandCnt; i++ {
// 				rankId := int(cfg.ShowRankStart + 1 + int32(rand.Intn(int(myrank-cfg.ShowRankStart)+1)))

// 				if slices.Contains(rankIds, int(rankId)) {
// 					continue
// 				}

// 				rankIds = append(rankIds, rankId)

// 				if len(rankIds) == int(frontCnt) {
// 					break
// 				}
// 			}
// 		}

// 		for i := 0; i < maxRandCnt; i++ {
// 			if len(rankIds) == 10 {
// 				break
// 			}

// 			rankId := int(myrank + 1 + int32(rand.Intn(int(cfg.ShowRankEnd-myrank)+1)))

// 			if slices.Contains(rankIds, rankId) {
// 				continue
// 			}

// 			rankIds = append(rankIds, rankId)
// 		}
// 	}

// 	for _, rankId := range rankIds {
// 		rank := a.arena_ranks[rankId-1]

// 		rankInfo := rankData2Proto(int32(rankId), rank)

// 		showList = append(showList, rankInfo)
// 	}

// 	sort.Slice(showList, func(i, j int) bool {
// 		return showList[i].RankId < showList[j].RankId
// 	})

// 	return showList
// }

// func rankData2Proto(rankId int32, rank *model.ArenaRankData) *msg.ArenaRankInfo {
// 	if rank == nil {
// 		return nil
// 	}

// 	rankCfg := template.GetZombieColiseumTemplate().GetCfgByRank(rankId)
// 	if rankCfg == nil {
// 		return nil
// 	}

// 	rankInfo := &msg.ArenaRankInfo{}
// 	rankInfo.RankId = rankId
// 	rankInfo.MonsterIds = rank.MonsterIds
// 	rankInfo.IsRobot = rank.IsRobot

// 	if !rank.IsRobot {
// 		_, player := ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(rank.AccountId))
// 		if player != nil {
// 			rankInfo.PlayerInfo = ToPlayerSimpleInfo(player)
// 		}
// 	} else {

// 		robotCfg := template.GetRobotTemplate().GetZombieColiseumCfg(uint32(rank.AccountId))
// 		if robotCfg != nil {
// 			rankInfo.MonsterIds = []int32{0, 0, 0, 0, 0, 0}
// 			for pos, id := range robotCfg.GetList() {
// 				if pos > len(rankInfo.MonsterIds) {
// 					continue
// 				}

// 				rankInfo.MonsterIds[pos-1] = int32(id)
// 			}

// 			rankInfo.PlayerInfo = &msg.PlayerSimpleInfo{}
// 			rankInfo.PlayerInfo.AccountId = uint32(rank.AccountId)
// 			rankInfo.PlayerInfo.Name = template.GetRandomNameTemplate().RandOne()
// 			rankInfo.PlayerInfo.Head = robotCfg.RandHead()
// 			rankInfo.PlayerInfo.HeadFrame = robotCfg.RandHeadFrame()
// 			rankInfo.PlayerInfo.ShipId = robotCfg.RandShip()
// 			rankInfo.PlayerInfo.Combat = randRobotCombat(rankCfg.RobotCombat[0], rankCfg.RobotCombat[1])
// 		}
// 	}

// 	return rankInfo
// }

// func (a *ArenaService) Pk(playerData *common.PlayerData, rankId int32) msg.ErrCode {
// 	cur := time.Now().Unix()

// 	log.Debug("ArenaService, pk", zap.Any("rankId", rankId), zap.Any("accountId", playerData.GetAccountId()))
// 	if int(rankId) > len(a.arena_ranks) {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	{
// 		myrank := a.FindPlayerRank(playerData.GetAccountId())
// 		cfg := template.GetZombieColiseumTemplate().GetCfgByRank(myrank)
// 		if !(cfg.ShowRankStart <= rankId && rankId <= cfg.ShowRankEnd) {
// 			log.Error("ArenaService, pk list err", zap.Any("rankId", rankId))
// 			// return msg.ErrCode_ERR_NONE
// 		}
// 	}

// 	for _, shipId := range playerData.Arena.TodayUseShips {
// 		if playerData.AccountInfo.ShipId == shipId {
// 			return msg.ErrCode_ARENA_SHIP_IS_USED
// 		}
// 	}

// 	rank := a.arena_ranks[rankId-1]

// 	defenderId := rank.AccountId
// 	{
// 		defenderSession := a.arena_playerSessions[defenderId]
// 		if defenderSession != nil && cur <= defenderSession.timeout {
// 			return msg.ErrCode_ARENA_IN_PK
// 		}
// 	}

// 	todayPkCnt := playerData.Arena.TodayPkCnt
// 	todayBuyPkCnt := playerData.Arena.TodayBuyPkCnt
// 	if todayPkCnt >= template.GetSystemItemTemplate().ArenaDayPkRecoverCnt+todayBuyPkCnt {
// 		return msg.ErrCode_ARENA_IN_PK
// 	}

// 	{
// 		oldSession := a.arena_playerSessions[playerData.GetAccountId()]
// 		if oldSession != nil {
// 			if cur < oldSession.timeout {
// 				return msg.ErrCode_ARENA_IN_PK
// 			} else {
// 				a.ClearSession(oldSession)
// 			}
// 		}
// 	}

// 	session := a.arena_fightSession[rankId]
// 	if session != nil {
// 		if cur < session.timeout {
// 			return msg.ErrCode_ARENA_IN_PK
// 		} else {
// 			a.ClearSession(session)
// 		}
// 	}

// 	session = &ArenaFightSession{}
// 	session.rankId = rankId
// 	session.attackerId = playerData.GetAccountId()
// 	session.defenserId = rank.AccountId
// 	session.timeout = cur + int64(3*60)

// 	a.arena_playerSessions[session.attackerId] = session
// 	a.arena_playerSessions[session.defenserId] = session

// 	playerData.Arena.TodayPkCnt++
// 	playerData.Arena.TotalPkCnt++
// 	playerData.Arena.TodayUseShips = append(playerData.Arena.TodayUseShips, playerData.AccountInfo.ShipId)

// 	dao.ArenaDao.SavePlayerData(playerData.Arena)

// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) ClearSession(session *ArenaFightSession) {
// 	if session == nil {
// 		return
// 	}

// 	delete(a.arena_fightSession, session.rankId)
// 	delete(a.arena_playerSessions, session.attackerId)
// 	delete(a.arena_playerSessions, session.defenserId)
// }

// func (a *ArenaService) SetDefend(playerData *common.PlayerData, monsterIds []int32) msg.ErrCode {

// 	if len(monsterIds) > len(playerData.Arena.DefendMonster) {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	// if len(playerData.Arena.DefendMonster) != 6 && playerData.Arena.DefendMonster[5] <= 0 {
// 	// 	return msg.ErrCode_ERR_NONE
// 	// }

// 	slices.Sort(playerData.Arena.UnlockMonsterIds)
// 	for pos, id := range monsterIds {
// 		if id <= 0 {
// 			continue
// 		}

// 		if playerData.Arena.DefendMonster[pos] < 0 {
// 			return msg.ErrCode_ARENA_POS_IS_LCOK
// 		}

// 		if _, b := slices.BinarySearch(playerData.Arena.UnlockMonsterIds, uint32(id)); !b {
// 			return msg.ErrCode_ARENA_MONSTER_IS_LCOK
// 		}
// 	}

// 	playerData.Arena.DefendMonster = monsterIds

// 	myRankId := a.FindPlayerRank(playerData.GetAccountId())
// 	if myRankId > 0 {
// 		rank := a.arena_ranks[myRankId-1]
// 		rank.MonsterIds = playerData.Arena.DefendMonster

// 		dao.ArenaRankDao.UpdateRank(rank)
// 	}

// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) GetAllRecord(playerData *common.PlayerData) []*msg.ArenaPkRecord {
// 	records := []*msg.ArenaPkRecord{}
// 	ec, datas := dao.ArenaDao.LoadPlayerPkRecord(playerData.GetAccountId())
// 	if ec != msg.ErrCode_SUCC {
// 		return records
// 	}

// 	for _, v := range datas {
// 		record := &msg.ArenaPkRecord{
// 			Win:      v.IsWin,
// 			Stamp:    v.Stamp,
// 			IsAttack: v.IsAttack,
// 			OldRank:  v.OldRank,
// 			NewRank:  v.NewRank,
// 		}

// 		enemyInfo := &msg.PlayerSimpleInfo{}
// 		enemyInfo.Name = v.EnemyInfo.Nick
// 		enemyInfo.Title = v.EnemyInfo.Title
// 		enemyInfo.Head = v.EnemyInfo.HeadImg
// 		enemyInfo.HeadFrame = v.EnemyInfo.HeadFrame
// 		enemyInfo.ShipId = v.EnemyInfo.ShipId
// 		enemyInfo.Combat = v.EnemyInfo.Combat
// 		record.EnemyInfo = enemyInfo

// 		records = append(records, record)
// 	}

// 	return records
// }

// func (a *ArenaService) CheckBattle(playerData *common.PlayerData) msg.ErrCode {
// 	session := a.arena_playerSessions[playerData.GetAccountId()]
// 	if session == nil {
// 		return msg.ErrCode_ERR_NONE
// 	}

// 	// template.GetZombieColiseumTemplate().GetCfg()
// 	return msg.ErrCode_SUCC
// }

// func (a *ArenaService) GetFightParam(playerData *common.PlayerData) *protocol.ZombieColiseumExtra {
// 	param := &protocol.ZombieColiseumExtra{}
// 	session := a.arena_playerSessions[playerData.GetAccountId()]
// 	if session == nil {
// 		return nil
// 	}

// 	rank := a.arena_ranks[session.rankId-1]
// 	if rank == nil {
// 		return nil
// 	}

// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(session.rankId)
// 	if cfg == nil {
// 		return nil
// 	}

// 	if !rank.IsRobot {
// 		_, defenser := ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(session.defenserId))
// 		if defenser != nil {
// 			param.Combat = defenser.Combat
// 		}
// 	} else {

// 		param.Combat = randRobotCombat(cfg.RobotCombat[0], cfg.RobotCombat[1]) //combatMin + 1 + uint32(rand.Intn(int(combatMax-combatMin)+1))
// 	}

// 	param.IsRobot = rank.IsRobot
// 	param.List = make(map[int]uint32)

// 	if needUseDefultMonsters(rank.MonsterIds) {
// 		rank.MonsterIds = template.GetSystemItemTemplate().ArenaInitDefendMonster
// 	}

// 	for index, id := range rank.MonsterIds {
// 		pos := index + 1
// 		param.List[pos] = 0
// 		if id > 0 {
// 			param.List[pos] = uint32(id)
// 		}
// 	}
// 	// for i := len(rank.MonsterIds) - 1; i >= 0; i-- {
// 	// 	param.List[pos] = 0
// 	// 	if rank.MonsterIds[i] > 0 {
// 	// 		param.List[pos] = uint32(rank.MonsterIds[i])
// 	// 	}

// 	// 	pos++
// 	// }

// 	param.Id = uint32(session.defenserId)

// 	// param.List = map[int]uint32{}

// 	return param
// }

// func (a *ArenaService) GmResetData(playerData *common.PlayerData) {
// 	if playerData == nil {
// 		a.LoadPlayerData(playerData)
// 	}

// 	a.DayReset(playerData)
// 	dao.ArenaDao.SavePlayerData(playerData.Arena)
// }

// func (a *ArenaService) OnClose(playerData *common.PlayerData) {
// 	log.Debug("ArenaService OnClose, ", zap.Any("accountid", playerData.GetAccountId()))
// 	a.FightResult(playerData, false, 0, 0)
// }

// func (a *ArenaService) FightResult(playerData *common.PlayerData, isWin bool, KillCnt uint64, fightTime int64) {
// 	// isWin = true
// 	log.Debug("arena FightResult, ", zap.Int64("accountId", playerData.GetAccountId()), zap.Bool("win", isWin))

// 	session := a.arena_playerSessions[playerData.GetAccountId()]
// 	if session == nil {
// 		return
// 	}

// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(session.rankId)
// 	if cfg == nil {
// 		log.Error("arena FightResult, cfg is nil", zap.Int64("accountId", playerData.GetAccountId()), zap.Int32("rankId", session.rankId))
// 		return
// 	}

// 	rankId := session.rankId
// 	attackerId := session.attackerId
// 	defenserId := session.defenserId
// 	oldRankId := a.FindPlayerRank(attackerId)
// 	curRankId := oldRankId

// 	reward := cfg.LoseRewards
// 	if isWin {

// 		reward = cfg.WinRewards
// 		if oldRankId < 0 || rankId < oldRankId {
// 			bDefenseRobot := a.arena_ranks[rankId-1].IsRobot
// 			a.arena_ranks[rankId-1].AccountId = attackerId
// 			a.arena_ranks[rankId-1].ServerId = config.Conf.ServerId
// 			a.arena_ranks[rankId-1].IsRobot = false

// 			curRankId = rankId
// 			a.arena_playerRank[attackerId] = curRankId

// 			if 1 <= oldRankId && oldRankId <= int32(len(a.arena_ranks)+1) {
// 				a.arena_ranks[oldRankId-1].AccountId = defenserId

// 				myMonsters := a.arena_ranks[oldRankId-1].MonsterIds
// 				a.arena_ranks[oldRankId-1].MonsterIds = a.arena_ranks[rankId-1].MonsterIds
// 				a.arena_ranks[rankId-1].MonsterIds = myMonsters

// 				if needUseDefultMonsters(a.arena_ranks[oldRankId-1].MonsterIds) {
// 					a.arena_ranks[oldRankId-1].MonsterIds = template.GetSystemItemTemplate().ArenaInitDefendMonster
// 				}

// 				if needUseDefultMonsters(a.arena_ranks[rankId-1].MonsterIds) {
// 					a.arena_ranks[rankId-1].MonsterIds = template.GetSystemItemTemplate().ArenaInitDefendMonster
// 				}

// 				if bDefenseRobot {
// 					a.arena_playerRank[defenserId] = oldRankId
// 					dao.ArenaRankDao.UpdateRank(a.arena_ranks[oldRankId-1])
// 				} else {

// 				}
// 			} else {
// 				if needUseDefultMonsters(playerData.Arena.DefendMonster) {
// 					a.arena_ranks[rankId-1].MonsterIds = template.GetSystemItemTemplate().ArenaInitDefendMonster
// 				} else {
// 					a.arena_ranks[rankId-1].MonsterIds = playerData.Arena.DefendMonster
// 				}

// 				dao.ArenaRankDao.DelRank(defenserId)
// 			}

// 			dao.ArenaRankDao.UpdateRank(a.arena_ranks[rankId-1])
// 		}

// 		if rankId == 1 {
// 			bannerMsg := &msg.InterNotifyBanner{}
// 			bannerMsg.BannerType = uint32(msg.BannerType_Banner_Game)
// 			bannerMsg.Content = fmt.Sprintf("%d", template.ArenaRankTickerID)
// 			bannerMsg.Params = append(bannerMsg.Params, playerData.AccountInfo.Nick)

// 			ServMgr.GetCommonService().SendBannerMsg(bannerMsg)
// 		}
// 	}

// 	para := "lose"
// 	if isWin {
// 		para = "win"
// 	}

// 	if reward != nil {
// 		var notifyItems []uint32
// 		rwdStr := ""
// 		for _, v := range reward {
// 			rwdStr = fmt.Sprintf("%d,%d;", v.ItemId, v.ItemNum)
// 			addItems := ServMgr.GetItemService().AddItem(playerData.AccountInfo.AccountId, v.ItemId, v.ItemNum,
// 				publicconst.ArenaBuyPk, false)
// 			notifyItems = append(notifyItems, ServMgr.GetItemService().GetSimpleItemIds(addItems)...)
// 		}

// 		ServMgr.GetItemService().updateClientItemsChange(playerData.AccountInfo.AccountId, notifyItems)

// 		para = fmt.Sprintf("%s|%s", para, rwdStr)
// 	}

// 	addRecord(session, oldRankId, isWin)

// 	a.ClearSession(session)

// 	ntf := &msg.ArenaResultNtf{}
// 	ntf.Win = isWin
// 	ntf.KillCnt = uint32(KillCnt)
// 	ntf.FightTime = uint32(fightTime)
// 	ntf.OldRank = oldRankId
// 	ntf.CurRank = curRankId
// 	ntf.RewardItem = TemplateSimpleItemToProtocolSImpleItems(reward)
// 	tools.SendMsg(playerData.PlayerAgent, ntf, 0, msg.ErrCode_SUCC)

// 	ServMgr.GetCommonService().AddStaticsData(playerData, publicconst.Statics_Arena_Pk, para)

// 	ServMgr.GetFightService().LeaveFight(playerData)
// }

// func needUseDefultMonsters(monsters []int32) bool {
// 	if monsters == nil || len(monsters) == 0 {
// 		return true
// 	}

// 	for _, v := range monsters {
// 		if v > 0 {
// 			return false
// 		}
// 	}

// 	return true
// }

// func addRecord(session *ArenaFightSession, oldRankId int32, isWin bool) {
// 	if session == nil {
// 		return
// 	}

// 	_, attacker := ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(session.attackerId))
// 	if attacker == nil {
// 		return
// 	}

// 	var defenser *PlayerSimpleInfo
// 	if session.defenserId > 5000 {
// 		_, defenser = ServMgr.GetSocialService().GetPlayerSimpleInfo(uint32(session.defenserId))
// 	} else {
// 		robot := randRobot(session.rankId)
// 		defenser = &PlayerSimpleInfo{}
// 		defenser.AccountId = robot.AccountId
// 		defenser.Head = robot.Head
// 		defenser.HeadFrame = robot.HeadFrame
// 		defenser.ShipId = robot.ShipId
// 		defenser.Title = robot.Title
// 		defenser.Name = robot.Name
// 		defenser.Combat = robot.Combat
// 	}

// 	if defenser == nil {
// 		return
// 	}

// 	cur := time.Now().Unix()

// 	attackerRecord := &model.ArenaPlayerPkRecordData{}
// 	attackerRecord.IsWin = isWin
// 	attackerRecord.IsAttack = true
// 	attackerRecord.OldRank = oldRankId
// 	attackerRecord.NewRank = oldRankId
// 	if isWin && session.rankId < oldRankId {
// 		attackerRecord.NewRank = session.rankId
// 	}

// 	attackerRecord.Stamp = cur
// 	attackerRecord.EnemyInfo = PlayerSimpleInfoToEnemyInfo(defenser)
// 	dao.ArenaDao.AddPlayerPkRecord(session.attackerId, attackerRecord)

// 	if session.defenserId > 5000 {
// 		defenserRecord := &model.ArenaPlayerPkRecordData{}
// 		defenserRecord.IsWin = !isWin
// 		defenserRecord.IsAttack = false
// 		defenserRecord.OldRank = session.rankId
// 		defenserRecord.NewRank = session.rankId
// 		if isWin && session.rankId > oldRankId {
// 			defenserRecord.NewRank = oldRankId
// 		}
// 		defenserRecord.Stamp = cur
// 		defenserRecord.EnemyInfo = PlayerSimpleInfoToEnemyInfo(attacker)
// 		dao.ArenaDao.AddPlayerPkRecord(session.defenserId, defenserRecord)
// 	}

// }

// func (a *ArenaService) FindPlayerRank(target int64) int32 {

// 	rankId, ok := a.arena_playerRank[target]
// 	if !ok {
// 		return -1
// 	}
// 	// left, right := 0, len(a.ranks)-1
// 	// for left <= right {
// 	// 	mid := left + (right-left)/2
// 	// 	if a.ranks[mid].AccountId == target {
// 	// 		return int32(mid)
// 	// 	} else if a.ranks[mid].AccountId < target {
// 	// 		left = mid + 1
// 	// 	} else {
// 	// 		right = mid - 1
// 	// 	}
// 	// }

// 	return rankId
// }

// func (a *ArenaService) RankRewardSettlement() error {

// 	ec, rankDatas := dao.ArenaRankDao.LoadAllRank(config.Conf.ServerId)
// 	if ec != msg.ErrCode_SUCC {
// 		return errors.New("data not find")
// 	}

// 	log.Info("Arena RankRewardSettlement", zap.Int("num", len(rankDatas)))
// 	for _, v := range rankDatas {
// 		sendRankRewardMail(v.AccountId, uint32(v.RankId), template.ArenaRank)
// 	}

// 	return nil
// }

// func sendRankRewardMail(accountId int64, ranking uint32, rankType template.RankType) {
// 	rankReward := template.GetRankingTemplate().GetRankReward(rankType, ranking)
// 	if rankReward == nil {
// 		return
// 	}

// 	log.Info("Arena sendRankRewardMail", zap.Int64("accountId", accountId), zap.Uint32("rank", ranking))

// 	// make rank reward
// 	var items []*model.SimpleItem
// 	for m := 0; m < len(rankReward.RewardItems); m++ {
// 		items = append(items, &model.SimpleItem{
// 			Id:  rankReward.RewardItems[m].ItemId,
// 			Num: rankReward.RewardItems[m].ItemNum,
// 		})
// 	}

// 	// send mail
// 	mailConfig := template.GetMailTemplate().GetMail(rankReward.MailId)
// 	endTime := time.Now().AddDate(0, 0, 60)
// 	title := fmt.Sprintf("%v", mailConfig.Title)
// 	content := fmt.Sprintf("%v", mailConfig.Content)

// 	mail := model.NewMail(common2.GenSnowFlake(), title, content, items, uint32(endTime.Unix()))
// 	mail.MailType = 1

// 	if player := common.PlayerMgr.FindPlayerData(accountId); player != nil {
// 		InterSendUserMsg(func(msg interface{}, playerData *common.PlayerData) {
// 			ServMgr.GetMailService().AddSystemMail(playerData, mail)
// 		}, mail, player)
// 	} else {
// 		dao.MailDao.AddMail(accountId, mail)
// 	}

// }

// func randRobot(rankId int32) *msg.PlayerSimpleInfo {
// 	info := &msg.PlayerSimpleInfo{}

// 	cfg := template.GetZombieColiseumTemplate().GetCfgByRank(rankId)
// 	if cfg == nil {
// 		log.Error("arena randRobot, rank config not find ", zap.Any("rankId", rankId))
// 		return info
// 	}

// 	index := rand.Intn(len(cfg.RobotId))
// 	robotId := cfg.RobotId[index]

// 	robotCfg := template.GetRobotTemplate().GetZombieColiseumCfg(robotId)
// 	if robotCfg == nil {
// 		log.Error("arena randRobot,  ZombieColiseumCfg not find ", zap.Any("robotId", robotId))
// 		return info
// 	}

// 	info.AccountId = robotId
// 	info.Name = template.GetRandomNameTemplate().RandOne()
// 	info.Head = robotCfg.RandHead()
// 	info.HeadFrame = robotCfg.RandHeadFrame()
// 	info.ShipId = robotCfg.RandShip()
// 	info.Combat = randRobotCombat(cfg.RobotCombat[0], cfg.RobotCombat[1])
// 	return info
// }

// func randRobotCombat(min, max uint32) uint32 {

// 	if min > max {
// 		temp := min
// 		min = max
// 		max = temp
// 	}

// 	return uint32(min + 1 + uint32(rand.Intn(int(max-min)+1)))
// }

// // 计算从1970年开始的第N周（周一为每周第一天，1970年第1周包含1970-01-01）
// func weekSince1970(t time.Time) int {
// 	// 1970年1月1日是周四（Thursday）
// 	// 1970年第1周的起始日（周一）：1969-12-29
// 	firstWeekStart := time.Date(1969, 12, 29, 0, 0, 0, 0, time.UTC)

// 	// 计算当前日期与第1周起始日的时间差（天数）
// 	days := int(t.Sub(firstWeekStart).Hours() / 24)

// 	// 总周数 = 天数差 / 7（向上取整，+6是为了避免浮点数计算）
// 	totalWeeks := (days + 6) / 7

// 	return totalWeeks
// }

// func PlayerSimpleInfoToEnemyInfo(data *PlayerSimpleInfo) *model.ArenaEnemyData {
// 	if data == nil {
// 		return nil
// 	}

// 	info := &model.ArenaEnemyData{}
// 	info.AccountId = int64(data.AccountId)
// 	info.HeadFrame = data.HeadFrame
// 	info.HeadImg = data.Head
// 	info.Nick = data.Name
// 	info.ShipId = data.ShipId
// 	info.Title = data.Title
// 	info.Combat = data.Combat
// 	return info
// }
