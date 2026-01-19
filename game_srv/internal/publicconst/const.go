package publicconst

import "time"

type PlayerState int
type ServiceId int
type ItemSource int

type TaskState int
type TaskType uint32
type TaskCond uint32

type MailState uint32
type UseItemType uint32

type ItemCode uint32

type OpenConditionType uint32
type FunctionType uint32

type BattleResult uint32
type BattlStar uint32

type ServerStatus uint32
type ShopRefreshType uint32

type ServEnv string

type ShipSource uint32
type ShipAttrType uint32

type RoleSource uint32
type RoleAttrType uint32
type ItemBigType uint32

type UnlockType uint32
type PlayMethodId uint32

type FunctionId uint32

type PayType uint32

type ActivityType uint32

var (
	MAX_USERID_LEN = 10

	GLOBAL_DB_NAME            = "global"
	GLOBAL_ACCOUNT_COLLECTION = "account"
	//GLOBAL_SERVERINFO_COLLECTION = "serverInfo"
	GLOBAL_SERVERINFO_COLLECTION = "server_info"
	GLOBAL_MAIL                  = "mail"

	LOACL_ACCOUNT              = "account"
	LOACL_ITEM                 = "item"
	LOCAL_TASK                 = "task"
	LOCAL_MAIL                 = "mail"
	LOCAL_ACHIEVE              = "achieve"
	LOCAL_MISSION              = "mission"
	LOCAL_SHOP                 = "shop"
	LOCAL_SHIP                 = "ship"
	LOCAL_ROLE                 = "role"
	LOCAL_TEAM                 = "team"
	LOCAL_EQUIP                = "equip"
	LOCAL_STARPORT             = "starport"
	LOCAL_PLAYMETHOD           = "playmethod"
	LOCAL_WEAPON               = "weapon"
	LOACL_TREASURE             = "treasure"
	LOCAL_ACTIVITY             = "activity"
	LOCAL_POKER                = "poker"
	LOCAL_CARDPOOL             = "cardpool"
	ORDER                      = "order"
	LOCAL_APPEARANCE           = "appearance"
	LOCAL_PET                  = "pet"
	LOCAL_MISSION_REWARD       = "missionreward"
	LOCAL_FRIEND               = "friend"
	LOCAL_GLOBAL               = "globalinfo"
	LOCAL_ALLIANCE             = "alliance"
	LOCAL_ALLIANCE_MEMBER      = "alliance_member"
	LOCAL_ALLIANCE_APPLICATION = "alliance_application"
	LOCAL_Alliance_RedPacket   = "alliance_redpacket"
	LOCAL_ALLIANCE_SHOP        = "alliance_shop"
	LOCAL_ALLIANCE_BOSS        = "alliance_boss"
	COLL_EXPLORE               = "explore"
	LOCAL_FIGHT                = "fight"
	LOCAL_PEAK_FIGHT           = "peak_fight"
	LOCAL_CONTRACT             = "contract"
	LOCAL_CUMULATIVE           = "cumulative"
	LOCAL_DESERT_FIGHT         = "desert_fight"
	LOCAL_COUNT                = "count"
	LOCAL_ARENA                = "arena"
	LOCAL_ATLAS                = "atlas"
	LOCAL_LUCK_SALE            = "luck_sale"
	LOCAL_FUNCTION_PREVIEW     = "function_preview"
	LOCAL_PET_RANK             = "rank_pet"
	LOCAL_ARENA_RANK           = "rank_arena"
	LOCAL_RANK                 = "rank_mission"
	LOCAL_RankAlliance         = "rank_alliance"
	LOCAL_SPECIAL_MISSION_RANK = "rank_specailmission"
	LOCAL_DESERT_RANK          = "rank_desert"

	ITEM_LOG       = "itemLog"
	DB_OP_TIME_OUT = 20 * time.Second

	CLIENT_HEART_INTERVAL = 15 // 客户端心跳间隔(s)
	MAX_CLIENT_HERART_NUM = 3  // 最大心跳包数量

	MAX_UPDATE_ORE_TOTAL_TIME = 10 // 更新矿洞总量时间

	MONGO_NO_RESULT = "mongo: no documents in result"

	MAX_RECYCLE_PLAYER_DATA = 3600 // 玩家数据保留1小时
	DAY_SECONDS             = 24 * 3600

	REFRESH_ORE_INTEVAL = 10 // 刷新矿洞总量间隔

	GS_SERVER_PREFIX = "/server/gs/"
	MAX_SERVER_TTL   = 5

	UID_SERVER_PREFIX = "/server/uid"

	MAX_MAIL_OP_COUNT = 10

	EVENT_POOL_SIZE uint32 = 100
	MAX_EVENT_SIZE  uint32 = 100

	MAX_MSG_COUN_MINUTE uint32 = 500 // 一分钟法术举报最大的数量

	MISSION_BATTLE_MIN_SECONDS = 0 // 关卡战斗最小时长

	MAX_NICK_LENGTH = 10

	MAX_HEAD_IMAGE_LENGTH = 10
	MAX_SIGN_LENGTH       = 50
	UPDATE_NICK_CD        = 10

	ACTIVE_SECONDS = 1800 // 活跃时间
)

