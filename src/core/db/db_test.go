package db

import (
	"testing"

	"github.com/ssyrota/frog-db/src/core/db/schema"
	errs "github.com/ssyrota/frog-db/src/core/err"
	"github.com/tj/assert"
)

// Test db table creation.
func TestTableCreate(t *testing.T) {
	t.Run("accepts and save schema with available types", func(t *testing.T) {
		db := NewDb()
		tableName := "frog"
		fields := map[string]string{"id": "integer", "name": "string"}
		sch, err := schema.New(fields)
		assert.Nil(t, err)
		err = db.CreateTable(tableName, sch)
		assert.Nil(t, err)
		resultSchema, err := db.TableSchema(tableName)
		assert.Nil(t, err)
		for k, resT := range resultSchema {
			assert.EqualValues(t, fields[k], string(resT))
		}
	})

	t.Run("fail on create already existed table", func(t *testing.T) {
		db := NewDb()
		tableName := "frog"
		fields := map[string]string{"id": "integer", "name": "string"}
		sch, err := schema.New(fields)
		assert.Nil(t, err)
		err = db.CreateTable(tableName, sch)
		assert.Nil(t, err)
		err = db.CreateTable(tableName, sch)
		assert.IsType(t, &errs.ErrTableAlreadyExists{}, err)
	})
}
