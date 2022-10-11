package table

import (
	"fmt"
	"hash/fnv"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/elliotchance/pie/v2"
	dbtypes "github.com/ssyrota/frog-db/src/core/db/dbtypes"
	"github.com/ssyrota/frog-db/src/core/db/deepcopy"
	"github.com/ssyrota/frog-db/src/core/db/schema"
	errs "github.com/ssyrota/frog-db/src/core/err"
	"golang.org/x/exp/slices"
)

type ColumnSet map[string]any

// Validate schema and create new table
func NewTable(sch schema.T) (*T, error) {
	for column, t := range sch {
		if !dbtypes.IsAvailableName(string(t)) {
			return nil, errs.NewErrInvalidTypeProvided(column, string(t))
		}
	}
	return &T{schema: sch}, nil
}

type T struct {
	mu     sync.RWMutex
	schema schema.T
	data   []ColumnSet
}

// Dump table.
type Dump struct {
	Schema schema.T    `json:"schema"`
	Data   []ColumnSet `json:"data"`
	Name   string      `json:"name"`
}

func (t *T) Dump(tableName string) (*Dump, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var dump Dump
	dump.Data = t.data
	dump.Schema = t.schema
	dump.Name = tableName
	return &dump, nil
}

func (t *T) LoadDump(data *[]ColumnSet) error {
	inserted, err := t.InsertRows(data)
	log.Printf("%v", len(*data) == int(inserted))
	return err
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
	rowsToInsert := make([]ColumnSet, len(*rows))
	requiredColumns := MapKeys(t.schema)
	for i, row := range *rows {
		rowColumns := MapKeys(row)
		// Check required columns
		omitted := pie.Filter(requiredColumns, func(a string) bool {
			return !slices.Contains(rowColumns, a)
		})
		if len(omitted) != 0 {
			return 0, errs.NewErrColumnsRequired(omitted)
		}
		// Check extra columns
		extra := pie.Filter(rowColumns, func(a string) bool {
			return !slices.Contains(requiredColumns, a)
		})
		if len(extra) != 0 {
			return 0, errs.NewErrColumnsNotFound(extra)
		}
		// Validate types
		rowToInsert := make(ColumnSet)
		for k, v := range row {
			typedVal, err := dbtypes.NewDataVal(t.schema[k], v)
			if err != nil {
				return 0, err
			}
			rowToInsert[k] = typedVal
		}
		rowsToInsert[i] = rowToInsert
	}
	t.data = append(t.data, rowsToInsert...)
	return uint(len(rowsToInsert)), nil
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

// Update rows from table
func (t *T) DeleteDuplicates() (uint, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	columnNames := make([]string, 0, len(t.schema))
	for k := range t.schema {
		columnNames = append(columnNames, k)
	}
	uniqueSet := Set{}
	deleted := uint(0)
	for i, row := range t.data {
		orderedColumnsStr := strings.Builder{}
		for _, columnName := range columnNames {
			_, err := orderedColumnsStr.WriteString(fmt.Sprint(row[columnName]))
			if err != nil {
				return 0, err
			}
		}
		rowHash := hash(orderedColumnsStr.String())
		if _, ok := uniqueSet[fmt.Sprint(rowHash)]; ok {
			t.data = append(t.data[:i], t.data[i+1:]...)
			deleted++
		} else {
			uniqueSet[fmt.Sprint(rowHash)] = void{}
		}
	}
	return deleted, nil
}

// Select data from table,
// empty columns list and empty conditions considered as "select all"
func (t *T) SelectRows(columns *[]string, conditions ColumnSet) (*[]ColumnSet, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	filteredIDs, err := t.filter(conditions)
	if err != nil {
		return nil, err
	}
	res := make([]ColumnSet, len(*filteredIDs))
	for i, id := range *filteredIDs {
		row, err := t.removeExtraFields(t.data[id], columns)
		if err != nil {
			return nil, err
		}
		res[i] = row
	}
	return &res, nil
}
func (t *T) removeExtraFields(row ColumnSet, requiredColumns *[]string) (ColumnSet, error) {
	copied, err := deepcopy.Map(row)
	if err != nil {
		return nil, err
	}
	if len(*requiredColumns) != 0 {
		for columnName := range copied {
			if !slices.Contains(*requiredColumns, columnName) {
				delete(copied, columnName)
			}
		}
	}
	return copied, nil
}

// Get indexes from data, that matches conditions.
// when condition is empty filter returns all data ids from table
func (t *T) filter(rawCondition ColumnSet) (*[]int, error) {
	condition, err := t.setFromRaw(rawCondition)
	if err != nil {
		return nil, err
	}
	res := []int{}
rows:
	for i, row := range t.data {
		for k, v := range condition {
			if !reflect.DeepEqual(row[k], v) {
				continue rows
			}
		}
		res = append(res, i)
	}
	return &res, nil
}

// Convert raw map to typed ColumnSet
func (t *T) setFromRaw(raw ColumnSet) (ColumnSet, error) {
	typedSet := make(ColumnSet, len(raw))
	for k, v := range raw {
		dataType, ok := t.schema[k]
		if !ok {
			return nil, errs.NewErrColumnsNotFound([]string{k})
		}
		val, err := dbtypes.NewDataVal(dataType, v)
		if err != nil {
			return nil, err
		}
		typedSet[k] = val
	}
	return typedSet, nil
}

// Create list from old copy without unwanted indexes
func removeIndexes[T any](slice []T, ids []int) []T {
	result := make([]T, len(slice)-len(ids))
	counter := 0
	for id, v := range slice {
		if !slices.Contains(ids, id) {
			result[counter] = v
			counter++
		}
	}
	return result
}

func MapKeys[T any](input map[string]T) []string {
	result := make([]string, 0, len(input))
	for k := range input {
		result = append(result, k)
	}
	return result

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type Set map[string]void
type void struct{}
