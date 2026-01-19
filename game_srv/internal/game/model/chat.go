package model

import (
	"time"
)

type ChatChannel struct {
	ID         string    `bson:"_id"`         // 频道ID（serverid-type-sponsorID）
	Name       string    `bson:"name"`        // 展示名称
	OwnerID    int64     `bson:"owner_id"`    // 主管AccountID
	Admins     []int64   `bson:"admins"`      // 管理员列表(AccountID)
	CreatedAt  time.Time `bson:"created_at"`  // 创建时间
	LastActive time.Time `bson:"last_active"` // 最后活跃时间
}

type ChatMessage struct {
	ChannelID string    `bson:"channel_id"` // 频道ID
	SenderID  int64     `bson:"sender_id"`  // 发送者AccountID
	Content   string    `bson:"content"`    // 消息内容
	Items     []uint32  `bson:"items"`      // 道具列表
	Timestamp time.Time `bson:"timestamp"`  // 发送时间
}
