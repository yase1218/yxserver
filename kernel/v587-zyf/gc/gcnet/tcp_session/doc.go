package tcp_session

import (
	"github.com/v587-zyf/gc/iface"
)

type Recv func(conn iface.ITcpSession, data any)

type Call func(ss iface.ITcpSession)
