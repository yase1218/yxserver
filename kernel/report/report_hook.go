package report

import "msg"

func ReportHookReward(channelId uint32, accountId uint64, serId, tp uint32, reward []*msg.Item) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["tp"] = tp //0： 免费  1：付费
	properties["reward"] = reward

	UploadKlLogServer(Report_Hook, properties)

}
