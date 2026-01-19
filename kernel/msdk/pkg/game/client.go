package game

//
//import (
//	"encoding/json"
//	"fmt"
//	"net/url"
//	"strconv"
//	"time"
//
//	"kernel/msdk/pkg/config"
//	"kernel/msdk/pkg/http"
//)
//
//// Client Game API客户端
//type Client struct {
//	config     *config.Config
//	httpClient *http.Client
//}
//
//// NewClient 创建新的Game客户端
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
//// RequestSendParams 发送请求的参数
//type RequestSendParams struct {
//	AccessToken string   `json:"access_token"`
//	Platform    uint8    `json:"platform,omitempty"`
//	ToFriends   []string `json:"to_friends"`
//	Title       string   `json:"title,omitempty"`
//	Message     string   `json:"message"`
//	Image       string   `json:"image,omitempty"`
//	Data        string   `json:"data,omitempty"`
//}
//
//// SendRequest 发送请求给好友
//func (c *Client) SendRequest(params RequestSendParams) error {
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("app_key", c.config.Game.GameKey)
//
//	if params.Platform > 0 {
//		formData.Set("platform", strconv.Itoa(int(params.Platform)))
//	}
//
//	if len(params.ToFriends) > 10 {
//		return fmt.Errorf("朋友ID数量不能超过10个")
//	}
//
//	// 拼接好友ID
//	if len(params.ToFriends) > 0 {
//		friendsStr := ""
//		for i, id := range params.ToFriends {
//			if i > 0 {
//				friendsStr += ","
//			}
//			friendsStr += id
//		}
//		formData.Set("to_friends", friendsStr)
//	} else {
//		return fmt.Errorf("朋友ID列表不能为空")
//	}
//
//	if params.Title != "" {
//		formData.Set("title", params.Title)
//	}
//
//	formData.Set("message", params.Message)
//
//	if params.Image != "" {
//		formData.Set("image", params.Image)
//	}
//
//	if params.Data != "" {
//		formData.Set("data", params.Data)
//	}
//
//	url := fmt.Sprintf("%s/game/user/request/send", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// GuestSwapParams 游客账户绑定参数
//type GuestSwapParams struct {
//	AccessToken      string `json:"access_token"`
//	GuestAccessToken string `json:"guest_access_token"`
//}
//
//// SwapGuestAccount 将游客账户绑定到平台账户
//func (c *Client) SwapGuestAccount(params GuestSwapParams) error {
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("guest_access_token", params.GuestAccessToken)
//
//	url := fmt.Sprintf("%s/game/guest/swap", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// BindPlatformParams 绑定平台账户参数
//type BindPlatformParams struct {
//	AppID                uint32 `json:"app_id"`
//	AccessToken          string `json:"access_token"`
//	SecondaryAccessToken string `json:"secondary_access_token"`
//}
//
//// BindPlatformAccount 将账户B绑定到应用C上的账户A
//func (c *Client) BindPlatformAccount(params BindPlatformParams) error {
//	formData := url.Values{}
//	formData.Set("app_id", strconv.Itoa(int(params.AppID)))
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("secondary_access_token", params.SecondaryAccessToken)
//
//	url := fmt.Sprintf("%s/bind/app/platform/create", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// BoundAccount 已绑定账户信息
//type BoundAccount struct {
//	Platform uint8  `json:"platform"`
//	UID      uint64 `json:"uid,omitempty"`
//	UserInfo struct {
//		Nickname string `json:"nickname"`
//		Gender   uint8  `json:"gender"`
//		Icon     string `json:"icon"`
//		Email    string `json:"email,omitempty"`
//	} `json:"user_info"`
//	CreateTime uint32 `json:"create_time"`
//}
//
//// BindPlatformInfoResponse 平台绑定信息响应
//type BindPlatformInfoResponse struct {
//	AvailablePlatforms []uint8        `json:"available_platforms"`
//	BoundedAccounts    []BoundAccount `json:"bounded_accounts"`
//}
//
//// GetBindPlatformInfo 获取平台绑定信息
//func (c *Client) GetBindPlatformInfo(appID uint32, accessToken string) (*BindPlatformInfoResponse, error) {
//	formData := url.Values{}
//	formData.Set("app_id", strconv.Itoa(int(appID)))
//	formData.Set("access_token", accessToken)
//
//	url := fmt.Sprintf("%s/bind/app/platform/info/get", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
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
//	// 转换为BindPlatformInfoResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var infoResp BindPlatformInfoResponse
//	if err := json.Unmarshal(jsonData, &infoResp); err != nil {
//		return nil, fmt.Errorf("反序列化绑定信息失败: %w", err)
//	}
//
//	return &infoResp, nil
//}
//
//// DeleteBindPlatformParams 删除平台绑定参数
//type DeleteBindPlatformParams struct {
//	AppID             uint32 `json:"app_id"`
//	AccessToken       string `json:"access_token"`
//	SecondaryPlatform uint8  `json:"secondary_platform"`
//	SecondaryUID      uint64 `json:"secondary_uid"`
//}
//
//// DeleteBindPlatform 删除平台绑定
//func (c *Client) DeleteBindPlatform(params DeleteBindPlatformParams) error {
//	formData := url.Values{}
//	formData.Set("app_id", strconv.Itoa(int(params.AppID)))
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("secondary_platform", strconv.Itoa(int(params.SecondaryPlatform)))
//	formData.Set("secondary_uid", strconv.FormatUint(params.SecondaryUID, 10))
//
//	url := fmt.Sprintf("%s/bind/app/platform/delete", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// LogoutAllDevicesParams 注销所有设备参数
//type LogoutAllDevicesParams struct {
//	AccessToken string `json:"access_token"`
//	AppID       uint32 `json:"app_id"`
//}
//
//// LogoutAllDevices 注销所有设备
//func (c *Client) LogoutAllDevices(params LogoutAllDevicesParams) error {
//	timestamp := time.Now().Unix()
//
//	// 构建请求体
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("app_id", strconv.Itoa(int(params.AppID)))
//	formData.Set("timestamp", strconv.FormatInt(timestamp, 10))
//
//	// 计算签名
//	message := fmt.Sprintf("access_token=%s&app_id=%d&timestamp=%d",
//		params.AccessToken, params.AppID, timestamp)
//	signature := http.GenerateHMACSHA256(c.config.Game.GameKey, message)
//
//	// 设置请求头
//	headers := map[string]string{
//		"Authorization": fmt.Sprintf("Signature %s", signature),
//	}
//
//	url := fmt.Sprintf("%s/game/logout_all_devices", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, headers, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// LocalRequirementParams 获取设置参数
//type LocalRequirementParams struct {
//	AppID  uint32 `json:"app_id"`
//	Region string `json:"region"`
//}
//
//// LocalRequirementResponse 设置响应
//type LocalRequirementResponse struct {
//	Code int `json:"code"`
//	Data struct {
//		IsFormRequired bool   `json:"is_form_required"`
//		GameMinAge     int    `json:"game_min_age"`
//		GovMinAge      int    `json:"gov_min_age"`
//		IsSkippable    bool   `json:"is_skippable"`
//		WebURL         string `json:"web_url"`
//	} `json:"data,omitempty"`
//	Error       string `json:"error,omitempty"`
//	ErrorDetail string `json:"error_detail,omitempty"`
//}
//
//// GetLocalRequirement 获取设置
//func (c *Client) GetLocalRequirement(params LocalRequirementParams) (*LocalRequirementResponse, error) {
//	values := url.Values{}
//	values.Set("app_id", strconv.Itoa(int(params.AppID)))
//	values.Set("region", params.Region)
//
//	url := fmt.Sprintf("%s/api/v1/game/local-requirement", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result LocalRequirementResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.Code != 0 {
//		return &result, fmt.Errorf("API错误: %s, 详情: %s", result.Error, result.ErrorDetail)
//	}
//
//	return &result, nil
//}
//
//// UserInfoParams 获取用户信息参数
//type UserInfoParams struct {
//	AppID           uint32 `json:"app_id"`
//	Region          string `json:"region"`
//	AccessToken     string `json:"access_token,omitempty"`
//	GpcSessionToken string `json:"gpc_session_token,omitempty"`
//}
//
//// UserInfoResponse 用户信息响应
//type UserInfoResponse struct {
//	Code int `json:"code"`
//	Data struct {
//		HasUserSkipped    bool   `json:"has_user_skipped"`
//		HasUserIdentified bool   `json:"has_user_identified"`
//		IsGuardian        bool   `json:"is_guardian"`
//		Name              string `json:"name,omitempty"`
//		Mobile            string `json:"mobile,omitempty"`
//		Dob               string `json:"dob,omitempty"`
//		PrefillMobile     string `json:"prefill_mobile,omitempty"`
//		Platform          uint8  `json:"platform"`
//	} `json:"data,omitempty"`
//	Error       string `json:"error,omitempty"`
//	ErrorDetail string `json:"error_detail,omitempty"`
//}
//
//// GetUserLocalInfo 获取用户本地信息
//func (c *Client) GetUserLocalInfo(params UserInfoParams) (*UserInfoResponse, error) {
//	values := url.Values{}
//	values.Set("app_id", strconv.Itoa(int(params.AppID)))
//	values.Set("region", params.Region)
//
//	if params.AccessToken != "" {
//		values.Set("access_token", params.AccessToken)
//	} else if params.GpcSessionToken != "" {
//		values.Set("gpc_session_token", params.GpcSessionToken)
//	} else {
//		return nil, fmt.Errorf("AccessToken和GpcSessionToken至少需要提供一个")
//	}
//
//	url := fmt.Sprintf("%s/api/v1/game/local-requirement/user-info", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Get(url, values, nil)
//	if err != nil {
//		return nil, fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result UserInfoResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return nil, fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.Code != 0 {
//		return &result, fmt.Errorf("API错误: %s, 详情: %s", result.Error, result.ErrorDetail)
//	}
//
//	return &result, nil
//}
//
//// ServerListResponse 服务器列表响应
//type ServerListResponse struct {
//	Servers []Server `json:"servers"`
//}
//
//// Server 服务器信息
//type Server struct {
//	ServerID   uint16 `json:"server_id"`
//	ServerName string `json:"server_name"`
//	Status     uint8  `json:"status"` // 0: 离线, 1: 在线, 2: 维护
//	CreateTime uint32 `json:"create_time"`
//}
//
//// GetServerList 获取服务器列表
//func (c *Client) GetServerList(appID uint32) (*ServerListResponse, error) {
//	values := url.Values{}
//	values.Set("app_id", strconv.Itoa(int(appID)))
//
//	url := fmt.Sprintf("%s/game/servers/list", c.config.GetOAuthDomain())
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
//	// 转换为ServerListResponse结构
//	jsonData, err := json.Marshal(result)
//	if err != nil {
//		return nil, fmt.Errorf("序列化响应失败: %w", err)
//	}
//
//	var serverList ServerListResponse
//	if err := json.Unmarshal(jsonData, &serverList); err != nil {
//		return nil, fmt.Errorf("反序列化服务器列表失败: %w", err)
//	}
//
//	return &serverList, nil
//}
//
//// RoleParams 角色参数
//type RoleParams struct {
//	AccessToken   string `json:"access_token"`
//	AppServerID   uint16 `json:"app_server_id"`
//	AppRoleID     uint8  `json:"app_role_id"`
//	Role          string `json:"role"`
//	ClientType    uint8  `json:"client_type,omitempty"`
//	AppIdentifier string `json:"app_identifier,omitempty"`
//}
//
//// CreateRole 创建角色
//func (c *Client) CreateRole(params RoleParams) error {
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("app_server_id", strconv.Itoa(int(params.AppServerID)))
//	formData.Set("app_role_id", strconv.Itoa(int(params.AppRoleID)))
//	formData.Set("role", params.Role)
//
//	if params.ClientType > 0 {
//		formData.Set("client_type", strconv.Itoa(int(params.ClientType)))
//	}
//
//	if params.AppIdentifier != "" {
//		formData.Set("app_identifier", params.AppIdentifier)
//	}
//
//	url := fmt.Sprintf("%s/game/role/create", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// UpdateRole 更新角色
//func (c *Client) UpdateRole(params RoleParams) error {
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("app_server_id", strconv.Itoa(int(params.AppServerID)))
//	formData.Set("app_role_id", strconv.Itoa(int(params.AppRoleID)))
//	formData.Set("role", params.Role)
//
//	if params.ClientType > 0 {
//		formData.Set("client_type", strconv.Itoa(int(params.ClientType)))
//	}
//
//	if params.AppIdentifier != "" {
//		formData.Set("app_identifier", params.AppIdentifier)
//	}
//
//	url := fmt.Sprintf("%s/game/role/update", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// DeleteRole 删除角色
//func (c *Client) DeleteRole(params RoleParams) error {
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("app_server_id", strconv.Itoa(int(params.AppServerID)))
//	formData.Set("app_role_id", strconv.Itoa(int(params.AppRoleID)))
//
//	if params.ClientType > 0 {
//		formData.Set("client_type", strconv.Itoa(int(params.ClientType)))
//	}
//
//	url := fmt.Sprintf("%s/game/role/delete", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
//
//// BatchRoleParams 批量角色操作参数
//type BatchRoleParams struct {
//	AccessToken string       `json:"access_token"`
//	Roles       []RoleParams `json:"roles"`
//	Operation   string       `json:"operation"` // "create", "update", "delete"
//}
//
//// BatchProcessRoles 批量处理角色
//func (c *Client) BatchProcessRoles(params BatchRoleParams) error {
//	// 转换角色数组为请求参数
//	formData := url.Values{}
//	formData.Set("access_token", params.AccessToken)
//	formData.Set("operation", params.Operation)
//
//	// 批量角色数据JSON编码
//	rolesData, err := json.Marshal(params.Roles)
//	if err != nil {
//		return fmt.Errorf("角色数据序列化失败: %w", err)
//	}
//	formData.Set("roles_data", string(rolesData))
//
//	url := fmt.Sprintf("%s/game/role/batch", c.config.GetOAuthDomain())
//	resp, err := c.httpClient.Post(url, formData, nil, true)
//	if err != nil {
//		return fmt.Errorf("发送请求失败: %w", err)
//	}
//
//	var result GenericResponse
//	if err := http.ParseJSONResponse(resp, &result); err != nil {
//		return fmt.Errorf("解析响应失败: %w", err)
//	}
//
//	if result.IsError() {
//		return fmt.Errorf("API错误: %s", result.GetError())
//	}
//
//	return nil
//}