const (
	Logining PlayerState = iota // 登录中
	Online                      // 在线
	Offline                     // 离线
)

const (
	GMService ServiceId = iota
	ItemService
	AccountService
	CommonService
	GlobalService
	NatsService
	InterService
	SystemService
	MissionService
	OnHookService
	ShipService
	RoleService
	TeamService
	TaskService
	EquipService
	StarPortService
	MailService
	ShopService
	PlayMethodService
	WeaponService
	TreasureService
	PayService
	ActivityService
	PokerService
	CardPoolService
	ListenRedisService
	ChargeService
	SocialService
	AdService
	RankService
	AttrService
	RsService
	AppearanceService
	PvPService
	BattleService
	PetService
	WorldChatService
	FriendService
	AllianceService
	ExploreService
	AllianceRankService
	FightService
	PeakFightService
	ContractService
	DesertService
	ArenaService
	AtlasService
	LuckSaleService
	FunctionPreviewService
)

const (
	OreAddItem                  ItemSource = iota // 挖矿获得
	InitAddItem                                   // 初始添加道具
	GMAddItem                                     // gm 获得
	UseCostItem                                   // 道具使用扣除
	UseAddItem                                    // 道具使用添加
	LevelAddItem                                  // 等级提升获得
	RecoveryItem                                  // 自然恢复
	BuyApCostItem                                 // 购买体力消耗
	BuyApAddItem                                  // 购买体力获得
	MissionCostItem                               // 关卡消耗体力
	PassMissionAddItem                            // 通关获得道具
	MissionBoxAddItem                             // 关卡宝箱获得
	OnHookAddItem                                 // 挂机奖励获得道具
	GMCostItem                                    // GM扣除道具
	QuickOnHookCostItem                           // 快速挂机消耗
	QuickOnHookAddItem                            // 快速挂机获得
	ShipUpgradeCostItem                           // 机甲升级消耗
	ShipUpgradeStarCostItem                       // 机甲升星
	ShipUpgradeStarAddItem                        // 机甲升星获得
	RoleUpgradeStarCostItem                       // 驾驶员升星消耗
	RoleUpgradeStartAddItem                       // 驾驶员升星获得
	RoleSendGiftCostItem                          // 赠送礼物消耗
	RoleSendGiftAddItem                           // 赠送礼物获得
	RoleBatchSendGiftCostItem                     // 一键赠送礼物消耗
	RoleBatchSendGiftAddItem                      // 一键赠送礼物获得
	RoleBatchSendGiftReturnItem                   // 一键赠送返还
	GetTaskRewardAddItem                          // 任务奖励获得
	GetTaskActiveRewardAddItem                    // 任务活跃度奖励
	EquipUpgradeCostItem                          // 装备升级消耗
	EquipAutoUpgradeCostItem                      // 装备自动升级
	RdDrawCardCostItem                            // 研发建筑抽卡消耗
	RdUpgradeLevel                                // 研发建筑升级
	MailAddItem                                   // 邮件获得
	ShopAddItem                                   // 商店获得
	ShopCostItem                                  // 商店消耗
	CdkAddItem                                    // cdk 获得
	PlayMethodAddItem                             // 玩法获得
	UpgradeCostItem                               // 升级消耗
	ShipExchangeCostItem                          // 机甲兑换消耗
	SuitRewardAddItem                             // 套装奖励获得
	UseItem                                       // 使用道具
	ExchangeShip                                  // 兑换机甲
	UpgradeWeaponCostItem                         // 升级武器消耗
	UpgradeWeaponAddItem                          // 升级武器获得
	ChargeBuyAddItem                              // 充值购买获得
	PassMissionAddWeapon                          // 通关获得武器
	CostTenLottery                                // 十次抽卡消耗
	AddTenLottery                                 // 十次抽卡获得
	CostOneLottery                                // 一次抽卡消耗
	AddOneLottery                                 // 一次抽卡获得
	AddFirstLottery                               // 首次抽卡获得
	LoginActivityAddItem                          // 登录活动获得
	QuestionRewardAddItem                         // 问卷调查获得
	ChargeAddItem                                 // 充值获得
	AdAddItem                                     // 广告获得
	MonthCardAddItem                              // 月卡获得
	MonthCardRewardAddItem                        // 月卡奖励获得
	MainFundAddItem                               // 主线基金获得
	BuyMainFundAddItem                            // 主线基金充值获得
	ActivePassAddItem                             // 战令获得
	BuyActivePassGradeCostItem                    // 购买通行证消耗
	ActivePassBuyAddItem                          // 购买战令获得
	TaskPassAddItem                               // 战令任务获得
	BuyTaskPassGradeCostItem                      // 购买战令消耗
	MonthCardDailyAddItem                         // 月卡每日获得
	OpenServerAddItem                             // 开服活动获得
	OpenServerActiveAddItem                       // 开服活动获得
	LikesMissionRankAddItem                       // 点赞排行榜获得
	DailApAddItem                                 // 日常体力获得
	TalentUpgradeCostItem                         // 天赋升级消耗
	AdMissionAddItem                              // 关卡结算获得
	ActiveAppearance                              // 激活外观消耗
	ConstructionLotteryCostItem                   // 交互物抽奖消耗
	PetGrowUpCostItem                             // 宠物成长加速扣道具
	PetAdventureAddItem                           // 冒险奖励获得
	BuyPetEggSlotCostItem                         // 购买宠物蛋槽消耗
	AcceleratePetEggCostItem                      // 宠物蛋加速消耗
	PetEggCultivateCostItem                       // 宠物蛋培养消耗
	PetEggAddItem                                 // 宠物蛋获得
	PetEggCostItem                                // 宠物蛋完成消耗
	PetMaterialAddItem                            // 宠物材料获得
	PetPartActiveCostItem                         // 宠物部位激活消耗
	AdventureAddItem                              // 奇遇引导奖励
	MissionExtraAddItem                           // 关卡额外奖励获得
	RankMissionRewardAddItem                      // 关卡排行榜奖励获得
	PetEvolutionCostItem                          // 宠物蛋进化消耗
	PetSuccinctCostItem                           // 宠物洗炼消耗
	PetLevelRewardAddItem                         // 宠物等级奖励获得
	DesertSignAddItem                             // 沙漠签到获得
	ActiveWeapon                                  // 激活武器获得
	ActiveWeaponCostItem                          // 激活武器消耗
	DesertSuppSignAddItem                         // 沙漠补签消耗
	DesrtSuppSignCostItem                         // 沙漠补签消耗
	LotteryTimesAddItem                           // 抽卡次数奖励获得
	DesertActPreviewAddItem                       // 沙漠活动预览获得
	CreateAllianceCostItem                        // 创建联盟消耗道具
	AllianceRedPacket                             // 联盟红包
	ExploreAddItem                                // 探索消耗
	ExploreUseItem                                // 探索消耗道具
	ExploreCollectItem                            // 探索收集建筑物道具
	ExploreReward                                 // 冒险-每天免费获得卡包的道具ID
	PeakFight                                     // 巅峰战场
	PeakFightMatch                                // 巅峰战场匹配
	PeakFightReset                                // 巅峰战场重置
	ArenaBuyPkCnt                                 //竞技场购买挑战次数
	ArenaBuyPk                                    //竞技场挑战
	AtlasActivate                                 // 图鉴激活
	LuckSaleExtract                               // 幸运售货机抽奖
	LuckSaleTaskReward                            // 幸运售货机任务完成奖励
	ContractRand                                  // 合约刷新
	ContractReward                                // 合约奖励
	FunctionPreviewReward                         // 功能预览奖励
	GemRefresh                                    // 宝石洗练消耗
	EquipStage                                    // 装备本关卡奖励
	EquipStagePick                                // 装备本掉落
	EquipStageBuyCount                            // 装备本购买领奖次数
	ResourcesPassCost                             // 资源本道具消耗
	ResourcesPassBuyItem                          // 资源本购买道具
	ResourcesPassAttackCost                       // 资源本攻击消耗
	ResourcesPassAttackReward                     // 资源本攻击获取
	PetActivate                                   // 宠物激活
	PetLvUp                                       // 宠物升级消耗
	MaxFirstPassReward                            // 领取最大首通奖励
	ComposeItem                                   // 道具合成
	CoatActive                                    // 皮肤激活消耗
)

