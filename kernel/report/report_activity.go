package report

import "msg"

func ReportLoginSignReward(channelId uint32, accountId uint64, serId, actId uint32, reward []*msg.Item) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["activityId"] = actId
	properties["reward"] = reward

	UploadKlLogServer(Report_Login_Sign, properties)

}

func ReportPassReward(channelId uint32, accountId uint64, serId, actId uint32, reward []*msg.Item) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["activityId"] = actId
	properties["reward"] = reward

	UploadKlLogServer(Report_Pass_Reward, properties)

}
