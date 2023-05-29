package internal

import (
	"fmt"
	"strconv"
)

func cmdFilterLen(c *cmdCtx) error {
	fields := c.c.StringSlice("fields")

	if !c.c.IsSet("min") && !c.c.IsSet("max") {
		return fmt.Errorf("at least one of --min and --max must be set")
	}

	var minLen *int
	var maxLen *int

	if c.c.IsSet("min") {
		min := c.c.Int("min")
		minLen = &min
	}

	if c.c.IsSet("max") {
		max := c.c.Int("max")
		maxLen = &max
	}

	var filtered []datum

	for i, d := range c.data {
		dataLen := calculateStringLength(d, fields...)

		if minLen != nil && dataLen <= *minLen {
			c.logger.Debug().
				Int("line", i+1).
				Int("length", dataLen).
				Strs("fields", fields).
				Msg("filtered by min length")
			continue
		}
		if maxLen != nil && dataLen >= *maxLen {
			c.logger.Debug().
				Int("line", i+1).
				Int("length", dataLen).
				Strs("fields", fields).
				Msg("filtered by max length")
			continue
		}

		filtered = append(filtered, d)
	}

	c.logger = c.logger.With().Int("filtered_count", len(filtered)).Logger()
	c.logger = c.logger.With().Int("out_of_bounds", len(c.data)-len(filtered)).Logger()
	c.logger.Info().Msg("finished filtering by length")

	c.logger.Info().Msg("writing output")
	err := write(c.outPath, filtered)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func calculateStringLength(m map[string]interface{}, keys ...string) int {
	length := 0
	for _, key := range keys {
		if value, ok := m[key]; ok {
			length += len(valueToString(value))
		}
	}
	return length
}

// TODO: this will explode with deeply nested objects.
func valueToString(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case []interface{}:
		var str string
		for _, element := range v {
			str += valueToString(element)
		}
		return str
	case map[string]interface{}:
		var str string
		for key, val := range v {
			str += key + ": " + valueToString(val)
		}
		return str
	default:
		return fmt.Sprintf("%v", v)
	}
}
