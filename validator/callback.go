package validator

import (
	"github.com/lokhman/gowl/types"
)

type callback struct {
	validate func(value interface{}, flags types.Flag) ErrorInterface
	strict   bool
}

func (c callback) Validate(value interface{}, flags types.Flag) ErrorInterface {
	err := c.validate(value, flags)
	if err != nil {
		for _, err := range err.Errors() {
			err.Constraint = c
		}
	}
	return err
}

func (c callback) Strict() bool {
	return c.strict
}

func (_ callback) Name() string {
	return "Callback"
}

func Callback(validate func(value interface{}, flags types.Flag) ErrorInterface, strict bool) ConstraintInterface {
	return callback{validate, strict}
}
