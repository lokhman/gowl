package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Length uint

func (c Length) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.String:
		if v.Len() != int(c) {
			return NewConstraintError(c, "this value should have exactly %d character(s)", c)
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		if v.Len() != int(c) {
			return NewConstraintError(c, "this value should contain exactly %d element(s)", c)
		}
	default:
		return UnexpectedTypeError(c, value)
	}
	return nil
}

func (_ Length) Strict() bool {
	return false
}

func (_ Length) Name() string {
	return "length"
}
