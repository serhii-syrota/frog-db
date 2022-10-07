// Package dbtypes provides supported db types
package dbtypes

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type Type string

var AvailableTypes = []Type{Integer, Real, Char, String, RealInv, Image}

func IsAvailableType(val string) bool {
	for _, t := range AvailableTypes {
		if string(t) == val {
			return true
		}
	}
	return false
}

const (
	Integer Type = "integer"
	Real    Type = "real"
	Char    Type = "char"
	String  Type = "string"
	RealInv Type = "realInv"
	Image   Type = "image"
)

func NewInteger(v any) (*int64, error) {
	res, err := cast.ToInt64E(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func NewReal(v any) (*float64, error) {
	res, err := cast.ToFloat64E(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func NewChar(v any) (*rune, error) {
	var val rune
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.String:
		if len(value.String()) != 1 {
			return nil, fmt.Errorf("%s must contain exact 1 symbol", v)
		}
		for _, runeValue := range value.String() {
			val = runeValue
		}
	case reflect.Int32:
		val = rune(value.Int())
	default:
		return nil, fmt.Errorf("unknown type provided for rune, %T", v)
	}
	return &val, nil
}

func NewString(v any) (*string, error) {
	res, err := cast.ToStringE(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type TRealInv struct {
	A float64
	B float64
}

func NewRealInv(a, b any) (*TRealInv, error) {
	aVal, err := NewReal(a)
	if err != nil {
		return nil, err
	}
	bVal, err := NewReal(b)
	if err != nil {
		return nil, err
	}
	if *aVal > *bVal {
		return nil, errs.NewErrInvalidRange(*aVal, *bVal)
	}

	return &TRealInv{*aVal, *bVal}, nil
}
