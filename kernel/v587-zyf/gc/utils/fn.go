package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFunctionName(i interface{}) string {

	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetFunNameByCaller(l int) (string, int) {
	pc, _, line, _ := runtime.Caller(l)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	return split[len(split)-1], line
}
