package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trip-bot/internal/config"
	"trip-bot/internal/handlers"
	"trip-bot/internal/repository"
	"trip-bot/internal/services"
	"trip-bot/internal/telegram"

	"github.com/go-telegram/bot"
)

func main() {
	cfg := config.Load()

	// Репозиторий состояния
	stateRepo, err := repository.NewStateRepo(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to init state repo: %v", err)
	}

	// Сервис Google Sheets
	sheetSvc, err := services.NewSheetService(cfg.GoogleCredsFile, cfg.SheetID, cfg.SheetRange)
	if err != nil {
		log.Fatalf("Failed to init sheet service: %v", err)
	}

	// Telegram клиент
	tgClient, err := telegram.New(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to init telegram client: %v", err)
	}

	// Сервис обновлений (основная логика)
	updateSvc := services.NewUpdateService(
		sheetSvc,
		stateRepo,
		tgClient,
		cfg.ChatID,
		cfg.ThreadID,
		cfg.SheetID,
		services.FormatMessage,
	)

	// HTTP обработчик вебхука
	webhookHandler := handlers.NewSheetWebhookHandler(updateSvc, cfg.WebhookSecret)

	// HTTP сервер
	httpSrv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/webhook/sheet-update" {
				webhookHandler.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	// Telegram long polling бот
	tgBot, err := bot.New(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create telegram bot: %v", err)
	}
	handlers.RegisterCommands(tgBot, updateSvc)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Telegram long polling")
		tgBot.Start(context.Background())
	}()

	<-stop
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	// tgBot.Stop() если требуется
	log.Println("Shutdown complete")
}