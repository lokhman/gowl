package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type length struct {
	name string
	min  int
	max  int
}

func (c length) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.String:
		if n := v.Len(); c.min == c.max && n != c.min {
			return NewConstraintError(c, "this value should have exactly %d character(s)", c)
		} else if c.min != -1 && n < c.min {
			return NewConstraintError(c, "this value should have %d character(s) or more", c)
		} else if c.max != -1 && n > c.max {
			return NewConstraintError(c, "this value should have %d character(s) or less", c)
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		if n := v.Len(); c.min == c.max && n != c.min {
			return NewConstraintError(c, "this value should contain exactly %d element(s)", c)
		} else if c.min != -1 && n < c.min {
			return NewConstraintError(c, "this value should contain %d element(s) or more", c)
		} else if c.max != -1 && n > c.max {
			return NewConstraintError(c, "this value should contain %d element(s) or less", c)
		}
	default:
		return UnexpectedTypeError(c, value)
	}
	return nil
}

func (_ length) Strict() bool {
	return false
}

func (c length) Name() string {
	return c.name
}

func Length(value uint) ConstraintInterface {
	return length{"Length", int(value), int(value)}
}

func MinLength(value uint) ConstraintInterface {
	return length{"MinLength", int(value), -1}
}

func MaxLength(value uint) ConstraintInterface {
	return length{"MaxLength", -1, int(value)}
}
