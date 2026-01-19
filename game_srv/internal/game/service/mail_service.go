package service

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/db"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/dao"
	"kernel/tools"
	"msg"
	"sort"
	"time"

	"encoding/json"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/zy/game_data/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

var (
	global_mail map[int64]*model.GlobalMail
)

func init() {
	global_mail = make(map[int64]*model.GlobalMail)
}

func InitGlobalMail(ctx context.Context) {
	cursor, err := db.LocalMongoReader.Find(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_GlobalMail,
		bson.M{},
	)
	if err != nil {
		panic(fmt.Errorf("load gloal mail failed,%s", err.Error()))
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		mail := &model.GlobalMail{}
		err := cursor.Decode(mail)
		if err != nil {
			panic(fmt.Errorf("load gloal mail decode failed,%s", err.Error()))
		}
		insert_global_mail(mail)
	}
}

func insert_global_mail(gmail *model.GlobalMail) error {
	if gmail == nil {
		return fmt.Errorf("mail nil")
	}

	if FindGlobalMail(gmail.MailId) != nil {
		return fmt.Errorf("mail exist,id:%d", gmail.MailId)
	}

	global_mail[gmail.MailId] = gmail
	return nil
}

func FindGlobalMail(global_id int64) *model.GlobalMail {
	return global_mail[global_id]
}

func AddGlobalMail(mail *model.GlobalMail) {
	if err := insert_global_mail(mail); err != nil {
		log.Error("AddGlobalMail failed", zap.Error(err))
		return
	}
	save_global_mail_db(mail, true)
}

func DelGlobalMail(global_id int64) {
	gmail := FindGlobalMail(global_id)
	if gmail == nil {
		log.Error("gmail nil")
		return
	}
	delete(global_mail, global_id)
	save_global_mail_db(gmail, false)
}

func AddUserGlobalMail(p *player.Player, gmail *model.GlobalMail) {
	if p.UserData.BaseInfo.CreateTime >= gmail.RoleStartTime {
		if gmail.RoleEndTime == 0 || p.UserData.BaseInfo.CreateTime <= gmail.RoleEndTime {
			if m := add_user_mail(p, gmail.FmtToUser()); m != nil {
				p.SaveMail()
				notifyClientMail(p, m.MailId, false)
			}
		}
	}
}

func save_global_mail_db(mail *model.GlobalMail, add bool) {
	db_name := config.Conf.LocalMongo.DB
	mongo_op := &dao.WriteOperation{
		Database:   db_name,
		Collection: db.CollectionName_GlobalMail,
		Filter: bson.M{
			"mail_id": mail.MailId,
		},
		Tms: time.Now().UnixMilli(),
	}
	if add {
		mongo_op.Type = dao.Upsert
		copy_value := dao.DeepCopy(mail) // use deep copy
		mongo_op.Update = bson.M{"$set": copy_value}
	} else {
		mongo_op.Type = dao.Delete
	}
	db.LocalMongoWriter.AsyncWrite(mongo_op)
}

func LoadUserMail(pid uint32, p *player.Player) {
	res := &msg.ResponseLoadMail{
		Result: msg.ErrCode_SUCC,
		Data:   make([]*msg.Mail, 0),
	}
	defer p.SendResponse(pid, res, res.Result)
	res.Data = ToProtocolMails(p.UserData.MailData.Mails)
}

func get_user_mail(p *player.Player, mailId int64) *model.Mail {
	return p.UserData.MailData.Mails[mailId]
}

func DelReadMail(pid uint32, p *player.Player, req *msg.RequestDelMail) {
	res := &msg.ResponseDelMail{
		Result:  msg.ErrCode_SUCC,
		MailIds: make([]int64, 0),
	}
	defer p.SendResponse(pid, res, res.Result)
	mailIds := make([]int64, 0)
	if req.MailId > 0 {
		mail := get_user_mail(p, req.MailId)
		if mail == nil {
			res.Result = msg.ErrCode_MAIL_NOT_EXIST
			return
		}

		if mail.Status == uint32(msg.MailStatus_Mail_Read) && len(mail.Items) > 0 {
			res.Result = msg.ErrCode_MAIL_NO_REWARD
			return
		}
		del_user_mails(p, []int64{req.MailId})
		res.MailIds = append(mailIds, req.MailId)
		return
	} else {
		for _, mail := range p.UserData.MailData.Mails {
			status := msg.MailStatus(mail.Status)
			switch status {
			case msg.MailStatus_Mail_Not_Read:
				continue
			case msg.MailStatus_Mail_Read:
				if len(mail.Items) > 0 {
					continue
				}
			}
			res.MailIds = append(res.MailIds, mail.MailId)
		}
		del_user_mails(p, res.MailIds)
	}
}

