package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetService struct {
	srv     *sheets.Service
	sheetID string
	range_  string
}

func NewSheetService(credsFile, sheetID, sheetRange string) (*SheetService, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		return nil, err
	}
	return &SheetService{srv: srv, sheetID: sheetID, range_: sheetRange}, nil
}

func (s *SheetService) GetData(ctx context.Context) (raw string, hash string, err error) {
	resp, err := s.srv.Spreadsheets.Values.Get(s.sheetID, s.range_).Context(ctx).Do()
	if err != nil {
		return "", "", err
	}

	var sb strings.Builder
	for _, row := range resp.Values {
		for _, cell := range row {
			sb.WriteString(fmt.Sprintf("%v", cell))
		}
	}
	raw = sb.String()
	hasher := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(hasher[:])
	return raw, hash, nil
}