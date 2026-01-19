package model

import (
	"kernel/tools"
)

type AccountRole struct {
	AccountId int64
	Roles     []*Role
}

// Role 驾驶员
type Role struct {
	Id         uint32
	Exp        uint32 // 经验
	FavorLevel uint32 // 好感度等级
	StarLevel  uint32 // 星级
	CreateTime uint32
	UpdateTime uint32
}

func NewRole(id uint32) *Role {
	return &Role{
		Id:         id,
		Exp:        0,
		FavorLevel: 1,
		StarLevel:  0,
		CreateTime: tools.GetCurTime(),
	}
}

func NewAccountRole(accountId int64) *AccountRole {
	ret := &AccountRole{
		AccountId: accountId,
	}
	ret.Roles = make([]*Role, 0, 0)
	return ret
}
