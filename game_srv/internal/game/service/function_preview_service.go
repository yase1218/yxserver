package service

import (
	"encoding/json"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/game/condition"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"math/rand"
	"msg"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
)

var (
	fp_conditions map[msg.ConditionType][]uint32
	fp_ids        map[uint32][]*template.UnlockCondition
)

const ( // 资源本常量 对应功能模块表id
	MoneyPass                     = 104100
	EquipPass                     = 104200
	SidearmPass                   = 104300
	PetPass                       = 104400
	MaxHitCount                   = 10 // 一轮最大攻击次数
	MaxRankSize                   = 100
	LastResourcesPassResetTimeKey = "last_resources_pass_reset_time"
)

var passTypeToRankTypeMap map[int]template.RankType

var AllPassTypes = []uint32{
	MoneyPass,
	EquipPass,
	SidearmPass,
	PetPass,
}

func init() {
	passTypeToRankTypeMap = make(map[int]template.RankType)
	passTypeToRankTypeMap[MoneyPass] = template.ResourcesPassCoinRank
	passTypeToRankTypeMap[EquipPass] = template.ResourcesPassEquipRank
	passTypeToRankTypeMap[SidearmPass] = template.ResourcesPassDeputyWeaponRank
	passTypeToRankTypeMap[PetPass] = template.ResourcesPassPetRank
}

func InitFuncPreview() {
	fp_conditions = make(map[msg.ConditionType][]uint32)
	fp_ids = make(map[uint32][]*template.UnlockCondition)
	for k := range template.GetFunctionPreviewTemplate().GetAllCfg() {
		functionCfg := template.GetFunctionTemplate().GetFunction(k)
		if functionCfg == nil {
			log.Error("function cfg nil", zap.Uint32("id", k))
			continue
		}

		for _, vv := range functionCfg.Conditions {
			conditionType := msg.ConditionType(vv.Id)
			if fp_conditions[conditionType] == nil {
				fp_conditions[conditionType] = make([]uint32, 0)
			}
			fp_conditions[conditionType] = append(fp_conditions[conditionType], k)
		}

		fp_ids[k] = functionCfg.Conditions
	}
}

func UpdateFunctionPreview(p *player.Player, conditionType msg.ConditionType) {
	ntf := &msg.FunctionPreviewNtf{Data: make(map[uint32]uint32)}
	for _, v := range fp_conditions[conditionType] {
		state, ok := p.UserData.FunctionPreview.Data[v]
		if !ok {
			p.UserData.FunctionPreview.Data[v] = msg.TaskState_Task_Accept
			state = p.UserData.FunctionPreview.Data[v]
		}

		if state != msg.TaskState_Task_Accept {
			continue
		}
		finish := true
		for _, vv := range fp_ids[v] {
			_, ok := condition.GetCondition().Check(p, vv.Parse())
			if !ok {
				finish = false
				break
			}
		}
		if finish {
			p.UserData.FunctionPreview.Data[v] = msg.TaskState_Task_Complete
			ntf.Data[v] = uint32(p.UserData.FunctionPreview.Data[v])
		}
	}

	if len(ntf.Data) > 0 {
		p.SaveFunctionPreview()
		p.SendNotify(ntf)
	}
}

