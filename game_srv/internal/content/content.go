package content

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gameserver/internal/config"
	"io"
	"kernel/kenum"
	"kernel/tools"
	"msg"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

const (
	ContentTypeChat = iota + 1
	ContentTypeNick
)

type ContentCallBack func(msg.ErrCode)
type PushContentFn func(*ContentCb) error
type ContentData struct {
	Content                                   []string
	UserId, ChannelNo                         string
	Os                                        int
	RoleId, RoleName, RoleLevel, ServerId, Ip string
	Cb                                        ContentCallBack
	ContentType                               int // 1聊天，2昵称修改
	Id                                        uint64
	PushTm                                    time.Time
}

type ContentCb struct {
	Ec msg.ErrCode
	Cb ContentCallBack
}

type ContentService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	stopChan chan struct{}
	wg       sync.WaitGroup
	state    atomic.Uint32

	queue      chan *ContentData
	panicFunc  func(string)
	pushBackFn PushContentFn
}

func (c *ContentService) start() error {
	if !c.state.CompareAndSwap(kenum.WorkState_Idle, kenum.WorkState_Running) {
		return errors.New("ContentService can't start, current state : " + kenum.StateToString(c.state.Load()))
	}

	log.Info("ContentService start running")
	c.wg.Add(1)
	go tools.GoSafePost("ContentService run", func() {
		c.run()
	}, c.panicFunc)

	c.wg.Add(1)
	go tools.GoSafePost("ContentService monitor", func() {
		c.monitor()
	}, c.panicFunc)

	log.Info("ContentService start success")
	return nil
}

func (c *ContentService) monitor() {
	defer c.wg.Done()
	defer func() {
		if c.cancel != nil {
			c.cancel()
		}
	}()
	for {
		select {
		case <-c.stopChan:
			log.Info("ContentService received stop signal, draining messages")
			return
		case <-c.ctx.Done():
			log.Info("ContentService context canceled, exiting")
			return
		case <-time.After(time.Second * 3):
			log.Info("ContentService queue size", zap.Int("size", len(c.queue)))
		}
	}
}

func (c *ContentService) run() {
	defer c.wg.Done()
	defer func() {
		if c.cancel != nil {
			c.cancel()
		}
	}()

	for {
		select {
		case <-c.stopChan:
			log.Info("ContentService received stop signal, draining messages")
			return
		case <-c.ctx.Done():
			log.Info("ContentService context canceled, exiting")
			return
		case m, ok := <-c.queue:
			now := time.Now()
			cost := now.Sub(m.PushTm).Milliseconds()
			log.Info("trace content from push to pull", zap.Int64("cost", cost), zap.Uint64("content id", m.Id))
			if !ok {
				log.Error("ContentService queue channel closed")
				return
			}

			go tools.GoSafePost("ContentService wait stop", func() {
				c.processContent(m)
			}, c.panicFunc)
		}
	}
}

func (c *ContentService) stop() error {
	timeOut := time.Second * 10
	if !c.state.CompareAndSwap(kenum.WorkState_Running, kenum.WorkState_Stopping) {
		return errors.New("ContentService can't stop, current state : " + kenum.StateToString(c.state.Load()))
	}

	log.Info("ContentService stopping")

	// 发送停止信号
	close(c.stopChan)

	// 等待goroutine退出
	stopped := make(chan struct{})
	go tools.GoSafePost("ContentService wait stop", func() {
		c.wg.Wait()
		close(stopped)
	}, c.panicFunc)

	select {
	case <-stopped:
		c.state.Store(kenum.WorkState_Stopped)
		log.Info("ContentService stopped")
		return nil
	case <-time.After(timeOut):
		if c.cancel != nil {
			c.cancel()
		}
		log.Warn("ContentService stop timeout, forcing shutdown", zap.Duration("timeout", timeOut))
		return errors.New("ContentService stop timeout after " + timeOut.String())
	}
}

func (c *ContentService) push(data *ContentData) error {
	if data == nil {
		return nil
	}

	state := c.state.Load()
	if state != kenum.WorkState_Running {
		return errors.New("ContentService cannot push gate msg, state is : " + kenum.StateToString(state))
	}

	select {
	case c.queue <- data:
		data.PushTm = time.Now()
		return nil
	case <-c.stopChan:
		return errors.New("ContentService is stopping, gate msg reject")
	case <-c.ctx.Done():
		return errors.New("ContentService ctx cancelled, gate msg reject")
	default:
		return errors.New("ContentService gate msg queue is full, gate msg reject")
	}
}

func (c *ContentService) processContent(data *ContentData) {
	if data == nil {
		log.Error("ContentService data nil", zap.Any("data", data))
		return
	}
	if data.Cb == nil {
		log.Error("ContentService cb nil", zap.Any("data", data))
		return
	}

	// check content valid return ec
	var ec msg.ErrCode
	for i := 0; i < len(data.Content); i++ {
		if ec = detectText(data.UserId, data.ChannelNo, data.Os, data.RoleId, data.RoleName, data.RoleLevel,
			data.ServerId, data.Content[i], data.Ip, 2, 0, data.ContentType, data.Id); ec != msg.ErrCode_SUCC {
			break
		}
	}
	c.pushBackFn(&ContentCb{
		Ec: ec,
		Cb: data.Cb,
	})
}

