package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"aviatickets_mati/pkg/postgres"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Использование: go run main.go <количество_рейсов>")
	}

	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Неверное количество рейсов: %v", err)
	}

	// Загрузка конфигурации БД (можно вынести в отдельный конфиг)
	cfg := postgres.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "avia_user"),
		Password: getEnv("DB_PASSWORD", "avia_pass"),
		Name:     getEnv("DB_NAME", "aviatickets"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Подключение к БД
	ctx := context.Background()
	db, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Генерация рейсов
	generator := NewFlightGenerator(db)
	if err := generator.GenerateFlights(count); err != nil {
		log.Fatalf("Ошибка генерации рейсов: %v", err)
	}

	log.Printf("Успешно сгенерировано %d рейсов", count)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
