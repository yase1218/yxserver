package go_tg_bot

import (
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var defTgBot *TgBot

func InitTgBot(ctx context.Context, opts ...any) (err error) {
	defTgBot = NewTgBot()
	if err = defTgBot.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func Get() *gotgbot.Bot {
	return defTgBot.Get()
}

func GetCtx() context.Context {
	return defTgBot.GetCtx()
}

func AddHandle(handler ext.Handler) {
	defTgBot.AddHandle(handler)
}

func Start() {
	defTgBot.Start()
}

func StartWebHook() {
	defTgBot.StartWebHook()
}

func ProcessUpdate(update *gotgbot.Update) error {
	return defTgBot.ProcessUpdate(update)
}