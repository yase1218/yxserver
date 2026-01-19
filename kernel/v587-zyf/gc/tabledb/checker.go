package tabledb

import (
	"errors"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"reflect"
	"strings"
)

type CheckIter interface {
	check() string
}
type tagChecker struct {
	fn   func(a ...any) bool
	desc string
}

var tagCheckerMap = make(map[string]tagChecker)

var itemTypeChecker = make(map[int]func(a ...any) bool)

func init() {
	tagCheckerMap["condition"] = newTagChecker(conditionExist, "条件中没有 %d (如果是32,可能是任务id不存在)")
}

func newTagChecker(fn func(a ...any) bool, desc string) tagChecker {
	return tagChecker{fn: fn, desc: desc}
}

func (t *TableDb) checkFunc(value any) (bool, []string) {
	allErrorMsg := make([]string, 0)
	hasCheck := false
	//如果结构体实现了Checkable
	if c, ok := value.(CheckIter); ok {
		errMsg := c.check()
		if len(errMsg) > 0 {
			allErrorMsg = append(allErrorMsg, errMsg)
		}
		hasCheck = true
	}
	//如果，去检查Tag 中配置的检查
	typeValue := reflect.ValueOf(value)
	if t.hasTagCheck(typeValue) {
		err := t.checkOneDataMember(typeValue)
		if err != nil {
			allErrorMsg = append(allErrorMsg, "Struct:"+reflect.TypeOf(value).String()+" "+err.Error())
		}
		hasCheck = true
	}
	return hasCheck, allErrorMsg
}

func (t *TableDb) checkTableBase() error {
	allErrorMsg := make([]string, 0)
	gameDbBaseValueOf := reflect.ValueOf(reflect.ValueOf(iTdb).Elem().FieldByName("TableBase").Interface()).Elem()

	for j := 0; j < gameDbBaseValueOf.NumField(); j++ {
		f := gameDbBaseValueOf.Field(j)
		fKind := f.Kind()
		if fKind == reflect.Map {
			for _, v := range f.MapKeys() {
				value := f.MapIndex(v).Interface()
				hasCheck, err := t.checkFunc(value)
				if !hasCheck {
					break
				}
				if len(err) > 0 {
					allErrorMsg = append(allErrorMsg, err...)
				}
			}
		}
	}

	if len(allErrorMsg) < 1 {
		return nil
	}

	return errors.New(strings.Join(allErrorMsg, "\n"))
}

