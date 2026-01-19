package iface

import (
	"database/sql"
	"gopkg.in/gorp.v1"
)

type IDBModel interface {
	SetDbMap(dbMap *gorp.DbMap)
	DbMap() *gorp.DbMap
	SetDb(db *sql.DB)
	Db() *sql.DB
}

type BaseTableInterface interface {
	TableName() string
	TableEngine() string
	TableEncode() string
	TableComment() string
	TableIndex() [][]string
	TableUnique() [][]string
}