// ReadMail 读取邮件
func ReadMail(pid uint32, p *player.Player, req *msg.RequestReadMail) {
	res := &msg.ResponseReadMail{
		Result: msg.ErrCode_SUCC,
		MailId: req.MailId,
	}
	defer p.SendResponse(pid, res, res.Result)
	curTime := tools.GetCurTime()
	mail := get_user_mail(p, req.MailId)
	if mail == nil {
		res.Result = msg.ErrCode_MAIL_NOT_EXIST
		return
	}

	if mail.EndTime > 0 && curTime >= mail.EndTime {
		del_user_mails(p, []int64{mail.MailId})
		res.Result = msg.ErrCode_MAIL_NOT_EXIST
		return
	}

	if mail.Status > uint32(msg.MailStatus_Mail_Read) {
		res.Result = msg.ErrCode_MAIL_HAS_GET_REWARD
		return

	}
	if mail.Status != uint32(msg.MailStatus_Mail_Read) {
		mail.Status = uint32(msg.MailStatus_Mail_Read)
		p.SaveMail()
	}
}

func BatchGetMailReward(pid uint32, p *player.Player, req *msg.RequestBatchGetMailReward) {
	res := &msg.ResponseBatchGetMailReward{
		Result:  msg.ErrCode_SUCC,
		MailIds: make([]int64, 0),
		AddItem: make([]*msg.Item, 0),
	}
	defer p.SendResponse(pid, res, res.Result)
	check_over_time_mail(p, time.Now())

	ntf_ids := make([]uint32, 0)
	for _, mail := range p.UserData.MailData.Mails {
		if ec := check_mail_reward(p, mail); ec != msg.ErrCode_SUCC {
			continue
		}
		res.MailIds = append(res.MailIds, mail.MailId)
		items, item_ids := get_mail_reward(p, mail)
		res.AddItem = append(res.AddItem, items...)
		ntf_ids = append(ntf_ids, item_ids...)
		mail.Status = uint32(msg.MailStatus_Mail_Get_Reward)
	}

	if len(ntf_ids) > 0 {
		updateClientItemsChange(p.GetUserId(), ntf_ids)
		p.SaveMail()
	}
}

func get_mail_reward(p *player.Player, mail *model.Mail) ([]*msg.Item, []uint32) {
	items := make([]*msg.Item, 0)
	item_ids := make([]uint32, 0)
	for i := 0; i < len(mail.Items); i++ {
		addItems := AddItem(p.GetUserId(),
			mail.Items[i].Id, int32(mail.Items[i].Num),
			publicconst.MailAddItem, false)
		item_ids = append(item_ids, GetSimpleItemIds(addItems)...)

		items = append(items, &msg.Item{
			ItemId:  mail.Items[i].Id,
			ItemNum: int64(mail.Items[i].Num)},
		)
	}

	return items, item_ids
}

func check_mail_reward(p *player.Player, mail *model.Mail) msg.ErrCode {
	if mail == nil {
		return msg.ErrCode_MAIL_NOT_EXIST
	}

	if len(mail.Items) == 0 {
		return msg.ErrCode_MAIL_NO_REWARD
	}

	if mail.Status >= uint32(msg.MailStatus_Mail_Get_Reward) {
		return msg.ErrCode_MAIL_HAS_GET_REWARD
	}
	return msg.ErrCode_SUCC
}

