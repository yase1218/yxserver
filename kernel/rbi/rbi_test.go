package rbi

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRbi(t *testing.T) {
	data := &PlayerRegister{
		GameSvrId:      "1",
		DtEventTime:    time.Now(),
		VGameAppid:     "",
		PlatID:         0,
		IZoneAreaID:    0,
		VOpenID:        "",
		VRoleID:        "",
		VRoleName:      "",
		VClientIP:      "",
		Region:         "",
		Country:        "",
		GarenaOpenID:   "",
		Timekey:        0,
		ClientVersion:  "",
		SystemSoftware: "",
		SystemHardware: "",
		TelecomOper:    "",
		Network:        "",
		ScreenWidth:    0,
		ScreenHight:    0,
		Density:        0,
		CpuHardware:    "",
		Memory:         0,
		GLRender:       "",
		GLVersion:      "",
		DeviceId:       "",
		GenderType:     0,
	}
	str := StructToPipeString(data)
	t.Log(str)
}

// 启动本地 UDP server 模拟接收方
func startTestUDPServer(t *testing.T) (chan string, func(), *net.UDPAddr) {
	received := make(chan string, 1)
	var wg sync.WaitGroup

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ResolveUDPAddr failed: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("ListenUDP failed: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("ReadFromUDP failed: %v", err)
			return
		}
		received <- string(buffer[:n])
	}()

	// 返回接收通道、关闭函数、以及实际监听地址
	return received, func() {
		conn.Close()
		close(received)
		wg.Wait() // 等待 goroutine 结束再返回
	}, conn.LocalAddr().(*net.UDPAddr)
}

// 发送 UDP 数据
func sendUDPData(t *testing.T, data string, targetAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		t.Fatalf("DialUDP failed: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

// ==================== 实际测试 ====================

func Test_SendUDPStruct(t *testing.T) {
	// 启动测试用 UDP Server
	receivedChan, closeFunc, serverAddr := startTestUDPServer(t)
	defer closeFunc()

	// 构造你的结构体
	data := &PlayerRegister{
		GameSvrId:      "1",
		DtEventTime:    time.Now(),
		VGameAppid:     "",
		PlatID:         0,
		IZoneAreaID:    0,
		VOpenID:        "",
		VRoleID:        "",
		VRoleName:      "",
		VClientIP:      "",
		Region:         "",
		Country:        "",
		GarenaOpenID:   "",
		Timekey:        0,
		ClientVersion:  "",
		SystemSoftware: "",
		SystemHardware: "",
		TelecomOper:    "",
		Network:        "",
		ScreenWidth:    0,
		ScreenHight:    0,
		Density:        0,
		CpuHardware:    "",
		Memory:         0,
		GLRender:       "",
		GLVersion:      "",
		DeviceId:       "",
		GenderType:     0,
	}

	// 转换为 | 分隔格式字符串
	packet := StructToPipeString(data)

	// 发送 UDP 数据包
	sendUDPData(t, packet, serverAddr)

	// 设置超时防止死锁
	select {
	case received := <-receivedChan:
		t.Logf("Received: %s", received)

		// 验证是否包含结构体名 + 正确字段数
		if !strings.HasPrefix(received, "PlayerRegister|") {
			t.Errorf("Expected prefix 'PlayerRegister|', got: %s", received)
		}

		fields := strings.Split(received, "|")
		expectedFieldCount := 28 // PlayerRegister 字段数量
		if len(fields) != expectedFieldCount {
			t.Errorf("Expected %d fields, got %d", expectedFieldCount, len(fields))
		}

	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for UDP message")
	}
}
