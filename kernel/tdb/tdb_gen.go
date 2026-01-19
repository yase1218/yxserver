package tdb

func (t *Tdb) Patch() {
	t.genMonsterRefresh()
}

func (t *Tdb) genMonsterRefresh() {
	t.MonsterRefreshMap = make(map[int][]*StageMonsterRefreshCfg)
	for _, cfg := range t.StageMonsterRefreshCfgs {
		if t.MonsterRefreshMap[cfg.StageId] == nil {
			t.MonsterRefreshMap[cfg.StageId] = make([]*StageMonsterRefreshCfg, 0)
		}
		t.MonsterRefreshMap[cfg.StageId] = append(t.MonsterRefreshMap[cfg.StageId], cfg)
	}
}
