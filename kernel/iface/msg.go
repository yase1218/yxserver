package iface

import "google.golang.org/protobuf/proto"

type IProtoMessage interface {
	proto.Message

	Reset()
	String() string
	ProtoMessage()
}
