package common

import (
	"errors"
	"gameserver/internal/config"
	"kernel/data/db_global"
	"kernel/kenum"
	"kernel/tools"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	//serverInfo *model.ServerInfo
	serverInfo *db_global.ServerInfo
)

func GetServerInfo() *db_global.ServerInfo {
	return serverInfo
}

func LoadServerInfo() {
	var err error
	if config.Conf.ServerId == 0 {
		log.Panic("server id is 0")
	}
	serverInfo, err = db_global.GetServerInfoModel().GetOne(uint64(config.Conf.ServerId))
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Error("get db server info err", zap.Error(err), zap.Uint32("sid", config.Conf.ServerId))
		return
	}
	if serverInfo == nil {
		SerID := uint64(config.Conf.ServerId)
		if SerID == 0 {
			SerID, _ = db_global.GenServerInfoIdSeq()
		}
		newServerInfo := &db_global.ServerInfo{
			SerID:         SerID,
			SerName:       config.Conf.ServerName,
			SerAddr:       config.Conf.Tcp.Addr,
			RegisterLimit: kenum.Server_Register_Limit,
			OpenTime:      time.Now().Unix(),
			DisplayTime:   time.Now().Unix(),
			Whitelist:     config.Conf.Whitelist,
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
		}

		if config.Conf.Ws.Addr != "" {
			newServerInfo.SerAddr = config.Conf.Ws.Addr
		}

		if _, err = db_global.GetServerInfoModel().Upsert(newServerInfo); err != nil {
			log.Error("upsert server info err", zap.Error(err))
			return
		}
		serverInfo = newServerInfo
	}

	//serverInfo = dao.ServerInfoDao.GetServerInfo(config.Conf.ServerId)
	//if serverInfo == nil {
	//	serverInfo = dao.ServerInfoDao.AddServerInfo(config.Conf.ServerId, config.Conf.ServerName)
	//}
}

func GetOpenServerDays() uint32 {
	info := GetServerInfo()
	return tools.GetDiffDay(time.Unix(info.OpenTime, 0), time.Now()) + 1
}

// func UpdateRegisterNum() {
// 	serverInfo.RegisterNum++
// 	dao.ServerInfoDao.UpdateRegistNum(config.Conf.ServerId)

// 	//serverInfo.RegistNum += 1
// 	//dao.ServerInfoDao.UpdateRegistNum(config.Conf.ServerId)
// }

// func RegisterFull() bool {
// 	//if serverInfo.RegistNum >= int64(config.Conf.MaxRegisterNum) {
// 	//	return true
// 	//}
// 	return false
// }