func GetMailReward(pid uint32, p *player.Player, req *msg.RequestGetMailReward) {
	res := &msg.ResponseGetMailReward{
		Result:     msg.ErrCode_SUCC,
		MailId:     req.MailId,
		RewardItem: make([]*msg.Item, 0),
	}
	defer p.SendResponse(pid, res, res.Result)

	check_over_time_mail(p, time.Now())
	mail := get_user_mail(p, req.MailId)

	ec := check_mail_reward(p, mail)
	if ec != msg.ErrCode_SUCC {
		res.Result = ec
		return
	}

	items, item_ids := get_mail_reward(p, mail)

	res.RewardItem = items

	updateClientItemsChange(p.GetUserId(), item_ids)
	mail.Status = uint32(msg.MailStatus_Mail_Get_Reward)
	p.SaveMail()
}

func refresh_global_mail_onlogin(p *player.Player, now time.Time) {
	change := false
	global_mails := make([]*model.GlobalMail, 0)
	curTime := uint32(now.Unix())

	for k, v := range global_mail {
		if k > p.UserData.MailData.GlobalMailId &&
			(v.EndTime <= 0 || v.EndTime > curTime) {
			global_mails = append(global_mails, v)
		}
	}
	sort.Slice(global_mails, func(i, j int) bool {
		return global_mails[i].MailId < global_mails[j].MailId
	})
	for _, v := range global_mails {
		if add_user_mail(p, v.FmtToUser()) != nil {
			change = true
		}
	}

	if check_over_time_mail(p, now) {
		change = true
	}
	if change {
		p.SaveMail()
	}
}

func check_over_time_mail(p *player.Player, now time.Time) bool {
	var ids []int64
	curTime := uint32(now.Unix())
	for k, v := range p.UserData.MailData.Mails {
		if v.EndTime > 0 && curTime >= v.EndTime {
			ids = append(ids, k)
		}
	}

	del_user_mails(p, ids)

	return len(ids) > 0
}

func del_user_mails(p *player.Player, ids []int64) {
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		delete(p.UserData.MailData.Mails, id)
	}
	p.SaveMail()
}

func add_user_mail(p *player.Player, mail *model.Mail) *model.Mail {
	if mail == nil {
		log.Error("mail nil", ZapUser(p))
		return nil
	}

	if mail.GlobalId > 0 {
		if p.UserData.MailData.GlobalMailId >= mail.GlobalId {
			log.Error("global mail added",
				zap.Int64("new global id", mail.GlobalId),
				zap.Int64("user global id", p.UserData.MailData.GlobalMailId),
				ZapUser(p),
			)
			return nil
		}
	}

	p.UserData.MailData.Mails[mail.MailId] = mail
	if uint32(len(p.UserData.MailData.Mails)) >= template.GetSystemItemTemplate().MAX_MAIL_NUM {
		var (
			del_mail_id             int64  // 最终要删除的邮件ID
			min_create_time         uint32 = 0
			min_no_item_create_time uint32 = 0
			min_unread_create_time  uint32 = 0
		)

		// 遍历所有邮件寻找删除候选
		for _, existingMail := range p.UserData.MailData.Mails {
			// 策略1: 优先寻找已读邮件中最旧的
			if existingMail.Status == 1 {
				if min_create_time == 0 || existingMail.CreateTime < min_create_time {
					min_create_time = existingMail.CreateTime
					del_mail_id = existingMail.MailId
				}
			} else {
				// 策略2: 其次寻找无附件未读邮件中最旧的
				if len(existingMail.Items) == 0 {
					if min_no_item_create_time == 0 || existingMail.CreateTime < min_no_item_create_time {
						min_no_item_create_time = existingMail.CreateTime
						del_mail_id = existingMail.MailId
					}
				}
				// 策略3: 最后寻找任意未读邮件中最旧的
				if min_unread_create_time == 0 || existingMail.CreateTime < min_unread_create_time {
					min_unread_create_time = existingMail.CreateTime
					// 只有当没有找到无附件邮件时，才更新候选
					if len(existingMail.Items) > 0 && min_no_item_create_time == 0 {
						del_mail_id = existingMail.MailId
					}
				}
			}
		}

		if del_mail_id > 0 {
			delete(p.UserData.MailData.Mails, del_mail_id)
			if del_mail_id == mail.MailId {
				log.Warn("new mail del for full", zap.Any("new mail", mail), ZapUser(p))
				return nil
			}
		}
	}
	if mail.GlobalId > 0 {
		p.UserData.MailData.GlobalMailId = mail.GlobalId
	}

	DebugLog("add user mail", zap.Any("mail", mail), zap.Any("user mail", p.UserData.MailData.Mails), ZapUser(p))
	return mail
}

