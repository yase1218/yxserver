package tdb

import (
	"fmt"
	"github.com/v587-zyf/gc/iface"
	zlog "github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"reflect"
)

func MapLoader(fieldName string, keyFieldName string) func(iface.ITableDb, []any) error {
	return func(TableDb iface.ITableDb, objs []any) error {
		// 检查 TableDb 是否为 nil
		if TableDb == nil {
			return fmt.Errorf("table db is nil")
		}

		// 获取字段值
		fieldV := reflect.ValueOf(TableDb).Elem().FieldByName(fieldName)
		if !fieldV.IsValid() {
			return fmt.Errorf("field %s does not exist", fieldName)
		}
		if fieldV.Kind() != reflect.Map {
			return fmt.Errorf("field %s is not a map", fieldName)
		}

		// 检查 map 是否为空或需要初始化
		if fieldV.IsNil() || fieldV.Len() > 0 {
			fieldV.Set(reflect.MakeMap(fieldV.Type()))
		}

		//// 获取 map 的键类型
		//mapKeyType := fieldV.Type().Key()
		//if mapKeyType.Kind() != reflect.Int {
		//	return fmt.Errorf("map field %s requires an integer key, but got %s", fieldName, mapKeyType.Kind())
		//}
		//
		//// 处理空对象列表的情况
		//if len(objs) == 0 {
		//	return nil
		//}

		for _, obj := range objs {
			objV := reflect.ValueOf(obj)
			if objV.Kind() != reflect.Ptr || objV.Elem().Kind() != reflect.Struct {
				zlog.Error("invalid object type", zap.Any("obj", obj))
				return fmt.Errorf("object must be a pointer to struct, but got %v", objV.Kind())
			}

			keyFieldV := objV.Elem().FieldByName(keyFieldName)
			if !keyFieldV.IsValid() {
				zlog.Error("key field does not exist", zap.String("keyFieldName", keyFieldName), zap.Any("obj", obj))
				return fmt.Errorf("key field %s does not exist in object", keyFieldName)
			}

			if keyFieldV.Kind() != reflect.Int {
				zlog.Error("key field type mismatch", zap.String("keyFieldName", keyFieldName), zap.Any("keyFieldV", keyFieldV))
				continue
			}

			// 检查是否重复
			if fieldV.MapIndex(keyFieldV).IsValid() {
				zlog.Error("值重复!!!", zap.String("fieldName", fieldName), zap.String("keyFieldName", keyFieldName), zap.Int64("keyFieldV", keyFieldV.Int()))
				continue
				return fmt.Errorf("duplicate value for key field %s: %v", keyFieldName, keyFieldV.Interface())
			}

			// 设置 map 值
			fieldV.SetMapIndex(keyFieldV, objV)
		}

		return nil
	}
}
