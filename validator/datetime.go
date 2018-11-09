package validator

import (
	"reflect"
	"time"

	"github.com/lokhman/gowl/helpers"
	"github.com/lokhman/gowl/types"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
)

type Timestamp string

func (c Timestamp) Validate(value interface{}, _ types.Flag) ErrorInterface {
	v := helpers.Indirect(reflect.ValueOf(value))
	if v.Kind() != reflect.String {
		return UnexpectedTypeError(c, value)
	}
	if _, err := time.Parse(string(c), v.String()); err != nil {
		return NewConstraintError(c, "this value is not a valid timestamp")
	}
	return nil
}

func (_ Timestamp) Strict() bool {
	return false
}

func (_ Timestamp) Name() string {
	return "Timestamp"
}

const (
	Date     = Timestamp(DateFormat)
	Time     = Timestamp(TimeFormat)
	DateTime = Timestamp(DateTimeFormat)
)