func FunctionPreviewReward(p *player.Player, req *msg.FunctionPreviewRewardReq, ack *msg.FunctionPreviewRewardAck) {
	functionPreviewCfg := template.GetFunctionPreviewTemplate().GetCfg(req.GetId())
	if functionPreviewCfg == nil {
		log.Error("function preview cfg nil", zap.Uint32("id", req.GetId()))
		ack.Err = msg.ErrCode_CONFIG_NIL
		return
	}

	state := p.UserData.FunctionPreview.Data[req.GetId()]
	if state != msg.TaskState_Task_Complete {
		ack.Err = msg.ErrCode_TASK_NOT_COMPLETE
		return
	}

	if functionPreviewCfg.Reward != nil {
		var notifyItems []uint32
		for _, item := range functionPreviewCfg.Reward {
			addItems := AddItem(p.GetUserId(), item.ItemId, int32(item.ItemNum), publicconst.FunctionPreviewReward, false)
			notifyItems = append(notifyItems, GetSimpleItemIds(addItems)...)
		}
		updateClientItemsChange(p.GetUserId(), notifyItems)
	}

	p.UserData.FunctionPreview.Data[req.GetId()] = msg.TaskState_Task_Done

	p.SaveFunctionPreview()

	ack.Id = req.GetId()
	ack.RewardItem = TemplateItemToProtocolItems(functionPreviewCfg.Reward)
	ack.State = uint32(p.UserData.FunctionPreview.Data[req.GetId()])
}

func GetResourcesPassBaseData(p *player.Player, passType uint32) *msg.GetResourcesPassBaseDataResp {
	if passType <= 0 && passType != MoneyPass && passType != EquipPass && passType != SidearmPass && passType != PetPass {
		log.Error("handle resouce pass data err!", zap.Uint32("passtype", passType))
		return nil
	}

	passInfo := p.UserData.ResourcesPass.PassList

	var resp = &msg.GetResourcesPassBaseDataResp{
		PassType:     passType,
		BuyCount:     0,
		HitCount:     0,
		Total:        0,
		ExpressState: false,
	}
	for _, v := range passInfo {
		if passType == v.PassType {
			resp.BuyCount = v.BuyCount
			resp.HitCount = v.HitCount
			resp.PassType = v.PassType
			resp.ExpressState = v.ExpressState
			resp.Total = v.Total
			break
		}
	}

	return resp
}

