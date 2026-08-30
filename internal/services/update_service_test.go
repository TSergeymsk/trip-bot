package services

import (
	"context"
	"errors"
	"testing"
)

type mockProvider struct {
	raw  string
	hash string
	err  error
}

func (m *mockProvider) GetData(ctx context.Context) (string, string, error) {
	return m.raw, m.hash, m.err
}

type mockStateRepo struct {
	hash   string
	msgID  int
	saveErr error
}

func (m *mockStateRepo) Get(sheetID string, threadID int) (string, int, error) {
	return m.hash, m.msgID, nil
}
func (m *mockStateRepo) Save(sheetID string, threadID int, hash string, msgID int) error {
	return m.saveErr
}

type mockNotifier struct {
	sendErr error
}

func (m *mockNotifier) SendOrUpdate(ctx context.Context, chatID int64, threadID int, msgID int, text string) (int, error) {
	if m.sendErr != nil {
		return 0, m.sendErr
	}
	return 999, nil
}

func TestUpdateService_Refresh(t *testing.T) {
	ctx := context.Background()
	formatter := func(s string) string { return "formatted: " + s }

	t.Run("no changes, msg exists", func(t *testing.T) {
		provider := &mockProvider{raw: "data", hash: "hash1"}
		state := &mockStateRepo{hash: "hash1", msgID: 123}
		notifier := &mockNotifier{}
		svc := NewUpdateService(provider, state, notifier, 1, 2, "sheet", formatter)

		err := svc.Refresh(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("changes detected", func(t *testing.T) {
		provider := &mockProvider{raw: "new", hash: "hash2"}
		state := &mockStateRepo{hash: "hash1", msgID: 123}
		notifier := &mockNotifier{}
		svc := NewUpdateService(provider, state, notifier, 1, 2, "sheet", formatter)

		err := svc.Refresh(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("provider error", func(t *testing.T) {
		provider := &mockProvider{err: errors.New("sheet error")}
		state := &mockStateRepo{hash: "hash1", msgID: 123}
		notifier := &mockNotifier{}
		svc := NewUpdateService(provider, state, notifier, 1, 2, "sheet", formatter)

		err := svc.Refresh(ctx)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("notifier error", func(t *testing.T) {
		provider := &mockProvider{raw: "new", hash: "hash2"}
		state := &mockStateRepo{hash: "hash1", msgID: 123}
		notifier := &mockNotifier{sendErr: errors.New("tg error")}
		svc := NewUpdateService(provider, state, notifier, 1, 2, "sheet", formatter)

		err := svc.Refresh(ctx)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}