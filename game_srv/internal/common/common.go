package common

import (
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

var (
	SnowFlake *utils.Snowflake
)

func InitSnowFlake(id int64) error {
	var err error
	SnowFlake, err = utils.NewSnowflake(id)
	if err != nil {
		log.Error("snow flake init err", zap.Error(err))
		return err
	}
	return nil
}

func GenSnowFlake() int64 {
	return int64(SnowFlake.Generate())
}

type UintPair struct {
	First  uint32
	Second uint32
}

func (u *UintPair) EqualWith(o *UintPair) bool {
	return u.First == o.First && u.Second == o.Second
}
