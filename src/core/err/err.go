// Package errs provides errors, that may occur on db usage
package errs

import "fmt"

type ErrTableAlreadyExists struct {
	error
}

func NewErrTableAlreadyExists(tableName string) *ErrTableAlreadyExists {
	return &ErrTableAlreadyExists{fmt.Errorf("table %s already exists", tableName)}
}

type ErrTableNotFound struct {
	error
}

func NewErrTableNotFound(tableName string) *ErrTableNotFound {
	return &ErrTableNotFound{fmt.Errorf("table %s not found", tableName)}
}

type ErrColumnNotFound struct {
	error
}

func NewErrColumnNotFound(columnName string) *ErrColumnNotFound {
	return &ErrColumnNotFound{fmt.Errorf("column %s not found", columnName)}
}

type ErrNoColumns struct {
	error
}

func NewErrNoColumns() *ErrNoColumns {
	return &ErrNoColumns{fmt.Errorf("cannot create table without columns")}
}

type ErrInvalidTypeProvided struct {
	error
}

func NewErrInvalidTypeProvided(columnName, t string) *ErrInvalidTypeProvided {
	return &ErrInvalidTypeProvided{fmt.Errorf("cannot create column %s with type %s", columnName, t)}
}

type ErrInvalidRange struct {
	error
}

func NewErrInvalidRange(a, b float64) *ErrInvalidRange {
	return &ErrInvalidRange{fmt.Errorf("invalid range %v>%v", a, b)}
}

type ErrInvalidRangeDeclaration struct {
	error
}

func NewErrInvalidRangeDeclaration() *ErrInvalidRangeDeclaration {
	return &ErrInvalidRangeDeclaration{fmt.Errorf("invalid range declaration, should be provided as map with \"from\",\"to\" fields")}
}

type ErrDbIO struct {
	error
}

func NewErrDbIO(err error) *ErrDbIO {
	return &ErrDbIO{fmt.Errorf("db io error: %s", err.Error())}
}
