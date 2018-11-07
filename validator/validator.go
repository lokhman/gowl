package validator

import (
	"fmt"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

const (
	_ types.Flag = 1 << iota

	StrictValidation
	FastValidation
	OmitEmpty

	AllowExtraFields
	AllowMissingFields
)

var Flags = StrictValidation | OmitEmpty | AllowMissingFields
var TagName = "validate"

type ConstraintInterface interface {
	Validate(value interface{}, flags types.Flag) ErrorInterface
	Strict() bool
	Name() string
}

func Validate(value interface{}, constraints ...ConstraintInterface) ValidationError {
	return ValidateSpecial(value, constraints, Flags)
}

func ValidateSpecial(value interface{}, constraints []ConstraintInterface, flags types.Flag) ValidationError {
	var strict bool
	if len(constraints) == 0 {
		constraints = []ConstraintInterface{Struct{}}
	} else if strict = checkStrict(constraints); !strict {
		if flags.Has(OmitEmpty) && helpers.IsEmpty(value) {
			return nil
		}
	}

	var ve ValidationError
	for _, constraint := range constraints {
		if err := constraint.Validate(value, flags); err != nil {
			ve = append(ve, err.Errors()...)
			if flags.Has(FastValidation) {
				break
			}
			// strict constraint always go first
			if strict && flags.Has(StrictValidation) {
				break
			}
		}
	}
	return ve
}

func checkStrict(constraints []ConstraintInterface) bool {
	for i := len(constraints) - 1; i >= 0; i-- {
		if constraint := constraints[i]; constraint.Strict() {
			if i > 0 {
				panic(fmt.Sprintf(`gowl/validator: constraint %s should go first`, constraint.Name()))
			}
			return true
		}
	}
	return false
}
