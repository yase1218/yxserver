package utils

import (
	"reflect"
	"strings"
)

// 校验结构体的字段是否为空
func ValidateColumn(src any, fields []string) (res bool) {
	if len(fields) < 1 {
		return
	}

	fieldMap := make(map[string]struct{})
	for _, v := range fields {
		fieldMap[strings.ToLower(v)] = struct{}{}
	}

	structValue := reflect.ValueOf(src).Elem() // 获取字段值
	structType := reflect.TypeOf(src).Elem()   // 获取字段类型
	for i := 0; i < structType.NumField(); i++ {
		n := strings.ToLower(structType.Field(i).Tag.Get("validate"))
		if _, ok := fieldMap[n]; !ok {
			continue
		}
		if structValue.Field(i).IsZero() {
			return true
		}
	}
	return false
}
