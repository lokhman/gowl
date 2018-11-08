package validator

import (
	"reflect"
	"time"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

const (
	rangeTooLowError      = "this value should be greater than %#v"
	rangeTooHighError     = "this value should be less than %#v"
	rangeTooLowInclError  = "this value should be %#v or more"
	rangeTooHighInclError = "this value should be %#v or less"
)

type range_ struct {
	name     string
	type_    reflect.Type
	min      *reflect.Value
	max      *reflect.Value
	included bool
}

func (c range_) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if !v.IsValid() || v.Type() != c.type_ {
		return UnexpectedTypeError(c, value)
	}

	switch c.type_.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := v.Int()
		if c.min != nil {
			if min := c.min.Int(); c.included && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.included && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Int(); c.included && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.included && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value := v.Uint()
		if c.min != nil {
			if min := c.min.Uint(); c.included && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.included && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Uint(); c.included && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.included && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Float32, reflect.Float64:
		value := v.Float()
		if c.min != nil {
			if min := c.min.Float(); c.included && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.included && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.Float(); c.included && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.included && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.String:
		value := v.String()
		if c.min != nil {
			if min := c.min.String(); c.included && value < min {
				return NewConstraintError(c, rangeTooLowInclError, min)
			} else if !c.included && value <= min {
				return NewConstraintError(c, rangeTooLowError, min)
			}
		}
		if c.max != nil {
			if max := c.max.String(); c.included && value > max {
				return NewConstraintError(c, rangeTooHighInclError, max)
			} else if !c.included && value >= max {
				return NewConstraintError(c, rangeTooHighError, max)
			}
		}
	case reflect.Struct:
		switch value := v.Interface().(type) {
		case time.Time:
			if c.min != nil {
				if min := c.min.Interface().(time.Time); c.included && value.Before(min) {
					return NewConstraintError(c, rangeTooLowInclError, min)
				} else if !c.included && (value.Equal(min) || value.Before(min)) {
					return NewConstraintError(c, rangeTooLowError, min)
				}
			}
			if c.min != nil {
				if min := c.min.Interface().(time.Time); c.included && value.After(min) {
					return NewConstraintError(c, rangeTooHighInclError, min)
				} else if !c.included && (value.Equal(min) || value.After(min)) {
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

func (c range_) Name() string {
	return c.name
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
	return range_{"LessThan", t, nil, &v, false}
}

func LessThanOrEqual(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{"LessThanOrEqual", t, nil, &v, true}
}

func GreaterThan(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{"GreaterThan", t, &v, nil, false}
}

func GreaterThanOrEqual(value interface{}) ConstraintInterface {
	v := reflect.ValueOf(value)
	t := checkRangeArgument(v)
	return range_{"GreaterThanOrEqual", t, &v, nil, true}
}

func Between(min, max interface{}) ConstraintInterface {
	minV, maxV := reflect.ValueOf(min), reflect.ValueOf(max)
	minT := checkRangeArgument(minV)
	if minT != checkRangeArgument(maxV) {
		panic("gowl/validator: min and max arguments should have similar types")
	}
	return range_{"Between", minT, &minV, &maxV, true}
}