func (t *TableDb) checkAllViaObj(objBeChecker any) error {
	values := reflect.ValueOf(objBeChecker).Elem()

	allErrorMsg := make([]string, 0)

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i)
		fKind := f.Kind()
		switch fKind {
		case reflect.Slice:
			for i := 0; i < f.Len(); i++ {
				value := f.Index(i).Interface()
				hasCheck, err := t.checkFunc(value)
				if !hasCheck {
					break
				}
				if len(err) > 0 {
					allErrorMsg = append(allErrorMsg, err...)
				}
			}
		case reflect.Map:
			for _, v := range f.MapKeys() {
				value := f.MapIndex(v).Interface()
				hasCheck, err := t.checkFunc(value)
				if !hasCheck {
					break
				}
				if len(err) > 0 {
					allErrorMsg = append(allErrorMsg, err...)
				}
			}
		case reflect.Ptr:
			hasCheck, err := t.checkFunc(f.Interface())
			if !hasCheck {
				break
			}
			if len(err) > 0 {
				allErrorMsg = append(allErrorMsg, err...)
			}

		default:
			//fmt.Println("unhandle type in gamedb:", fKind)
		}
	}
	if len(allErrorMsg) < 1 {
		return nil
	}

	return errors.New(strings.Join(allErrorMsg, "\n"))
}
func (t *TableDb) hasTagCheck(v reflect.Value) bool {
	if (v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr) && v.Elem().Kind() == reflect.Struct {
		for j := 0; j < v.Elem().NumField(); j++ {
			if len(v.Elem().Type().Field(j).Tag.Get("checker")) > 0 {
				return true
			}
		}
	}
	return false
}
func (t *TableDb) Check() error {
	checkers := []func(db iface.ITableDb) error{
		//(*GameDb).checkConf,
	}
	for _, checker := range checkers {
		err := checker(iTdb)
		if err != nil {
			return err
		}
	}
	err := t.checkAll()
	if err != nil {
		return err
	}

	return nil
}
func (t *TableDb) CheckConf(c any) (err error) {
	return t.checkAllViaObj(c)
}
func (t *TableDb) checkAll() error {
	err := t.checkTableBase()
	if err != nil {
		return err
	}
	return nil
}
func (t *TableDb) checkOneDataMember(v reflect.Value) error {
	errMsgsSlice := make([]string, 0)
	errMsgs := &errMsgsSlice

	if (v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr) && v.Elem().Kind() == reflect.Struct {
		for j := 0; j < v.Elem().NumField(); j++ {
			checkerName := v.Elem().Type().Field(j).Tag.Get("checker")
			fieldType := v.Elem().Type().Field(j).Type.Name()
			fieldName := v.Elem().Type().Field(j).Name
			checkerName = strings.TrimSpace(checkerName)
			if len(checkerName) < 1 {
				continue
			}
			checkerNameArr := strings.Split(checkerName, "_")
			checkerFnName := checkerNameArr[0]
			checkerNameExtra := checkerNameArr[1:]
			checkerFnName = strings.ToLower(checkerFnName)
			var tagc tagChecker
			if _, ok := tagCheckerMap[checkerFnName]; ok {
				tagc = tagCheckerMap[checkerFnName]
			} else {
				fmt.Printf("No checker for %s\n", checkerFnName)
				continue
			}

			switch fieldType {
			case "int":
				id := int(v.Elem().Field(j).Int())
				if id == 0 {
					continue
				}
				checkOneId(errMsgs, tagc, fieldName, id)
			case "IntSlice":
				data := v.Elem().Field(j).Interface().(IntSlice)
				for m := 0; m < len(data); m++ {
					checkOneId(errMsgs, tagc, fieldName, data[m])
				}
			case "IntMap":
				data := v.Elem().Field(j).Interface().(IntMap)
				for k, v := range data {
					checkOneId(errMsgs, tagc, fieldName, k, v)
				}
			case "string":
				str := v.Elem().Field(j).String()
				checkOneId(errMsgs, tagc, fieldName, str, checkerNameExtra)
			case "StringSlice":
				data := v.Elem().Field(j).Interface().(StringSlice)
				for m := 0; m < len(data); m++ {
					checkOneId(errMsgs, tagc, fieldName, data[m], checkerNameExtra)
				}
			case "MedalCheckId":
				id := int(v.Elem().Field(j).Int())
				checkOneId(errMsgs, tagc, fieldName, id)
			default:
				//fmt.Printf("fieldType = %+v no handler ,for fieldName %s \n", fieldType, fieldName)
			}
		}
	}
	if len(*errMsgs) < 1 {
		return nil
	}
	return errors.New(strings.Join(*errMsgs, "\n"))
}

func conditionExist(args ...any) bool {
	//id := args[0].(int)
	//if gameDb.GetCondition(id) == nil {
	//	return false
	//}
	//if len(args) < 2 {
	//	fmt.Printf("conditionExist:why args is %v", args)
	//	return false
	//}
	//value := args[1].(int)
	//
	//if len(args) > 2 {
	//	subs := args[2].(map[int]int)
	//	for subType, v := range subs {
	//		if gameDb.ConditionSubTypes[id*100+subType] == nil || v < 1 {
	//			return false
	//		}
	//	}
	//}
	//
	//switch id {
	//case pb.CONDITIONTYPE_FINISH_MAIN_LINE_TASK:
	//	return taskExist(value)
	//default:
	//	return true
	//}
	return true
}
func checkOneId(errMsgs *[]string, tc tagChecker, fieldName string, args ...any) {
	if !tc.fn(args...) {
		*errMsgs = append(*errMsgs, fmt.Sprintf("字段 %s 在%v", fieldName, fmt.Sprintf(tc.desc, args)))
	}
}
