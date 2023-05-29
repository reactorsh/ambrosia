package providers

import "fmt"

type DryRun struct{}

func (d *DryRun) Infer(req *InferRequest) (*InferResponse, error) {
	fmt.Printf(
		"--BEGIN\nSystem Prompt: %s\nPrompt: %s\n--END\n",
		req.SystemPrompt,
		req.Prompt,
	)

	return &InferResponse{
		ID:     req.ID,
		Resp:   req.Prompt,
		Tokens: len(req.Prompt) * 2,
	}, nil
}

func (d *DryRun) Ping() error {
	return nil
}

func NewDryRun() *DryRun {
	return &DryRun{}
}
