package mysql

import (
	"context"
	"github.com/v587-zyf/gc/enums"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Mysql struct {
	options       *MysqlOption
	db            *gorm.DB
	autoCreateDbs []any

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMysql() *Mysql {
	m := &Mysql{
		options: NewMysqlOption(),
	}

	return m
}

func (m *Mysql) Init(ctx context.Context, opts ...any) (err error) {
	m.ctx, m.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(m.options)
		}
	}

	dialector := mysql.Open(m.options.uri)
	gormConf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),

		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if m.db, err = gorm.Open(dialector, gormConf); err != nil {
		return
	}

	if err = m.db.Use(&OptimisticLocker{}); err != nil {
		return
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	{
		maxConn := enums.DB_MAX_CONN
		if m.options.max_conn != 0 {
			maxConn = m.options.max_conn
		}
		sqlDB.SetMaxOpenConns(maxConn)
	}
	{
		maxIdleConn := enums.DB_MAX_IDLE_CONN
		if m.options.max_idle_conn != 0 {
			maxIdleConn = m.options.max_idle_conn
		}
		sqlDB.SetMaxIdleConns(maxIdleConn)
	}
	{
		connLifetime := enums.DB_CONN_LIFETIME
		if m.options.conn_max_lifetime != 0 {
			connLifetime = m.options.conn_max_lifetime
		}
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(connLifetime))
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

func (m *Mysql) GetDB() *gorm.DB {
	return m.db
}
func (m *Mysql) GetCtx() context.Context {
	return m.ctx
}

func (m *Mysql) AddAutoCreateDb(db any) {
	m.autoCreateDbs = append(m.autoCreateDbs, db)
}
func (m *Mysql) AutoCreateDbs() (err error) {
	return m.db.AutoMigrate(m.autoCreateDbs...)
}
