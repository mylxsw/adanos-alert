package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIConfig struct {
	APIKey       string
	Endpoint     string
	Organization string
}

type OpenAILLM struct {
	config OpenAIConfig
	llm    *openai.LLM
}

func NewOpenAILLM(config OpenAIConfig) (LLM, error) {
	options := make([]openai.Option, 0)
	if config.APIKey != "" {
		options = append(options, openai.WithToken(config.APIKey))
	}

	if config.Endpoint != "" {
		options = append(options, openai.WithBaseURL(config.Endpoint))
	}

	if config.Organization != "" {
		options = append(options, openai.WithOrganization(config.Organization))
	}

	llm, err := openai.New(options...)
	if err != nil {
		return nil, err
	}

	return &OpenAILLM{config: config, llm: llm}, nil
}

func (o *OpenAILLM) Chat(ctx context.Context, model string, prompt string) (*ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	completion, err := o.llm.Call(ctx, prompt, llms.WithTemperature(0.8))
	if err != nil {
		return nil, err
	}

	return &ChatResponse{Content: completion}, nil
}

func (o *OpenAILLM) Summarize(ctx context.Context, model string, text string) (*ChatResponse, error) {
	return o.Chat(ctx, model, fmt.Sprintf(`Help me summarize this passage:\n\n<content>%s</content>`, text))
}
