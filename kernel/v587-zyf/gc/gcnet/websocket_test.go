package gcnet

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/v587-zyf/gc/buffer_pool"
	"github.com/v587-zyf/gc/gcnet/ws_handler"
	"github.com/v587-zyf/gc/gcnet/ws_server"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"sync"
	"testing"
	"time"
)

type MsgHandler struct{}

var _ MsgHandler

func Recv(s iface.IWsSession, data any) {
	//fmt.Println("recv---", string(data.([]byte)))
	s.SendMsg(func(args ...any) ([]byte, error) {
		return data.([]byte), nil
	})

}
func TestConcurrentWebSocketRequests(t *testing.T) {
	ctx := context.Background()

	err := buffer_pool.Init(
		context.Background(),
		buffer_pool.WithSize(100),
		buffer_pool.WithBufferSize(1024),
		buffer_pool.WithMaxSize(100),
		buffer_pool.WithAutoCleanup(true),
		buffer_pool.WithCleanupPeriod(1*time.Second),
	)
	if err != nil {
		t.Errorf("Failed to initialize buffer pool: %v", err)
		return
	}
	if err = log.Init(ctx, log.WithSerName("Test"), log.WithSkipCaller(2)); err != nil {
		panic("Log Init err" + err.Error())
	}
	if err = ws_handler.Init(ctx, ws_handler.WithName("Test"), ws_handler.WithRecvFn(Recv)); err != nil {
		t.Errorf("Failed to initialize ws_handler: %v", err)
	}

	serverAddr := ":8080"
	serverURL := "ws://localhost:8080/ws"
	numUsers := 5000
	rate := 200          // 每秒发送200个请求
	messagesPerUser := 3 // 每人发送3条消息

	ws := ws_server.NewWsServer()
	if err = ws.Init(ctx, ws_server.WithAddr(serverAddr), ws_server.WithMethod(ws_handler.Get()), ws_server.WithHttps(false)); err != nil {
		t.Errorf("ws init err:%v", err)
		return
	}

	// 启动 WebSocket 服务器
	go ws.Start()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 并发请求
	var wg sync.WaitGroup
	wg.Add(numUsers)

	startTime := time.Now()

	// 控制并发率
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	userCh := make(chan struct{}, rate)

	for i := 0; i < numUsers; i++ {
		userCh <- struct{}{}
		go func(i int) {
			defer func() {
				<-userCh
				wg.Done()
			}()

			//fmt.Println("user-", i)
			conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
			if err != nil {
				t.Errorf("用户 %d 连接失败: %v", i, err)
				return
			}
			defer conn.Close()

			for j := 0; j < messagesPerUser; j++ {
				select {
				case <-ticker.C:
					message := fmt.Sprintf("Hello-User-%d", i)
					err = conn.WriteMessage(websocket.BinaryMessage, []byte(message))
					if err != nil {
						t.Errorf("用户 %d 发送消息失败: %v", i, err)
						continue
					}

					_, response, err := conn.ReadMessage()
					if err != nil {
						t.Errorf("用户 %d 读取消息失败: %v", i, err)
						continue
					}
					fmt.Println("response:", string(response))
					//if string(response) != message {
					//	t.Errorf("用户 %d 期望接收的消息为 '%s', 实际接收 '%s'", i, message, string(response))
					//}
				default:
					time.Sleep(1 * time.Millisecond)
				}
			}
		}(i)
	}

	// 等待所有并发操作完成
	wg.Wait()

	duration := time.Since(startTime)
	t.Logf("所有用户已处理，耗时: %s\n", duration)
}
