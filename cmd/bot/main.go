package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aviatickets_mati/configs"
	"aviatickets_mati/internal/bot"
	"aviatickets_mati/internal/repository"
	"aviatickets_mati/internal/usecase"
	"aviatickets_mati/pkg/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Загрузка конфигурации
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация базы данных
	pgPool, err := postgres.NewPool(ctx, postgres.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgPool.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(pgPool)
	flightRepo := repository.NewFlightRepository(pgPool)
	bookingRepo := repository.NewBookingRepository(pgPool)

	// Инициализация use cases
	userUC := usecase.NewUserUseCase(userRepo)
	flightUC := usecase.NewFlightUseCase(flightRepo)
	bookingUC := usecase.NewBookingUseCase(bookingRepo, flightRepo)

	// Инициализация Telegram бота
	botAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	botAPI.Debug = cfg.Telegram.Debug
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	// Создание и запуск бота
	aviaBot := bot.NewAviaTicketsBot(botAPI, userUC, flightUC, bookingUC, cfg.AdminKey)
	go aviaBot.Start(ctx)

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second) // Даем время на завершение работы
}
