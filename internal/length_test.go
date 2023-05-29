package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueToString(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		assert.Equal(t, "123.456", valueToString(123.456))
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "hello world", valueToString("hello world"))
	})

	t.Run("bool", func(t *testing.T) {
		assert.Equal(t, "true", valueToString(true))
		assert.Equal(t, "false", valueToString(false))
	})

	t.Run("slice", func(t *testing.T) {
		assert.Equal(t, "123.456hello worldtrue", valueToString([]interface{}{123.456, "hello world", true}))
	})

	t.Run("map", func(t *testing.T) {
		result := valueToString(map[string]interface{}{"key1": 123.456, "key2": "hello world"})
		assert.Contains(t, result, "key1: 123.456")
		assert.Contains(t, result, "key2: hello world")
		assert.Equal(t, len("key1: 123.456key2: hello world"), len(result))
	})

	t.Run("other", func(t *testing.T) {
		assert.Equal(t, "123", valueToString(123))
	})
}

func TestCalculateStringLength(t *testing.T) {
	m := map[string]interface{}{
		"key1": 123.456,
		"key2": "hello world",
		"key3": true,
		"key4": []interface{}{123.456, "hello world", true},
		"key5": map[string]interface{}{"key1": 123.456, "key2": "hello world"},
	}

	t.Run("single key", func(t *testing.T) {
		assert.Equal(t, 7, calculateStringLength(m, "key1")) // accounting for "123.456"
	})

	t.Run("multiple keys", func(t *testing.T) {
		// accounting for "123.456" (7 chars), "hello world" (11 chars), "true" (4 chars)
		// array with "123.456", "hello world", "true" (22 chars)
		// map with "key1: 123.456" (14 chars) and "key2: hello world" (16 chars)
		assert.Equal(t, 7+11+4+22+14+16, calculateStringLength(m, "key1", "key2", "key3", "key4", "key5"))
	})

	t.Run("non-existing keys", func(t *testing.T) {
		assert.Equal(t, 0, calculateStringLength(m, "key6"))
	})

	t.Run("empty key", func(t *testing.T) {
		assert.Equal(t, 0, calculateStringLength(m, ""))
	})

	t.Run("empty map", func(t *testing.T) {
		assert.Equal(t, 0, calculateStringLength(map[string]interface{}{}, "key1"))
	})

	t.Run("no keys", func(t *testing.T) {
		assert.Equal(t, 0, calculateStringLength(m))
	})
}
