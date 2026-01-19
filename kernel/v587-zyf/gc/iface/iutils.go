package iface

type Integer interface {
	int | int64 | int32 | uint | uint64 | uint32
}
