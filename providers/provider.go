package providers

type InferRequest struct {
	ID           int
	SystemPrompt string
	Prompt       string
}

func (r *InferRequest) ByteCnt() int {
	return len(r.Prompt) + len(r.SystemPrompt)
}

type InferResponse struct {
	ID     int
	Resp   string
	Tokens int
}

type Provider interface {
	Infer(*InferRequest) (*InferResponse, error)
	Ping() error
}