// 计算签名
func calculateSign(game, userId, key string) string {
	sign := game + userId + "%" + key
	hash := md5.Sum([]byte(sign))
	return hex.EncodeToString(hash[:])
}

/**
 * 敏感词检测接口 - 全文本检测
 * 应用场景：
 * 除去角色昵称和聊天的场景，其他需要内容审核的地方需要接入这个接口，例如公告，个性签名等非角色名昵称和聊天的场景。
 * 注意：该接口为游戏服务端调用，请勿在客户端调用该接口。如果有客户端调用需求请联系@易健
 */
func detectText(userId, channelNo string, os int, roleId, roleName, roleLevel, serverId, content, ip string, isYd, isReplace int, contentType int, id uint64) msg.ErrCode {
	now := time.Now()
	if config.Conf.Debug && channelNo == "" {
		channelNo = "110001" // 雷霆（方便测试） android
		if os == 2 {         // ios
			channelNo = "210009" // 雷霆（方便测试）
		}
	}
	// 计算签名
	sign := calculateSign(config.Conf.Leiting.Game, userId, config.Conf.Leiting.Key)
	signCost := time.Since(now).Milliseconds()
	log.Info("trace content sign cost", zap.Int64("cost", signCost), zap.Uint64("content id", id))

	// 创建请求体数据
	jsonData := fmt.Sprintf(`{"game":"%s","userId":"%s","channelNo":"%s","os":%d,"roleId":"%s","roleName":"%s","roleLevel":"%s","serverId":"%s","content":"%s","ip":"%s","isYd":%d,"isReplace":%d,"sign":"%s"}`,
		config.Conf.Leiting.Game, // 游戏标识
		userId,                   // 手游账号
		channelNo,                // 渠道ID。例如110003：九游；130009：4399
		os,                       // 系统 1:android 2 ios 3 web 4 pc端
		roleId,                   // 玩家角色ID
		roleName,                 // 角色昵称
		roleLevel,                // 游戏等级
		serverId,                 // 游戏区组
		content,                  // 文本内容 长度小于2000用户发表内容，建议对内容中JSON、表情符,回车，换行、HTML标签、UBB标签等做过滤，只传递纯文本，以减少误判概率
		ip,                       // 玩家客户端IP地址
		isYd,                     // 默认 填2
		isReplace,                // 是否将拦截敏感词替换 1替换 0不替换  默认0，该值需传1。特别提醒：会存在三方服务未返回敏感词，从而导致敏感文本无法被替换为*的情况。
		sign,                     // 签名 MD5(game+userId+%+KEY)转小写，key使用ltconfig->leiting->api.key.{game}
	)
	log.Debug("detectText req ", zap.Any("jsonData", jsonData))
	var url = config.Conf.Leiting.TextCheckUrl + "/common/sensitive/detection/text"
	if contentType == ContentTypeNick {
		url = config.Conf.Leiting.TextCheckUrl + "/common/sensitive/detection/nickname"
	}
	// 发起 POST 请求
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonData))
	if err != nil {
		log.Error("detectText creating request err:", zap.Error(err))
		return msg.ERRCODE_SYSTEM_ERROR
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("detectText sending request:", zap.Error(err))
		return msg.ERRCODE_SYSTEM_ERROR
	}
	defer resp.Body.Close()

	reqCost := time.Since(now).Milliseconds() - signCost
	log.Info("trace content http req cost ", zap.Int64("cost", reqCost), zap.Uint64("content id", id))

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("detectText reading response body:", zap.Error(err))
		return msg.ERRCODE_SYSTEM_ERROR
	}
	log.Debug("detectText response:", zap.ByteString("body", body))

	readCost := time.Since(now).Milliseconds() - signCost - reqCost
	log.Info("trace content http read response cost ", zap.Int64("cost", readCost), zap.Uint64("content id", id))

	// 处理接口返回结果
	res := new(TextCheckResponse)
	if err = json.Unmarshal(body, res); err != nil {
		log.Error("detectText unmarshal response body err:", zap.Error(err))
		return msg.ERRCODE_SYSTEM_ERROR
	}
	if res.Data.Code != 1 {
		if contentType == ContentTypeNick {
			return msg.ErrCode_NICK_HAS_FORBIDDEN
		}
		return msg.ErrCode_CONTENT_HAS_FORBIDDEN // 敏感词检测不通过
	}

	totalCost := time.Since(now).Milliseconds()
	log.Info("trace content http total cost", zap.Int64("cost", totalCost), zap.Uint64("content id", id))
	return msg.ErrCode_SUCC
}

type TextCheckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"data"`
}
