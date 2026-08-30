package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-telegram/bot"
)

type Client struct {
	bot     *bot.Bot
	httpCli *http.Client
	token   string
}

func New(token string) (*Client, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, err
	}
	return &Client{
		bot:     b,
		httpCli: &http.Client{},
		token:   token,
	}, nil
}

func (c *Client) SendOrUpdate(ctx context.Context, chatID int64, threadID int, msgID int, text string) (int, error) {
	if msgID == 0 {
		// Отправка нового сообщения через библиотеку (работает)
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

	// Редактирование через прямой HTTP-запрос
	url := fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", c.token)
	payload := map[string]interface{}{
		"chat_id":           chatID,
		"message_id":        msgID,
		"text":              text,
		"parse_mode":        "Markdown",
		"message_thread_id": threadID,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			MessageID int `json:"message_id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if !result.Ok {
		return 0, fmt.Errorf("telegram API error")
	}
	// ID сообщения при редактировании не меняется
	return msgID, nil
}