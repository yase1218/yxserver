package discover

import (
	"context"
	"kernel/tools"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type ServiceRegister struct {
	// etcd client
	cli *clientv3.Client
	// service register key
	serviceKey string
	// // service register serverInfo
	// serverInfo string
	// leaseID
	leaseID atomic.Int64
	// mutex for thread-safe access to serverInfo
	//infoMutex sync.RWMutex
	// context for cancellation
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewServiceRegister(
	endpoints []string,
	serviceKey string,
) *ServiceRegister {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("etcd connect err", zap.Error(err))
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	serviceReg := &ServiceRegister{
		cli:        cli,
		serviceKey: serviceKey,
		ctx:        ctx,
		cancelFunc: cancel,
	}

	return serviceReg
}

// Register registers the service with etcd and starts keep-alive routine
func (s *ServiceRegister) Register(ttl int64, info string, post func(string)) error {
	err := s.registerServer(ttl, info)
	if err != nil {
		log.Error("register err", zap.Error(err))
		return err
	}

	// renewal
	go tools.GoSafePost("service register ListenKeepAliveChan", func() {
		s.ListenKeepAliveChan(ttl, info)
	}, post)

	return nil
}

// UpdateServerInfo updates the server information and optionally syncs with etcd immediately
func (s *ServiceRegister) UpdateServerInfo(newInfo string) error {
	leaseID := s.leaseID.Load()

	// Check if we have a valid lease
	if leaseID == 0 {
		return s.registerServer(30, newInfo) // Default TTL if not registered
	}

	// Update the server info in etcd with the current lease
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	_, err := s.cli.Put(ctx, s.serviceKey, newInfo, clientv3.WithLease(clientv3.LeaseID(leaseID)))
	if err != nil {
		log.Error("Failed to update server info in etcd", zap.Error(err))
		return err
	}

	log.Debug("Server info updated successfully",
		zap.String("key", s.serviceKey),
		zap.String("info", newInfo))
	return nil
}

func (s *ServiceRegister) ListenKeepAliveChan(ttl int64, info string) {
	lease := clientv3.NewLease(s.cli)

	leaseID := s.leaseID.Load()

	keepAlive, err := lease.KeepAlive(s.ctx, clientv3.LeaseID(leaseID))
	if err != nil {
		log.Error("service register keepAlive err", zap.Error(err))
		return
	}

	for {
		select {
		case res, ok := <-keepAlive:
			if !ok {
				log.Warn("KeepAlive channel closed, attempting to re-register")
				err = s.registerServer(ttl, info)
				if err != nil {
					log.Error("ListenKeepAliveChan err", zap.Error(err))
					// Wait before retrying
					time.Sleep(2 * time.Second)
					continue
				}
				newLeaseID := s.leaseID.Load()
				// Recreate keepAlive channel with new lease
				keepAlive, err = lease.KeepAlive(s.ctx, clientv3.LeaseID(newLeaseID))
				if err != nil {
					log.Error("keepAlive err", zap.Error(err))
					// Wait before retrying
					time.Sleep(2 * time.Second)
				}
			} else if res == nil {
				log.Warn("KeepAlive response is nil, channel might be closing")
			}
		case <-s.ctx.Done():
			log.Info("Stopping keep-alive listener")
			return
		}
	}
}

func (s *ServiceRegister) registerServer(ttl int64, info string) error {
	grant, err := s.cli.Grant(s.ctx, ttl)
	if err != nil {
		log.Error("lease application failed", zap.Error(err))
		return err
	}
	s.leaseID.Store(int64(grant.ID))
	// 进行注册 设置服务 并绑定租约
	if _, err := s.cli.Put(s.ctx, s.serviceKey, info, clientv3.WithLease(grant.ID)); err != nil {
		log.Error("register failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *ServiceRegister) Stop() {
	s.deregister()
	s.cancelFunc()
	s.cli.Close()
}

func (s *ServiceRegister) deregister() error {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	ctx := context.Background()

	_, err := s.cli.Delete(context.Background(), s.serviceKey)
	if err != nil {
		log.Error("Failed to delete service key from etcd",
			zap.String("key", s.serviceKey),
			zap.Error(err))
		return err
	}

	leaseID := s.leaseID.Load()

	if leaseID != 0 {
		_, err = s.cli.Revoke(ctx, clientv3.LeaseID(leaseID))
		if err != nil {
			log.Warn("Failed to revoke lease, but key was deleted",
				zap.Int64("leaseID", int64(leaseID)),
				zap.Error(err))
			// 非致命错误
		}
	}

	log.Info("Service deregistered successfully", zap.String("key", s.serviceKey))
	return nil
}
