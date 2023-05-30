package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatumString(t *testing.T) {
	t.Run("Test with different types", func(t *testing.T) {
		testCases := []struct {
			name     string
			input    datum
			keys     []string
			expected string
		}{
			{
				name:     "Test with string type",
				input:    datum{"foo": "bar", "baz": "quux"},
				keys:     []string{"foo"},
				expected: "foo: bar\n",
			},
			{
				name:     "Test with integer type",
				input:    datum{"one": 1, "two": 2},
				keys:     []string{"one", "two"},
				expected: "one: 1\ntwo: 2\n",
			},
			{
				name:     "Test with float type",
				input:    datum{"onePointOne": 1.1, "twoPointTwo": 2.2},
				keys:     []string{"onePointOne"},
				expected: "onePointOne: 1.1\n",
			},
			{
				name:     "Test with bool type",
				input:    datum{"trueVal": true, "falseVal": false},
				keys:     []string{"falseVal"},
				expected: "falseVal: false\n",
			},
			{
				name:     "Test with nil value",
				input:    datum{"nilVal": nil},
				keys:     []string{"nilVal"},
				expected: "nilVal: <nil>\n",
			},
			{
				name:     "Test with selected keys",
				input:    datum{"foo": "bar", "baz": "quux", "blern": 3},
				keys:     []string{"foo", "blern"},
				expected: "foo: bar\nblern: 3\n",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, tc.input.String(tc.keys, true))
			})
		}
	})
}

func BenchmarkIsEqual(b *testing.B) {
	testData := []struct {
		name string
		a, b datum
	}{
		{
			name: "equal data",
			a: datum{
				"field1": "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
			b: datum{
				"field1": "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
		},
		{
			name: "equal data ambrosia",
			a: datum{
				"ambrosia": "foo bar baz",
				"field1":   "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2":   "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3":   "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
			b: datum{
				"field1": "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
		},
		{
			name: "unequal data",
			a: datum{
				"field1": "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
			b: datum{
				"field1": "It was the best of times, it was the worst of times.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
		},
		{
			name: "unequal data ambrosia",
			a: datum{
				"ambrosia": "foo bar baz",
				"field1":   "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
				"field2":   "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3":   "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
			b: datum{
				"field1": "It was the best of times, it was the worst of times.",
				"field2": "I am no bird; and no net ensnares me: I am a free human being with an independent will.",
				"field3": "Happy families are all alike; every unhappy family is unhappy in its own way.",
			},
		},
	}

	for _, tt := range testData {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				isEqual(tt.a, tt.b)
			}
		})
	}
}

func TestDatumFunctions(t *testing.T) {
	t.Run("isEqual correctly compares two datums", func(t *testing.T) {
		// Setup
		datum1 := datum{"key1": "value1", "key2": 2}
		datum2 := datum{"key1": "value1", "key2": 2}
		datum3 := datum{"key1": "value1", "key2": 3}
		datum4 := datum{"ambrosia": "foo", "key1": "value1", "key2": 2}

		// Verify
		assert.True(t, isEqual(datum1, datum2))
		assert.True(t, isEqual(datum1, datum4))

		assert.False(t, isEqual(datum1, datum3))
		assert.False(t, isEqual(datum3, datum4))
	})

	t.Run("contains correctly checks if a datum is in a slice", func(t *testing.T) {
		// Setup
		datum1 := datum{"key1": "value1", "key2": 2}
		datum2 := datum{"key1": "value1", "key2": 3}
		datumSlice := []datum{
			{"key1": "value1", "key2": 2},
			{"key1": "value2", "key2": 4},
		}

		// Verify
		assert.True(t, contains(datumSlice, datum1))
		assert.False(t, contains(datumSlice, datum2))
	})

	t.Run("datumSub correctly subtracts one slice from another", func(t *testing.T) {
		// Setup
		datumSlice1 := []datum{
			{"key1": "value1", "key2": 2},
			{"key1": "value2", "key2": 4},
		}
		datumSlice2 := []datum{
			{"key1": "value1", "key2": 2},
		}
		datumSlice3 := []datum{
			{"key1": "value2", "key2": 4},
		}

		// Execute
		result := datumSub(datumSlice1, datumSlice2)

		// Verify
		assert.Equal(t, datumSlice3, result)
	})
}
