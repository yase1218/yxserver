package rbi

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// join 工具函数
func join(strs []string, sep string) string {
	var res []string
	for _, s := range strs {
		if s != "" {
			res = append(res, s)
		}
	}
	if len(res) == 0 {
		return ""
	}
	return strings.Join(res, sep)
}

// GetFieldValue 获取字段值，返回字符串形式，并处理空值为默认值
func GetFieldValue(field reflect.Value, fieldType reflect.Type) string {
	switch field.Kind() {
	case reflect.String:
		val := field.String()
		if val == "" {
			return "0"
		}
		return val

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := field.Int()
		if val == 0 {
			return "0"
		}
		return fmt.Sprintf("%d", val)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := field.Uint()
		if val == 0 {
			return "0"
		}
		return fmt.Sprintf("%d", val)

	case reflect.Float32, reflect.Float64:
		val := field.Float()
		if val == 0.0 {
			return "0"
		}
		return fmt.Sprintf("%.2f", val)

	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time)
			if t.IsZero() {
				return "0"
			}
			return t.Format("2006-01-02 15:04:05")
		}
		var nestedValues []string
		nestedV := field
		nestedT := fieldType
		for i := 0; i < nestedV.NumField(); i++ {
			f := nestedT.Field(i)
			val := nestedV.Field(i)
			nestedValues = append(nestedValues, GetFieldValue(val, f.Type))
		}
		return join(nestedValues, "|") // 直接合并进主列表，不带结构体名
	default:
		return ""
	}
}

// StructToPipeString 将任意 struct 转为 | 分隔的字符串，格式为：StructName|field1|field2|...
func StructToPipeString(s interface{}) string {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	structType := v.Type()
	structName := structType.Name()

	var values []string
	values = append(values, structName)

	for i := 0; i < v.NumField(); i++ {
		field := structType.Field(i)
		value := v.Field(i)

		strVal := GetFieldValue(value, field.Type)
		values = append(values, strVal)
	}

	return join(values, "|")
}

func encodeGob(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func decodeGob(data []byte, target interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(target)
}
