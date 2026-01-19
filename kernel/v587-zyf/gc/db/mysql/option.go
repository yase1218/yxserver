package mysql

type MysqlOption struct {
	uri               string
	max_conn          int
	max_idle_conn     int
	conn_max_lifetime int
}

type Option func(o *MysqlOption)

func NewMysqlOption() *MysqlOption {
	return &MysqlOption{}
}

func WithUri(uri string) Option {
	return func(o *MysqlOption) {
		o.uri = uri
	}
}

func WithMaxConn(max_conn int) Option {
	return func(o *MysqlOption) {
		o.max_conn = max_conn
	}
}

func WithMaxIdleConn(max_idle_conn int) Option {
	return func(o *MysqlOption) {
		o.max_idle_conn = max_idle_conn
	}
}

func WithConnMaxLifetime(conn_max_lifetime int) Option {
	return func(o *MysqlOption) {
		o.conn_max_lifetime = conn_max_lifetime
	}
}
