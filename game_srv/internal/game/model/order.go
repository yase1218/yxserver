package model

type Order struct {
	OrderId    string `json:"orderId"`
	AccountId  uint32 `json:"accountId"`
	Nick       string `json:"nick"`
	ChannelId  uint32 `json:"channelId"`
	ServerId   uint32 `json:"serverId"`
	PayType    uint32 `json:"payType"`
	Pid        uint32 `json:"pid"`
	Num        uint32 `json:"num"`
	Status     uint32 `json:"status"`
	Remark     string `json:"remark"`
	CreateTime uint32 `json:"createTime"`
	UpdateTime uint32 `json:"updateTime"`
}
