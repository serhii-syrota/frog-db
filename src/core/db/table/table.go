package table

import (
	"sync"

	"github.com/ssyrota/frog-db/src/core/db/schema"
	"github.com/ssyrota/frog-db/src/core/dbtypes"
)

type T struct {
	sync.RWMutex
	schema *schema.T
}

func (t *T) Schema() map[string]dbtypes.Type {
	t.RLock()
	defer t.RUnlock()
	return t.schema.Val
}

func NewTable(sch *schema.T) *T {
	return &T{schema: sch}
}
