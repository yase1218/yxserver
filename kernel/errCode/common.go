package errCode

import (
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
)

func init() {
	errcode.SetLanguage(enums.LANG_CN)
}

var (
	ERR_MARSHAL   = errcode.CreateErrCode(10001, errcode.NewCodeLang("marshal失败", enums.LANG_CN))
	ERR_UNMARSHAL = errcode.CreateErrCode(10002, errcode.NewCodeLang("unmarshal失败", enums.LANG_CN))
)
