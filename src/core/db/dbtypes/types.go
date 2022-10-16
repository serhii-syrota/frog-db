// Package dbtypes provides supported db types
package dbtypes

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type Type string

var AvailableTypeNames = []Type{Integer, Real, Char, String, RealInv, Image}

func IsAvailableName(val string) bool {
	for _, t := range AvailableTypeNames {
		if string(t) == val {
			return true
		}
	}
	return false
}

// Parse and return pointer to parsed value
func NewDataVal(dataType Type, val any) (any, error) {
	switch dataType {
	case Integer:
		return NewInteger(val)
	case Real:
		return NewReal(val)
	case Char:
		return NewChar(val)
	case String:
		return NewString(val)
	case RealInv:
		return NewRealInv(val)
	case Image:
		return NewImage(val)
	default:
		return nil, fmt.Errorf("%s is invalid data type", dataType)
	}
}

const (
	Integer Type = "integer"
	Real    Type = "real"
	Char    Type = "char"
	String  Type = "string"
	RealInv Type = "realInv"
	Image   Type = "image"
)

func NewInteger(v any) (int64, error) {
	res, err := cast.ToInt64E(v)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func NewReal(v any) (float64, error) {
	res, err := cast.ToFloat64E(v)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func NewChar(v any) (rune, error) {
	var val rune
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.String:
		if len(value.String()) != 1 {
			return 0, fmt.Errorf("%s must contain exact 1 symbol", v)
		}
		for _, runeValue := range value.String() {
			val = runeValue
		}
	case reflect.Int32:
		val = rune(value.Int())
	default:
		return 0, fmt.Errorf("unknown type provided for rune, %T", v)
	}
	return val, nil
}

func NewString(v any) (string, error) {
	res, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}
	return res, nil
}

func NewImage(v any) (string, error) {
	res, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}
	return res, nil
}

func NewRealInv(val any) ([]float64, error) {
	data := make([]any, 2)
	switch typedVal := val.(type) {
	case []float64:
		if len(typedVal) != 2 {
			return nil, errs.NewErrInvalidRangeDeclaration()
		}
		data[0] = typedVal[0]
		data[1] = typedVal[1]
	case []any:
		if len(typedVal) != 2 {
			return nil, errs.NewErrInvalidRangeDeclaration()
		}
		data[0] = typedVal[0]
		data[1] = typedVal[1]
	default:
		return nil, errs.NewErrInvalidRangeDeclaration()
	}
	aVal, err := NewReal(data[0])
	if err != nil {
		return nil, err
	}
	bVal, err := NewReal(data[1])
	if err != nil {
		return nil, err
	}
	if aVal > bVal {
		return nil, errs.NewErrInvalidRange(aVal, bVal)
	}
	return []float64{aVal, bVal}, nil
}
