package services

import "fmt"

func FormatMessage(data string) string {
	// Здесь вы можете парсить data (например, CSV или структурированный текст)
	// и возвращать красивое сообщение в Markdown.
	return fmt.Sprintf("📅 *Актуальные данные поездки:*\n\n%s", data)
}