const (
	TASK_ACCEPT   TaskState = iota
	TASK_COMPLETE           // 完成未领奖
	TASK_DONE               // 完成领奖
)

const (
	MAIN_TASK    TaskType = iota + 1 // 主线任务
	ACHIEVE_TASK                     // 成就任务
	DAILY_TASK                       // 日常任务
	WEEKLY_TASK                      // 周长任务

	ALLIANCE_WEEKLY_TASK = 10 // 联盟每周任务
	LUCK_SALE_TASK       = 11 // 幸运售货机任务
)

const (
	TASK_COND_PASS_MISSION_ID              TaskCond = 1001 // 通关指定关卡
	TASK_COND_BATTLE_MISSIONTYPE                    = 1002 // 挑战、扫荡X类型关卡X次
	TASK_COND_ANY_BATTLE                            = 1003 // 参与任意战斗
	TASK_COND_ANY_ONE_MISSION                       = 1004 //
	TASK_COND_MISSION_BOX_REWARD                    = 1005 // 领取X次通关宝箱
	TASK_COND_FIRST_PASS_CHALLENGE_MISSION          = 1006 // 首通任意经营难度N次
	TASK_COND_PASS_TYPE_MISSION                     = 1007 // 完成XX关卡类型XX次

	TASK_COND_ADD_ITEM  = 2001 // 获得道具
	TASK_COND_COST_ITEM = 2002 // 消耗道具

	TASK_COND_LOGIN = 3001

	TASK_COND_SHIP_LEVEL        = 4001
	TASK_COND_SHIP_RARITY_NUM   = 4002
	TASK_COND_SHIP_STAR_NUM     = 4003
	TASK_COND_SPECIFY_SHIP_STAR = 4004
	TASK_COND_SHIP_SWITCH       = 4005 // 更换X次上阵库鲁
	TASK_COND_ACTIVE_SHIP       = 4006 // 激活XX个XX品质库鲁

	TASK_COND_GET_EQUIP_NUM       = 5001
	TASK_COND_UPGRADE_EQUIP       = 5002
	TASK_COND_ANY_EQUIP_LEVEL     = 5004
	TASK_COND_ALL_EQUIP_POS_LEVEL = 5006
	TASK_COND_EQUIP_RARITY_NUM    = 5007
	TASK_COND_COMPOSE_EQUIP       = 5008
	TASK_COND_PUT_ON_EQUIP        = 5009

	TASK_COND_PUT_ON_DISK     = 5101 // 镶嵌XX个磁盘
	TASK_COND_ADD_DISK        = 5102 // 获得XX个磁盘
	TASK_COND_ADD_RARITY_DISK = 5103 // 获得X个XX品质磁盘

	TASK_COND_UPGRADE_WEAPON       = 6001
	TASK_COND_ANY_WEAPON_LEVEL     = 6002
	TASK_COND_SPECIFY_WEAPON_LEVEL = 6003
	TASK_COND_ALL_WEAPON_MIN_LEVEL = 6004
	TASK_COND_WEAPON_LIB_LEVEL     = 6005
	TASK_COND_GET_WEAPON_NUM       = 6006

	TASK_COND_KILL_MONSTER_NUM            = 7001
	TASK_COND_KILL_ELITE_MONSTER_BOSS_NUM = 7003

	TASK_COND_COMPLETE_DAILY_TASK_COUNT  = 8001
	TASK_COND_COMPLETE_WEEKLY_TASK_COUNT = 8002
	TASK_COND_GET_DAILY_BOX              = 8003 // 领取活跃度宝箱

	TASK_COND_DAILY_ACTIVE_SCORE  = 9001
	TASK_COND_WEEKLY_ACTIVE_SCORE = 9002

	TASK_COND_GET_ON_HOOK_REWARD = 10001
	TASK_COND_QUICK_ON_HOOK      = 10002

	TASK_COND_BUY_AP = 11001

	TASK_COND_CLICK_UPDATE_NICK = 13001

	TASK_COND_GET_CHIP_NUM = 14001 // 获得筹码数XXXX
	TASK_COND_POKER_NUM    = 14002 // 凑成XX牌型XX次

	TASK_COND_LOTTERY                = 15001 // 累计在任意卡池招募XX次
	TASK_COND_SHIP_LOTTERY           = 15002 // 累计进行XX次库鲁招募
	TASK_COND_DISK_BIND_BOX_LOTTERY  = 15003 // 累计进行XX次磁盘盲盒抽奖
	TASK_COND_SHIP_BEAST_EGG_LOTTERY = 15004 // 累计进行XX次库鲁兽扭蛋

	TASK_COND_GET_PET = 16001 // 获得宠物 激活XX个XX品质库鲁兽
	//TASK_COND_PET_ADVENTURE = 16002 // 宠物冒险
	TASK_COND_PET_LEVEL = 16002 // 库鲁兽等级达到XX级

	TASK_COND_PET_PART_CHARM = 17001 //潮流度（魅力值）达到XXX

	TASK_COND_ALLIANCE_JOIN        = 18001 //加入联盟
	TASK_COND_ALLIANCE_BOSS_FIGHT  = 18002 //挑战联盟BOSS，N次
	TASK_COND_ALLIANCE_FINISH_TASK = 18003 //完成联盟任务N个

	TASK_COND_PEAK_FIGHT_PK   = 19001 //进行PVP战斗N次
	TASK_COND_PEAK_FIGHT_RANK = 19002 //赛季最高达到X段位

	TASK_COND_EXPLORE_USE_CARD     = 20001 //探索使用卡牌X次
	TASK_COND_EXPLORE_OCCUPY_BUILD = 20002 //占领指定ID的建筑
	TASK_COND_EXPLORE_PROGRESS     = 20003 //达成XX地图的XX%探索度 TODO

	TASK_COND_RESUME_DIAMOND_BUY_SHOP = 21001 // 消耗钻石购买商品x次
	TASK_COND_REWARD_EQUIP            = 22001 // 装备本累计领取奖励1次
	TASK_COND_RESOURCES_PASS          = 23001 // 进行XX次素材本
)

