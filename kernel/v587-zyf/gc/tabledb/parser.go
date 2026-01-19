package tabledb

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/v587-zyf/gc/utils"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	reLineBreak = regexp.MustCompile("\n")

	lang  string        = "zh-cn"
	cn_tw map[rune]rune //简体对应繁体
)

func changeLangStr(str string) string {
	if lang == "zh-cn" {
		return str
	}
	str_rune := []rune(str)

	for i := 0; i < len(str_rune); i++ {
		if v, ok := cn_tw[str_rune[i]]; ok {
			str_rune[i] = v
		}
	}
	return string(str_rune)
}

type Decoder interface {
	Decode(str string) error
}

func ReadXlsxSheet(sheet *xlsx.Sheet, obj interface{}, startRow int, startCol int, groupFinder func(groupName string, fieldName string) (int, error)) ([]interface{}, error) {
	objT := reflect.TypeOf(obj)

	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		return nil, errors.New("readSheet must pass a struct type")
	}
	if len(sheet.Rows) <= startRow || len(sheet.Cols) <= startCol {
		return nil, errors.New("empty sheet " + sheet.Name + ",row:" + strconv.Itoa(len(sheet.Rows)) + ",col:" + strconv.Itoa(len(sheet.Cols)))
	}
	type FieldInfo struct {
		Index   int
		Field   *reflect.StructField
		Group   string
		ColName string
	}

	var maxColumnIndex = 0
	serverNeed := make(map[int]bool)
	for i, cell := range sheet.Rows[startRow].Cells {
		if i < startCol-1 {
			continue
		} else if cell == nil || len(strings.TrimSpace(cell.Value)) == 0 {
			continue
		}
		maxColumnIndex = i
		serverNeed[i] = true
	}

	objT = objT.Elem()
	var colMap = make(map[int]*FieldInfo)
	columnFound := make(map[string]bool)
	for i, cell := range sheet.Rows[startRow].Cells {
		if i < startCol-1 {
			continue
		} else if cell == nil || len(strings.TrimSpace(cell.Value)) == 0 {
			break
		}
		if !serverNeed[i] {
			continue
		}
		cellValue := strings.TrimSpace(cell.Value)
		for j := 0; j < objT.NumField(); j++ {
			field := objT.Field(j)
			if field.Tag.Get("col") == cellValue {
				colMap[i] = &FieldInfo{Index: j, Field: &field, Group: field.Tag.Get("group"), ColName: cellValue}
				columnFound[cellValue] = true
			}
		}
	}

	if len(colMap) == 0 {
		return nil, errors.New("no column found for sheet " + sheet.Name)
	}
	//检查是否缺少column配置
	for j := 0; j < objT.NumField(); j++ {
		field := objT.Field(j)
		colInStruct := field.Tag.Get("col")
		if len(colInStruct) < 1 {
			continue
		}
		if !columnFound[colInStruct] {
			//fmt.Printf("no column found for sheet %s, column %s \n", sheet.Name, colInStruct)
			return nil, errors.New(fmt.Sprintf("表格 %s,中没有列%s,更新checkconfig.exe再试试", sheet.Name, colInStruct))
		}
	}

	errFunc := func(elem reflect.Type, fieldIndex, i, j int, sheet *xlsx.Sheet, err error) error {
		return fmt.Errorf("field %s at %c%d error for sheet %s: %s", elem.Field(fieldIndex).Name, 'A'+j%26, i+1, sheet.Name, err.Error())
	}

	rowCount := len(sheet.Rows) - 3
	columnCount := 0
	sizeAll := 0
	var result = make([]interface{}, 0)
	emptyRowCount := 0
	emptRow := 0
	for i, row := range sheet.Rows {
		if i < startRow+4 {
			continue
		} else if row == nil || len(row.Cells) == 1 {
			if emptyRowCount >= 4 {
				break
			}
			emptRow = i
			emptyRowCount += 1
			continue
		}

		if emptyRowCount >= 1 {
			return nil, fmt.Errorf("错误：空行, 表[%v] blank row in %v", sheet.Name, emptRow+1)
		}

		objInstance := reflect.New(objT).Interface()
		objV := reflect.ValueOf(objInstance).Elem()

		columnCount = 0
		for j, cell := range row.Cells {
			if j < startCol-1 {
				continue
			}
			fieldInfo := colMap[j]
			if fieldInfo == nil {
				continue
			}
			cellString, err := cell.FormattedValue()
			if err != nil {
				return nil, fmt.Errorf("get column %s error for sheet %s,err:%v,cell:%v", fieldInfo.ColName, sheet.Name, err, cell)
			}
			cellString = strings.TrimSpace(cellString)
			if j == startCol && i >= startRow && (cell == nil || len(cellString) == 0) {
				goto exit //finish when meet first empty row (the first column of this row is empty)
			}
			if j > maxColumnIndex {
				break
			}
			fieldV := objV.Field(fieldInfo.Index)
			if !fieldV.CanSet() {
				return nil, fmt.Errorf("field %s can not be set for sheet %s", objT.Field(fieldInfo.Index).Name, sheet.Name)
			}
			if decoder, ok := fieldV.Addr().Interface().(Decoder); ok {
				err := decoder.Decode(cellString)
				if err != nil {
					return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
				}
				continue
			}
			if len(cellString) == 0 {
				continue
			}

			columnCount++
			if objT.Field(fieldInfo.Index).Tag.Get("client") != "" {
				sizeAll += len(cellString)
			}

			switch objT.Field(fieldInfo.Index).Type.Kind() {
			case reflect.Bool:
				if cellString == "1" {
					fieldV.SetBool(true)
				} else if cellString == "0" {
					fieldV.SetBool(false)
				} else {
					b, err := strconv.ParseBool(cellString)
					if err != nil {
						return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
					}
					fieldV.SetBool(b)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if len(fieldInfo.Group) > 0 && groupFinder != nil {
					x, err := groupFinder(fieldInfo.Group, cellString)
					if err != nil {
						return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
					}
					fieldV.SetInt(int64(x))
				} else {
					cellFloat, err := strconv.ParseFloat(cellString, 64)
					if err != nil {
						return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
					}
					fieldV.SetInt(int64(utils.RoundFloat(cellFloat, 0)))
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				x, err := strconv.ParseUint(cellString, 10, 64)
				if err != nil {
					return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
				}
				fieldV.SetUint(x)
			case reflect.Float32, reflect.Float64:
				x, err := strconv.ParseFloat(cellString, 64)
				if err != nil {
					return nil, errFunc(objT, fieldInfo.Index, i, j, sheet, err)
				}
				fieldV.SetFloat(x)
			case reflect.String:
				s1 := reLineBreak.ReplaceAllString(cellString, "")
				s1 = changeLangStr(s1)
				fieldV.SetString(strings.Replace(s1, `"`, `\"`, -1))
			default:
				columnCount--
				continue
			}

		}
		result = append(result, objInstance)
	}

exit:
	if rowCount < 1 {
		rowCount = 1
	}
	if columnCount < 1 {
		columnCount = 1
	}

	return result, nil
}
