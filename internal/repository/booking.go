package repository

import (
	"context"
	_ "embed"
	"errors"

	"aviatickets_mati/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/booking/seatCheck.sql
	seatCheckQuery string

	//go:embed sql/booking/booking.sql
	bookingQuery string

	//go:embed sql/booking/updateSeats.sql
	updateQuery string

	//go:embed sql/booking/findByUserID.sql
	findByUserID string

	//go:embed sql/booking/findByFlightID.sql
	findByFlightID string

	//go:embed sql/booking/getFlightIDBeforeCancel.sql
	getFlightIDBeforeCancel string

	//go:embed sql/booking/updateBooking.sql
	updateBooking string

	//go:embed sql/booking/updateFlight.sql
	updateFlight string
)

var ErrSeatAlreadyTaken = errors.New("seat is already taken")

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking entity.Booking) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Проверяем доступность места
	var seatTaken bool
	seatArgs := pgx.NamedArgs{
		"flight_id":   booking.FlightID,
		"seat_number": booking.SeatNumber,
	}

	if err := tx.QueryRow(ctx, seatCheckQuery, seatArgs).Scan(&seatTaken); err != nil {
		return err
	}

	if seatTaken {
		return ErrSeatAlreadyTaken
	}

	bookingArgs := pgx.NamedArgs{
		"user_id":        booking.UserID,
		"flight_id":      booking.FlightID,
		"passenger_name": booking.PassengerName,
		"seat_number":    booking.SeatNumber,
		"booking_status": booking.BookingStatus,
	}

	// Создаем бронирование
	if err := tx.QueryRow(ctx, bookingQuery, bookingArgs).Scan(&booking.ID, &booking.BookingDate); err != nil {
		return err
	}

	updateArgs := pgx.NamedArgs{
		"flight_id": booking.FlightID,
	}

	if _, err := tx.Exec(ctx, updateQuery, updateArgs); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *BookingRepository) FindByUserID(ctx context.Context, userID int) ([]entity.Booking, error) {
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := r.db.Query(ctx, findByUserID, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]entity.Booking, 0)

	for rows.Next() {
		var (
			booking entity.Booking
			flight  entity.Flight
		)
		if err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.FlightID,
			&booking.PassengerName,
			&booking.SeatNumber,
			&booking.BookingStatus,
			&booking.BookingDate,
			&flight.FlightNumber,
			&flight.DepartureAirport,
			&flight.ArrivalAirport,
			&flight.DepartureTime,
			&flight.ArrivalTime,
			&flight.Price,
		); err != nil {
			return nil, err
		}

		booking.Flight = flight
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, id int) (*entity.Booking, error) {
	args := pgx.NamedArgs{
		"id": id,
	}

	var booking entity.Booking
	var flight entity.Flight
	err := r.db.QueryRow(ctx, findByFlightID, args).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.FlightID,
		&booking.PassengerName,
		&booking.SeatNumber,
		&booking.BookingStatus,
		&booking.BookingDate,
		&flight.FlightNumber,
		&flight.DepartureAirport,
		&flight.ArrivalAirport,
		&flight.DepartureTime,
		&flight.ArrivalTime,
		&flight.Price,
	)

	if err != nil {
		return nil, err
	}

	booking.FlightID = flight.ID
	return &booking, nil
}

func (r *BookingRepository) Cancel(ctx context.Context, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var flightID int
	getArgs := pgx.NamedArgs{"id": id}

	// Получаем информацию о бронировании перед отменой
	if err = tx.QueryRow(ctx, getFlightIDBeforeCancel, getArgs).Scan(&flightID); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, updateBooking, getArgs); err != nil {
		return err
	}

	flightArgs := pgx.NamedArgs{"flight_id": flightID}
	if _, err := tx.Exec(ctx, updateFlight, flightArgs); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
