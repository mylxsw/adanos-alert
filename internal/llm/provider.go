package llm

import (
	"github.com/mylxsw/adanos-alert/configs"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (s Provider) Register(app infra.Binder) {
	app.MustSingleton(func(conf *configs.Config) (LLM, error) {
		return NewOpenAILLM(OpenAIConfig{
			Endpoint:     conf.OpenAI.Endpoint,
			APIKey:       conf.OpenAI.APIKey,
			Organization: conf.OpenAI.Organization,
		})
	})
}
