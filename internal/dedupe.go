package internal

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func cmdDedupe(c *cmdCtx) error {
	var deduped []datum
	var err error
	if c.c.Bool("rl") {
		deduped, err = dedupeRL(c, c.data)
		if err != nil {
			return fmt.Errorf("failed to dedupe data: %w", err)
		}
	} else {
		deduped = dedupe(c, c.data)
	}

	c.logger = c.logger.With().Int("duplicates_found", len(c.data)-len(deduped)).Logger()
	c.logger = c.logger.With().Int("out_record_count", len(deduped)).Logger()
	c.logger.Info().Msg("deduped data")

	c.logger.Info().Msg("writing deduped data")

	return write(c.outPath, deduped)
}

// dedupe could be made more efficient, but it's not worth it for now.
func dedupe(c *cmdCtx, data []datum) []datum {
	fields := c.c.StringSlice("fields")
	ignoreCase := c.c.Bool("ignore-case")

	// Create keys for each datum based on the fields to dedupe on.
	var keys []string
	for _, d := range data {
		if ignoreCase {
			keys = append(keys, strings.ToLower(d.String(fields, true)))
		} else {
			keys = append(keys, d.String(fields, true))
		}
	}

	// Create dedupe set
	seen := make(map[string]struct{})
	exists := struct{}{}

	ret := make([]datum, 0)

	for i, d := range keys {
		if _, found := seen[d]; found {
			switch {
			// This is a bug in urfave
			// https://github.com/urfave/cli/issues/1737
			case c.c.Count("debug") == 2:
				c.logger.Debug().
					Int("line", i+1).
					Msg("duplicate found")
			case c.c.Count("debug") >= 3:
				c.logger.Debug().
					Int("line", i+1).
					Interface("data", data[i].String(fields, true)).
					Msg("duplicate found")
			}
			continue
		}

		seen[d] = exists
		ret = append(ret, data[i])
	}

	return ret
}

// dedupeRL is not currently very efficient; this can be improved.
func dedupeRL(c *cmdCtx, data []datum) ([]datum, error) {
	pbar := progressbar.DefaultSilent(0)
	if c.c.Bool("progress") {
		pbar = progressbar.Default(int64(len(data)))
	}

	fields := c.c.StringSlice("fields")
	thresh := c.c.Float64("rl-threshold")

	var loweredFields [][]string
	var deduped []datum

	for _, d := range data {
		strD := d.String(fields, true)
		lstrD := strings.ToLower(strD)
		lfS := strings.Fields(lstrD)
		loweredFields = append(loweredFields, lfS)
	}

	var wg sync.WaitGroup
	dedupedChan := make(chan datum, len(data))

	numCPU := runtime.NumCPU()
	chunkSize := (len(loweredFields) + numCPU - 1) / numCPU
	for i := 0; i < len(loweredFields); i += chunkSize {
		end := i + chunkSize
		if end > len(loweredFields) {
			end = len(loweredFields)
		}
		wg.Add(1)
		go func(i, end int) {
			defer wg.Done()
			for ; i < end; i++ {
				if c.c.Bool("progress") {
					pbar.Add(1)
				}
				d1 := loweredFields[i]
				if len(d1) == 0 {
					c.logger.Debug().
						Int("line", i).
						Msg("empty field, not comparing")
					dedupedChan <- data[i]
					continue
				}

				isDuplicate := false
				for j := i + 1; j < len(loweredFields); j++ {
					d2 := loweredFields[j]

					rl := rougeL(d1, d2)
					if rl > thresh {
						c.logger.Debug().
							Str("d1", strings.Join(d1, " ")).
							Str("d2", strings.Join(d2, " ")).
							Float64("roguel", rl).
							Msg("duplicate found")
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					dedupedChan <- data[i]
				}
			}
		}(i, end)
	}
	go func() {
		wg.Wait()
		close(dedupedChan)
	}()

	for d := range dedupedChan {
		deduped = append(deduped, d)
	}

	return deduped, nil
}
