package cmd

import (
	"fmt"

	"github.com/rammtw/tgposter/internal/api"
	"github.com/rammtw/tgposter/internal/config"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запустить HTTP API сервер для постинга",
	Example: `  tg-poster serve --port 8080
  tg-poster serve -p 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")

		cfg, err := config.Load(cmd)
		if err != nil {
			return err
		}

		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("🚀 API сервер запущен на http://localhost%s\n", addr)
		return api.ListenAndServe(addr, cfg.Token)
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "Порт для API сервера")
	rootCmd.AddCommand(serveCmd)
}
