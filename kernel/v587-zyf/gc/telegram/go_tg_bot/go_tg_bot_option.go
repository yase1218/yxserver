package go_tg_bot

type TgBotOption struct {
	token string

	webHookHost string
	webHookDir  string
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

func WithWebHookHost(webHookHost string) Option {
	return func(opts *TgBotOption) {
		opts.webHookHost = webHookHost
	}
}

func WithWebHookDir(webHookDir string) Option {
	return func(opts *TgBotOption) {
		opts.webHookDir = webHookDir
	}
}
