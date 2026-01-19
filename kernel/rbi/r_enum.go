package rbi

type IRbi interface {
	Name() string
}

const (
	Rbi_Server_Url       = "10.128.20.41"
	Rbi_Server_Port      = "37500"
	Rbi_Server_Port_Test = "37510"
)

var PlatMap = map[string]int{
	"ios":       0,
	"android":   1,
	"pc":        2,
	"harmonyos": 3,
}
