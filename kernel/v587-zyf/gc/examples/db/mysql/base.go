package mysql

import (
	"gorm.io/gorm"
	"time"
)

type ModelBase struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"createAt"`
	UpdatedAt time.Time      `gorm:"updateAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" swaggertype:"string"`
	Version   int            `gorm:"column:version;type:int;default:1"`
}
