// Package types provides supported db types
package types

import (
	"reflect"
	"strconv"

	"github.com/ssyrota/frog-db/src/core/errs"
)

type Integer struct {
	Val int64
}

func NewInteger(v any, columnName string) (*Integer, error) {
	var val int64
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = int64(value.Uint())
	case reflect.Float32, reflect.Float64:
		val = int64(value.Float())
	case reflect.String:
		if parsed, err := strconv.ParseInt(value.String(), 10, 64); err != nil {
			return nil, errs.NewErrConvertStringToNum(columnName, "int64", value, err)
		} else {
			val = parsed
		}
	default:
		return nil, errs.NewErrValueTypeMismatch(columnName, "int64", v)

	}
	return &Integer{val}, nil
}

type Real struct {
	Val float64
}

func NewReal(v any, columnName string) (*Real, error) {
	var val float64
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = float64(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = float64(value.Uint())
	case reflect.Float32, reflect.Float64:
		val = value.Float()
	case reflect.String:
		if parsed, err := strconv.ParseFloat(value.String(), 64); err != nil {
			return nil, errs.NewErrConvertStringToNum(columnName, "float64", value, err)
		} else {
			val = parsed
		}
	default:
		return nil, errs.NewErrValueTypeMismatch(columnName, "float64", v)

	}
	return &Real{val}, nil
}

type Char struct {
	Val rune
}

func NewChar(v any, columnName string) (*Char, error) {
	var val rune
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.String:
		if len(value.String()) != 1 {
			return nil, errs.NewErrValueTypeMismatch(columnName, "rune", v)
		}
		for _, runeValue := range value.String() {
			val = runeValue
		}
	case reflect.Int32:
		val = rune(value.Int())
	default:
		return nil, errs.NewErrValueTypeMismatch(columnName, "rune", v)
	}
	return &Char{val}, nil
}

type String struct {
	Val string
}

func NewString(v any, columnName string) (*String, error) {
	var val string
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.String:
		val = value.String()
	case reflect.Int32:
		val = string(rune(value.Int()))
	default:
		return nil, errs.NewErrValueTypeMismatch(columnName, "string", v)
	}
	return &String{val}, nil
}

type RealInv struct {
	A float64
	B float64
}

func NewRealInv(a, b any, columnName string) (*RealInv, error) {
	aVal, err := NewReal(a, columnName)
	if err != nil {
		return nil, err
	}
	bVal, err := NewReal(b, columnName)
	if err != nil {
		return nil, err
	}
	if aVal.Val > bVal.Val {
		return nil, errs.NewErrInvalidRange(columnName, aVal.Val, bVal.Val)
	}

	return &RealInv{aVal.Val, bVal.Val}, nil
}
