package oauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInspectToken(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求路径和参数
		if r.URL.Path != "/oauth/token/inspect" {
			t.Errorf("预期路径 /oauth/token/inspect, 实际获取 %s", r.URL.Path)
		}

		token := r.URL.Query().Get("token")
		if token != "test_token" {
			t.Errorf("预期token=test_token, 实际获取 %s", token)
		}

		// 返回模拟的成功响应
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

	// 创建配置和客户端
	//cfg := &config.Config{
	//	OAuth: config.OAuthConfig{
	//		TestDomain: server.URL,
	//		UseTestEnv: true,
	//	},
	//}

	client := NewClient()

	// 调用被测试的方法
	params := TokenInspectParams{
		Token: "test_token",
	}

	tokenInfo, err := client.InspectToken(params)
	if err != nil {
		t.Fatalf("InspectToken方法返回错误: %v", err)
	}

	// 验证结果
	if tokenInfo.OpenID != "test_open_id" {
		t.Errorf("预期OpenID=test_open_id, 实际获取 %s", tokenInfo.OpenID)
	}

	if tokenInfo.Platform != 1 {
		t.Errorf("预期Platform=1, 实际获取 %d", tokenInfo.Platform)
	}

	if tokenInfo.AppID != 100003 {
		t.Errorf("预期AppID=100003, 实际获取 %d", tokenInfo.AppID)
	}

	if len(tokenInfo.Scope) != 2 {
		t.Errorf("预期Scope长度=2, 实际获取 %d", len(tokenInfo.Scope))
	}
}

func TestGetUserInfo(t *testing.T) {
	// 创建一个模拟的HTTP服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求路径和参数
		if r.URL.Path != "/oauth/user/info/get" {
			t.Errorf("预期路径 /oauth/user/info/get, 实际获取 %s", r.URL.Path)
		}

		token := r.URL.Query().Get("access_token")
		if token != "test_token" {
			t.Errorf("预期access_token=test_token, 实际获取 %s", token)
		}

		// 返回模拟的成功响应
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

	// 创建配置和客户端
	//cfg := &config.Config{
	//	OAuth: config.OAuthConfig{
	//		TestDomain: server.URL,
	//		UseTestEnv: true,
	//	},
	//}

	client := NewClient()

	// 调用被测试的方法
	userInfo, err := client.GetUserInfo("test_token")
	if err != nil {
		t.Fatalf("GetUserInfo方法返回错误: %v", err)
	}

	// 验证结果
	if userInfo.Nickname != "Test User" {
		t.Errorf("预期Nickname=Test User, 实际获取 %s", userInfo.Nickname)
	}

	if userInfo.Platform != 1 {
		t.Errorf("预期Platform=1, 实际获取 %d", userInfo.Platform)
	}

	if userInfo.Icon != "http://example.com/avatar.jpg" {
		t.Errorf("预期Icon=http://example.com/avatar.jpg, 实际获取 %s", userInfo.Icon)
	}

	if userInfo.Gender != 1 {
		t.Errorf("预期Gender=1, 实际获取 %d", userInfo.Gender)
	}

	if userInfo.OpenID != "test_open_id" {
		t.Errorf("预期OpenID=test_open_id, 实际获取 %s", userInfo.OpenID)
	}
}

// 基准测试示例
func BenchmarkGetUserInfo(b *testing.B) {
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

	//cfg := &config.Config{
	//	OAuth: config.OAuthConfig{
	//		TestDomain: server.URL,
	//		UseTestEnv: true,
	//	},
	//}

	client := NewClient()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetUserInfo("test_token")
	}
}
