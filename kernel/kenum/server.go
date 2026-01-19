package kenum

const (
	SER_GATE_PREFIX        = "/server/gate/"
	SER_GATE_PREFIX_Format = "/server/gate/%d"

	SER_PREFIX              = "/server/"
	SER_GAME_PREFIX         = "/server/game/"
	SER_GAME_PREFIX_Format  = "/server/game/%d"
	SER_FIGHT_PREFIX        = "/server/fight/"
	SER_FIGHT_PREFIX_Format = "/server/fight/%d"
)

const (
	SER_LOGIN_REGISTER_LIMIT = 3000
)

// 工作状态
type WorkState = uint32

const (
	WorkState_Idle     WorkState = 0
	WorkState_Running  WorkState = 1
	WorkState_Stopping WorkState = 2
	WorkState_Stopped  WorkState = 3
)

func StateToString(state WorkState) string {
	switch state {
	case WorkState_Idle:
		return "idle"
	case WorkState_Running:
		return "running"
	case WorkState_Stopping:
		return "stopping"
	case WorkState_Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}
