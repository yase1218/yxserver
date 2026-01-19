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
	"msg"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// 需要入库的好友相关操作
const (
	FriendOp_Apply uint32 = iota
	FriendOp_Accept
	FriendOp_Black

	FriendOp_Max
)

var (
	FriendOpMap map[uint32]map[uint64]map[uint64]*model.FriendOp
)

func init() {
	FriendOpMap = make(map[uint32]map[uint64]map[uint64]*model.FriendOp)
	for i := uint32(0); i < FriendOp_Max; i++ {
		FriendOpMap[i] = make(map[uint64]map[uint64]*model.FriendOp)
	}
}

func InitFiendOps(ctx context.Context) {
	cursor, err := db.LocalMongoReader.Find(
		ctx,
		config.Conf.LocalMongo.DB,
		db.CollectionName_FriendOp,
		bson.M{},
	)
	if err != nil {
		panic(fmt.Errorf("load friend op failed,%s", err.Error()))
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		op := &model.FriendOp{}
		err := cursor.Decode(op)
		if err != nil {
			panic(fmt.Errorf("load friend op decode failed,%s", err.Error()))
		}
		insert_friend_op(op)
	}
}

func insert_friend_op(op *model.FriendOp) {
	if _, ok := FriendOpMap[op.OpType][op.TarId]; !ok {
		FriendOpMap[op.OpType][op.TarId] = make(map[uint64]*model.FriendOp)
	}
	FriendOpMap[op.OpType][op.TarId][op.OpId] = op
}

func FindFriendOp(op_uid, tar_uid uint64, op_type uint32) *model.FriendOp {
	if _, ok := FriendOpMap[op_type]; !ok {
		return nil
	}
	if _, ok := FriendOpMap[op_type][tar_uid]; !ok {
		return nil
	}
	return FriendOpMap[op_type][tar_uid][op_uid]
}

func RemoveFriendOp(op_uid, tar_uid uint64, op_type uint32) {
	op := FindFriendOp(op_uid, tar_uid, op_type)
	if op == nil {
		return
	}
	delete(FriendOpMap[op_type][tar_uid], op_uid)

	save_friend_op_db(op, false)
}

func AddFriendOp(op_uid, tar_uid uint64, op_type uint32) {
	op := &model.FriendOp{
		TarId:  tar_uid,
		OpId:   op_uid,
		OpType: op_type,
		OpTime: uint64(time.Now().Unix()),
	}

	insert_friend_op(op)
	save_friend_op_db(op, true)

	switch op_type {
	case FriendOp_Apply:
		tar_player := player.FindByUserId(tar_uid)
		if tar_player != nil {
			if _, black := tar_player.UserData.FriendData.Black[op_uid]; !black {
				_, op_simple := GetPlayerSimpleInfo(op_uid)
				if op_simple != nil {
					tar_player.SendNotify(&msg.NotifyApplyFriend{
						Account: ToPlayerSimpleInfo(op_simple),
					})
				} else {
					log.Error("GetPlayerSimpleInfo nil", zap.Uint64("op id", op_uid), zap.Uint64("tar id", tar_uid))
				}
			}
		}
	}
}

func GetFriendList(pid uint32, p *player.Player) {
	//UpdateFriend(p)
	res := &msg.ResponseFriList{
		Result: msg.ErrCode_SUCC,
		Data:   make([]*msg.FriendInfo, 0),
	}
	defer p.SendResponse(pid, res, res.Result)

	if err := FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
		res.Result = err
		return
	}

	for fid := range p.UserData.FriendData.Friend {
		_, ps := GetPlayerSimpleInfo(fid)
		if ps == nil {
			log.Error("GetPlayerSimpleInfo nil", zap.Uint64("target id", fid), zap.Uint64("uid", p.GetUserId()))
			continue
		}

		res.Data = append(res.Data, &msg.FriendInfo{
			Account:        ToPlayerSimpleInfo(ps),
			LastOnlineTime: ps.LastOnlineTime,
		})
	}

	sort.Slice(res.Data, func(i, j int) bool {
		return res.Data[i].LastOnlineTime > res.Data[j].LastOnlineTime
	})
}

