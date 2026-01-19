package model

// // AccountFriend 好友
// type AccountFriend struct {
// 	AccountId uint64
// 	Friend    []uint64
// 	ApplyList []uint64
// 	BlackList []uint64
// }

// func NewAccountFriend(accountId uint64) *AccountFriend {
// 	return &AccountFriend{
// 		AccountId: accountId,
// 		Friend:    make([]uint64, 0, 0),
// 		ApplyList: make([]uint64, 0, 0),
// 		BlackList: make([]uint64, 0, 0),
// 	}
// }

// 玩家好友列表 黑名单 好友申请列表
type UserFriend struct {
	Friend map[uint64]struct{}
	Black  map[uint64]struct{}
}

// 待处理好友操作记录
type FriendOp struct {
	TarId  uint64 `bson:"tar_id"`  // 目标id
	OpId   uint64 `bson:"op_id"`   // 操作者id
	OpType uint32 `bson:"op_type"` // 操作
	OpTime uint64 `bson:"op_time"` // 操作时间
}
