package errcode

import (
	"fmt"
	"github.com/v587-zyf/gc/enums"
)

var (
	language    = enums.DEFAULT_LANGUAGE
	defaultErrs = make(ErrGroup)
)

type ErrCode int32
type ErrGroup map[ErrCode]map[enums.LANGUAGE]string

func SetLanguage(lang enums.LANGUAGE) {
	language = lang
}

type CodeLang struct {
	desc string
	lang enums.LANGUAGE
}

func NewCodeLang(desc string, lang enums.LANGUAGE) CodeLang {
	return CodeLang{
		desc: desc,
		lang: lang,
	}
}

func CreateErrCode(code int32, args ...CodeLang) ErrCode {
	if len(args) <= 0 {
		panic("create err code must have at least one language")
	}

	errCode := ErrCode(code)
	_, ok := defaultErrs[errCode]
	if !ok {
		defaultErrs[errCode] = make(map[enums.LANGUAGE]string)
	}

	for _, arg := range args {
		if _, ok = defaultErrs[errCode][arg.lang]; ok {
			msg := fmt.Sprintf("duplicate create err code, code:%d msg:%s lang:%s", code, arg.desc, arg.lang)
			panic(msg)
		}
		defaultErrs[errCode][arg.lang] = arg.desc
	}

	return errCode
}

func (code ErrCode) Error() string {
	if v, ok := defaultErrs[code][language]; !ok {
		return fmt.Sprintf("UNKNOW_ERR_CODE[%d]", code)
	} else {
		return v
	}
}

func (code ErrCode) Int() int {
	return int(code)
}
func (code ErrCode) Int32() int32 {
	return int32(code)
}
func (code ErrCode) UInt32() uint32 {
	return uint32(code)
}
