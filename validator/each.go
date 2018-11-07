package validator

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Each []ConstraintInterface

func (c Each) Validate(value interface{}, flags types.Flag) ErrorInterface {
	if len(c) == 0 {
		return nil
	}
	strict := checkStrict(c)

	v := helpers.Indirect(reflect.ValueOf(value))
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		var ve ValidationError
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i).Interface()
			if !strict && flags.Has(OmitEmpty) && helpers.IsEmpty(value) {
				continue
			}

			for _, constraint := range c {
				if err := constraint.Validate(value, flags); err != nil {
					for _, e := range err.Errors() {
						if key := strconv.Itoa(i); e.VarName != "" {
							e.VarName = key + "." + e.VarName
						} else {
							e.VarName = key
						}
					}
					ve = append(ve, err.Errors()...)
					if flags.Has(FastValidation) {
						return ve
					}
					// strict constraints always go first
					if strict && flags.Has(StrictValidation) {
						break
					}
				}
			}
		}
		return ve
	case reflect.Map:
		var ve ValidationError
		for _, kv := range v.MapKeys() {
			value := v.MapIndex(kv).Interface()
			if !strict && flags.Has(OmitEmpty) && helpers.IsEmpty(value) {
				continue
			}

			for _, constraint := range c {
				if err := constraint.Validate(value, flags); err != nil {
					for _, e := range err.Errors() {
						key := fmt.Sprintf("%v", kv.Interface())
						if e.VarName != "" {
							e.VarName = key + "." + e.VarName
						} else {
							e.VarName = key
						}
					}
					ve = append(ve, err.Errors()...)
					if flags.Has(FastValidation) {
						return ve
					}
					// strict constraints always go first
					if strict && flags.Has(StrictValidation) {
						break
					}
				}
			}
		}
		return ve
	}
	return UnexpectedTypeError(c, value)
}

func (_ Each) Strict() bool {
	return false
}

func (_ Each) Name() string {
	return "each"
}
