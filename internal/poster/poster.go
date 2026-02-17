package poster

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Poster struct {
	bot *bot.Bot
}

func New(token string) (*Poster, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}
	return &Poster{bot: b}, nil
}

// Send публикует текст в указанный канал и возвращает message_id.
func (p *Poster) Send(ctx context.Context, channel, text string) (int, error) {
	msg, err := p.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    channel,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return 0, fmt.Errorf("telegram API error: %w", err)
	}
	return msg.ID, nil
}
