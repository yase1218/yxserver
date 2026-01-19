package model

import (
	"gameserver/internal/common"
	"kernel/tools"
)

type UserMail struct {
	GlobalMailId int64 // 当前最大全局邮件id
	Mails        map[int64]*Mail
}

type Mail struct {
	MailId     int64
	Title      string
	Content    string
	Items      []*SimpleItem
	Status     uint32
	CreateTime uint32
	EndTime    uint32
	MailType   uint32
	// RoleStartTime uint32
	// RoleEndTime   uint32
	GlobalId int64 // 如果是全局邮件 记录原始id
}

func (m *Mail) Clone() *Mail {
	ret := &Mail{
		MailId:     m.MailId,
		Title:      m.Title,
		Content:    m.Content,
		Status:     m.Status,
		CreateTime: m.CreateTime,
		EndTime:    m.EndTime,
		// RoleStartTime: m.RoleStartTime,
		// RoleEndTime:   m.RoleEndTime,
	}
	if len(m.Items) > 0 {
		ret.Items = make([]*SimpleItem, len(m.Items))
		copy(ret.Items, m.Items)
	}
	return ret
}

func NewMail(mailId int64, title, content string, item []*SimpleItem, endTime uint32) *Mail {
	return &Mail{
		MailId:     mailId,
		Title:      title,
		Content:    content,
		Items:      item,
		CreateTime: tools.GetCurTime(),
		EndTime:    endTime,
	}
}

type GlobalMail struct {
	MailId        int64 `bson:"mail_id"`
	Title         string
	Content       string
	Items         []*SimpleItem
	CreateTime    uint32
	EndTime       uint32
	RoleStartTime uint32
	RoleEndTime   uint32
}

func CreateGlobalMail(mailId int64, title, content string, items []*SimpleItem, endTime, roleStart, roleEnd uint32) *GlobalMail {
	return &GlobalMail{
		MailId:        mailId,
		Title:         title,
		Content:       content,
		Items:         items,
		CreateTime:    tools.GetCurTime(),
		EndTime:       endTime,
		RoleStartTime: roleStart,
		RoleEndTime:   roleEnd,
	}
}

func (m *GlobalMail) FmtToUser() *Mail {
	return &Mail{
		MailId:     common.GenSnowFlake(),
		Title:      m.Title,
		Content:    m.Content,
		Items:      m.Items,
		CreateTime: m.CreateTime,
		EndTime:    m.EndTime,
		GlobalId:   m.MailId,
	}
}
