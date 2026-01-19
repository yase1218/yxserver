package mysql

import "time"

const TestTableName = "test"

type TestModel struct {
	ModelBase
	Time time.Time `gorm:"time"`
}

func (TestModel) TableName() string {
	return TestTableName
}
