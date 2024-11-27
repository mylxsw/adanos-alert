package llm

import "context"

type ChatResponse struct {
	Content string
}

type LLM interface {
	Chat(ctx context.Context, model, prompt string) (*ChatResponse, error)
	Summarize(ctx context.Context, model string, text string) (*ChatResponse, error)
}
