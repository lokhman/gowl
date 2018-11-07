package validator

import (
	"fmt"
	"reflect"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type Map map[interface{}][]ConstraintInterface

func (c Map) Validate(value interface{}, flags types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if v.Kind() != reflect.Map { // || v.Type().Key().Kind() != reflect.String {
		return UnexpectedTypeError(c, value)
	}

	var ve ValidationError
	for _, kv := range v.MapKeys() {
		key := kv.Interface()
		constraints, ok := c[key]
		if !ok && !flags.Has(AllowMissingFields) {
			return NewConstraintError(c, "key `%v` is missing in the constraint", key)
		} else if len(constraints) == 0 {
			continue
		}

		strict := checkStrict(constraints)
		value := v.MapIndex(kv).Interface()
		if !strict && flags.Has(OmitEmpty) && helpers.IsEmpty(value) {
			continue
		}

		for _, constraint := range constraints {
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
		kt := v.Type().Key()
		for key := range c {
			kv := reflect.ValueOf(key)
			if !kv.Type().AssignableTo(kt) || !v.MapIndex(kv).IsValid() {
				return NewConstraintError(c, "key `%v` is not expected in the constraint", key)
			}
		}
	}
	return ve
}

func (_ Map) Strict() bool {
	return false
}

func (_ Map) Name() string {
	return "map"
}
