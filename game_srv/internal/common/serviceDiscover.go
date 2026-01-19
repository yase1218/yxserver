package common

// import (
// 	"context"
// 	"github.com/v587-zyf/gc/log"
// 	clientv3 "go.etcd.io/etcd/client/v3"
// 	"go.uber.org/zap"
// 	"kernel/tools"
// 	"sync"
// 	"time"
// )

// type UpdateServerFunc func(key, value string)
// type RemoveServerFunc func(key string)
// type ListServerFunc func(key, value string)

// type ServiceDiscover struct {
// 	//etcd client
// 	cli *clientv3.Client

// 	// lock
// 	lock sync.Mutex

// 	ListFunc   ListServerFunc
// 	UpdateFunc UpdateServerFunc
// 	RemoveFunc RemoveServerFunc
// }

// func NewServiceDiscover(endpoint []string) *ServiceDiscover {
// 	cli, err := clientv3.New(clientv3.Config{
// 		Endpoints:   endpoint,
// 		DialTimeout: 5 * time.Second,
// 	})
// 	if err != nil {
// 		log.Error("etcd connect err", zap.Error(err))
// 		return nil
// 		//log.Errorf("etcd client connect faith...")
// 	}

// 	return &ServiceDiscover{
// 		cli:  cli,
// 		lock: sync.Mutex{},
// 	}
// }

// // DiscoverService 服务发现
// func (s *ServiceDiscover) DiscoverService(servicePrefix string) (err error) {
// 	kvCli := clientv3.NewKV(s.cli)
// 	// 先通过前缀查询所有val
// 	valResp, err := kvCli.Get(context.Background(), servicePrefix, clientv3.WithPrefix())
// 	if err != nil {
// 		log.Error("etcd err", zap.Error(err))
// 		//log.Errorf("etcd error")
// 		return err
// 	}

// 	// 将所有服务遍历进list
// 	for _, kv := range valResp.Kvs {
// 		if s.ListFunc != nil {
// 			s.ListFunc(string(kv.Key), string(kv.Value))
// 		}
// 	}

// 	// 获取etcd的watch
// 	watcher := clientv3.NewWatcher(s.cli)

// 	// 通过watch住服务的前缀 对应上文的 "/web/userService/"
// 	// 会检测该前缀下的所有变化
// 	watchChan := watcher.Watch(context.Background(), servicePrefix, clientv3.WithPrefix())

// 	// 启一个协程进行监视
// 	go tools.GoSafe("watch service discover", func() {
// 		for watchResponse := range watchChan {
// 			for _, e := range watchResponse.Events {
// 				switch e.Type {
// 				// [新增 | 修改]
// 				case clientv3.EventTypePut:
// 					if s.UpdateFunc != nil {
// 						s.UpdateFunc(string(e.Kv.Key), string(e.Kv.Value))
// 					}
// 				case clientv3.EventTypeDelete:
// 					if s.RemoveFunc != nil {
// 						s.RemoveFunc(string(e.Kv.Key))
// 					}
// 				}
// 			}
// 		}
// 	})
// 	return nil
// }
