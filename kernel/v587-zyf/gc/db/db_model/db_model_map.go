package db_model

import (
	"context"
	"database/sql"
	"github.com/astaxie/beego/orm"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/iface"
	"gopkg.in/gorp.v1"
	"time"
)

type DBModelMapUnit struct {
	model  iface.IDBModel
	initer func(dbMap *gorp.DbMap)
}
type DBModelMap struct {
	options *DBModelOption

	modelMap map[string][]DBModelMapUnit

	ctx    context.Context
	cancel context.CancelFunc
}

func NewDBModelMap() *DBModelMap {
	m := &DBModelMap{
		options: NewDBModelOption(),
	}

	return m
}

func (dmm *DBModelMap) Init(ctx context.Context, opts ...any) (err error) {
	dmm.ctx, dmm.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.(Option)(dmm.options)
		}
	}

	return
}

func (dmm *DBModelMap) Register(dbKey string, model iface.IDBModel, initer func(dbMap *gorp.DbMap)) {
	mapItems, ok := dmm.modelMap[dbKey]
	if !ok {
		mapItems = make([]DBModelMapUnit, 0, 2)
		dmm.modelMap[dbKey] = mapItems
	}

	for _, mi := range mapItems {
		if mi.model == model {
			return
		}
	}

	dmm.modelMap[dbKey] = append(mapItems, DBModelMapUnit{model: model, initer: initer})
}

func (dmm *DBModelMap) Start() error {
	for dbKey := range dmm.modelMap {
		dbInfo, ok := dmm.options.databases[dbKey]
		if !ok || dbInfo.uri == "" {
			continue
		}

		if err := dmm.link(dbKey, dbInfo); err != nil {
			return err
		}
	}

	return nil
}

func (dmm *DBModelMap) link(dbKey string, dbInfo *DBModelOptionUnit) error {
	db, err := sql.Open("mysql", dbInfo.uri)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	{
		maxIdle := dbInfo.maxIdle
		if maxIdle <= 0 {
			maxIdle = enums.DB_MAX_IDLE_CONN
		}
		db.SetMaxIdleConns(maxIdle)
	}
	{
		maxOpenCon := dbInfo.maxOpenConn
		if maxOpenCon <= 0 {
			maxOpenCon = enums.DB_MAX_CONN
		}
		db.SetMaxOpenConns(maxOpenCon)
	}
	{
		maxLifetime := dbInfo.conn_max_lifetime
		if maxLifetime <= 0 {
			maxLifetime = enums.DB_CONN_LIFETIME
		}
		db.SetConnMaxLifetime(time.Duration(maxLifetime))
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	if dmm.options.enableDbTrace {
		dbMap.TraceOn("[gorp]", &DbLogger{})
	}

	mapItems, ok := dmm.modelMap[dbKey]
	if ok {
		for _, mi := range mapItems {
			mi.model.SetDbMap(dbMap)
			mi.model.SetDb(db)
			mi.initer(dbMap)
		}
	}

	if dbInfo.ormDefaultDb {
		if err = orm.AddAliasWthDB("default", "mysql", db); err != nil {
			return err
		}
	}

	if dbInfo.tableAutoCheck {
		if err = orm.AddAliasWthDB(dbKey, "mysql", db); err != nil {
			return err
		}
		if err = orm.RunSyncdb(dbKey, false, true); err != nil {
			return err
		}
	}

	return nil
}
