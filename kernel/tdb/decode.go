package tdb

import (
	"errors"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

const (
	SEMICOLON = ";"
	COMMA     = ","
	COLON     = ":"
	PIPE      = "|"
	SPACE     = " "
	HLINE     = "-"
)

// old
type (
	FloatSlice   []float64
	IntSlice     []int
	IntSlice2    [][]int
	StringSlice  []string
	StringSlice2 [][]string
	IntMap       map[int]int
)

// new
type (
	FloatSlice2 [][]float64
	IntFloatMap map[int]float64
	IntStrMap   map[int]string
	StrIntMap   map[string]int
	StrFloatMap map[string]float64
	StrStrMap   map[string]string
)

/***************************************************************/
/************************数据类型解析*****************************/
/***************************************************************/
func (this *IntSlice) Decode(str string) error {
	ints, err := utils.IntSliceFromString(str, COMMA)
	if err != nil {
		return err
	}
	*this = IntSlice(ints)
	return nil
}
func (this IntSlice) ToInt32Slice() []int32 {
	l := len(this)
	ret := make([]int32, l)
	if l == 0 {
		return ret
	}
	for i := 0; i < l; i++ {
		ret[i] = int32(this[i])
	}
	return ret
}
func (this IntSlice) String(sep string) string {
	var arrStr = make([]string, len(this))
	for i, v := range this {
		arrStr[i] = strconv.Itoa(v)
	}
	return strings.Join(arrStr, sep)
}

func (this *FloatSlice) Decode(str string) error {
	ints, err := utils.Float64SliceFromString(str, COMMA)
	if err != nil {
		return err
	}
	*this = FloatSlice(ints)
	return nil
}

func (this *StringSlice) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		*this = make([]string, 0)
		return nil
	}
	*this = strings.Split(str, COMMA)
	return nil
}

func (this *IntMap) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	//fmt.Printf("str = %+v\n", str)
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	*this = make(IntMap)
	for _, v := range infoList {
		list := strings.Split(strings.TrimSpace(v), COMMA)
		if len(list) != 2 {
			log.Warn("IntMap 属性信息格式错误", zap.String("v", v))
			continue
			return errors.New(v + "IntMap 属性信息格式错误")
		}

		k, err := strconv.Atoi(list[0])
		if err != nil {
			return err
		}
		if _, ok := (*this)[k]; ok {
			log.Warn("IntMap 属性重复", zap.String("v", v))
			continue
			return errors.New(v + "IntMap 属性重复")
		}
		v, err := strconv.Atoi(list[1])
		if err != nil {
			return err
		}
		(*this)[k] = v
	}
	//fmt.Printf("decode string is %s intmap:%v", str, *this)
	return nil

}
func (this *IntMap) Add(delta IntMap) {
	for k, v := range delta {
		(*this)[k] += v
	}
}
func (this IntMap) Clone() IntMap {
	ret := make(IntMap, len(this))
	for k, v := range this {
		ret[k] = v
	}
	return ret
}

func (this *IntSlice2) Decode(str string) error {
	*this = make(IntSlice2, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		ints, err := utils.IntSliceFromString(v, COMMA)
		if err != nil {
			return err
		}
		*this = append(*this, ints)
	}
	return nil
}

func (this *StringSlice2) Decode(str string) error {
	*this = make(StringSlice2, 0)
	if len(str) == 0 {
		return nil
	}
	stringSlice := strings.Split(strings.TrimSpace(str), SEMICOLON)
	for _, v := range stringSlice {
		allData := make(StringSlice, 0)
		data := strings.Split(strings.TrimSpace(v), COMMA)
		for _, v1 := range data {
			allData = append(allData, v1)
		}
		*this = append(*this, allData)
	}
	return nil
}

func (this *FloatSlice2) Decode(str string) error {
	*this = make(FloatSlice2, 0)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		ints, err := utils.Float64SliceFromString(v, COMMA)
		if err != nil {
			return err
		}
		*this = append(*this, ints)
	}
	return nil
}

func (this *IntFloatMap) Decode(str string) error {
	*this = make(IntFloatMap)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		strs := strings.Split(v, COMMA)
		(*this)[utils.StrToInt(strs[0])] = utils.StrToFloat(strs[1])
	}
	return nil
}

func (this *IntStrMap) Decode(str string) error {
	*this = make(IntStrMap)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		strs := strings.Split(v, COMMA)
		(*this)[utils.StrToInt(strs[0])] = strs[1]
	}
	return nil
}

func (this *StrIntMap) Decode(str string) error {
	*this = make(StrIntMap)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		strs := strings.Split(v, COMMA)
		(*this)[strs[0]] = utils.StrToInt(strs[1])
	}
	return nil
}

func (this *StrFloatMap) Decode(str string) error {
	*this = make(StrFloatMap)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		strs := strings.Split(v, COMMA)
		(*this)[strs[0]] = utils.StrToFloat(strs[1])
	}
	return nil
}

func (this *StrStrMap) Decode(str string) error {
	*this = make(StrStrMap)
	if len(str) == 0 {
		return nil
	}
	infoList := strings.Split(strings.Trim(strings.TrimSpace(str), PIPE), PIPE)
	if len(infoList) == 0 {
		return nil
	}

	for _, v := range infoList {
		strs := strings.Split(v, COMMA)
		(*this)[strs[0]] = strs[1]
	}
	return nil
}