const (
	MAIL_NONE MailState = iota
	MAIL_OPEN
	MAIL_GET_REWARD
)

const (
	USE_ITEM_NONE            UseItemType = iota
	USE_ITEM_COMPOSE                     // 合成
	USE_ITEM_DECOMPOSE                   // 分解
	USE_ITEM_COST            = 3         // 道具消耗
	USE_ITEM_SELECT          = 6         // 自选道具
	USE_ITEM_FROM_ITEM_GROUP = 7         // 随机获取一个
	USE_ITEM_GET_SHIP        = 8         // 获得机甲

	USE_ITEM_ADD_PLAYMETHOD_TIMES = 12 // 添加玩法次数
	USE_ITEM_ADD_Equip            = 13 // 增加装备
	USE_ITEM_ADD_PET              = 19 // 加宠物

	USE_ITEM_ADV_CARD_PACKAGE = 21 // 世界探索-使用行动牌卡包
	USE_ITEM_ADV_CARD         = 22 // 世界探索-使用行动牌

	USE_ITEM_ADD_GEM = 25 // 增加宝石
)

const (
	ITEM_CODE_NONE    ItemCode = iota
	ITEM_CODE_GOLD    ItemCode = 600001 // 金币
	ITEM_CODE_DIAMOND ItemCode = 600002 // 钻石
	ITEM_CODE_AP      ItemCode = 600003 // 体力
	ITEM_CODE_EXP     ItemCode = 600004 // 账号经验
	ITEM_ALLIANCE_EXP ItemCode = 800004 // 联盟经验
)

