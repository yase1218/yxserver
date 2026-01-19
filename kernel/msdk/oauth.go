package msdk

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"io"
	"kernel/kenum"
	"net"
	"net/http"
	"net/url"
	"time"
)

type (
	OauthInspectReq struct {
		Token      string `json:"token"`       // 访问令牌 r
		Country    string `json:"country"`     // 要检查令牌的国家，最多6个字符
		IncludeUid bool   `json:"include_uid"` // 回传中是否应包含 UID
	}

	OauthInspectAck struct {
		Error              string   `json:"error,omitempty"`
		CreateTime         uint32   `json:"create_time,omitempty"`          // 令牌创建时间
		Uid                uint64   `json:"uid,omitempty"`                  // 用户 ID
		OpenId             string   `json:"open_id,omitempty"`              // 开放 ID，在 OPEN_ID 应用标志被选中时出现
		Platform           uint8    `json:"platform,omitempty"`             // 令牌平台
		LoginPlatform      uint8    `json:"login_platform,omitempty"`       // 请求令牌时玩家使用的登录平台
		MainActivePlatform uint8    `json:"main_active_platform,omitempty"` // 表示玩家当前主要的登录账号平台。如果玩家进行了账号找回 (Swap) 或使用了游客账号绑定v2 (Guest Bind v2)，该字段将指向其账号找回后或游客绑定后的平台。如果玩家从未进行过账号找回或游客账号绑定v2，则表示其主账户平台。
		AppId              uint32   `json:"app_id,omitempty"`               // 应用 ID
		ExpiryTime         uint32   `json:"expiry_time,omitempty"`          // 令牌过期时间
		Scope              []string `json:"scope,omitempty"`                // 令牌范围
		LoginType          uint8    `json:"login_type,omitempty"`           // 此参数用于识别登录来源，使 API 调用方能够根据不同的登录类型处理相应的逻辑或校验
	}
)

func NewCustomHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 连接超时
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 20 * time.Second, // TLS 握手超时
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // 不推荐跳过证书验证
			},
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 9 * time.Second, // 整个请求的最大超时
	}
}

func OauthInspect(req OauthInspectReq) (*OauthInspectAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Inspect_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}

	httpCli := NewCustomHTTPClient()

	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpCli.Do(httpReq)
	if err != nil {
		fmt.Println("msdk oauth inspect request failed err:", err)
		//log.Error("msdk oauth inspect request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		fmt.Printf("msdk oauth inspect bad response code:%d\n", resp.StatusCode)
		//log.Error("msdk oauth inspect bad response", zap.Error(err))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("read response body failed", zap.Error(err))
		return nil, err
	}

	//body, err := utils.HttpGet(reqUrl, params)
	//if err != nil {
	//	log.Error("msdk oauth inspect err", zap.Error(err))
	//	return nil, err
	//}

	ack := new(OauthInspectAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_inspect_ack err", zap.Error(err))
		return nil, err
	}

	return ack, nil
}

type (
	OauthGetUserReq struct {
		AccessToken string `json:"access_token"`
	}
	OauthGetUserAck struct {
		Nickname           string `json:"nickname"`
		UID                uint64 `json:"uid,omitempty"`
		Platform           uint8  `json:"platform"`
		MainActivePlatform uint8  `json:"main_active_platform"`
		Icon               string `json:"icon"`
		Gender             uint8  `json:"gender"`
		Email              string `json:"email,omitempty"`
		OpenID             string `json:"open_id,omitempty"`
	}
)

func OauthGetUser(req *OauthGetUserReq) (*OauthGetUserAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Get_User_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}

	body, err := utils.HttpGet(reqUrl, params)
	if err != nil {
		log.Error("msdk oauth get user err", zap.Error(err))
		return nil, err
	}

	ack := new(OauthGetUserAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_get_user_ack err", zap.Error(err))
		return nil, err
	}

	return ack, nil
}

type (
	OauthGetFriendsReq struct {
		AccessToken string `json:"access_token"`
		ShowUID     uint32 `json:"show_uid"`
	}
	OauthGetFriendsUnit struct {
		Platform uint8    `json:"platform"`
		Friends  []string `json:"friends"`
		UIDs     []uint64 `json:"uids,omitempty"`
	}
	OauthGetFriendsAck struct {
		FriendsGroups []OauthGetFriendsUnit `json:"friends_groups"`
	}
)

func OauthGetFriends(req *OauthGetFriendsReq) (*OauthGetFriendsAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Get_Friend_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}
	body, err := utils.HttpGet(reqUrl, params)
	if err != nil {
		log.Error("msdk oauth get friends err", zap.Error(err))
		return nil, err
	}
	ack := new(OauthGetFriendsAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_get_friends_ack err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

func OauthGetFriendsInapp(req *OauthGetFriendsReq) (*OauthGetFriendsAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Get_Friend_Inapp_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}
	body, err := utils.HttpGet(reqUrl, params)
	if err != nil {
		log.Error("msdk oauth get friends inapp err", zap.Error(err))
		return nil, err
	}
	ack := new(OauthGetFriendsAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_get_friends_ack err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

type (
	OauthGetFriendInfoReq struct {
		AccessToken string   `json:"access_token"`
		Platform    uint8    `json:"platform"`
		Friends     []string `json:"friends"`
	}
	OauthGetFriendInfoUnit struct {
		Platform uint8  `json:"platform"`
		OpenID   string `json:"open_id,omitempty"`
		Nickname string `json:"nickname"`
		Icon     string `json:"icon"`
		Gender   uint8  `json:"gender"`
		Uid      uint64 `json:"uid,omitempty"`
	}
	OauthGetFriendInfoAck struct {
		Friends []OauthGetFriendInfoUnit `json:"friends"`
	}
)

func OauthGetFriendInfo(req *OauthGetFriendInfoReq) (*OauthGetFriendInfoAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Get_Friend_Info_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}
	body, err := utils.HttpGet(reqUrl, params)
	if err != nil {
		log.Error("msdk oauth get friend info err", zap.Error(err))
		return nil, err
	}
	ack := new(OauthGetFriendInfoAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_get_friend_info_ack err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

type (
	OauthGetRoleReq struct {
		AccessToken string `json:"access_token"`
		AppID       uint32 `json:"app_id"`
	}
	OauthGetRoleUnit struct {
		AppServerID   uint16 `json:"app_server_id"`
		Server        string `json:"server"`
		AppRoleID     uint8  `json:"app_role_id"`
		ClientType    uint8  `json:"client_type"`
		Role          string `json:"role"`
		AppIdentifier string `json:"app_identifier"`
	}
	OauthGetRoleAck map[string][]OauthGetRoleUnit
)

func OauthGetRole(req *OauthGetRoleReq) (*OauthGetRoleAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Oauth_Get_Role_Url
	params, err := utils.StructToValuesByKey(req, "json")
	if err != nil {
		log.Error("struct to values err", zap.Error(err))
		return nil, err
	}
	body, err := utils.HttpGet(reqUrl, params)
	if err != nil {
		log.Error("msdk oauth get role err", zap.Error(err))
		return nil, err
	}
	ack := new(OauthGetRoleAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal oauth_get_role_ack err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}
