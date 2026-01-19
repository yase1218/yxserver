package iface

import "google.golang.org/protobuf/proto"

type IProtoMessage interface {
	proto.Message

	Reset()
	String() string
	ProtoMessage()
	//Marshal() ([]byte, error)
	//MarshalTo([]byte) (int, error)
	//Unmarshal([]byte) error
	//Size() int
}

type MessageFrame struct {
	Len    uint32
	MsgID  uint16
	Tag    uint32
	UserID uint64
	Body   IProtoMessage
}
