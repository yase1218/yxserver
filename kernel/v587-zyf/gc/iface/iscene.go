package iface

type TileType int

const (
	Empty    TileType = iota // 空
	Wall                     // 墙
	Player                   // 玩家
	Treasure                 // 宝箱
	Thunder                  // 雷
	Flag                     // 旗子
)
