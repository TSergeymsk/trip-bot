package services

import (
	"context"
	"log"

	"trip-bot/internal/interfaces"
)

type UpdateService struct {
	provider   interfaces.DataProvider
	stateRepo  interfaces.StateRepository
	notifier   interfaces.TelegramNotifier
	chatID     int64
	threadID   int
	sheetID    string
	formatter  func(string) string
}

func NewUpdateService(
	provider interfaces.DataProvider,
	stateRepo interfaces.StateRepository,
	notifier interfaces.TelegramNotifier,
	chatID int64,
	threadID int,
	sheetID string,
	formatter func(string) string,
) *UpdateService {
	return &UpdateService{
		provider:  provider,
		stateRepo: stateRepo,
		notifier:  notifier,
		chatID:    chatID,
		threadID:  threadID,
		sheetID:   sheetID,
		formatter: formatter,
	}
}

func (u *UpdateService) Refresh(ctx context.Context) error {
	raw, currentHash, err := u.provider.GetData(ctx)
	if err != nil {
		return err
	}

	lastHash, msgID, err := u.stateRepo.Get(u.sheetID, u.threadID)
	if err != nil {
		return err
	}

	if currentHash == lastHash && msgID != 0 {
		log.Println("No changes detected")
		return nil
	}

	text := u.formatter(raw)
	newMsgID, err := u.notifier.SendOrUpdate(ctx, u.chatID, u.threadID, msgID, text)
	if err != nil {
		return err
	}

	if err := u.stateRepo.Save(u.sheetID, u.threadID, currentHash, newMsgID); err != nil {
		return err
	}

	log.Println("Message updated successfully")
	return nil
}