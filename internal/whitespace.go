package internal

import (
	"strings"
)

func cmdWhitespace(c *cmdCtx) error {
	fields := c.c.StringSlice("fields")

	c.logger.Info().Msg("trimming whitespace")
	var trimmed []datum
	var trimCnt int
	for i, datum := range c.data {
		for _, field := range fields {
			s, ok := datum[field].(string)
			if !ok {
				continue
			}
			sTrm := strings.TrimSpace(s)
			if s != sTrm {
				trimCnt++
				c.logger.Debug().
					Int("line", i+1).
					Str("field", field).
					Str("before", s).
					Str("after", sTrm).
					Msg("trimmed whitespace")
			}
			datum[field] = sTrm

		}
		trimmed = append(trimmed, datum)
	}
	c.logger = c.logger.With().Int("trim_count", trimCnt).Logger()
	c.logger.Info().Msg("trimmed whitespace")

	c.logger.Info().Msg("writing trimmed data")
	return write(c.outPath, trimmed)
}
