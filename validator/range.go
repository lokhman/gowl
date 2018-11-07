package validator

import (
	"reflect"
	"time"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

const (
	rangeTooLowError      = "this value should be greater than `%v`"
	rangeTooHighError     = "this value should be less than `%v`"
	rangeTooLowInclError  = "this value should be `%v` or more"
	rangeTooHighInclError = "this value should be `%v` or less"
)

type range_ struct {
	typ  reflect.Type
	min  *reflect.Value
	max  *reflect.Value
	incl bool
}

func (c range_) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if !v.IsValid() || v.Type() != c.typ {
		return UnexpectedTypeError(c, value)
	}

	switch c.typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := v.Int()
		if c.min != nil {
			if min := c.min.Int(); c.incl && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.incl && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Int(); c.incl && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.incl && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value := v.Uint()
		if c.min != nil {
			if min := c.min.Uint(); c.incl && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.incl && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Uint(); c.incl && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.incl && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Float32, reflect.Float64:
		value := v.Float()
		if c.min != nil {
			if min := c.min.Float(); c.incl && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.incl && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Float(); c.incl && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.incl && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.String:
		value := v.String()
		if c.min != nil {
			if min := c.min.String(); c.incl && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.incl && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.String(); c.incl && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.incl && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Struct:
		switch value := v.Interface().(type) {
		case time.Time:
			if c.min != nil {
				if min := c.min.Interface().(time.Time); c.incl && value.Before(min) {
					return NewConstraintError(c, rangeTooLowInclError, min)
				} else if !c.incl && (value.Equal(min) || value.Before(min)) {
					return NewConstraintError(c, rangeTooLowError, min)
				}
			}
			if c.min != nil {
				if min := c.min.Interface().(time.Time); c.incl && value.After(min) {
					return NewConstraintError(c, rangeTooHighInclError, min)
				} else if !c.incl && (value.Equal(min) || value.After(min)) {
					return NewConstraintError(c, rangeTooHighError, min)
				}
			}
		}
	}
	return nil
}

func (_ range_) Strict() bool {
	return false
}

func (_ range_) Name() string {
	return "range"
}

func checkRangeArgument(v reflect.Value) reflect.Type {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return v.Type()
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			return v.Type()
		}
	}
	panic("gowl/validator: incompatible argument type")
}

func LessThan(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{t, nil, &v, false}
}

func LessThanOrEqual(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{t, nil, &v, true}
}

func GreaterThan(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{t, &v, nil, false}
}

func GreaterThanOrEqual(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{t, &v, nil, true}
}

func Between(min, max interface{}) ConstraintInterface {
	minV, maxV := reflect.ValueOf(min), reflect.ValueOf(max)
	minT := checkRangeArgument(minV)
	if minT != checkRangeArgument(maxV) {
		panic("gowl/validator: min and max arguments should have similar types")
	}
	return range_{minT, &minV, &maxV, true}
}
