// Package types provides supported db types
package types

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
	"github.com/ssyrota/frog-db/src/core/errs"
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

type RealInv struct {
	A float64
	B float64
}

func NewRealInv(a, b any) (*RealInv, error) {
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

	return &RealInv{*aVal, *bVal}, nil
}
