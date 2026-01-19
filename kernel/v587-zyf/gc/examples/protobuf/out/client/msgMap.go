package pb
import (
	"reflect"
)

var msgProtoTypes = make(map[uint16]reflect.Type)
var msgNames = make(map[uint16]string)

func init() {
	msgProtoTypes[MsgID_MyMessageId]=reflect.TypeOf((*MyMessage)(nil)).Elem()
	msgNames[MsgID_MyMessageId]="MyMessage"
}

func GetMsgProtoType(key uint16) reflect.Type {
	return msgProtoTypes[key]
}

func GetMsgName(key uint16) string {
	return msgNames[key]
}

const (
	MsgID_MyMessageId=0
)

func GetMsgIdFromType(i interface{}) uint16 {
	switch i.(type) {
	case *MyMessage:
		return 0
	default:
		return 0
	}
}
