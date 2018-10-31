package gowl

import (
	"reflect"
	"runtime"
)

func getTypeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
