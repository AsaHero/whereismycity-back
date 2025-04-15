package embeddings

import "context"

type Client interface {
	Generate(ctx context.Context, text string) ([]float64, error)
}
