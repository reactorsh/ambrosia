package internal

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCmdWhitespace(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	logger := zerolog.New(os.Stdout)

	t.Run("Test with nonexistent field in datum", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		inputPath := tempDir + "/input"
		outputPath := tempDir + "/output"

		ctx := &cmdCtx{
			inPath:  inputPath,
			outPath: outputPath,
			logger:  logger,
			data:    []datum{{"test": " test "}},
		}

		app := &cli.App{}
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("nonexistent"), "fields", "doc")
		ctx.c = cli.NewContext(app, set, nil)

		err := cmdWhitespace(ctx)
		assert.NoError(t, err, "Unexpected error in cmdWhitespace")

		file, _ := os.Open(outputPath)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var data []datum
		for scanner.Scan() {
			var d datum
			err := json.Unmarshal([]byte(scanner.Text()), &d)
			require.NoError(t, err, "Unexpected error in json.Unmarshal")
			data = append(data, d)
		}

		assert.Equal(t, ctx.data, data, "Unexpected data in output file")
	})

	t.Run("Test with non-string entry in datum", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		inputPath := tempDir + "/input"
		outputPath := tempDir + "/output"

		ctx := &cmdCtx{
			inPath:  inputPath,
			outPath: outputPath,
			logger:  logger,
			data:    []datum{{"test": 123.01}},
		}

		app := &cli.App{}
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("test"), "fields", "doc")
		ctx.c = cli.NewContext(app, set, nil)

		err := cmdWhitespace(ctx)
		assert.NoError(t, err, "Unexpected error in cmdWhitespace")

		file, _ := os.Open(outputPath)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var data []datum
		for scanner.Scan() {
			var d datum
			err := json.Unmarshal([]byte(scanner.Text()), &d)
			require.NoError(t, err, "Unexpected error in json.Unmarshal")
			data = append(data, d)
		}

		assert.Equal(t, ctx.data, data, "Unexpected data in output file")
	})

	t.Run("Test with valid parameters", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		inputPath := tempDir + "/input"
		outputPath := tempDir + "/output"

		ctx := &cmdCtx{
			inPath:  inputPath,
			outPath: outputPath,
			logger:  logger,
			data:    []datum{{"test": " test "}},
		}

		app := &cli.App{}
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("test"), "fields", "doc")
		ctx.c = cli.NewContext(app, set, nil)

		err := cmdWhitespace(ctx)
		assert.NoError(t, err, "Unexpected error in cmdWhitespace")

		file, _ := os.Open(outputPath)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var data []datum
		for scanner.Scan() {
			var d datum
			err := json.Unmarshal([]byte(scanner.Text()), &d)
			require.NoError(t, err, "Unexpected error in json.Unmarshal")
			data = append(data, d)
		}

		assert.Equal(t, ctx.data, data, "Unexpected data in output file")
	})

	t.Run("Test with whitespace trimming", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "test")
		defer os.RemoveAll(tempDir)

		inputPath := tempDir + "/input"
		outputPath := tempDir + "/output"

		ctx := &cmdCtx{
			inPath:  inputPath,
			outPath: outputPath,
			logger:  logger,
			data:    []datum{{"test": " test "}},
		}

		app := &cli.App{}
		set := flag.NewFlagSet("test", 0)
		set.Var(cli.NewStringSlice("test"), "fields", "doc")
		ctx.c = cli.NewContext(app, set, nil)

		err := cmdWhitespace(ctx)
		assert.NoError(t, err, "Unexpected error in cmdWhitespace")

		file, _ := os.Open(outputPath)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var data []datum
		for scanner.Scan() {
			var d datum
			err := json.Unmarshal([]byte(scanner.Text()), &d)
			require.NoError(t, err, "Unexpected error in json.Unmarshal")
			data = append(data, d)
		}

		assert.Equal(t, "test", data[0]["test"], "Whitespace was not trimmed correctly")
	})
}
