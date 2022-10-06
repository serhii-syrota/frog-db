package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerType(t *testing.T) {
	t.Run("recognizes int8", func(t *testing.T) {
		val := int8(121)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Value())
	})
	t.Run("recognizes int32", func(t *testing.T) {
		val := int32(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Value())
	})
	t.Run("recognizes int", func(t *testing.T) {
		val := 1
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Value())
	})
	t.Run("recognizes int64", func(t *testing.T) {
		val := int64(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, val, res.Value())
	})
	t.Run("recognizes uint64", func(t *testing.T) {
		val := uint64(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Value())
	})
	t.Run("recognizes float", func(t *testing.T) {
		val := 1.1
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.Value())
	})
}
