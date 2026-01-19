package report

import (
	"fmt"
)

const (
	Report_Hook        = "kl_hook_reward"
	Report_Friend      = "kl_friend"
	Report_Login_Sign  = "kl_login_sign"
	Report_Pass_Reward = "kl_pass_reward"
	Report_weapon_add  = "kl_weapon_add"
	Report_weapon_lvup = "kl_weapon_lvup"
	Report_Lottery     = "kl_lottery"
	Report_Poker_Push  = "kl_poker_push"
)

type KlLogReq struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

var (
	url = fmt.Sprintf("https://log.jumpoyo.com/log/uploadStrategyLog")
)

type KlLogAck struct {
}

func UploadKlLogServer(eventName string, properties map[string]interface{}) {

	// jsonData, err := json.Marshal(properties)
	// if err != nil {
	// 	fmt.Println("JSON 编码失败:", err)
	// 	return
	// }

	// req := &KlLogReq{
	// 	Cmd:  eventName,
	// 	Data: string(jsonData),
	// }
	// reqBytes, err := json.Marshal(req)
	// if err != nil {
	// 	log.Error(fmt.Sprintf("marshal err:%v ", err))
	// 	return
	// }

	// go tools.GoSafe("UploadKlLogServer", func() {
	// 	resp, err := utils.PostJson(url, reqBytes)
	// 	if err != nil {
	// 		log.Error(fmt.Sprintf("http post err:%v", err))
	// 		return
	// 	}

	// 	ack := new(KlLogAck)
	// 	if err = json.Unmarshal(resp, &ack); err != nil {
	// 		log.Error(fmt.Sprintf("unmarshal err:%v", err))
	// 		return
	// 	}
	// })

	// if loginAck.Data != nil {
	// 	//r.SetToken(loginAck.Data.Token)
	// 	r.SetUserID(uint64(loginAck.Data.AccountId))
	// 	r.SetGateAddr(loginAck.Data.ServerInfo.ServerAddr)
	// }
}
