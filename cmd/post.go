package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rammtw/tgposter/internal/config"
	"github.com/rammtw/tgposter/internal/converter"
	"github.com/rammtw/tgposter/internal/poster"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Немедленно опубликовать .md файл в канал",
	Example: `  tg-poster post --file article.md --channel @my_channel
  tg-poster post -f post.md -c @nom_nom_notes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		channel, _ := cmd.Flags().GetString("channel")

		cfg, err := config.Load(cmd)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("не удалось прочитать файл: %w", err)
		}

		tgText := converter.MarkdownToTelegram(string(content))

		p, err := poster.New(cfg.Token)
		if err != nil {
			return fmt.Errorf("не удалось инициализировать бота: %w", err)
		}

		msgID, err := p.Send(context.Background(), channel, tgText)
		if err != nil {
			return fmt.Errorf("ошибка отправки: %w", err)
		}

		fmt.Printf("✅ Сообщение опубликовано в %s (message_id: %d)\n", channel, msgID)
		return nil
	},
}

func init() {
	postCmd.Flags().StringP("file", "f", "", "Путь к .md файлу (обязательный)")
	postCmd.Flags().StringP("channel", "c", "", "Название канала, например @my_channel (обязательный)")
	_ = postCmd.MarkFlagRequired("file")
	_ = postCmd.MarkFlagRequired("channel")
	rootCmd.AddCommand(postCmd)
}
