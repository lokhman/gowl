package helpers

import (
	"reflect"
	"runtime"
)

func GetTypeName(i interface{}) string {
	if t := reflect.TypeOf(i); t != nil {
		return t.String()
	}
	return "nil"
}

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return Indirect(v.Elem())
	}
	return v
}

func IsEmpty(i interface{}) bool {
	v := Indirect(reflect.ValueOf(i))
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface:
		return v.IsNil()
	case reflect.Struct:
		// using DeepEqual to fix incomparable interface comparison
		return reflect.DeepEqual(reflect.Zero(v.Type()).Interface(), v.Interface())
	}
	return !v.IsValid()
}

func IsNil(i interface{}) bool {
	v := Indirect(reflect.ValueOf(i))
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return !v.IsValid()
}
