package providers

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/sashabaranov/go-openai"
)

type OAIModel string

const (
	ModelGPT3Dot5Turbo OAIModel = openai.GPT3Dot5Turbo
	ModelGPT4          OAIModel = openai.GPT4
)

type OAI struct {
	c         *openai.Client
	model     OAIModel
	token     string
	logger    zerolog.Logger
	maxtokens int
}

type OAIConfig struct {
	Token     string
	BaseURL   string
	Timeout   time.Duration
	Model     OAIModel
	Logger    zerolog.Logger
	MaxTokens int
}

func NewOAI(c OAIConfig) *OAI {
	clientConf := openai.DefaultConfig(c.Token)
	if c.BaseURL != "" {
		clientConf.BaseURL = c.BaseURL
	}
	if c.Timeout != 0 {
		clientConf.HTTPClient = &http.Client{
			Timeout: c.Timeout,
		}
	}

	client := openai.NewClientWithConfig(clientConf)

	return &OAI{
		c:         client,
		model:     c.Model,
		token:     c.Token,
		logger:    c.Logger,
		maxtokens: c.MaxTokens,
	}
}

func (o *OAI) Infer(req *InferRequest) (*InferResponse, error) {
	// Build messages with optional sysprompt
	var messages []openai.ChatCompletionMessage
	if req.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: req.SystemPrompt,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Prompt,
	})

	// Fire off request
	resp, err := o.c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     string(o.model),
			Messages:  messages,
			MaxTokens: o.maxtokens,
		},
	)
	o.logger.Debug().
		Interface("resp", resp).
		Err(err).
		Msg("response from oai")

	if err != nil {
		return nil, err
	}

	// Parse response
	ret := &InferResponse{
		ID:     req.ID,
		Resp:   resp.Choices[0].Message.Content,
		Tokens: resp.Usage.TotalTokens,
	}

	return ret, nil
}

func (o *OAI) Ping() error {
	o.logger.Debug().Msg("pinging openai")
	_, err := o.c.ListModels(context.Background())
	o.logger.Debug().Err(err).Msg("pinged openai")
	return err
}
