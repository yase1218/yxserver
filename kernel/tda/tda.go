package tda

import (
	"kernel/kenum"

	"github.com/ThinkingDataAnalytics/go-sdk/thinkingdata"
)

type Tda struct {
	ta *thinkingdata.TDAnalytics
}

var tda *Tda

func init() {
	tda = new(Tda)
}

func GetTda() *Tda {
	return tda
}

func GetTa() *thinkingdata.TDAnalytics {
	return tda.ta
}

func Init() error {
	// todo 换consumer 看要不要改成缓冲池形式上报
	//{
	//// 创建 BatchConsumer, 指定接收端地址、APP ID
	//consumer, err := thinkingdata.NewBatchConsumer("SERVER_URL", "APP_ID")
	//// 创建 BatchConsumer, 设置数据不压缩，默认gzip压缩，可在内网传输
	//consumer, err := thinkingdata.NewBatchConsumerWithCompress("SERVER_URL", "APP_ID",false)
	//}

	consumer, err := thinkingdata.NewDebugConsumer(kenum.Tda_Server_Url, kenum.Tda_App_Id)
	if err != nil {
		return err
	}

	ta := thinkingdata.New(consumer)
	tda.ta = &ta

	return nil
}

func Send() bool {
	return Send_Switch
}

func (t *Tda) SetUser(accountId, distinctId string, props map[string]any) error {
	return t.ta.UserSet(accountId, distinctId, props)
}
