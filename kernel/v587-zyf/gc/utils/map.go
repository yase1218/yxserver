package utils

func ConvertIntMapToFloat64(origin map[int]int) map[int32]float64 {
	ret := make(map[int32]float64, len(origin))
	for i, v := range origin {
		ret[int32(i)] = float64(v)
	}
	return ret
}
func ConvertFloat64MapToFloat64(origin map[int]float64) map[int32]float64 {
	ret := make(map[int32]float64, len(origin))
	for i, v := range origin {
		ret[int32(i)] = v
	}
	return ret
}
