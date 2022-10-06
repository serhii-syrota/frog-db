// Package types provides supported db types
package types

import (
	"reflect"
	"strconv"

	"github.com/ssyrota/frog-db/src/core/errs"
)

// DataType provides validation and matching.
type DataType interface {
	Validate() error
	IsIt() bool
}

type Integer struct {
	val int64
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
			return nil, errs.NewErrValueTypeMismatch(columnName, "int64", err.Error())
		} else {
			val = parsed
		}
	default:
		return nil, errs.NewErrValueTypeMismatch(columnName, "int64", v)

	}
	return &Integer{val}, nil
}
func (i *Integer) Value() int64 {
	return i.val
}
