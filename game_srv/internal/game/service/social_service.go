package service

import (
	"context"
	"fmt"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"msg"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zeromicro/go-zero/core/collection"
	"go.uber.org/zap"
)

// SocialService 社交服务
var (
	playerCache       *collection.Cache
	playerNickCache   *collection.Cache
	playerDetailCache *collection.Cache
)

func init() {
	var err error
	// 缓存5分钟
	playerCache, err = collection.NewCache(5 * time.Minute)
	if err != nil {
		panic(fmt.Errorf("init playerCache err:%s", err.Error()))
	}
	playerNickCache, err = collection.NewCache(5 * time.Minute)
	if err != nil {
		panic(fmt.Errorf("init playerNickCache err:%s", err.Error()))
	}
	playerDetailCache, err = collection.NewCache(5 * time.Minute)
	if err != nil {
		panic(fmt.Errorf("init playerNickCache err:%s", err.Error()))
	}
}

// PlayerDetailInfo 玩家详细信息
type PlayerDetailInfo struct {
	AccountId    uint64
	Name         string
	Level        uint32
	MissionId    int
	Combat       uint32
	ShipId       uint32
	SupportInfos []*model.Ship
	Equips       []*model.EquipPos
	Pets         map[uint32]*model.Pet
	Head         uint32
	HeadFrame    uint32
	Title        uint32
	AllianceId   uint32
	AllianceName string
	ServerId     uint32
	ServerName   string
	CoatId       int
}

// PlayerSimpleInfo 玩家简单信息
type PlayerSimpleInfo struct {
	Uid            uint64
	Name           string
	Level          uint32
	MissionId      int
	Head           uint32
	HeadFrame      uint32
	Title          uint32
	ShipId         uint32
	UsePet         uint32
	LastOnlineTime uint32
	Combat         uint32
	AllianceId     uint32
	AllianceName   string
	ServerId       uint32
	ServerName     string
	CoatId         int
}

