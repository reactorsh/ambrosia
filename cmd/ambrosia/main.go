package main

import (
	"os"
	"runtime"
	"time"

	"github.com/reactorsh/ambrosia/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// Commit is the git commit hash of the binary, set during build.
	version string
)

func main() {
	// Windows doesn't support colored output.
	if runtime.GOOS == "windows" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	app := cli.App{
		Name:    "ambrosia",
		Usage:   "prepare your text datasets",
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "oz",
				Email: "oz@reactor.sh",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				EnvVars: []string{"AMBROSIA_DEBUG", "DEBUG"},
				Usage:   "set log level to debug",
				Action: func(c *cli.Context, d bool) error {
					if d {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				EnvVars: []string{"AMBROSIA_JSON", "JSON"},
				Usage:   "use json for logging",
				Action: func(c *cli.Context, j bool) error {
					if j {
						log.Logger = log.Output(os.Stderr)
					}
					return nil
				},
				Value: false,
			},
			&cli.BoolFlag{
				Name:   "cpuprofile",
				Usage:  "enable cpu profiling",
				Hidden: true,
			},
			&cli.BoolFlag{
				Name:   "memprofile",
				Usage:  "enable memory profiling",
				Hidden: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "whitespace",
				ArgsUsage: "INFILE.jsonl [OUTFILE.jsonl]",
				Usage:     "remove whitespace from data",
				Action:    internal.CmdInit,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "fields",
						Aliases:  []string{"f"},
						EnvVars:  []string{"AMBROSIA_FIELDS", "FIELDS"},
						Usage:    "the comma-separated json `FIELD`(s) to remove whitespace from",
						Required: true,
						Category: "required:",
					},
				},
			},
			{
				Name:      "dedupe",
				ArgsUsage: "INFILE.jsonl [OUTFILE.jsonl]",
				Usage:     "remove duplicate content",
				Action:    internal.CmdInit,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "fields",
						Aliases:  []string{"f"},
						EnvVars:  []string{"AMBROSIA_FIELDS", "FIELDS"},
						Usage:    "the comma-separated json `FIELD`(s) in each piece of data to compare.",
						Required: true,
						Category: "required:",
					},
					&cli.BoolFlag{
						Name:    "ignore-case",
						Aliases: []string{"i"},
						EnvVars: []string{"AMBROSIA_IGNORE_CASE", "IGNORE_CASE"},
						Usage:   "ignore case when comparing lines",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "rougel",
						Aliases: []string{"rl"},
						EnvVars: []string{"AMBROSIA_RL", "RL"},
						Usage:   "use ROUGE-L to compare lines",
						Value:   false,
					},
					&cli.Float64Flag{
						Name:    "rl-threshold",
						Aliases: []string{"rlt"},
						EnvVars: []string{"AMBROSIA_RLT", "RLT"},
						Usage:   "if --rougel is set, the threshold for comparison",
						Value:   0.7,
					},
					&cli.BoolFlag{
						Name:    "progress",
						Aliases: []string{"p"},
						EnvVars: []string{"AMBROSIA_PROGRESS", "PROGRESS"},
						Usage:   "show progress bar",
						Value:   false,
					},
				},
			},
			{
				Name:      "length",
				ArgsUsage: "INFILE.jsonl [OUTFILE.jsonl]",
				Usage:     "filter data based on length",
				Action:    internal.CmdInit,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "fields",
						Aliases:  []string{"f"},
						EnvVars:  []string{"AMBROSIA_FIELDS", "FIELDS"},
						Usage:    "the comma-separated json `FIELD`(s) to check length of, multiple fields will be summed",
						Required: true,
						Category: "required:",
					},
					&cli.IntFlag{
						Name:        "min",
						EnvVars:     []string{"AMBROSIA_MIN", "MIN"},
						Usage:       "the minimum length of a field (<=)",
						DefaultText: "nil",
					},
					&cli.IntFlag{
						Name:        "max",
						EnvVars:     []string{"AMBROSIA_MAX", "MAX"},
						Usage:       "the maximum length of a field (>=)",
						DefaultText: "nil",
					},
				},
			},
			{
				Name:      "filter",
				ArgsUsage: "INFILE.jsonl [OUTFILE.jsonl]",
				Usage:     "filter data based on a string or wordlist",
				Action:    internal.CmdInit,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "fields",
						Aliases:  []string{"f"},
						EnvVars:  []string{"AMBROSIA_FIELDS", "FIELDS"},
						Usage:    "the comma-separated json `FIELD`(s) in each piece of data to filter on",
						Required: true,
						Category: "required:",
					},
					&cli.BoolFlag{
						Name:     "regex",
						Aliases:  []string{"r"},
						EnvVars:  []string{"AMBROSIA_REGEX", "REGEX"},
						Usage:    "treat string as a regular expression",
						Value:    false,
						Category: "optional:",
					},
					&cli.StringFlag{
						Name:      "wordlist",
						Aliases:   []string{"w"},
						EnvVars:   []string{"AMBROSIA_WORDLIST", "WORDLIST"},
						Usage:     "a wordlist file to filter on",
						TakesFile: true,
						Category:  "type:",
					},
					&cli.StringFlag{
						Name:     "string",
						Aliases:  []string{"s"},
						EnvVars:  []string{"AMBROSIA_STRING", "STRING"},
						Usage:    "a string to filter on",
						Category: "type:",
					},
				},
			},
			{
				Name:      "psort",
				ArgsUsage: "INFILE.jsonl",
				Usage:     "sort data based on LLM prompts and responses",
				Action:    internal.CmdInit,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "model",
						Aliases: []string{"m"},
						EnvVars: []string{"AMBROSIA_MODEL", "MODEL"},
						Usage:   "the `MODEL` to use, supported: ['gpt-3.5-turbo', 'gpt-4']",
						Value:   "gpt-3.5-turbo",
					},
					&cli.StringFlag{
						Name:    "token",
						Aliases: []string{"t"},
						EnvVars: []string{"AMBROSIA_TOKEN", "TOKEN"},
						Usage:   "the auth `TOKEN` to use with models that require it",
					},
					&cli.StringFlag{
						Name:    "baseurl",
						Aliases: []string{"b"},
						EnvVars: []string{"AMBROSIA_BASEURL", "BASEURL"},
						Usage:   "the `BASEURL` to use with models that require it, if not the default",
					},
					&cli.StringFlag{
						Name:    "instruction",
						Aliases: []string{"i"},
						EnvVars: []string{"AMBROSIA_INSTRUCTION", "INSTRUCTION"},
						Usage:   "the `INSTRUCTION` to use for inference",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "end-instruction",
						Aliases: []string{"e"},
						EnvVars: []string{"AMBROSIA_END_INSTRUCTION", "END_INSTRUCTION"},
						Usage:   "`END_INSTRUCTION` will be placed after the data in the prompt",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "sysprompt",
						Aliases: []string{"sp"},
						EnvVars: []string{"AMBROSIA_SYSPROMPT", "SYSPROMPT"},
						Usage:   "the `SYSPROMPT` to use for inference, on models that support it",
						Value:   "",
					},
					&cli.StringSliceFlag{
						Name:    "fields",
						Aliases: []string{"f"},
						EnvVars: []string{"AMBROSIA_FIELDS", "FIELDS"},
						Usage:   "the json `FIELD`(s) to use for prompts.  All fields used in random order if not specified.",
					},
					&cli.BoolFlag{
						Name:    "json",
						Aliases: []string{"j"},
						EnvVars: []string{"AMBROSIA_JSON", "JSON"},
						Usage:   "send data portion of prompt as a json object, instead of a string with fields",
					},
					&cli.BoolFlag{
						Name:    "include-resp",
						Aliases: []string{"ir"},
						EnvVars: []string{"AMBROSIA_INCLUDE_RESP", "INCLUDE_RESP"},
						Usage:   "include the LLM response as a new field in the output entry named 'ambrosia'",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "dry-run",
						Aliases: []string{"d"},
						EnvVars: []string{"AMBROSIA_DRY_RUN", "DRY_RUN"},
						Usage:   "don't actually perform inference, just print the prompts that would be used",
						Value:   false,
					},
					&cli.IntFlag{
						Name:    "concurrency",
						Aliases: []string{"c"},
						EnvVars: []string{"AMBROSIA_CONCURRENCY", "CONCURRENCY"},
						Usage:   "the number of concurrent requests to make to the model",
						Value:   10,
					},
					&cli.IntFlag{
						Name:    "rpm",
						EnvVars: []string{"AMBROSIA_RPM", "RPM"},
						Usage:   "the maximum number of requests per minute to make to the model",
						Value:   3150,
					},
					&cli.IntFlag{
						Name:    "tpm",
						EnvVars: []string{"AMBROSIA_TPM", "TPM"},
						Usage:   "the maximum number of tokens per minute to and from the model",
						Value:   81000,
					},
					&cli.IntFlag{
						Name:    "max-tokens",
						Aliases: []string{"mt"},
						EnvVars: []string{"AMBROSIA_MAX_TOKENS", "MAX_TOKENS"},
						Usage:   "the (requested) max number of tokens to generate, 0 for unlimited",
						Value:   5,
					},
					&cli.DurationFlag{
						Name:    "timeout",
						Aliases: []string{"to"},
						EnvVars: []string{"AMBROSIA_TIMEOUT", "TIMEOUT"},
						Usage:   "the maximum amount of time to wait for a response from the model before retrying",
						Value:   15 * time.Second,
					},
					&cli.BoolFlag{
						Name:    "progress",
						Aliases: []string{"p"},
						EnvVars: []string{"AMBROSIA_PROGRESS", "PROGRESS"},
						Usage:   "show progress bar",
						Value:   false,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
