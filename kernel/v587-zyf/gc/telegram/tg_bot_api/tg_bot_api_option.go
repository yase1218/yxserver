package tg_bot_api

type TgBotOption struct {
	token string
}

type Option func(opts *TgBotOption)

func NewGrpcOption() *TgBotOption {
	o := &TgBotOption{}

	return o
}

func WithToken(token string) Option {
	return func(opts *TgBotOption) {
		opts.token = token
	}
}
