package ai

import "github.com/sashabaranov/go-openai"

func NewClient(apiKey string, opts ...func(*openai.ClientConfig)) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	for _, opt := range opts {
		opt(&cfg)
	}
	return openai.NewClientWithConfig(cfg)
}

func WithOrgID(orgID string) func(*openai.ClientConfig) {
	return func(cfg *openai.ClientConfig) {
		cfg.OrgID = orgID
	}
}
