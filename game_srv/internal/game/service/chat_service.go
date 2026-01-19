package service

import (
	"gameserver/internal/common"
	"gameserver/internal/config"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"gameserver/internal/publicconst"
	"kernel/tools"
	"msg"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

type WorldMsg struct {
	Account *PlayerSimpleInfo
	MsgId   uint32
	Content string
	Para    string
	Tp      uint32
}

var (
	history_chats map[uint32]*WorldMsg
	new_chats     []*WorldMsg
	cur_chat_id   uint32
)

func init() {
	history_chats = make(map[uint32]*WorldMsg)
	new_chats = make([]*WorldMsg, 0)
}

func UpdateChat(now time.Time) {
	brocastWorldChat()
}

func brocastWorldChat() {
	msgs := getSendMsg()
	msgNum := len(msgs)
	if msgNum == 0 {
		return
	}

	msg1 := &msg.BrocastWorldChat{
		Tp:   msg.ChatChannelType_Chat_Channel_World,
		Data: ToProtocolWorldChatInfos(msgs),
	}

	msg2 := &msg.BrocastWorldChat{
		Tp:   msg.ChatChannelType_Chat_Channel_World,
		Data: ToProtocolWorldChatInfos(msgs[msgNum-1:]),
	}

	users := player.AllPlayers()
	for _, p := range users {
		if p.InWorldChannel {
			p.SendNotify(msg1)
		} else {
			p.SendNotify(msg2)
		}
	}
}

// AddWorldChat 添加世界聊天
func addWorldChat(p *player.Player, content, para string) *WorldMsg {
	cur_chat_id++
	if len(history_chats) >= 1000 {
		delete(history_chats, cur_chat_id-1000)
	}
	data := &WorldMsg{
		Account: getSimpleInfoFromUser(p.UserData),
		MsgId:   cur_chat_id,
		Content: content,
		Para:    para,
	}
	history_chats[cur_chat_id] = data
	addSendMsg(data)
	return data
}

func addSendMsg(data *WorldMsg) {
	new_chats = append(new_chats, data)
}

func getSendMsg() []*WorldMsg {
	ret := new_chats
	new_chats = new_chats[0:0]
	return ret
}

func EnterWorldChat(p *player.Player) []*WorldMsg {
	ret := GetWorldChat(0)
	p.InWorldChannel = true
	return ret
}

func LeaveWorldChat(p *player.Player) {
	p.InWorldChannel = false
}

func GetWorldChat(startId uint32) []*WorldMsg {
	startMsgId := startId
	if startId == 0 {
		startMsgId = cur_chat_id
	}
	var ret []*WorldMsg
	for i := int(startMsgId); i >= 0; i-- {
		if data, ok := history_chats[uint32(i)]; ok {
			ret = append(ret, data)
		}
		if len(ret) >= 10 {
			break
		}
	}
	return ret
}

func GetWorldChatData(startId uint32) []*WorldMsg {
	return GetWorldChat(startId)
}

func chat_cmd(p *player.Player, content string) {
	params := strings.Split(content, ",")
	if len(params) < 1 {
		return
	}
	switch params[0] {
	case "mail":
		if len(params) < 2 {
			return
		}
		gmail := &model.GlobalMail{
			MailId:  common.GenSnowFlake(),
			Title:   "mail cmd",
			Content: content,
			Items: []*model.SimpleItem{
				{
					Id:  600002,
					Num: 10000,
				},
			},
			CreateTime: tools.GetCurTime(),
		}
		uid := utils.StrToUInt64(params[1])
		if uid == 0 {
			AddGlobalMail(gmail)
			players := player.AllPlayers()
			for _, u := range players {
				AddUserGlobalMail(u, gmail)
			}
		} else {
			AddOfflineMail(uid, gmail.FmtToUser())
		}
	case "charge":
		if len(params) < 2 {
			return
		}
		errCode, order := CreateOrder(p, &msg.RequestCreateOrder{ChargeId: utils.StrToInt32(params[1])})
		if errCode != msg.ErrCode_SUCC {
			log.Error("gm create order fail", zap.Int("errCode", int(errCode)))
			return
		}
		order.Status = model.OrderShipment
		processOnlinePay(order, p)
		ClearOutPut(p, order.ChargeID, true)
		order.SaveStatus()
	default:
		log.Error("unhandler cmd", zap.String("cmd", params[0]), ZapUser(p))
	}
}

func SendWorldChat(p *player.Player, content, para string) (msg.ErrCode, *WorldMsg) {
	//if config.Conf.IsDebug() {
	//	log.Info("process chat cmd", zap.String("content", content), ZapUser(p))
	//	key := "/inst,"
	//	if strings.HasPrefix(content, key) {
	//		chat_cmd(p, content[len(key):])
	//		return msg.ErrCode_SUCC, &WorldMsg{
	//			Account: getSimpleInfoFromUser(p.UserData),
	//			MsgId:   cur_chat_id,
	//		}
	//	}
	//}

	p.LastWorldChatTime = tools.GetCurTime()
	//return msg.ErrCode_SUCC, addWorldChat(p, template.GetForbiddenTemplate().Filter(content), template.GetForbiddenTemplate().Filter(para))
	return msg.ErrCode_SUCC, addWorldChat(p, content, para)
}

func SendGmHandle(p *player.Player, content string) (bool, *WorldMsg) {
	if config.Conf.IsDebug() {
		log.Info("process chat cmd", zap.String("content", content), ZapUser(p))
		key := "/inst,"
		if strings.HasPrefix(content, key) {
			chat_cmd(p, content[len(key):])
			return true, &WorldMsg{
				Account: getSimpleInfoFromUser(p.UserData),
				MsgId:   cur_chat_id,
			}
		}
	}
	return false, nil
}

func SendWorldCheck(p *player.Player, content, para string) msg.ErrCode {
	if err := FunctionOpen(p, publicconst.World_Chat); err != msg.ErrCode_SUCC {
		return err
	}

	if p.UserData.BaseInfo.ForbiddenChat == 1 {
		return msg.ErrCode_HAS_FORBIDDEN_CHAT
	}

	content = strings.Trim(content, " ")
	if len(content) == 0 {
		return msg.ErrCode_INVALID_DATA
	}
	if utf8.RuneCountInString(content) > 500 {
		return msg.ErrCode_CHAT_MSG_TO_LONG
	}
	para = strings.Trim(para, " ")
	if utf8.RuneCountInString(para) > 500 {
		return msg.ErrCode_CHAT_MSG_TO_LONG
	}
	curTime := tools.GetCurTime()
	if curTime-p.LastWorldChatTime < 5 {
		return msg.ErrCode_WORLD_CHAT_IN_CD
	}
	return msg.ErrCode_SUCC
}

func EmoteSend(p *player.Player, req *msg.EmoteSendReq) error {
	// switch req.GetType() {
	// case msg.EmoteSendType_Emote_Send_PeakFight:
	// 	return ServMgr.GetPeakFightService().BroadcastEmote(p, req)
	// case msg.EmoteSendType_Emote_Send_Chat:
	// 	// todo chat send emote
	// 	fallthrough
	// default:
	// 	log.Error("emote send type not set", zap.Int("type", int(req.GetType())))
	// }
	return nil
}

func ToProtocolWorldChatInfo(data *WorldMsg) *msg.WorldChatInfo {
	ret := &msg.WorldChatInfo{
		Account:     ToPlayerSimpleInfo(data.Account),
		Content:     data.Content,
		MsgId:       data.MsgId,
		Para:        data.Para,
		ContentType: msg.ChatContentType(data.Tp),
	}
	return ret
}

func ToProtocolWorldChatInfos(data []*WorldMsg) []*msg.WorldChatInfo {
	var ret []*msg.WorldChatInfo
	for i := 0; i < len(data); i++ {
		ret = append(ret, ToProtocolWorldChatInfo(data[i]))
	}
	return ret
}
