package tapping

type TapType = uint32

const (
	TapType_Nil TapType = iota
	TapType_Tda
	TapType_Local

	TapType_All
)

type ITapEnvet interface {
	EventName() string
}

type TapLogin struct {
	AccountId string
	Uid       uint64
}

func (t *TapLogin) EventName() string {
	return "login_event"
}
