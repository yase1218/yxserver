package tg_bot_api

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

type TgBot struct {
	options *TgBotOption

	ctx    context.Context
	cancel context.CancelFunc

	Bot *tgbotapi.BotAPI
}

func NewTgBot() *TgBot {
	t := &TgBot{
		options: NewGrpcOption(),
	}

	return t
}

func (t *TgBot) Init(ctx context.Context, option ...any) (err error) {
	t.ctx, t.cancel = context.WithCancel(ctx)

	for _, opt := range option {
		opt.(Option)(t.options)
	}

	t.Bot, err = tgbotapi.NewBotAPI(t.options.token)
	if err != nil {
		log.Error("Failed to connect to Telegram", zap.Error(err))
		return err
	}

	return nil
}

func (t *TgBot) Get() *tgbotapi.BotAPI {
	return t.Bot
}

func (t *TgBot) GetCtx() context.Context {
	return t.ctx
}

/* 获取玩家头像 && 下载下来
 * tgUpCfg := tgbotapi.UserProfilePhotosConfig{
		UserID: int64(tgUser.UserID),
		Limit:  1,
	}
	tgPhotos, err := tg_bot_api.Get().GetUserProfilePhotos(tgUpCfg)
	if err != nil {
		return nil, err
	}
	if len(tgPhotos.Photos) > 0 && len(tgPhotos.Photos[0]) > 0 {
		headFile, err := tg_bot_api.Get().GetFile(tgbotapi.FileConfig{FileID: tgPhotos.Photos[0][len(tgPhotos.Photos[0])-1].FileID})
		if err == nil {
			telegram.Head = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", tg_bot_api.Get().Token, headFile.FilePath)
		}
	}

	if len(tgPhotos.Photos) > 0 {
		headFile, err := tg_bot_api.Get().GetFile(tgbotapi.FileConfig{FileID: tgPhotos.Photos[0][len(tgPhotos.Photos[0])-1].FileID})
		if err == nil {
			resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", tg_bot_api.Get().Token, headFile.FilePath))
			if err == nil {
				defer resp.Body.Close()
				data, err := io.ReadAll(resp.Body)
				if err == nil {
					targetDir := GetHandleOps().Tg_Head_Path
					if _, err := os.Stat(targetDir); os.IsNotExist(err) {
						os.MkdirAll(targetDir, 0755)
					}
					headName := fmt.Sprintf("%s/%v.jpg", targetDir, uid)
					if err = os.WriteFile(headName, data, 0644); err != nil {
						log.Error("save tg user head err", zap.Error(err))
					} else {
						telegram.Head = GetHandleOps().Tg_Head_Link + strings.Replace(headName, ".", "", 1)
					}
				}
			}
		}
	}
*/
