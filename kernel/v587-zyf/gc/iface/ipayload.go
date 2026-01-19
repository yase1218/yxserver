package iface

import "io"

type IPayload interface {
	io.WriterTo
	Len() int
	CheckEncoding(enabled bool, opcode uint8) bool
}
