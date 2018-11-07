package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type MaxLength uint

func (c MaxLength) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.String:
		if v.Len() > int(c) {
			return NewConstraintError(c, "this value should have %d character(s) or less", c)
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		if v.Len() > int(c) {
			return NewConstraintError(c, "this value should contain %d element(s) or less", c)
		}
	default:
		return UnexpectedTypeError(c, value)
	}
	return nil
}

func (_ MaxLength) Strict() bool {
	return false
}

func (_ MaxLength) Name() string {
	return "maxlength"
}
