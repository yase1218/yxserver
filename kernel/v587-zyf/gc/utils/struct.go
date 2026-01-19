package utils

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// StructToValues 将结构体转换为 url.Values
func StructToValuesByKey(v any, key string) (url.Values, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct")
	}

	values := url.Values{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)

		tag := structField.Tag.Get(key)
		if tag == "" {
			continue
		}

		var value string
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(field.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.String:
			value = field.String()
		case reflect.Slice, reflect.Array:
			if reflect.TypeOf(field.Interface()).Implements(reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()) {
				text, err := field.Interface().(encoding.TextMarshaler).MarshalText()
				if err != nil {
					return nil, err
				}
				value = string(text)
			} else {
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					elemValue := fmt.Sprintf("%v", elem.Interface())
					values.Add(tag, elemValue)
				}
				continue
			}
		default:
			value = fmt.Sprintf("%v", field.Interface())
		}

		values.Add(tag, value)
	}

	return values, nil
}
