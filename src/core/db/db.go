package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type Dump []table.Dump
type Db interface {
	Execute(command any) (*[]table.ColumnSet, error)
	IntrospectSchema() (map[string]schema.T, error)
	StoreDump() error
	JsonDump() <-chan DumpMsg
	FromDump(dumpPath string) error
}
type Database struct {
	tables map[string]*table.T
	path   string
}

func New(path string, dumpInterval time.Duration) (*Database, error) {
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	db := &Database{tables: make(map[string]*table.T), path: path}
	// Run store dump interval job
	go func() {
		ticker := time.NewTicker(dumpInterval)
		for range ticker.C {
			err := db.StoreDump()
			if err != nil {
				log.Printf("error: %s", err.Error())
			}
		}
	}()
	return db, nil
}

var _ Db = new(Database)

// FromDump implementation.
func (db *Database) FromDump(dumpPath string) error {
	// Save dump before delete data
	err := db.StoreDump()
	if err != nil {
		return err
	}
	dumpRaw, err := os.ReadFile(dumpPath)
	if err != nil {
		return err
	}
	var dump Dump
	err = json.Unmarshal(dumpRaw, &dump)
	if err != nil {
		return err
	}
	// Clean up tables if exists
	db.tables = make(map[string]*table.T)
	for _, dumpTable := range dump {
		if _, err := db.Execute(&CommandCreateTable{dumpTable.Name, dumpTable.Schema}); err != nil {
			return err
		}
		storedTable, err := db.table(dumpTable.Name)
		if err != nil {
			return err
		}
		if err := storedTable.LoadDump(&dumpTable.Data); err != nil {
			return err
		}
	}
	return nil
}

// StoreDump implementation.
func (db *Database) StoreDump() error {
	file, err := os.OpenFile(db.path, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	dumpCh := db.JsonDump()
	for value := range dumpCh {
		if value.Err != nil {
			return err
		}
		_, err := writer.Write(value.Payload)
		if err != nil {
			return err
		}
		writer.Flush()
	}
	return nil
}

type DumpMsg struct {
	Payload []byte
	Err     error
}

// JsonDump implementation.
func (db *Database) JsonDump() <-chan DumpMsg {
	ch := make(chan DumpMsg)
	go func() {
		ch <- DumpMsg{[]byte("["), nil}
		tableNames := table.MapKeys(db.tables)
		for i, tableName := range tableNames {
			dump, err := db.tables[tableName].Dump(tableName)
			if err != nil {
				ch <- DumpMsg{nil, err}
				close(ch)
				return
			}
			bytes, err := json.Marshal(dump)
			if err != nil {
				ch <- DumpMsg{nil, err}
				close(ch)
				return
			}
			ch <- DumpMsg{bytes, nil}
			if i != len(tableNames)-1 {
				ch <- DumpMsg{[]byte(","), nil}
			}
		}

		ch <- DumpMsg{[]byte("]"), nil}
		close(ch)
	}()
	return ch
}

// IntrospectSchema implementation.
func (db *Database) IntrospectSchema() (map[string]schema.T, error) {
	dbSchema := map[string]schema.T{}
	for k, t := range db.tables {
		dbSchema[k] = t.Schema()
	}
	return dbSchema, nil
}

// Execute implementation.
func (db *Database) Execute(command any) (*[]table.ColumnSet, error) {
	switch typedCommand := command.(type) {
	case *CommandDropTable:
		return db.dropTable(*typedCommand)
	case *CommandCreateTable:
		return db.createTable(*typedCommand)
	case *CommandInsert:
		return db.runInsert(*typedCommand)
	case *CommandSelect:
		return db.runSelect(*typedCommand)
	case *CommandUpdate:
		return db.runUpdate(*typedCommand)
	case *CommandDelete:
		return db.runDelete(*typedCommand)
	case *CommandRemoveDuplicates:
		return db.runRemoveDuplicates(*typedCommand)
	default:
		return nil, fmt.Errorf("unknown command type: %T", typedCommand)
	}
}

type CommandDropTable struct {
	Name string
}

// Drop table from db
func (d *Database) dropTable(command CommandDropTable) (*[]table.ColumnSet, error) {
	_, ok := d.tables[command.Name]
	if !ok {
		return nil, errs.NewErrTableNotFound(command.Name)
	}
	delete(d.tables, command.Name)
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully dropped table %s", command.Name)}}, nil
}

type CommandCreateTable struct {
	Name   string
	Schema schema.T
}

// Create new table in db
func (d *Database) createTable(command CommandCreateTable) (*[]table.ColumnSet, error) {
	if _, ok := d.tables[command.Name]; ok {
		return nil, errs.NewErrTableAlreadyExists(command.Name)
	}
	createdTable, err := table.NewTable(command.Schema)
	if err != nil {
		return nil, err
	}
	d.tables[command.Name] = createdTable
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully created table %s", command.Name)}}, nil
}

type CommandInsert struct {
	To   string
	Data *[]table.ColumnSet
}

// Insert rows to db table
func (d *Database) runInsert(command CommandInsert) (*[]table.ColumnSet, error) {
	to, err := d.table(command.To)
	if err != nil {
		return nil, err
	}
	inserted, err := to.InsertRows(command.Data)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {
			"message": fmt.Sprintf("successfully inserted %d %s to table %s",
				inserted,
				english.PluralWord(int(inserted), "row", ""),
				command.To)}},
		nil
}

type CommandSelect struct {
	From       string
	Fields     *[]string
	Conditions table.ColumnSet
}

// Select rows from db table
func (d *Database) runSelect(command CommandSelect) (*[]table.ColumnSet, error) {
	to, err := d.table(command.From)
	if err != nil {
		return nil, err
	}
	rows, err := to.SelectRows(command.Fields, command.Conditions)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type CommandUpdate struct {
	TableName  string
	Conditions table.ColumnSet
	Data       table.ColumnSet
}

// Update rows in db table
func (d *Database) runUpdate(command CommandUpdate) (*[]table.ColumnSet, error) {
	to, err := d.table(command.TableName)
	if err != nil {
		return nil, err
	}
	rowsCount, err := to.UpdateRows(command.Conditions, command.Data)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully updated %d %s in table %s",
		rowsCount,
		english.PluralWord(int(rowsCount), "row", ""),
		command.TableName)}}, nil
}

type CommandDelete struct {
	From       string
	Conditions table.ColumnSet
}

// Delete rows from db table
func (d *Database) runDelete(command CommandDelete) (*[]table.ColumnSet, error) {
	to, err := d.table(command.From)
	if err != nil {
		return nil, err
	}
	rowsCount, err := to.DeleteRows(command.Conditions)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully deleted %d %s from table %s",
		rowsCount,
		english.PluralWord(int(rowsCount), "row", ""),
		command.From)}}, nil
}

type CommandRemoveDuplicates struct {
	From string
}

// Delete duplicate rows from db table
func (d *Database) runRemoveDuplicates(command CommandRemoveDuplicates) (*[]table.ColumnSet, error) {
	to, err := d.table(command.From)
	if err != nil {
		return nil, err
	}
	rowsCount, err := to.DeleteDuplicates()
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully deleted %d %s from table %s",
		rowsCount,
		english.PluralWord(int(rowsCount), "row", ""),
		command.From)}}, nil
}

func (d *Database) table(name string) (*table.T, error) {
	table, ok := d.tables[name]
	if !ok {
		return nil, errs.NewErrTableNotFound(name)
	}
	return table, nil
}
