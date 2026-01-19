package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

// 8001-8000
var (
	ERR_USER_IN_FIGHT                 = errcode.CreateErrCode(80001, errcode.NewCodeLang("玩家在战斗中", enums.LANG_CN))
	ERR_FIGHT_BEFORE_CHECK            = errcode.CreateErrCode(80002, errcode.NewCodeLang("战前检查未通过", enums.LANG_CN))
	ERR_USER_NOT_IN_FIGHT             = errcode.CreateErrCode(80003, errcode.NewCodeLang("用户不在战斗中", enums.LANG_CN))
	ERR_MATCH                         = errcode.CreateErrCode(80004, errcode.NewCodeLang("匹配失败", enums.LANG_CN))
	ERR_AP_NOT_ENOUGH                 = errcode.CreateErrCode(80005, errcode.NewCodeLang("体力不足", enums.LANG_CN))
	ERR_BAG_NOT_ENOUGH                = errcode.CreateErrCode(80006, errcode.NewCodeLang("背包物品不足", enums.LANG_CN))
	ERR_PEAK_FIGHT_TICKETS_NOT_ENOUGH = errcode.CreateErrCode(80007, errcode.NewCodeLang("巅峰战场门票不足", enums.LANG_CN))
	ERR_REPEATE_REWARD                = errcode.CreateErrCode(80008, errcode.NewCodeLang("重复领取奖励", enums.LANG_CN))
	ERR_TASK_NOT_COMPLETE             = errcode.CreateErrCode(80009, errcode.NewCodeLang("任务未完成", enums.LANG_CN))
)
