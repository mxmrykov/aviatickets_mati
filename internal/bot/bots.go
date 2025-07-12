package bot

import (
	"context"
	"log"

	"aviatickets_mati/internal/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AviaTicketsBot struct {
	bot       *tgbotapi.BotAPI
	userUC    *usecase.UserUseCase
	flightUC  *usecase.FlightUseCase
	bookingUC *usecase.BookingUseCase
	adminKey  string // Добавляем поле для хранения ключа
}

func NewAviaTicketsBot(
	bot *tgbotapi.BotAPI,
	userUC *usecase.UserUseCase,
	flightUC *usecase.FlightUseCase,
	bookingUC *usecase.BookingUseCase,
	adminKey string,
) *AviaTicketsBot {
	return &AviaTicketsBot{
		bot:       bot,
		userUC:    userUC,
		flightUC:  flightUC,
		bookingUC: bookingUC,
		adminKey:  adminKey,
	}
}

func (b *AviaTicketsBot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		user, err := b.userUC.GetOrCreate(
			ctx,
			update.Message.From.ID,
			update.Message.From.UserName,
			update.Message.From.FirstName,
			update.Message.From.LastName,
		)
		if err != nil {
			log.Printf("Error getting user: %v", err)
			continue
		}

		switch update.Message.Command() {
		case "start":
			b.handleStart(update.Message.Chat.ID, user)
		case "search":
			b.handleSearchCommand(ctx, update.Message)
		case "book":
			b.handleBookCommand(ctx, update.Message, user.ID)
		case "mybookings":
			b.handleMyBookingsCommand(ctx, update.Message, user.ID)
		case "cancel":
			b.handleCancelCommand(ctx, update.Message, user.ID)
		case "admin":
			b.handleAdminCommand(ctx, update.Message)
		default:
			b.sendMessage(update.Message.Chat.ID, "Неизвестная команда")
		}
	}
}

func (b *AviaTicketsBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
