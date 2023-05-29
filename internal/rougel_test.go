package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongestCommonSubsequence(t *testing.T) {
	testCases := []struct {
		s1          string
		s2          string
		expectedLCS int
		description string
	}{
		{
			s1:          "This is an example",
			s2:          "This is another example",
			expectedLCS: 3,
			description: "Common words between two sentences",
		},
		{
			s1:          "Hello world",
			s2:          "Goodbye world",
			expectedLCS: 1,
			description: "One common word",
		},
		{
			s1:          "Different words here",
			s2:          "No match at all",
			expectedLCS: 0,
			description: "No common words",
		},
		{
			s1:          "",
			s2:          "Some words here",
			expectedLCS: 0,
			description: "Empty string",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			result := longestCommonSubsequence(strings.Fields(testCase.s1), strings.Fields(testCase.s2))
			assert.Equal(t, testCase.expectedLCS, result)
		})
	}
}

func TestRougeL(t *testing.T) {
	testCases := []struct {
		s1             string
		s2             string
		expectedRougeL float64
		description    string
	}{
		{
			s1:             "This is an example",
			s2:             "This is another example",
			expectedRougeL: 0.75,
			description:    "Common words between two sentences",
		},
		{
			s1:             "Hello world",
			s2:             "Goodbye world",
			expectedRougeL: 0.5,
			description:    "One common word",
		},
		{
			s1:             "Different words here",
			s2:             "No match at all",
			expectedRougeL: 0.0,
			description:    "No common words",
		},
		{
			s1:             "",
			s2:             "Some words here",
			expectedRougeL: 0.0,
			description:    "Empty string",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			result := rougeL(strings.Fields(testCase.s1), strings.Fields(testCase.s2))
			assert.InDelta(t, testCase.expectedRougeL, result, 1e-9)
		})
	}
}
