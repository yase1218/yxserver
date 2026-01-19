package handle

import (
	"fmt"
	"gameserver/internal/game/player"
	"kernel/tools"
	"msg"

	"github.com/v587-zyf/gc/log"
	"github.com/zy/game_data/template"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"gameserver/internal/config"
	"gameserver/internal/game/model"
	"gameserver/internal/game/service"
)

type MsgHandle func(serverId int64, data []byte)

var (
	msgHandleMap = make(map[uint32]MsgHandle)
)

func init() {
	//msgHandleMap[msg.MsgId_ID_ResponseInterRegistUid] = ResponseInterRegistUid
	msgHandleMap[uint32(msg.InterMsgId_ID_InterRequestGetServerInfo)] = InterRequestGetServerInfo
	msgHandleMap[uint32(msg.InterMsgId_ID_InterNotifyMail)] = InterNotifyMail
	msgHandleMap[uint32(msg.InterMsgId_ID_InterNotifyBanner)] = InterNotifyBanner
	msgHandleMap[uint32(msg.InterMsgId_ID_InterRequestGmMsg)] = GMCommand
	msgHandleMap[uint32(msg.InterMsgId_ID_InterDelMail)] = InterDelMail
	msgHandleMap[uint32(msg.InterMsgId_ID_InterRequestPayMsg)] = InterRequestPay
	msgHandleMap[uint32(msg.InterMsgId_ID_InterModifyShopItemVer)] = InterModifyShopItemVer
	msgHandleMap[uint32(msg.InterMsgId_ID_InterResponseUseCdk)] = InterMResponseUseCdk
	msgHandleMap[uint32(msg.InterMsgId_ID_InterRequestAd)] = InterRequestAd
	msgHandleMap[uint32(msg.InterMsgId_ID_InterQuestion)] = InterQuestion
}

func RouteInterMsg(req *msg.RequestCommonInterMsg) {
	if h, ok := msgHandleMap[req.MsgId]; ok {
		h(req.ServerId, req.Data)
	} else {
		log.Error("RouteInterMsg no handle", zap.Uint32("msgID", req.MsgId))
	}
}

func InterRequestGetServerInfo(serverId int64, data []byte) {
	subject := fmt.Sprintf("%v", serverId)
	interRes := &msg.InterResponseGetServerInfo{}
	interRes.Result = msg.InterErrCode_SUCC
	interRes.Env = config.Conf.Env
	interRes.OnlineNum = player.OnlineNum()
	service.PublisInterMsg(subject, uint32(msg.InterMsgId_ID_InterResponseGetServerInfo), interRes)
}

func InterNotifyMail(serverId int64, data []byte) {
	mailMsg := msg.InterNotifyMail{}
	err := proto.Unmarshal(data, &mailMsg)
	if err != nil {
		log.Error("InterNotifyMail err", zap.Error(err))
		//log.Errorf("InterNotifyMail err:%v", err)
		return
	}

	var items []*model.SimpleItem
	for i := 0; i < len(mailMsg.Attachment); i++ {
		items = append(items, &model.SimpleItem{
			Id:  mailMsg.Attachment[i].ItemId,
			Num: mailMsg.Attachment[i].ItemNum,
		})
	}

	gmail := model.CreateGlobalMail(mailMsg.MailId, mailMsg.Title, mailMsg.Content, items, mailMsg.EndTime, mailMsg.RoleStartTime, mailMsg.RoleEndTime)
	if mailMsg.AccountId == 0 {
		service.AddGlobalMail(gmail)
		players := player.AllPlayers()
		for _, p := range players {
			service.AddUserGlobalMail(p, gmail)
		}
	} else {
		service.AddOfflineMail(uint64(mailMsg.AccountId), gmail.FmtToUser())
	}
	log.Info("InterNotifyMail err", zap.Int64("accountId", mailMsg.AccountId), zap.Int64("mailId", mailMsg.MailId),
		zap.String("title", mailMsg.Title), zap.String("content", mailMsg.Content), zap.Any("attachment", mailMsg.Attachment))
	//log.Infof("InterNotifyMail accountId:%v mailId:%v title:%v content:%v attachment:%v", mailMsg.AccountId, mailMsg.MailId, mailMsg.Title, mailMsg.Content, mailMsg.Attachment)
}

// InterDelMail 删除邮件
func InterDelMail(serverId int64, data []byte) {
	delMailMsg := msg.InterDelMail{}
	err := proto.Unmarshal(data, &delMailMsg)
	if err != nil {
		log.Error("InterDelMail err", zap.Error(err))
		//log.Errorf("InterNotifyMail err:%v", err)
		return
	}

	if delMailMsg.AccountId == 0 {
		service.DelGlobalMail(delMailMsg.MailId)
	} else {
		// fix 这个接口就不要提供了
		//service.GMDeleteMailById(uint64(delMailMsg.AccountId), delMailMsg.MailId)
	}
	log.Info("InterDelMail", zap.Int64("accountId", delMailMsg.AccountId), zap.Int64("mailId", delMailMsg.MailId))
}

