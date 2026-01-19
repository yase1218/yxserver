package enums

import "time"

const (
	RDB_KEY_RECONNECT         = "Reconnect"                    // 重连
	RDB_KEY_RECONNECT_TIME    = time.Minute * 2                // 重连时间
	RDB_KEY_SER_GATE          = "Thunder{service}GateStress"   // 网关服压力
	RDB_KEY_SER_GATE_ReqCount = "Thunder{service}GateReqCount" // 网关请求数量
)
