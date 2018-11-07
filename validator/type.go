package validator

import (
	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Type string

func (c Type) Validate(value interface{}, _ types.Flag) ErrorInterface {
	if helpers.GetTypeName(value) != string(c) {
		return NewConstraintError(c, `this value should be of type "%s"`, c)
	}
	return nil
}

func (_ Type) Strict() bool {
	return true
}

func (_ Type) Name() string {
	return "Type"
}
