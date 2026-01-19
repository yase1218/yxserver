package discover

import (
	"context"
	"fmt"
	"kernel/kenum"
	"kernel/tools"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

//var serverDiscover *ServerDiscover

type UpdateFn func(string, string)
type RemoveFn func(string)

type ServerDiscover struct {
	prefix  string
	options *DiscoverOption
	cli     *clientv3.Client

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// func init() {
// 	serverDiscover = &ServerDiscover{
// 		options: NewDiscoverOption(),
// 		done:    make(chan struct{}, 1),
// 	}
// }

// func GetServerDiscover() *ServerDiscover {
// 	return serverDiscover
// }

func NewServerDiscover() *ServerDiscover {
	sd := &ServerDiscover{
		options: NewDiscoverOption(),
		done:    make(chan struct{}, 1),
	}
	return sd
}

func (s *ServerDiscover) Init(ctx context.Context, opts ...any) (err error) {
	s.ctx, s.cancel = context.WithCancel(ctx)

	for _, opt := range opts {
		opt.(Option)(s.options)
	}

	etcdCfg := clientv3.Config{
		Endpoints:   s.options.endpoints,
		DialTimeout: s.options.dialTimeout,
	}
	if s.cli, err = clientv3.New(etcdCfg); err != nil {
		return
	}

	return
}

func (s *ServerDiscover) Start(prefix string, post func(string)) {
	s.prefix = prefix
	kvCli := clientv3.NewKV(s.cli)
	valAck, err := kvCli.Get(s.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		panic(fmt.Errorf("start discover failed, err:%s", err.Error()))
	}

	for _, kv := range valAck.Kvs {
		s.listenServerFn(string(kv.Key), string(kv.Value))
	}

	go tools.GoSafePost("discover watch", func() {
		s.Watch()
	}, post)
}

func (s *ServerDiscover) Watch() {
	watcher := clientv3.NewWatcher(s.cli)
	watchCh := watcher.Watch(s.ctx, kenum.SER_PREFIX, clientv3.WithPrefix())

LOOP:
	for {
		select {
		case <-s.done:
			break LOOP
		default:
			select {
			case watchAck := <-watchCh:
				for _, event := range watchAck.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						s.updateServerFn(string(event.Kv.Key), string(event.Kv.Value))
					case clientv3.EventTypeDelete:
						s.removeServerFn(string(event.Kv.Key))
					}
				}
			case <-s.done:
				break LOOP
			}
		}
	}
}

func (s *ServerDiscover) Stop() {
	s.done <- struct{}{}
}

func (s *ServerDiscover) listenServerFn(k, v string) {
	s.updateServerFn(k, v)
}

func (s *ServerDiscover) updateServerFn(k, v string) {
	if strings.Contains(k, s.prefix) {
		s.options.updateFn(k, v)
	}
}

func (s *ServerDiscover) removeServerFn(k string) {
	if strings.Contains(k, s.prefix) {
		s.options.removeFn(k)
	}
}
