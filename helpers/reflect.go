package helpers

import (
	"reflect"
	"runtime"
)

func GetTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
