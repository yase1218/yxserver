package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Float64Slice []float64
type IntSlice []int
type StringSlice []string
type IntMap map[int]int

func IntSliceFromString(str string, sep string) (IntSlice, error) {
	if len(str) == 0 {
		return IntSlice(make([]int, 0)), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make(IntSlice, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.Atoi(strs[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Float64SliceFromString(str string, sep string) (Float64Slice, error) {
	if len(str) == 0 {
		return Float64Slice(make([]float64, 0)), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make(Float64Slice, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.ParseFloat(strs[i], 64)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s IntSlice) Index(element int) int {
	for i, v := range s {
		if v == element {
			return i
		}
	}
	return -1
}
func (s IntSlice) RemoveIndex(index int) IntSlice {
	if index < 0 || index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}
func (s IntSlice) RemoveElement(element int) IntSlice {
	for i, v := range s {
		if v == element {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
func (s IntSlice) Add(element int) IntSlice {
	return append(s, element)
}
func (s IntSlice) AddUnique(element int) IntSlice {
	if s.Index(element) < 0 {
		return s
	}
	return append(s, element)
}
func (s IntSlice) String(sep string) string {
	var arrStr = make([]string, len(s))
	for i, v := range s {
		arrStr[i] = strconv.Itoa(v)
	}
	return strings.Join(arrStr, sep)
}
func (s IntSlice) Len() int {
	return len(s)
}
func (s IntSlice) Less(i, j int) bool {
	if s[j] != s[i] {
		return s[j] > s[i]
	}
	return false
}
func (s IntSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func ConvertInt32SliceToIntSlice(origin []int32) []int {
	ret := make([]int, len(origin))
	for i, v := range origin {
		ret[i] = int(v)
	}
	return ret
}

func ConvertUInt32SliceToIntSlice(origin []uint32) []int {
	ret := make([]int, len(origin))
	for i, v := range origin {
		ret[i] = int(v)
	}
	return ret
}

func ConvertIntSlice2Int32Slice(origin []int) []int32 {
	ret := make([]int32, len(origin))
	for i, v := range origin {
		ret[i] = int32(v)
	}
	return ret
}

func ConvertIntSlice2UInt32Slice(origin []int) []uint32 {
	ret := make([]uint32, len(origin))
	for i, v := range origin {
		ret[i] = uint32(v)
	}
	return ret
}

func ConvertMapIntToInt32(origin map[int]int) map[int32]int32 {
	ret := make(map[int32]int32, len(origin))
	for i, v := range origin {
		ret[int32(i)] = int32(v)
	}
	return ret
}

func ConvertMapIntToUInt32(origin map[int]int) map[uint32]uint32 {
	ret := make(map[uint32]uint32, len(origin))
	for i, v := range origin {
		ret[uint32(i)] = uint32(v)
	}
	return ret
}

func ConvertMapInt32ToInt(origin map[int32]int32) map[int]int {
	ret := make(map[int]int, len(origin))
	for i, v := range origin {
		ret[int(i)] = int(v)
	}
	return ret
}

func ConvertMapUInt32ToInt(origin map[uint32]uint32) map[int]int {
	ret := make(map[int]int, len(origin))
	for i, v := range origin {
		ret[int(i)] = int(v)
	}
	return ret
}

func SliceIntUnique(origin []int) []int {
	ret := make([]int, 0)
	tempMap := make(map[int]struct{})
	for _, v := range origin {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

func SliceInt32Unique(arr []int32) []int32 {
	ret := make([]int32, 0)
	tempMap := make(map[int32]struct{})
	for _, v := range arr {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

func SliceUInt32Unique(arr []uint32) []uint32 {
	ret := make([]uint32, 0)
	tempMap := make(map[uint32]struct{})
	for _, v := range arr {
		if _, ok := tempMap[v]; ok {
			continue
		}
		ret = append(ret, v)
		tempMap[v] = struct{}{}
	}
	return ret
}

// 2维int转string
func SliceInt2ToString(arr [][]int, sep1 string, sep2 string) string {
	slice1 := make([]string, len(arr))
	for k, v := range arr {
		slice1[k] = JoinIntSlice(v, sep1)
	}
	return strings.Join(slice1, sep2)
}

// 2维int转1维string
func SliceInt2ToSliceString1(arr [][]int, sep string) []string {
	slice1 := make([]string, len(arr))
	for k, v := range arr {
		slice1[k] = JoinIntSlice(v, sep)
	}
	return slice1
}

func JoinIntSlice(a []int, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func JoinInt32Slice(a []int32, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, sep)
}

func JoinUInt32Slice(a []uint32, sep string) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range a {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, sep)
}

// anySlice2StringSlice 将anySlice转换为stringSlice
func AnySlice2StringSlice(arr []any) []string {
	strArr := make([]string, len(arr))
	for index, v := range arr {
		strArr[index] = fmt.Sprint(v)
	}
	return strArr
}

// IntMap2ToString 将IntMap转换为字符串
// 1,2;3,4;
func IntMap2ToString(arr IntMap) string {
	slice1 := make([]string, 0, len(arr))
	for k, v := range arr {
		slice1 = append(slice1, fmt.Sprintf("%d,%d", k, v))
	}
	return strings.Join(slice1, ";")
}

// InCollection 判断字符串 elem 是否在字符串切片 elems 中
// Checks if the given string elem is in the string slice elems.
func InCollection(elem string, elems []string) bool {
	for _, item := range elems {
		if item == elem {
			return true
		}
	}
	return false
}

// GetIntersectionElem 获取两个字符串切片 a 和 b 的交集中的一个元素
// Gets an element in the intersection of two string slices a and b
func GetIntersectionElem(a, b []string) string {
	for _, item := range a {
		if InCollection(item, b) {
			return item
		}
	}
	return ""
}

// Split 分割给定的字符串 s，使用 sep 作为分隔符。空值将会被过滤掉。
// Splits the given string s using sep as the separator. Empty values will be filtered out.
func Split(s string, sep string) []string {
	var list = strings.Split(s, sep)
	var j = 0
	for _, v := range list {
		if v = strings.TrimSpace(v); v != "" {
			list[j] = v
			j++
		}
	}
	return list[:j]
}
