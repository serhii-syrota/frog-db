package table

import (
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/ssyrota/frog-db/src/core/db/schema"
	dbtypes "github.com/ssyrota/frog-db/src/core/db/types"
	errs "github.com/ssyrota/frog-db/src/core/err"
	"golang.org/x/exp/slices"
)

type ColumnSet map[string]any

func NewTable(sch schema.T) (*T, error) {
	for column, t := range sch {
		if !dbtypes.IsAvailableName(string(t)) {
			return nil, errs.NewErrInvalidTypeProvided(column, string(t))
		}
	}
	return &T{schema: sch}, nil
}

// Table with schema, data and crud commands
type T struct {
	mu     sync.RWMutex
	schema schema.T
	data   []ColumnSet
}

// Introspect schema
func (t *T) Schema() schema.T {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.schema
}

// Insert rows to table
func (t *T) InsertRows(rows *[]ColumnSet) (uint, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data = append(t.data, *rows...)
	return uint(len(*rows)), nil
}

// Update rows in table
func (t *T) UpdateRows(rawCondition ColumnSet, newRawData ColumnSet) (uint, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	newData, err := t.setFromRaw(newRawData)
	if err != nil {
		return 0, err
	}
	ids, err := t.filter(rawCondition)
	if err != nil {
		return 0, err
	}
	for _, v := range *ids {
		rawToUpdate := t.data[v]
		for column, updatedValue := range newData {
			rawToUpdate[column] = updatedValue
		}

	}
	return uint(len(*ids)), nil
}

// Update rows from table
func (t *T) DeleteRows(rawCondition ColumnSet) (uint, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	ids, err := t.filter(rawCondition)
	if err != nil {
		return 0, err
	}
	t.data = removeIndexes(t.data, *ids)
	return uint(len(*ids)), nil
}

// Select data from table
func (t *T) SelectRows(rows *[]string, conditions ColumnSet) (*[]ColumnSet, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	filtered, err := t.filter(conditions)
	if err != nil {
		return nil, err
	}
	res := make([]ColumnSet, len(*filtered))
	for _, i := range *filtered {
		v := t.data[i]
		// Filter only required fields
		for columnName := range v {
			if !slices.Contains(*rows, columnName) {
				delete(v, columnName)
			}
		}
		res = append(res, v)
	}

	return &res, nil

}

// Get filtered data indexes from table
func (t *T) filter(rawCondition ColumnSet) (*[]int, error) {
	condition, err := t.setFromRaw(rawCondition)
	if err != nil {
		return nil, err
	}
	res := []int{}
rows:
	for i, row := range t.data {
		for k, v := range condition {
			if !cmp.Equal(row[k], v) {
				continue rows
			}
		}
		res = append(res, i)
	}
	return &res, nil
}

func (t *T) setFromRaw(raw ColumnSet) (ColumnSet, error) {
	condition := make(ColumnSet, len(raw))
	for k, v := range raw {
		dataType, ok := t.schema[k]
		if !ok {
			return nil, errs.NewErrColumnNotFound(k)
		}
		val, err := dbtypes.NewDataVal(dataType, v)
		if err != nil {
			return nil, err
		}
		condition[k] = val
	}
	return condition, nil
}

func removeIndexes[T any](slice []T, ids []int) []T {
	result := make([]T, len(slice)-len(ids))

	counter := 0
	for id, v := range slice {
		if slices.Contains(ids, id) {
			result[counter] = v
			counter++
		}
	}
	return result
}
