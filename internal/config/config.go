package config

import (
	"os"
	"strconv"
)

type Config struct {
	TelegramToken   string
	GoogleCredsFile string
	SheetID         string
	SheetRange      string
	ChatID          int64
	ThreadID        int
	WebhookSecret   string
	Port            string
	DBPath          string
}

func Load() *Config {
	return &Config{
		TelegramToken:   getenv("TELEGRAM_TOKEN", ""),
		GoogleCredsFile: getenv("GOOGLE_CREDENTIALS", ""),
		SheetID:         getenv("SHEET_ID", ""),
		SheetRange:      getenv("SHEET_RANGE", "A1:E20"),
		ChatID:          mustParseInt64(getenv("CHAT_ID", "0")),
		ThreadID:        mustParseInt(getenv("THREAD_ID", "0")),
		WebhookSecret:   getenv("WEBHOOK_SECRET", ""),
		Port:            getenv("PORT", "8080"),
		DBPath:          getenv("DB_PATH", "state.db"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustParseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func mustParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}