func InterNotifyBanner(serverId int64, data []byte) {
	bannerMsg := msg.InterNotifyBanner{}
	err := proto.Unmarshal(data, &bannerMsg)
	if err != nil {
		log.Error("InterNotifyBanner err", zap.Error(err))
		//log.Errorf("InterNotice err:%v", err)
		return
	}

	service.SendBannerMsg(&bannerMsg)
}

func InterRequestPay(serverId int64, data []byte) {
	req := new(msg.InterRequestPayMsg)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("InterRequestPay err", zap.Error(err))
		//log.Errorf("InterRequestPay err:%v", err)
		return
	}
	log.Info("InterRequestPay", zap.Any("req", req))
	//log.Infof("InterRequestPay data:%+v", req)
	service.PayCallBack(req)
}

func InterModifyShopItemVer(serverId int64, data []byte) {
	req := new(msg.InterModifyShopItemVer)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("InterModifyShopItemVer err", zap.Error(err))
		//log.Errorf("InterModifyShopItemVer err:%v", err)
		return
	}
	log.Info("InterModifyShopItemVer", zap.Any("req", req))
	//log.Infof("InterModifyShopItemVer data:%v", req)
	if req.IsDelete == 1 {
		service.DelShopItemVer(req.ShopItemId)
	} else {
		service.UpdateShopItemVer(req.ShopItemId, req.Ver)
	}
}

func InterMResponseUseCdk(serverId int64, data []byte) {
	req := new(msg.InterResponseUseCdk)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("InterMResponseUseCdk err", zap.Error(err))
		return
	}
	log.Info("InterMResponseUseCdk", zap.Any("req", req))
	p := player.FindByUserId(uint64(req.AccountId))
	if p != nil {
		// TODO 此时玩家下线 会丢单 改为异步
		service.InterUseCdkResponse(req, p)
	}
}

func InterRequestAd(serverId int64, data []byte) {
	req := new(msg.InterRequestAd)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("InterRequestAd err", zap.Error(err))
		//log.Errorf("InterRequestAd err:%v", err)
		return
	}

	log.Info("InterRequestAd", zap.Any("req", req))
	//log.Infof("InterRequestAd data:%v", req)
	p := player.FindByUserId(uint64(req.AccountId))

	if p != nil {
		service.InterEndAd(req, p)
	} else {
		log.Error("InterRequestAd not online", zap.Uint32("accountId", req.AccountId))
		//log.Errorf("InterRequestAd %v not online", req.AccountId)
	}
}

func GMCommand(serverId int64, data []byte) {
	gmMsg := msg.InterRequestGmMsg{}
	err := proto.Unmarshal(data, &gmMsg)
	if err != nil {
		log.Error("GMCommand err", zap.Error(err))
		//log.Errorf("GMCommand %v", err)
		return
	}

	if p := player.FindByUserId(uint64(gmMsg.AccountId)); p != nil {
		service.ProcessCommand(msg.GMCommandId(gmMsg.CmdId), gmMsg.Content, p)
	} else {
		log.Error("GMCommand not online", zap.Int64("accountId", gmMsg.AccountId))
	}
}

func InterQuestion(serverId int64, data []byte) {
	req := new(msg.InterQuestion)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("InterQuestion err", zap.Error(err))
		return
	}

	p := player.FindByUserId(uint64(req.AccountId))
	if p == nil {
		log.Error("InterQuestion err", zap.Error(fmt.Errorf("player not find, accountId:%d", req.AccountId)))
		return
	}

	cfg := template.GetSystemItemTemplate().GetQuestionInfo(req.QuestionId)
	if cfg == nil {
		log.Error("InterQuestion err,", zap.Error(fmt.Errorf("cfg not find , accountId:%d, QuestionId:%s,", req.AccountId, req.QuestionId)))
		return
	}

	if tools.ListStrContain(p.UserData.BaseInfo.QuestionIds, req.QuestionId) {
		log.Error("InterQuestion err,", zap.Error(fmt.Errorf("question is exist, accountId:%d, QuestionId:%s,", req.AccountId, req.QuestionId)))
		return
	}

	p.UserData.BaseInfo.QuestionIds = append(p.UserData.BaseInfo.QuestionIds, req.QuestionId)
	p.SaveBaseInfo()

	ntf := &msg.QuestionFinishNtf{
		QuestionId: req.QuestionId,
	}
	p.SendNotify(ntf)
}
