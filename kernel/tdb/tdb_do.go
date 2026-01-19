package tdb

func GetMonsterRefresh(stageId int) []*StageMonsterRefreshCfg {
	return tdb.MonsterRefreshMap[stageId]
}
