package handle

import (
	"msg"

	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
)

// RequestLoadEquipHandle 加载装备
func RequestLoadEquipHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.ResponseLoadEquip{Result: msg.ErrCode_SUCC}
	if len(p.UserData.Equip.EquipData) > 0 {
		res.Equips = service.ToProtocolEquips(p.UserData.Equip.EquipData)
	}
	if len(p.UserData.Equip.EquipPosData) > 0 {
		res.UseEquips = service.ToProtocolEquipPosList(p.UserData.Equip.EquipPosData)
	}
	res.SuitReward = service.ToProtocolSuit(p.UserData.Equip)
	res.Suits = service.ToProtocolEquipSuits(p.UserData.Equip.EquipSuits)
	res.UseSuitId = p.UserData.Equip.UseEquipSuit
	p.SendResponse(packetId, res, res.Result)
}

// RequestEquipPosHandle 安装装备
func RequestEquipPosHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestEquipPos)
	err, _, equipPos := service.EquipPos(p, req.Pos, req.EquipId)
	res := &msg.ResponseEquipPos{Result: err}
	res.Result = err
	if err == msg.ErrCode_SUCC {
		ntf_msg := &msg.NotifyEquipSlotChange{}
		ntf_msg.Data = append(ntf_msg.Data, service.ToProtocolEquipPos(equipPos))
		p.SendNotify(ntf_msg)

	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestEquipPosUpgradeHandle 装备部位升级
func RequestEquipPosUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestEquipPosUpgrade)
	res := &msg.ResponseEquipPosUpgrade{}
	var err msg.ErrCode
	var equipPos *model.EquipPos
	if req.IsAuto == 1 {
		err, equipPos = service.AutoUpgradeEquipPos(p, req.Pos)
	} else {
		err, equipPos = service.UpgradeEquipPos(p, req.Pos)
	}
	res.Result = err
	if res.Result == msg.ErrCode_SUCC {
		ntf_msg := &msg.NotifyEquipSlotChange{}
		service.ToProtocolEquipPos(equipPos)
		ntf_msg.Data = append(ntf_msg.Data, service.ToProtocolEquipPos(equipPos))
		p.SendNotify(ntf_msg)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestEquipUpgradeStageHandle 装备升阶
func RequestEquipUpgradeStageHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestEquipUpgradeStage)
	res := &msg.ResponseEquipUpgradeStage{}
	var result msg.ErrCode
	var equipPos []*msg.UseEquip
	var changeEquip []*model.Equip
	var getEquip []*msg.Equip
	if req.IsAuto == 1 {
		service.OneKeyUpgradeEquip(p)
		return
		//result, equipPos, changeEquip, getEquip = service.ServMgr.GetEquipService().AutoUpgradeStageEquip(playerData, req.EquipId)
	} else {
		result, equipPos, changeEquip, getEquip = service.UpgradeStageEquip(p, req.EquipId, req.CostEquips)
	}
	res.Result = result
	if result == msg.ErrCode_SUCC {
		res.Result = result
		res.OldEquipId = req.EquipId
		res.GetEquip = getEquip
		res.IsAuto = req.IsAuto

		ntf_msg := &msg.NotifyEquipChange{}
		ntf_msg.Data = service.ToProtocolEquips(changeEquip)
		p.SendNotify(ntf_msg)

		service.SyncSuit(p)

		if equipPos != nil {
			pos_ntf_msg := &msg.NotifyEquipSlotChange{}
			pos_ntf_msg.Data = equipPos
			p.SendNotify(pos_ntf_msg)
		}
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestSuitRewardHandle 请求套装奖励
func RequestSuitRewardHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSuitReward)
	res := &msg.ResponseSuitReward{}
	err, items, equip := service.GetSuitPosReward(p, req.SuitId, req.SuitPos)
	res.Result = err
	if err == msg.ErrCode_SUCC {
		ntf_msg := &msg.NotifySuitPosReward{}
		ntf_msg.SuitId = req.SuitId
		ntf_msg.Data = &msg.SuitPosReward{}
		ntf_msg.Data.Pos = req.SuitPos
		ntf_msg.Data.EquipId = equip.Data.Id
		p.SendNotify(ntf_msg)

		res.SuitId = req.SuitId
		res.GetItems = service.TemplateSimpleItemToProtocolSImpleItems(items)
	}
	p.SendResponse(packetId, res, res.Result)
}

// RequestPutInSuitHandle 请求放入套装
func RequestPutInSuitHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestPutInSuit)
	res := &msg.ResponsePutInSuit{}
	res.Result = service.PutInSuit(p, req.SuitId)
	res.SuitId = req.SuitId
	p.SendResponse(packetId, res, res.Result)
}

// RequestUseSuitHandle 请求使用套装
func RequestUseSuitHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestUseSuit)
	res := &msg.ResponseUseSuit{}
	res.Result = service.UseSuit(p, req.SuitId)
	res.SuitId = req.SuitId
	p.SendResponse(packetId, res, res.Result)
}

// RequestAllEquipUpgradeHandle 所有装备部位一键升级
func RequestAllEquipUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	service.OneKeyUpLvEquip(packetId, p)
}

func LoadGemReqHandle(packetId uint32, args interface{}, p *player.Player) {
	res := &msg.LoadGemAck{Result: msg.ErrCode_SUCC}
	if len(p.UserData.Equip.GemBag) > 0 {
		res.Gems = service.ToProtocolGemmap(p.UserData.Equip.GemBag)
	}
	if len(p.UserData.Equip.GemPos) > 0 {
		res.GemPos = service.ToProtocolGemPos(p.UserData.Equip.GemPos)
	}
	p.SendResponse(packetId, res, res.Result)
}

func SocketGemReqHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.SocketGemReq)
	err, slot := service.SocketGem(p, req.Pos, req.Uuid, req.Slot)
	res := &msg.SocketGemAck{Result: err}
	if err == msg.ErrCode_SUCC {
		res.Pos = req.Pos
		res.Uuid = req.Uuid
		res.Slot = slot
	}
	p.SendResponse(packetId, res, res.Result)
}

func UnSocketGemReqHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.UnSocketGemReq)
	err := service.UnSocketGem(p, req.Pos, req.Slot)
	res := &msg.UnSocketGemAck{Result: err}
	if err == msg.ErrCode_SUCC {
		res.Pos = req.Pos
		res.Slot = req.Slot
	}
	p.SendResponse(packetId, res, res.Result)
}

func GemLockReqHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GemLockReq)
	err := service.LockGem(p, req.Uuid, req.Lock)
	res := &msg.GemLockAck{Result: err}
	if err == msg.ErrCode_SUCC {
		res.Uuid = req.Uuid
		res.Lock = req.Lock
	}
	p.SendResponse(packetId, res, res.Result)
}

func GemComposeReqHandle(packetId uint32, args interface{}, p *player.Player) {
	err, gems := service.ComposeGem(p)
	res := &msg.GemComposeAck{Result: err}
	if err == msg.ErrCode_SUCC {
		res.Gems = gems
	}
	p.SendResponse(packetId, res, res.Result)
}

func GemRefreshReqHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.GemRefreshReq)
	err, uuid := service.RefreshGem(p, req.Uuid)
	res := &msg.GemRefreshAck{Result: err}
	if err == msg.ErrCode_SUCC {
		res.NewUuid = uuid
		res.OldUuid = req.Uuid
	}

	p.SendResponse(packetId, res, res.Result)
}
