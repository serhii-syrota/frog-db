package schema

import (
	dbtypes "github.com/ssyrota/frog-db/src/core/db/dbtypes"
	errs "github.com/ssyrota/frog-db/src/core/err"
)

type T = map[string]dbtypes.Type

func New(data map[string]string) (T, error) {
	schema := make(T)
	for k, v := range data {
		if !dbtypes.IsAvailableName(v) {
			return nil, errs.NewErrInvalidTypeProvided(k, v)
		}
		schema[k] = dbtypes.Type(v)
	}
	return schema, nil
}
