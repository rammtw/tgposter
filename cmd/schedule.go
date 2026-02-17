package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rammtw/tgposter/internal/config"
	"github.com/rammtw/tgposter/internal/converter"
	"github.com/rammtw/tgposter/internal/scheduler"

	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Запланировать публикацию .md файла на определённое время",
	Example: `  tg-poster schedule --file article.md --channel @my_channel --time "2026-02-18 14:00"
  tg-poster schedule -f post.md -c @nom_nom_notes -T "2026-02-18 09:30"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		channel, _ := cmd.Flags().GetString("channel")
		timeStr, _ := cmd.Flags().GetString("time")
		tz, _ := cmd.Flags().GetString("tz")

		cfg, err := config.Load(cmd)
		if err != nil {
			return err
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			return fmt.Errorf("неизвестная таймзона %q: %w", tz, err)
		}

		postTime, err := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
		if err != nil {
			return fmt.Errorf("неверный формат времени (ожидается YYYY-MM-DD HH:MM): %w", err)
		}

		if postTime.Before(time.Now()) {
			return fmt.Errorf("время публикации %s уже прошло", postTime.Format(time.RFC3339))
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("не удалось прочитать файл: %w", err)
		}

		tgText := converter.MarkdownToTelegram(string(content))

		s := scheduler.New(cfg.Token)
		s.Schedule(context.Background(), channel, tgText, postTime)

		fmt.Printf("⏰ Публикация запланирована в %s → %s\n",
			postTime.Format("2006-01-02 15:04 MST"), channel)
		fmt.Println("Ожидание... (нажмите Ctrl+C для отмены)")

		// Блокируем до завершения
		select {}
	},
}

func init() {
	scheduleCmd.Flags().StringP("file", "f", "", "Путь к .md файлу (обязательный)")
	scheduleCmd.Flags().StringP("channel", "c", "", "Канал, например @my_channel (обязательный)")
	scheduleCmd.Flags().StringP("time", "T", "", "Время публикации: \"2026-02-18 14:00\" (обязательный)")
	scheduleCmd.Flags().String("tz", "Europe/Moscow", "Таймзона (по умолчанию Europe/Moscow)")
	_ = scheduleCmd.MarkFlagRequired("file")
	_ = scheduleCmd.MarkFlagRequired("channel")
	_ = scheduleCmd.MarkFlagRequired("time")
	rootCmd.AddCommand(scheduleCmd)
}
