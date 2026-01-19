package player

// 以下全位异步
func (p *Player) SaveNick() {
	p.set_dirty_async("nick")
}

func (p *Player) SaveLevel() {
	p.set_dirty_async("level")
}

func (p *Player) SaveHeadImg() {
	p.set_dirty_async("headimage")
}

func (p *Player) SaveHeadFrame() {
	p.set_dirty_async("headframe")
}

func (p *Player) SaveTitle() {
	p.set_dirty_async("title")
}

func (p *Player) SaveBaseInfo() {
	p.set_dirty_async("baseinfo")
}

func (p *Player) SaveStageInfo() {
	p.set_dirty_async("stageinfo")
}

func (p *Player) SaveItems() {
	//log.Debug("save items", zap.Any("items", p.UserData.Items))
	p.set_dirty_async("items")
}

func (p *Player) SaveTask() {
	p.set_dirty_async("task")
}

func (p *Player) SaveMission() {
	p.set_dirty_async("mission")
}

func (p *Player) SaveShips() {
	p.set_dirty_async("ships")
}

func (p *Player) SaveTeam() {
	p.set_dirty_async("team")
}

func (p *Player) SaveEquip() {
	p.set_dirty_async("equip")
}

func (p *Player) SaveShop() {
	p.set_dirty_async("shop")
}

func (p *Player) SavePlayMethod() {
	p.set_dirty_async("playmethod")
}

func (p *Player) SaveWeapon() {
	p.set_dirty_async("weapon")
}

func (p *Player) SaveTreasure() {
	p.set_dirty_async("treasure")
}

func (p *Player) SavePoker() {
	p.set_dirty_async("poker")
}

func (p *Player) SaveAccountActivity() {
	p.set_dirty_async("accountactivity")
}

func (p *Player) SaveCardPool() {
	p.set_dirty_async("cardpool")
}

func (p *Player) SaveAppearance() {
	p.set_dirty_async("appearance")
}

func (p *Player) SavePetData() {
	p.set_dirty_async("petdata")
}

func (p *Player) SaveAccountFri() {
	p.set_dirty_async("frienddata")
}

func (p *Player) SaveFight() {
	p.set_dirty_async("fight")
}

func (p *Player) SavePeakFight() {
	p.set_dirty_async("peakfight")
}

func (p *Player) SaveContract() {
	p.set_dirty_async("contract")
}

func (p *Player) SaveDesert() {
	p.set_dirty_async("desert")
}

func (p *Player) SaveArena() {
	p.set_dirty_async("arena")
}

func (p *Player) SaveAtlas() {
	p.set_dirty_async("atlas")
}

func (p *Player) SaveLuckSale() {
	p.set_dirty_async("lucksale")
}

func (p *Player) SaveFunctionPreview() {
	p.set_dirty_async("functionpreview")
}

func (p *Player) SaveMail() {
	p.set_dirty_async("mail")
}

func (p *Player) SaveRegister() {
	p.set_dirty_async("isregister")
}

func (p *Player) SaveEquipStage() {
	p.set_dirty_async("equipstage")
}

func (p *Player) SaveResourcesPass() {
	p.set_dirty_async("resources_pass")
}

func (p *Player) SaveLikes() {
	p.set_dirty_async("likes")
}

func (p *Player) SaveWeekPass() {
	p.set_dirty_async("weekpass")
}

func (p *Player) SavePersonalized() {
	p.set_dirty_async("personalized")
}
