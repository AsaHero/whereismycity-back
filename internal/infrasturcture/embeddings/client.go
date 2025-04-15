package embeddings

import (
	"context"
	"fmt"
	"time"

	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type apiClient struct {
	cfg    *config.Config
	client *openai.Client
}

func New(cfg *config.Config) (Client, error) {
	timeoutDuration, err := time.ParseDuration(cfg.OpenAI.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeout duration: %w", err)
	}

	client := openai.NewClient(option.WithAPIKey(cfg.OpenAI.APIKey), option.WithRequestTimeout(timeoutDuration))
	return &apiClient{
		cfg:    cfg,
		client: &client,
	}, nil
}

func (c *apiClient) Generate(ctx context.Context, text string) ([]float64, error) {
	response, err := c.client.Embeddings.New(
		ctx,
		openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(text),
			},
			Model: openai.EmbeddingModelTextEmbedding3Small,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	return response.Data[0].Embedding, nil
}
