package tools

import (
	"errors"
	"math/rand"
)

// ListContain 是否存在指定的值
func ListContain(data []uint32, target uint32) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			return true
		}
	}
	return false
}

// ListContain 是否存在指定的值
func ListContain64(data []uint64, target uint64) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			return true
		}
	}
	return false
}

// 移除指定下标的元素
func ListIntRemoveIndex(data []int, index int) []int {
	if index < 0 || index >= len(data) {
		return data
	}

	return append(data[:index], data[index+1:]...)
}

func ListIntContain(data []int, target int) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			return true
		}
	}
	return false
}

func ListStrContain(data []string, target string) bool {
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			return true
		}
	}
	return false
}

// ListRemove
func ListRemove(data []uint32, target uint32) []uint32 {
	pos := -1
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			pos = i
			break
		}
	}

	if pos != -1 {
		data = append(data[:pos], data[pos+1:]...)
	}
	return data
}

// ListRemove
func ListIntRemove(data []int, target int) []int {
	pos := -1
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			pos = i
			break
		}
	}

	if pos != -1 {
		data = append(data[:pos], data[pos+1:]...)
	}
	return data
}

// SetBit 设置位
func SetBit(value *uint32, pos uint32) {
	*value |= 1 << pos
}

// GetBit 获得位置
func GetBit(value uint32, pos uint32) bool {
	if value&(1<<pos) > 0 {
		return true
	}
	return false
}

func SetBitFlag(val *uint32, pos uint32, flag uint32) {
	*val |= flag << pos
}

func GetBitFlag(val uint32, pos uint32, flag uint32) bool {
	if val&(flag<<pos) > 0 {
		return true
	}
	return false
}

// ListUint32Equal 两个uint32 list 是否相等
func ListUint32Equal(l1 []uint32, l2 []uint32) bool {
	if len(l1) != len(l2) {
		return false
	}

	for i := 0; i < len(l1); i++ {
		if l1[i] != l2[i] {
			return false
		}
	}

	return true
}

func ListUint32AddNoRepeat(l []uint32, target uint32) []uint32 {
	exist := false
	for i := 0; i < len(l); i++ {
		if l[i] == target {
			exist = true
			break
		}
	}
	if !exist {
		l = append(l, target)
	}
	return l
}

func ListIntAddNoRepeat(l []int, target int) []int {
	exist := false
	for i := 0; i < len(l); i++ {
		if l[i] == target {
			exist = true
			break
		}
	}
	if !exist {
		l = append(l, target)
	}
	return l
}

func ListUint32AddNoRepeats(l []uint32, target []uint32) []uint32 {
	var ret []uint32
	ret = l
	for i := 0; i < len(target); i++ {
		ret = ListUint32AddNoRepeat(ret, target[i])
	}
	return ret
}

func RandArray(r *rand.Rand, weight []uint32, times uint32) (error, []int) {
	if r == nil || len(weight) == 0 || times == 0 {
		return errors.New("invalid data"), nil
	}

	if len(weight) == 1 {
		return nil, []int{0}
	}

	var totalWeight uint32 = 0
	for i := 0; i < len(weight); i++ {
		totalWeight += weight[i]
	}

	var ret []int
	for i := 0; i < int(times); i++ {
		randWeight := r.Int31n(int32(totalWeight))
		var temp uint32 = 0
		for k := 0; k < len(weight); k++ {
			temp += weight[k]
			if randWeight < int32(temp) {
				ret = append(ret, k)
				break
			}
		}
	}
	return nil, ret
}
