package msdk

import (
	"encoding/json"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"kernel/kenum"
)

type (
	GameSendUserReq struct {
		AccessToken string   `json:"access_token"`
		Platform    uint8    `json:"platform,omitempty"`
		ToFriends   []string `json:"to_friends"`
		Title       string   `json:"title,omitempty"`
		Message     string   `json:"message"`
		Image       string   `json:"image,omitempty"`
		Data        string   `json:"data,omitempty"`
	}
	GameSendUserAck struct {
		Code uint8 `json:"result"`
	}
)

func GameSendUser(req *GameSendUserReq) (*GameSendUserAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Game_Send_Friend_Url
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
	ack := new(GameSendUserAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal user_send_friend err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}

type (
	GameGuestSwapReq struct {
		AccessToken      string `json:"access_token"`
		GuestAccessToken string `json:"guest_access_token"`
	}
	GameGuestSwapAck struct {
		Result uint8 `json:"result"`
	}
)

func GameGuestSwap(req *GameGuestSwapReq) (*GameGuestSwapAck, error) {
	reqUrl := MsdkUrl + kenum.Msdk_Game_Guest_Swap_Url
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
	ack := new(GameGuestSwapAck)
	if err = json.Unmarshal(body, ack); err != nil {
		log.Error("unmarshal user_send_friend err", zap.Error(err))
		return nil, err
	}
	return ack, nil
}
