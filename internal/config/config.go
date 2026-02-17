package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Token string
}

func Load(cmd *cobra.Command) (*Config, error) {
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("telegram Bot token не задан: используйте --token или TELEGRAM_BOT_TOKEN")
	}
	return &Config{Token: token}, nil
}
