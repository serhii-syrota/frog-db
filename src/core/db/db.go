package db

import (
	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/db/table"
	"github.com/ssyrota/frog-db/src/core/dbtypes"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type Database struct {
	tables map[string]*table.T
}

func NewDb() *Database {
	return &Database{tables: make(map[string]*table.T)}
}

var _ Db = new(Database)

// CreateTable implementation.
func (d *Database) CreateTable(name string, schema *schema.T) error {
	if _, ok := d.tables[name]; ok {
		return errs.NewErrTableAlreadyExists(name)
	}
	d.tables[name] = table.NewTable(schema)
	return nil
}

// TableSchema implementation.
func (d *Database) TableSchema(name string) (map[string]dbtypes.Type, error) {
	table, ok := d.tables[name]
	if !ok {
		return nil, errs.NewErrTableNotFound(name)
	}
	return table.Schema(), nil
}

type Db interface {
	CreateTable(name string, schema *schema.T) error
	TableSchema(name string) (map[string]dbtypes.Type, error)
}
