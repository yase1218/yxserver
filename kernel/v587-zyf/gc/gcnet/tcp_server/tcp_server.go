package tcp_server

import (
	"context"
	"github.com/v587-zyf/gc/gcnet/tcp_session"
	"github.com/v587-zyf/gc/gcnet/tcp_session_mgr"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/tools"
	"net"
	"sync"
)

type TcpServer struct {
	options *TcpOption

	ctx    context.Context
	cancel context.CancelFunc

	listener net.Listener

	wg sync.WaitGroup
}

func NewTcpServer() *TcpServer {
	s := &TcpServer{
		options: NewTcpOption(),
	}

	return s
}

func (s *TcpServer) Init(ctx context.Context, option ...Option) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt(s.options)
	}

	s.listener, err = net.Listen("tcp", s.options.listenAddr)
	if err != nil {
		log.Error("net listen err", zap.Error(err))
		return
	}

	return nil
}

func (s *TcpServer) Start() {
	s.wg.Add(1)

	go tools.GoSafe("tcp_server start listen", func() {
		defer s.wg.Done()

	LOOP:
		for {
			c, err := s.listener.Accept()
			if err != nil {
				//log.Error("tcp listen err", zap.Error(err))
				break LOOP
			}
			ss := tcp_session.NewSession(context.Background(), c)
			ss.Hooks().OnMethod(s.options.method)
			ss.Start()

			tcp_session_mgr.GetSessionMgr().AllAdd(ss)
		}
	})

	s.Wait()
}

func (s *TcpServer) Stop() {
	s.listener.Close()
}

func (s *TcpServer) Wait() {
	s.wg.Wait()
}
