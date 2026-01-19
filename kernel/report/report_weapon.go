package report

func ReportWeaponAdd(channelId uint32, accountId uint64, serId, weaponId, lv uint32) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["weapon_id"] = weaponId
	properties["lv"] = lv

	UploadKlLogServer(Report_weapon_add, properties)

}

func ReportWeaponLvup(channelId uint32, accountId uint64, serId, weaponId, lv uint32) {

	properties := make(map[string]interface{})
	properties["channel_id"] = channelId
	properties["server_id"] = serId
	properties["account_id"] = accountId
	properties["weapon_id"] = weaponId
	properties["lv"] = lv

	UploadKlLogServer(Report_weapon_lvup, properties)

}
