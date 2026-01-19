package tools

import (
	"encoding/binary"
	"math/rand"
	"strconv"
	"strings"

	"github.com/zy/game_data/template"
)

// FilterString 过滤字符串
func FilterString(str string) string {
	str = strings.Replace(str, "drop", "*", -1)
	str = strings.Replace(str, "remove", "*", -1)
	return str
}

func Str2UInt32(str string) uint32 {
	temp, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return uint32(temp)
}

func ListRemoveRepeat(data []uint32) []uint32 {
	temp := make(map[uint32]interface{})
	var ret []uint32
	for i := 0; i < len(data); i++ {
		temp[data[i]] = struct{}{}
	}
	for id, _ := range temp {
		ret = append(ret, id)
	}
	return ret
}

func LimitUint32(a, maxValue uint32) uint32 {
	if a > maxValue {
		return maxValue
	}
	return a
}

func MergeToMapItem(data map[uint32]uint32, list []*template.SimpleItem) {
	for k := 0; k < len(list); k++ {
		if count, ok := data[list[k].ItemId]; ok {
			data[list[k].ItemId] = count + list[k].ItemNum
		} else {
			data[list[k].ItemId] = list[k].ItemNum
		}
	}
}

// RandomPercentage 随机一个百分比 返回 [0, 100]
func RandomPercentage(r *rand.Rand) uint32 {
	return uint32(r.Int31n(101))
}

func Uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, num)
	return bytes
}

func BytesToUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

func FindMaxElemForArray(data []int) int {
	maxElem := 0
	for i := range data {
		if data[i] > maxElem {
			maxElem = data[i]
		}
	}

	return maxElem
}
