package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenOutPath tests the genOutPath function.
func TestGenOutPath(t *testing.T) {
	t.Run("test with one argument", func(t *testing.T) {
		outPath := genOutPath("cmd", []string{"/path/to/file.txt"})
		assert.Equal(t, "/path/to/file_cmd.txt", outPath)
	})

	t.Run("test with more than one argument", func(t *testing.T) {
		outPath := genOutPath("cmd", []string{"/path/to/file.txt", "/another/path/to/file.txt"})
		assert.Equal(t, "/another/path/to/file.txt", outPath)
	})
}

func TestPrefixPathTmpl(t *testing.T) {
	assert := assert.New(t)

	// Test Case 1: Regular file path
	filePath := "/home/user/test.txt"
	expected := "/home/user/test_psort_%c.txt"
	result := prefixPathTmpl(filePath)
	assert.Equal(expected, result, "They should be equal")

	// Test Case 2: File path without extension
	filePath = "/home/user/test"
	expected = "/home/user/test_psort_%c"
	result = prefixPathTmpl(filePath)
	assert.Equal(expected, result, "They should be equal")

	// Test Case 3: File in root directory
	filePath = "/test.txt"
	expected = "/test_psort_%c.txt"
	result = prefixPathTmpl(filePath)
	assert.Equal(expected, result, "They should be equal")

	// Test Case 4: File path with spaces and special characters
	filePath = "/home/user/test file@.txt"
	expected = "/home/user/test file@_psort_%c.txt"
	result = prefixPathTmpl(filePath)
	assert.Equal(expected, result, "They should be equal")

	// Test Case 5: Empty file path
	filePath = ""
	expected = "_psort_%c."
	result = prefixPathTmpl(filePath)
	assert.Equal(expected, result, "They should be equal")
}
