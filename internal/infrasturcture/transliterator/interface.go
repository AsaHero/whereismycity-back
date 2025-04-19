package transliterator

import "context"

type Client interface {
	Transliterate(ctx context.Context, text string) (string, error)
}