func GetPlayerDetailInfo(uid uint64) (msg.ErrCode, *PlayerDetailInfo) {
	ret, err := playerDetailCache.Take(fmt.Sprintf("%v", uid), func() (any, error) {
		var u *model.UserData
		p := player.FindByUserId(uid)
		if p != nil {
			u = p.UserData
		} else {
			var er error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			u, er = player.LoadUser(ctx, uid)
			if er != nil {
				return nil, er
			}
		}
		if u == nil {
			return nil, fmt.Errorf("user:%v not exist", uid)
		}
		supportInfos := make([]*model.Ship, 0, len(u.BaseInfo.SupportId))
		for _, v := range u.Ships.Ships {
			for _, vv := range u.BaseInfo.SupportId {
				if vv == v.Id {
					supportInfos = append(supportInfos, v)
				}
			}
		}
		data := &PlayerDetailInfo{
			AccountId:    u.UserId,
			Name:         u.Nick,
			Level:        u.Level,
			MissionId:    int(u.StageInfo.MissionId),
			Combat:       u.BaseInfo.Combat,
			ShipId:       u.BaseInfo.ShipId,
			SupportInfos: supportInfos,
			Equips:       u.Equip.EquipPosData,
			Pets:         u.PetData.Pets,
			Head:         u.HeadImg,
			HeadFrame:    u.HeadFrame,
			Title:        u.Title,
			ServerId:     u.ServerId,
			ServerName:   u.ServerName,
			CoatId:       u.Ships.GetShipCoatId(u.BaseInfo.ShipId),
		}

		//  TODO: 公会联盟信息
		// allianceInfo := ServMgr.GetAllianceService().GetAllianceInfo(int64(accountId))
		// if allianceInfo != nil {
		// 	data.AllianceId = allianceInfo.ID
		// 	data.AllianceName = allianceInfo.Name
		// }

		return data, nil
	})
	if err != nil {
		log.Error("GetPlayerDetailInfo OnInit err", zap.Uint64("accountId", uid), zap.Error(err))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	return msg.ErrCode_SUCC, ret.(*PlayerDetailInfo)
}

func GetPlayerBasic(uid uint64) (msg.ErrCode, *model.AccBasic) {
	err, ret := GetPlayerSimpleInfo(uid)
	if err != msg.ErrCode_SUCC {
		return err, nil
	}
	return msg.ErrCode_SUCC, PlayerSimpleInfoToAccBacic(ret)
}

func GetPlayerSimpleInfo(uid uint64) (msg.ErrCode, *PlayerSimpleInfo) {
	ret, err := playerCache.Take(fmt.Sprintf("%v", uid), func() (any, error) {
		var u *model.UserData
		p := player.FindByUserId(uid)
		if p != nil {
			u = p.UserData
		} else {
			var er error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			u, er = player.LoadUser(ctx, uid) // TODO redis + 布隆
			if er != nil {
				return nil, er
			}
		}
		if u == nil {
			return nil, fmt.Errorf("user:%v by uid not exist", uid)
		}

		data := &PlayerSimpleInfo{
			Uid:            u.UserId,
			Name:           u.Nick,
			Level:          u.Level,
			MissionId:      int(u.StageInfo.MissionId),
			Head:           u.HeadImg,
			HeadFrame:      u.HeadFrame,
			Title:          u.Title,
			ShipId:         u.BaseInfo.ShipId,
			UsePet:         u.BaseInfo.UsePet,
			LastOnlineTime: u.BaseInfo.LogoutTime,
			Combat:         u.BaseInfo.Combat,
			ServerId:       u.ServerId,
			ServerName:     u.ServerName,
			CoatId:         u.Ships.GetShipCoatId(u.BaseInfo.ShipId),
		}

		//  TODO: 公会联盟信息
		// allianceInfo := ServMgr.GetAllianceService().GetAllianceInfo(account.AccountId)
		// if allianceInfo != nil {
		// 	info.AllianceId = allianceInfo.ID
		// 	info.AllianceName = allianceInfo.Name
		// }
		return data, nil
	})

	if err != nil {
		log.Error("GetPlayerSimpleInfo OnInit err", zap.Uint64("accountId", uid), zap.Error(err))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	return msg.ErrCode_SUCC, ret.(*PlayerSimpleInfo)
}

func GetPlayerSimpleInfoByNick(nick string) (msg.ErrCode, *PlayerSimpleInfo) {
	ret, err := playerNickCache.Take(fmt.Sprintf("%s", nick), func() (any, error) {
		var u *model.UserData
		p := player.FindByNick(nick)
		if p != nil {
			u = p.UserData
		} else {
			var er error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			u, er = player.LoadByNickname(ctx, nick) // TODO redis + 布隆
			if er != nil {
				return nil, er
			}
		}
		if u == nil {
			return nil, fmt.Errorf("user:%v by nick not exist", nick)
		}
		data := &PlayerSimpleInfo{
			Uid:            u.UserId,
			Name:           u.Nick,
			Level:          u.Level,
			MissionId:      int(u.StageInfo.MissionId),
			Head:           u.HeadImg,
			HeadFrame:      u.HeadFrame,
			Title:          u.Title,
			ShipId:         u.BaseInfo.ShipId,
			UsePet:         u.BaseInfo.UsePet,
			LastOnlineTime: u.BaseInfo.LogoutTime,
			Combat:         u.BaseInfo.Combat,
			ServerId:       u.ServerId,
			ServerName:     u.ServerName,
			CoatId:         u.Ships.GetShipCoatId(u.BaseInfo.ShipId),
		}

		//  TODO: 公会联盟信息
		// allianceInfo := ServMgr.GetAllianceService().GetAllianceInfo(account.AccountId)
		// if allianceInfo != nil {
		// 	info.AllianceId = allianceInfo.ID
		// 	info.AllianceName = allianceInfo.Name
		// }
		return data, nil
	})

	if err != nil {
		log.Error("GetPlayerSimpleInfo OnInit err", zap.String("nick", nick), zap.Error(err))
		return msg.ErrCode_SYSTEM_ERROR, nil
	}

	return msg.ErrCode_SUCC, ret.(*PlayerSimpleInfo)
}

func GetPlayerSimpleInfos(ids []uint64) (msg.ErrCode, map[uint64]*PlayerSimpleInfo) {
	ret := make(map[uint64]*PlayerSimpleInfo)
	for _, id := range ids {
		er, v := GetPlayerSimpleInfo(id)
		if er != msg.ErrCode_SUCC {
			return er, nil
		}
		ret[id] = v
	}
	return msg.ErrCode_SUCC, ret
}

func ToProtocolPlayerDetail(data *PlayerDetailInfo, p *player.Player) *msg.PlayerDetailInfo {
	ret := &msg.PlayerDetailInfo{
		AccountId: uint32(data.AccountId),
		Name:      data.Name,
		Level:     data.Level,
		MissionId: uint32(data.MissionId),
		Combat:    data.Combat,
		ShipId:    data.ShipId,
		Equips:    ToProtocolEquipPosList(data.Equips),
		Head:      data.Head,
		HeadFrame: data.HeadFrame,
		Title:     data.Title,
		// TODO 老宠物移除暂留
		// Pets:         ToProtocolPets(data.Pets),
		ServerId:     data.ServerId,
		ServerName:   data.ServerName,
		AllianceId:   data.AllianceId,
		AllianceName: data.AllianceName,
		CoatId:       int32(data.CoatId),
	}

	for i := 0; i < len(data.SupportInfos); i++ {
		shipInfo := ToProtocolShip(data.SupportInfos[i], p)
		//shipInfo.SupportAttr = ToProtocolAttrs(data.SupportInfos[i].SupportAttr)
		ret.SupportInfos = append(ret.SupportInfos, shipInfo)
	}
	return ret
}

func ToPlayerSimpleInfo(data *PlayerSimpleInfo) *msg.PlayerSimpleInfo {
	if data == nil {
		return nil
	}
	var online uint32 = 0
	if player := player.FindByUserId(data.Uid); player != nil && player.IsOnline() {
		online = 1
	}

	return &msg.PlayerSimpleInfo{
		AccountId:    data.Uid,
		Name:         data.Name,
		Level:        data.Level,
		MissionId:    uint32(data.MissionId),
		ShipId:       data.ShipId,
		Head:         data.Head,
		HeadFrame:    data.HeadFrame,
		Title:        data.Title,
		PetId:        data.UsePet,
		Online:       online,
		Combat:       data.Combat,
		ServerId:     data.ServerId,
		ServerName:   data.ServerName,
		AllianceId:   data.AllianceId,
		AllianceName: data.AllianceName,
		CoatId:       int32(data.CoatId),
	}
}

func ToPlayerSimpleInfos(data []*PlayerSimpleInfo) []*msg.PlayerSimpleInfo {
	var ret []*msg.PlayerSimpleInfo
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToPlayerSimpleInfo(data[i]))
	}
	return ret
}

