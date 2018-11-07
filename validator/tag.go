package validator

import (
	"fmt"
	"regexp"
	"strconv"
)

var reTag = regexp.MustCompile(`([^,=]+)(?:=(\[.*?]|".*?"|[^,]+))?`)

func ParseTag(tag string) (constraints []ConstraintInterface, err error) {
	for _, match := range reTag.FindAllStringSubmatch(tag, -1) {
		ci, ok := TagOptions[match[1]]
		if !ok {
			return nil, fmt.Errorf(`gowl/validator: unknown tag option "%s"`, match[1])
		}
		constraint, err := ci(TagOption{match[1], match[2]})
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, constraint)
	}
	return
}

type TagOption struct {
	Name  string
	Value string
}

func (o TagOption) NoValue(lambda func() ConstraintInterface) (ConstraintInterface, error) {
	if o.Value != "" {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should not have value`, o.Name)
	}
	return lambda(), nil
}

func (o TagOption) WithBool(lambda func(v bool) ConstraintInterface, defaultValue bool) (ConstraintInterface, error) {
	if o.Value == "" {
		return lambda(defaultValue), nil
	}
	v, err := strconv.ParseBool(o.Value)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid bool value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithInt(lambda func(v int) ConstraintInterface) (ConstraintInterface, error) {
	v, err := strconv.Atoi(o.Value)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid int value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithUint(lambda func(v uint) ConstraintInterface) (ConstraintInterface, error) {
	v, err := strconv.ParseUint(o.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid uint value: %s`, o.Name, err)
	}
	return lambda(uint(v)), nil
}

func (o TagOption) WithString(lambda func(v string) ConstraintInterface) (ConstraintInterface, error) {
	n := len(o.Value)
	if n < 2 || o.Value[0] != '"' || o.Value[n-1] != '"' {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should be in format "..."`, o.Name)
	}
	return lambda(o.Value[1 : n-1]), nil
}

func (o TagOption) WithTag(lambda func(v []ConstraintInterface) ConstraintInterface) (ConstraintInterface, error) {
	n := len(o.Value)
	if n < 2 || o.Value[0] != '[' || o.Value[n-1] != ']' {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should be in format [x,y=z,...]`, o.Name)
	}
	v, err := ParseTag(o.Value[1 : n-1])
	if err != nil {
		return nil, err
	}
	return lambda(v), nil
}

var TagOptions = map[string]func(o TagOption) (constraint ConstraintInterface, err error){
	"required": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithBool(func(v bool) ConstraintInterface { return Required(v) }, true)
	},
	"valid": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithBool(func(v bool) ConstraintInterface { return Valid(v) }, true)
	},
	"len": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint) ConstraintInterface { return Length(v) })
	},
	"minlen": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint) ConstraintInterface { return MinLength(v) })
	},
	"maxlen": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint) ConstraintInterface { return MaxLength(v) })
	},
}

func init() {
	// avoid initialization loop
	TagOptions["each"] = func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithTag(func(v []ConstraintInterface) ConstraintInterface { return Each(v) })
	}
}
