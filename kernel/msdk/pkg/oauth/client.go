package oauth

//
//import (
//	"encoding/json"
//	"fmt"
//	"net/url"
//	"strconv"
//	"strings"
//
//	"kernel/msdk/pkg/config"
//	"kernel/msdk/pkg/http"
//)
//
//// Client OAuth API客户端
//type Client struct {
//	config     *config.Config
//	httpClient *http.Client
//}
//
//// NewClient 创建新的OAuth客户端
//func NewClient() *Client {
//	return &Client{
//		config:     config.GetConfig(),
//		httpClient: http.NewClient(config.GetConfig()),
//	}
//}
//
//// GenericResponse 一般的API响应
//type GenericResponse map[string]interface{}
//
//// IsError 检查响应是否包含错误
//func (r GenericResponse) IsError() bool {
//	_, hasError := r["error"]
//	return hasError
//}
//
//// GetError 获取错误信息
//func (r GenericResponse) GetError() string {
//	if err, ok := r["error"].(string); ok {
//		return err
//	}
//	return ""
//}
//
//// TokenInspectParams 检查Token的参数
//type TokenInspectParams struct {
//	Token      string `json:"token"`
//	Country    string `json:"country,omitempty"`
//	IncludeUID bool   `json:"include_uid,omitempty"`
//}
//
//// TokenInfo Token信息结构
//type TokenInfo struct {
//	CreateTime         uint32   `json:"create_time"`
//	UID                uint64   `json:"uid,omitempty"`
//	OpenID             string   `json:"open_id,omitempty"`
//	Platform           uint8    `json:"platform"`
//	LoginPlatform      uint8    `json:"login_platform"`
//	MainActivePlatform uint8    `json:"main_active_platform"`
//	AppID              uint32   `json:"app_id"`
//	ExpiryTime         uint32   `json:"expiry_time"`
//	Scope              []string `json:"scope"`
//}
//
//// InspectToken 检查Token有效性
//func (c *Client) InspectToken(params TokenInspectParams) (*TokenInfo, error) {
//	values := url.Values{}
//	values.Set("token", params.Token)
//
//	if params.Country != "" {
//		values.Set("country", params.Country)
//	}
//
//	if params.IncludeUID {
//		values.Set("include_uid", "1")
//	}
//
//	url := fmt.Sprintf("%s/oauth/token/inspect", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为TokenInfo结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var tokenInfo TokenInfo
//	if err := json.Unmarshal(jsonData, &tokenInfo); err != nil {
//		return nil, fmt.Errorf("反序列化Token信息失败: %w", err)
//	}
//
//	return &tokenInfo, nil
//}
//
//// UserInfo 用户信息结构
//type UserInfo struct {
//	Nickname           string `json:"nickname"`
//	UID                uint64 `json:"uid,omitempty"`
//	Platform           uint8  `json:"platform"`
//	MainActivePlatform uint8  `json:"main_active_platform"`
//	Icon               string `json:"icon"`
//	Gender             uint8  `json:"gender"`
//	Email              string `json:"email,omitempty"`
//	OpenID             string `json:"open_id,omitempty"`
//}
//
//// GetUserInfo 获取用户信息
//func (c *Client) GetUserInfo(accessToken string) (*UserInfo, error) {
//	values := url.Values{}
//	values.Set("access_token", accessToken)
//
//	url := fmt.Sprintf("%s/oauth/user/info/get", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为UserInfo结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var userInfo UserInfo
//	if err := json.Unmarshal(jsonData, &userInfo); err != nil {
//		return nil, fmt.Errorf("反序列化用户信息失败: %w", err)
//	}
//
//	return &userInfo, nil
//}
//
//// FriendsGroup 好友分组结构
//type FriendsGroup struct {
//	Platform uint8    `json:"platform"`
//	Friends  []string `json:"friends"`
//	UIDs     []uint64 `json:"uids,omitempty"`
//}
//
//// FriendsResponse 好友列表响应
//type FriendsResponse struct {
//	FriendsGroups []FriendsGroup `json:"friends_groups"`
//}
//
//// GetFriends 获取用户好友ID列表
//func (c *Client) GetFriends(accessToken string, showUID bool) (*FriendsResponse, error) {
//	values := url.Values{}
//	values.Set("access_token", accessToken)
//
//	if showUID {
//		values.Set("show_uid", "1")
//	}
//
//	url := fmt.Sprintf("%s/oauth/user/friends/get/v2", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为FriendsResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var friendsResp FriendsResponse
//	if err := json.Unmarshal(jsonData, &friendsResp); err != nil {
//		return nil, fmt.Errorf("反序列化好友信息失败: %w", err)
//	}
//
//	return &friendsResp, nil
//}
//
//// GetInAppFriends 获取用户应用好友ID列表
//func (c *Client) GetInAppFriends(accessToken string, showUID bool) (*FriendsResponse, error) {
//	values := url.Values{}
//	values.Set("access_token", accessToken)
//
//	if showUID {
//		values.Set("show_uid", "1")
//	}
//
//	url := fmt.Sprintf("%s/oauth/user/friends/inapp/get/v2", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为FriendsResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var friendsResp FriendsResponse
//	if err := json.Unmarshal(jsonData, &friendsResp); err != nil {
//		return nil, fmt.Errorf("反序列化好友信息失败: %w", err)
//	}
//
//	return &friendsResp, nil
//}
//
//// FriendInfo 好友信息结构
//type FriendInfo struct {
//	Platform uint8  `json:"platform"`
//	OpenID   string `json:"open_id,omitempty"`
//	Nickname string `json:"nickname"`
//	Icon     string `json:"icon"`
//	Gender   uint8  `json:"gender"`
//	Uid      uint64 `json:"uid,omitempty"`
//}
//
//// FriendInfoResponse 好友信息响应
//type FriendInfoResponse struct {
//	Friends []FriendInfo `json:"friends"`
//}
//
//// GetFriendsInfo 获取用户好友信息
//func (c *Client) GetFriendsInfo(accessToken string, friendIDs []string, platform uint8) (*FriendInfoResponse, error) {
//	values := url.Values{}
//	values.Set("access_token", accessToken)
//
//	// 拼接好友ID
//	if len(friendIDs) > 0 {
//		for i, id := range friendIDs {
//			values.Set(fmt.Sprintf("friend_id[%d]", i), id)
//		}
//	}
//
//	if platform > 0 {
//		values.Set("platform", strconv.Itoa(int(platform)))
//	}
//
//	url := fmt.Sprintf("%s/oauth/user/friends/info/get", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为FriendInfoResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var friendsInfoResp FriendInfoResponse
//	if err := json.Unmarshal(jsonData, &friendsInfoResp); err != nil {
//		return nil, fmt.Errorf("反序列化好友信息失败: %w", err)
//	}
//
//	return &friendsInfoResp, nil
//}
//
//// UserRole 用户角色结构
//type UserRole struct {
//	AppServerID   uint16 `json:"app_server_id"`
//	Server        string `json:"server"`
//	AppRoleID     uint8  `json:"app_role_id"`
//	ClientType    uint8  `json:"client_type"`
//	Role          string `json:"role"`
//	AppIdentifier string `json:"app_identifier"`
//}
//
//// UserRolesResponse 用户角色响应
//type UserRolesResponse map[string][]UserRole
//
//// GetUserRoles 获取用户角色
//func (c *Client) GetUserRoles(accessToken string, appID uint32) (*UserRolesResponse, error) {
//	values := url.Values{}
//	values.Set("access_token", accessToken)
//	values.Set("app_id", strconv.Itoa(int(appID)))
//
//	url := fmt.Sprintf("%s/oauth/user/roles", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return nil, fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	// 转换为UserRolesResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var rolesResp UserRolesResponse
//	if err := json.Unmarshal(jsonData, &rolesResp); err != nil {
//		return nil, fmt.Errorf("反序列化角色信息失败: %w", err)
//	}
//
//	return &rolesResp, nil
//}
//
//// AuthorizeURLParams 授权URL参数
//type AuthorizeURLParams struct {
//	ResponseType string   // 响应类型，通常为"code"
//	Scope        []string // 授权范围
//	RedirectURI  string   // 回调地址
//	State        string   // 状态参数，用于防CSRF攻击
//}
//
//// GenerateAuthorizeURL 生成授权URL
//func (c *Client) GenerateAuthorizeURL(params AuthorizeURLParams) string {
//	values := url.Values{}
//	values.Set("response_type", params.ResponseType)
//	values.Set("client_id", strconv.Itoa(int(c.config.OAuth.AppID)))
//	values.Set("redirect_uri", params.RedirectURI)
//
//	if len(params.Scope) > 0 {
//		values.Set("scope", strings.Join(params.Scope, " "))
//	}
//
//	if params.State != "" {
//		values.Set("state", params.State)
//	}
//
//	return fmt.Sprintf("%s/oauth/authorize?%s", c.config.GetOAuthDomain(), values.Encode())
//}
//
//// TokenExchangeParams 令牌交换参数
//type TokenExchangeParams struct {
//	Code        string // 授权码
//	RedirectURI string // 回调地址，须与获取授权码时一致
//	GrantType   string // 授权类型，通常为"authorization_code"
//}
//
//// TokenResponse 令牌响应
//type TokenResponse struct {
//	AccessToken  string `json:"access_token"`
//	TokenType    string `json:"token_type"`
//	ExpiresIn    int    `json:"expires_in"`
//	RefreshToken string `json:"refresh_token,omitempty"`
//	Scope        string `json:"scope,omitempty"`
//}
//
//// ExchangeToken 使用授权码交换令牌
//func (c *Client) ExchangeToken(params TokenExchangeParams) (*TokenResponse, error) {
//	formData := url.Values{}
//	formData.Set("code", params.Code)
//	formData.Set("redirect_uri", params.RedirectURI)
//	formData.Set("grant_type", params.GrantType)
//	formData.Set("client_id", strconv.Itoa(int(c.config.OAuth.AppID)))
//	formData.Set("client_secret", c.config.OAuth.AppKey)
//
//	url := fmt.Sprintf("%s/oauth/token", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var tokenResp TokenResponse
//	if err := http.ParseJSONResponse(resp, &tokenResp); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	return &tokenResp, nil
//}
//
//// RefreshTokenParams 刷新令牌参数
//type RefreshTokenParams struct {
//	RefreshToken string // 刷新令牌
//	Scope        string // 授权范围，可选
//}
//
//// RefreshToken 刷新访问令牌
//func (c *Client) RefreshToken(params RefreshTokenParams) (*TokenResponse, error) {
//	formData := url.Values{}
//	formData.Set("refresh_token", params.RefreshToken)
//	formData.Set("grant_type", "refresh_token")
//	formData.Set("client_id", strconv.Itoa(int(c.config.OAuth.AppID)))
//	formData.Set("client_secret", c.config.OAuth.AppKey)
//
//	if params.Scope != "" {
//		formData.Set("scope", params.Scope)
//	}
//
//	url := fmt.Sprintf("%s/oauth/token", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var tokenResp TokenResponse
//	if err := http.ParseJSONResponse(resp, &tokenResp); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	return &tokenResp, nil
//}