func PlayerSimpleInfoToAccBacic(data *PlayerSimpleInfo) *model.AccBasic {
	if data == nil {
		return nil
	}

	info := &model.AccBasic{}
	info.AccountId = int64(data.Uid)
	info.HeadFrame = data.HeadFrame
	info.HeadImg = data.Head
	info.Nick = data.Name
	info.ShipId = data.ShipId
	info.Title = data.Title
	return info
}

func getSimpleInfoFromUser(user_data *model.UserData) *PlayerSimpleInfo {
	if user_data == nil {
		return nil
	}

	info := &PlayerSimpleInfo{
		Uid:            user_data.UserId,
		Name:           user_data.Nick,
		Level:          user_data.Level,
		MissionId:      int(user_data.StageInfo.MissionId),
		Head:           user_data.HeadImg,
		HeadFrame:      user_data.HeadFrame,
		Title:          user_data.Title,
		ShipId:         user_data.BaseInfo.ShipId,
		UsePet:         user_data.BaseInfo.UsePet,
		LastOnlineTime: user_data.BaseInfo.LogoutTime,
		Combat:         user_data.BaseInfo.Combat,
		ServerId:       user_data.ServerId,
		ServerName:     user_data.ServerName,
		CoatId:         user_data.Ships.GetShipCoatId(user_data.BaseInfo.ShipId),
	}

	// TODO : 公会联盟
	// allianceInfo := ServMgr.GetAllianceService().GetAllianceInfo(account.AccountId)
	// if allianceInfo != nil {
	// 	info.AllianceId = allianceInfo.ID
	// 	info.AllianceName = allianceInfo.Name
	// }

	return info
}
