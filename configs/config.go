package configs

import (
	"os"
	"strconv"
)

type Config struct {
	DB       DBConfig
	Telegram TelegramConfig
	AdminKey string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type TelegramConfig struct {
	BotToken string
	Debug    bool
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	debug, err := strconv.ParseBool(getEnv("TELEGRAM_DEBUG", "false"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "db_main"),
			Port:     port,
			User:     getEnv("DB_USER", "avia_user"),
			Password: getEnv("DB_PASSWORD", "avia_pass"),
			Name:     getEnv("DB_NAME", "aviatickets"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Telegram: TelegramConfig{
			BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
			Debug:    debug,
		},
		AdminKey: getEnv("ADMIN_KEY", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
