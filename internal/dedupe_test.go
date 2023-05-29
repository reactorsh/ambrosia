package internal

import (
	"flag"
	"os"
	"strconv"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func BenchmarkDedupeRL(b *testing.B) {
	// Mock the cli.Context
	set := flag.NewFlagSet("test", 0)
	set.String("field", "name", "doc")
	set.Float64("rl-threshold", 0.5, "doc")
	ctx := cli.NewContext(nil, set, nil)

	data := make([]datum, 100)
	for i := 0; i < 100; i++ {
		data[i] = datum{"name": "Test" + strconv.Itoa(i)}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dedupeRL(&cmdCtx{c: ctx}, data)
	}
}

func TestDedupeRL(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.Var(cli.NewStringSlice("testField"), "fields", "doc")
	set.Float64("rl-threshold", 0.8, "doc")
	ctx := cli.NewContext(app, set, nil)

	t.Run("field is an empty string", func(t *testing.T) {
		data := []datum{
			{"testField": ""},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("no duplicate fields", func(t *testing.T) {
		data := []datum{
			{"testField": "test1"},
			{"testField": "test2"},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("duplicate fields", func(t *testing.T) {
		data := []datum{
			{"testField": "test"},
			{"testField": "test"},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Equal(t, []datum{{"testField": "test"}}, result)
	})

	t.Run("duplicate fields with different case", func(t *testing.T) {
		data := []datum{
			{"testField": "test"},
			{"testField": "TEST"},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("near duplicate sentences that should be treated as duplicates", func(t *testing.T) {
		data := []datum{
			{"testField": "This is a test."},
			{"testField": "This is also a test."},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("field not a string", func(t *testing.T) {
		data := []datum{
			{"testField": 1234},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("field not found in data", func(t *testing.T) {
		data := []datum{
			{"otherField": "test"},
		}
		result, err := dedupeRL(&cmdCtx{c: ctx}, data)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})
}

func TestDedupe(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	t.Run("dedupe should remove duplicate data", func(t *testing.T) {
		data := []datum{
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
			{"name": "Alice", "age": 25},
		}

		app := cli.NewApp()
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("name", "age"), "fields", "doc")
		ctx := &cmdCtx{
			c:      cli.NewContext(app, set, nil),
			logger: logger,
			data:   data,
		}

		expected := []datum{
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
		}

		assert.Equal(t, expected, dedupe(ctx, data))
	})

	t.Run("dedupe should apply ignore-case flag", func(t *testing.T) {
		data := []datum{
			{"name": "Alice", "age": 25},
			{"name": "alice", "age": 25},
		}

		app := cli.NewApp()
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("name", "age"), "fields", "doc")
		set.Bool("ignore-case", true, "doc")
		ctx := &cmdCtx{
			c:      cli.NewContext(app, set, nil),
			logger: logger,
			data:   data,
		}

		expected := []datum{
			{"name": "Alice", "age": 25},
		}

		assert.Equal(t, expected, dedupe(ctx, data))
	})

	t.Run("dedupe should handle empty slice", func(t *testing.T) {
		data := []datum{}

		app := cli.NewApp()
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("name", "age"), "fields", "doc")
		ctx := &cmdCtx{
			c:      cli.NewContext(app, set, nil),
			logger: logger,
			data:   data,
		}

		assert.Equal(t, data, dedupe(ctx, data))
	})

	t.Run("dedupe should handle no duplicate entries", func(t *testing.T) {
		data := []datum{
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
		}

		app := cli.NewApp()
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("name", "age"), "fields", "doc")
		ctx := &cmdCtx{
			c:      cli.NewContext(app, set, nil),
			logger: logger,
			data:   data,
		}

		assert.Equal(t, data, dedupe(ctx, data))
	})

	t.Run("dedupe should check only the requested fields", func(t *testing.T) {
		data := []datum{
			{"name": "Alice", "age": 25, "city": "New York"},
			{"name": "Bob", "age": 30, "city": "New York"},
			{"name": "Alice", "age": 25, "city": "San Francisco"},
		}

		app := cli.NewApp()
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("name", "age"), "fields", "doc")
		ctx := &cmdCtx{
			c:      cli.NewContext(app, set, nil),
			logger: logger,
			data:   data,
		}

		// Even though the 'city' is different for the first and third datum,
		// they should be considered duplicates because we're only deduping on 'name' and 'age'.
		expected := []datum{
			{"name": "Alice", "age": 25, "city": "New York"},
			{"name": "Bob", "age": 30, "city": "New York"},
		}

		assert.Equal(t, expected, dedupe(ctx, data))
	})
}
