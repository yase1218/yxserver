package kenum

const (
	Msdk_Dev_Domain  = "https://%v.connect.garena.com"
	Msdk_Test_Domain = "https://testconnect.garena.com"

	// oauth
	Msdk_Oauth_Inspect_Url          = "/oauth/token/inspect"            // 检查token
	Msdk_Oauth_Get_User_Url         = "oauth/user/info/get"             // 获取用户信息
	Msdk_Oauth_Get_Friend_Url       = "oauth/user/friends/get/v2"       // 获取好友列表
	Msdk_Oauth_Get_Friend_Inapp_Url = "oauth/user/friends/inapp/get/v2" // 获取应用好友列表
	Msdk_Oauth_Get_Friend_Info_Url  = "oauth/user/friends/info/get/v2"  // 获取应用好友信息
	Msdk_Oauth_Get_Role_Url         = "oauth/user/role"                 // 获取用户角色

	// game
	Msdk_Game_Send_Friend_Url = "/game/user/request/send"  // 向朋友发送
	Msdk_Game_Guest_Swap_Url  = "/game/guest/swap"         // 绑定游戏账户
	Msdk_Game_Logout_All_Url  = "/game/logout_all_devices" // 注销所有设备

	// bind
	Msdk_Bind_Platform_Create_Url   = "/bind/app/platform/create"   // 绑定平台账户
	Msdk_Bind_Platform_Get_Info_Url = "/bind/app/platform/info/get" // 查询绑定平台账户
	Msdk_Bind_Platform_Del_Url      = "/bind/app/platform/delete"   // 删除绑定平台账户
)

var Msdk_Platform_Name = map[uint32]string{
	Msdk_Platform_Test:     "test",
	Msdk_Platform_Garena:   "Garena",
	Msdk_Platform_BeeTalk:  "BeeTalk",
	Msdk_Platform_Facebook: "Facebook",
	Msdk_Platform_Guest:    "Guest",
	Msdk_Platform_VK:       "VK",
	Msdk_Platform_Line:     "Line",
	Msdk_Platform_Huawei:   "Huawei",
	Msdk_Platform_Google:   "Google",
	Msdk_Platform_WeChat:   "WeChat",
	Msdk_Platform_Apple:    "Apple",
	Msdk_Platform_Twitter:  "Twitter",
	Msdk_Platform_Email:    "Email",
	Msdk_Platform_PGS:      "PGS",
	Msdk_Platform_TAP:      "TAPTAP",
	Msdk_Platform_Leiting:  "Leiting",
}

const (
	Msdk_Platform_Test = iota // 测试
	Msdk_Platform_Garena
	Msdk_Platform_BeeTalk
	Msdk_Platform_Facebook
	Msdk_Platform_Guest
	Msdk_Platform_VK
	Msdk_Platform_Line
	Msdk_Platform_Huawei
	Msdk_Platform_Google
	Msdk_Platform_WeChat
	Msdk_Platform_Apple
	Msdk_Platform_Twitter
	Msdk_Platform_Email
	Msdk_Platform_PGS

	Msdk_Platform_TAP     = 1000
	Msdk_Platform_Leiting = 2000
)

var Msdk_Client_Name = map[uint8]string{
	Msdk_Client_Unkown:  "Unkown",
	Msdk_Client_iOS:     "iOS",
	Msdk_Client_Android: "Android",
	Msdk_Client_PC:      "PC",
}

const (
	Msdk_Client_Unkown = iota + 1
	Msdk_Client_iOS
	Msdk_Client_Android
	Msdk_Client_PC
)

var Msdk_Gender_Name = map[uint8]string{
	Msdk_Gender_Unknown: "Unknown",
	Msdk_Gender_Male:    "Male",
	Msdk_Gender_Female:  "Female",
}

const (
	Msdk_Gender_Unknown = iota + 1
	Msdk_Gender_Male
	Msdk_Gender_Female
)
