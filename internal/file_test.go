package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLoad tests the load function.
func TestLoad(t *testing.T) {
	t.Run("test non-existent file", func(t *testing.T) {
		_, err := load("nonexistentfile")
		assert.Error(t, err)
	})

	t.Run("test invalid JSON", func(t *testing.T) {
		file, _ := os.CreateTemp("", "prefix")
		defer os.Remove(file.Name())

		file.WriteString("invalid json")

		_, err := load(file.Name())
		assert.Error(t, err)
	})

	t.Run("test valid JSON", func(t *testing.T) {
		file, _ := os.CreateTemp("", "prefix")
		defer os.Remove(file.Name())

		file.WriteString(`{"key":"value"}`)

		data, err := load(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(data))
		assert.Equal(t, "value", data[0]["key"])
	})
}

// TestWrite tests the write function.
func TestWrite(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "prefix")
	defer os.RemoveAll(tempDir)

	t.Run("test file already exists", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "existing_file.txt")
		file, _ := os.Create(filePath)
		file.Close()

		err := write(filePath, []datum{{"key": "value"}})
		assert.Error(t, err)
	})

	t.Run("test JSON marshal error", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "output_file.txt")

		err := write(filePath, []datum{{"key": make(chan int)}})
		assert.Error(t, err)
	})

	t.Run("test successful write", func(t *testing.T) {
		fileName := "output_file_" + time.Now().Format("20060102150405") + ".txt"
		filePath := filepath.Join(tempDir, fileName)

		err := write(filePath, []datum{{"key": "value"}})
		assert.NoError(t, err)

		content, _ := os.ReadFile(filePath)
		assert.Equal(t, "{\"key\":\"value\"}\n", string(content))
	})
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

func TestLoadResumable(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test")
	assert.NoError(err)

	// Ensure the temporary directory is cleaned up
	defer os.RemoveAll(tmpDir)

	// Create test data
	datum1 := datum{"key1": "value1", "key2": "value2"}
	datum2 := datum{"key3": "value3", "key4": "value4"}
	datum3 := datum{"key5": "value5", "key6": "value6"}

	// Create resumable files
	for i, d := range []datum{datum1, datum2} {
		f, err := os.Create(filepath.Join(tmpDir, "testFile_"+strconv.Itoa(i)+".jsonl"))
		assert.NoError(err)

		data, err := json.Marshal(d)
		assert.NoError(err)

		_, err = f.WriteString(string(data) + "\n")
		assert.NoError(err)

		f.Close()
	}

	// Create non-resumable file
	f, err := os.Create(filepath.Join(tmpDir, "unrelatedFile.jsonl"))
	assert.NoError(err)

	data, err := json.Marshal(datum3)
	assert.NoError(err)

	_, err = f.WriteString(string(data) + "\n")
	assert.NoError(err)

	f.Close()

	// Call loadResumable
	res, err := loadResumable("cmd", filepath.Join(tmpDir, "testFile.jsonl"))
	assert.NoError(err)
	assert.Equal([]datum{datum1, datum2}, res)

	// Remove the resumable files
	for i := 0; i < 2; i++ {
		os.Remove(filepath.Join(tmpDir, "testFile_"+strconv.Itoa(i)+".jsonl"))
	}

	// Case when no resumable files are found
	res, err = loadResumable("cmd", filepath.Join(tmpDir, "testFile.jsonl"))
	assert.NoError(err)
	if res == nil {
		assert.Nil(res)
	} else {
		assert.Equal([]datum{}, res)
	}
}

func TestLoadWordlist(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write some data to it
	data := []string{"word1", "word2", "word3"}
	for _, word := range data {
		tmpfile.WriteString(word + "\n")
	}

	// Ensure the file data is written to disk
	tmpfile.Sync()
	tmpfile.Close()

	// Test loadWordlist
	words, err := loadWordlist(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, data, words, "The loaded wordlist did not match the expected data")
}

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
