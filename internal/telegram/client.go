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

	params := &bot.EditMessageTextParams{
		ChatID:          chatID,
		MessageThreadID: threadID, // теперь поле поддерживается
		MessageID:       msgID,
		Text:            text,
		ParseMode:       "Markdown",
	}
	_, err := c.bot.EditMessageText(ctx, params)
	return msgID, err
}