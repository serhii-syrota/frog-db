package db

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
	dbtypes "github.com/ssyrota/frog-db/src/core/db/types"
	errs "github.com/ssyrota/frog-db/src/core/err"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	dumpPath := ".test.json"
	os.Setenv("DUMP_PATH", dumpPath)
	defer os.Remove(dumpPath)

	validTableSchema := schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}
	tableName := "frog"
	validCreateCommand := &CommandCreateTable{tableName, validTableSchema}
	invalidSchema := schema.T{"invalid_type_column": "unknown_type"}

	t.Run("fails on unknown command type", func(t *testing.T) {
		db, err := New(dumpPath, time.Second)
		assert.Nil(t, err)
		assert.NoError(t, err, "")
		_, err = db.Execute("unknown smooth command")
		assert.NotNil(t, err)
		assert.EqualError(t, err, fmt.Sprintf("unknown command type: %T", "unknown smooth command"))
	},
	)

	t.Run("CreateTable with IntrospectSchema", func(t *testing.T) {
		t.Run("accepts schema with valid data types and with introspect returns provided schema",
			func(t *testing.T) {
				db, err := New(dumpPath, time.Second)
				assert.Nil(t, err)
				assert.NotNil(t, db)
				createRes, err := db.Execute(validCreateCommand)
				assert.Nil(t, err)
				assert.Equal(t, (*createRes)[0]["message"], fmt.Sprintf("successfully created table %s", tableName))
				introspectionRes, err := db.IntrospectSchema()
				assert.Nil(t, err)
				assert.Equal(t, introspectionRes[tableName], validTableSchema)
			},
		)
		t.Run("fails on create table with duplicate name",
			func(t *testing.T) {
				db, err := New(dumpPath, time.Second)
				assert.Nil(t, err)
				assert.NotNil(t, db)
				_, err = db.Execute(validCreateCommand)
				assert.Nil(t, err)
				_, err = db.Execute(validCreateCommand)
				assert.NotNil(t, err)
				assert.EqualError(t, err, fmt.Sprintf("table %s already exists", tableName))
			},
		)
		t.Run("fails on invalid dataType in schema provided",
			func(t *testing.T) {
				db, err := New(dumpPath, time.Second)
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
			db, err := New(dumpPath, time.Second)
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
			db, err := New(dumpPath, time.Second)
			assert.Nil(t, err)
			assert.NotNil(t, db)
			dropResult, err := db.Execute(&CommandDropTable{"frog"})
			assert.Nil(t, dropResult)
			assert.NotNil(t, err)
			assert.EqualError(t, err, "table frog not found")
		})
	})

	t.Run("Insert", func(t *testing.T) {
		t.Run("accepts and save input with required columns and valid types", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			insertResult, err := db.Execute(&CommandInsert{"frog", rows})
			assert.Nil(t, err)
			assert.Equal(t, fmt.Sprintf("successfully inserted %d rows to table frog", len(*rows)), (*insertResult)[0]["message"])

			selectResult, err := db.Execute(&CommandSelect{"frog", &[]string{}, table.ColumnSet{}})
			assert.Nil(t, err)
			assert.Equal(t, *rows, (*selectResult))
		})
		t.Run("fail input without required columns", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{{"leg_length": 1}}
			_, err := db.Execute(&CommandInsert{"frog", rows})
			assert.NotNil(t, err)
			assert.IsType(t, &errs.ErrColumnsRequired{}, err)
		})
		t.Run("fail input with unexpected columns", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"unknown": 1, "leg_length": 2, "jump": []float64{2.5, 3.5}}}
			_, err := db.Execute(&CommandInsert{"frog", rows})
			assert.NotNil(t, err)
			assert.IsType(t, &errs.ErrColumnsNotFound{}, err)
		})
		t.Run("fail input with columns type mismatch", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": "short", "jump": []float64{2.5, 3.5}}}
			_, err := db.Execute(&CommandInsert{"frog", rows})
			assert.NotNil(t, err)
			assert.Error(t, errors.New(""), err)
		})
	})

	t.Run("Select", func(t *testing.T) {
		t.Run("accepts valid conditions and fields and return data, that matches conditions", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			db.Execute(&CommandInsert{"frog", rows})
			selectResult, err := db.Execute(&CommandSelect{"frog", &[]string{"jump"}, table.ColumnSet{"leg_length": 1}})
			assert.Nil(t, err)
			assert.NotNil(t, selectResult)
			assert.Equal(t, selectResult, &[]table.ColumnSet{{"jump": []float64{2.2, 3.3}}})
		})
		t.Run("is idempotent", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			db.Execute(&CommandInsert{"frog", rows})
			for i := 0; i < 1000; i++ {
				selectResult, _ := db.Execute(&CommandSelect{"frog", &[]string{"jump"}, table.ColumnSet{"leg_length": 1}})
				assert.Equal(t, selectResult, &[]table.ColumnSet{{"jump": []float64{2.2, 3.3}}})
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("fail on invalid update data", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			db.Execute(&CommandInsert{"frog", rows})

			updateFields := table.ColumnSet{"jump": []float64{10, 9}}
			_, err := db.Execute(&CommandUpdate{"frog", table.ColumnSet{"leg_length": 1}, updateFields})
			assert.NotNil(t, err)
		})
		t.Run("accepts valid conditions and updates table rows", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			tableName := "frog"
			db.Execute(&CommandCreateTable{tableName, schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			db.Execute(&CommandInsert{tableName, rows})

			updateConditions := table.ColumnSet{"leg_length": 1}
			updateFields := table.ColumnSet{"jump": []float64{10, 11}}
			updateResult, err := db.Execute(&CommandUpdate{tableName, updateConditions, updateFields})
			assert.Nil(t, err)
			assert.NotNil(t, updateResult)
			assert.Equal(t, &[]table.ColumnSet{{"message": "successfully updated 1 row in table frog"}}, updateResult)
			selectResult, _ := db.Execute(&CommandSelect{tableName, &[]string{"jump"}, updateConditions})
			assert.Equal(t, selectResult, &[]table.ColumnSet{{"jump": []float64{10, 11}}})
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("delete data by valid conditions", func(t *testing.T) {
			db, _ := New(dumpPath, time.Second)
			db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
			rows := &[]table.ColumnSet{
				{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
				{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}
			db.Execute(&CommandInsert{"frog", rows})

			deleteRes, err := db.Execute(&CommandDelete{"frog", table.ColumnSet{"leg_length": float64(1)}})
			assert.NoError(t, err)
			assert.Equal(t, &[]table.ColumnSet{{"message": "successfully deleted 1 row from table frog"}}, deleteRes)

			selectResult, err := db.Execute(&CommandSelect{"frog", &[]string{}, table.ColumnSet{}})
			assert.NoError(t, err)
			assert.Equal(t, []table.ColumnSet{{"leg_length": float64(2), "jump": []float64{2.5, 3.5}}}, (*selectResult))
		})
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		db, _ := New(dumpPath, time.Second)
		db.Execute(&CommandCreateTable{"frog", schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
		rows := &[]table.ColumnSet{
			{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
			{"leg_length": float64(1), "jump": []float64{2.2, 3.3}},
			{"leg_length": float64(1), "jump": []float64{2.5, 3.5}}}
		db.Execute(&CommandInsert{"frog", rows})
		deleteRes, err := db.Execute(&CommandRemoveDuplicates{"frog"})
		assert.NoError(t, err)
		assert.Equal(t, &[]table.ColumnSet{{"message": "successfully deleted 1 row from table frog"}}, deleteRes)
	})
}

// Test save dump to file and create db from dump.
func TestDump(t *testing.T) {
	t.Run("save and upload", func(t *testing.T) {
		dumpPath := ".test_dump.json"
		os.Setenv("DUMP_PATH", dumpPath)
		database, err := New(dumpPath, time.Second)
		defer os.Remove(dumpPath)

		assert.NoError(t, err)
		tables := []string{"frog", "leg"}
		database.Execute(&CommandCreateTable{Name: tables[0], Schema: schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})
		database.Execute(&CommandCreateTable{Name: tables[1], Schema: schema.T{"leg_length": dbtypes.Real, "jump": dbtypes.RealInv}})

		rows := &[]table.ColumnSet{{"leg_length": float64(1), "jump": []float64{2.2, 3.3}}}
		for i := 0; i < 10; i++ {
			_, err = database.Execute(&CommandInsert{To: "frog", Data: rows})
			assert.NoError(t, err)
		}
		for i := 0; i < 10; i++ {
			_, err = database.Execute(&CommandInsert{To: "leg", Data: rows})
			assert.NoError(t, err)
		}
		err = database.StoreDump()
		assert.NoError(t, err)

		os.Setenv("DUMP_PATH", ".new_dump.json")
		defer os.Remove(".new_dump.json")
		newDb, err := New(dumpPath, time.Second)
		assert.NoError(t, err)
		err = newDb.FromDump(dumpPath)
		assert.NoError(t, err)
		selectRes, err := newDb.Execute(&CommandSelect{"frog", &[]string{}, make(table.ColumnSet)})
		assert.NoError(t, err)
		assert.Equal(t, 10, len(*selectRes))
	})
}
