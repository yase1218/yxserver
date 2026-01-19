package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

// 9001-10000
var (
	ERR_FIGHT_NOT_FOUND            = errcode.CreateErrCode(93001, errcode.NewCodeLang("战斗未找到", enums.LANG_CN))
	ERR_MSG_NOT_HANDLER            = errcode.CreateErrCode(90002, errcode.NewCodeLang("战斗协议未找到", enums.LANG_CN))
	ERR_SCENE_DATA                 = errcode.CreateErrCode(93003, errcode.NewCodeLang("场景数据有误", enums.LANG_CN))
	ERR_SCENE_CANT_NOT_WALK        = errcode.CreateErrCode(93004, errcode.NewCodeLang("该地点无法移动", enums.LANG_CN))
	ERR_ACTOR_NOT_FOUND            = errcode.CreateErrCode(93005, errcode.NewCodeLang("未找到Actor", enums.LANG_CN))
	ERR_ACTOR_CANT_NOT_MOVE        = errcode.CreateErrCode(93006, errcode.NewCodeLang("Actor不能移动", enums.LANG_CN))
	ERR_MOVE                       = errcode.CreateErrCode(93007, errcode.NewCodeLang("移动失败", enums.LANG_CN))
	ERR_SCENE_POINT_NIL            = errcode.CreateErrCode(93008, errcode.NewCodeLang("地图指定点位为空", enums.LANG_CN))
	ERR_SKILL_CD                   = errcode.CreateErrCode(93009, errcode.NewCodeLang("技能CD", enums.LANG_CN))
	ERR_SKILL_NOT_FOUND            = errcode.CreateErrCode(93010, errcode.NewCodeLang("技能未找到", enums.LANG_CN))
	ERR_SKILL_BEING_RELEASE        = errcode.CreateErrCode(93011, errcode.NewCodeLang("已有技能在释放", enums.LANG_CN))
	ERR_SKILL_CAN_NOT_USE          = errcode.CreateErrCode(93012, errcode.NewCodeLang("技能无法释放", enums.LANG_CN))
	ERR_SKILL_CONDITION_ERROR      = errcode.CreateErrCode(93013, errcode.NewCodeLang("技能condition参数异常", enums.LANG_CN))
	ERR_SKILL_CONDITION_FAIL       = errcode.CreateErrCode(93014, errcode.NewCodeLang("技能condition未通过", enums.LANG_CN))
	ERR_BUFF_TYPE                  = errcode.CreateErrCode(93015, errcode.NewCodeLang("Buff类型错误", enums.LANG_CN))
	ERR_ATTACK                     = errcode.CreateErrCode(93016, errcode.NewCodeLang("攻击失败", enums.LANG_CN))
	ERR_CREATE_SKILL               = errcode.CreateErrCode(93017, errcode.NewCodeLang("创建技能", enums.LANG_CN))
	ERR_CAN_NOT_INTERACT           = errcode.CreateErrCode(93018, errcode.NewCodeLang("不能交互", enums.LANG_CN))
	ERR_CHIP_NOT_ENOUGH            = errcode.CreateErrCode(93019, errcode.NewCodeLang("筹码不足", enums.LANG_CN))
	ERR_FIGHT_SHOP_CAN_NOT_REFRESH = errcode.CreateErrCode(93020, errcode.NewCodeLang("战斗商店无法刷新", enums.LANG_CN))
	ERR_NO_POKER_REWARD            = errcode.CreateErrCode(93021, errcode.NewCodeLang("当前牌组无奖励", enums.LANG_CN))
	ERR_FIGHT_GM_ERR               = errcode.CreateErrCode(93022, errcode.NewCodeLang("GM指令无效", enums.LANG_CN))
	ERR_REVIVE_NO_NUM              = errcode.CreateErrCode(93023, errcode.NewCodeLang("无复活次数", enums.LANG_CN))
	ERR_PLAYER_DIE                 = errcode.CreateErrCode(93024, errcode.NewCodeLang("玩家已死亡", enums.LANG_CN))
	ERR_RAND_PLAYER_NIL            = errcode.CreateErrCode(93025, errcode.NewCodeLang("随机玩家为空", enums.LANG_CN))
	ERR_ILLEGAL_OPERATION          = errcode.CreateErrCode(93026, errcode.NewCodeLang("非法操作", enums.LANG_CN))
	ERR_BUFF_ALREADY_ADD           = errcode.CreateErrCode(93027, errcode.NewCodeLang("BUFF已添加", enums.LANG_CN))
)
