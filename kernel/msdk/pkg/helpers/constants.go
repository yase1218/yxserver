package helpers

// Platform 平台定义
const (
	PlatformGarena   uint8 = 1
	PlatformBeeTalk  uint8 = 2
	PlatformFacebook uint8 = 3
	PlatformGuest    uint8 = 4
	PlatformVK       uint8 = 5
	PlatformLine     uint8 = 6
	PlatformHuawei   uint8 = 7
	PlatformGoogle   uint8 = 8
	PlatformWeChat   uint8 = 9
	PlatformApple    uint8 = 10
	PlatformTwitter  uint8 = 11
	PlatformEmail    uint8 = 13
	PlatformPGS      uint8 = 14
)

// Gender 性别定义
const (
	GenderUnknown uint8 = 0
	GenderMale    uint8 = 1
	GenderFemale  uint8 = 2
)

// ClientType 客户端类型定义
const (
	ClientTypeUnknown uint8 = 0
	ClientTypeIOS     uint8 = 1
	ClientTypeAndroid uint8 = 2
	ClientTypePC      uint8 = 3
)

// 平台名称映射
var PlatformNames = map[uint8]string{
	PlatformGarena:   "Garena",
	PlatformBeeTalk:  "BeeTalk",
	PlatformFacebook: "Facebook",
	PlatformGuest:    "Guest",
	PlatformVK:       "VK",
	PlatformLine:     "Line",
	PlatformHuawei:   "Huawei",
	PlatformGoogle:   "Google",
	PlatformWeChat:   "WeChat",
	PlatformApple:    "Apple",
	PlatformTwitter:  "Twitter",
	PlatformEmail:    "Email",
	PlatformPGS:      "PGS",
}

// 性别名称映射
var GenderNames = map[uint8]string{
	GenderUnknown: "未知",
	GenderMale:    "男性",
	GenderFemale:  "女性",
}

// 客户端类型名称映射
var ClientTypeNames = map[uint8]string{
	ClientTypeUnknown: "未知",
	ClientTypeIOS:     "iOS",
	ClientTypeAndroid: "Android",
	ClientTypePC:      "PC",
}
