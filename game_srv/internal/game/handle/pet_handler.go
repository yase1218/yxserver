package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadPetHandle 加载宠物
func RequestLoadPetHandle(packetId uint32, args interface{}, p *player.Player) {
	totalPower, fightList, actList := service.GetPlayerPets(p)
	resp := &msg.GetPlayerPetsInfoResp{}
	resp.Lv = totalPower
	resp.FightList = fightList
	resp.ActPets = actList
	p.SendResponse(packetId, resp, msg.ErrCode_SUCC)
}

// RequestActPet 激活宠物
func RequestActPet(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.ActavitePetReq)
	code, list := service.ActTargetPet(p, req.PetId)
	resp := &msg.ActavitePetResp{}
	resp.Code = code
	resp.ActPets = list
	p.SendResponse(packetId, resp, resp.Code)
}

// RequestUpdatePetStarLv 宠物升星
func RequestUpdatePetStarLv(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.UpdatePetStarLvReq)
	code, lv := service.UpdatePetStarLv(p, req.PetId)
	resp := &msg.UpdatePetStarLvResp{}
	resp.Code = code
	resp.PetId = req.PetId
	resp.StarLv = lv
	p.SendResponse(packetId, resp, code)
}

// RequestUpdatePetLv 宠物升级
func RequestUpdatePetLv(packetId uint32, args interface{}, p *player.Player) {
	code, lv := service.UpdatePetLv(p)
	resp := &msg.UpdatePetLvResp{}
	resp.Code = code
	resp.Lv = lv
	p.SendResponse(packetId, resp, code)
}

// RequestUpdatePetsState 更新宠物编队
func RequestUpdatePetsState(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.UpdatePetStateReq)
	code, list := service.ChangePetState(p, req.Pets)
	resp := &msg.UpdatePetStateResp{}
	resp.Code = code
	resp.Pets = list
	p.SendResponse(packetId, resp, code)
}
