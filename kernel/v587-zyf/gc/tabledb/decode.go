package tabledb

import (
	"errors"
	"fmt"
	"github.com/v587-zyf/gc/utils"
	"reflect"
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

type FloatSlice []float64
type IntSlice []int
type IntSlice2 [][]int
type StringSlice []string
type StringSlice2 [][]string

type HmsTime struct {
	Hour   int `client:"hour"`
	Minute int `client:"minute"`
	Second int `client:"second"`
}
type HmsTimes []*HmsTime
type IntMap map[int]int

/***************************************************************/
/************************数据类型解析*****************************/
/***************************************************************/
func (this *IntSlice) Decode(str string) error {
	ints, err := utils.IntSliceFromString(str, PIPE)
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
	ints, err := utils.Float64SliceFromString(str, PIPE)
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
	*this = strings.Split(str, PIPE)
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
			return errors.New(v + "IntMap 属性信息格式错误")
		}

		k, err := strconv.Atoi(list[0])
		if err != nil {
			return err
		}
		if _, ok := (*this)[k]; ok {
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

func (this *HmsTime) String() string {
	return fmt.Sprintf("%d:%d:%d", this.Hour, this.Minute, this.Second)
}
func (this *HmsTimes) Decode(str string) error {
	*this = make(HmsTimes, 0)
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(strings.TrimSpace(str), PIPE)
	if len(infoList) < 1 {
		return errors.New(str + "属性信息格式错误")
	}
	for _, v := range infoList {
		one := &HmsTime{}
		err := one.Decode(v)
		if err != nil {
			return err
		}
		*this = append(*this, one)
	}
	return nil
}
func (this *HmsTime) Decode(str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return nil
	}
	infoList := strings.Split(strings.TrimSpace(str), COLON)
	if len(infoList) < 1 {
		return errors.New(str + "属性信息格式错误")
	}
	var hms HmsTime
	hms.Hour, _ = strconv.Atoi(infoList[0])
	if len(infoList) > 1 {
		hms.Minute, _ = strconv.Atoi(infoList[1])
	}
	if len(infoList) > 2 {
		hms.Second, _ = strconv.Atoi(infoList[2])
	}
	if hms.Hour < 0 || hms.Hour > 23 || hms.Minute < 0 || hms.Minute > 59 || hms.Second < 0 || hms.Second > 59 {
		return errors.New(str + "时分秒不对")
	}
	*this = hms
	return nil
}
func (this *HmsTime) GetSecondsFromZero() int { //从0点到该时刻的秒数
	return this.Hour*60*60 + this.Minute*60 + this.Second
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

func DecodeConfValues(obj interface{}, gameConfigs map[string]*GlobalBaseCfg) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("%v must be a struct pointer", obj)
	}
	var values = make(map[string]string, 0)
	for _, v := range gameConfigs {
		values[v.Name] = v.Value
	}

	objT = objT.Elem()
	objV = objV.Elem()
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)

		if fieldT.Type.Kind() == reflect.Ptr {
			//param := reflect.New(fieldT.Type.Elem())
			//groupConfigName := strings.TrimSpace(fieldT.Tag.Get("confgroup"))
			//if len(groupConfigName) > 0 {
			//	err := DecodeConfValues(groupConfigName, param.Interface(), gameConfigs)
			//	if err != nil {
			//		return err
			//	}
			//	fieldV.Set(param)
			//}
			fmt.Printf("未知的指针类型name:%v,type:%v", fieldT.Name, fieldT.Type.Kind())
			continue
		}

		configName := strings.TrimSpace(fieldT.Tag.Get("conf"))
		if len(configName) == 0 {
			continue
		}
		defaultDefine := strings.TrimSpace(fieldT.Tag.Get("default"))
		value := values[configName]
		if len(value) == 0 {
			value = defaultDefine
		}
		cellString := strings.TrimSpace(value)
		if decoder, ok := fieldV.Addr().Interface().(Decoder); ok {
			err := decoder.Decode(cellString)
			if err != nil {
				return err
			}
			continue
		}
		switch fieldT.Type.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				//fmt.Println("-------------------", configName)
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				//fmt.Println("-------------------", configName)
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			fieldV.SetString(value)
		}
	}
	return nil
}
