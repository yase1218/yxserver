package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gameserver/internal/config"
	"gameserver/internal/game/charge"
	"gameserver/internal/game/model"
	"gameserver/internal/game/player"
	"io/ioutil"
	"msg"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

const (
	ChargeVerifyOk        = 0  // 发货成功
	ChargeVerifyMoney     = 9  // 订单金额不一致
	ChargeVerifyProductId = 12 // productId校验失败
	ChargeVerifyUserId    = 14 // userId不一致
	ChargeVerifyUpLimit   = 90 // 购买超限
	ChargeVerifySign      = 96 // 签名不一致
	ChargeVerifyChannelNo = 97 // channelNo不一致
	ChargeVerifyWhite     = 98 // 白名单异常
	ChargeVerifyOther     = 99 // 其他情况
)

func PayCallBack(req *msg.InterRequestPayMsg) {
	if req.Status != ChargeVerifyOk {
		reportFailOrder(req, int(req.Status))
		log.Error("PayCallBack order status not ok", zap.Any("req", req))
		return
	}
	order := charge.GetOrderManager().GetOrderById(req.GameOrderNo)
	if order == nil {
		reportFailOrder(req, ChargeVerifyOther)
		log.Error("PayCallBack order not found", zap.Any("req", req))
		return
	}
	if order.Status == model.OrderShipment || order.Status == model.OrderVerified { // 订单已验证或已发货
		//reportFailOrder(req, ChargeVerifyOk)
		log.Debug("PayCallBack order repeat verify", zap.Any("order", order))
		return
	}
	if status := LeitingVerifyCharge(req, order); status != ChargeVerifyOk {
		reportFailOrder(req, int(status))
		return
	}

	p := player.FindByAccount(req.OpenID)
	if p != nil { // 在线处理
		processOnlinePay(order, p)
		order.Status = model.OrderShipment // 已发货
		// 清零推送倒计时
		ClearOutPut(p, order.ChargeID, true)
	} else { // 不在线
		order.Status = model.OrderVerified // 验证通过
		log.Debug("PayCallBack player offline", zap.Any("order", order))
	}
	order.UpdateAt = time.Now()
	//if b, err = json.Marshal(order); err != nil {
	//	log.Error("PayCallBack Marshal order err", zap.String("accountId", req.OpenID), zap.Any("req", req), zap.Error(err))
	//	return
	//}
	////err = rc.SetEx(rdb_single.GetCtx(), key, string(b), model.RedisOrderTime).Err()
	//err = rc.HSet(rdb_single.GetCtx(), key, req.GameOrderNo, string(b)).Err()
	//if err != nil {
	//	log.Error("PayCallBack redis set order err", zap.String("accountId", req.OpenID), zap.Any("req", req), zap.Error(err))
	//	return
	//}
	//err = rc.HExpire(rdb_single.GetCtx(), key, model.RedisOrderTime, req.GameOrderNo).Err()
	//if err != nil {
	//	log.Error("PayCallBack redis expire order err", zap.String("accountId", req.OpenID), zap.Any("req", req), zap.Error(err))
	//	return
	//}
	order.SaveStatus()
	//log.Debug("PayCallBack reportFailOrder ChargeVerifyOk", zap.Any("order", order))
	//reportFailOrder(req, ChargeVerifyOk)
}

func processOnlinePay(inMsg interface{}, p *player.Player) {
	//req := inMsg.(*msg.InterRequestPayMsg)
	//ChargeCallBack(p, chargeID, 1)
	// 订单发货了 TODO 重构充值
	// dao.OrderDao.UpdateOrderStatus(req.OrderId, 2)

	//req := inMsg.(*msg.InterRequestPayMsg)
	//chargeID, _ := strconv.Atoi(req.ProductId)
	order := inMsg.(*model.OrderInfo)
	ChargeCallBack(p, order.ChargeID, 1)
}

func LeitingVerifyCharge(req *msg.InterRequestPayMsg, order *model.OrderInfo) (status int32) {
	var money int
	var err error
	if money, err = strconv.Atoi(req.Money); err != nil {
		log.Warn("LeitingVerifyCharge money not int ", zap.Any("req", req), zap.Any("order", order))
		return ChargeVerifyMoney
	}
	if money != order.Money {
		log.Warn("LeitingVerifyCharge money fail", zap.Any("req", req), zap.Any("order", order))
		return ChargeVerifyMoney
	}
	if req.ProductId != order.ProductId {
		log.Warn("LeitingVerifyCharge ProductId fail", zap.Any("req", req), zap.Any("order", order))
		return ChargeVerifyProductId
	}
	if req.OpenID != order.AccountId {
		log.Warn("LeitingVerifyCharge OpenID fail", zap.Any("req", req), zap.Any("order", order))
		return ChargeVerifyUserId
	}
	if req.ChannelNo != order.ChannelNo {
		log.Warn("LeitingVerifyCharge ChannelNo fail", zap.Any("req", req), zap.Any("order", order))
		return ChargeVerifyChannelNo
	}
	return ChargeVerifyOk
}

// 计算 MD5 签名
func calculateMD5Signature(game, channelNo, gameOrderNo string, status int, key string) string {
	data := game + channelNo + gameOrderNo + fmt.Sprint(status) + key
	hash := md5.Sum([]byte(data))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

// 订单异常上报接口调用示例
func reportFailOrder(req *msg.InterRequestPayMsg, status int) {
	// 参数
	params := make(map[string]interface{})
	params["game"] = config.Conf.Leiting.Game
	params["channelNo"] = req.ChannelNo
	params["gameOrderNo"] = req.GameOrderNo
	params["status"] = status
	params["ip"] = req.Ip // ip 白名单校验失败（status=98）时，需要把异常 ip 传过来
	params["sign"] = calculateMD5Signature(config.Conf.Leiting.Game, req.ChannelNo, req.GameOrderNo, status, config.Conf.Leiting.Key)

	// 构建 JSON 数据
	jsonParams, _ := json.Marshal(params)

	// 发送 POST 请求
	resp, err := http.Post(config.Conf.Leiting.Domain+"/game_service/deal_fail_bill.do", "application/json", strings.NewReader(string(jsonParams)))
	if err != nil {
		log.Error("reportFailOrder post err:", zap.Int("status", status), zap.Any("req", req), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("reportFailOrder read body:", zap.Any("req", req), zap.Error(err))
		return
	}
	log.Debug("reportFailOrder success body:", zap.Int("status", status), zap.Any("req", req), zap.String("body", string(body)))
}
