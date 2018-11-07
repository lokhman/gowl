package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Length uint

// Length
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

// MinLength
type MinLength uint

func (c MinLength) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.String:
		if v.Len() < int(c) {
			return NewConstraintError(c, "this value should have %d character(s) or more", c)
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		if v.Len() < int(c) {
			return NewConstraintError(c, "this value should contain %d element(s) or more", c)
		}
	default:
		return UnexpectedTypeError(c, value)
	}
	return nil
}

func (_ MinLength) Strict() bool {
	return false
}

func (_ MinLength) Name() string {
	return "min_length"
}

// MaxLength
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
	return "max_length"
}
