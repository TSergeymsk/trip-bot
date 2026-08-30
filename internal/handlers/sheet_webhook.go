package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"trip-bot/internal/services"
)

type SheetWebhookHandler struct {
	updateSvc *services.UpdateService
	secret    string
}

func NewSheetWebhookHandler(updateSvc *services.UpdateService, secret string) *SheetWebhookHandler {
	return &SheetWebhookHandler{updateSvc: updateSvc, secret: secret}
}

func (h *SheetWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("X-Webhook-Secret") != h.secret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var payload map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Failed to parse payload: %v", err)
		}
		// можно извлечь sheetId, threadId и т.д. если нужно
	}

	// Асинхронно обновляем
	go func() {
		ctx := r.Context()
		if err := h.updateSvc.Refresh(ctx); err != nil {
			log.Printf("Update failed: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}