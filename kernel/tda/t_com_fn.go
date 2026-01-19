package tda

import (
	"encoding/json"
	"fmt"
	"leaf/log"
	"reflect"
	"strings"
	"time"

	"github.com/v587-zyf/gc/utils"
)

// FlattenStructToMap 将结构体转换为 map[string]any，并将嵌套结构体字段“平铺”到顶层 map 中
// 同时使用 isEmptyValue 判断字段是否为空，空字段将被过滤
func FlattenStructToMap(obj interface{}) (map[string]any, error) {
	result := make(map[string]any)

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// 跳过未导出字段
		if field.PkgPath != "" {
			continue
		}

		// 获取 json tag，只处理顶层字段的 tag
		jsonTag := field.Tag.Get("json")
		name := field.Name
		if jsonTag != "" && jsonTag != "-" {
			name = strings.Split(jsonTag, ",")[0]
		}

		// 判断字段是否为空
		if isEmptyValue(value) && field.Tag.Get("canEmpty") == "" {
			//fmt.Printf("字段 %s 被跳过（判断为空）\n", name)
			continue
		}

		// 处理字段值
		switch value.Kind() {
		case reflect.Struct:
			if value.Type().String() == "time.Time" {
				//fmt.Printf("嵌套结构体为 time.Time 类型，将转换为字符串\n")
				result[name] = value.Interface().(time.Time).Format("2006-01-02 15:04:05")
				//result[name] = value.Interface().(time.Time).String()
			} else {
				//fmt.Printf("嵌套结构体为普通结构体，将转换为 map\n")
				nestedMap, err := FlattenStructToMap(value.Interface())
				if err != nil {
					fmt.Printf("嵌套结构体转换出错: %s %v\n", name, err)
					return nil, err
				}
				for k, v := range nestedMap {
					result[k] = v
				}
			}
		case reflect.Ptr:
			if value.IsNil() {
				continue
			}
			elem := value.Elem()
			if elem.Kind() == reflect.Struct {
				nestedMap, err := FlattenStructToMap(elem.Interface())
				if err != nil {
					fmt.Printf("嵌套结构体转换出错: %s %v\n", name, err)
					return nil, err
				}
				// 合并嵌套 map 到顶层
				for k, v := range nestedMap {
					result[k] = v
				}
			} else {
				// 判断指针指向的基本类型是否为空
				if isEmptyValue(elem) {
					//fmt.Printf("字段 %s 被跳过（判断为空）\n", name)
					continue
				}
				result[name] = elem.Interface()
			}
		default:
			result[name] = value.Interface()
		}
	}

	return result, nil
}

// isEmptyValue 判断字段是否为空（零值）
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
		return v.IsNil()
	case reflect.Struct:
		// 特别处理 time.Time 类型
		if v.Type().String() == "time.Time" {
			t := v.Interface().(time.Time)
			//fmt.Printf("检测到 time.Time: %v, 是否为零值: %v\n", t, t.Equal(time.Time{}))
			return t.Equal(time.Time{})
		}

		// 普通结构体字段递归判断是否为空
		for i := 0; i < v.NumField(); i++ {
			if !isEmptyValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

func UploadLog(channelId uint32, accountId, distinctId, eventName string, properties map[string]interface{}) {

	UploadKlLogServer(channelId, accountId, distinctId, eventName, properties)

	return
}

type KlLogReq struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type KlLogAck struct {
}

func UploadKlLogServer(channelId uint32, accountId, distinctId, eventName string, properties map[string]interface{}) {

	properties["channel_id"] = channelId
	properties["account_id"] = accountId
	properties["distinct_id"] = distinctId

	jsonData, err := json.Marshal(properties)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return
	}

	req := &KlLogReq{
		Cmd:  eventName,
		Data: string(jsonData),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Error(fmt.Sprintf("marshal err:%v ", err))
		return
	}

	url := fmt.Sprintf("https://log.jumpoyo.com/log/uploadStrategyLog")
	resp, err := utils.PostJson(url, reqBytes)
	if err != nil {
		log.Error(fmt.Sprintf("http post err:%v", err))
		return
	}

	ack := new(KlLogAck)
	if err = json.Unmarshal(resp, &ack); err != nil {
		log.Error(fmt.Sprintf("loginAck unmarshal err:%v", err))
		return
	}

	// if loginAck.Data != nil {
	// 	//r.SetToken(loginAck.Data.Token)
	// 	r.SetUserID(uint64(loginAck.Data.AccountId))
	// 	r.SetGateAddr(loginAck.Data.ServerInfo.ServerAddr)
	// }
}
