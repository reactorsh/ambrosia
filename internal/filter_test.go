package internal

import (
	"flag"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestExtractField(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("field1", "field3"), "fields", "doc")
		c := &cmdCtx{
			c: cli.NewContext(nil, set, nil),
			data: []datum{
				{"field1": "value1", "field2": "value2"},
				{"field1": "value3", "field2": "value4", "field3": "value5"},
			},
			logger: zerolog.Nop(),
		}
		res := extractFields(c)
		assert.Equal(t, []string{"value1\n", "value3\nvalue5\n"}, res)
	})

	t.Run("field not exist", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("field1", "field3"), "fields", "doc")
		c := &cmdCtx{
			c: cli.NewContext(nil, set, nil),
			data: []datum{
				{"field1": "value1", "field2": "value2"},
				{"field1": "value3", "field2": "value4"},
			},
			logger: zerolog.Nop(),
		}
		res := extractFields(c)
		assert.Equal(t, []string{"value1\n", "value3\n"}, res)
	})

	t.Run("empty data", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("field1"), "fields", "doc")
		c := &cmdCtx{
			c:      cli.NewContext(nil, set, nil),
			data:   []datum{},
			logger: zerolog.Nop(),
		}
		res := extractFields(c)
		assert.Empty(t, res)
	})
}

func TestFilterRegex(t *testing.T) {
	t.Run("valid regex", func(t *testing.T) {
		filters := []string{"abc", "123"}
		filterFn, err := filterRegex(filters)
		assert.NoError(t, err)

		assert.True(t, filterFn("abcdef"))
		assert.True(t, filterFn("987123"))
		assert.False(t, filterFn("xyz987"))
	})

	t.Run("invalid regex", func(t *testing.T) {
		filters := []string{"abc", "["} // "[" is an invalid regex
		_, err := filterRegex(filters)
		assert.Error(t, err)
	})

	t.Run("no match", func(t *testing.T) {
		filters := []string{"abc", "123"}
		filterFn, err := filterRegex(filters)
		assert.NoError(t, err)

		assert.False(t, filterFn("xyz987"))
	})

	t.Run("empty filters", func(t *testing.T) {
		var filters []string
		filterFn, err := filterRegex(filters)
		assert.NoError(t, err)

		assert.False(t, filterFn("xyz987"))
	})
}

func TestFilterStr(t *testing.T) {
	t.Run("string contains substring", func(t *testing.T) {
		filters := []string{"abc", "123"}
		filterFn := filterStr(filters)

		assert.True(t, filterFn("abcdef"))
		assert.True(t, filterFn("987123"))
		assert.False(t, filterFn("xyz987"))
	})

	t.Run("no match", func(t *testing.T) {
		filters := []string{"abc", "123"}
		filterFn := filterStr(filters)

		assert.False(t, filterFn("xyz987"))
	})

	t.Run("empty filters", func(t *testing.T) {
		var filters []string
		filterFn := filterStr(filters)

		assert.False(t, filterFn("xyz987"))
	})
}
