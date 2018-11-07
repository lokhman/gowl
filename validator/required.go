package validator

import (
	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Required bool

func (c Required) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if empty := helpers.IsEmpty(value); bool(c) && empty {
		return NewConstraintError(c, "this value should not be empty")
	} else if !bool(c) && !empty {
		return NewConstraintError(c, "this value should be empty")
	}
	return nil
}

func (_ Required) Strict() bool {
	return true
}

func (_ Required) Name() string {
	return "Required"
}

const (
	Empty    = Required(false)
	NotEmpty = Required(true)
)
