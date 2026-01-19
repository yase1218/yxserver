package tabledb

import (
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"reflect"
)

func MapLoader(fieldName string, keyFieldName string) func(iface.ITableDb, []interface{}) error {
	return func(TableDb iface.ITableDb, objs []interface{}) error {
		fieldV := reflect.ValueOf(TableDb).Elem().FieldByName(fieldName)
		if fieldV.Kind() != reflect.Map {
			return fmt.Errorf("field %s is not a map", fieldName)
		}

		if fieldV.IsNil() || fieldV.Len() > 0 {
			fieldV.Set(reflect.MakeMap(fieldV.Type()))
		}
		for _, obj := range objs {
			objV := reflect.ValueOf(obj)
			keyFieldV := objV.Elem().FieldByName(keyFieldName)
			if !keyFieldV.IsValid() {
				return fmt.Errorf("key field %s wrong filedV:%v, when setting %s\n", keyFieldName, fieldName, keyFieldV)
			}
			if keyFieldV.Kind() != reflect.Int {
				fmt.Printf("key field %s wrong filedV:%v, when setting %s\n", keyFieldName, fieldName, keyFieldV)
				continue
			}
			if fieldV.MapIndex(keyFieldV).IsValid() {
				//return fmt.Errorf(" >%v<. The value of field >%s< in sheet >%s< is duplicate",
				return fmt.Errorf("表 %s 列 %s 值->%v 重复了",
					fieldName, keyFieldName, keyFieldV)
			}
			fieldV.SetMapIndex(keyFieldV, objV)
		}
		return nil
	}
}
