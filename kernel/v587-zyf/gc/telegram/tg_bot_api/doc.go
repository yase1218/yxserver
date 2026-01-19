package tg_bot_api

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var defTgBot *TgBot

func InitTgBot(ctx context.Context, opts ...any) (err error) {
	defTgBot = NewTgBot()
	if err = defTgBot.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *tgbotapi.BotAPI {
	return defTgBot.Get()
}

func GetCtx() context.Context {
	return defTgBot.GetCtx()
}
