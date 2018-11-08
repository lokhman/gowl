package validator

import (
	"reflect"
	"time"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

// Equal
type equal struct {
	name     string
	value    interface{}
	inverted bool
}

func (c equal) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if equal := v.IsValid() && isEqual(v.Interface(), c.value); equal && c.inverted {
		return NewConstraintError(c, "this value should not be equal to `%v`", c.value)
	} else if !equal {
		return NewConstraintError(c, "this value should be equal to `%v`", c.value)
	}
	return nil
}

func (_ equal) Strict() bool {
	return false
}

func (c equal) Name() string {
	return c.name
}

func Equal(value interface{}) ConstraintInterface {
	return equal{"Equal", value, false}
}

func NotEqual(value interface{}) ConstraintInterface {
	return equal{"NotEqual", value, true}
}

// Identical
type identical struct {
	name     string
	value    interface{}
	inverted bool
}

func (c identical) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if equal := isEqual(value, c.value); equal && c.inverted {
		return NewConstraintError(c, "this value should not be identical to `%v`", c.value)
	} else if !equal {
		return NewConstraintError(c, "this value should be identical to `%v`", c.value)
	}
	return nil
}

func (_ identical) Strict() bool {
	return false
}

func (c identical) Name() string {
	return c.name
}

func Identical(value interface{}) ConstraintInterface {
	return identical{"Identical", value, false}
}

func NotIdentical(value interface{}) ConstraintInterface {
	return identical{"NotIdentical", value, true}
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

var (
	True  = Equal(true)
	False = Equal(false)
)