func GetFriApplyList(pid uint32, p *player.Player) {
	res := &msg.ResponseFriApplyList{
		Result: msg.ErrCode_SUCC,
		Data:   make([]*msg.PlayerSimpleInfo, 0),
	}
	defer p.SendResponse(pid, res, res.Result)

	u_map, op_ok := FriendOpMap[FriendOp_Apply]
	if !op_ok {
		return
	}
	tar_map, u_ok := u_map[p.GetUserId()]
	if !u_ok {
		return
	}

	for tar_id := range tar_map {
		_, ps := GetPlayerSimpleInfo(tar_id)
		if ps == nil {
			log.Error("GetPlayerSimpleInfo nil", zap.Uint64("target id", tar_id), zap.Uint64("uid", p.GetUserId()))
			continue
		}
		_, p_s := GetPlayerSimpleInfo(tar_id)
		if p_s == nil {
			log.Error("GetPlayerSimpleInfo nil", zap.Uint64("target id", tar_id), zap.Uint64("uid", p.GetUserId()))
			continue
		}
		res.Data = append(res.Data, ToPlayerSimpleInfo(p_s))
	}
}

func AddFriend(pid uint32, p *player.Player, req *msg.RequestAddFriend) {
	res := &msg.ResponseAddFriend{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	tar_id := uint64(req.AccountId)

	if p.GetUserId() == tar_id {
		res.Result = msg.ErrCode_SYSTEM_ERROR
		return
	}

	if req.AccountId == 0 {
		res.Result = msg.ErrCode_INVALID_DATA
		return
	}

	if _, ok := p.UserData.FriendData.Black[tar_id]; ok {
		res.Result = msg.ErrCode_SUCC
		return
	}

	friNum := len(p.UserData.FriendData.Friend)
	if uint32(friNum) >= template.GetSystemItemTemplate().FriMaxNum {
		res.Result = msg.ErrCode_OVER_FRIEND_LIMIT
		return
	}

	if _, already := p.UserData.FriendData.Friend[tar_id]; already {
		res.Result = msg.ErrCode_OTHER_IS_FRIEND
		return
	}

	_, p_s := GetPlayerSimpleInfo(tar_id)
	if p_s == nil {
		res.Result = msg.ErrCode_PLAYER_NOT_EXIST
		return
	}

	// 已经申请过
	if FindFriendOp(p.GetUserId(), tar_id, FriendOp_Apply) != nil {
		return
	}

	AddFriendOp(p.GetUserId(), tar_id, FriendOp_Apply)
}

func FriendApplyOp(pid uint32, p *player.Player, req *msg.RequestFriApplyOp) {
	op_uid := uint64(req.AccountId)
	res := &msg.ResponseFriApplyOp{
		Result:    msg.ErrCode_SUCC,
		AccountId: req.AccountId,
		Op:        req.Op,
	}
	defer p.SendResponse(pid, res, res.Result)

	if err := FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
		res.Result = err
		return
	}

	accept := false
	if req.Op == 0 {
		accept = true
	}

	if req.AccountId == 0 {
		process_all_friend_apply(p, accept)
	} else {
		op := FindFriendOp(op_uid, p.GetUserId(), FriendOp_Apply)
		if op == nil {
			log.Error("op nil", zap.Uint64("op id", op_uid), zap.Uint64("tar id", p.GetUserId()))
			res.Result = msg.ErrCode_SYSTEM_ERROR
			return
		}
		process_friend_apply(p, op, accept)
	}
}

func DelFriend(pid uint32, p *player.Player, req *msg.RequestDelFriend) {
	res := &msg.ResponseDelFriend{
		Result:    msg.ErrCode_SUCC,
		AccountId: req.AccountId,
	}
	defer p.SendResponse(pid, res, res.Result)

	for _, uid := range req.AccountId {
		del_friend(uint64(uid), p)
	}
}

