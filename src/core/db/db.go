package db

import (
	"errors"
	"fmt"

	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type Db interface {
	Execute(command any) (*[]table.ColumnSet, error)
	StoreDump() error
	// IntrospectSchema()
	// ExportData()
	// ImportData()
}
type Database struct {
	tables map[string]*table.T
	path   string
}

func New(path string) (*Database, error) {
	return &Database{tables: make(map[string]*table.T), path: path}, nil
}

var _ Db = new(Database)

// StoreDump implementation.
func (db *Database) StoreDump() error {
	return nil
}

// Execute implementation.
func (db *Database) Execute(command any) (*[]table.ColumnSet, error) {
	switch typedCommand := command.(type) {
	case CommandDropTable:
		return db.dropTable(typedCommand)
	case CommandCreateTable:
		return db.createTable(typedCommand)
	case CommandInsert:
		return db.runInsert(typedCommand)
	case CommandSelect:
		return db.runSelect(typedCommand)
	case CommandUpdate:
		return db.runUpdate(typedCommand)
	case CommandDelete:
		return db.runDelete(typedCommand)
	default:
		return nil, errors.New("unknown command")
	}
}

type CommandDropTable struct {
	Name string
}

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
	Schema *schema.T
}

func (d *Database) createTable(command CommandCreateTable) (*[]table.ColumnSet, error) {
	if _, ok := d.tables[command.Name]; ok {
		return nil, errs.NewErrTableAlreadyExists(command.Name)
	}
	d.tables[command.Name] = table.NewTable(command.Schema)
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully created table %s", command.Name)}}, nil
}

type CommandInsert struct {
	To   string
	Data *[]table.ColumnSet
}

func (d *Database) runInsert(command CommandInsert) (*[]table.ColumnSet, error) {
	to, err := d.table(command.To)
	if err != nil {
		return nil, err
	}
	inserted, err := to.InsertRows(command.Data)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully inserted %d rows to table %s", inserted, command.To)}}, nil
}

type CommandSelect struct {
	From       string
	Fields     *[]string
	Conditions table.ColumnSet
}

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

func (d *Database) runUpdate(command CommandUpdate) (*[]table.ColumnSet, error) {
	to, err := d.table(command.TableName)
	if err != nil {
		return nil, err
	}
	rowsCount, err := to.UpdateRows(command.Conditions, command.Data)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully updated %d rows in table %s", rowsCount, command.TableName)}}, nil
}

type CommandDelete struct {
	From       string
	Conditions table.ColumnSet
}

func (d *Database) runDelete(command CommandDelete) (*[]table.ColumnSet, error) {
	to, err := d.table(command.From)
	if err != nil {
		return nil, err
	}
	rowsCount, err := to.DeleteRows(command.Conditions)
	if err != nil {
		return nil, err
	}
	return &[]table.ColumnSet{0: {"message": fmt.Sprintf("successfully deleted %d rows from table %s", rowsCount, command.From)}}, nil
}

func (d *Database) table(name string) (*table.T, error) {
	table, ok := d.tables[name]
	if !ok {
		return nil, errs.NewErrTableNotFound(name)
	}
	return table, nil
}
