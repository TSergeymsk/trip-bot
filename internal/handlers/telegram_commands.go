package handlers

import (
	"context"
	"log"

	"trip-bot/internal/services"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RegisterCommands(b *bot.Bot, updateSvc *services.UpdateService) {
	// Команда /status
	b.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			// Можно показать текущий хеш или просто ответить, что бот работает
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "🤖 Бот работает, отслеживает изменения.",
			})
		})

	// Команда /refresh - принудительное обновление
	b.RegisterHandler(bot.HandlerTypeMessageText, "/refresh", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if err := updateSvc.Refresh(ctx); err != nil {
				log.Printf("Refresh error: %v", err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "❌ Ошибка обновления: " + err.Error(),
				})
				return
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "✅ Данные обновлены",
			})
		})
}