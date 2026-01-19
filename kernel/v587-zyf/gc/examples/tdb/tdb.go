package tdb

import (
	"encoding/gob"
	"github.com/v587-zyf/gc/iface"
	c "github.com/v587-zyf/gc/tabledb"
)

type Tdb struct {
	*c.TableDb
	*TableBase
	TableDbPath string
}

var tdb *Tdb

func init() {
	tdb = &Tdb{
		TableBase: &TableBase{},
		TableDb: &c.TableDb{
			FileModTime: make(map[string]int64),
			InitConf:    &InitConf{},
		},
	}
}

func Init(tableDbPath string) error {
	tdb.Init(tableDbPath)

	return tdb.Load(tdb)
}

func (t *Tdb) Load(tdb iface.ITableDb) (err error) {
	t.TableDb.FileInfos = fileInfos

	otherReg := []any{
		InitConf{},
	}
	for _, a := range otherReg {
		gob.Register(a)
	}
	//gob.Register(TextErrorTextCfg{})
	gob.Register(c.GlobalBaseCfg{})
	for _, info := range t.FileInfos {
		for _, i := range info.SheetInfos {
			gob.Register(i.ObjPropType)
		}
	}

	err = t.TableDb.Load(tdb)
	if err != nil {
		return err
	}

	err = t.CheckConf(t.InitConf)

	t.Patch()

	return
}
func (t *Tdb) CheckConf(c any) (err error) {
	return t.TableDb.CheckConf(tdb)
}

func Conf() *InitConf {
	if conf, ok := tdb.InitConf.(InitConf); ok {
		return &conf
	} else {
		return tdb.InitConf.(*InitConf)
	}
}

func Get() *Tdb {
	return tdb
}