const (
	LOSE BattleResult = iota
	WIN
)

const (
	BATTLE_ONE_STAR BattlStar = iota + 1
	bATTLE_TWO_STAR
	BATTLE_THREE_STAR
	BATTLE_FOUR_STAR
)

const (
	SERVER_CLOSE        ServerStatus = iota
	SERVER_INSIDE_OPEN               //  服务器内部开放
	SERVER_OUTSIDE_OPEN              // 服务器对外开放
)

const (
	Env_Test ServEnv = "test"
	Env_Dev  ServEnv = "dev"
	Env_Prod ServEnv = "prod"
)

const (
	Ship_Level      ShipAttrType = iota // 机甲等级
	Ship_Star_Level                     // 机甲星级
)

const (
	InitRole RoleSource = iota
)

const (
	Role_Star_Level  RoleAttrType = iota // 驾驶员星级
	Role_Favor_Level                     // 驾驶员好感度等级f
)

const (
	Item_BigType_NONE              ItemBigType = iota
	Item_BigType_Favor                         = 7  // 好感度
	Item_BigType_Peak_Fight_energy             = 31 // pvp能量
)

const (
	Unlock_Mission           UnlockType = iota + 1 // 关卡解锁
	Unlock_Player_Level                            // 玩家等级解锁
	Unlock_Time                                    // 绝对时间解锁
	Unlock_Create_Role_Time                        // 创角时间解锁
	Unlock_Challenge_Mission                       // 挑战关卡解锁
	Unlock_Shop_Item
	Unlock_Open_Server_Time // 开服时间解锁
	Unlock_Pet_Id
)

