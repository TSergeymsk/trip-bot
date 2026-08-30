package repository

import (
	"os"
	"testing"
)

func TestStateRepo(t *testing.T) {
	tmpDB := "test_state.db"
	defer os.Remove(tmpDB)

	repo, err := NewStateRepo(tmpDB)
	if err != nil {
		t.Fatal(err)
	}

	sheetID := "test_sheet"
	threadID := 123

	// Пустое состояние
	hash, msgID, err := repo.Get(sheetID, threadID)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "" || msgID != 0 {
		t.Errorf("expected empty, got hash=%s msgID=%d", hash, msgID)
	}

	// Сохраняем
	err = repo.Save(sheetID, threadID, "abc123", 42)
	if err != nil {
		t.Fatal(err)
	}

	hash, msgID, err = repo.Get(sheetID, threadID)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "abc123" || msgID != 42 {
		t.Errorf("expected abc123/42, got %s/%d", hash, msgID)
	}

	// Обновляем
	err = repo.Save(sheetID, threadID, "def456", 99)
	if err != nil {
		t.Fatal(err)
	}
	hash, msgID, _ = repo.Get(sheetID, threadID)
	if hash != "def456" || msgID != 99 {
		t.Errorf("expected def456/99, got %s/%d", hash, msgID)
	}
}