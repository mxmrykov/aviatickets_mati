package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"aviatickets_mati/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Обработчик админ-команд
func (b *AviaTicketsBot) handleAdminCommand(ctx context.Context, msg *tgbotapi.Message) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) == 0 {
		b.sendMessage(msg.Chat.ID, "Использование:\n"+
			"/admin <ключ> - авторизация\n"+
			"/admin create_flight <данные> - создать рейс\n"+
			"/admin delete_flight <ID> - удалить рейс\n"+
			"/admin flight_stats <ID> - статистика по рейсу")
		return
	}

	// Проверка ключа
	if args[0] != b.adminKey {
		b.sendMessage(msg.Chat.ID, "Неверный админ-ключ")
		return
	}

	if len(args) < 2 {
		b.sendMessage(msg.Chat.ID, "Доступные админ-команды:\n"+
			"create_flight\n"+
			"delete_flight\n"+
			"flight_stats")
		return
	}

	switch args[1] {
	case "create_flight":
		b.handleCreateFlight(ctx, msg, args[2:])
	case "delete_flight":
		b.handleDeleteFlight(ctx, msg, args[2:])
	case "flight_stats":
		b.handleFlightStats(ctx, msg, args[2:])
	default:
		b.sendMessage(msg.Chat.ID, "Неизвестная админ-команда")
	}
}

// Создание рейса
func (b *AviaTicketsBot) handleCreateFlight(ctx context.Context, msg *tgbotapi.Message, args []string) {
	if len(args) < 7 {
		b.sendMessage(msg.Chat.ID, "Формат: /admin <ключ> create_flight "+
			"<номер> <откуда> <куда> <вылет(2006-01-02T15:04)> <прилет> <цена> <места>")
		return
	}

	flightNumber := args[0]
	departure := args[1]
	arrival := args[2]

	departureTime, err := time.Parse(time.RFC3339, args[3])
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Неверный формат даты вылета. Используйте RFC3339 (2006-01-02T15:04:00Z)")
		return
	}

	arrivalTime, err := time.Parse(time.RFC3339, args[4])
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Неверный формат даты прилета. Используйте RFC3339 (2006-01-02T15:04:00Z)")
		return
	}

	price, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Неверный формат цены")
		return
	}

	seats, err := strconv.Atoi(args[6])
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Неверный формат количества мест")
		return
	}

	flight := entity.Flight{
		FlightNumber:     flightNumber,
		DepartureAirport: departure,
		ArrivalAirport:   arrival,
		DepartureTime:    departureTime,
		ArrivalTime:      arrivalTime,
		Price:            price,
		AvailableSeats:   seats,
	}

	if err := b.flightUC.CreateFlight(ctx, flight); err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка создания рейса: "+err.Error())
		return
	}

	b.sendMessage(msg.Chat.ID, "Рейс успешно создан!")
}

// Удаление рейса
func (b *AviaTicketsBot) handleDeleteFlight(ctx context.Context, msg *tgbotapi.Message, args []string) {
	if len(args) < 1 {
		b.sendMessage(msg.Chat.ID, "Формат: /admin <ключ> delete_flight <номер_рейса>")
		return
	}

	flightNumber := args[0]

	// Находим рейс по номеру
	flight, err := b.flightUC.GetFlightByNumber(ctx, flightNumber)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Рейс не найден")
		return
	}

	if err := b.flightUC.DeleteFlight(ctx, flight.ID); err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка удаления рейса: "+err.Error())
		return
	}

	b.sendMessage(msg.Chat.ID, fmt.Sprintf("Рейс %s успешно удален", flightNumber))
}

// Статистика по рейсу
func (b *AviaTicketsBot) handleFlightStats(ctx context.Context, msg *tgbotapi.Message, args []string) {
	if len(args) < 1 {
		b.sendMessage(msg.Chat.ID, "Формат: /admin <ключ> flight_stats <номер_рейса>")
		return
	}

	flightNumber := args[0]

	// Находим рейс по номеру
	flight, err := b.flightUC.GetFlightByNumber(ctx, flightNumber)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Рейс не найден")
		return
	}

	stats, err := b.flightUC.GetFlightStats(ctx, flight.ID)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка получения статистики: "+err.Error())
		return
	}

	response := fmt.Sprintf(
		"📊 Статистика по рейсу %s\n\n"+
			"Всего мест: %d\n"+
			"Забронировано: %d\n"+
			"Свободно: %d\n"+
			"Процент загрузки: %.2f%%",
		flightNumber,
		stats.TotalSeats,
		stats.BookedSeats,
		stats.AvailableSeats,
		stats.LoadFactor*100,
	)

	b.sendMessage(msg.Chat.ID, response)
}
