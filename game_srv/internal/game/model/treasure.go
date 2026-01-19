package model

type ShipTreasure struct {
	ShipId          uint32
	WarTreasure     []uint32 // 出战秘宝
	SupportTreasure []uint32 // 支援秘宝
}

type WeaponTreasure struct {
	WeaponId uint32
	Treasure []uint32
}

type AccountTreasure struct {
	AccountId  int64
	ShipData   []*ShipTreasure
	MissData   []uint32
	CommData   []uint32
	WeaponData []*WeaponTreasure
}

func NewAccountTreasure(accountId int64, commData []uint32) *AccountTreasure {
	return &AccountTreasure{
		AccountId:  accountId,
		ShipData:   make([]*ShipTreasure, 0, 0),
		MissData:   make([]uint32, 0, 0),
		CommData:   commData,
		WeaponData: make([]*WeaponTreasure, 0, 0),
	}
}

func NewWeaponTreasure(weaponId uint32, treasure []uint32) *WeaponTreasure {
	return &WeaponTreasure{
		WeaponId: weaponId,
		Treasure: treasure,
	}
}

func NewShipTreasure(shipId uint32) *ShipTreasure {
	return &ShipTreasure{
		ShipId: shipId,
	}
}
