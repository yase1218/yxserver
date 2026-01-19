package model

type ShipPoker struct {
	ShipId uint32
	Poker  []uint32 // 扑克
}

type WeaponPoker struct {
	WeaponId uint32
	Poker    []uint32
}

type AccountPoker struct {
	AccountId  int64
	ShipData   []*ShipPoker
	MissData   []int
	CommData   []uint32
	WeaponData []*WeaponPoker
}

func NewAccountPoker(accountId int64, commData []uint32) *AccountPoker {
	return &AccountPoker{
		AccountId:  accountId,
		ShipData:   make([]*ShipPoker, 0, 0),
		MissData:   make([]int, 0, 0),
		CommData:   commData,
		WeaponData: make([]*WeaponPoker, 0, 0),
	}
}

func NewWeaponPoker(weaponId uint32, poker uint32) *WeaponPoker {
	ret := &WeaponPoker{
		WeaponId: weaponId,
	}
	ret.Poker = append(ret.Poker, poker)
	return ret
}

func NewShipPoker(shipId uint32) *ShipPoker {
	return &ShipPoker{
		ShipId: shipId,
	}
}
