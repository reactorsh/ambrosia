package internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/reactorsh/ambrosia/limiter"
	"github.com/reactorsh/ambrosia/providers"
	"github.com/schollz/progressbar/v3"
)

func cmdPSort(c *cmdCtx) error {
	c.logger.Info().Msg("checking for resumable outputs in output path")

	datumCompleted, err := loadResumable(c.c.Command.Name, c.inPath)
	if err != nil {
		return err
	}

	var todo []datum
	if len(datumCompleted) > 0 {
		c.logger.Debug().Int("count", len(datumCompleted)).Msg("resumable data found")
		todo = datumSub(c.data, datumCompleted)
		c.logger.Info().Int("remaining", len(todo)).Msg("resuming")
	} else {
		c.logger.Debug().
			Msg("no resumable outputs found")
		todo = c.data
	}

	var prompter providers.Provider

	switch {

	case c.c.Bool("dry-run"):
		c.logger.Info().Msg("dry-run enabled")
		prompter = providers.NewDryRun()

	case c.c.String("model") == "gpt-3.5-turbo" || c.c.String("model") == "gpt-4":
		c.logger.Info().Str("model", c.c.String("model")).Msg("using openai model")
		prompter = oaiPrompter(c)

	default:
		return fmt.Errorf("dry-run not set and no valid model specified")
	}

	err = prompter.Ping()
	if err != nil {
		return fmt.Errorf("error with model: %w", err)
	}

	reqC := make(chan providers.InferRequest, c.c.Int("concurrency")*2)
	go submitPrompts(c, reqC, todo)

	respC := make(chan providers.InferResponse, c.c.Int("concurrency")*2)
	go inference(c, prompter, reqC, respC)

	pbar := progressbar.DefaultSilent(0)
	if c.c.Bool("progress") {
		pbar = progressbar.Default(int64(len(todo)))
	}

	appender := newPrefixAppender(prefixPathTmpl(c.inPath))
	defer appender.close()
	for resp := range respC {
		if c.c.Bool("dry-run") {
			continue
		}

		var err error
		if c.c.Bool("include-resp") {
			err = appender.appendWithResponse(resp.Resp, todo[resp.ID])
		} else {
			err = appender.append(resp.Resp, todo[resp.ID])
		}

		if err != nil {
			return fmt.Errorf("error appending: %w", err)
		}

		if c.c.Bool("progress") {
			pbar.Add(1)
		}
	}

	return nil
}

func inference(c *cmdCtx, p providers.Provider, in <-chan providers.InferRequest, out chan<- providers.InferResponse) {
	lim := limiter.New(c.c.Int("rpm"), c.c.Int("tpm"), &c.logger)

	wg := sync.WaitGroup{}
	wg.Add(c.c.Int("concurrency"))

	for i := 0; i < c.c.Int("concurrency"); i++ {
		go func() {
			defer wg.Done()
			for query := range in {
				c.logger.Debug().Interface("query", query).Msg("inference request")

				// Very naive retry, good enough for now.
				var resp *providers.InferResponse
				var err error
				for {
					if !c.c.Bool("dry-run") {
						lim.Wait(query.ByteCnt())
					}
					resp, err = p.Infer(&query)
					if err != nil {
						c.logger.Debug().
							Err(err).
							Interface("resp", resp).
							Msg("inference error, retrying")
						time.Sleep(1 * time.Second)
						continue
					}

					break
				}
				c.logger.Debug().Interface("resp", resp).Msg("inference response")
				if !c.c.Bool("dry-run") {
					lim.TPMReconcile(query.ByteCnt(), resp.Tokens)
				}
				out <- *resp
			}
		}()
	}

	wg.Wait()
	close(out)
}

func submitPrompts(c *cmdCtx, queue chan<- providers.InferRequest, data []datum) {
	for i, d := range data {
		var b strings.Builder

		// Handle prompt
		if c.c.String("instruction") != "" {
			fmt.Fprintf(&b, "%s\n\n", c.c.String("instruction"))
		}

		if c.c.Bool("json") {
			jb, err := d.JSON(c.c.StringSlice("fields"))
			if err != nil {
				c.logger.Fatal().Err(err).Msg("error marshalling json")
			}
			fmt.Fprintf(&b, "%s\n", string(jb))
		} else {
			// Handle 'all fields' case
			if len(c.c.StringSlice("fields")) == 0 {
				for field, value := range d {
					fmt.Fprintf(&b, "%s: %v\n", field, value)
				}
			}

			// Handle specific fields
			for _, field := range c.c.StringSlice("fields") {
				val, ok := d[field]
				if ok {
					fmt.Fprintf(&b, "%s: %v\n", field, val)
				}
			}
		}

		if c.c.String("end-instruction") != "" {
			fmt.Fprintf(&b, "\n%s\n", c.c.String("end-instruction"))
		}

		prompt := strings.TrimSpace(b.String())

		req := providers.InferRequest{
			ID:           i,
			SystemPrompt: c.c.String("sysprompt"),
			Prompt:       prompt,
		}

		// Enqueue built prompt
		queue <- req
	}
	close(queue)
}

func oaiPrompter(c *cmdCtx) *providers.OAI {
	oaiConf := providers.OAIConfig{
		Token:     c.c.String("token"), // Will error on the ping if invalid
		Logger:    c.logger,
		Timeout:   c.c.Duration("timeout"),
		MaxTokens: c.c.Int("max-tokens"),
	}

	var model providers.OAIModel
	switch c.c.String("model") {
	case "gpt-3.5-turbo":
		model = providers.ModelGPT3Dot5Turbo
	case "gpt-4":
		model = providers.ModelGPT4
	}
	oaiConf.Model = model

	if c.c.String("baseurl") != "" {
		oaiConf.BaseURL = c.c.String("baseurl")
	}

	return providers.NewOAI(oaiConf)
}
