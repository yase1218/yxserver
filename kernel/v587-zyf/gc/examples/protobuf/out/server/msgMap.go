package server
import (
	"reflect"
)

var msgProtoTypes = make(map[uint16]reflect.Type)
var msgNames = make(map[uint16]string)

func init() {
	msgProtoTypes[MsgID_Hello_RequestId]=reflect.TypeOf((*HelloRequest)(nil)).Elem()
	msgProtoTypes[MsgID_Hello_ResponseId]=reflect.TypeOf((*HelloResponse)(nil)).Elem()
	msgNames[MsgID_Hello_RequestId]="HelloRequest"
	msgNames[MsgID_Hello_ResponseId]="HelloResponse"
}

func GetMsgProtoType(key uint16) reflect.Type {
	return msgProtoTypes[key]
}

func GetMsgName(key uint16) string {
	return msgNames[key]
}

const (
	MsgID_Hello_RequestId=0
	MsgID_Hello_ResponseId=1
)

func GetMsgIdFromType(i interface{}) uint16 {
	switch i.(type) {
	case *HelloRequest:
		return 0
	case *HelloResponse:
		return 1
	default:
		return 0
	}
}
