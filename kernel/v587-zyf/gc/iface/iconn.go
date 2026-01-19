package iface

type IConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type IWsConnMethod interface {
	Name() string
	Start(ss IWsSession)
	Recv(conn IWsSession, data any)
	Stop(ss IWsSession)
}
