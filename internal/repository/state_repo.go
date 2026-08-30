package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type StateRepo struct {
	db *sql.DB
}

func NewStateRepo(dbPath string) (*StateRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Создаём таблицу
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS state (
			sheet_id TEXT,
			thread_id INTEGER,
			last_hash TEXT,
			message_id INTEGER,
			PRIMARY KEY (sheet_id, thread_id)
		)
	`)
	if err != nil {
		return nil, err
	}
	return &StateRepo{db: db}, nil
}

func (r *StateRepo) Get(sheetID string, threadID int) (hash string, msgID int, err error) {
	row := r.db.QueryRow(
		"SELECT last_hash, message_id FROM state WHERE sheet_id=? AND thread_id=?",
		sheetID, threadID,
	)
	err = row.Scan(&hash, &msgID)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}
	return
}

func (r *StateRepo) Save(sheetID string, threadID int, hash string, msgID int) error {
	_, err := r.db.Exec(
		"INSERT OR REPLACE INTO state (sheet_id, thread_id, last_hash, message_id) VALUES (?, ?, ?, ?)",
		sheetID, threadID, hash, msgID,
	)
	return err
}