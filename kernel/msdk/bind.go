package msdk

import (
	"encoding/json"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"kernel/kenum"
)

type (
	BindPlatformCreateReq struct {
		AppID                uint32 `json:"app_id"`                 // 与访问令牌对应的应用 ID
		AccessToken          string `json:"access_token"`           // 主用户的访问令牌
		SecondaryAccessToken string `json:"secondary_access_token"` // 辅助用户的访问令牌
	}
	BindPlatformCreateAck struct {
		Result uint8 `json:"result"` // 如果出现，则始终为 0，表示请求成功
	}
)

func BindPlatformCreate(req *BindPlatformCreateReq) (*BindPlatformCreateAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Bind_Platform_Create_Url
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
	ack := new(BindPlatformCreateAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal user_send_friend err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

type (
	BindPlatformGetInfoReq struct {
		AppID       uint32 `json:"app_id"`       // 与访问令牌对应的应用 ID
		AccessToken string `json:"access_token"` // 主用户的访问令牌
	}
	BindPlatformGetBoundedAccountsUserInfo struct {
		Nickname string `json:"nickname"`
		Gender   uint8  `json:"gender"`
		Icon     string `json:"icon"`
		Email    string `json:"email,omitempty"`
	}
	BindPlatformGetBoundedAccounts struct {
		Platform   uint8                                  `json:"platform"`
		UID        uint64                                 `json:"uid,omitempty"` // uint64 整数, 游戏不鼓励存取/使用此字段
		UserInfo   BindPlatformGetBoundedAccountsUserInfo `json:"user_info"`
		CreateTime uint32                                 `json:"create_time"`
	}
	BindPlatformGetInfoAck struct {
		AvailablePlatforms []uint8                          `json:"available_platforms"`
		BoundedAccounts    []BindPlatformGetBoundedAccounts `json:"bounded_accounts"`
	}
)

func BindPlatformGetInfo(req *BindPlatformGetInfoReq) (*BindPlatformGetInfoAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Bind_Platform_Get_Info_Url
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
	ack := new(BindPlatformGetInfoAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal user_send_friend err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

type (
	BindPlatformDelReq struct {
		AppID             uint32 `json:"app_id"`
		AccessToken       string `json:"access_token"`
		SecondaryPlatform uint8  `json:"secondary_platform"`
		SecondaryUID      uint64 `json:"secondary_uid"`
	}
	BindPlatformDelAck struct {
		Result uint8 `json:"result"`
	}
)

func BindPlatformDel(req *BindPlatformDelReq) (*BindPlatformDelAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Bind_Platform_Del_Url
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
	ack := new(BindPlatformDelAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal user_send_friend err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}
