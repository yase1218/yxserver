package db_model

import (
	"database/sql"
	"gopkg.in/gorp.v1"
)

type DBBaseModel struct {
	dbMap *gorp.DbMap
	db    *sql.DB
}

func (dbm *DBBaseModel) SetDbMap(dbMap *gorp.DbMap) {
	dbm.dbMap = dbMap
}
func (dbm *DBBaseModel) DbMap() *gorp.DbMap {
	return dbm.dbMap
}
func (dbm *DBBaseModel) SetDb(db *sql.DB) {
	dbm.db = db
}
func (dbm *DBBaseModel) Db() *sql.DB {
	return dbm.db
}
