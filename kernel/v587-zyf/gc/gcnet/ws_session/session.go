package ws_session

import (
	"context"
	"kernel/tools"
	"math"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/buffer_pool"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/ws_session_mgr"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type Session struct {
	id   uint64
	conn *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc

	hooks *Hooks
	cache sync.Map
	//method iface.IWsSessionMethod

	outChan chan []byte
	isClose bool

	heartbeatTime time.Time
}

func NewSession(ctx context.Context, conn *websocket.Conn) *Session {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ctx:    ctx,
		cancel: cancel,

		outChan: make(chan []byte, 1024),

		hooks: NewHooks(),

		heartbeatTime: time.Now(),
	}
	s.conn = conn

	return s
}

func (s *Session) Start() {
	s.hooks.ExecuteStart(s)

	go tools.GoSafe("ws_session read pump", func() {
		s.readPump()
	})

	go tools.GoSafe("ws_session io pump", func() {
		s.IOPump()
	})
}

func (s *Session) Hooks() *Hooks {
	return s.hooks
}

func (s *Session) Set(key string, value any) {
	s.cache.Store(key, value)
}
func (s *Session) Get(key string) (any, bool) {
	return s.cache.Load(key)
}
func (s *Session) Remove(key string) {
	s.cache.Delete(key)
}

func (s *Session) GetID() uint64 {
	return s.id
}
func (s *Session) SetID(id uint64) {
	if id <= 0 {
		id = 0
	}
	s.id = id
}

func (s *Session) Close() error {
	log.Info("seesion close", zap.Uint64("session id", s.GetID()))
	if !s.isClose {
		s.isClose = true

		s.hooks.ExecuteStop(s)

		close(s.outChan)

		s.cancel()
		s.conn.Close()

		ws_session_mgr.GetSessionMgr().UnRegisterCh <- s
	}

	return nil
}

func (s *Session) GetConn() iface.IConn {
	return s.conn
}

func (s *Session) GetRemoteAddr() string {
	if s.conn == nil {
		return ""
	}
	return s.conn.RemoteAddr().String()
}

func (s *Session) GetCtx() context.Context {
	return s.ctx
}

func (s *Session) Login() {
	ws_session_mgr.GetSessionMgr().LoginCh <- s
}

func (s *Session) DoSomething(fn func(args ...any) bool) bool {
	return fn()
}
func (s *Session) CheckSomething(fn func(args ...any) bool) bool {
	return fn()
}
func (s *Session) Heartbeat() {
	s.heartbeatTime = time.Now()
}

func (s *Session) IsHeartbeatTimeout(now time.Time) bool {
	return now.After(s.heartbeatTime.Add(enums.HEARTBEAT_TIMEOUT))
}

func (s *Session) SendMsg(fn func(args ...any) ([]byte, error), args ...any) error {
	if s.isClose {
		return nil
	}
	sendBytes, err := fn(args...)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		select {
		case s.outChan <- sendBytes:
			return nil
		default:
			backoff := time.Duration(100) * time.Millisecond * time.Duration(2^i)
			time.Sleep(backoff)
		}
	}

	return errcode.ERR_NET_SEND_TIMEOUT
}

func (s *Session) readPump() {
LOOP:
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
			}
			log.Error("ws_server read err", zap.Uint64("session id", s.GetID()), zap.Error(err))
			break LOOP
		}
		if len(message) > enums.MAX_MSG_SIZE {
			log.Warn("消息超过最大长度，忽略", zap.Int("length", len(message)))
			continue
		}

		if message != nil && len(message) > 0 {
			buf := buffer_pool.GetBuffer()
			if buf == nil {
				log.Error("buffer_pool get err", zap.Uint64("session id", s.GetID()))
				break LOOP
			}

			buf.Data = ensureCapacity(buf.Data, len(message))
			copy(buf.Data, message)

			s.hooks.ExecuteRecv(s, buf.Data)
			buffer_pool.Put(buf)
		}
	}
	log.Info("readPump exit", zap.Uint64("session id", s.GetID()))
	s.cancel()
}

func ensureCapacity(slice []byte, size int) []byte {
	if cap(slice) >= size {
		return slice[:size]
	}
	newSlice := make([]byte, size)
	if len(slice) < size {
		copy(newSlice, slice)
	}
	return newSlice
}

func (s *Session) IOPump() {
	var (
		err     error
		backoff time.Duration
	)

LOOP:
	for {
		select {
		case data := <-s.outChan:
			for i := 0; i < 3; i++ {
				if err = s.conn.WriteMessage(websocket.BinaryMessage, data); err == nil {
					break
				}
				backoff = calculateBackoff(i)
				time.Sleep(backoff)
			}
			if err != nil {
				log.Info("WriteMessage err", zap.Uint64("session id", s.GetID()), zap.Error(err))
				break LOOP
			}
		case <-s.ctx.Done():
			log.Info("context done", zap.Uint64("session id", s.GetID()), zap.Error(err))
			break LOOP
		}
	}

	log.Info("session loop exit", zap.Uint64("session id", s.GetID()))
	s.Close()
}

func calculateBackoff(attempt int) time.Duration {
	return time.Duration(100) * time.Millisecond * time.Duration(math.Min(math.Pow(2, float64(attempt)), float64(time.Second/time.Millisecond)))
}
