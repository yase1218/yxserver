package report

func ReportFriend(channelId uint32, accountId int64, serId, tp uint32, num int) {
	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["tp"] = tp //0:添加  1：删除
	properties["num"] = num

	UploadKlLogServer(Report_Friend, properties)
}
