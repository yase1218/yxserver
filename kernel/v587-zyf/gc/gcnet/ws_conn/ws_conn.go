package ws_conn

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/buffer_pool"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/ws_session_mgr"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/tools"
	"math"
	"sync"
	"time"
)

type Conn struct {
	id   uint64
	conn iface.IConn

	ctx    context.Context
	cancel context.CancelFunc

	hooks *Hooks
	cache sync.Map
	//method iface.IWsConnMethod

	outChan chan []byte
	isClose bool
	mu      sync.RWMutex

	heartbeatTime time.Time
}

func NewConn(ctx context.Context, conn iface.IConn) *Conn {
	ctx, cancel := context.WithCancel(ctx)
	s := &Conn{
		ctx:    ctx,
		cancel: cancel,

		outChan: make(chan []byte, 1024),

		hooks: NewHooks(),

		heartbeatTime: time.Now(),
	}
	s.conn = conn

	return s
}

func (s *Conn) Start() {
	s.hooks.ExecuteStart(s)

	go tools.GoSafe("ws_conn read pump", func() {
		s.readPump()
	})
	go tools.GoSafe("ws_conn io pump", func() {
		s.IOPump()
	})
}

func (s *Conn) Hooks() *Hooks {
	return s.hooks
}

func (s *Conn) Set(key string, value any) {
	s.cache.Store(key, value)
}
func (s *Conn) Get(key string) (any, bool) {
	return s.cache.Load(key)
}
func (s *Conn) Remove(key string) {
	s.cache.Delete(key)
}

func (s *Conn) GetID() uint64 {
	return s.id
}
func (s *Conn) SetID(id uint64) {
	if id <= 0 {
		id = 0
	}
	s.id = id
}

func (s *Conn) Close() error {
	//fmt.Println("close-----------")
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *Conn) GetConn() iface.IConn {
	return s.conn
}

func (s *Conn) GetCtx() context.Context {
	return s.ctx
}

func (s *Conn) Login() {
	ws_session_mgr.GetSessionMgr().LoginCh <- s
}

func (s *Conn) DoSomething(fn func(args ...any) bool) bool {
	return fn()
}
func (s *Conn) CheckSomething(fn func(args ...any) bool) bool {
	return fn()
}
func (s *Conn) Heartbeat() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.heartbeatTime = time.Now()
}

func (s *Conn) IsHeartbeatTimeout(now time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return time.Now().After(s.heartbeatTime.Add(enums.HEARTBEAT_TIMEOUT))
}

func (s *Conn) SendMsg(fn func(args ...any) ([]byte, error), args ...any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *Conn) readPump() {
LOOP:
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				//log.Info("ws_server read err", zap.Error(err))
			}
			break LOOP
		}
		if len(message) > enums.MAX_MSG_SIZE {
			log.Warn("消息超过最大长度，忽略", zap.Int("length", len(message)))
			continue
		}

		if message != nil && len(message) > 0 {
			buf := buffer_pool.GetBuffer()
			if buf == nil {
				log.Error("buffer_pool get err")
				break LOOP
			}

			buf.Data = ensureCapacity(buf.Data, len(message))
			copy(buf.Data, message)

			s.hooks.ExecuteRecv(s, buf.Data)
			buffer_pool.Put(buf)
		}
	}

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

func (s *Conn) IOPump() {
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
				break LOOP
			}
		case <-s.ctx.Done():
			for len(s.outChan) > 0 {
				select {
				case data := <-s.outChan:
					if err = s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
						break LOOP
					}
				case <-time.After(time.Second):
					log.Warn("timeout waiting for messages to send")
					break LOOP
				}
			}
			break LOOP
		}
	}

	s.Close()
}

func calculateBackoff(attempt int) time.Duration {
	return time.Duration(100) * time.Millisecond * time.Duration(math.Min(math.Pow(2, float64(attempt)), float64(time.Second/time.Millisecond)))
}
