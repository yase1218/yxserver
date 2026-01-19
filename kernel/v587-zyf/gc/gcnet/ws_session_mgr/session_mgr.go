package ws_session_mgr

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
	allClients sync.Map // iface.IWsSession:struct{}
	allClientN int64

	onlineClients sync.Map // uint64:iface.IWsSession
	onlineClientN int64

	RegisterCh   chan iface.IWsSession
	LoginCh      chan iface.IWsSession
	UnRegisterCh chan iface.IWsSession
}

func GetSessionMgr() *SessionMgr {
	return sessionMgr
}

func NewSessionMgr() *SessionMgr {
	s := &SessionMgr{
		RegisterCh:   make(chan iface.IWsSession, 10240),
		LoginCh:      make(chan iface.IWsSession, 10240),
		UnRegisterCh: make(chan iface.IWsSession, 10240),
	}

	return s
}

func (s *SessionMgr) AllLength() int {
	return int(atomic.LoadInt64(&s.allClientN))
}

func (s *SessionMgr) IsConn(ss iface.IWsSession) (ok bool) {
	_, ok = s.allClients.Load(ss)

	return
}

func (s *SessionMgr) GetAll() (allSS map[iface.IWsSession]struct{}) {
	allSS = make(map[iface.IWsSession]struct{})

	s.AllRange(func(ss iface.IWsSession) (result bool) {
		allSS[ss] = struct{}{}
		return true
	})

	return
}

func (s *SessionMgr) AllRange(fn func(ss iface.IWsSession) (result bool)) {
	s.allClients.Range(func(key, value any) bool {
		ss := value.(iface.IWsSession)
		result := fn(ss)
		if !result {
			return false
		}
		return true
	})
}

func (s *SessionMgr) AllAdd(ss iface.IWsSession) {
	log.Info("add session to mgr", zap.Uint64("session id", ss.GetID()))
	s.allClients.Store(ss, struct{}{})
	atomic.AddInt64(&s.allClientN, 1)
}

func (s *SessionMgr) AllDel(ss iface.IWsSession) {
	log.Info("del session to mgr", zap.Uint64("session id", ss.GetID()))
	s.allClients.Delete(ss)
	atomic.AddInt64(&s.allClientN, -1)
}

func (s *SessionMgr) OnlineLen() int {
	return int(atomic.LoadInt64(&s.onlineClientN))
}

func (s *SessionMgr) OnlineAdd(userID uint64, ss iface.IWsSession) {
	s.onlineClients.Store(userID, ss)
	atomic.AddInt64(&s.onlineClientN, 1)
}

func (s *SessionMgr) OnlineDel(userID uint64) {
	s.onlineClients.Delete(userID)
	atomic.AddInt64(&s.onlineClientN, -1)
}

func (s *SessionMgr) OnlineGetOne(userID uint64) (ss iface.IWsSession) {
	if v, ok := s.onlineClients.Load(userID); ok {
		ss = v.(iface.IWsSession)
	}

	return
}

func (s *SessionMgr) OnlineOnce(userID uint64, fn func(ss iface.IWsSession)) {
	cli := s.OnlineGetOne(userID)
	if cli == nil {
		return
	}

	fn(cli)
}

func (s *SessionMgr) OnlineRange(fn func(userID uint64, ss iface.IWsSession)) {
	s.onlineClients.Range(func(key, value any) bool {
		userID := key.(uint64)
		ss := value.(iface.IWsSession)
		fn(userID, ss)
		return true
	})

	return
}

func (s *SessionMgr) IsOnline(userID uint64) (ss iface.IWsSession, ok bool) {
	if val, ok := s.onlineClients.Load(userID); ok {
		ss = val.(iface.IWsSession)
		return ss, true
	}
	return
}

func (s *SessionMgr) Login(ss iface.IWsSession) {
	if s.IsConn(ss) {
		s.OnlineAdd(ss.GetID(), ss)
	}
}
func (s *SessionMgr) Disconnect(ss iface.IWsSession) {
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
			return session.IsHeartbeatTimeout(currentTime)
		}
		if session.CheckSomething(fn) {
			log.Info("session be cleared", zap.Uint64("session id", session.GetID()))
			session.Close()
		}
	}
}
