package db_server

import "fmt"

var (
	DB string
)

const (
	DB_SPACE_WAR = "space_war_local%d"
)

const (
	COL_COUNTER = "counter"
	COL_USER    = "user"
)

func SetDB(SID int64) {
	DB = fmt.Sprintf(DB_SPACE_WAR, SID)
}

func GetDB() string {
	return DB
}
