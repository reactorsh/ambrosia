package internal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileAppender(t *testing.T) {
	t.Run("append data correctly", func(t *testing.T) {
		// Setup
		testFile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		appender, err := newFileAppender(testFile.Name())
		assert.NoError(t, err)
		testDatum := map[string]interface{}{"key": "value"}

		// Execute
		err = appender.append(testDatum)

		// Verify
		assert.NoError(t, err)
		content, err := os.ReadFile(testFile.Name())
		assert.NoError(t, err)
		var unmarshalledDatum map[string]interface{}
		err = json.Unmarshal(content, &unmarshalledDatum)
		assert.NoError(t, err)
		assert.Equal(t, testDatum, unmarshalledDatum)

		// Teardown
		err = os.Remove(testFile.Name())
		assert.NoError(t, err)
	})

	t.Run("close correctly", func(t *testing.T) {
		// Setup
		testFile, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		appender, err := newFileAppender(testFile.Name())
		assert.NoError(t, err)

		// Execute
		err = appender.close()

		// Verify
		assert.NoError(t, err)

		// Check if the file is really closed by trying to remove it
		err = os.Remove(testFile.Name())
		assert.NoError(t, err)
	})
}

func TestPrefixAppender(t *testing.T) {
	t.Run("append data correctly", func(t *testing.T) {
		// Setup
		pathTmpl := "test_%[1]c.json"
		appender := newPrefixAppender(pathTmpl)
		testDatum := map[string]interface{}{"key": "value"}

		// Execute
		err := appender.append("test response", testDatum)

		// Verify
		assert.NoError(t, err)
		content, err := os.ReadFile("test_t.json")
		assert.NoError(t, err)
		var unmarshalledDatum map[string]interface{}
		err = json.Unmarshal(content, &unmarshalledDatum)
		assert.NoError(t, err)
		assert.Equal(t, testDatum, unmarshalledDatum)

		// Teardown
		err = os.Remove("test_t.json")
		assert.NoError(t, err)
	})

	t.Run("return error if response is empty", func(t *testing.T) {
		// Setup
		pathTmpl := "test_%[1]c.json"
		appender := newPrefixAppender(pathTmpl)
		testDatum := map[string]interface{}{"key": "value"}

		// Execute
		err := appender.append("", testDatum)

		// Verify
		assert.Error(t, err)
		assert.Equal(t, errNilResponse, err)
	})

	t.Run("close correctly", func(t *testing.T) {
		// Setup
		pathTmpl := "test_%[1]c.json"
		appender := newPrefixAppender(pathTmpl)
		testDatum := map[string]interface{}{"key": "value"}
		err := appender.append("test response", testDatum)
		assert.NoError(t, err)

		// Execute
		err = appender.close()

		// Verify
		assert.NoError(t, err)

		// Check if the file is really closed by trying to remove it
		err = os.Remove("test_t.json")
		assert.NoError(t, err)
	})
}
