package errcode

import (
	"github.com/v587-zyf/gc/enums"
)

var (
	ERR_SUCCEED      = CreateErrCode(0, NewCodeLang("成功", enums.LANG_CN), NewCodeLang("succeed", enums.LANG_EN))
	ERR_STANDARD_ERR = CreateErrCode(1, NewCodeLang("失败", enums.LANG_CN), NewCodeLang("failed", enums.LANG_EN))
	ERR_SIGN         = CreateErrCode(2, NewCodeLang("验证未通过", enums.LANG_CN), NewCodeLang("Verification failed", enums.LANG_EN))
	ERR_PARAM        = CreateErrCode(3, NewCodeLang("参数错误", enums.LANG_CN), NewCodeLang("The parameter is incorrect", enums.LANG_EN))
	ERR_CONFIG_NIL   = CreateErrCode(4, NewCodeLang("配置为空", enums.LANG_CN), NewCodeLang("The Config Is Nil", enums.LANG_EN))

	ERR_NET_SEND_TIMEOUT   = CreateErrCode(11, NewCodeLang("发送数据超时", enums.LANG_CN), NewCodeLang("The sending data timed out", enums.LANG_EN))
	ERR_NET_PKG_LEN_LIMIT  = CreateErrCode(12, NewCodeLang("数据包长度限制", enums.LANG_CN), NewCodeLang("Packet length limit", enums.LANG_EN))
	ERR_SERVER_INTERNAL    = CreateErrCode(13, NewCodeLang("服务器内部错误", enums.LANG_CN), NewCodeLang("Server internal error", enums.LANG_EN))
	ERR_WP_TOO_MANY_WORKER = CreateErrCode(14, NewCodeLang("工作池任务太多", enums.LANG_CN), NewCodeLang("There are too many work pool tasks", enums.LANG_EN))
	ERR_JSON_MARSHAL_ERR   = CreateErrCode(15, NewCodeLang("json打包错误", enums.LANG_CN), NewCodeLang("JSON packaging error", enums.LANG_EN))
	ERR_JSON_UNMARSHAL_ERR = CreateErrCode(16, NewCodeLang("json解包错误", enums.LANG_CN), NewCodeLang("JSON unpacking error", enums.LANG_EN))

	ERR_EVENT_PARAM_INVALID     = CreateErrCode(31, NewCodeLang("事件参数错误", enums.LANG_CN), NewCodeLang("Event parameter error", enums.LANG_EN))
	ERR_EVENT_LISTENER_LIMIT    = CreateErrCode(32, NewCodeLang("事件监听器数量限制", enums.LANG_CN), NewCodeLang("Event listener limit", enums.LANG_EN))
	ERR_EVENT_LISTENER_EMPTY    = CreateErrCode(33, NewCodeLang("事件监听器为空", enums.LANG_CN), NewCodeLang("Event listener is empty", enums.LANG_EN))
	ERR_EVENT_LISTENER_NOT_FIND = CreateErrCode(34, NewCodeLang("事件监听器未找到", enums.LANG_CN), NewCodeLang("Event listener not found", enums.LANG_EN))

	ERR_MQ_REPLY_HEAD_LEN      = CreateErrCode(51, NewCodeLang("mq回复头长度错误", enums.LANG_CN), NewCodeLang("mq reply header length error", enums.LANG_EN))
	ERR_MQ_BUFF_WRITE          = CreateErrCode(52, NewCodeLang("mq写入缓冲区错误", enums.LANG_CN), NewCodeLang("mq write buffer error", enums.LANG_EN))
	ERR_MQ_REPLY_EMPTY         = CreateErrCode(53, NewCodeLang("mq回复数据为空", enums.LANG_CN), NewCodeLang("mq reply data is empty", enums.LANG_EN))
	ERR_MQ_REPLY_PB            = CreateErrCode(54, NewCodeLang("mq回复数据pb错误", enums.LANG_CN), NewCodeLang("mq reply data pb error", enums.LANG_EN))
	ERR_MQ_RECV_DATA_UNMARSHAL = CreateErrCode(55, NewCodeLang("mq接收数据解析错误", enums.LANG_CN), NewCodeLang("mq receive data unmarshal error", enums.LANG_EN))
	ERR_MQ_MSG_ID_NOT_REGISTER = CreateErrCode(56, NewCodeLang("mq消息未注册", enums.LANG_CN), NewCodeLang("mq message not registered", enums.LANG_EN))
	ERR_MQ_CONNECT_FAIL        = CreateErrCode(57, NewCodeLang("mq连接失败", enums.LANG_CN), NewCodeLang("mq connect failed", enums.LANG_EN))
	ERR_MQ_REQ_TIMEOUT         = CreateErrCode(58, NewCodeLang("mq请求超时", enums.LANG_CN), NewCodeLang("mq request timeout", enums.LANG_EN))
	ERR_MQ_SERVER_NOT_FOUND    = CreateErrCode(59, NewCodeLang("mq服务器未找到", enums.LANG_CN), NewCodeLang("mq server not found", enums.LANG_EN))

	ERR_USER_DATA_NOT_FOUND  = CreateErrCode(3001, NewCodeLang("用户信息未找到", enums.LANG_CN), NewCodeLang("User information not found", enums.LANG_EN))
	ERR_USER_DATA_INVALID    = CreateErrCode(3002, NewCodeLang("用户信息错误", enums.LANG_CN), NewCodeLang("The user information is incorrect", enums.LANG_EN))
	ERR_REDIS_UPDATE_USER    = CreateErrCode(3003, NewCodeLang("redis更新玩家数据错误", enums.LANG_CN), NewCodeLang("Redis update player data error", enums.LANG_EN))
	ERR_REDIS_LOGIN_DATA_NIL = CreateErrCode(3004, NewCodeLang("redis登陆数据数据为空", enums.LANG_CN), NewCodeLang("The Redis login data is empty", enums.LANG_EN))
	ERR_MONGO_UPSERT         = CreateErrCode(3005, NewCodeLang("upsert错误", enums.LANG_CN), NewCodeLang("upsert error", enums.LANG_EN))
	ERR_MONGO_FIND           = CreateErrCode(3006, NewCodeLang("未找到数据", enums.LANG_CN), NewCodeLang("Data not found", enums.LANG_EN))
)
