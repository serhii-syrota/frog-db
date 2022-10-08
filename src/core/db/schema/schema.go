package schema

import (
	"github.com/ssyrota/frog-db/src/core/dbtypes"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type T struct {
	Val map[string]dbtypes.Type
}

func New(data map[string]string) (*T, error) {
	schema := make(map[string]dbtypes.Type)
	for k, v := range data {
		if !dbtypes.IsAvailableName(v) {
			return nil, errs.NewErrInvalidTypeProvided(k, v)
		}
		schema[k] = dbtypes.Type(v)
	}
	return &T{schema}, nil
}
