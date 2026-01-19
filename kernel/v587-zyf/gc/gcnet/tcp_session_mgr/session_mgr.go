package tcp_session_mgr

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

var sessionMgr *SessionMgr

func init() {
	sessionMgr = NewSessionMgr()
}

type SessionMgr struct {
	allClients sync.Map // iface.ITcpSession:struct{}
	allClientN int64

	onlineClients sync.Map // uint64:iface.ITcpSession
	onlineClientN int64

	RegisterCh   chan iface.ITcpSession
	LoginCh      chan iface.ITcpSession
	UnRegisterCh chan iface.ITcpSession
}

func GetSessionMgr() *SessionMgr {
	return sessionMgr
}

func NewSessionMgr() *SessionMgr {
	s := &SessionMgr{}

	return s
}

func (s *SessionMgr) AllLength() int {
	return int(atomic.LoadInt64(&s.allClientN))
}

func (s *SessionMgr) IsConn(ss iface.ITcpSession) (ok bool) {
	_, ok = s.allClients.Load(ss)

	return
}

func (s *SessionMgr) GetAll() (allSS map[iface.ITcpSession]struct{}) {
	allSS = make(map[iface.ITcpSession]struct{})

	s.AllRange(func(ss iface.ITcpSession) (result bool) {
		allSS[ss] = struct{}{}
		return true
	})

	return
}

func (s *SessionMgr) AllRange(fn func(ss iface.ITcpSession) (result bool)) {
	s.allClients.Range(func(key, value any) bool {
		ss := value.(iface.ITcpSession)
		result := fn(ss)
		if !result {
			return false
		}
		return true
	})
}

func (s *SessionMgr) AllAdd(ss iface.ITcpSession) {
	s.allClients.Store(ss, struct{}{})
	atomic.AddInt64(&s.allClientN, 1)
}

func (s *SessionMgr) AllDel(ss iface.ITcpSession) {
	s.allClients.Delete(ss)
	atomic.AddInt64(&s.allClientN, -1)
}

func (s *SessionMgr) OnlineLen() int {
	return int(atomic.LoadInt64(&s.onlineClientN))
}

func (s *SessionMgr) OnlineAdd(userID uint64, ss iface.ITcpSession) {
	s.onlineClients.Store(userID, ss)
	atomic.AddInt64(&s.onlineClientN, 1)
}

func (s *SessionMgr) OnlineDel(userID uint64) {
	s.onlineClients.Delete(userID)
	atomic.AddInt64(&s.onlineClientN, -1)
}

func (s *SessionMgr) OnlineGetOne(userID uint64) (ss iface.ITcpSession) {
	if v, ok := s.onlineClients.Load(userID); ok {
		ss = v.(iface.ITcpSession)
	}

	return
}

func (s *SessionMgr) OnlineOnce(userID uint64, fn func(ss iface.ITcpSession)) {
	cli := s.OnlineGetOne(userID)
	if cli == nil {
		return
	}

	fn(cli)
}

func (s *SessionMgr) OnlineRange(fn func(userID uint64, ss iface.ITcpSession)) {
	s.onlineClients.Range(func(key, value any) bool {
		userID := key.(uint64)
		ss := value.(iface.ITcpSession)
		fn(userID, ss)
		return true
	})

	return
}

func (s *SessionMgr) IsOnline(userID uint64) (ss iface.ITcpSession, ok bool) {
	if val, ok := s.onlineClients.Load(userID); ok {
		ss = val.(iface.ITcpSession)
		return ss, true
	}
	return
}

func (s *SessionMgr) Login(ss iface.ITcpSession) {
	if s.IsConn(ss) {
		s.OnlineAdd(ss.GetID(), ss)
	}
}
func (s *SessionMgr) Disconnect(ss iface.ITcpSession) {
	log.Info("session disconnect", zap.Uint64("session id", ss.GetID()))
	if ss.GetID() != 0 {
		s.OnlineDel(ss.GetID())
	}

	s.AllDel(ss)
}

func (s *SessionMgr) Start() {
	for {
		select {
		case ss := <-s.RegisterCh:
			s.AllAdd(ss)
		case ss := <-s.LoginCh:
			s.Login(ss)
		case ss := <-s.UnRegisterCh:
			s.Disconnect(ss)
		}
	}
}

func (s *SessionMgr) ClearTimeout() {
	currentTime := time.Now()

	allClients := s.GetAll()
	for session := range allClients {
		var fn = func(args ...any) bool {
			log.Info("session time out", zap.Uint64("session id", session.GetID()))
			return session.IsHeartbeatTimeout(currentTime)
		}
		if session.CheckSomething(fn) {
			session.Close()
		}
	}
}
