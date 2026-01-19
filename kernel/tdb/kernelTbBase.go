package tdb

var fileInfos = []FileInfo{

	{"1.0-角色-机甲_ship.xlsx", []SheetInfo{
		{SheetName: "ship_list", Initer: MapLoader("ShipListCfgs", "Id"), ObjPropType: ShipListCfg{}},
	}},
	{"1.1-角色-装备_equip.xlsx", []SheetInfo{
		{SheetName: "equip", Initer: MapLoader("EquipCfgs", "Id"), ObjPropType: EquipCfg{}},
	}},
	{"1.3-角色-怪物_monster.xlsx", []SheetInfo{
		{SheetName: "monster_list", Initer: MapLoader("MonsterListCfgs", "Id"), ObjPropType: MonsterListCfg{}},
	}},
	{"1.4-角色-交互物_construction.xlsx", []SheetInfo{
		{SheetName: "construction_list", Initer: MapLoader("ConstructionListCfgs", "Id"), ObjPropType: ConstructionListCfg{}},
	}},
	{"10-属性表_attr_list.xlsx", []SheetInfo{
		{SheetName: "attr_list", Initer: MapLoader("AttrListAttrListCfgs", "Id"), ObjPropType: AttrListAttrListCfg{}},
	}},
	{"101-功能模块_function.xlsx", []SheetInfo{
		{SheetName: "function", Initer: MapLoader("FunctionFunctionCfgs", "Id"), ObjPropType: FunctionFunctionCfg{}},
	}},
	{"102-活动总表-ActivityMethod.xlsx", []SheetInfo{
		{SheetName: "ActivityMethod", Initer: MapLoader("ActivityMethodCfgs", "Id"), ObjPropType: ActivityMethodCfg{}},
	}},
	{"103-充值-Charge.xlsx", []SheetInfo{
		{SheetName: "Charge", Initer: MapLoader("ChargeCfgs", "Id"), ObjPropType: ChargeCfg{}},
	}},
	{"104-屏蔽字.xlsx", []SheetInfo{
		{SheetName: "forbidden", Initer: MapLoader("ForbiddenCfgs", "Id"), ObjPropType: ForbiddenCfg{}},
	}},
	{"105-广告-Advert.xlsx", []SheetInfo{
		{SheetName: "Advert", Initer: MapLoader("AdvertCfgs", "Id"), ObjPropType: AdvertCfg{}},
	}},
	{"106-机器人-robot.xlsx", []SheetInfo{
		{SheetName: "robot", Initer: MapLoader("RobotCfgs", "Id"), ObjPropType: RobotCfg{}},
	}},
	{"107-随机名字-RandomName.xlsx", []SheetInfo{
		{SheetName: "RandomName", Initer: MapLoader("RandomNameCfgs", "Id"), ObjPropType: RandomNameCfg{}},
	}},
	{"11-通用设置_gameOption.xlsx", []SheetInfo{
		{SheetName: "gameOption", Initer: MapLoader("GameOptionGameOptionCfgs", "Id"), ObjPropType: GameOptionGameOptionCfg{}},
	}},
	{"12-物品系统_item.xlsx", []SheetInfo{
		{SheetName: "item", Initer: MapLoader("ItemItemCfgs", "Id"), ObjPropType: ItemItemCfg{}},
	}},
	{"13.1-关卡配置_stage.xlsx", []SheetInfo{
		{SheetName: "stage", Initer: MapLoader("StageStageCfgs", "Id"), ObjPropType: StageStageCfg{}},
	}},
	{"13.2-地图配置_map.xlsx", []SheetInfo{
		{SheetName: "map", Initer: MapLoader("MapMapCfgs", "Id"), ObjPropType: MapMapCfg{}},
	}},
	{"13.4-战场效果表_stage_buff.xlsx", []SheetInfo{
		{SheetName: "stage_buff", Initer: MapLoader("StageBuffStageBuffCfgs", "Id"), ObjPropType: StageBuffStageBuffCfg{}},
	}},
	{"16-任务系统-task.xlsx", []SheetInfo{
		{SheetName: "task", Initer: MapLoader("TaskCfgs", "Id"), ObjPropType: TaskCfg{}},
	}},
	{"17-活跃度-Activation.xlsx", []SheetInfo{
		{SheetName: "Activation", Initer: MapLoader("ActivationCfgs", "Id"), ObjPropType: ActivationCfg{}},
	}},
	{"18-商店-Shop.xlsx", []SheetInfo{
		{SheetName: "shop", Initer: MapLoader("ShopCfgs", "Id"), ObjPropType: ShopCfg{}},
	}},
	{"2.0-宠物_Pet.xlsx", []SheetInfo{
		{SheetName: "Pet", Initer: MapLoader("PetPetCfgs", "Id"), ObjPropType: PetPetCfg{}},
	}},
	{"2.1-宠物_蛋_PetEgg.xlsx", []SheetInfo{
		{SheetName: "PetEgg", Initer: MapLoader("PetEggPetEggCfgs", "Id"), ObjPropType: PetEggPetEggCfg{}},
	}},
	{"2.2-宠物-时装_PetClothes.xlsx", []SheetInfo{
		{SheetName: "PetClothes", Initer: MapLoader("PetClothesCfgs", "Id"), ObjPropType: PetClothesCfg{}},
	}},
	{"2.3-宠物-套装_PetSuit.xlsx", []SheetInfo{
		{SheetName: "PetSuit", Initer: MapLoader("PetSuitCfgs", "Id"), ObjPropType: PetSuitCfg{}},
	}},
	{"2.4-宠物_觉醒_PetEvolution.xlsx", []SheetInfo{
		{SheetName: "PetEvolution", Initer: MapLoader("PetEvolutionPetEvolutionCfgs", "Level"), ObjPropType: PetEvolutionPetEvolutionCfg{}},
	}},
	{"2.5-宠物_洗炼_PetRefresh.xlsx", []SheetInfo{
		{SheetName: "PetRefresh", Initer: MapLoader("PetRefreshPetRefreshCfgs", "Id"), ObjPropType: PetRefreshPetRefreshCfg{}},
	}},
	{"21-关卡配置-关卡宝箱-RuleBox.xlsx", []SheetInfo{
		{SheetName: "RuleBox", Initer: MapLoader("RuleBoxCfgs", "Id"), ObjPropType: RuleBoxCfg{}},
	}},
	{"22-关卡配置-怪物刷新_stage_monster_refresh.xlsx", []SheetInfo{
		{SheetName: "stage_monster_refresh", Initer: MapLoader("StageMonsterRefreshCfgs", "Id"), ObjPropType: StageMonsterRefreshCfg{}},
	}},
	{"3-技能表_skill_list.xlsx", []SheetInfo{
		{SheetName: "skill_list", Initer: MapLoader("SkillListSkillListCfgs", "Id"), ObjPropType: SkillListSkillListCfg{}},
	}},
	{"3.1-技能组合表_skillCombo.xlsx", []SheetInfo{
		{SheetName: "skillCombo", Initer: MapLoader("SkillComboSkillComboCfgs", "Id"), ObjPropType: SkillComboSkillComboCfg{}},
	}},
	{"31-玩家等级-PlayerLevel.xlsx", []SheetInfo{
		{SheetName: "PlayerLevel", Initer: MapLoader("PlayerLevelCfgs", "Id"), ObjPropType: PlayerLevelCfg{}},
	}},
	{"33-随机表-LootPool.xlsx", []SheetInfo{
		{SheetName: "lootPool", Initer: MapLoader("LootPoolCfgs", "Id"), ObjPropType: LootPoolCfg{}},
	}},
	{"34-玩法总表_playingMethod.xlsx", []SheetInfo{
		{SheetName: "playingMethod", Initer: MapLoader("PlayingMethodPlayingMethodCfgs", "Id"), ObjPropType: PlayingMethodPlayingMethodCfg{}},
	}},
	{"39-武器-Weapon.xlsx", []SheetInfo{
		{SheetName: "Weapon", Initer: MapLoader("WeaponCfgs", "Id"), ObjPropType: WeaponCfg{}},
	}},
	{"40-角色等级-RoleLevel.xlsx", []SheetInfo{
		{SheetName: "RoleLevel", Initer: MapLoader("RoleLevelCfgs", "Id"), ObjPropType: RoleLevelCfg{}},
	}},
	{"41-角色升星-RoleRank.xlsx", []SheetInfo{
		{SheetName: "RoleRank", Initer: MapLoader("RoleRankCfgs", "Id"), ObjPropType: RoleRankCfg{}},
	}},
	{"43-装备等级-equip_level.xlsx", []SheetInfo{
		{SheetName: "equip_level", Initer: MapLoader("EquipLevelCfgs", "Id"), ObjPropType: EquipLevelCfg{}},
	}},
	{"44-武器升级-WeaponLevel.xlsx", []SheetInfo{
		{SheetName: "WeaponLevel", Initer: MapLoader("WeaponLevelCfgs", "Id"), ObjPropType: WeaponLevelCfg{}},
	}},
	{"45-武器库-WeaponLibrary.xlsx", []SheetInfo{
		{SheetName: "WeaponLibrary", Initer: MapLoader("WeaponLibraryCfgs", "Id"), ObjPropType: WeaponLibraryCfg{}},
	}},
	{"46-支援属性-support_value.xlsx", []SheetInfo{
		{SheetName: "support_value", Initer: MapLoader("SupportValueCfgs", "Id"), ObjPropType: SupportValueCfg{}},
	}},
	{"48-机甲招募-lottery_ship.xlsx", []SheetInfo{
		{SheetName: "LotteryShip", Initer: MapLoader("LotteryShipCfgs", "Id"), ObjPropType: LotteryShipCfg{}},
	}},
	{"49-七日登录-sevenday_login.xlsx", []SheetInfo{
		{SheetName: "SevendayLogin", Initer: MapLoader("SevendayLoginCfgs", "Id"), ObjPropType: SevendayLoginCfg{}},
	}},
	{"5-特效表_effect_list.xlsx", []SheetInfo{
		{SheetName: "effect_list", Initer: MapLoader("EffectListEffectListCfgs", "Id"), ObjPropType: EffectListEffectListCfg{}},
	}},
	{"50-邮件-mail.xlsx", []SheetInfo{
		{SheetName: "mail", Initer: MapLoader("MailCfgs", "Id"), ObjPropType: MailCfg{}},
	}},
	{"51-月卡-MonthlyCard.xlsx", []SheetInfo{
		{SheetName: "MonthlyCard", Initer: MapLoader("MonthlyCardCfgs", "Id"), ObjPropType: MonthlyCardCfg{}},
	}},
	{"52.0-关卡基金-StageFund.xlsx", []SheetInfo{
		{SheetName: "StageFund", Initer: MapLoader("StageFundCfgs", "Id"), ObjPropType: StageFundCfg{}},
	}},
	{"52.1-签到基金-SignFund.xlsx", []SheetInfo{
		{SheetName: "SignFund", Initer: MapLoader("SignFundCfgs", "Id"), ObjPropType: SignFundCfg{}},
	}},
	{"53.0-活跃战令-ActivePass.xlsx", []SheetInfo{
		{SheetName: "ActivePass", Initer: MapLoader("ActivePassCfgs", "Id"), ObjPropType: ActivePassCfg{}},
	}},
	{"53.1-任务战令-TaskPass.xlsx", []SheetInfo{
		{SheetName: "TaskPass", Initer: MapLoader("TaskPassCfgs", "Id"), ObjPropType: TaskPassCfg{}},
	}},
	{"56-开服活动-SevendaysTask.xlsx", []SheetInfo{
		{SheetName: "SevendaysTask", Initer: MapLoader("SevendaysTaskCfgs", "Id"), ObjPropType: SevendaysTaskCfg{}},
	}},
	{"57-领取体力-GetiAp.xlsx", []SheetInfo{
		{SheetName: "GetiAp", Initer: MapLoader("GetiApCfgs", "Id"), ObjPropType: GetiApCfg{}},
	}},
	{"58-天赋-Talent.xlsx", []SheetInfo{
		{SheetName: "Talent", Initer: MapLoader("TalentCfgs", "Id"), ObjPropType: TalentCfg{}},
	}},
	{"59-外观-Appearance.xlsx", []SheetInfo{
		{SheetName: "Appearance", Initer: MapLoader("AppearanceCfgs", "Id"), ObjPropType: AppearanceCfg{}},
	}},
	{"6-配件表_skillpro_list.xlsx", []SheetInfo{
		{SheetName: "skillpro_list", Initer: MapLoader("SkillproListSkillproListCfgs", "Id"), ObjPropType: SkillproListSkillproListCfg{}},
	}},
	{"60-巅峰挑战_BattlePass.xlsx", []SheetInfo{
		{SheetName: "BattlePass", Initer: MapLoader("BattlePassBattlePassCfgs", "Id"), ObjPropType: BattlePassBattlePassCfg{}},
	}},
	{"61-巅峰挑战匹配_BattleMatch.xlsx", []SheetInfo{
		{SheetName: "BattleMatch", Initer: MapLoader("BattleMatchBattleMatchCfgs", "Id"), ObjPropType: BattleMatchBattleMatchCfg{}},
	}},
	{"62-巅峰道具_BattleItem.xlsx", []SheetInfo{
		{SheetName: "BattleItem", Initer: MapLoader("BattleItemBattleItemCfgs", "Id"), ObjPropType: BattleItemBattleItemCfg{}},
	}},
	{"63-表情_Emote.xlsx", []SheetInfo{
		{SheetName: "Emote", Initer: MapLoader("EmoteEmoteCfgs", "Id"), ObjPropType: EmoteEmoteCfg{}},
	}},
	{"64-奇遇_Adventure.xlsx", []SheetInfo{
		{SheetName: "Adventure", Initer: MapLoader("AdventureAdventureCfgs", "Id"), ObjPropType: AdventureAdventureCfg{}},
	}},
	{"65.0-排行榜_Rank.xlsx", []SheetInfo{
		{SheetName: "Rank", Initer: MapLoader("RankRankCfgs", "Id"), ObjPropType: RankRankCfg{}},
	}},
	{"65.1-排行榜-首通奖励_RankReward.xlsx", []SheetInfo{
		{SheetName: "RankReward", Initer: MapLoader("RankRewardCfgs", "Id"), ObjPropType: RankRewardCfg{}},
	}},
	{"65.2-排行榜-排名奖励_RankingReward.xlsx", []SheetInfo{
		{SheetName: "RankingReward", Initer: MapLoader("RankingRewardCfgs", "Id"), ObjPropType: RankingRewardCfg{}},
	}},
	{"66.0-联盟_Guild.xlsx", []SheetInfo{
		{SheetName: "Guild", Initer: MapLoader("GuildGuildCfgs", "Id"), ObjPropType: GuildGuildCfg{}},
	}},
	{"67.0-冒险地图_AdvMap.xlsx", []SheetInfo{
		{SheetName: "AdvMap", Initer: MapLoader("AdvMapAdvMapCfgs", "Id"), ObjPropType: AdvMapAdvMapCfg{}},
	}},
	{"67.1-冒险路点_AdvPoint.xlsx", []SheetInfo{
		{SheetName: "AdvPoint", Initer: MapLoader("AdvPointAdvPointCfgs", "Id"), ObjPropType: AdvPointAdvPointCfg{}},
	}},
	{"67.2-冒险卡牌_AdvCard.xlsx", []SheetInfo{
		{SheetName: "AdvCard", Initer: MapLoader("AdvCardAdvCardCfgs", "Id"), ObjPropType: AdvCardAdvCardCfg{}},
	}},
	{"67.3-冒险卡包_AdvCardBag.xlsx", []SheetInfo{
		{SheetName: "AdvCardBag", Initer: MapLoader("AdvCardBagAdvCardBagCfgs", "Id"), ObjPropType: AdvCardBagAdvCardBagCfg{}},
	}},
	{"67.4-冒险怪物战力_AdvValue.xlsx", []SheetInfo{
		{SheetName: "AdvValue", Initer: MapLoader("AdvValueAdvValueCfgs", "Id"), ObjPropType: AdvValueAdvValueCfg{}},
	}},
	{"69-疯狂合约-Contract.xlsx", []SheetInfo{
		{SheetName: "Contract", Initer: MapLoader("ContractCfgs", "Id"), ObjPropType: ContractCfg{}},
	}},
	{"70-广播-Broadcast.xlsx", []SheetInfo{
		{SheetName: "broadcast", Initer: MapLoader("BroadcastCfgs", "Id"), ObjPropType: BroadcastCfg{}},
	}},
	{"71-红包_RedPacket.xlsx", []SheetInfo{
		{SheetName: "RedPacket", Initer: MapLoader("RedPacketRedPacketCfgs", "Id"), ObjPropType: RedPacketRedPacketCfg{}},
	}},
	{"72-战外地图_OutMap.xlsx", []SheetInfo{
		{SheetName: "OutMap", Initer: MapLoader("OutMapOutMapCfgs", "Id"), ObjPropType: OutMapOutMapCfg{}},
	}},
	{"73-僵尸竞技场_ZombieColiseum.xlsx", []SheetInfo{
		{SheetName: "ZombieColiseum", Initer: MapLoader("ZombieColiseumZombieColiseumCfgs", "Id"), ObjPropType: ZombieColiseumZombieColiseumCfg{}},
	}},
	{"74-图鉴_Handbook.xlsx", []SheetInfo{
		{SheetName: "Handbook", Initer: MapLoader("HandbookHandbookCfgs", "Id"), ObjPropType: HandbookHandbookCfg{}},
	}},
	{"98-帮助小说明-HelpTips.xlsx", []SheetInfo{
		{SheetName: "HelpTips", Initer: MapLoader("HelpTipsCfgs", "Id"), ObjPropType: HelpTipsCfg{}},
	}},
	{"99-多语言表-language.xlsx", []SheetInfo{
		{SheetName: "Language", Initer: MapLoader("LanguageCfgs", "Id"), ObjPropType: LanguageCfg{}},
	}},
}

