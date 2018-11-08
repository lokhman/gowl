package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lokhman/gowl/helpers"
)

type TagOption struct {
	Name   string
	Value  string
	Prefix string
}

func (o TagOption) NoValue(lambda func() ConstraintInterface) (ConstraintInterface, error) {
	if o.Value != "" || o.Prefix != "" {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should not have value`, o.Name)
	}
	return lambda(), nil
}

func (o TagOption) WithBool(lambda func(v bool) ConstraintInterface) (ConstraintInterface, error) {
	if o.Prefix != "" {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid bool prefix "%s"`, o.Name, o.Prefix)
	}
	if o.Value == "" {
		return lambda(true), nil
	}
	v, err := strconv.ParseBool(o.Value)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid bool value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithInt(lambda func(v int64) ConstraintInterface) (ConstraintInterface, error) {
	var bitSize int
	switch o.Prefix {
	case "i8":
		bitSize = 8
	case "i16":
		bitSize = 16
	case "i32":
		bitSize = 32
	case "i64":
		bitSize = 64
	case "i", "": // 0
	default:
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid int prefix "%s"`, o.Name, o.Prefix)
	}
	v, err := strconv.ParseInt(o.Value, 10, bitSize)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid int value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithUint(lambda func(v uint64) ConstraintInterface) (ConstraintInterface, error) {
	var bitSize int
	switch o.Prefix {
	case "u8":
		bitSize = 8
	case "u16":
		bitSize = 16
	case "u32":
		bitSize = 32
	case "u64":
		bitSize = 64
	case "u", "": // 0
	default:
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid uint prefix "%s"`, o.Name, o.Prefix)
	}
	v, err := strconv.ParseUint(o.Value, 10, bitSize)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid uint value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithFloat(lambda func(v float64) ConstraintInterface) (ConstraintInterface, error) {
	var bitSize int
	switch o.Prefix {
	case "f32":
		bitSize = 32
	case "f64", "f", "":
		bitSize = 64
	default:
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid float prefix "%s"`, o.Name, o.Prefix)
	}
	v, err := strconv.ParseFloat(o.Value, bitSize)
	if err != nil {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid float value: %s`, o.Name, err)
	}
	return lambda(v), nil
}

func (o TagOption) WithString(lambda func(v string) ConstraintInterface) (ConstraintInterface, error) {
	if o.Prefix != "" {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid string prefix "%s"`, o.Name, o.Prefix)
	}
	n := len(o.Value)
	if n < 2 || o.Value[0] != '\'' || o.Value[n-1] != '\'' {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should be in format '...'`, o.Name)
	}
	return lambda(helpers.StripSlashes(o.Value[1:n-1], `'"`)), nil
}

func (o TagOption) WithValue(lambda func(v interface{}) ConstraintInterface) (ConstraintInterface, error) {
	if o.Value == "" {
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" should have value`, o.Name)
	}
	if o.Value[0] == '\'' {
		return o.WithString(func(v string) ConstraintInterface { return lambda(v) })
	}
	switch o.Prefix {
	case "i", "":
		return o.WithInt(func(v int64) ConstraintInterface { return lambda(int(v)) })
	case "i8":
		return o.WithInt(func(v int64) ConstraintInterface { return lambda(int8(v)) })
	case "i16":
		return o.WithInt(func(v int64) ConstraintInterface { return lambda(int16(v)) })
	case "i32":
		return o.WithInt(func(v int64) ConstraintInterface { return lambda(int32(v)) })
	case "i64":
		return o.WithInt(func(v int64) ConstraintInterface { return lambda(int64(v)) })
	case "u":
		return o.WithUint(func(v uint64) ConstraintInterface { return lambda(uint(v)) })
	case "u8":
		return o.WithUint(func(v uint64) ConstraintInterface { return lambda(uint8(v)) })
	case "u16":
		return o.WithUint(func(v uint64) ConstraintInterface { return lambda(uint16(v)) })
	case "u32":
		return o.WithUint(func(v uint64) ConstraintInterface { return lambda(uint32(v)) })
	case "u64":
		return o.WithUint(func(v uint64) ConstraintInterface { return lambda(uint(v)) })
	case "f32":
		return o.WithFloat(func(v float64) ConstraintInterface { return lambda(float32(v)) })
	case "f64", "f":
		return o.WithFloat(func(v float64) ConstraintInterface { return lambda(float64(v)) })
	default:
		return nil, fmt.Errorf(`gowl/validator: tag option "%s" has invalid prefix "%s"`, o.Name, o.Prefix)
	}
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

func ParseTag(tag string) (constraints []ConstraintInterface, err error) {
	quotes, brackets := 0, 0
	tokenStart := 0

	tag += ","
	for i, c := range tag {
		switch {
		case c == '[' && quotes%2 == 0:
			brackets++
		case c == ']' && quotes%2 == 0:
			brackets--
		case c == '\'' && i > 0 && tag[i-1] != '\\':
			quotes++
		case c == ',' && brackets == 0 && quotes%2 == 0:
			tn, tv, tp := tag[tokenStart:i], "", ""
			if p := strings.IndexByte(tn, '='); p != -1 {
				tn, tv = tn[:p], tn[p+1:]
			}
			ci, ok := TagOptions[tn]
			if !ok {
				return nil, fmt.Errorf(`gowl/validator: unknown tag option "%s"`, tn)
			}
			if tv != "" && tv[0] != '[' && tv[0] != '\'' {
				// extract value prefix ("u8:", "f32:", etc)
				if p := strings.IndexByte(tv, ':'); p != -1 {
					tv, tp = tv[p+1:], tv[:p]
				}
			}
			var constraint ConstraintInterface
			if constraint, err = ci(TagOption{tn, tv, tp}); err != nil {
				return nil, err
			}
			constraints = append(constraints, constraint)
			tokenStart = i + 1
		}
	}
	return
}

var TagOptions = map[string]func(o TagOption) (constraint ConstraintInterface, err error){
	"required": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithBool(func(v bool) ConstraintInterface { return Required(v) })
	},
	"valid": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithBool(func(v bool) ConstraintInterface { return Valid(v) })
	},
	"type": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithString(func(v string) ConstraintInterface { return Type(v) })
	},
	"len": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint64) ConstraintInterface { return Length(v) })
	},
	"minlen": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint64) ConstraintInterface { return MinLength(v) })
	},
	"maxlen": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithUint(func(v uint64) ConstraintInterface { return MaxLength(v) })
	},
	"eq": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithValue(func(v interface{}) ConstraintInterface { return Equal(v) })
	},
	"true": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.NoValue(func() ConstraintInterface { return True })
	},
	"false": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.NoValue(func() ConstraintInterface { return False })
	},
	"lt": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithValue(func(v interface{}) ConstraintInterface { return LessThan(v) })
	},
	"lte": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithValue(func(v interface{}) ConstraintInterface { return LessThanOrEqual(v) })
	},
	"gt": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithValue(func(v interface{}) ConstraintInterface { return GreaterThan(v) })
	},
	"gte": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithValue(func(v interface{}) ConstraintInterface { return GreaterThanOrEqual(v) })
	},
	"re": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithString(func(v string) ConstraintInterface { return Regexp(v) })
	},
	"repsx": func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithString(func(v string) ConstraintInterface { return RegexpPOSIX(v) })
	},
}

func init() {
	// avoid initialization loop
	TagOptions["each"] = func(o TagOption) (constraint ConstraintInterface, err error) {
		return o.WithTag(func(v []ConstraintInterface) ConstraintInterface { return Each(v) })
	}
}
