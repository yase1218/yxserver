package iface

type ITableDb interface {
	Init(path string)
	Patch()
	Load(tdb ITableDb) error
	GetConf() any
	CheckConf(c any) error
}
