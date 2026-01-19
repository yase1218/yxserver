package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GetAllFieldsAsString(obj any) string {
	return GetAllFieldsAsStringWithTableName(obj, "")
}
func GetAllFieldsAsStringWithTableName(obj any, tableName string) string {
	objT := reflect.TypeOf(obj)
	var fields []string
	for i := 0; i < objT.NumField(); i++ {
		fieldT := objT.Field(i)
		tag := fieldT.Tag.Get("db")
		if tag == "" {
			continue
		}
		oneFileName := fmt.Sprintf("`%s`", tag)
		if tableName != "" {
			oneFileName = fmt.Sprintf("%s.`%s`", tableName, tag)
		}
		fields = append(fields, oneFileName)
	}
	return strings.Join(fields, ",")
}
