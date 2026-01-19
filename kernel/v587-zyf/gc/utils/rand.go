package utils

import (
	"math"
	"math/rand"
	"time"
)

/**
*随机权重
*randData map[int]int{索引:权重，索引：权重}
*return 索引
 */
func RandWeightByMap(randData map[int]int) int {
	sum := 0
	for _, v := range randData {
		sum += v
	}
	if sum <= 0 {
		return -1
	}
	randNum := rand.Intn(sum)
	count := 0
	for k, v := range randData {
		count += v
		if randNum < count {
			return k
		}
	}
	return -1
}

// 随机选择一个 GoodsWeight 的索引
func RandomWeightIndex(weights []int) int {
	total := 0
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return -1
	}
	r := rand.Intn(total) + 1
	sum := 0
	for i, w := range weights {
		if w <= 0 {
			continue
		}
		sum += w
		if r <= sum {
			return i
		}
	}
	return len(weights) - 1
}

func RandomWeightIndexCanNoResult(weights []int) int {
	if len(weights) == 0 {
		return -1
	}

	// 生成一次随机数 1-100
	r := rand.Intn(100) + 1

	// 直接遍历原数组，检查权重
	for i, weight := range weights {
		if weight > 0 && weight <= 100 && r <= weight {
			return i
		}
	}

	return -1
}

/*
 * IntWeightedRandom
 *  @Description: 	加权随机选择（int类型）
 *  @param weights	key=选项, value=权重(非负)
 *  @return int		选中的key, 权重是否有效	-1为空
 */
func IntWeightedRandom(weights map[int]int) int {
	if len(weights) == 0 {
		return -1
	}

	total := 0
	valid := make(map[int]int)
	for k, w := range weights {
		if w > 0 {
			total += w
			valid[k] = w
		}
	}

	if total == 0 {
		keys := make([]int, 0, len(weights))
		for k := range weights {
			keys = append(keys, k)
		}
		return keys[rand.Intn(len(keys))]
	}

	r := rand.Intn(total)
	for k, w := range valid {
		r -= w
		if r < 0 {
			return k
		}
	}

	for k := range valid {
		return k
	}
	return -1
}

func RandSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RoundFloat 使用 math.Round 将浮点数 f 四舍五入到小数点后 n 位。
func RoundFloat(f float64, n int) float64 {
	if f < 0 {
		return f
	}
	shift := math.Pow(10, float64(n))
	// 将浮点数乘以10的n次方，四舍五入到最近的整数，然后再除以10的n次方。
	return math.Round(f*shift) / shift
}
func RoundUp(num float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Ceil(num*shift) / shift
}

// 包含上下限 [min, max]
func RandomWithAll(min, max int) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(rand.Intn(max-min+1) + min)
}

// 不包含上限 [min, max)
func RandomWithMin(min, max int) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(rand.Intn(max-min) + min)
}

// 不包含下限 (min, max]
func RandomWithMax(min, max int) int64 {
	var res int64
	rand.New(rand.NewSource(time.Now().UnixNano()))
Restart:
	res = int64(rand.Intn(max-min+1) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}

// 都不包含 (min, max)
func RandomWithNo(min, max int) int64 {
	var res int64
	rand.New(rand.NewSource(time.Now().UnixNano()))
Restart:
	res = int64(rand.Intn(max-min) + min)
	if res == int64(min) {
		goto Restart
	}
	return res
}
