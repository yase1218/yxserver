package enum

import "kernel/iface"

// accountId, protoMsg
type GameGrpcMsgHandler func(int64, iface.IProtoMessage) (iface.IProtoMessage, error)
