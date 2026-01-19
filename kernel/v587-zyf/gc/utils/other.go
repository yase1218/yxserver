package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/v587-zyf/gc/iface"
	"reflect"
	"unsafe"
)

// IsZero 零值判断
func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func ToBinaryNumber[T iface.Integer](n T) T {
	var x T = 1
	for x < n {
		x *= 2
	}
	return x
}

// MethodExists
// if nil return false
func MethodExists(in any, method string) (reflect.Value, bool) {
	if in == nil || method == "" {
		return reflect.Value{}, false
	}
	p := reflect.TypeOf(in)
	if p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	if p.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	object := reflect.ValueOf(in)
	newMethod := object.MethodByName(method)
	if !newMethod.IsValid() {
		return reflect.Value{}, false
	}
	return newMethod, true
}

// MaskXOR 计算掩码
func MaskXOR(b []byte, key []byte) {
	var maskKey = binary.LittleEndian.Uint32(key)
	var key64 = uint64(maskKey)<<32 + uint64(maskKey)

	for len(b) >= 64 {
		v := binary.LittleEndian.Uint64(b)
		binary.LittleEndian.PutUint64(b, v^key64)
		v = binary.LittleEndian.Uint64(b[8:16])
		binary.LittleEndian.PutUint64(b[8:16], v^key64)
		v = binary.LittleEndian.Uint64(b[16:24])
		binary.LittleEndian.PutUint64(b[16:24], v^key64)
		v = binary.LittleEndian.Uint64(b[24:32])
		binary.LittleEndian.PutUint64(b[24:32], v^key64)
		v = binary.LittleEndian.Uint64(b[32:40])
		binary.LittleEndian.PutUint64(b[32:40], v^key64)
		v = binary.LittleEndian.Uint64(b[40:48])
		binary.LittleEndian.PutUint64(b[40:48], v^key64)
		v = binary.LittleEndian.Uint64(b[48:56])
		binary.LittleEndian.PutUint64(b[48:56], v^key64)
		v = binary.LittleEndian.Uint64(b[56:64])
		binary.LittleEndian.PutUint64(b[56:64], v^key64)
		b = b[64:]
	}

	for len(b) >= 8 {
		v := binary.LittleEndian.Uint64(b[:8])
		binary.LittleEndian.PutUint64(b[:8], v^key64)
		b = b[8:]
	}

	var n = len(b)
	for i := 0; i < n; i++ {
		idx := i & 3
		b[i] ^= key[idx]
	}
}

// BufferReset 重置buffer底层的切片
// Reset the buffer's underlying slice
// 注意：修改后面的属性一定要加偏移量，否则可能会导致未定义的行为。
// Note: Be sure to add an offset when modifying the following properties, otherwise it may lead to undefined behavior.
func BufferReset(b *bytes.Buffer, p []byte) { *(*[]byte)(unsafe.Pointer(b)) = p }

// WithDefault 如果原值为零值, 返回新值, 否则返回原值
func WithDefault[T comparable](rawValue, newValue T) T {
	if IsZero(rawValue) {
		return newValue
	}
	return rawValue
}

func SelectValue[T any](ok bool, a, b T) T {
	if ok {
		return a
	}
	return b
}
