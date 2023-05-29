package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type filterFn func(string) bool

func cmdFilter(c *cmdCtx) error {
	if c.c.IsSet("wordlist") && c.c.IsSet("string") {
		return errors.New("cannot use both --wordlist and --string")
	}

	if !c.c.IsSet("wordlist") && !c.c.IsSet("string") {
		return errors.New("must specify either --wordlist or --string")
	}

	var filterList []string
	var match filterFn
	var err error

	if c.c.IsSet("wordlist") {
		filterList, err = loadWordlist(c.c.String("wordlist"))
		if err != nil {
			return fmt.Errorf("failed to load wordlist: %w", err)
		}
	} else {
		filterList = []string{c.c.String("string")}
	}

	if c.c.Bool("regex") {
		match, err = filterRegex(filterList)
		if err != nil {
			return fmt.Errorf("failed to compile regexp: %w", err)
		}
	} else {
		match = filterStr(filterList)
	}

	var filtered []datum
	extracted := extractFields(c)

	for i, e := range extracted {
		if !match(e) {
			filtered = append(filtered, c.data[i])
			continue
		}
		switch {
		// This is a bug in urfave
		// https://github.com/urfave/cli/issues/1737
		case c.c.Count("debug") == 2:
			c.logger.Debug().
				Int("line", i+1).
				Msg("filtered data")
		case c.c.Count("debug") >= 3:
			c.logger.Debug().
				Int("line", i+1).
				Interface("data", c.data[i]).
				Msg("duplicate found")
		}
		continue
	}

	c.logger = c.logger.With().Int("filtered_count", len(filtered)).Logger()
	c.logger = c.logger.With().Int("filter_hits", len(c.data)-len(filtered)).Logger()

	c.logger.Info().Msg("filtered data")

	c.logger.Info().Msg("writing data")

	return write(c.outPath, filtered)
}

func extractFields(c *cmdCtx) []string {
	fields := c.c.StringSlice("fields")
	var ret []string
	for _, d := range c.data {
		ret = append(ret, d.String(fields, false))
	}
	return ret
}

func filterStr(filters []string) filterFn {
	return func(s string) bool {
		for _, f := range filters {
			if strings.Contains(s, f) {
				return true
			}
		}
		return false
	}
}

func filterRegex(filters []string) (filterFn, error) {
	var rFilters []*regexp.Regexp
	for _, f := range filters {
		r, err := regexp.Compile(f)
		if err != nil {
			return nil, err
		}
		rFilters = append(rFilters, r)
	}

	return func(s string) bool {
		for _, f := range rFilters {
			if f.MatchString(s) {
				return true
			}
		}
		return false
	}, nil
}
