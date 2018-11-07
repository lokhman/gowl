package validator

import (
	"fmt"
	"strings"

	"github.com/lokhman/gowl/helpers"
)

type ErrorInterface interface {
	Error() string
	Errors() ValidationError
}

type ConstraintError struct {
	Constraint ConstraintInterface
	VarName    string
	Message    string
	Args       []interface{}
}

func (e *ConstraintError) Error() string {
	if len(e.Args) > 0 {
		return fmt.Sprintf(e.Message, e.Args...)
	}
	return e.Message
}

func (e *ConstraintError) Errors() ValidationError {
	return ValidationError{e}
}

func NewConstraintError(constraint ConstraintInterface, message string, args ...interface{}) *ConstraintError {
	return &ConstraintError{constraint, "", message, args}
}

type ValidationError []*ConstraintError

func (e ValidationError) Error() string {
	buf := new(strings.Builder)
	for i, err := range e {
		if i > 0 {
			buf.WriteString("; ")
		}
		if err.VarName != "" {
			buf.WriteByte('[')
			buf.WriteString(err.VarName)
			buf.WriteString("] ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

func (e ValidationError) Errors() ValidationError {
	return e
}

func UnexpectedTypeError(constraint ConstraintInterface, value interface{}) *ConstraintError {
	return NewConstraintError(constraint, `unexpected value of type "%s"`, helpers.GetTypeName(value))
}
