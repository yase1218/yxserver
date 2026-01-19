package config

//
//import (
//	"encoding/json"
//	"fmt"
//	"os"
//)
//
//var msdkConfig *Config
//
//// Config 包含所有API配置信息
//type Config struct {
//	// OAuth API配置
//	OAuth OAuthConfig `json:"oauth"`
//
//	// Game API配置
//	Game GameConfig `json:"game"`
//
//	// HTTP客户端配置
//	HTTP HTTPConfig `json:"http"`
//
//	// 定义配置
//	Definitions DefinitionsConfig `json:"definitions"`
//}
//
//// OAuthConfig OAuth API的配置项
//type OAuthConfig struct {
//	// AppID 应用ID
//	AppID uint32 `json:"app_id"`
//
//	// AppKey 应用密钥，用于加密签名
//	AppKey string `json:"app_key"`
//
//	// 生产环境域名，一般为 https://[appid].connect.garena.com
//	ProductionDomain string `json:"production_domain"`
//
//	// 测试环境域名，一般为 https://testconnect.garena.com
//	TestDomain string `json:"test_domain"`
//
//	// 是否使用测试环境
//	UseTestEnv bool `json:"use_test_env"`
//}
//
//// GameConfig Game API的配置项
//type GameConfig struct {
//	// GameID 游戏ID，通常与OAuth的AppID相同
//	GameID uint32 `json:"game_id"`
//
//	// 游戏密钥，一般与OAuth的AppKey相同
//	GameKey string `json:"game_key"`
//}
//
//// HTTPConfig HTTP客户端配置
//type HTTPConfig struct {
//	// 是否使用HTTPS
//	UseHTTPS bool `json:"use_https"`
//
//	// 超时时间（秒）
//	Timeout int `json:"timeout"`
//
//	// 重试次数
//	RetryCount int `json:"retry_count"`
//
//	// 自定义请求头
//	Headers map[string]string `json:"headers"`
//}
//
//// GetOAuthDomain 获取OAuth API的域名
//func (c *Config) GetOAuthDomain() string {
//	if c.OAuth.UseTestEnv {
//		return c.OAuth.TestDomain
//	}
//	return c.OAuth.ProductionDomain
//}
//
//// GetProtocol 获取使用的协议（http/https）
//func (c *Config) GetProtocol() string {
//	if c.HTTP.UseHTTPS {
//		return "https"
//	}
//	return "http"
//}
//
//// LoadFromFile 从文件加载配置
//func LoadFromFile(cfg *Config) error {
//	if msdkConfig == nil {
//		msdkConfig = cfg
//	}
//
//	return nil
//}
//
//// GetConfig 返回全局配置
//func GetConfig() *Config {
//	if msdkConfig == nil {
//		msdkConfig = DefaultConfig()
//	}
//	return msdkConfig
//}
//
//// SaveToFile 将配置保存到文件
//func (c *Config) SaveToFile(path string) error {
//	data, err := json.MarshalIndent(c, "", "  ")
//	if err != nil {
//		return fmt.Errorf("序列化配置失败: %w", err)
//	}
//
//	err = os.WriteFile(path, data, 0644)
//	if err != nil {
//		return fmt.Errorf("保存配置文件失败: %w", err)
//	}
//
//	return nil
//}
//
//// DefinitionsConfig 包含各种枚举值的定义
//type DefinitionsConfig struct {
//	// 平台定义
//	Platforms map[uint8]string `json:"platforms"`
//
//	// 性别定义
//	Genders map[uint8]string `json:"genders"`
//
//	// 客户端类型定义
//	ClientTypes map[uint8]string `json:"client_types"`
//}
//
//// DefaultConfig 创建默认配置
//func DefaultConfig() *Config {
//	return &Config{
//		OAuth: OAuthConfig{
//			AppID:            0,
//			AppKey:           "",
//			ProductionDomain: "https://%d.connect.garena.com",
//			TestDomain:       "https://testconnect.garena.com",
//			UseTestEnv:       true,
//		},
//		Game: GameConfig{
//			GameID:  0,
//			GameKey: "",
//		},
//		HTTP: HTTPConfig{
//			UseHTTPS:   true,
//			Timeout:    30,
//			RetryCount: 3,
//			Headers:    map[string]string{},
//		},
//		Definitions: DefinitionsConfig{
//			Platforms: map[uint8]string{
//				1:  "Garena",
//				2:  "BeeTalk",
//				3:  "Facebook",
//				4:  "Guest",
//				5:  "VK",
//				6:  "Line",
//				7:  "Huawei",
//				8:  "Google",
//				9:  "WeChat",
//				10: "Apple",
//				11: "Twitter",
//				13: "Email",
//				14: "PGS",
//			},
//			Genders: map[uint8]string{
//				0: "未知",
//				1: "男性",
//				2: "女性",
//			},
//			ClientTypes: map[uint8]string{
//				0: "未知",
//				1: "iOS",
//				2: "Android",
//				3: "PC",
//			},
//		},
//	}
//}
