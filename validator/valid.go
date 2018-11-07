package validator

import (
	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Valid bool

func (c Valid) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if nil := helpers.IsNil(value); bool(c) && nil {
		return NewConstraintError(c, "this value should not be nil")
	} else if !bool(c) && !nil {
		return NewConstraintError(c, "this value should be nil")
	}
	return nil
}

func (_ Valid) Strict() bool {
	return true
}

func (_ Valid) Name() string {
	return "valid"
}

const (
	Nil    = Valid(false)
	NotNil = Valid(true)
)
