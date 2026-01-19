package db_model

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/v587-zyf/gc/db/db_model"
	"github.com/v587-zyf/gc/utils"
	"sync"
	"time"

	"gopkg.in/gorp.v1"
)

type IdSeq struct {
	curId int
	maxId int
	Name  string
	mu    sync.Mutex
}

func (is *IdSeq) Next() (int, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	if is.curId > is.maxId || is.maxId == 0 {
		number, step, err := GetIdModel().GetIds(is.Name)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("generate new id form %s failed: %s", is.Name, err.Error()))
		}
		is.curId = number * step
		is.maxId = is.curId + step - 1
	}

	returnId := is.curId
	is.curId++

	return returnId, nil
}

type Id struct {
	Id        int       `db:"id"`
	Number    int       `db:"number"`
	Name      string    `db:"name"`
	Step      int       `db:"step"`
	CreatedAt time.Time `db:"createAt"`
	UpdateAt  time.Time `db:"updateAt"`
}

type IdModel struct {
	db_model.DBBaseModel
}

var (
	idModel  = &IdModel{}
	idFields = utils.GetAllFieldsAsString(Id{})
)

func init() {
	db_model.Register(DB_TEST, idModel, func(dbMap *gorp.DbMap) {
		dbMap.AddTableWithName(Id{}, "ids").SetKeys(true, "id")
	})
}

func GetIdModel() *IdModel {
	return idModel
}

func (m *IdModel) GetIds(name string) (number int, step int, err error) {
	if err = m.DbMap().Db.Ping(); err != nil {
		return 0, 0, err
	}

	tx, err := m.DbMap().Db.Begin()
	if err != nil {
		return 0, 0, err
	}

	row := tx.QueryRow("select number, step from ids where name = ? for update", name)
	if err = row.Scan(&number, &step); err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	number++
	if _, err = tx.Exec("update ids set number = ? where name = ?", number, name); err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	tx.Commit()

	return
}

func (m *IdModel) Create(idSeq *Id) error {
	return m.DbMap().Insert(idSeq)
}

func (m *IdModel) CheckIdsCfg(ids []string) {
	logs.Info("CheckIdsCfg  ids:%v ", ids)
	for _, idName := range ids {
		num, _, _ := m.GetIds(idName)
		if num == 0 {
			idInfo := &Id{Number: 1, Name: idName, Step: 1000, CreatedAt: time.Now(), UpdateAt: time.Now()}
			m.Create(idInfo)
		}
	}
}
