package http_server

type HttpOption struct {
	listenAddr string

	isHttps bool
	pem     string
	key     string

	//allowOrigins []string
}

type Option func(opts *HttpOption)

func NewHttpOption() *HttpOption {
	o := &HttpOption{}

	return o
}

func WithListenAddr(addr string) Option {
	return func(opts *HttpOption) {
		opts.listenAddr = addr
	}
}

func WithIsHttps(isHttps bool) Option {
	return func(opts *HttpOption) {
		opts.isHttps = isHttps
	}
}
func WithPem(pem string) Option {
	return func(opts *HttpOption) {
		opts.pem = pem
	}
}
func WithKey(key string) Option {
	return func(opts *HttpOption) {
		opts.key = key
	}
}

//func WithAllOrigins(allowOrigins string) Option {
//	return func(opts *HttpOption) {
//		opts.allowOrigins = allowOrigins
//	}
//}