func RecommandPlayer(pid uint32, p *player.Player) {
	res := &msg.ResponseRecommandFriend{
		Result: msg.ErrCode_SUCC,
		Data:   make([]*msg.PlayerSimpleInfo, 0),
	}
	defer p.SendResponse(pid, res, res.Result)

	players := player.AllPlayers()
	list := make([]uint64, 0, len(players))
	for _, r_p := range players {
		// 过滤自己
		if r_p.GetUserId() == p.GetUserId() {
			continue
		}

		// 没有开启好友功能的
		if msg.ErrCode_SUCC != FunctionOpen(r_p, publicconst.Friend) {
			continue
		}

		// 已经是好友
		if _, ok := p.UserData.FriendData.Friend[r_p.GetUserId()]; ok {
			continue
		}

		// 在黑名单
		if _, ok := p.UserData.FriendData.Black[r_p.GetUserId()]; ok {
			continue
		}

		list = append(list, r_p.GetUserId())
		if len(list) >= 100 {
			break
		}
	}

	template.Shuffle_FisherYates_Forward64(list, p.Rand)

	for _, id := range list {
		_, p_s := GetPlayerSimpleInfo(id)
		if p_s == nil {
			continue
		}
		res.Data = append(res.Data, ToPlayerSimpleInfo(p_s))
		if len(res.Data) >= 5 {
			break
		}
	}
}

func GetBlackList(pid uint32, p *player.Player) {
	res := &msg.ResponseFriBlackList{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	for uid := range p.UserData.FriendData.Black {
		_, p_s := GetPlayerSimpleInfo(uid)
		if p_s == nil {
			continue
		}
		res.Data = append(res.Data, ToPlayerSimpleInfo(p_s))
	}
}

func BlackOp(pid uint32, p *player.Player, req *msg.RequestBlackListOp) {
	res := &msg.ResponseBlackListOp{
		Result: msg.ErrCode_SUCC,
		Data:   make([]*msg.PlayerSimpleInfo, 0),
		Op:     req.Op,
	}
	defer p.SendResponse(pid, res, res.Result)

	if len(req.AccountId) == 0 {
		res.Result = msg.ErrCode_INVALID_DATA
		return
	}
	for _, id := range req.AccountId {
		tar_id := uint64(id)
		if req.Op == 0 { // 加入黑名单
			delete(p.UserData.FriendData.Friend, tar_id)
			p.UserData.FriendData.Black[tar_id] = struct{}{}

			// 删除对方好友
			tar_player := player.FindByUserId(tar_id)
			if tar_player != nil {
				delete(tar_player.UserData.FriendData.Friend, tar_id)
				tar_player.SaveAccountFri()
				tar_player.SendNotify(&msg.NotifyAddBlack{
					AccountId: uint32(p.GetUserId()),
				})
			} else {
				AddFriendOp(p.GetUserId(), tar_id, FriendOp_Black)
			}

		} else { // 移除黑名单
			delete(p.UserData.FriendData.Black, tar_id)
		}
		_, p_s := GetPlayerSimpleInfo(tar_id)
		if p_s != nil {
			res.Data = append(res.Data, ToPlayerSimpleInfo(p_s))
		}
	}

	p.SaveAccountFri()
}

func SearchPlayer(pid uint32, p *player.Player, req *msg.RequestSearchPlayer) {
	res := &msg.ResponseSearchPlayer{
		Result: msg.ErrCode_SUCC,
	}
	defer p.SendResponse(pid, res, res.Result)

	if err := FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
		res.Result = err
		return
	}

	content := strings.Trim(req.Content, " ")
	if len(content) > 100 {
		res.Result = msg.ErrCode_PLAYER_NOT_EXIST
		return
	}

	id := template.GetUInt64(content)
	if id > 0 {
		if id == p.GetUserId() {
			res.Result = msg.ErrCode_NOT_ADD_SELF_FRIEND
			return
		}
		if _, ok := p.UserData.FriendData.Black[id]; ok {
			res.Result = msg.ErrCode_PLAYER_NOT_EXIST
			return
		}

		_, p_s := GetPlayerSimpleInfo(id)
		if p_s == nil {
			log.Error("GetPlayerSimpleInfo nil", zap.Uint64("target id", id), zap.Uint64("uid", p.GetUserId()))
			res.Result = msg.ErrCode_PLAYER_NOT_EXIST
			return
		}

		res.Data = ToPlayerSimpleInfo(p_s)
	} else {
		if content == p.UserData.Nick {
			res.Result = msg.ErrCode_NOT_ADD_SELF_FRIEND
			return
		}
		_, p_s := GetPlayerSimpleInfoByNick(content)
		if p_s == nil {
			res.Result = msg.ErrCode_PLAYER_NOT_EXIST
			return
		}

		if _, ok := p.UserData.FriendData.Black[p_s.Uid]; ok {
			res.Result = msg.ErrCode_PLAYER_NOT_EXIST
			return
		}

		res.Data = ToPlayerSimpleInfo(p_s)
	}
}

