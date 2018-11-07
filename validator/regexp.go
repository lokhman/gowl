package validator

import (
	"io"
	"reflect"
	"regexp"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

type regexp_ struct {
	*regexp.Regexp
}

func (c regexp_) Validate(value interface{}, _ types.Flag) ErrorInterface {
	var match bool
	v := helpers.Indirect(reflect.ValueOf(value))
	switch i := v.Interface().(type) {
	case []byte:
		match = c.Match(i)
	case io.RuneReader:
		match = c.MatchReader(i)
	default:
		if v.Kind() != reflect.String {
			return UnexpectedTypeError(c, value)
		}
		match = c.MatchString(v.String())
	}
	if !match {
		return NewConstraintError(c, "this value is not valid")
	}
	return nil
}

func (_ regexp_) Strict() bool {
	return false
}

func (_ regexp_) Name() string {
	return "regexp"
}

func Regexp(expr string) ConstraintInterface {
	re, err := regexp.Compile(expr)
	if err != nil {
		panic("gowl/validator: " + err.Error())
	}
	return regexp_{re}
}

func RegexpPOSIX(expr string) ConstraintInterface {
	re, err := regexp.CompilePOSIX(expr)
	if err != nil {
		panic("gowl/validator: " + err.Error())
	}
	return regexp_{re}
}

func RegexpCompiled(re *regexp.Regexp) ConstraintInterface {
	if re == nil {
		panic("gowl/validator: regexp is not compiled")
	}
	return regexp_{re}
}
