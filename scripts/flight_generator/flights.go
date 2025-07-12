package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"aviatickets_mati/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightGenerator struct {
	db *pgxpool.Pool
}

func NewFlightGenerator(db *pgxpool.Pool) *FlightGenerator {
	return &FlightGenerator{db: db}
}

func (g *FlightGenerator) GenerateFlights(count int) error {
	ctx := context.Background()

	for i := 0; i < count; i++ {
		flight := g.generateRandomFlight()
		if err := g.insertFlight(ctx, flight); err != nil {
			return fmt.Errorf("ошибка вставки рейса: %w", err)
		}
	}

	return nil
}

func (g *FlightGenerator) generateRandomFlight() entity.Flight {
	rand.Seed(time.Now().UnixNano())

	// Выбираем случайные аэропорты (исключаем полеты из одного аэропорта в тот же)
	var departure, arrival string
	for {
		departure = cisoAirports[rand.Intn(len(cisoAirports))]
		arrival = cisoAirports[rand.Intn(len(cisoAirports))]
		if departure != arrival {
			break
		}
	}

	// Генерируем даты после июля 2025 года
	baseDate := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)
	daysOffset := rand.Intn(365) // В пределах года
	hoursOffset := rand.Intn(24)
	minutesOffset := rand.Intn(12) * 5 // Кратно 5 минутам

	departureTime := baseDate.AddDate(0, 0, daysOffset).
		Add(time.Duration(hoursOffset) * time.Hour).
		Add(time.Duration(minutesOffset) * time.Minute)

	// Длительность полета от 1 до 8 часов
	duration := time.Duration(1+rand.Intn(8)) * time.Hour
	arrivalTime := departureTime.Add(duration)

	// Генерируем номер рейса (пример: SU1234)
	airlineCode := []string{"SU", "U6", "TK", "KC", "HY", "DV", "5N", "6K", "J2", "4G"}[rand.Intn(10)]
	flightNumber := fmt.Sprintf("%s%d", airlineCode, 1000+rand.Intn(9000))

	// Цена от 5000 до 50000 руб
	price := 5000 + rand.Float64()*45000

	// Количество мест от 50 до 300
	seats := 50 + rand.Intn(251)

	return entity.Flight{
		FlightNumber:     flightNumber,
		DepartureAirport: departure,
		ArrivalAirport:   arrival,
		DepartureTime:    departureTime,
		ArrivalTime:      arrivalTime,
		Price:            price,
		AvailableSeats:   seats,
	}
}

func (g *FlightGenerator) insertFlight(ctx context.Context, flight entity.Flight) error {
	query := `
        INSERT INTO flights (
            flight_number, departure_airport, arrival_airport,
            departure_time, arrival_time, price, available_seats
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )
    `

	_, err := g.db.Exec(ctx, query,
		flight.FlightNumber,
		flight.DepartureAirport,
		flight.ArrivalAirport,
		flight.DepartureTime,
		flight.ArrivalTime,
		flight.Price,
		flight.AvailableSeats,
	)

	if err != nil {
		return err
	}

	log.Printf("Добавлен рейс: %s %s-%s %s %.2f руб",
		flight.FlightNumber,
		flight.DepartureAirport,
		flight.ArrivalAirport,
		flight.DepartureTime.Format("02.01.2006 15:04"),
		flight.Price,
	)

	return nil
}
