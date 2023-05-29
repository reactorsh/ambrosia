package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfer(t *testing.T) {
	d := NewDryRun()

	// Test case 1: Check if function correctly populates response.
	req := &InferRequest{
		ID:           1,
		SystemPrompt: "system prompt",
		Prompt:       "test prompt",
	}

	expectedResponse := &InferResponse{
		ID:     req.ID,
		Resp:   req.Prompt,
		Tokens: len(req.Prompt) * 2,
	}

	resp, err := d.Infer(req)

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)

	// Test case 2: Check if function correctly handles empty prompts.
	req = &InferRequest{
		ID:           2,
		SystemPrompt: "system prompt",
		Prompt:       "",
	}

	expectedResponse = &InferResponse{
		ID:     req.ID,
		Resp:   req.Prompt,
		Tokens: len(req.Prompt) * 2,
	}

	resp, err = d.Infer(req)

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestPing(t *testing.T) {
	d := NewDryRun()

	// Since Ping() just returns nil, we only need to check if it returns error as nil.
	err := d.Ping()

	assert.Nil(t, err)
}