const (
	PlayMethod           FunctionId = 12000
	PlayMethod_CHALLENGE            = 12001
	PlayMethod_COIN                 = 12002
	PlayMethod_EQUIP                = 12003
	PlayMethod_WEAPON               = 12004
	PEAK_FIGHT                      = 37000 // 巅峰战场
	PET_Succinct                    = 38101 // 宠物洗练
	PET_Evolution                   = 38102 // 宠物进化
	World_Chat                      = 41000
	Friend                          = 43000 // 好友
	Add_Friend                      = 43001
	Arena                           = 101000
)

const (
	DiamondPay PayType = iota + 1 // 充钻石
	ShopPay                       // 商城支付
)

const (
	CardPoolActivity ActivityType = iota + 100 // 卡池活动
	LoginActivity
	ActivePass
	TaskPass
	OpenServer
	Desert   = 200 // 沙漠活动
	WeekPass       // 周常
)

const (
	Statics_Login_Id            uint32 = 306
	Guide_Id                           = 400
	Statics_Main_Mission_Id            = 500
	Statics_Challege_Mission_Id        = 501
	Statics_Coin_Mission_Id            = 502
	Statics_Equip_Mission_Id           = 503
	Statics_Weapon_Mission_Id          = 504
	Statics_Desert_Mission_Id          = 505
	Statisc_Union_Mission_Id           = 506
	Statics_Exit_Battle                = 507
	Statics_Buy_Shop_Item_Id           = 601
	Statics_Buy_Ap_Id                  = 602
	Statics_Lottey                     = 603
	Statics_Buy_OnHook                 = 604
	Statics_Arena_Pk                   = 605
	Statics_Add_New_Pet                = 606
	Statics_Pet_Evolution              = 607
	Statics_Add_Atlas                  = 608
	Statics_Add_Luck_Sale              = 609
)

type ChargeType uint32

const (
	Charge_Diamond   ChargeType = iota + 1 // 充钻石
	Charge_First                           // 首充
	Charge_MonthCard                       // 月卡
	Charge_MainFund
	Charge_ActivePass
	Charge_TaskPass
	Charge_MissionGift  = 7
	Charge_PerGift      = 8
	Charge_ActivitySign = 9
)

type AdType uint32

const (
	AD_Battle_Life                 = 1  // 战内复活
	Ad_Battle_Select_Weapon        = 2  // 战内武器重选
	Ad_Battle_Settle_Chip          = 3  // 战内结算筹码
	Ad_Add_Ap                      = 4  // 添加体力
	Ad_Add_OnHook                  = 5  // 快速挂机
	Ad_Coin_Sweep_Timess           = 6  // 战外-金币本扫荡次数*2
	Ad_Equip_Sweep_Times           = 7  // 战外-装备本扫荡次数*2
	Ad_Weapon_Sweep_times          = 8  // 战外-武器本扫荡次数*2
	Ad_Add_Free_Lottery            = 9  // 战外-招募免费单抽*1
	Ad_Shop_Buy                    = 10 // 商店购买
	Ad_Battle_Part_Shop            = 11 //配件商店
	Ad_Battle_Main_Settle          = 12 // 主线关卡结算
	Ad_Battle_Construction_Lottery = 13 // 交互物抽奖
)

type AttrType uint32

