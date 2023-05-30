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
