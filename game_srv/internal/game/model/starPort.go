package model

type AccountStarPort struct {
	AccountId int64
	RDData    *RDBuild // 研发建筑
}

type RDAccessoryData struct {
	ID           int64
	Pos          uint32
	AccessoryId  uint32  // 配件模版id
	Rarity       uint32  // 稀有度
	RarityFactor float32 // 稀有度系数
	ExpFactor    float32 // 回收系数
	Level        uint32  // 等级
	LevelFactor  float32 // 等级系数
}

// RDBuild 研发建筑
type RDBuild struct {
	Level        uint32
	TechLevel    uint32
	TechExp      uint32
	PosData      []*RDAccessoryData // 配件数据
	DrawCardData *RDAccessoryData   // 抽卡数据
}

func NewRDBuild() *RDBuild {
	return &RDBuild{
		Level:     1,
		TechLevel: 1,
		TechExp:   0,
		PosData:   make([]*RDAccessoryData, 0, 0),
	}
}

func NewRDAccessory(id int64) *RDAccessoryData {
	return &RDAccessoryData{
		ID: id,
	}
}

func NewAccountStarPort(accountId int64) *AccountStarPort {
	return &AccountStarPort{
		AccountId: accountId,
		RDData:    NewRDBuild(),
	}
}
