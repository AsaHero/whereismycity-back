package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/go-resty/resty/v2"
)

type Bot struct {
	Token  string
	ChatID string
	client *resty.Client
}

func New(cfg *config.Config) (*Bot, error) {
	if cfg.Telegram.Token == "" || cfg.Telegram.ChatID == "" {
		return nil, fmt.Errorf("telegram token or chat id is empty")
	}

	client := resty.New().
		SetBaseURL(fmt.Sprintf("https://api.telegram.org/bot%s", cfg.Telegram.Token)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	return &Bot{
		Token:  cfg.Telegram.Token,
		ChatID: cfg.Telegram.ChatID,
		client: client,
	}, nil
}

func (b *Bot) SendContacts(ctx context.Context, req models.SendContactsRequest) error {
	fMessage := `
Contact Form Submission:

📅 *Date:* %s
👤 *Name:* %s
📧 *Email:* %s
🏢 *Company:* %s

📌 *Message:* %s
`

	msg := fmt.Sprintf(fMessage, time.Now().Format("2006-01-02 15:04"), req.Name, req.Email, req.Company, req.Message)

	r, err := b.client.R().
		SetContext(ctx).
		SetBody(map[string]string{
			"chat_id":    b.ChatID,
			"text":       msg,
			"parse_mode": "Markdown",
		}).
		Post("/sendMessage")

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if r.StatusCode() != 200 {
		return fmt.Errorf("failed to send message: %s", r.Status())
	}

	return err
}
