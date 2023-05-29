package internal

import (
	"errors"
	"fmt"

	"github.com/pkg/profile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type cmdCtx struct {
	c       *cli.Context
	inPath  string
	outPath string
	logger  zerolog.Logger
	data    []datum
}

func CmdInit(c *cli.Context) error {
	if c.Bool("cpuprofile") {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	if c.Bool("memprofile") {
		defer profile.Start(profile.MemProfile).Stop()
	}

	logger := log.With().Str("command", c.Command.Name).Logger()

	if len(c.Args().Slice()) < 1 {
		return errors.New("missing required argument, review the help documentation with -h")
	}

	inPath := c.Args().First()
	logger = logger.With().
		Str("infile", inPath).
		Logger()

	outPath := genOutPath(c.Command.Name, c.Args().Slice())
	if c.Command.Name != "psort" {
		logger = logger.With().
			Str("outfile", outPath).
			Logger()
	}

	data, err := load(inPath)
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	logger = logger.With().Int("in_record_count", len(data)).Logger()
	logger.Info().Msg("loaded data")

	ctx := &cmdCtx{
		c:       c,
		inPath:  inPath,
		outPath: outPath,
		logger:  logger,
		data:    data,
	}

	switch c.Command.Name {
	case "dedupe":
		err = cmdDedupe(ctx)
	case "length":
		err = cmdFilterLen(ctx)
	case "filter":
		err = cmdFilter(ctx)
	case "psort":
		err = cmdPSort(ctx)
	case "whitespace":
		err = cmdWhitespace(ctx)
	}

	if err != nil {
		return err
	}

	ctx.logger.Info().Msg("done")

	return nil
}
