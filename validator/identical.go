package validator

import (
	"reflect"

	"github.com/lokhman/gowl/types"
)

// Identical
type identical struct {
	value interface{}
}

func (c identical) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if !reflect.DeepEqual(value, c.value) {
		return NewConstraintError(c, `this value should be identical to "%v"`, c.value)
	}
	return nil
}

func (_ identical) Strict() bool {
	return false
}

func (_ identical) Name() string {
	return "identical"
}

func Identical(value interface{}) ConstraintInterface {
	return identical{value}
}

// NotIdentical
type notIdentical struct {
	value interface{}
}

func (c notIdentical) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if reflect.DeepEqual(value, c.value) {
		return NewConstraintError(c, `this value should not be identical to "%v"`, c.value)
	}
	return nil
}

func (_ notIdentical) Strict() bool {
	return false
}

func (_ notIdentical) Name() string {
	return "not_identical"
}

func NotIdentical(value interface{}) ConstraintInterface {
	return notIdentical{value}
}
