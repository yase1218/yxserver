package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

// 7001-8000
var (
	ERR_RATE      = errcode.CreateErrCode(70001, errcode.NewCodeLang("请求过快", enums.LANG_CN))
	ERR_NO_SERVER = errcode.CreateErrCode(70002, errcode.NewCodeLang("暂无服务器可以连接", enums.LANG_CN))
	//ERR_MAINTENANCE    = errcode.CreateErrCode(70003, errcode.NewCodeLang("服务器维护中", enums.LANG_CN))
	ERR_WHITE_LIST     = errcode.CreateErrCode(70004, errcode.NewCodeLang("白名单未通过", enums.LANG_CN))
	ERR_REGISTER_LIMIT = errcode.CreateErrCode(70005, errcode.NewCodeLang("服务器已满", enums.LANG_CN))

	ERR_SDK_LOGIN_FAIL         = errcode.CreateErrCode(1000, errcode.NewCodeLang("sdk登录失败", enums.LANG_CN))
	ERR_REGIST_FULL            = errcode.CreateErrCode(1001, errcode.NewCodeLang("服务器已满，请前往新服体验", enums.LANG_CN))
	ERR_SYSTEM_FAIL            = errcode.CreateErrCode(1002, errcode.NewCodeLang("系统错误", enums.LANG_CN))
	ERR_MAINTENANCE            = errcode.CreateErrCode(1003, errcode.NewCodeLang("系统维护中", enums.LANG_CN))
	ERR_WHITELIST              = errcode.CreateErrCode(1004, errcode.NewCodeLang("已开启白名单登录", enums.LANG_CN))
	ERR_SDK_CHARGE_VERIFY_FAIL = errcode.CreateErrCode(1005, errcode.NewCodeLang("sdk订单验证失败", enums.LANG_CN))
	ERR_ILLEGAL_SERVER         = errcode.CreateErrCode(1006, errcode.NewCodeLang("无效服务器", enums.LANG_CN))
)