type TableBase struct {
	ShipListCfgs                     map[int]*ShipListCfg
	EquipCfgs                        map[int]*EquipCfg
	MonsterListCfgs                  map[int]*MonsterListCfg
	ConstructionListCfgs             map[int]*ConstructionListCfg
	AttrListAttrListCfgs             map[int]*AttrListAttrListCfg
	FunctionFunctionCfgs             map[int]*FunctionFunctionCfg
	ActivityMethodCfgs               map[int]*ActivityMethodCfg
	ChargeCfgs                       map[int]*ChargeCfg
	ForbiddenCfgs                    map[int]*ForbiddenCfg
	AdvertCfgs                       map[int]*AdvertCfg
	RobotCfgs                        map[int]*RobotCfg
	RandomNameCfgs                   map[int]*RandomNameCfg
	GameOptionGameOptionCfgs         map[int]*GameOptionGameOptionCfg
	ItemItemCfgs                     map[int]*ItemItemCfg
	StageStageCfgs                   map[int]*StageStageCfg
	MapMapCfgs                       map[int]*MapMapCfg
	StageBuffStageBuffCfgs           map[int]*StageBuffStageBuffCfg
	TaskCfgs                         map[int]*TaskCfg
	ActivationCfgs                   map[int]*ActivationCfg
	ShopCfgs                         map[int]*ShopCfg
	PetPetCfgs                       map[int]*PetPetCfg
	PetEggPetEggCfgs                 map[int]*PetEggPetEggCfg
	PetClothesCfgs                   map[int]*PetClothesCfg
	PetSuitCfgs                      map[int]*PetSuitCfg
	PetEvolutionPetEvolutionCfgs     map[int]*PetEvolutionPetEvolutionCfg
	PetRefreshPetRefreshCfgs         map[int]*PetRefreshPetRefreshCfg
	RuleBoxCfgs                      map[int]*RuleBoxCfg
	StageMonsterRefreshCfgs          map[int]*StageMonsterRefreshCfg
	SkillListSkillListCfgs           map[int]*SkillListSkillListCfg
	SkillComboSkillComboCfgs         map[int]*SkillComboSkillComboCfg
	PlayerLevelCfgs                  map[int]*PlayerLevelCfg
	LootPoolCfgs                     map[int]*LootPoolCfg
	PlayingMethodPlayingMethodCfgs   map[int]*PlayingMethodPlayingMethodCfg
	WeaponCfgs                       map[int]*WeaponCfg
	RoleLevelCfgs                    map[int]*RoleLevelCfg
	RoleRankCfgs                     map[int]*RoleRankCfg
	EquipLevelCfgs                   map[int]*EquipLevelCfg
	WeaponLevelCfgs                  map[int]*WeaponLevelCfg
	WeaponLibraryCfgs                map[int]*WeaponLibraryCfg
	SupportValueCfgs                 map[int]*SupportValueCfg
	LotteryShipCfgs                  map[int]*LotteryShipCfg
	SevendayLoginCfgs                map[int]*SevendayLoginCfg
	EffectListEffectListCfgs         map[int]*EffectListEffectListCfg
	MailCfgs                         map[int]*MailCfg
	MonthlyCardCfgs                  map[int]*MonthlyCardCfg
	StageFundCfgs                    map[int]*StageFundCfg
	SignFundCfgs                     map[int]*SignFundCfg
	ActivePassCfgs                   map[int]*ActivePassCfg
	TaskPassCfgs                     map[int]*TaskPassCfg
	SevendaysTaskCfgs                map[int]*SevendaysTaskCfg
	GetiApCfgs                       map[int]*GetiApCfg
	TalentCfgs                       map[int]*TalentCfg
	AppearanceCfgs                   map[int]*AppearanceCfg
	SkillproListSkillproListCfgs     map[int]*SkillproListSkillproListCfg
	BattlePassBattlePassCfgs         map[int]*BattlePassBattlePassCfg
	BattleMatchBattleMatchCfgs       map[int]*BattleMatchBattleMatchCfg
	BattleItemBattleItemCfgs         map[int]*BattleItemBattleItemCfg
	EmoteEmoteCfgs                   map[int]*EmoteEmoteCfg
	AdventureAdventureCfgs           map[int]*AdventureAdventureCfg
	RankRankCfgs                     map[int]*RankRankCfg
	RankRewardCfgs                   map[int]*RankRewardCfg
	RankingRewardCfgs                map[int]*RankingRewardCfg
	GuildGuildCfgs                   map[int]*GuildGuildCfg
	AdvMapAdvMapCfgs                 map[int]*AdvMapAdvMapCfg
	AdvPointAdvPointCfgs             map[int]*AdvPointAdvPointCfg
	AdvCardAdvCardCfgs               map[int]*AdvCardAdvCardCfg
	AdvCardBagAdvCardBagCfgs         map[int]*AdvCardBagAdvCardBagCfg
	AdvValueAdvValueCfgs             map[int]*AdvValueAdvValueCfg
	ContractCfgs                     map[int]*ContractCfg
	BroadcastCfgs                    map[int]*BroadcastCfg
	RedPacketRedPacketCfgs           map[int]*RedPacketRedPacketCfg
	OutMapOutMapCfgs                 map[int]*OutMapOutMapCfg
	ZombieColiseumZombieColiseumCfgs map[int]*ZombieColiseumZombieColiseumCfg
	HandbookHandbookCfgs             map[int]*HandbookHandbookCfg
	HelpTipsCfgs                     map[int]*HelpTipsCfg
	LanguageCfgs                     map[int]*LanguageCfg
}

func GetShipListCfg(Id int) *ShipListCfg {
	return tdb.ShipListCfgs[Id]
}

