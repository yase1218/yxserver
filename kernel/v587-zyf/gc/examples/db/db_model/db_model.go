package db_model

import (
	"fmt"
	"github.com/v587-zyf/gc/db/db_model"
	"github.com/v587-zyf/gc/utils"
	"gopkg.in/gorp.v1"
	"time"
)

type Test struct {
	db_model.BaseTable
	Id   int       `db:"id" orm:"pk;auto"`
	Time time.Time `db:"time" orm:"comment(时间)"`
}

func (t *Test) TableName() string {
	return "test"
}

type TestModel struct {
	db_model.DBBaseModel
}

var (
	testDb     = new(Test)
	testModel  = new(TestModel)
	testIdSeq  = &IdSeq{Name: "test"}
	TestFields = utils.GetAllFieldsAsString(Test{})
)

func init() {
	db_model.Register(DB_TEST, testModel, func(dbMap *gorp.DbMap) {
		dbMap.AddTableWithName(Test{}, testDb.TableName()).SetKeys(true, "id")
	})
}

func GetTestModel() *TestModel {
	return testModel
}

func (m *TestModel) Create(data *Test) error {
	return m.DbMap().Insert(data)
}

func (m *TestModel) Update(data *Test) (int, error) {
	count, err := m.DbMap().Update(data)
	if err != nil {
		fmt.Println("err", err)
	}
	return int(count), err
}

func (m *TestModel) GetOne(id int) (*Test, error) {
	var data Test
	err := m.DbMap().SelectOne(&data, fmt.Sprintf("select %s from %s where id = ?", TestFields, testDb.TableName()), id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
func (m *TestModel) GetSeqId() (int, error) {
	return testIdSeq.Next()
}
