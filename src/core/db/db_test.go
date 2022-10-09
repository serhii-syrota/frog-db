package db

import (
	"fmt"
	"testing"

	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
	dbtypes "github.com/ssyrota/frog-db/src/core/db/types"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	validTableSchema := schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}
	tableName := "frog"
	validCreateCommand := &CommandCreateTable{tableName, validTableSchema}
	invalidSchema := schema.T{"invalid_type_column": "unknown_type"}

	t.Run("fails on unknown command type", func(t *testing.T) {
		db, err := New("")
		assert.Nil(t, err)
		assert.NotNil(t, db)
		_, err = db.Execute("unknown smooth command")
		assert.NotNil(t, err)
		assert.EqualError(t, err, fmt.Sprintf("unknown command type: %T", "unknown smooth command"))
	},
	)

	t.Run("CreateTable with IntrospectSchema", func(t *testing.T) {
		t.Run(
			"accepts schema with valid data types and with introspect returns provided schema",
			func(t *testing.T) {
				db, err := New("")
				assert.Nil(t, err)
				assert.NotNil(t, db)
				createRes, err := db.Execute(validCreateCommand)
				assert.Nil(t, err)
				assert.Equal(t, (*createRes)[0]["message"], fmt.Sprintf("successfully created table %s", tableName))
				introspectionRes, err := db.IntrospectSchema(tableName)
				assert.Nil(t, err)
				assert.Equal(t, introspectionRes, validTableSchema)
			},
		)

		t.Run(
			"fails on create table with duplicate name",
			func(t *testing.T) {
				db, err := New("")
				assert.Nil(t, err)
				assert.NotNil(t, db)
				_, err = db.Execute(validCreateCommand)
				assert.Nil(t, err)
				_, err = db.Execute(validCreateCommand)
				assert.NotNil(t, err)
				assert.EqualError(t, err, fmt.Sprintf("table %s already exists", tableName))
			},
		)

		t.Run(
			"fails on invalid dataType in schema provided",
			func(t *testing.T) {
				db, err := New("")
				assert.Nil(t, err)
				assert.NotNil(t, db)
				_, err = db.Execute(&CommandCreateTable{"frog", invalidSchema})
				assert.NotNil(t, err)
				assert.EqualError(t, err, fmt.Sprintf("cannot create column %s with type %s", "invalid_type_column", "unknown_type"))
			},
		)
	})

	t.Run("DropTable", func(t *testing.T) {
		t.Run("drops existed table", func(t *testing.T) {
			db, err := New("")
			assert.Nil(t, err)
			assert.NotNil(t, db)
			_, err = db.Execute(validCreateCommand)
			assert.Nil(t, err)
			existedTable, err := db.table(tableName)
			assert.Nil(t, err)
			assert.NotNil(t, existedTable)
			dropResult, err := db.Execute(&CommandDropTable{"frog"})
			assert.Nil(t, err)
			assert.Equal(t, (*dropResult)[0]["message"], "successfully dropped table frog")
			removedTable, err := db.table(tableName)
			assert.Nil(t, removedTable)
			assert.NotNil(t, err)
		})

		t.Run("fails on drop non existed table", func(t *testing.T) {
			db, err := New("")
			assert.Nil(t, err)
			assert.NotNil(t, db)
			dropResult, err := db.Execute(&CommandDropTable{"frog"})
			assert.Nil(t, dropResult)
			assert.NotNil(t, err)
			assert.EqualError(t, err, "table frog not found")
		})
	})

	t.Run("Insert row", func(t *testing.T) {
		db, _ := New("")
		db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
		rows := &[]table.ColumnSet{{"leg_length": []float64{2.2, 3.3}}}
		insertResult, err := db.Execute(&CommandInsert{"frog", rows})
		assert.Nil(t, err)
		assert.NotNil(t, insertResult)
		assert.Equal(t, fmt.Sprintf("successfully inserted %d row to table frog", len(*rows)), (*insertResult)[0]["message"])

		selectResult, err := db.Execute(&CommandSelect{"frog", &[]string{}, table.ColumnSet{}})
		assert.Nil(t, err)
		assert.NotNil(t, selectResult)
		assert.Equal(t, []table.ColumnSet{{"leg_length": []float64{2.2, 3.3}}}, (*selectResult))
	})

}
