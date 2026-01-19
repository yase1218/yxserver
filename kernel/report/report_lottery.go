package report

func ReportLottery(channelId uint32, accountId int64, serId, cardId, times uint32) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["card_id"] = cardId
	properties["times"] = times

	UploadKlLogServer(Report_Lottery, properties)

}
