package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"aviatickets_mati/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *AviaTicketsBot) handleStart(chatID int64, user *entity.User) {
	msg := tgbotapi.NewMessage(chatID,
		"Добро пожаловать в систему бронирования авиабилетов!\n\n"+
			"Доступные команды:\n"+
			"/search - Поиск рейсов\n"+
			"/book - Забронировать билет\n"+
			"/mybookings - Мои бронирования\n"+
			"/cancel - Отменить бронирование")

	b.bot.Send(msg)
}

func (b *AviaTicketsBot) handleSearchCommand(ctx context.Context, msg *tgbotapi.Message) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 2 {
		b.sendMessage(msg.Chat.ID, "Использование: /search <откуда> <куда>")
		return
	}

	origin := args[0]
	destination := args[1]

	flights, err := b.flightUC.SearchFlights(ctx, origin, destination)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при поиске рейсов")
		return
	}

	if len(flights) == 0 {
		b.sendMessage(msg.Chat.ID, "Рейсы не найдены")
		return
	}

	var response strings.Builder
	response.WriteString("Найденные рейсы:\n\n")
	for _, flight := range flights {
		response.WriteString(
			fmt.Sprintf("✈️ Рейс %s\n🛫 %s -> %s\n⏱ %s - %s\n💵 Цена: %.2f ₽\n🪑 Свободных мест: %d\n\n",
				flight.FlightNumber,
				flight.DepartureAirport,
				flight.ArrivalAirport,
				flight.DepartureTime.Format("02.01.2006 15:04"),
				flight.ArrivalTime.Format("02.01.2006 15:04"),
				flight.Price,
				flight.AvailableSeats,
			))
	}

	b.sendMessage(msg.Chat.ID, response.String())
}

func (b *AviaTicketsBot) handleBookCommand(ctx context.Context, msg *tgbotapi.Message, userID int) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 3 {
		b.sendMessage(msg.Chat.ID, "Использование: /book <номер_рейса> <имя_пассажира> <номер_места>")
		return
	}

	flightNumber := args[0]
	passengerName := strings.Join(args[1:len(args)-1], " ")
	seatNumber := args[len(args)-1]

	// Находим рейс по номеру
	flight, err := b.flightUC.GetFlightByNumber(ctx, flightNumber)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Рейс не найден")
		return
	}

	// Создаем бронирование
	booking := entity.Booking{
		UserID:        userID,
		FlightID:      flight.ID,
		PassengerName: passengerName,
		SeatNumber:    seatNumber,
		BookingStatus: "confirmed",
	}

	if err := b.bookingUC.CreateBooking(ctx, booking); err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при бронировании: "+err.Error())
		return
	}

	response := fmt.Sprintf(
		"✅ Бронирование подтверждено!\n\n"+
			"Пассажир: %s\n"+
			"Рейс: %s\n"+
			"Место: %s\n"+
			"Дата вылета: %s",
		passengerName,
		flight.FlightNumber,
		seatNumber,
		flight.DepartureTime.Format("02.01.2006 15:04"),
	)

	b.sendMessage(msg.Chat.ID, response)
}

func (b *AviaTicketsBot) handleMyBookingsCommand(ctx context.Context, msg *tgbotapi.Message, userID int) {
	bookings, err := b.bookingUC.GetUserBookings(ctx, userID)

	if err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при получении бронирований")
		return
	}

	if len(bookings) == 0 {
		b.sendMessage(msg.Chat.ID, "У вас нет активных бронирований")
		return
	}

	var response strings.Builder
	response.WriteString("Ваши бронирования:\n\n")

	for _, booking := range bookings {
		response.WriteString(
			fmt.Sprintf("📌 Бронирование #%d\n\n"+
				"Маршрут: %s - %s\n\n"+
				"Пассажир: %s\n"+
				"Рейс: %s\n"+
				"Место: %s\n\n"+
				"Дата и время вылета: %s\n"+
				"Дата и время посадки: %s\n\n"+
				"Статус: %s\n",
				booking.ID,
				booking.Flight.DepartureAirport,
				booking.Flight.ArrivalAirport,
				booking.PassengerName,
				booking.Flight.FlightNumber,
				booking.SeatNumber,
				booking.Flight.DepartureTime,
				booking.Flight.ArrivalTime,
				booking.BookingStatus,
			))
	}

	b.sendMessage(msg.Chat.ID, response.String())
}

func (b *AviaTicketsBot) handleCancelCommand(ctx context.Context, msg *tgbotapi.Message, userID int) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 1 {
		b.sendMessage(msg.Chat.ID, "Использование: /cancel <ID бронирования>")
		return
	}

	bookingID, err := strconv.Atoi(args[0])
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Неверный ID бронирования")
		return
	}

	// Проверяем, что бронирование принадлежит пользователю
	booking, err := b.bookingUC.GetBookingByID(ctx, bookingID)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Бронирование не найдено")
		return
	}

	if booking.UserID != userID {
		b.sendMessage(msg.Chat.ID, "Это не ваше бронирование")
		return
	}

	if err := b.bookingUC.CancelBooking(ctx, bookingID); err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при отмене бронирования")
		return
	}

	b.sendMessage(msg.Chat.ID, "Бронирование успешно отменено")
}
