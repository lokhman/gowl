package validator

import (
	"reflect"
	"time"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

// Equal
type equal struct {
	value interface{}
}

func (c equal) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if !v.IsValid() || !isEqual(v.Interface(), c.value) {
		return NewConstraintError(c, "this value should be equal to `%v`", c.value)
	}
	return nil
}

func (_ equal) Strict() bool {
	return false
}

func (_ equal) Name() string {
	return "Equal"
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
	if v.IsValid() && isEqual(v.Interface(), c.value) {
		return NewConstraintError(c, "this value should not be equal to `%v`", c.value)
	}
	return nil
}

func (_ notEqual) Strict() bool {
	return false
}

func (_ notEqual) Name() string {
	return "NotEqual"
}

func NotEqual(value interface{}) ConstraintInterface {
	return notEqual{value}
}

// Identical
type identical struct {
	value interface{}
}

func (c identical) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if !isEqual(value, c.value) {
		return NewConstraintError(c, "this value should be identical to `%v`", c.value)
	}
	return nil
}

func (_ identical) Strict() bool {
	return false
}

func (_ identical) Name() string {
	return "Identical"
}

func Identical(value interface{}) ConstraintInterface {
	return identical{value}
}

// NotIdentical
type notIdentical struct {
	value interface{}
}

func (c notIdentical) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if isEqual(value, c.value) {
		return NewConstraintError(c, "this value should not be identical to `%v`", c.value)
	}
	return nil
}

func (_ notIdentical) Strict() bool {
	return false
}

func (_ notIdentical) Name() string {
	return "NotIdentical"
}

func NotIdentical(value interface{}) ConstraintInterface {
	return notIdentical{value}
}

// ...
func isEqual(x, y interface{}) bool {
	if xt, ok := x.(time.Time); ok {
		if yt, ok := y.(time.Time); ok {
			return xt.Equal(yt)
		}
	}
	return reflect.DeepEqual(x, y)
}
