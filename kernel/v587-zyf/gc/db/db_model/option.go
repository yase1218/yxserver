package db_model

type DBModelOptionUnit struct {
	uri               string
	maxIdle           int
	maxOpenConn       int
	conn_max_lifetime int

	ormDefaultDb   bool
	tableAutoCheck bool
}
type DBModelOption struct {
	databases     map[string]*DBModelOptionUnit
	enableDbTrace bool
}

type Option func(o *DBModelOption)

func NewDBModelOption() *DBModelOption {
	return &DBModelOption{}
}

func WithDatabase(dbName string, dbu *DBModelOptionUnit) Option {
	return func(o *DBModelOption) {
		o.databases[dbName] = dbu
	}
}
func WithEnableDbTrace(enableDbTrace bool) Option {
	return func(o *DBModelOption) {
		o.enableDbTrace = enableDbTrace
	}
}