func RangShipListCfgs(f func(conf *ShipListCfg) bool) {
	for _, v := range tdb.ShipListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetEquipCfg(Id int) *EquipCfg {
	return tdb.EquipCfgs[Id]
}

func RangEquipCfgs(f func(conf *EquipCfg) bool) {
	for _, v := range tdb.EquipCfgs {
		if !f(v) {
			return
		}
	}
}

func GetMonsterListCfg(Id int) *MonsterListCfg {
	return tdb.MonsterListCfgs[Id]
}

func RangMonsterListCfgs(f func(conf *MonsterListCfg) bool) {
	for _, v := range tdb.MonsterListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetConstructionListCfg(Id int) *ConstructionListCfg {
	return tdb.ConstructionListCfgs[Id]
}

func RangConstructionListCfgs(f func(conf *ConstructionListCfg) bool) {
	for _, v := range tdb.ConstructionListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAttrListAttrListCfg(Id int) *AttrListAttrListCfg {
	return tdb.AttrListAttrListCfgs[Id]
}

func RangAttrListAttrListCfgs(f func(conf *AttrListAttrListCfg) bool) {
	for _, v := range tdb.AttrListAttrListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetFunctionFunctionCfg(Id int) *FunctionFunctionCfg {
	return tdb.FunctionFunctionCfgs[Id]
}

func RangFunctionFunctionCfgs(f func(conf *FunctionFunctionCfg) bool) {
	for _, v := range tdb.FunctionFunctionCfgs {
		if !f(v) {
			return
		}
	}
}

func GetActivityMethodCfg(Id int) *ActivityMethodCfg {
	return tdb.ActivityMethodCfgs[Id]
}

func RangActivityMethodCfgs(f func(conf *ActivityMethodCfg) bool) {
	for _, v := range tdb.ActivityMethodCfgs {
		if !f(v) {
			return
		}
	}
}

func GetChargeCfg(Id int) *ChargeCfg {
	return tdb.ChargeCfgs[Id]
}

func RangChargeCfgs(f func(conf *ChargeCfg) bool) {
	for _, v := range tdb.ChargeCfgs {
		if !f(v) {
			return
		}
	}
}

func GetForbiddenCfg(Id int) *ForbiddenCfg {
	return tdb.ForbiddenCfgs[Id]
}

func RangForbiddenCfgs(f func(conf *ForbiddenCfg) bool) {
	for _, v := range tdb.ForbiddenCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvertCfg(Id int) *AdvertCfg {
	return tdb.AdvertCfgs[Id]
}

func RangAdvertCfgs(f func(conf *AdvertCfg) bool) {
	for _, v := range tdb.AdvertCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRobotCfg(Id int) *RobotCfg {
	return tdb.RobotCfgs[Id]
}

func RangRobotCfgs(f func(conf *RobotCfg) bool) {
	for _, v := range tdb.RobotCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRandomNameCfg(Id int) *RandomNameCfg {
	return tdb.RandomNameCfgs[Id]
}

func RangRandomNameCfgs(f func(conf *RandomNameCfg) bool) {
	for _, v := range tdb.RandomNameCfgs {
		if !f(v) {
			return
		}
	}
}

func GetGameOptionGameOptionCfg(Id int) *GameOptionGameOptionCfg {
	return tdb.GameOptionGameOptionCfgs[Id]
}

func RangGameOptionGameOptionCfgs(f func(conf *GameOptionGameOptionCfg) bool) {
	for _, v := range tdb.GameOptionGameOptionCfgs {
		if !f(v) {
			return
		}
	}
}

func GetItemItemCfg(Id int) *ItemItemCfg {
	return tdb.ItemItemCfgs[Id]
}

func RangItemItemCfgs(f func(conf *ItemItemCfg) bool) {
	for _, v := range tdb.ItemItemCfgs {
		if !f(v) {
			return
		}
	}
}

func GetStageStageCfg(Id int) *StageStageCfg {
	return tdb.StageStageCfgs[Id]
}

func RangStageStageCfgs(f func(conf *StageStageCfg) bool) {
	for _, v := range tdb.StageStageCfgs {
		if !f(v) {
			return
		}
	}
}

func GetMapMapCfg(Id int) *MapMapCfg {
	return tdb.MapMapCfgs[Id]
}

func RangMapMapCfgs(f func(conf *MapMapCfg) bool) {
	for _, v := range tdb.MapMapCfgs {
		if !f(v) {
			return
		}
	}
}

func GetStageBuffStageBuffCfg(Id int) *StageBuffStageBuffCfg {
	return tdb.StageBuffStageBuffCfgs[Id]
}

func RangStageBuffStageBuffCfgs(f func(conf *StageBuffStageBuffCfg) bool) {
	for _, v := range tdb.StageBuffStageBuffCfgs {
		if !f(v) {
			return
		}
	}
}

func GetTaskCfg(Id int) *TaskCfg {
	return tdb.TaskCfgs[Id]
}

func RangTaskCfgs(f func(conf *TaskCfg) bool) {
	for _, v := range tdb.TaskCfgs {
		if !f(v) {
			return
		}
	}
}

func GetActivationCfg(Id int) *ActivationCfg {
	return tdb.ActivationCfgs[Id]
}

func RangActivationCfgs(f func(conf *ActivationCfg) bool) {
	for _, v := range tdb.ActivationCfgs {
		if !f(v) {
			return
		}
	}
}

func GetShopCfg(Id int) *ShopCfg {
	return tdb.ShopCfgs[Id]
}

func RangShopCfgs(f func(conf *ShopCfg) bool) {
	for _, v := range tdb.ShopCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetPetCfg(Id int) *PetPetCfg {
	return tdb.PetPetCfgs[Id]
}

func RangPetPetCfgs(f func(conf *PetPetCfg) bool) {
	for _, v := range tdb.PetPetCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetEggPetEggCfg(Id int) *PetEggPetEggCfg {
	return tdb.PetEggPetEggCfgs[Id]
}

func RangPetEggPetEggCfgs(f func(conf *PetEggPetEggCfg) bool) {
	for _, v := range tdb.PetEggPetEggCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetClothesCfg(Id int) *PetClothesCfg {
	return tdb.PetClothesCfgs[Id]
}

func RangPetClothesCfgs(f func(conf *PetClothesCfg) bool) {
	for _, v := range tdb.PetClothesCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetSuitCfg(Id int) *PetSuitCfg {
	return tdb.PetSuitCfgs[Id]
}

func RangPetSuitCfgs(f func(conf *PetSuitCfg) bool) {
	for _, v := range tdb.PetSuitCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetEvolutionPetEvolutionCfg(Level int) *PetEvolutionPetEvolutionCfg {
	return tdb.PetEvolutionPetEvolutionCfgs[Level]
}

func RangPetEvolutionPetEvolutionCfgs(f func(conf *PetEvolutionPetEvolutionCfg) bool) {
	for _, v := range tdb.PetEvolutionPetEvolutionCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPetRefreshPetRefreshCfg(Id int) *PetRefreshPetRefreshCfg {
	return tdb.PetRefreshPetRefreshCfgs[Id]
}

func RangPetRefreshPetRefreshCfgs(f func(conf *PetRefreshPetRefreshCfg) bool) {
	for _, v := range tdb.PetRefreshPetRefreshCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRuleBoxCfg(Id int) *RuleBoxCfg {
	return tdb.RuleBoxCfgs[Id]
}

func RangRuleBoxCfgs(f func(conf *RuleBoxCfg) bool) {
	for _, v := range tdb.RuleBoxCfgs {
		if !f(v) {
			return
		}
	}
}

func GetStageMonsterRefreshCfg(Id int) *StageMonsterRefreshCfg {
	return tdb.StageMonsterRefreshCfgs[Id]
}

func RangStageMonsterRefreshCfgs(f func(conf *StageMonsterRefreshCfg) bool) {
	for _, v := range tdb.StageMonsterRefreshCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSkillListSkillListCfg(Id int) *SkillListSkillListCfg {
	return tdb.SkillListSkillListCfgs[Id]
}

func RangSkillListSkillListCfgs(f func(conf *SkillListSkillListCfg) bool) {
	for _, v := range tdb.SkillListSkillListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSkillComboSkillComboCfg(Id int) *SkillComboSkillComboCfg {
	return tdb.SkillComboSkillComboCfgs[Id]
}

func RangSkillComboSkillComboCfgs(f func(conf *SkillComboSkillComboCfg) bool) {
	for _, v := range tdb.SkillComboSkillComboCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPlayerLevelCfg(Id int) *PlayerLevelCfg {
	return tdb.PlayerLevelCfgs[Id]
}

func RangPlayerLevelCfgs(f func(conf *PlayerLevelCfg) bool) {
	for _, v := range tdb.PlayerLevelCfgs {
		if !f(v) {
			return
		}
	}
}

func GetLootPoolCfg(Id int) *LootPoolCfg {
	return tdb.LootPoolCfgs[Id]
}

func RangLootPoolCfgs(f func(conf *LootPoolCfg) bool) {
	for _, v := range tdb.LootPoolCfgs {
		if !f(v) {
			return
		}
	}
}

func GetPlayingMethodPlayingMethodCfg(Id int) *PlayingMethodPlayingMethodCfg {
	return tdb.PlayingMethodPlayingMethodCfgs[Id]
}

func RangPlayingMethodPlayingMethodCfgs(f func(conf *PlayingMethodPlayingMethodCfg) bool) {
	for _, v := range tdb.PlayingMethodPlayingMethodCfgs {
		if !f(v) {
			return
		}
	}
}

func GetWeaponCfg(Id int) *WeaponCfg {
	return tdb.WeaponCfgs[Id]
}

func RangWeaponCfgs(f func(conf *WeaponCfg) bool) {
	for _, v := range tdb.WeaponCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRoleLevelCfg(Id int) *RoleLevelCfg {
	return tdb.RoleLevelCfgs[Id]
}

func RangRoleLevelCfgs(f func(conf *RoleLevelCfg) bool) {
	for _, v := range tdb.RoleLevelCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRoleRankCfg(Id int) *RoleRankCfg {
	return tdb.RoleRankCfgs[Id]
}

func RangRoleRankCfgs(f func(conf *RoleRankCfg) bool) {
	for _, v := range tdb.RoleRankCfgs {
		if !f(v) {
			return
		}
	}
}

func GetEquipLevelCfg(Id int) *EquipLevelCfg {
	return tdb.EquipLevelCfgs[Id]
}

func RangEquipLevelCfgs(f func(conf *EquipLevelCfg) bool) {
	for _, v := range tdb.EquipLevelCfgs {
		if !f(v) {
			return
		}
	}
}

func GetWeaponLevelCfg(Id int) *WeaponLevelCfg {
	return tdb.WeaponLevelCfgs[Id]
}

func RangWeaponLevelCfgs(f func(conf *WeaponLevelCfg) bool) {
	for _, v := range tdb.WeaponLevelCfgs {
		if !f(v) {
			return
		}
	}
}

func GetWeaponLibraryCfg(Id int) *WeaponLibraryCfg {
	return tdb.WeaponLibraryCfgs[Id]
}

func RangWeaponLibraryCfgs(f func(conf *WeaponLibraryCfg) bool) {
	for _, v := range tdb.WeaponLibraryCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSupportValueCfg(Id int) *SupportValueCfg {
	return tdb.SupportValueCfgs[Id]
}

func RangSupportValueCfgs(f func(conf *SupportValueCfg) bool) {
	for _, v := range tdb.SupportValueCfgs {
		if !f(v) {
			return
		}
	}
}

func GetLotteryShipCfg(Id int) *LotteryShipCfg {
	return tdb.LotteryShipCfgs[Id]
}

func RangLotteryShipCfgs(f func(conf *LotteryShipCfg) bool) {
	for _, v := range tdb.LotteryShipCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSevendayLoginCfg(Id int) *SevendayLoginCfg {
	return tdb.SevendayLoginCfgs[Id]
}

func RangSevendayLoginCfgs(f func(conf *SevendayLoginCfg) bool) {
	for _, v := range tdb.SevendayLoginCfgs {
		if !f(v) {
			return
		}
	}
}

func GetEffectListEffectListCfg(Id int) *EffectListEffectListCfg {
	return tdb.EffectListEffectListCfgs[Id]
}

func RangEffectListEffectListCfgs(f func(conf *EffectListEffectListCfg) bool) {
	for _, v := range tdb.EffectListEffectListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetMailCfg(Id int) *MailCfg {
	return tdb.MailCfgs[Id]
}

func RangMailCfgs(f func(conf *MailCfg) bool) {
	for _, v := range tdb.MailCfgs {
		if !f(v) {
			return
		}
	}
}

func GetMonthlyCardCfg(Id int) *MonthlyCardCfg {
	return tdb.MonthlyCardCfgs[Id]
}

func RangMonthlyCardCfgs(f func(conf *MonthlyCardCfg) bool) {
	for _, v := range tdb.MonthlyCardCfgs {
		if !f(v) {
			return
		}
	}
}

func GetStageFundCfg(Id int) *StageFundCfg {
	return tdb.StageFundCfgs[Id]
}

func RangStageFundCfgs(f func(conf *StageFundCfg) bool) {
	for _, v := range tdb.StageFundCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSignFundCfg(Id int) *SignFundCfg {
	return tdb.SignFundCfgs[Id]
}

func RangSignFundCfgs(f func(conf *SignFundCfg) bool) {
	for _, v := range tdb.SignFundCfgs {
		if !f(v) {
			return
		}
	}
}

func GetActivePassCfg(Id int) *ActivePassCfg {
	return tdb.ActivePassCfgs[Id]
}

func RangActivePassCfgs(f func(conf *ActivePassCfg) bool) {
	for _, v := range tdb.ActivePassCfgs {
		if !f(v) {
			return
		}
	}
}

func GetTaskPassCfg(Id int) *TaskPassCfg {
	return tdb.TaskPassCfgs[Id]
}

func RangTaskPassCfgs(f func(conf *TaskPassCfg) bool) {
	for _, v := range tdb.TaskPassCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSevendaysTaskCfg(Id int) *SevendaysTaskCfg {
	return tdb.SevendaysTaskCfgs[Id]
}

func RangSevendaysTaskCfgs(f func(conf *SevendaysTaskCfg) bool) {
	for _, v := range tdb.SevendaysTaskCfgs {
		if !f(v) {
			return
		}
	}
}

func GetGetiApCfg(Id int) *GetiApCfg {
	return tdb.GetiApCfgs[Id]
}

func RangGetiApCfgs(f func(conf *GetiApCfg) bool) {
	for _, v := range tdb.GetiApCfgs {
		if !f(v) {
			return
		}
	}
}

func GetTalentCfg(Id int) *TalentCfg {
	return tdb.TalentCfgs[Id]
}

func RangTalentCfgs(f func(conf *TalentCfg) bool) {
	for _, v := range tdb.TalentCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAppearanceCfg(Id int) *AppearanceCfg {
	return tdb.AppearanceCfgs[Id]
}

func RangAppearanceCfgs(f func(conf *AppearanceCfg) bool) {
	for _, v := range tdb.AppearanceCfgs {
		if !f(v) {
			return
		}
	}
}

func GetSkillproListSkillproListCfg(Id int) *SkillproListSkillproListCfg {
	return tdb.SkillproListSkillproListCfgs[Id]
}

func RangSkillproListSkillproListCfgs(f func(conf *SkillproListSkillproListCfg) bool) {
	for _, v := range tdb.SkillproListSkillproListCfgs {
		if !f(v) {
			return
		}
	}
}

func GetBattlePassBattlePassCfg(Id int) *BattlePassBattlePassCfg {
	return tdb.BattlePassBattlePassCfgs[Id]
}

func RangBattlePassBattlePassCfgs(f func(conf *BattlePassBattlePassCfg) bool) {
	for _, v := range tdb.BattlePassBattlePassCfgs {
		if !f(v) {
			return
		}
	}
}

func GetBattleMatchBattleMatchCfg(Id int) *BattleMatchBattleMatchCfg {
	return tdb.BattleMatchBattleMatchCfgs[Id]
}

func RangBattleMatchBattleMatchCfgs(f func(conf *BattleMatchBattleMatchCfg) bool) {
	for _, v := range tdb.BattleMatchBattleMatchCfgs {
		if !f(v) {
			return
		}
	}
}

func GetBattleItemBattleItemCfg(Id int) *BattleItemBattleItemCfg {
	return tdb.BattleItemBattleItemCfgs[Id]
}

func RangBattleItemBattleItemCfgs(f func(conf *BattleItemBattleItemCfg) bool) {
	for _, v := range tdb.BattleItemBattleItemCfgs {
		if !f(v) {
			return
		}
	}
}

func GetEmoteEmoteCfg(Id int) *EmoteEmoteCfg {
	return tdb.EmoteEmoteCfgs[Id]
}

func RangEmoteEmoteCfgs(f func(conf *EmoteEmoteCfg) bool) {
	for _, v := range tdb.EmoteEmoteCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdventureAdventureCfg(Id int) *AdventureAdventureCfg {
	return tdb.AdventureAdventureCfgs[Id]
}

func RangAdventureAdventureCfgs(f func(conf *AdventureAdventureCfg) bool) {
	for _, v := range tdb.AdventureAdventureCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRankRankCfg(Id int) *RankRankCfg {
	return tdb.RankRankCfgs[Id]
}

func RangRankRankCfgs(f func(conf *RankRankCfg) bool) {
	for _, v := range tdb.RankRankCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRankRewardCfg(Id int) *RankRewardCfg {
	return tdb.RankRewardCfgs[Id]
}

func RangRankRewardCfgs(f func(conf *RankRewardCfg) bool) {
	for _, v := range tdb.RankRewardCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRankingRewardCfg(Id int) *RankingRewardCfg {
	return tdb.RankingRewardCfgs[Id]
}

func RangRankingRewardCfgs(f func(conf *RankingRewardCfg) bool) {
	for _, v := range tdb.RankingRewardCfgs {
		if !f(v) {
			return
		}
	}
}

func GetGuildGuildCfg(Id int) *GuildGuildCfg {
	return tdb.GuildGuildCfgs[Id]
}

func RangGuildGuildCfgs(f func(conf *GuildGuildCfg) bool) {
	for _, v := range tdb.GuildGuildCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvMapAdvMapCfg(Id int) *AdvMapAdvMapCfg {
	return tdb.AdvMapAdvMapCfgs[Id]
}

func RangAdvMapAdvMapCfgs(f func(conf *AdvMapAdvMapCfg) bool) {
	for _, v := range tdb.AdvMapAdvMapCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvPointAdvPointCfg(Id int) *AdvPointAdvPointCfg {
	return tdb.AdvPointAdvPointCfgs[Id]
}

func RangAdvPointAdvPointCfgs(f func(conf *AdvPointAdvPointCfg) bool) {
	for _, v := range tdb.AdvPointAdvPointCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvCardAdvCardCfg(Id int) *AdvCardAdvCardCfg {
	return tdb.AdvCardAdvCardCfgs[Id]
}

func RangAdvCardAdvCardCfgs(f func(conf *AdvCardAdvCardCfg) bool) {
	for _, v := range tdb.AdvCardAdvCardCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvCardBagAdvCardBagCfg(Id int) *AdvCardBagAdvCardBagCfg {
	return tdb.AdvCardBagAdvCardBagCfgs[Id]
}

func RangAdvCardBagAdvCardBagCfgs(f func(conf *AdvCardBagAdvCardBagCfg) bool) {
	for _, v := range tdb.AdvCardBagAdvCardBagCfgs {
		if !f(v) {
			return
		}
	}
}

func GetAdvValueAdvValueCfg(Id int) *AdvValueAdvValueCfg {
	return tdb.AdvValueAdvValueCfgs[Id]
}

func RangAdvValueAdvValueCfgs(f func(conf *AdvValueAdvValueCfg) bool) {
	for _, v := range tdb.AdvValueAdvValueCfgs {
		if !f(v) {
			return
		}
	}
}

func GetContractCfg(Id int) *ContractCfg {
	return tdb.ContractCfgs[Id]
}

func RangContractCfgs(f func(conf *ContractCfg) bool) {
	for _, v := range tdb.ContractCfgs {
		if !f(v) {
			return
		}
	}
}

func GetBroadcastCfg(Id int) *BroadcastCfg {
	return tdb.BroadcastCfgs[Id]
}

func RangBroadcastCfgs(f func(conf *BroadcastCfg) bool) {
	for _, v := range tdb.BroadcastCfgs {
		if !f(v) {
			return
		}
	}
}

func GetRedPacketRedPacketCfg(Id int) *RedPacketRedPacketCfg {
	return tdb.RedPacketRedPacketCfgs[Id]
}

func RangRedPacketRedPacketCfgs(f func(conf *RedPacketRedPacketCfg) bool) {
	for _, v := range tdb.RedPacketRedPacketCfgs {
		if !f(v) {
			return
		}
	}
}

func GetOutMapOutMapCfg(Id int) *OutMapOutMapCfg {
	return tdb.OutMapOutMapCfgs[Id]
}

func RangOutMapOutMapCfgs(f func(conf *OutMapOutMapCfg) bool) {
	for _, v := range tdb.OutMapOutMapCfgs {
		if !f(v) {
			return
		}
	}
}

func GetZombieColiseumZombieColiseumCfg(Id int) *ZombieColiseumZombieColiseumCfg {
	return tdb.ZombieColiseumZombieColiseumCfgs[Id]
}

func RangZombieColiseumZombieColiseumCfgs(f func(conf *ZombieColiseumZombieColiseumCfg) bool) {
	for _, v := range tdb.ZombieColiseumZombieColiseumCfgs {
		if !f(v) {
			return
		}
	}
}

func GetHandbookHandbookCfg(Id int) *HandbookHandbookCfg {
	return tdb.HandbookHandbookCfgs[Id]
}

func RangHandbookHandbookCfgs(f func(conf *HandbookHandbookCfg) bool) {
	for _, v := range tdb.HandbookHandbookCfgs {
		if !f(v) {
			return
		}
	}
}

func GetHelpTipsCfg(Id int) *HelpTipsCfg {
	return tdb.HelpTipsCfgs[Id]
}

func RangHelpTipsCfgs(f func(conf *HelpTipsCfg) bool) {
	for _, v := range tdb.HelpTipsCfgs {
		if !f(v) {
			return
		}
	}
}

func GetLanguageCfg(Id int) *LanguageCfg {
	return tdb.LanguageCfgs[Id]
}

func RangLanguageCfgs(f func(conf *LanguageCfg) bool) {
	for _, v := range tdb.LanguageCfgs {
		if !f(v) {
			return
		}
	}
}

type ShipListCfg struct {
	Id                int         `col:"Id" client:"Id"`                               // 角色ID
	Rarity            int         `col:"rarity" client:"rarity"`                       // 稀有度
	ShipDebris        IntMap      `col:"ShipDebris" client:"ShipDebris"`               // 机甲激活碎片
	Shiprepeat        IntMap      `col:"Shiprepeat" client:"Shiprepeat"`               // 机甲抽重转化
	Skill             IntSlice    `col:"skill" client:"skill"`                         // 技能
	Skillpro          IntSlice    `col:"skillpro" client:"skillpro"`                   // 初始配件
	InitialAttributes IntFloatMap `col:"InitialAttributes" client:"InitialAttributes"` // 初始属性
	BornDuration      int         `col:"BornDuration" client:"BornDuration"`           // 出生持续时间
	Show              int         `col:"show" client:"show"`                           // 机甲是否显示
}

type EquipCfg struct {
	Id                int         `col:"Id" client:"Id"`                               // 系统ID
	Rarity            int         `col:"rarity" client:"rarity"`                       // 稀有度
	SmallStage        int         `col:"smallStage" client:"smallStage"`               // 小阶段
	RankCondition     IntSlice    `col:"rankCondition" client:"rankCondition"`         // 升阶条件
	Rank              int         `col:"rank" client:"rank"`                           // 装备升阶结果
	OpenRank          int         `col:"openRank" client:"openRank"`                   // 升阶开放
	Position          int         `col:"position" client:"position"`                   // 装备位置
	LevelMax          int         `col:"levelMax" client:"levelMax"`                   // 等级上限
	InitialAttributes IntFloatMap `col:"InitialAttributes" client:"InitialAttributes"` // 初始属性
	LevelAttributes   IntFloatMap `col:"LevelAttributes" client:"LevelAttributes"`     // 升级属性
	Decompose         IntSlice2   `col:"decompose" client:"decompose"`                 // 分解产出
	OutfitsID         int         `col:"outfitsID" client:"outfitsID"`                 // 套装判定
	OutfitsAttribute  FloatSlice2 `col:"outfitsAttribute" client:"outfitsAttribute"`   // 套装属性
	RewardResources   IntMap      `col:"RewardResources" client:"RewardResources"`     // 装备品质奖励
	SkillAttributes   int         `col:"SkillAttributes" client:"SkillAttributes"`     // 套装图鉴技能
}

type MonsterListCfg struct {
	Id                int         `col:"Id" client:"Id"`                               // 索引
	Type              int         `col:"type" client:"type"`                           // 类型
	InitialAttributes IntFloatMap `col:"InitialAttributes" client:"InitialAttributes"` // 初始属性
	Skill             IntSlice    `col:"skill" client:"skill"`                         // 技能ID
	RelationBuff      IntSlice    `col:"relationBuff" client:"relationBuff"`           // 附加buff
	DeadRefresh       IntSlice    `col:"DeadRefresh" client:"DeadRefresh"`             // 死亡召唤
	DeadSkill         IntSlice    `col:"DeadSkill" client:"DeadSkill"`                 // 死亡技能
	BornDuration      int         `col:"BornDuration" client:"BornDuration"`           // 出生持续时间
	Rewards           IntSlice    `col:"rewards" client:"rewards"`                     // 掉落奖励
}

type ConstructionListCfg struct {
	Id            int      `col:"Id" client:"Id"`                       // 索引
	Group         int      `col:"Group" client:"Group"`                 // 交互物组
	Type          int      `col:"type" client:"type"`                   // 脚本类型
	ScriptKey     IntSlice `col:"ScriptKey" client:"ScriptKey"`         // 脚本参数
	ScriptUpgrade IntSlice `col:"ScriptUpgrade" client:"ScriptUpgrade"` // 升级参数
}

type AttrListAttrListCfg struct {
	Id            int      `col:"Id" client:"Id"`                       // 属性ID
	Total         int      `col:"total" client:"total"`                 // 是否是总值
	Formulapara   IntSlice `col:"formulapara" client:"formulapara"`     // 公式参数-战外
	Formulatype   int      `col:"formulatype" client:"formulatype"`     // 公式类型-战外
	Formulapara2  IntSlice `col:"formulapara2" client:"formulapara2"`   // 公式参数-战内
	Formulatype2  int      `col:"formulatype2" client:"formulatype2"`   // 公式类型-战内
	IsDisplayed   int      `col:"isDisplayed" client:"isDisplayed"`     // 是否显示
	DisplayPara   IntSlice `col:"displayPara" client:"displayPara"`     // 显示参数
	Effectivetype int      `col:"effectivetype" client:"effectivetype"` // 生效范围类型
	Valuetype     int      `col:"valuetype" client:"valuetype"`         // 数值类型(1-固定值, 2-百分比)
	Initialvalue  float64  `col:"initialvalue" client:"initialvalue"`   // 初始数值
	Value         float64  `col:"value" client:"value"`                 // 战斗力换算系数
}

type FunctionFunctionCfg struct {
	Id        int    `col:"Id" client:"Id"`               // id(不可修改)
	Condition IntMap `col:"condition" client:"condition"` // 开启条件
}

type ActivityMethodCfg struct {
	Id           int       `col:"Id" client:"Id"`                     // ID
	ActivityType int       `col:"ActivityType" client:"ActivityType"` // 活动类型
	Condition    IntSlice2 `col:"condition" client:"condition"`       // 开启条件（满足条件活动就开启）
	Param        IntSlice2 `col:"param" client:"param"`               // 参数
}

type ChargeCfg struct {
	Id                          int       `col:"Id" client:"Id"`                                                   // ID
	Type                        int       `col:"Type" client:"Type"`                                               // 充值类型
	CostRMB                     int       `col:"CostRMB" client:"CostRMB"`                                         // 购买所需人民币
	ItemID                      IntSlice2 `col:"ItemID" client:"ItemID"`                                           // 充值获得
	FirstPurchaseExtra          IntSlice2 `col:"FirstPurchaseExtra" client:"FirstPurchaseExtra"`                   // 首次购买额外奖励
	FirstPurchaseExtraResetTime int       `col:"FirstPurchaseExtraResetTime" client:"FirstPurchaseExtraResetTime"` // 首次奖励重置时间
	PurchaseExtra               IntSlice2 `col:"PurchaseExtra" client:"PurchaseExtra"`                             // 购买额外奖励
}

type ForbiddenCfg struct {
	Id   int    `col:"Id" client:"Id"`     // 唯一ID
	Text string `col:"text" client:"text"` // 文字
}

type AdvertCfg struct {
	Id       int       `col:"Id" client:"Id"`             // ID
	Type     int       `col:"Type" client:"Type"`         // 类型
	ItemID   IntSlice2 `col:"ItemID" client:"ItemID"`     // 奖励
	Time     int       `col:"Time" client:"Time"`         // 广告最小时间
	MaxTimes int       `col:"MaxTimes" client:"MaxTimes"` // 最大次数
}

type RobotCfg struct {
	Id        int      `col:"Id" client:"Id"`               // ID
	Type      int      `col:"Type" client:"Type"`           // 类型
	Head      IntSlice `col:"Head" client:"Head"`           // 头像
	HeadFrame IntSlice `col:"HeadFrame" client:"HeadFrame"` // 头像框
	Ship      IntSlice `col:"Ship" client:"Ship"`           // 出战机甲
	Param     string   `col:"Param" client:"Param"`         // 参数
	PlayParam string   `col:"PlayParam" client:"PlayParam"` // 参数2
}

type RandomNameCfg struct {
	Id   int    `col:"Id" client:"Id"`     // ID
	Name string `col:"Name" client:"Name"` // 名字
}

type GameOptionGameOptionCfg struct {
	Id    int    `col:"Id" client:"Id"`       // ID
	Value string `col:"value" client:"value"` // 数值
}

type ItemItemCfg struct {
	Id         int      `col:"Id" client:"Id"`                 // 编号
	Type       int      `col:"type" client:"type"`             // 物品类型
	Use        int      `col:"use" client:"use"`               // 物品获取时行为
	Quality    int      `col:"quality" client:"quality"`       // 品质
	Price      int      `col:"price" client:"price"`           // 出售价格
	EffectType int      `col:"effectType" client:"effectType"` // 效果类型
	EffectArgs IntSlice `col:"effectArgs" client:"effectArgs"` // 效果参数1
}

type StageStageCfg struct {
	Id              int         `col:"Id" client:"Id"`                           // 索引
	NextId          int         `col:"nextId" client:"nextId"`                   // 下一关id
	Type            int         `col:"type" client:"type"`                       // 关卡类型
	MainLevelPoint  string      `col:"mainLevelPoint" client:"mainLevelPoint"`   // 打点数据文件
	StageId         int         `col:"stageId" client:"stageId"`                 // 关卡子ID
	Unlock          IntMap      `col:"unlock" client:"unlock"`                   // 解锁条件
	ShowPower       int         `col:"ShowPower" client:"ShowPower"`             // 推荐战力
	Charper         int         `col:"charper" client:"charper"`                 // 章节
	PowerCost       int         `col:"powerCost" client:"powerCost"`             // 体力消耗
	Condition       IntMap      `col:"condition" client:"condition"`             // 解锁条件
	WeatherType     IntSlice    `col:"weatherType" client:"weatherType"`         // 天气
	StageBuff       IntSlice2   `col:"stageBuff" client:"stageBuff"`             // 关卡BUFF
	Box             IntSlice2   `col:"box" client:"box"`                         // 奖励箱
	StageSummonID   IntSlice    `col:"stageSummonID" client:"stageSummonID"`     // 交互物组id
	MaximumMonster  int         `col:"MaximumMonster" client:"MaximumMonster"`   // 刷怪上限
	Clearancetion   IntSlice2   `col:"clearancetion" client:"clearancetion"`     // 通关条件
	Failure         IntSlice2   `col:"failure" client:"failure"`                 // 失败条件
	RewardComplete  IntSlice    `col:"rewardComplete" client:"rewardComplete"`   // 关卡宝箱
	RewardQuit      IntFloatMap `col:"rewardQuit" client:"rewardQuit"`           // 关卡未通关奖励(每分钟)
	RewardQuitLimit int         `col:"rewardQuitLimit" client:"rewardQuitLimit"` // 未通关奖励上限(分钟)
	ReconnectLimit  float64     `col:"reconnectLimit" client:"reconnectLimit"`   // 断线重连上限
	Reward          IntMap      `col:"reward" client:"reward"`                   // 关卡通关奖励
	PatrolReward    IntFloatMap `col:"patrolReward" client:"patrolReward"`       // 挂机奖励
	HpRate          float64     `col:"hpRate" client:"hpRate"`                   // 怪物血量系数
	DefRate         float64     `col:"defRate" client:"defRate"`                 // 怪物防御系数
	AttackRate      float64     `col:"attackRate" client:"attackRate"`           // 怪物攻击系数
	TimeHpRate      FloatSlice2 `col:"timeHpRate" client:"timeHpRate"`           // 怪物波次血量系数
	TimeAtkRate     FloatSlice2 `col:"timeAtkRate" client:"timeAtkRate"`         // 怪物波次攻击系数
	WeaponID        IntSlice    `col:"weaponID" client:"weaponID"`               // 解锁武器ID
	UnlockSkillpro  IntSlice    `col:"unlockSkillpro" client:"unlockSkillpro"`   // 解锁秘宝
	UniqueSkillpro  IntSlice    `col:"uniqueSkillpro" client:"uniqueSkillpro"`   // 专用秘宝
	UnlockPoker     IntSlice    `col:"unlockPoker" client:"unlockPoker"`         // 解锁扑克
	UniquePoker     IntSlice    `col:"uniquePoker" client:"uniquePoker"`         // 专用扑克
	BossFuryTime    int         `col:"bossFuryTime" client:"bossFuryTime"`       // boss狂暴时间
	BossFuryBuff    IntSlice    `col:"bossFuryBuff" client:"bossFuryBuff"`       // boss狂暴赋予BUFF
	SkillproShop    int         `col:"skillproShop" client:"skillproShop"`       // 配件商店
	ExtraReward     IntMap      `col:"ExtraReward" client:"ExtraReward"`         // 通关额外奖励
	ProgressReward  IntSlice2   `col:"ProgressReward" client:"ProgressReward"`   // BOSS伤害奖励
}

type MapMapCfg struct {
	Id             int      `col:"Id" client:"Id"`                         // 地图索引
	MapGroup       int      `col:"mapGroup" client:"mapGroup"`             // 地图组
	Weight         int      `col:"Weight" client:"Weight"`                 // 权重
	Type           int      `col:"Type" client:"Type"`                     // 房间类型
	ConnectionType int      `col:"connectionType" client:"connectionType"` // 拼接方向
	BornPoint      IntSlice `col:"bornPoint" client:"bornPoint"`           // 出生点
}

type StageBuffStageBuffCfg struct {
	Id            int `col:"Id" client:"Id"`                       // ID
	Weight        int `col:"Weight" client:"Weight"`               // 权重
	ShipID        int `col:"ShipID" client:"ShipID"`               // 机甲
	ShipBuffID    int `col:"ShipBuffID" client:"ShipBuffID"`       // 机甲BUFF
	ShipMonsterID int `col:"ShipMonsterID" client:"ShipMonsterID"` // 小怪BUFF
}

type TaskCfg struct {
	Id             int      `col:"Id" client:"Id"`                         // 任务ID
	TaskType       int      `col:"taskType" client:"taskType"`             // 任务刷新规则
	TaskCondition  int      `col:"taskCondition" client:"taskCondition"`   // 任务条件
	Param          int      `col:"param" client:"param"`                   // 参数
	TaskMax        int      `col:"taskMax" client:"taskMax"`               // 任务进度条最大值
	Data           int      `col:"data" client:"data"`                     // 是否继承前置数据
	HistoryData    int      `col:"historyData" client:"historyData"`       // 是否读取历史数据
	FrontTask      int      `col:"frontTask" client:"frontTask"`           // 前置任务
	FrontStage     int      `col:"frontStage" client:"frontStage"`         // 前置关卡
	FrontPlayLevel int      `col:"frontPlayLevel" client:"frontPlayLevel"` // 玩家等级
	DayActivation  int      `col:"DayActivation" client:"DayActivation"`   // 日活跃度
	WeekActivation int      `col:"WeekActivation" client:"WeekActivation"` // 周活跃度
	PassExp        int      `col:"PassExp" client:"PassExp"`               // 提供的经验
	TaskReward     IntMap   `col:"taskReward" client:"taskReward"`         // 任务奖励
	Effect1        IntSlice `col:"effect1" client:"effect1"`               // 条件参数1
}

type ActivationCfg struct {
	Id          int    `col:"Id" client:"Id"`                   // ID
	RefreshType int    `col:"RefreshType" client:"RefreshType"` // 刷新类型
	Activations int    `col:"Activations" client:"Activations"` // 活跃度数值
	Reward      IntMap `col:"Reward" client:"Reward"`           // 奖励
}

type ShopCfg struct {
	Id               int       `col:"Id" client:"Id"`                             // 商品ID
	Type             int       `col:"type" client:"type"`                         // 商品类型
	ItemID           IntSlice2 `col:"ItemID" client:"ItemID"`                     // 关联道具表
	ShopType         int       `col:"shopType" client:"shopType"`                 // 归属商店
	CostItem         IntSlice2 `col:"costItem" client:"costItem"`                 // 购买所需道具
	ChargeID         int       `col:"ChargeID" client:"ChargeID"`                 // 充值ID
	CostAd           int       `col:"CostAd" client:"CostAd"`                     // 广告购买ID
	Discount         int       `col:"discount" client:"discount"`                 // 商品折扣
	Renovate         IntSlice  `col:"renovate" client:"renovate"`                 // 刷新机制
	Limited          int       `col:"limited" client:"limited"`                   // 限购数量
	UnlockLevel      int       `col:"unlockLevel" client:"unlockLevel"`           // 解锁条件-关卡
	Unlockrank       int       `col:"unlockrank" client:"unlockrank"`             // 解锁条件-等级
	Prepose          int       `col:"prepose" client:"prepose"`                   // 解锁条件-前置
	Edition          int       `col:"edition" client:"edition"`                   // 刷新校验
	UnlockGuildLevel int       `col:"unlockGuildLevel" client:"unlockGuildLevel"` // 解锁条件-联盟等级
	RedPacket        int       `col:"RedPacket" client:"RedPacket"`               // 红包
}

type PetPetCfg struct {
	Id             int         `col:"Id" client:"Id"`                         // 宠物ID
	BaseId         int         `col:"baseId" client:"baseId"`                 // 基类ID
	SuitId         int         `col:"SuitId" client:"SuitId"`                 // 套装ID
	Time           int         `col:"Time" client:"Time"`                     // 体型变化时间
	NextId         int         `col:"NextId" client:"NextId"`                 // 成长后ID
	LevelReward    IntSlice    `col:"LevelReward" client:"LevelReward"`       // 升级奖励
	AdvReward      int         `col:"AdvReward" client:"AdvReward"`           // 冒险奖励
	Attr           IntFloatMap `col:"Attr" client:"Attr"`                     // 基础属性
	SkillAttrNum   IntFloatMap `col:"SkillAttrNum" client:"SkillAttrNum"`     // 职业属性
	StartSkill     int         `col:"StartSkill" client:"StartSkill"`         // 初始技能
	SkillNum       int         `col:"SkillNum" client:"SkillNum"`             // 拥有技能数量
	PetRefreshItem IntSlice2   `col:"PetRefreshItem" client:"PetRefreshItem"` // 洗练道具
}

type PetEggPetEggCfg struct {
	Id    int `col:"Id" client:"Id"`       // 宠物蛋ID
	Time  int `col:"Time" client:"Time"`   // 孵化时间
	PetId int `col:"PetId" client:"PetId"` // 宠物ID
}

type PetClothesCfg struct {
	Id        int      `col:"Id" client:"Id"`               // ID
	Type      int      `col:"Type" client:"Type"`           // 类型
	BaseId    int      `col:"baseId" client:"baseId"`       // 基类ID
	Condition int      `col:"Condition" client:"Condition"` // 激活条件
	Parm      IntSlice `col:"Parm" client:"Parm"`           // 条件参数
	Charm     int      `col:"Charm" client:"Charm"`         // 魅力值
}

type PetSuitCfg struct {
	Id       int      `col:"Id" client:"Id"`             // 套装ID
	BaseId   int      `col:"baseId" client:"baseId"`     // 基类ID
	PetModle IntSlice `col:"PetModle" client:"PetModle"` // 模型组建
}

type PetEvolutionPetEvolutionCfg struct {
	Level         int       `col:"Level" client:"Level"`                 // 进化等级
	EvolutionItem IntSlice2 `col:"EvolutionItem" client:"EvolutionItem"` // 进化道具
	Rate          int       `col:"Rate" client:"Rate"`                   // 成功概率
	AttrNum       IntSlice  `col:"AttrNum" client:"AttrNum"`             // 解锁职业属性
	SkillNum      int       `col:"SkillNum" client:"SkillNum"`           // 增加技能数量
}

type PetRefreshPetRefreshCfg struct {
	Id            int       `col:"Id" client:"Id"`                       // 序号
	PetId         int       `col:"PetId" client:"PetId"`                 // 宠物ID
	SkillInterval IntSlice2 `col:"SkillInterval" client:"SkillInterval"` // 随机技能
}

type RuleBoxCfg struct {
	Id      int       `col:"Id" client:"Id"`           // 索引
	TypeID  IntSlice2 `col:"typeID" client:"typeID"`   // 类型,参数
	IteamID IntMap    `col:"iteamID" client:"iteamID"` // 道具ID
}

type StageMonsterRefreshCfg struct {
	Id                 int      `col:"Id" client:"Id"`                                 // id
	StageId            int      `col:"StageId" client:"StageId"`                       // 关卡id
	RefreshType        int      `col:"RefreshType" client:"RefreshType"`               // 刷新类型
	RefreshStart       float64  `col:"RefreshStart" client:"RefreshStart"`             // 刷新开始时间单位：秒
	RefreshEnd         float64  `col:"RefreshEnd" client:"RefreshEnd"`                 // 刷新结束时间单位：秒
	RefreshDelay       float64  `col:"RefreshDelay" client:"RefreshDelay"`             // 刷新延迟单位：秒
	MonsterID          IntMap   `col:"MonsterID" client:"MonsterID"`                   // 刷新怪物id
	RefreshRange       IntSlice `col:"RefreshRange" client:"RefreshRange"`             // 刷新距离
	RefreshAngle       IntSlice `col:"RefreshAngle" client:"RefreshAngle"`             // 刷新角度
	RefreshNum         int      `col:"RefreshNum" client:"RefreshNum"`                 // 刷新数量
	RefreshMaxTimes    int      `col:"RefreshMaxTimes" client:"RefreshMaxTimes"`       // 刷新最大次数
	RefreshInterval    float64  `col:"RefreshInterval" client:"RefreshInterval"`       // 刷新间隔
	RefreshCD          float64  `col:"RefreshCD" client:"RefreshCD"`                   // 刷新触发间隔单位：秒
	RefreshAddInterval int      `col:"RefreshAddInterval" client:"RefreshAddInterval"` // 递增触发间隔单位：次数
	RefreshAddCD       int      `col:"RefreshAddCD" client:"RefreshAddCD"`             // 递增触发间隔单位：秒
	RefreshAddNum      int      `col:"RefreshAddNum" client:"RefreshAddNum"`           // 刷新数量递增单位：个
	IsRandom           int      `col:"IsRandom" client:"IsRandom"`                     // 超出范围是否不随机
	RefreshPoint       IntSlice `col:"RefreshPoint" client:"RefreshPoint"`             // 刷新点
	Rewards            IntSlice `col:"rewards" client:"rewards"`                       // 掉落奖励
	Wave               int      `col:"wave" client:"wave"`                             // 波次显示
	FirstAll           int      `col:"firstAll" client:"firstAll"`                     // 首次全部刷新
}

type SkillListSkillListCfg struct {
	Id                     int        `col:"Id" client:"Id"`                                         // 技能ID
	BaseId                 int        `col:"baseId" client:"baseId"`                                 // 基类ID
	Level                  int        `col:"level" client:"level"`                                   // 等级
	IsSuper                int        `col:"isSuper" client:"isSuper"`                               // 是否是超武
	Levelpoints            int        `col:"levelpoints" client:"levelpoints"`                       // 三选一权重
	Color                  int        `col:"color" client:"color"`                                   // 花色
	SaveWithDead           int        `col:"saveWithDead" client:"saveWithDead"`                     // 施法者死亡后是否要保留（0：不保留，1：保留)
	DamageStatisticsID     int        `col:"damageStatisticsID" client:"damageStatisticsID"`         // 伤害统计技能id
	Type                   int        `col:"type" client:"type"`                                     // 技能类型
	SkillProbability       int        `col:"skillProbability" client:"skillProbability"`             // 技能优先级
	TargetGroup            int        `col:"targetGroup" client:"targetGroup"`                       // 目标阵营选择
	TargetSelectType       int        `col:"targetSelectType" client:"targetSelectType"`             // 目标选择类型
	TargetRange            FloatSlice `col:"targetRange" client:"targetRange"`                       // 目标选择类型和范围
	TargetSelectNum        int        `col:"targetSelectNum" client:"targetSelectNum"`               // 目标选择数量
	ActionToward           int        `col:"actionToward" client:"actionToward"`                     // 施法朝向
	ReleaseLimit           int        `col:"releaseLimit" client:"releaseLimit"`                     // 释放次数
	PublicCD               int        `col:"PublicCD" client:"PublicCD"`                             // 大冷却
	Cd                     int        `col:"cd" client:"cd"`                                         // 技能冷却时间
	Duration               int        `col:"duration" client:"duration"`                             // 持续时间
	InitBuff               IntSlice   `col:"initBuff" client:"initBuff"`                             // 获得技能附加buff
	ReleaseBuff            IntSlice   `col:"releaseBuff" client:"releaseBuff"`                       // 技能释放附加buff
	ProbabilityReleaseBuff IntSlice   `col:"probabilityReleaseBuff" client:"probabilityReleaseBuff"` // 释放时概率附加buff
	Probability            int        `col:"probability" client:"probability"`                       // 附加概率
	ReleaseSkill           IntSlice   `col:"releaseSkill" client:"releaseSkill"`                     // 技能结束时释放技能
	DamageGroup            int        `col:"damageGroup" client:"damageGroup"`                       // 判定目标阵营
	FixationTime           int        `col:"FixationTime" client:"FixationTime"`                     // 定帧时间
	FixColor               IntSlice   `col:"FixColor" client:"FixColor"`                             // 定帧变色
	RepelValue             float64    `col:"RepelValue" client:"RepelValue"`                         // 击退力度
	RepelTime              float64    `col:"RepelTime" client:"RepelTime"`                           // 击退时间
	FloatingValue          float64    `col:"FloatingValue" client:"FloatingValue"`                   // 击飞力度
	TrackType              int        `col:"trackType" client:"trackType"`                           // 轨迹类型
	TrackPro               FloatSlice `col:"trackPro" client:"trackPro"`                             // 轨迹参数
	TrackStartTime         int        `col:"trackStartTime" client:"trackStartTime"`                 // 轨迹触发时间
	TrackTriggerInterval   int        `col:"trackTriggerInterval" client:"trackTriggerInterval"`     // 轨迹触发间隔
	TrackMaxNum            int        `col:"trackMaxNum" client:"trackMaxNum"`                       // 轨迹触发次数
	SkillEffect            IntSlice   `col:"skillEffect" client:"skillEffect"`                       // 轨迹特效
	SkillEffectDuration    int        `col:"skillEffectDuration" client:"skillEffectDuration"`       // 轨迹特效时间
	DamageRatio1           int        `col:"damageRatio1" client:"damageRatio1"`                     // 伤害系数
	MaxHpRatio1            int        `col:"MaxHpRatio1" client:"MaxHpRatio1"`                       // 附加目标最大生命千分比
	TargetBuff1            IntSlice   `col:"targetBuff1" client:"targetBuff1"`                       // 伤害附加buff
	TrackEndSkill          IntSlice   `col:"trackEndSkill" client:"trackEndSkill"`                   // 二段技能效果
	TrackEndEffect         int        `col:"trackEndEffect" client:"trackEndEffect"`                 // 轨迹结束特效
	TEndEffectDuration     int        `col:"tEndEffectDuration" client:"tEndEffectDuration"`         // 轨迹结束特效时间
	DamageRatio2           int        `col:"damageRatio2" client:"damageRatio2"`                     // 伤害系数
	MaxHpRatio2            int        `col:"MaxHpRatio2" client:"MaxHpRatio2"`                       // 附加目标最大生命千分比
	TargetBuff2            IntSlice   `col:"targetBuff2" client:"targetBuff2"`                       // 伤害附加buff
	ConditionType          IntSlice   `col:"ConditionType" client:"ConditionType"`                   // 条件类型
	Condition              IntSlice2  `col:"Condition" client:"Condition"`                           // 条件参数
	RelatedSkillpro        IntSlice   `col:"relatedSkillpro" client:"relatedSkillpro"`               // 超武秘宝
	CarrySkillpro          IntSlice   `col:"carrySkillpro" client:"carrySkillpro"`                   // 非超武秘宝
}

type SkillComboSkillComboCfg struct {
	Id             int      `col:"Id" client:"Id"`                         // 招式id
	RelatedSkill   int      `col:"RelatedSkill" client:"RelatedSkill"`     // 关联技能
	ComboLv        int      `col:"ComboLv" client:"ComboLv"`               // 组合等级
	ComboType      int      `col:"ComboType" client:"ComboType"`           // 组合类型
	ScriptKey      IntSlice `col:"ScriptKey" client:"ScriptKey"`           // 组合参数
	BanCombo       IntSlice `col:"BanCombo" client:"BanCombo"`             // 招式黑名单
	ChangeSkill    IntSlice `col:"changeSkill" client:"changeSkill"`       // 替换技能
	RemoveSkillpro IntSlice `col:"RemoveSkillpro" client:"RemoveSkillpro"` // 移除配件
	Skillpro       IntSlice `col:"skillpro" client:"skillpro"`             // 赋予配件
}

type PlayerLevelCfg struct {
	Id          int     `col:"Id" client:"Id"`                   // 表格ID
	PlayerLevel int     `col:"playerLevel" client:"playerLevel"` // 玩家等级
	IApMax      float64 `col:"iApMax" client:"iApMax"`           // 当前等级体力上限
}

type LootPoolCfg struct {
	Id      int       `col:"Id" client:"Id"`           // 随机池ID
	Type    int       `col:"type" client:"type"`       // 随机池类型
	SubType int       `col:"SubType" client:"SubType"` // 子类型
	Iteam   IntSlice2 `col:"iteam" client:"iteam"`     // 产出道具
}

type PlayingMethodPlayingMethodCfg struct {
	Id       int    `col:"Id" client:"Id"`             // 索引
	Type     int    `col:"type" client:"type"`         // 关卡类型
	Cost     IntMap `col:"cost" client:"cost"`         // 消耗资源，数量
	Limit    int    `col:"limit" client:"limit"`       // 刷新上限
	TimeType int    `col:"timeType" client:"timeType"` // 开放时间类型
	Time     string `col:"time" client:"time"`         // 开放时间
}

type WeaponCfg struct {
	Id          int    `col:"Id" client:"Id"`                   // 索引ID
	WeaponId    int    `col:"WeaponId" client:"WeaponId"`       // 武器ID
	SkillDebris int    `col:"SkillDebris" client:"SkillDebris"` // 关联技能
	Item        IntMap `col:"Item" client:"Item"`               // 激活道具
}

type RoleLevelCfg struct {
	Id              int         `col:"Id" client:"Id"`                           // 唯一ID
	RoleLevel       int         `col:"roleLevel" client:"roleLevel"`             // 等级数
	IteamID         IntMap      `col:"iteamID" client:"iteamID"`                 // 升级所需道具
	BreakThrough    int         `col:"BreakThrough" client:"BreakThrough"`       // 突破节点
	LevelAttributes IntFloatMap `col:"LevelAttributes" client:"LevelAttributes"` // 升级属性
}

type RoleRankCfg struct {
	Id              int         `col:"Id" client:"Id"`                           // 唯一ID
	Type            int         `col:"type" client:"type"`                       // 类型
	RoleId          int         `col:"roleId" client:"roleId"`                   // 英雄ID
	RoleRank        int         `col:"roleRank" client:"roleRank"`               // 阶级数
	IteamID         IntMap      `col:"iteamID" client:"iteamID"`                 // 升星所需道具
	RankAttributes  IntFloatMap `col:"RankAttributes" client:"RankAttributes"`   // 升星属性
	StarCoefficient float64     `col:"starCoefficient" client:"starCoefficient"` // 星级系数
	RewardIteam     IntMap      `col:"rewardIteam" client:"rewardIteam"`         // 升星获得奖励
	PickSkillpro    IntSlice    `col:"pickSkillpro" client:"pickSkillpro"`       // 出战配件
	DefaultSkillpro IntSlice    `col:"defaultSkillpro" client:"defaultSkillpro"` // 助战配件
	BanSkillpro     IntSlice    `col:"banSkillpro" client:"banSkillpro"`         // 禁用配件
	PickPoker       int         `col:"pickPoker" client:"pickPoker"`             // 解锁扑克
	PickSkillcombo  IntSlice    `col:"pickSkillcombo" client:"pickSkillcombo"`   // 解锁组合
}

type EquipLevelCfg struct {
	Id       int       `col:"Id" client:"Id"`             // 唯一ID
	Level    int       `col:"level" client:"level"`       // 等级数
	Position int       `col:"position" client:"position"` // 装备位置
	IteamID  IntSlice2 `col:"iteamID" client:"iteamID"`   // 升级所需道具
}

type WeaponLevelCfg struct {
	Id               int         `col:"Id" client:"Id"`                             // 索引ID
	WeaponId         int         `col:"WeaponId" client:"WeaponId"`                 // 武器ID
	Level            int         `col:"Level" client:"Level"`                       // 当前等级
	Consume          IntMap      `col:"Consume" client:"Consume"`                   // 升级消耗
	UpgradeOutput    int         `col:"UpgradeOutput" client:"UpgradeOutput"`       // 升级产出库经验
	SkillCoefficient int         `col:"SkillCoefficient" client:"SkillCoefficient"` // 技能系数
	RankAttributes   IntFloatMap `col:"RankAttributes" client:"RankAttributes"`     // 升级属性
	PickSkillpro     IntSlice    `col:"pickSkillpro" client:"pickSkillpro"`         // 出战配件
	BanSkillpro      IntSlice    `col:"banSkillpro" client:"banSkillpro"`           // 禁用秘宝
	PickPoker        int         `col:"pickPoker" client:"pickPoker"`               // 解锁扑克
	PickSkillcombo   IntSlice    `col:"pickSkillcombo" client:"pickSkillcombo"`     // 解锁组合
}

type WeaponLibraryCfg struct {
	Id                int         `col:"Id" client:"Id"`                               // 当前等级
	Consume           int         `col:"Consume" client:"Consume"`                     // 升级所需经验
	InitialAttributes IntFloatMap `col:"InitialAttributes" client:"InitialAttributes"` // 全局属性
}

type SupportValueCfg struct {
	Id                int         `col:"Id" client:"Id"`                               // id
	Rarity            int         `col:"rarity" client:"rarity"`                       // 稀有度
	RoleRank          int         `col:"roleRank" client:"roleRank"`                   // 星级
	InitialAttributes IntFloatMap `col:"InitialAttributes" client:"InitialAttributes"` // 支援属性
}

type LotteryShipCfg struct {
	Id             int       `col:"Id" client:"Id"`                         // 卡池ID
	Type           int       `col:"Type" client:"Type"`                     // 类型
	Cost           int       `col:"Cost" client:"Cost"`                     // 消耗道具
	ShopId         int       `col:"ShopId" client:"ShopId"`                 // 道具商品ID
	Times          int       `col:"Times" client:"Times"`                   // 招募次数
	FreeTimes      int       `col:"FreeTimes" client:"FreeTimes"`           // 每X日免费1次单抽
	ProgressTimes  int       `col:"ProgressTimes" client:"ProgressTimes"`   // 每XX抽
	ProgressReward IntMap    `col:"ProgressReward" client:"ProgressReward"` // 获得一份奖励
	LootPool       IntSlice2 `col:"LootPool" client:"LootPool"`             // 随机池
	MinRewardTimes int       `col:"MinRewardTimes" client:"MinRewardTimes"` // 小保底次数
	TenReward      IntSlice  `col:"TenReward" client:"TenReward"`           // 小保底道具
	TenLootPool    int       `col:"TenLootPool" client:"TenLootPool"`       // 小保底池
	BigRewardTimes int       `col:"BigRewardTimes" client:"BigRewardTimes"` // 大奖保底次数
	BigReward      IntSlice  `col:"BigReward" client:"BigReward"`           // 大奖保底道具
	BigLootPool    int       `col:"BigLootPool" client:"BigLootPool"`       // 大奖保底池
}

type SevendayLoginCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // ID
	ActivityId int       `col:"ActivityId" client:"ActivityId"` // 活动ID
	Day        int       `col:"Day" client:"Day"`               // 第X天
	Reward     IntSlice2 `col:"Reward" client:"Reward"`         // 登录奖励
}

type EffectListEffectListCfg struct {
	Id                 int        `col:"Id" client:"Id"`                                 // 特效ID
	DamageRectType     FloatSlice `col:"damageRectType" client:"damageRectType"`         // 伤害包围盒类型
	DamageAddType      int        `col:"damageAddType" client:"damageAddType"`           // 伤害变化类型
	DamageAddRate      FloatSlice `col:"damageAddRate" client:"damageAddRate"`           // 伤害变化参数
	MaxDamageTime      int        `col:"maxDamageTime" client:"maxDamageTime"`           // 最大有效伤害次数
	DamagerTriggerTime IntSlice   `col:"damagerTriggerTime" client:"damagerTriggerTime"` // 判定触发（开始计算）时间
	DamageDisTime      int        `col:"damageDisTime" client:"damageDisTime"`           // 判定触发间隔
}

type MailCfg struct {
	Id     int       `col:"Id" client:"Id"`         // 邮件ID
	Title  int       `col:"Title" client:"Title"`   // 标题
	Text   int       `col:"Text" client:"Text"`     // 正文
	Reward IntSlice2 `col:"Reward" client:"Reward"` // 道具附件
}

type MonthlyCardCfg struct {
	Id              int       `col:"Id" client:"Id"`                           // ID
	Tpye            int       `col:"Tpye" client:"Tpye"`                       // 类型
	Time            int       `col:"Time" client:"Time"`                       // 持续时间
	DayReward       IntSlice2 `col:"DayReward" client:"DayReward"`             // 每日道具
	SkipAd          int       `col:"SkipAd" client:"SkipAd"`                   // 免广告
	MainlevelReward int       `col:"MainlevelReward" client:"MainlevelReward"` // 主线关卡产出增加
	OnHookReward    int       `col:"onHookReward" client:"onHookReward"`       // 挂机、快速挂机产出增加
	QuestTimes      int       `col:"QuestTimes" client:"QuestTimes"`           // 日常本*3每日扫荡次数增加
}

type StageFundCfg struct {
	Id        int       `col:"Id" client:"Id"`               // ID
	ChargeID  int       `col:"ChargeID" client:"ChargeID"`   // 充值ID
	Tpye      int       `col:"Tpye" client:"Tpye"`           // 类型
	StageId   int       `col:"StageId" client:"StageId"`     // 关卡ID
	Reward    IntSlice2 `col:"Reward" client:"Reward"`       // 免费道具
	PayReward IntSlice2 `col:"PayReward" client:"PayReward"` // 付费道具
}

type SignFundCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // ID
	ChargeID   int       `col:"ChargeID" client:"ChargeID"`     // 充值ID
	ActivityID int       `col:"ActivityID" client:"ActivityID"` // 活动ID
	Day        int       `col:"Day" client:"Day"`               // 第X天
	Reward     IntSlice2 `col:"Reward" client:"Reward"`         // 免费道具
	PayReward  IntSlice2 `col:"PayReward" client:"PayReward"`   // 付费道具
}

type ActivePassCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // ID
	ActivityID int       `col:"ActivityID" client:"ActivityID"` // 活动ID
	Lv         int       `col:"Lv" client:"Lv"`                 // 等级
	Exp        int       `col:"Exp" client:"Exp"`               // 升级所需经验
	Reward     IntSlice2 `col:"Reward" client:"Reward"`         // 免费道具
	PayReward  IntSlice2 `col:"PayReward" client:"PayReward"`   // 付费道具
}

type TaskPassCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // ID
	ActivityID int       `col:"ActivityID" client:"ActivityID"` // 活动ID
	Lv         int       `col:"Lv" client:"Lv"`                 // 等级
	Exp        int       `col:"Exp" client:"Exp"`               // 升级所需经验
	Reward     IntSlice2 `col:"Reward" client:"Reward"`         // 免费道具
	PayReward  IntSlice2 `col:"PayReward" client:"PayReward"`   // 付费道具
}

type SevendaysTaskCfg struct {
	Id     int       `col:"Id" client:"Id"`         // ID
	Point  int       `col:"Point" client:"Point"`   // 所需积分
	Reward IntSlice2 `col:"Reward" client:"Reward"` // 道具
}

type GetiApCfg struct {
	Id     int       `col:"Id" client:"Id"`         // ID
	Time   string    `col:"Time" client:"Time"`     // 领取时间段
	Reward IntSlice2 `col:"Reward" client:"Reward"` // 道具
}

type TalentCfg struct {
	Id           int         `col:"Id" client:"Id"`                     // ID
	Type         int         `col:"Type" client:"Type"`                 // 类型
	Attributes   IntFloatMap `col:"Attributes" client:"Attributes"`     // 属性
	PickSkillpro IntSlice    `col:"pickSkillpro" client:"pickSkillpro"` // 出战配件
	Reward       IntSlice2   `col:"Reward" client:"Reward"`             // 消耗道具
}

type AppearanceCfg struct {
	Id         int         `col:"Id" client:"Id"`                 // ID
	Type       int         `col:"Type" client:"Type"`             // 类型
	Condition  int         `col:"Condition" client:"Condition"`   // 激活条件
	Parm       IntSlice    `col:"Parm" client:"Parm"`             // 条件参数
	Attributes IntFloatMap `col:"Attributes" client:"Attributes"` // 外观属性
}

type SkillproListSkillproListCfg struct {
	Id             int      `col:"Id" client:"Id"`                         // 配件ID
	BaseId         int      `col:"BaseId" client:"BaseId"`                 // 配件组id
	Type           int      `col:"Type" client:"Type"`                     // 类型
	Quality        int      `col:"Quality" client:"Quality"`               // 品质
	Lv             int      `col:"Lv" client:"Lv"`                         // 等级
	Weight         int      `col:"Weight" client:"Weight"`                 // 是否初始生效
	RelatedSkill   int      `col:"RelatedSkill" client:"RelatedSkill"`     // 关联技能id
	DeBlock        IntSlice `col:"deBlock" client:"deBlock"`               // 局内解锁配件
	DeTach         IntSlice `col:"deTach" client:"deTach"`                 // 黑名单
	SkillId        IntSlice `col:"SkillId" client:"SkillId"`               // 附加技能
	BuffId         IntSlice `col:"BuffId" client:"BuffId"`                 // 附加buff
	EventId        IntSlice `col:"EventId" client:"EventId"`               // 触发事件
	ConstructionId IntSlice `col:"ConstructionId" client:"ConstructionId"` // 刷新交互物
	Value          int      `col:"Value" client:"Value"`                   // 购买价格
}

type BattlePassBattlePassCfg struct {
	Id          int      `col:"Id" client:"Id"`                   // 索引
	Season      int      `col:"Season" client:"Season"`           // 赛季
	BattleTime  IntSlice `col:"BattleTime" client:"BattleTime"`   // 赛季时间
	WeekTime    IntSlice `col:"WeekTime" client:"WeekTime"`       // 周时间（开始，结束）
	StageBuff   int      `col:"StageBuff" client:"StageBuff"`     // 战斗环境
	BOSSRefresh int      `col:"BOSSRefresh" client:"BOSSRefresh"` // BOSS刷新ID
	Broadcast   IntSlice `col:"Broadcast" client:"Broadcast"`     // 适用广播
}

type BattleMatchBattleMatchCfg struct {
	Id              int       `col:"Id" client:"Id"`                           // Id
	CupStart        int       `col:"CupStart" client:"CupStart"`               // 奖杯区间-开始（闭区间）
	CupEnd          int       `col:"CupEnd" client:"CupEnd"`                   // 奖杯区间-结束（闭区间）
	BattleLevel     int       `col:"BattleLevel" client:"BattleLevel"`         // 段位
	Star            int       `col:"Star" client:"Star"`                       // 星数
	MatchRangeStart int       `col:"MatchRangeStart" client:"MatchRangeStart"` // 匹配区间-开始（闭区间）
	MatchRangeEnd   int       `col:"MatchRangeEnd" client:"MatchRangeEnd"`     // 匹配区间-结束（闭区间）
	BattleMap       IntSlice  `col:"BattleMap" client:"BattleMap"`             // 比赛关卡
	BOSS            int       `col:"BOSS" client:"BOSS"`                       // 触发BOSS能量
	ClearanceEnergy int       `col:"clearanceEnergy" client:"clearanceEnergy"` // 结算能量
	WinCup          IntSlice  `col:"WinCup" client:"WinCup"`                   // 奖杯数量
	WinExp          IntSlice  `col:"WinExp" client:"WinExp"`                   // 通行证经验
	WinReward       IntSlice2 `col:"WinReward" client:"WinReward"`             // 胜利奖励
	RobotId         IntSlice  `col:"RobotId" client:"RobotId"`                 // 机器人ID
}

type BattleItemBattleItemCfg struct {
	Id         int      `col:"Id" client:"Id"`                 // 道具Id
	Type       int      `col:"Type" client:"Type"`             // 道具类型
	EffectType int      `col:"effectType" client:"effectType"` // 效果类型
	EffectArgs IntSlice `col:"effectArgs" client:"effectArgs"` // 效果参数1
	Cd         int      `col:"cd" client:"cd"`                 // 使用冷却
}

type EmoteEmoteCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // 表情Id
	Type       int       `col:"Type" client:"Type"`             // 表情类型
	UnlockType IntSlice  `col:"UnlockType" client:"UnlockType"` // 解锁条件
	Conditions IntSlice2 `col:"Conditions" client:"Conditions"` // 解锁参数
	UnlockDes  int       `col:"UnlockDes" client:"UnlockDes"`   // 解锁文本
}

type AdventureAdventureCfg struct {
	Id      int    `col:"Id" client:"Id"`           // 奇遇Id
	Rewards IntMap `col:"Rewards" client:"Rewards"` // 奖励
}

type RankRankCfg struct {
	Id int `col:"Id" client:"Id"` // 排行榜ID
}

type RankRewardCfg struct {
	Id      int    `col:"Id" client:"Id"`           // Id
	RankID  int    `col:"RankID" client:"RankID"`   // 排行榜ID
	Stage   int    `col:"Stage" client:"Stage"`     // 关卡ID
	Rewards IntMap `col:"Rewards" client:"Rewards"` // 首通奖励
}

type RankingRewardCfg struct {
	Id        int       `col:"Id" client:"Id"`               // Id
	RankID    int       `col:"RankID" client:"RankID"`       // 排行榜ID
	SubType   int       `col:"SubType" client:"SubType"`     // 类型
	RankStart int       `col:"RankStart" client:"RankStart"` // 开始名次
	RankEnd   int       `col:"RankEnd" client:"RankEnd"`     // 结束名次
	Rewards   IntSlice2 `col:"Rewards" client:"Rewards"`     // 排名奖励
	MailID    int       `col:"mailID" client:"mailID"`       // 发放邮件ID
}

type GuildGuildCfg struct {
	Id      int `col:"Id" client:"Id"`           // Id
	Lv      int `col:"Lv" client:"Lv"`           // 联盟等级
	Exp     int `col:"Exp" client:"Exp"`         // 升级经验
	Num     int `col:"Num" client:"Num"`         // 联盟人数
	ItemMax int `col:"ItemMax" client:"ItemMax"` // 每周勋章上限
}

type AdvMapAdvMapCfg struct {
	Id        int `col:"Id" client:"Id"`               // Id
	Map       int `col:"Map" client:"Map"`             // 大地图
	Region    int `col:"Region" client:"Region"`       // 区域
	RegionMap int `col:"RegionMap" client:"RegionMap"` // 区域地图文件
}

type AdvPointAdvPointCfg struct {
	Id            int         `col:"Id" client:"Id"`                       // Id
	Map           int         `col:"Map" client:"Map"`                     // 大地图
	Region        int         `col:"Region" client:"Region"`               // 区域
	Type          int         `col:"Type" client:"Type"`                   // 类型
	Value         int         `col:"Value" client:"Value"`                 // 战斗战力
	FixItem       IntFloatMap `col:"FixItem" client:"FixItem"`             // 修复所需道具
	OnhookReward  IntFloatMap `col:"OnhookReward" client:"OnhookReward"`   // 单次产出道具
	OnhookTime    int         `col:"OnhookTime" client:"OnhookTime"`       // 单次产出时间(单位：秒)
	OnhookType    int         `col:"OnhookType" client:"OnhookType"`       // 自动/手动收获
	OnhookTimeMax int         `col:"OnhookTimeMax" client:"OnhookTimeMax"` // 收获时间上限(单位：秒)
	Mystery       int         `col:"Mystery" client:"Mystery"`             // 第N个秘境
	Reward        IntSlice2   `col:"Reward" client:"Reward"`               // 战斗胜利奖励
}

type AdvCardAdvCardCfg struct {
	Id      int       `col:"Id" client:"Id"`           // Id
	BaseId  int       `col:"BaseId" client:"BaseId"`   // 卡牌组ID
	Type    int       `col:"Type" client:"Type"`       // 类型
	Move    int       `col:"Move" client:"Move"`       // 前进格数
	Event   int       `col:"Event" client:"Event"`     // 事件效果类型
	Effect1 IntSlice2 `col:"effect1" client:"effect1"` // 效果参数
	Value   int       `col:"Value" client:"Value"`     // 战力要求
	Weight  int       `col:"Weight" client:"Weight"`   // 权重
	Reward  IntSlice2 `col:"Reward" client:"Reward"`   // 战斗胜利奖励
}

type AdvCardBagAdvCardBagCfg struct {
	Id            int      `col:"Id" client:"Id"`                       // 唯一id
	BaseId        IntSlice `col:"BaseId" client:"BaseId"`               // 卡牌组
	Num           IntSlice `col:"Num" client:"Num"`                     // 卡牌个数
	QualityWeight IntSlice `col:"QualityWeight" client:"QualityWeight"` // 品质随机
	CradID        IntSlice `col:"CradID" client:"CradID"`               // 指定卡牌ID
}

type AdvValueAdvValueCfg struct {
	Id        int `col:"Id" client:"Id"`               // Id
	MapId     int `col:"MapId" client:"MapId"`         // 地图ID
	CellStart int `col:"CellStart" client:"CellStart"` // 从第N格（闭区间）
	CellEnd   int `col:"CellEnd" client:"CellEnd"`     // 到第N格（闭区间）
	Value     int `col:"Value" client:"Value"`         // 战力要求
}

type ContractCfg struct {
	Id         int       `col:"Id" client:"Id"`                 // 任务ID
	Group      int       `col:"group" client:"group"`           // 分组
	TaskType   IntSlice  `col:"taskType" client:"taskType"`     // 任务条件
	TaskReward IntSlice2 `col:"taskReward" client:"taskReward"` // 任务奖励
}

type BroadcastCfg struct {
	Id        int      `col:"id" client:"id"`               // 索引
	Type      int      `col:"type" client:"type"`           // 广播类型
	ScriptKey IntSlice `col:"ScriptKey" client:"ScriptKey"` // 条件参数
}

type RedPacketRedPacketCfg struct {
	Id          int       `col:"Id" client:"Id"`                   // Id
	ExtraReward IntSlice2 `col:"ExtraReward" client:"ExtraReward"` // 奖励
}

type OutMapOutMapCfg struct {
	Id   int `col:"Id" client:"Id"`     // 地图索引
	Type int `col:"Type" client:"Type"` // 房间类型
}

type ZombieColiseumZombieColiseumCfg struct {
	Id          int       `col:"Id" client:"Id"`                   // Id
	RankStart   int       `col:"RankStart" client:"RankStart"`     // 名次区间-开始
	RankEnd     int       `col:"RankEnd" client:"RankEnd"`         // 名次区间-结束
	BattleMap   IntSlice  `col:"BattleMap" client:"BattleMap"`     // 关卡参数
	WinReward   IntSlice2 `col:"WinReward" client:"WinReward"`     // 胜利奖励
	LoseReward  IntSlice2 `col:"LoseReward" client:"LoseReward"`   // 失败奖励
	RobotId     IntSlice  `col:"RobotId" client:"RobotId"`         // 机器人ID
	RobotCombat IntSlice  `col:"RobotCombat" client:"RobotCombat"` // 机器人战力区间
}

type HandbookHandbookCfg struct {
	Id             int       `col:"Id" client:"Id"`                         // Id
	Type           int       `col:"Type" client:"Type"`                     // 类型
	Condition      int       `col:"Condition" client:"Condition"`           // 激活条件
	Param          int       `col:"Param" client:"Param"`                   // 参数
	Reward         IntSlice2 `col:"Reward" client:"Reward"`                 // 激活奖励
	ArenaMonster   int       `col:"ArenaMonster" client:"ArenaMonster"`     // 是否竞技场怪物
	ArenaCondition int       `col:"ArenaCondition" client:"ArenaCondition"` // 竞技场激活条件
	ArenaParam     int       `col:"ArenaParam" client:"ArenaParam"`         // 竞技场激活参数
	ArenacostItem  IntSlice2 `col:"ArenacostItem" client:"ArenacostItem"`   // 竞技场激活道具
}

type HelpTipsCfg struct {
	Id     int `col:"Id" client:"Id"`         // ID
	Type   int `col:"type" client:"type"`     // 类型
	Weight int `col:"weight" client:"weight"` // 出现权重
}

type LanguageCfg struct {
	Id int    `col:"Id" client:"Id"` // 索引id
	Zh string `col:"zh" client:"zh"` // 中文文本
}