func BuyResoucesPassItem(p *player.Player, passType uint32, itemNum uint32) (msg.ErrCode, uint32) {
	if passType <= 0 || itemNum <= 0 {
		return msg.ErrCode_ResoucesPass_Type_Err, 0
	}

	if passType != MoneyPass && passType != EquipPass && passType != SidearmPass && passType != PetPass {
		log.Error("handle resource pass data err!", zap.Uint32("passtype", passType))
		return msg.ErrCode_ResoucesPass_Type_Err, 0
	}

	passInfoList := p.UserData.ResourcesPass.PassList
	if len(passInfoList) <= 0 {
		return msg.ErrCode_ResoucesPass_Data_Not_Found, 0
	}

	var res *model.ResourcesPass
	var recordIdx int
	for i := range passInfoList {
		if passInfoList[i].PassType == passType {
			res = passInfoList[i]
			recordIdx = i
			break
		}
	}

	if res == nil || recordIdx < 0 {
		return msg.ErrCode_ResoucesPass_Data_Not_Found, 0
	}

	buyAndCostListCfg := template.GetSystemItemTemplate().ResourcesPassBuyAndCostList
	var costItemId uint32 = 600002
	var getItemId uint32
	for i := range buyAndCostListCfg {
		if buyAndCostListCfg[i][0] == res.PassType {
			getItemId = buyAndCostListCfg[i][1]
			break
		}
	}

	if getItemId <= 0 {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	costNumArrCfg := template.GetSystemItemTemplate().ResourcesPassBuyPrice
	costNumLimitCfg := template.GetSystemItemTemplate().ResourcesPassBuyLimit

	if res.BuyCount >= costNumLimitCfg || (res.BuyCount+itemNum) > costNumLimitCfg {
		return msg.ErrCode_Buy_Count_Limit, 0
	}

	var buyContIdx uint32 = 0
	if res.BuyCount > 0 {
		buyContIdx = res.BuyCount - 1
	}

	if len(costNumArrCfg) <= 0 || res.BuyCount > uint32(len(costNumArrCfg)) {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	costList := make([]uint32, 0)
	for i := 0; i < int(itemNum); i++ {
		tempIdx := buyContIdx + uint32(i)
		tempCostNum := costNumArrCfg[tempIdx]
		costList = append(costList, tempCostNum)
	}

	var totalCostNum uint32
	for i := 0; i < len(costList); i++ {
		totalCostNum += costList[i]
	}

	if totalCostNum <= 0 {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	costCode := CostItem(p.GetUserId(), costItemId, totalCostNum, publicconst.ResourcesPassCost, true)
	if costCode != msg.ErrCode_SUCC {
		return msg.ErrCode_NO_ENOUGH_ITEM, 0
	}

	totalAddItemNum := itemNum * template.GetSystemItemTemplate().ResourcesPassBuyNum
	AddItem(p.GetUserId(), getItemId, int32(totalAddItemNum), publicconst.ResourcesPassBuyItem, true)

	passInfoList[recordIdx].BuyCount += itemNum
	p.UserData.ResourcesPass.PassList = passInfoList
	p.SaveAccountActivity()
	return msg.ErrCode_SUCC, passInfoList[recordIdx].BuyCount
}

func ResourcePassAttack(p *player.Player, passType uint32) (msg.ErrCode, uint32) {
	isPass := _checkIsResourcesPass(passType)
	if !isPass {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	passInfoList := p.UserData.ResourcesPass.PassList
	targetIdx := _findTargetResourcesPassIndex(passInfoList, passType)
	if targetIdx < 0 {
		return msg.ErrCode_ResoucesPass_Data_Not_Found, 0
	}

	cfg := template.GetSystemItemTemplate().ResourcesPassBuyAndCostList
	var consumeItemId uint32
	var getItemId uint32
	for i := 0; i < len(cfg); i++ {
		if passType == cfg[i][0] {
			consumeItemId = cfg[i][1]
			getItemId = cfg[i][2]
			break
		}
	}

	if consumeItemId <= 0 || getItemId <= 0 {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	passInfo := passInfoList[targetIdx]
	missions := p.UserData.Mission.Missions
	var maxMissionId int = 0
	for i := 0; i < len(missions); i++ {
		if missions[i].IsPass && missions[i].MissionId > maxMissionId {
			maxMissionId = missions[i].MissionId
		}
	}
	resStageCfg := template.GetResStageTemplate().GetMaxPassConfig(maxMissionId, passType)
	if resStageCfg == nil {
		return msg.ErrCode_CONFIG_NIL, 0
	}

	item := findItem(p, consumeItemId)
	if item == nil || item.Num <= 0 {
		return msg.ErrCode_ITEM_NOT_EXIST, 0
	}

	isExpress := passInfo.ExpressState
	totalReward := make([][]int32, 0)
	var realAttackNum uint32 = 1 // 真正攻击次数
	if isExpress {
		curAttackCount := passInfo.HitCount
		remainderCount := MaxHitCount - curAttackCount
		realAttackNum = remainderCount
		if uint32(item.Num) < remainderCount {
			realAttackNum = uint32(item.Num)
		}

		costRes := CostItem(p.GetUserId(), consumeItemId, realAttackNum, publicconst.ResourcesPassAttackCost, true)
		if costRes != msg.ErrCode_SUCC {
			return costRes, 0
		}

		for i := 0; i < int(realAttackNum); i++ {
			passInfo.HitCount++
			if passInfo.HitCount >= MaxHitCount {
				breakReward := resStageCfg.ParseBreakReward
				for i := 0; i < len(breakReward); i++ {
					totalReward = append(totalReward, breakReward[i])
				}
			} else {
				reward := _genRewardInfo(resStageCfg)
				totalReward = append(totalReward, reward)
			}

		}

	} else {
		costRes := CostItem(p.GetUserId(), consumeItemId, realAttackNum, publicconst.ResourcesPassAttackCost, true)
		if costRes != msg.ErrCode_SUCC {
			return costRes, 0
		}

		passInfo.HitCount++
		if passInfo.HitCount >= MaxHitCount {
			breakReward := resStageCfg.ParseBreakReward
			for i := 0; i < len(breakReward); i++ {
				totalReward = append(totalReward, breakReward[i])
			}
		}
		reward := _genRewardInfo(resStageCfg)
		totalReward = append(totalReward, reward)

	}

	if passInfo.HitCount >= MaxHitCount {
		passInfo.HitCount = 0
		passInfo.Total++
	}

	var totalRewardNum uint32
	for i := 0; i < len(totalReward); i++ {
		totalRewardNum += uint32(totalReward[i][1])
		AddPlayerItem(p, uint32(totalReward[i][0]), totalReward[i][1], publicconst.ResourcesPassAttackReward, true)
	}

	// _updatePlayerRankInfo(p, passType, passInfo)
	UpdateCommonRankInfo(p, passInfo.Total, passTypeToRankTypeMap[int(passType)])
	p.UserData.ResourcesPass.PassList[targetIdx] = passInfo
	p.SaveResourcesPass()
	UpdateTask(p, true, publicconst.TASK_COND_RESOURCES_PASS, realAttackNum) // 进行XX次素材本

	return msg.ErrCode_SUCC, totalRewardNum
}

func ResourcePassStateUpdate(p *player.Player, passType uint32, state bool) msg.ErrCode {
	isPass := _checkIsResourcesPass(passType)
	if !isPass {
		return msg.ErrCode_CONFIG_NIL
	}

	passInfoList := p.UserData.ResourcesPass.PassList
	passIdx := _findTargetResourcesPassIndex(passInfoList, passType)
	if passIdx < 0 {
		return msg.ErrCode_ResoucesPass_Data_Not_Found
	}

	passInfo := passInfoList[passIdx]
	passInfo.ExpressState = state
	p.SaveResourcesPass()

	return msg.ErrCode_SUCC
}

func ResourcesPassRankList(p *player.Player, passType uint32) (msg.ErrCode, []*msg.ResourcePassRankListInfo) {
	isPass := _checkIsResourcesPass(passType)
	if !isPass {
		return msg.ErrCode_CONFIG_NIL, nil
	}

	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
		key  = fmt.Sprintf("%v:%v", config.Conf.ServerId, passType)
	)

	res, err := rc.ZRevRangeByScoreWithScores(rCtx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  MaxRankSize,
	}).Result()
	if err != nil {
		return msg.ErrCode_SERVER_NOT_OPEN, nil
	}

	rankList := make([]*msg.ResourcePassRankListInfo, 0, len(res))
	for _, v := range res {
		memberStr, ok := v.Member.(string)
		if !ok {
			log.Error("member is not string type", zap.Any("member", v.Member))
			continue
		}

		var playerInfo msg.PlayerSimpleInfo
		if err := json.Unmarshal([]byte(memberStr), &playerInfo); err != nil {
			log.Error("unmarshal to player simple info err", zap.Error(err),
				zap.String("memberStr", memberStr))
			continue
		}

		data := &msg.ResourcePassRankListInfo{
			Score:      uint32(v.Score),
			PlayerInfo: &playerInfo,
		}

		rankList = append(rankList, data)
	}
	return msg.ErrCode_SUCC, rankList
}

func _genRewardInfo(cfg *template.JResStage) []int32 {
	randomNum := rand.Intn(10001)

	reward := make([]int32, 0)
	if randomNum <= int(cfg.BaseWeight) {
		for i := 0; i < len(cfg.ParseBaseReward); i++ {
			reward = append(reward, cfg.ParseBaseReward[i]...)
		}
	} else {
		for i := 0; i < len(cfg.ParseBigReward); i++ {
			reward = append(reward, cfg.ParseBigReward[i]...)
		}
	}

	return reward
}

func _checkIsResourcesPass(passType uint32) bool {
	if passType <= 0 {
		return false
	}

	if passType != MoneyPass && passType != EquipPass && passType != SidearmPass && passType != PetPass {
		log.Error("handle resource pass data err!", zap.Uint32("passtype", passType))
		return false
	}

	return true
}

func _findTargetResourcesPassIndex(list []*model.ResourcesPass, passType uint32) int {
	for i := 0; i < len(list); i++ {
		if list[i].PassType == passType {
			return i
		}
	}

	return -1
}

func _updatePlayerRankInfo(p *player.Player, passType uint32, passInfo *model.ResourcesPass) {

	var (
		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
		key  = fmt.Sprintf("%v:%v", config.Conf.ServerId, passType)
	)

	code, playerInfo := GetPlayerSimpleInfo(p.GetUserId())
	if code != msg.ErrCode_SUCC {
		return
	}
	simpleInfo := ToPlayerSimpleInfo(playerInfo)
	data, err := json.Marshal(simpleInfo)

	_, err1 := rc.ZAdd(rCtx, key, redis.Z{Score: float64(passInfo.Total), Member: data}).Result()
	if err1 != nil {
		log.Error("ZAdd player resouce pass rank err", zap.Error(err), zap.Uint32("passType", passType))
		return
	}

	rc.ZRemRangeByRank(rCtx, key, MaxRankSize, -1)
}

// func ClearPlayerRankInfo() {
// 	now := time.Now()
// 	if now.Weekday() != time.Sunday || now.Hour() != 0 {
// 		return
// 	}

// 	var (
// 		rc   = rdb_single.Get()
// 		rCtx = rdb_single.GetCtx()
// 	)

// 	lastResetTimeStr, err := rc.Get(rCtx, LastResourcesPassResetTimeKey).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			lastResetTimeStr = "0"
// 		} else {
// 			log.Error("Get last resource pass reset time error", zap.Error(err))
// 			return
// 		}
// 	}

// 	var lastResetTime int64
// 	if lastResetTimeStr == "" {
// 		lastResetTime = 0
// 	} else {
// 		lastResetTime, err = strconv.ParseInt(lastResetTimeStr, 10, 64)
// 		if err != nil {
// 			log.Error("reset resources pass time trans error",
// 				zap.String("lastResetTimeStr", lastResetTimeStr),
// 				zap.Error(err))
// 			lastResetTime = 0
// 		}
// 	}

// 	currentTime := time.Now().Unix()
// 	if utils.IsSameDay(currentTime, lastResetTime) {
// 		return
// 	}

// 	cleanSuccess := true
// 	for i := 0; i < len(AllPassTypes); i++ {
// 		key := fmt.Sprintf("%v:%v", config.Conf.ServerId, AllPassTypes[i])
// 		_, err := rc.ZRemRangeByRank(rCtx, key, 0, -1).Result()
// 		if err != nil {
// 			log.Error("Clear player resource pass rank err",
// 				zap.Error(err),
// 				zap.Uint32("passType", AllPassTypes[i]))
// 			cleanSuccess = false
// 			continue
// 		}

// 		log.Info("Clear player rank successfully",
// 			zap.Uint32("passType", AllPassTypes[i]),
// 			zap.String("key", key))
// 	}

// 	if cleanSuccess {
// 		_, err := rc.Set(rCtx, LastResourcesPassResetTimeKey, currentTime, 0).Result()
// 		if err != nil {
// 			log.Error("set resources pass rank reset time error", zap.Error(err))
// 		} else {
// 			log.Info("Reset time updated successfully", zap.Int64("resetTime", currentTime))
// 		}
// 	}
// }

func handleResetPlayerResourcePass(p *player.Player) {
	now := time.Now()
	resetHour := template.GetSystemItemTemplate().RefreshHour
	if now.Hour() >= int(resetHour) {
		p.UserData.ResourcesPass.LastResetTime = time.Now()
		for _, v := range p.UserData.ResourcesPass.PassList {
			v.BuyCount = 0
		}

		p.SaveResourcesPass()
	}
}
