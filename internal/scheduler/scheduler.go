package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/rammtw/tgposter/internal/poster"
)

type Scheduler struct {
	token string
}

func New(token string) *Scheduler {
	return &Scheduler{token: token}
}

func (s *Scheduler) Schedule(ctx context.Context, channel, text string, postTime time.Time) {
	delay := time.Until(postTime)

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-timer.C:
			p, err := poster.New(s.token)
			if err != nil {
				fmt.Printf("❌ Ошибка инициализации бота: %v\n", err)
				return
			}
			msgID, err := p.Send(ctx, channel, text)
			if err != nil {
				fmt.Printf("❌ Ошибка отправки: %v\n", err)
				return
			}
			fmt.Printf("✅ Отложенное сообщение опубликовано в %s (message_id: %d)\n", channel, msgID)
		case <-ctx.Done():
			fmt.Println("⚠️  Публикация отменена")
		}
	}()
}