func del_friend(tar_id uint64, p *player.Player) {
	if _, ok := p.UserData.FriendData.Friend[tar_id]; !ok {
		return
	}
	delete(p.UserData.FriendData.Friend, tar_id)
	p.SaveAccountFri()
}

func process_friend_apply(tar_player *player.Player, op *model.FriendOp, accept bool) {
	if accept { // 接受
		// 是否已经有该好友
		if _, ok := tar_player.UserData.FriendData.Friend[op.OpId]; !ok {
			tar_player.UserData.FriendData.Friend[op.OpId] = struct{}{}
			// 添加好友
			tar_player.SaveAccountFri()
		} else {
			log.Warn("repeated friend", zap.Uint64("op id", op.OpId), zap.Uint64("tar id", op.TarId))
		}
		op_player := player.FindByUserId(op.OpId)
		if op_player != nil { // 在线
			if _, ok := op_player.UserData.FriendData.Friend[op.TarId]; !ok {
				// 添加好友
				op_player.UserData.FriendData.Friend[op.TarId] = struct{}{}
				op_player.SaveAccountFri()

				_, tar_simple := GetPlayerSimpleInfo(op.TarId)
				if tar_simple == nil {
					log.Error("GetPlayerSimpleInfo nil", zap.Uint64("op id", op.OpId), zap.Uint64("tar id", op.TarId))
				} else {
					tar_player.SendNotify(&msg.NotifyFriApplyOp{
						Data: &msg.FriendInfo{
							Account:        ToPlayerSimpleInfo(tar_simple),
							LastOnlineTime: tar_simple.LastOnlineTime,
						},
						Op: 0,
					})
				}
			} else {
				log.Warn("repeated friend", zap.Uint64("op id", op.OpId), zap.Uint64("tar id", op.TarId))
			}
		} else { // 不在线加入待处理列表
			AddFriendOp(op.TarId, op.OpId, FriendOp_Accept)
		}
	} else { // 拒绝
		// 对方在线才通知
		op_player := player.FindByUserId(op.OpId)
		if op_player != nil { // 在线
			_, tar_simple := GetPlayerSimpleInfo(op.TarId)
			if tar_simple == nil {
				log.Error("GetPlayerSimpleInfo nil", zap.Uint64("op id", op.OpId), zap.Uint64("tar id", op.TarId))
			} else {
				op_player.SendNotify(&msg.NotifyFriApplyOp{
					Data: &msg.FriendInfo{
						Account:        ToPlayerSimpleInfo(tar_simple),
						LastOnlineTime: tar_simple.LastOnlineTime,
					},
					Op: 1,
				})
			}
		}
	}

	// 删除申请
	delete(FriendOpMap[op.OpType][op.TarId], op.OpId)
	save_friend_op_db(op, false)
}

func process_all_friend_apply(tar_player *player.Player, accept bool) {
	apply_map, a_ok := FriendOpMap[FriendOp_Apply]
	if !a_ok {
		return
	}

	tar_map, tar_ok := apply_map[tar_player.GetUserId()]
	if !tar_ok {
		return
	}

	for _, op := range tar_map {
		process_friend_apply(tar_player, op, accept)
	}
}

func save_friend_op_db(op *model.FriendOp, add bool) {
	db_name := config.Conf.LocalMongo.DB

	mongo_op := &dao.WriteOperation{
		Database:   db_name,
		Collection: db.CollectionName_FriendOp,
		Filter: bson.M{
			"tar_id":  op.TarId,
			"op_id":   op.OpId,
			"op_type": op.OpType,
		},
		Uuid: op.OpId,
		Tms:  time.Now().UnixMilli(),
	}

	if add {
		mongo_op.Type = dao.Upsert
		copy_value := dao.DeepCopy(op) // use deep copy
		mongo_op.Update = bson.M{"$set": copy_value}

	} else {
		mongo_op.Type = dao.Delete
	}
	db.LocalMongoWriter.AsyncWrite(mongo_op)
}

