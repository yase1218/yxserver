package tdb

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"io"
	"kernel/tools"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

type Tdb struct {
	*TableBase
	TableDbPath string

	Ver         string
	FileModTime map[string]int64
	FileInfos   []FileInfo

	MonsterRefreshMap map[int][]*StageMonsterRefreshCfg
}

var tdb *Tdb

func init() {
	tdb = &Tdb{
		TableBase:   &TableBase{},
		FileModTime: make(map[string]int64),
	}
}

func Init(path string) error {
	if path == "" {
		return errors.New("tableDbPath is empty")
	}

	tdb.Init(path)

	return tdb.Load(nil)
}

func (t *Tdb) Init(tableDbPath string) {
	tdb.TableDbPath = tableDbPath
}

func (t *Tdb) Load(tdb iface.ITableDb) (err error) {
	t.FileInfos = fileInfos

	otherReg := []any{
		//InitConf{},
	}
	for _, a := range otherReg {
		gob.Register(a)
	}
	for _, info := range t.FileInfos {
		for _, i := range info.SheetInfos {
			gob.Register(i.ObjPropType)
		}
	}

	f, err := os.Stat(t.TableDbPath)
	if err != nil {
		return err
	}

	dat := "tdb.dat"
	isCheck := false
	if f.IsDir() {
		if err = t.load(t.TableDbPath, dat); err != nil {
			return err
		}
		isCheck = true
	} else {
		if err = t.loadDat(t.TableDbPath); err != nil {
			return err
		}
	}

	t.Patch()

	if isCheck {
		err = t.Check()
		if err != nil {
			return err
		}
	}
	return
}

func (t *Tdb) load(path string, dat string) error {
	// 如果存在dat文件,则先载入.然后对比修改时间,重新加载时间不一致的文件
	datFile := filepath.Join(path, dat)
	f, err := os.Stat(datFile)
	if err == nil && !f.IsDir() {
		if err = t.loadDat(datFile); err != nil {
			log.Error("load dat file err", zap.String("datFile", datFile), zap.Error(err))
			t.Init(datFile)
		}
	}

	change, err := t.loadExcel(filepath.Join(path, "excel"))
	if err != nil {
		return err
	}

	if change {
		t.Ver = time.Now().Format("060102150405")
		t.genDat(datFile)
	}
	return nil
}

func (t *Tdb) loadDat(dat string) error {
	//now := time.Now()
	//defer func() {
	//	fmt.Println("loadGod use time", time.Since(now).Seconds())
	//}()

	f, err := os.Open(dat)
	if err != nil && err != io.EOF {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	dec := gob.NewDecoder(r)
	return dec.Decode(tdb)
}

func (t *Tdb) loadExcel(baseDir string) (bool, error) {
	//now := time.Now()
	//defer func() {
	//	fmt.Println("load all use time", time.Since(now).Seconds())
	//}()

	var wg sync.WaitGroup
	var loadErr error
	var num int
	for _, v := range t.FileInfos {
		fName := filepath.Join(baseDir, v.FileName)
		fInfo, err := os.Stat(fName)
		if err != nil {
			loadErr = errors.New("filename not found, " + err.Error())
			break
		}
		fModT := fInfo.ModTime().UnixNano()
		if t.FileModTime[v.FileName] == fModT {
			//fmt.Printf("文件 %v 未修改!\n", fileInfos[i].FileName)
			continue
		}
		t.FileModTime[v.FileName] = fModT
		num++

		wg.Add(1)
		go tools.GoSafe("tabledb load file", func() {
			if err = t.loadFile(fName, v.SheetInfos); err != nil {
				log.Error("load err", zap.String("fileName", v.FileName), zap.Error(err))
				loadErr = err
			}
			//fmt.Printf("加载完成: %v\n", fileInfos[index].FileName)
			wg.Done()
		})
	}

	wg.Wait()
	if loadErr != nil {
		return false, loadErr
	}
	//fmt.Printf("共加载: %v 个文件", num)
	if num == 0 {
		return false, nil
	}
	//fmt.Println("All TableDb rem cal is ok")
	return true, loadErr
}

func (t *Tdb) genDat(dat string) error {
	//now := time.Now()
	//defer func() {
	//	fmt.Println("generateGob use time", time.Since(now).Seconds())
	//}()

	f, err := os.Create(dat)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	if err = enc.Encode(tdb); err != nil {
		log.Error("encode err", zap.Error(err))
		return err
	}

	return w.Flush()
}

func (t *Tdb) GetConf() any {
	return nil
}

func (t *Tdb) loadFile(filename string, sheetInfos []SheetInfo) error {
	xlsFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	for _, sheetInfo := range sheetInfos {
		sheet, ok := xlsFile.Sheet[sheetInfo.SheetName]
		if !ok {
			log.Error("sheet not found", zap.String("tableName", sheetInfo.SheetName))
			return fmt.Errorf("no %s sheet found", sheetInfo.SheetName)
		}
		objPropType := reflect.New(reflect.TypeOf(sheetInfo.ObjPropType)).Interface()

		objs, err := ReadXlsxSheet(sheet, objPropType, 1, 1, nil)
		if err != nil {
			log.Error("read sheet err", zap.String("tableName", sheetInfo.SheetName), zap.Error(err))
			return err
		}
		if err = sheetInfo.Initer(tdb, objs); err != nil {
			log.Error("sheet init err", zap.String("tableName", sheetInfo.SheetName), zap.Error(err))
			return err
		}
	}
	return nil
}
