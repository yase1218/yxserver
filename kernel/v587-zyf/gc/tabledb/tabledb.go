package tabledb

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/v587-zyf/gc/errcode"
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

type TableDb struct {
	TableDbPath string

	Ver         string
	FileModTime map[string]int64

	FileInfos []FileInfo
	InitConf  any
}

var iTdb iface.ITableDb

func (t *TableDb) Init(path string) {
	t.TableDbPath = path
}

func (t *TableDb) Patch() {}

func (t *TableDb) Load(tdb iface.ITableDb) error {
	iTdb = tdb

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

	return nil
}

func (t *TableDb) load(path string, dat string) error {
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

func (t *TableDb) loadDat(dat string) error {
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
	return dec.Decode(iTdb)
}

func (t *TableDb) loadExcel(baseDir string) (bool, error) {
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

func (t *TableDb) genDat(dat string) error {
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
	if err = enc.Encode(iTdb); err != nil {
		log.Error("encode err", zap.Error(err))
		return err
	}

	return w.Flush()
}

func (t *TableDb) GetConf() any {
	return t.InitConf
}

func (t *TableDb) loadFile(filename string, sheetInfos []SheetInfo) error {
	xlsFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	for _, sheetInfo := range sheetInfos {
		sheet, ok := xlsFile.Sheet[sheetInfo.SheetName]
		if !ok {
			return fmt.Errorf("no %s sheet found", sheetInfo.SheetName)
		}
		objPropType := reflect.New(reflect.TypeOf(sheetInfo.ObjPropType)).Interface()

		objs, err := ReadXlsxSheet(sheet, objPropType, 0, 2, nil)
		if err != nil {
			return err
		}
		if err = sheetInfo.Initer(iTdb, objs); err != nil {
			return err
		}
	}
	return nil
}

var errSlice = make([]errcode.ErrCode, 0)

//var TextMap = make(map[string]string)

func InitError(code int32, codeLang ...errcode.CodeLang) errcode.ErrCode {
	e := errcode.CreateErrCode(code, codeLang...)
	errSlice = append(errSlice, e)
	return e
}

//func codeTextSign(constName, text string) string {
//	return text
//}

func LoadGlobalConf(tableDb iface.ITableDb, objs []any) error {
	tableConfs := make(map[string]*GlobalBaseCfg)
	for _, obj := range objs {
		game := obj.(*GlobalBaseCfg)
		if _, ok := tableConfs[game.Name]; ok {
			return errors.New(fmt.Sprintf("tableconf key:%d namd:%s 重复了", game.Id, game.Name))
		}
		tableConfs[game.Name] = game
	}
	return DecodeConfValues(tableDb.GetConf(), tableConfs)
}
