package kenum

const (
	Redis_Key_User      = "User:%v:%d"
	Redis_Key_Whitelist = "game:white"
	Redis_Key_Token     = "token"
	//Redis_Key_Whitelist = "white_list"
)

const (
	Redis_Login_Token         = "Token"
	Redis_Login_Gate          = "Gate"
	Redis_Login_UID           = "UID"
	Redis_Login_Tda_Comm_Attr = "Tda_Comm_Attr"
)

const (
	Redis_Login_Token_Index = iota
	Redis_Login_Gate_Index
	Redis_Login_Uid_Index
	Redis_Login_Tda_Comm_Attr_Index
)

const (
	RECONNECT_MAX_TIMES = 3
)
