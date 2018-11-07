package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

// Equal
type equal struct {
	value interface{}
}

func (c equal) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if !v.IsValid() || !reflect.DeepEqual(v.Interface(), c.value) {
		return NewConstraintError(c, `this value should be equal to "%v"`, c.value)
	}
	return nil
}

func (_ equal) Strict() bool {
	return false
}

func (_ equal) Name() string {
	return "equal"
}

func Equal(value interface{}) ConstraintInterface {
	return equal{value}
}

// NotEqual
type notEqual struct {
	value interface{}
}

func (c notEqual) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if v.IsValid() && reflect.DeepEqual(v.Interface(), c.value) {
		return NewConstraintError(c, `this value should not be equal to "%v"`, c.value)
	}
	return nil
}

func (_ notEqual) Strict() bool {
	return false
}

func (_ notEqual) Name() string {
	return "not_equal"
}

func NotEqual(value interface{}) ConstraintInterface {
	return notEqual{value}
}
