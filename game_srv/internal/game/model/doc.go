package model

import (
	"context"
	"gameserver/internal/publicconst"
)

func GetDBCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), publicconst.DB_OP_TIME_OUT)
}
