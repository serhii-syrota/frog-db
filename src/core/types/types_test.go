package types

import (
	"testing"

	"github.com/ssyrota/frog-db/src/core/errs"
	"github.com/stretchr/testify/assert"
)

func TestIntegerType(t *testing.T) {
	t.Run("recognizes int8", func(t *testing.T) {
		val := int8(121)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Val)
	})
	t.Run("recognizes int32", func(t *testing.T) {
		val := int32(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Val)
	})
	t.Run("recognizes int", func(t *testing.T) {
		val := 1
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Val)
	})
	t.Run("recognizes int64", func(t *testing.T) {
		val := int64(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, val, res.Val)
	})
	t.Run("recognizes uint64", func(t *testing.T) {
		val := uint64(123)
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(val), res.Val)
	})
	t.Run("recognizes float", func(t *testing.T) {
		val := 1.1
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.Val)
	})
	t.Run("recognizes valid string", func(t *testing.T) {
		val := "10"
		res, err := NewInteger(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, int64(10), res.Val)
	})
	t.Run("fails on invalid string", func(t *testing.T) {
		val := "10.0"
		res, err := NewInteger(val, "testing_column")
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestRealType(t *testing.T) {
	t.Run("recognizes float32", func(t *testing.T) {
		val := float32(121.1)
		res, err := NewReal(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, float64(val), res.Val)
	})
	t.Run("recognizes float64", func(t *testing.T) {
		val := float64(121.1)
		res, err := NewReal(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, float64(val), res.Val)
	})
	t.Run("recognizes int64", func(t *testing.T) {
		val := int64(12)
		res, err := NewReal(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, float64(val), res.Val)
	})
	t.Run("recognizes uint64", func(t *testing.T) {
		val := uint64(12)
		res, err := NewReal(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, float64(val), res.Val)
	})
	t.Run("recognizes string", func(t *testing.T) {
		val := "10.00"
		res, err := NewReal(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, float64(10.00), res.Val)
	})
	t.Run("fails on invalid string", func(t *testing.T) {
		val := "01,00"
		res, err := NewReal(val, "testing_column")
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestCharType(t *testing.T) {
	t.Run("recognizes string", func(t *testing.T) {
		val := "a"
		res, err := NewChar(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, 'a', res.Val)
	})
	t.Run("fails on invalid string", func(t *testing.T) {
		val := "ab"
		res, err := NewChar(val, "testing_column")
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("recognizes char", func(t *testing.T) {
		val := 'a'
		res, err := NewChar(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, 'a', res.Val)
	})
}

func TestStringType(t *testing.T) {
	t.Run("recognizes string", func(t *testing.T) {
		val := "asd"
		res, err := NewString(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, val, res.Val)
	})
	t.Run("recognizes char", func(t *testing.T) {
		val := 'a'
		res, err := NewString(val, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, "a", res.Val)
	})
}

func TestRealInvType(t *testing.T) {
	t.Run("accepts first value is lower", func(t *testing.T) {
		a := "1.1"
		b := "2.2"
		res, err := NewRealInv(a, b, "testing_column")
		assert.Nil(t, err)
		assert.Equal(t, 1.1, res.A)
		assert.Equal(t, 2.2, res.B)
	})
	t.Run("fails on second value is lower", func(t *testing.T) {
		a := "2.2"
		b := "1.1"
		_, err := NewRealInv(a, b, "testing_column")
		assert.NotNil(t, err)
		assert.IsType(t, &errs.ErrInvalidRange{}, err)
	})
}
