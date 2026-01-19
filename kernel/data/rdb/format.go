package rdb

import (
	"fmt"
	"kernel/kenum"
)

func FormatUserID(userID, serverId any) string {
	return fmt.Sprintf(kenum.Redis_Key_User, userID, serverId)
}

// {UserID:1234567890}HttpLogin
func FormatUserLogin(userID, serverId any) string {
	return fmt.Sprint(FormatUserID(userID, serverId), "HttpLogin")
}

// {UserDump}
func FormatUserDumpKey() string {
	return fmt.Sprintf("{UserDump}")
}
