// Package errs provides errors, that may occur on db usage
package errs

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize/english"
)

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

type ErrColumnsRequired struct {
	error
}

func NewErrColumnsRequired(columnNames []string) *ErrColumnsRequired {
	return &ErrColumnsRequired{
		fmt.Errorf("%s %s required",
			english.PluralWord(len(columnNames), "column", ""),
			strings.Join(columnNames, ", ")),
	}
}

type ErrColumnsNotFound struct {
	error
}

func NewErrColumnsNotFound(columnNames []string) *ErrColumnsNotFound {
	return &ErrColumnsNotFound{
		fmt.Errorf("%s %s not found",
			english.PluralWord(len(columnNames), "column", ""),
			strings.Join(columnNames, ", ")),
	}
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
