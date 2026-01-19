package enum

type FightPeakEventType int

const (
	Fight_Peak_Event_Type_Join_Room   FightPeakEventType = iota // 进入房间
	Fight_Peak_Event_Type_Exit_Room                             // 离开房间
	Fight_Peak_Event_Type_Enter_Fight                           // 进入战斗
	Fight_Peak_Event_Type_Pick_Item                             // 拾取物品
	Fight_Peak_Event_Type_Attack_Boss                           // 攻击BOSS
	Fight_Peak_Event_Type_Exit_Fight                            // 离开战斗
)

const (
	Fight_Peak_Room_Robot_Check_Seconds         = 3       // 每次计算机器人进度秒数
	Fight_Peak_Room_Timeout_Seconds             = 3       // 房间匹配最大时间,超出时间自动添加机器人
	Fight_Peak_Room_User_Num_Max                = 4       // 房间最大人数
	Fight_Peak_Room_Destroy_Seconds             = 60 * 10 // 房间最大销毁时间
	Fight_Peak_Room_Robot_Start_Timeout_Seconds = 10      // 这个时间内没有玩家更新进度,机器人也开始计算进度

	Fight_Peak_Robot_Speed = 3 // 机器人计算速度
)

type FightPeakStairType int

const (
	Fight_Peak_Stair_Pick_Item   = iota + 1 // 拾取阶段
	Fight_Peak_Stair_Attack_Boss            // 攻击boss阶段
)

//type FightPeakExitType int
//
//const (
//	Fight_Peak_Exit_Type_Exit FightPeakExitType = iota
//	Fight_Peak_Exit_Type_Timeout
//	Fight_Peak_Exit_Type_Dead
//)

type FightPeakStatus uint32

const (
	Fight_Peak_Status_Fighting FightPeakStatus = iota + 1 // 战斗中
	Fight_Peak_Status_Dead                                // 死亡
	Fight_Peak_Status_Finish                              // 完成
)
