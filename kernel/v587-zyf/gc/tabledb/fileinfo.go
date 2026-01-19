package tabledb

import "github.com/v587-zyf/gc/iface"

type SheetInfo struct {
	SheetName   string
	Initer      func(iface.ITableDb, []any) error
	ObjPropType any
}
type FileInfo struct {
	FileName   string
	SheetInfos []SheetInfo
}

func (f *FileInfo) GetFileName() string {
	return f.FileName
}

func (f *FileInfo) GetSheetInfos() []SheetInfo {
	return f.SheetInfos
}