func AddSystemMail(p *player.Player, mail *model.Mail) {
	if mail == nil {
		return
	}
	if add_user_mail(p, mail) == nil {
		log.Error("add user mail failed", zap.Any("mail", mail), ZapUser(p))
		return
	}
	p.SaveMail()
	notifyClientMail(p, mail.MailId, false)
}

func AddOfflineMail(uid uint64, mail *model.Mail) {
	if mail == nil {
		return
	}

	// 不需GlobalId 个人离线都是单独发放 设置GlobalId会影响前置全局邮件收取
	mail.GlobalId = 0

	mail_json, err := json.Marshal(mail)
	if err != nil {
		log.Error("add offline mail failed", zap.Uint64("uid", uid), zap.Any("mail", mail), zap.Error(err))
		return
	}

	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()
	key := RdbOfflineMailKey(uid)

	_, err = rc.HSet(rCtx, key, mail.MailId, mail_json).Result()
	if err != nil {
		log.Error("add offline mail failed", zap.Uint64("uid", uid), zap.Any("mail", mail), zap.Error(err))
		return
	}
}

func LoadOfflineMail(p *player.Player) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()
	key := RdbOfflineMailKey(p.GetUserId())

	hash, err := rc.HGetAll(rCtx, key).Result()
	mails := make([]*model.Mail, 0, len(hash))
	if err != nil {
		log.Error("load offline mail failed", ZapUser(p), zap.Error(err))
		return
	}
	for _, v := range hash {
		mail := &model.Mail{}
		err = json.Unmarshal([]byte(v), mail)
		if err != nil {
			log.Error("load offline mail failed", ZapUser(p), zap.Any("mail json", v), zap.Error(err))
			continue
		}

		// if mail.GlobalId > 0 {
		// 	gmail := FindGlobalMail(mail.GlobalId)
		// 	if gmail == nil { // 全局已删除
		// 		continue
		// 	}
		// 	if p.UserData.BaseInfo.CreateTime >= gmail.RoleStartTime {
		// 		if gmail.RoleEndTime > 0 || p.UserData.BaseInfo.CreateTime > gmail.RoleEndTime {
		// 			continue
		// 		} else {
		// 			continue
		// 		}
		// 	} else {
		// 		continue
		// 	}
		// }
		mails = append(mails, mail)
	}

	if len(mails) > 0 {
		add_cnt := 0
		for _, mail := range mails {
			if add_user_mail(p, mail) == nil {
				continue
			}
			add_cnt++
		}
		if add_cnt > 0 {
			p.SaveMail()
		}
	}

	if len(hash) > 0 {
		_, err = rc.Del(rCtx, key).Result()
		if err != nil {
			log.Error("load offline mail del failed", ZapUser(p), zap.Error(err))
		}
	}
}

// notifyClientMail 通知客户端邮件
func notifyClientMail(p *player.Player, mailId int64, isDelete bool) {
	ntf := &msg.NotifyMail{
		Data: &msg.Mail{},
	}
	if isDelete {
		ntf.Data.MailId = mailId
		ntf.IsDelete = 1
	} else {
		if mail := get_user_mail(p, mailId); mail != nil {
			ntf.Data = ToProtocolMail(mail)
		}
	}
	p.SendNotify(ntf)
}

func ToProtocolMail(mail *model.Mail) *msg.Mail {
	var ret msg.Mail
	ret.MailId = mail.MailId
	ret.Title = mail.Title
	ret.Content = mail.Content
	ret.Status = msg.MailStatus(mail.Status)
	ret.CreateTime = mail.CreateTime
	ret.EndTime = mail.EndTime
	ret.MailType = mail.MailType
	for i := 0; i < len(mail.Items); i++ {
		ret.Attachment = append(ret.Attachment, &msg.Item{
			ItemId:  mail.Items[i].Id,
			ItemNum: int64(mail.Items[i].Num),
		})
	}
	return &ret
}

func ToProtocolMails(mails map[int64]*model.Mail) []*msg.Mail {
	var ret []*msg.Mail
	for _, v := range mails {
		ret = append(ret, ToProtocolMail(v))
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].CreateTime <= ret[j].CreateTime
	})

	return ret
}
