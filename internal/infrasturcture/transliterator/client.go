package transliterator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/go-resty/resty/v2"
)

type apiClinet struct {
	config     *config.Config
	httpClient *resty.Client
}

func New(cfg *config.Config) (Client, error) {
	contextDeadline, err := time.ParseDuration(cfg.Transliterator.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeout duration: %w", err)
	}

	client := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s:%s", cfg.Transliterator.Host, cfg.Transliterator.Port)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTimeout(contextDeadline)

	return &apiClinet{
		config:     cfg,
		httpClient: client,
	}, nil
}

func (c *apiClinet) Transliterate(ctx context.Context, text string) (string, error) {
	data := map[string]any{
		"text": text,
	}

	response, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(data).
		Post("/transliterate")
	if err != nil {
		return "", fmt.Errorf("failed to transliterate: %w", err)
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("failed to transliterate: %s", response.Status())
	}

	var result struct {
		Transliteration string `json:"transliteration"`
	}

	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to transliterate: %w", err)
	}

	return result.Transliteration, nil
}
