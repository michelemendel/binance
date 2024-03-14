package main

import (
	"log/slog"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/michelemendel/binance/tui"
)

// https://github.com/binance/binance-connector-go

func init() {
	envFile := filepath.Join("", ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		slog.Error("error loading file ", "file", envFile, "error", err)
	}
}

func main() {
	// client.Run()
	tui.Run()
}
