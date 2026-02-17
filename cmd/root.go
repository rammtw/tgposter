package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tg-poster",
	Short: "CLI утилита для постинга Markdown в Telegram каналы",
	Long:  "Читает .md файлы и публикует их содержимое в Telegram каналы с поддержкой отложенного постинга и REST API.",
}

func Execute() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "warning: .env file not found, using environment variables")
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("token", "t", "", "Telegram Bot API token (или переменная TELEGRAM_BOT_TOKEN)")
}
