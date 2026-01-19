package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Array[T any] []T

// Scan 实现 sql.Scanner 接口
func (a *Array[T]) Scan(value any) error {
	if a == nil {
		return errors.New("Array[T]: Scan called on nil pointer")
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal json value: %v", value)
	}
	var temp []T
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*a = temp
	return nil
}

// Value 实现 driver.Valuer 接口
func (a *Array[T]) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(*a)
}

type Map[K comparable, V any] map[K]V

// Scan 实现 sql.Scanner 接口
func (m *Map[K, V]) Scan(input any) error {
	if m == nil {
		return errors.New("Map[K, V]: Scan called on nil pointer")
	}
	if input == nil {
		*m = make(Map[K, V])
		return nil
	}
	bytes, ok := input.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed: %v", input)
	}
	var temp Map[K, V]
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}
	*m = temp
	return nil
}

// Value 实现 driver.Valuer 接口
func (m *Map[K, V]) Value() (driver.Value, error) {
	if m == nil || len(*m) == 0 {
		return nil, nil
	}
	return json.Marshal(*m)
}
