package validator

import (
	"reflect"
	"strings"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type lookup struct {
	name     string
	value    interface{}
	inverted bool
}

func (c lookup) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		if c.inverted {
			for i := 0; i < v.Len(); i++ {
				if isEqual(v.Index(i).Interface(), c.value) {
					return NewConstraintError(c, "this value contains %#v", c.value)
				}
			}
		} else {
			for i := 0; i < v.Len(); i++ {
				if isEqual(v.Index(i).Interface(), c.value) {
					return nil
				}
			}
			return NewConstraintError(c, "this value does not contain %#v", c.value)
		}
	case reflect.Map:
		if c.inverted {
			for _, key := range v.MapKeys() {
				if isEqual(v.MapIndex(key).Interface(), c.value) {
					return NewConstraintError(c, "this value contains %#v", c.value)
				}
			}
		} else {
			for _, key := range v.MapKeys() {
				if isEqual(v.MapIndex(key).Interface(), c.value) {
					return nil
				}
			}
			return NewConstraintError(c, "this value does not contain %#v", c.value)
		}
	case reflect.String:
		cValue := reflect.ValueOf(c.value)
		if cValue.Kind() != reflect.String {
			return UnexpectedTypeError(c, value)
		}
		if contains := strings.Contains(v.String(), cValue.String()); c.inverted && contains {
			return NewConstraintError(c, "this value contains %#v", c.value)
		} else if !c.inverted && !contains {
			return NewConstraintError(c, "this value does not contain %#v", c.value)
		}
	default:
		return UnexpectedTypeError(c, value)
	}
	return nil
}

func (_ lookup) Strict() bool {
	return false
}

func (c lookup) Name() string {
	return c.name
}

func Contains(value interface{}) ConstraintInterface {
	return lookup{"Contains", value, false}
}

func Excludes(value interface{}) ConstraintInterface {
	return lookup{"Excludes", value, true}
}
