package report

func ReportPokerPush(channelId uint32, accountId string, serId uint32, pushIds, blacks, availableIds []int, idWeight map[int]int, retId, randNum int) {
	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId

	properties["push_ids"] = pushIds           // 推送id列表
	properties["blacks"] = blacks              // 黑名单id列表
	properties["available_ids"] = availableIds // 本次推送有效列表(排除黑名单)
	properties["id_weight"] = idWeight         // 推送id权重
	properties["ret_id"] = retId               // 返回的id
	properties["rand_num"] = randNum           // 随机数

	UploadKlLogServer(Report_Poker_Push, properties)
}
