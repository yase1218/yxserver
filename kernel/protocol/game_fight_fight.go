package protocol

import "encoding/json"

// 玩家战斗信息 断线时把数据存数据库,断线重连取出来继续战斗
type (
	FightInfoBasic struct {
		ShipId         uint32 // 机甲
		ShipLv         uint32 // 机甲等级
		CurExp         uint32 // 当前经验
		Hp             uint32 // 当前血量
		Chip           uint32 // 筹码
		ReliveNum      uint32 // 复活次数
		KillMonsterNum uint32 // 杀怪数量
	}
	FightInfoPoker struct {
	}
	FightInfoShop struct {
		Id        uint32 // 配件id
		Status    uint32 // 状态
		BuyStatus uint32 // 购买状态
	}
	FightInfoInteractive struct {
		Id          uint32 // id
		RemainTimes uint32 // 剩余时间
		NextRefTime uint32 // 下次刷新时间
	}
	FightInfo struct {
		StageId      uint32                // 关卡id
		FightSeconds int64                 // 战斗时长(秒)
		BasicInfo    *FightInfoBasic       // 基础数据
		PokerInfo    *FightInfoPoker       // 卡牌
		ShopInfo     *FightInfoShop        // 商店
		Interactive  *FightInfoInteractive // 交互物
		RandEvents   []uint32              // 随机事件
		// todo skill
		// todo buff
		Extra []byte // 额外数据,根据不同战斗,内容不同
	}
)

func (f *FightInfo) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FightInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, f)
}
