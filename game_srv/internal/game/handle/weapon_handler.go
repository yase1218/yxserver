package handle

import (
	"gameserver/internal/game/player"
	"gameserver/internal/game/service"
	"msg"
)

// RequestLoadWeaponHandle 加载武器
func RequestLoadWeaponHandle(packetId uint32, args interface{}, p *player.Player) {
	retMsg := &msg.ResponseLoadWeapon{Result: msg.ErrCode_SUCC}
	retMsg.LibData = service.ToProtocolWeaponLib(p.UserData.Weapon)
	retMsg.WeaponData = service.ToProtocolWeapons(p.UserData.Weapon.Weapons)
	retMsg.SecondaryWeapon = service.ToProtocolSecondaryWeapons(p.UserData.Weapon.SecondaryWeapons)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestWeaponUpgradeHandle 武器升级
func RequestWeaponUpgradeHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestWeaponUpgrade)
	err, level := service.UpgradeWeapon(p, req.WeaponId)
	retMsg := &msg.ResponseWeaponUpgrade{Result: err}
	if err == msg.ErrCode_SUCC {
		retMsg.WeaponLibLevel = p.UserData.Weapon.WeaponLibLevel
		retMsg.WeaponLibExp = p.UserData.Weapon.WeaponLibExp
		retMsg.WeaponId = req.WeaponId
		retMsg.WeaponLevel = level
	}
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestSetSecondaryWeaponHandle 请求设置副武器
func RequestSetSecondaryWeaponHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestSetSecondaryWeapon)
	retMsg := &msg.ResponseSetSecondaryWeapon{
		Pos:      req.Pos,
		WeaponId: req.WeaponId,
	}
	retMsg.Result = service.SetSecondaryWeapon(p, req.Pos, req.WeaponId)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}

// RequestActiveWeaponHandle 请求激活武器
func RequestActiveWeaponHandle(packetId uint32, args interface{}, p *player.Player) {
	req := args.(*msg.RequestActiveWeapon)
	retMsg := &msg.ResponseActiveWeapon{}
	retMsg.Result = service.ActiveWeapon(p, req.WeaponId)
	p.SendResponse(packetId, retMsg, retMsg.Result)
}
