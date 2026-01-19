package db_global

import (
	"fmt"
)

var (
	DB string
)

const (
	DB_SPACE_WAR_GLOBAL = "space_war_global"
)

const (
	COL_COUNTER     = "counter"
	COL_SERVER_INFO = "server_info"
	COL_ACCOUNT     = "account"
	COL_ACCOUNT_LOG = "account_log"
	COL_WHITE       = "white"
	COL_CHANNEL     = "channel"
)

func SetDBBySid(SID int64) {
	DB = fmt.Sprintf("%v%d", DB_SPACE_WAR_GLOBAL, SID)
}

func SetDBByName(dbName string) {
	DB = dbName
}

func GetDB() string {
	return DB
}

func CreateIndex() {
	AccountCreateIndex()
	AccountLogCreateIndex()
}
