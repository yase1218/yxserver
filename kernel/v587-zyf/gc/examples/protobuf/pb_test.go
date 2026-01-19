package protobuf

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	t.Run("StartServer", func(t *testing.T) {
		if err := StartServer(); err != nil {
			t.Errorf("Failed to start server: %v", err)
		}

		time.Sleep(1 * time.Second) // 等待服务器启动
	})

	t.Run("SayHello", func(t *testing.T) {
		message, err := SayHello("World 123 456")
		if err != nil {
			t.Errorf("Error calling SayHello: %v", err)
		}

		expected := "Hello World 123 456"
		if message != expected {
			t.Errorf("Expected %q, got %q", expected, message)
		}
		t.Logf("msg:%s", message)
	})
}
