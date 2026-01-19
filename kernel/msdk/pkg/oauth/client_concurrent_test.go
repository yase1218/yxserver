package oauth

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestConcurrentRequests 测试客户端并发发送请求的能力
func TestConcurrentRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"nickname": "Test User",
			"platform": 1,
			"main_active_platform": 1,
			"icon": "http://example.com/avatar.jpg",
			"gender": 1,
			"open_id": "test_open_id"
		}`))
	}))
	defer server.Close()

	client := NewClient()

	// 并发请求数
	concurrency := 100

	// 创建等待组
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// 错误通道
	errChan := make(chan error, concurrency)

	// 并发发送请求
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			_, err := client.GetUserInfo("test_token")
			if err != nil {
				errChan <- err
			}
		}()
	}

	// 等待所有请求完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		t.Errorf("并发请求错误: %v", err)
	}
}

// TestConcurrentTokenInspect 测试并发检查Token
func TestConcurrentTokenInspect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"create_time": 1397559010,
			"expiry_time": 1398855010,
			"scope": ["get_user_info", "get_friends"],
			"app_id": 100003,
			"platform": 1,
			"login_platform": 1,
			"main_active_platform": 1,
			"open_id": "test_open_id"
		}`))
	}))
	defer server.Close()

	client := NewClient()

	// 并发请求数
	concurrency := 50

	// 创建等待组
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// 结果通道
	resultChan := make(chan *TokenInfo, concurrency)
	errChan := make(chan error, concurrency)

	// 并发发送请求
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			params := TokenInspectParams{
				Token: "test_token",
			}

			tokenInfo, err := client.InspectToken(params)
			if err != nil {
				errChan <- err
				return
			}

			resultChan <- tokenInfo
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()
	close(resultChan)
	close(errChan)

	// 检查错误
	for err := range errChan {
		t.Errorf("并发检查Token错误: %v", err)
	}

	// 验证结果
	resultCount := 0
	for tokenInfo := range resultChan {
		resultCount++

		if tokenInfo.AppID != 100003 {
			t.Errorf("预期AppID=100003, 实际获取 %d", tokenInfo.AppID)
		}
	}

	if resultCount != concurrency {
		t.Errorf("预期结果数量=%d, 实际获取 %d", concurrency, resultCount)
	}
}
