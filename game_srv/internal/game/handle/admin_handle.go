package handle

import (
	"context"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"msg"
	"strconv"
	"time"

	"github.com/v587-zyf/gc/log"
	"github.com/zeromicro/go-zero/core/collection"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type AsyncMsgHandle func(serverId int64, data []byte) proto.Message

var (
	asyncMsgHandleMap   = make(map[uint32]AsyncMsgHandle)
	playerCacheByUserId *collection.Cache
	playerCacheByNick   *collection.Cache
)

func init() {
	asyncMsgHandleMap[uint32(msg.InterMsgId_ID_InterGetPlayerDataByNameOrId)] = GetPlayerData

	var err error

	playerCacheByUserId, err = collection.NewCache(5 * time.Minute)
	if err != nil {
		panic(fmt.Errorf("init playerCacheByUserId err:%s", err.Error()))
	}

	playerCacheByNick, err = collection.NewCache(5 * time.Minute)
	if err != nil {
		panic(fmt.Errorf("init playerCacheByNick err:%s", err.Error()))
	}

}

func RouteAdminMsg(req *msg.RequestCommonInterMsg) proto.Message {
	if h, ok := asyncMsgHandleMap[req.MsgId]; ok {
		return h(req.ServerId, req.Data)
	}

	log.Error("RouteAdminMsg no handle", zap.Uint32("msgID", req.MsgId))
	return nil
}

func GetPlayerData(serverId int64, data []byte) proto.Message {
	req := &msg.InterPlayerDataRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		return nil
	}

	var userData *model.UserData
	var err1 error
	if req.UserId == 0 {
		userData, err1 = getPlayerCacheByNick(req.Nick)
	} else {
		userData, err1 = getPlayerCacheByUserId(req.UserId)
	}

	if err1 != nil {
		return nil
	}

	var maxMissionId = 0
	userMission := userData.Mission.Missions
	for _, v := range userMission {
		if v.MissionId > maxMissionId && v.IsPass {
			maxMissionId = v.MissionId
		}
	}

	shipsId := make([]*msg.InterCommonIdAndLv, 0)
	for _, v := range userData.Ships.Ships {
		CommonInfo := &msg.InterCommonIdAndLv{
			Id: v.Id,
			Lv: v.Level,
		}
		shipsId = append(shipsId, CommonInfo)
	}

	weaponIds := make([]*msg.InterCommonIdAndLv, 0)
	for _, v := range userData.Weapon.Weapons {
		CommonInfo := &msg.InterCommonIdAndLv{
			Id: v.Id,
			Lv: v.Level,
		}
		weaponIds = append(weaponIds, CommonInfo)
	}

	loginTimeStr := time.Unix(int64(userData.BaseInfo.LoginTime), 0).Format("2006-01-02 15:04:05")
	logoutTimeStr := time.Unix(int64(userData.BaseInfo.LogoutTime), 0).Format("2006-01-02 15:04:05")
	createTimeStr := time.Unix(int64(userData.BaseInfo.CreateTime), 0).Format("2006-01-02 15:04:05")

	info := &msg.InterPlayerDataResponse{
		UserId:         userData.UserId,
		AccountId:      userData.AccountId,
		ServerId:       userData.ServerId,
		Nick:           userData.Nick,
		Level:          userData.Level,
		MissionId:      uint32(maxMissionId),
		ShipIds:        shipsId,
		WeaponIds:      weaponIds,
		RegistDay:      userData.BaseInfo.LoginCnt,
		LastLoginTime:  loginTimeStr,
		LastLogoutTime: logoutTimeStr,
		CreateTime:     createTimeStr,
		ServerName:     config.Conf.ServerName,
	}
	// return data
	return info
}

func getPlayerCacheByUserId(uid uint64) (*model.UserData, error) {
	cacheKey := strconv.FormatUint(uid, 10)
	ret, err := playerCacheByUserId.Take(cacheKey, func() (any, error) {
		var u *model.UserData
		p := player.FindByUserId(uid)
		if p != nil {
			u = p.UserData
		} else {
			var er error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			u, er = player.LoadUser(ctx, uint64(uid))
			if er != nil {
				return nil, er
			}
		}

		if u == nil {
			return nil, fmt.Errorf("user:%v by uid not exist", uid)
		}

		return u, nil
	})
	if ret == nil {
		return nil, err
	}
	return ret.(*model.UserData), err
}

func getPlayerCacheByNick(nick string) (*model.UserData, error) {
	ret, err := playerCacheByNick.Take(nick, func() (any, error) {
		var u *model.UserData
		p := player.FindByNick(nick)
		if p != nil {
			u = p.UserData
		} else {
			var er error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			u, er = player.LoadByNickname(ctx, nick)
			if er != nil {
				return nil, er
			}
		}
		if u == nil {
			return nil, fmt.Errorf("user:%v by nick not exist", nick)
		}

		return u, nil
	})
	if ret == nil {
		return nil, err
	}
	return ret.(*model.UserData), err
}