func UpdateFriend(p *player.Player) {
	op_list := make([]*model.FriendOp, 0)

	for _, tar_map := range FriendOpMap {
		op_map, ok := tar_map[p.GetUserId()]
		if !ok {
			continue
		}

		for _, op := range op_map {
			op_list = append(op_list, op)
		}
	}
	sort.Slice(op_list, func(i, j int) bool {
		return op_list[i].OpTime <= op_list[j].OpTime
	})

	for _, op := range op_list {
		_, op_s := GetPlayerSimpleInfo(op.OpId)
		if op_s == nil {
			continue
		}
		switch op.OpType {
		case FriendOp_Apply:
			// p.SendMsgNoPacket(&msg.NotifyApplyFriend{
			// 	Account: ToPlayerSimpleInfo(op_s),
			// })
			continue
		case FriendOp_Accept:
			p.UserData.FriendData.Friend[op.OpId] = struct{}{}
			p.SaveAccountFri()
			p.SendNotify(&msg.NotifyFriApplyOp{Data: &msg.FriendInfo{
				Account: ToPlayerSimpleInfo(op_s),
			}, Op: 0})
		case FriendOp_Black:
			delete(p.UserData.FriendData.Friend, op.OpId)
			p.SaveAccountFri()
			p.SendNotify(&msg.NotifyAddBlack{
				AccountId: uint32(op.OpId),
			})
		}
		RemoveFriendOp(op.OpId, op.TarId, op.OpType)
	}
}

func PrivateChat(pid uint32, p *player.Player, req *msg.RequestPrivateChat) {
	res := &msg.ResponsePrivateChat{
		Result:  msg.ErrCode_SUCC,
		Content: req.Content,
	}
	defer p.SendResponse(pid, res, res.Result)

	targetPlayer := player.FindByUserId(uint64(req.AccountId))
	if targetPlayer == nil {
		res.Result = msg.ErrCode_OTHER_NOT_ONLINE
		return
	}
	content := strings.Trim(req.Content, " ")
	res.Content = content
	para := strings.Trim(req.Para, " ")
	res.Para = para

	_, targetSimple := GetPlayerSimpleInfo(targetPlayer.GetUserId())
	res.TargetPlayer = ToPlayerSimpleInfo(targetSimple)

	if !InBlack(targetPlayer, p.GetUserId()) {
		_, p_s := GetPlayerSimpleInfo(p.GetUserId())
		notifyMsg := &msg.NotifyPrivateChat{
			Data:    ToPlayerSimpleInfo(p_s),
			Content: content,
			Para:    para,
		}
		targetPlayer.SendNotify(notifyMsg)
	}
}

func PrivateChatCheck(p *player.Player, req *msg.RequestPrivateChat) msg.ErrCode {
	if err := FunctionOpen(p, publicconst.Friend); err != msg.ErrCode_SUCC {
		return err
	}

	if p.UserData.BaseInfo.ForbiddenChat == 1 {
		return msg.ErrCode_HAS_FORBIDDEN_CHAT
	}

	if p.GetUserId() == uint64(req.AccountId) {
		return msg.ErrCode_INVALID_DATA
	}
	targetPlayer := player.FindByUserId(uint64(req.AccountId))
	if targetPlayer == nil {
		return msg.ErrCode_OTHER_NOT_ONLINE
	}
	content := strings.Trim(req.Content, " ")
	if len(content) == 0 {
		return msg.ErrCode_INVALID_DATA
	}
	if utf8.RuneCountInString(content) > 500 {
		return msg.ErrCode_CHAT_MSG_TO_LONG
	}
	//content = template.GetForbiddenTemplate().Filter(content)

	para := strings.Trim(req.Para, " ")
	if utf8.RuneCountInString(para) > 500 {
		return msg.ErrCode_CHAT_MSG_TO_LONG
	}
	return msg.ErrCode_SUCC
}

func InBlack(p *player.Player, uid uint64) bool {
	if _, ok := p.UserData.FriendData.Black[uid]; ok {
		return true
	}
	return false
}
