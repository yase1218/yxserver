package elastic

type ElasticOption struct {
	https bool

	host string
	port int

	userName string
	password string
}

type Option func(o *ElasticOption)

func NewMysqlOption() *ElasticOption {
	return &ElasticOption{}
}

func WithHttps(https bool) Option {
	return func(o *ElasticOption) {
		o.https = https
	}
}
func WithHost(host string) Option {
	return func(o *ElasticOption) {
		o.host = host
	}
}

func WithPort(port int) Option {
	return func(o *ElasticOption) {
		o.port = port
	}
}

func WithUserName(userName string) Option {
	return func(o *ElasticOption) {
		o.userName = userName
	}
}

func WithPassword(password string) Option {
	return func(o *ElasticOption) {
		o.password = password
	}
}
