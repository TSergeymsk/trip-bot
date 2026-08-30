package interfaces

import "context"

type DataProvider interface {
	GetData(ctx context.Context) (raw string, hash string, err error)
}

type StateRepository interface {
	Get(sheetID string, threadID int) (hash string, msgID int, err error)
	Save(sheetID string, threadID int, hash string, msgID int) error
}

type TelegramNotifier interface {
	SendOrUpdate(ctx context.Context, chatID int64, threadID int, msgID int, text string) (int, error)
}