const (
	Attack          AttrType = 1   //攻击
	AttackRatio     AttrType = 2   //攻击系数
	Hp              AttrType = 3   //生命
	HpRatio         AttrType = 4   //生命系数
	Defense         AttrType = 5   //防御
	DefenseRatio    AttrType = 6   //防御系数
	AttackBonus     AttrType = 7   //攻击加成（自身）
	HpBonus         AttrType = 8   //生命加成（自身）
	DefenseBonus    AttrType = 9   //防御加成（自身）
	ColorBlack      AttrType = 11  //黑桃类型伤害加成
	ColorRed        AttrType = 12  //红桃类型伤害加成
	colorMeiHua     AttrType = 13  //草花类型伤害加成
	ColorFangKuai   AttrType = 14  //方块类型伤害加成
	AttackSpeed     AttrType = 15  //攻速
	Speed           AttrType = 16  //移动速度
	CritRate        AttrType = 18  //暴击率
	CritDamageRate  AttrType = 19  //暴击伤害
	Damage          AttrType = 20  //伤害加成
	MainSkillHurt   AttrType = 21  //主武器伤害加成
	FuSkillHurt     AttrType = 22  //副武器伤害加成
	LeaderHurt      AttrType = 25  //首领/精英伤害加成
	MonsterHurt     AttrType = 26  //小怪伤害加成
	SkillRange      AttrType = 27  //技能范围
	BulletSpeed     AttrType = 28  //子弹速度
	SkillTime       AttrType = 29  //技能持续时间
	Reduce          AttrType = 34  //免伤
	ReduceRate      AttrType = 37  //免伤率
	HpRecover       AttrType = 40  //生命恢复
	HpRecoverPer    AttrType = 41  //生命恢复百分比
	MpRecover       AttrType = 42  //能量每秒恢复值
	Collect         AttrType = 44  //拾取范围
	ExpRate         AttrType = 46  //经验获取效率
	MaxMp           AttrType = 47  //最大mp
	AttackBonusSup  AttrType = 107 //攻击加成（助战）
	HpBonusSup      AttrType = 108 //生命加成（助战）
	DefenseBonusSup AttrType = 109 //防御加成（助战）
	AttackBonusAll  AttrType = 207 //攻击加成（全局）
	HpBonusAll      AttrType = 208 //生命加成（全局）
	DefenseBonusAll AttrType = 209 //防御加成（全局）
)

type PvPStatus uint32

const (
	PvPDefault  PvPStatus = iota // 默认
	PvPMatching                  // 匹配中
	PvPBattle                    // 战斗中
)

type PetPart uint32

const (
	Pet_Part_Body PetPart = iota + 1
	Pet_Part_Head
	Pet_Part_Up
	Pet_Part_Down
	Pet_Part_Hand
	Pet_Part_Foot
)

// Alliance constants
const (
	ALLIANCE_INIT_MAX_MEMBER = 30   // 联盟初始最大成员数
	ALLIANCE_MAX_MEMBER      = 50   // 联盟最大成员数
	ALLIANCE_QUIT_CD         = 7200 // 退盟CD时间(秒)
	ALLIANCE_APPLY_CD        = 7200 // 申请CD时间(秒)
	ALLIANCE_CREATE_COST     = 1000 // 创建联盟消耗的钻石
)

// Alliance member positions
const (
	ALLIANCE_POS_LEADER = 1 // 指挥官
	ALLIANCE_POS_VICE   = 2 // 副指挥官
	ALLIANCE_POS_ELITE  = 3 // 精英
	ALLIANCE_POS_MEMBER = 4 // 普通成员
)

// Alliance application status
const (
	ALLIANCE_APPLY_PENDING = 0 // 待处理
	ALLIANCE_APPLY_ACCEPT  = 1 // 已同意
	ALLIANCE_APPLY_REJECT  = 2 // 已拒绝
)

// Alliance shop types
const (
	ALLIANCE_SHOP_DAILY = 1 // 每日商店
	ALLIANCE_SHOP_MONTH = 2 // 每月商店
)

// Alliance rank types
const (
	ALLIANCE_RANK_POWER  = 1 // 成员战力榜
	ALLIANCE_RANK_ACTIVE = 2 // 成员活跃榜
	ALLIANCE_RANK_BOSS   = 3 // BOSS伤害榜
	ALLIANCE_RANK_DAMAGE = 4 // 联盟总伤害榜
)

const (
	CURRENT_ALLIANCE_BOSS_ID uint32 = 1 // 当前联盟BOSS ID
)

type FlagType uint32

const (
	CanGet FlagType = iota + 1 // 可领取
	Recived
)
