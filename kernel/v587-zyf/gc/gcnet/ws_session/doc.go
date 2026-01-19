package ws_session

import (
	"github.com/v587-zyf/gc/iface"
)

type Recv func(conn iface.IWsSession, data any)

type Call func(ss iface.IWsSession)
