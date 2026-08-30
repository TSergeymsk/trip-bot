package telegram

import (
	"context"
	"github.com/go-telegram/bot"
)

type Client struct {
	bot *bot.Bot
}

func New(token string) (*Client, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}
	return &Client{bot: b}, nil
}

func (c *Client) SendOrUpdate(ctx context.Context, chatID int64, threadID int, msgID int, text string) (int, error) {
	if msgID == 0 {
		params := &bot.SendMessageParams{
			ChatID:          chatID,
			MessageThreadID: threadID,
			Text:            text,
			ParseMode:       "Markdown",
		}
		resp, err := c.bot.SendMessage(ctx, params)
		if err != nil {
			return 0, err
		}
		return resp.ID, nil
	}

	// Используем прямой вызов API для редактирования, чтобы гарантированно передать message_thread_id
	params := map[string]interface{}{
		"chat_id":           chatID,
		"message_id":        msgID,
		"text":              text,
		"parse_mode":        "Markdown",
		"message_thread_id": threadID,
	}
	var result struct{ Ok bool }
	err := c.bot.Request(ctx, "editMessageText", params, &result)
	if err != nil {
		return 0, err
	}
	return msgID, nil
}