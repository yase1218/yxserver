package nats

import (
	"fmt"
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var nc *nats.Conn

func InitService(addr string) error {
	n, err := nats.Connect(
		addr,
		nats.MaxReconnects(-1),            // 无限重试
		nats.ReconnectWait(2*time.Second), // 重连等待时间
		nats.ReconnectHandler(func(*nats.Conn) {
			log.Warn("Reconnected to NATS")
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Error("Nats Disconnected", zap.Error(err))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Error("NATS connection permanently closed")
		}),
		nats.PingInterval(20*time.Second), // 发送Ping的间隔
		nats.MaxPingsOutstanding(3),
	)
	if err != nil {
		return err
	}
	nc = n
	return nil
}

func Subscribe(subject string, f nats.MsgHandler) error {
	_, err := nc.Subscribe(subject, f)
	return err
}

func Publish(subject string, message proto.Message) error {
	if message == nil {
		return fmt.Errorf("msg is nil, sub : %s", subject)
	}
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal msg failed when publish, msg type: %s, err: %s",
			reflect.TypeOf(message).String(), err.Error())
	}
	return nc.Publish(subject, data)
}

func Stop() {
	log.Info("nats close")
	nc.Close()
}
