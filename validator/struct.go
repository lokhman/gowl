package validator

import (
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Struct map[string][]ConstraintInterface

func (c Struct) Validate(value interface{}, flags types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if v.Kind() != reflect.Struct {
		return UnexpectedTypeError(c, value)
	}
	t := v.Type()

	var ve ValidationError
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		constraints, ok := c[f.Name]
		if !ok {
			if tag := f.Tag.Get(TagName); tag != "" {
				var err error
				constraints, err = ParseTag(tag)
				if err != nil {
					panic(err)
				}
			} else if !flags.Has(AllowMissingFields) {
				return NewConstraintError(c, `field "%s" is missing in the constraint`, f.Name)
			}
		}
		if len(constraints) == 0 {
			continue
		}

		strict := checkStrict(constraints)
		value := v.Field(i).Interface()
		if !strict && flags.Has(OmitEmpty) && helpers.IsEmpty(value) {
			continue
		}

		for _, constraint := range constraints {
			if err := constraint.Validate(value, flags); err != nil {
				for _, e := range err.Errors() {
					if e.VarName != "" {
						e.VarName = f.Name + "." + e.VarName
					} else {
						e.VarName = f.Name
					}
				}
				ve = append(ve, err.Errors()...)
				if flags.Has(FastValidation) {
					goto fast
				}
				// strict constraints always go first
				if strict && flags.Has(StrictValidation) {
					break
				}
			}
		}
	}

fast:
	if !flags.Has(AllowExtraFields) {
		for name := range c {
			if _, ok := t.FieldByName(name); !ok {
				return NewConstraintError(c, `field "%s" is not expected in the constraint`, name)
			}
		}
	}
	return ve
}

func (_ Struct) Strict() bool {
	return false
}

func (_ Struct) Name() string {
	return "Struct"
}
