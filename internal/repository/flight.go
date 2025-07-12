package repository

import (
	"context"
	_ "embed"

	"aviatickets_mati/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/flight/findAll.sql
	findAll string

	//go:embed sql/flight/findByID.sql
	findByID string

	//go:embed sql/flight/search.sql
	search string

	//go:embed sql/flight/updateSeats.sql
	updateSeats string

	//go:embed sql/admin/getStats.sql
	getStats string

	//go:embed sql/flight/findByFlightNum.sql
	findByFlightNum string
)

type FlightRepository struct {
	db *pgxpool.Pool
}

func NewFlightRepository(db *pgxpool.Pool) *FlightRepository {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) FindAll(ctx context.Context) ([]entity.Flight, error) {
	rows, err := r.db.Query(ctx, findAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []entity.Flight
	for rows.Next() {
		var flight entity.Flight
		if err := rows.Scan(
			&flight.ID,
			&flight.FlightNumber,
			&flight.DepartureAirport,
			&flight.ArrivalAirport,
			&flight.DepartureTime,
			&flight.ArrivalTime,
			&flight.Price,
			&flight.AvailableSeats,
			&flight.CreatedAt,
		); err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

func (r *FlightRepository) FindById(ctx context.Context, id int) (*entity.Flight, error) {
	args := pgx.NamedArgs{
		"id": id,
	}

	var flight entity.Flight
	err := r.db.QueryRow(ctx, findByID, args).Scan(
		&flight.ID,
		&flight.FlightNumber,
		&flight.DepartureAirport,
		&flight.ArrivalAirport,
		&flight.DepartureTime,
		&flight.ArrivalTime,
		&flight.Price,
		&flight.AvailableSeats,
		&flight.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &flight, nil
}

func (r *FlightRepository) Search(ctx context.Context, origin, destination string) ([]entity.Flight, error) {
	args := pgx.NamedArgs{
		"origin":      "%" + origin + "%",
		"destination": "%" + destination + "%",
	}

	rows, err := r.db.Query(ctx, search, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []entity.Flight
	for rows.Next() {
		var flight entity.Flight
		if err := rows.Scan(
			&flight.ID,
			&flight.FlightNumber,
			&flight.DepartureAirport,
			&flight.ArrivalAirport,
			&flight.DepartureTime,
			&flight.ArrivalTime,
			&flight.Price,
			&flight.AvailableSeats,
			&flight.CreatedAt,
		); err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

func (r *FlightRepository) UpdateSeats(ctx context.Context, flightID int, seats int) error {
	args := pgx.NamedArgs{
		"flight_id": flightID,
		"seats":     seats,
	}

	_, err := r.db.Exec(ctx, updateSeats, args)
	return err
}

func (r *FlightRepository) Create(ctx context.Context, flight entity.Flight) error {
	query := `
        INSERT INTO flights (
            flight_number, departure_airport, arrival_airport,
            departure_time, arrival_time, price, available_seats
        ) VALUES (
            @flight_number, @departure_airport, @arrival_airport,
            @departure_time, @arrival_time, @price, @available_seats
        )
    `

	args := pgx.NamedArgs{
		"flight_number":     flight.FlightNumber,
		"departure_airport": flight.DepartureAirport,
		"arrival_airport":   flight.ArrivalAirport,
		"departure_time":    flight.DepartureTime,
		"arrival_time":      flight.ArrivalTime,
		"price":             flight.Price,
		"available_seats":   flight.AvailableSeats,
	}

	_, err := r.db.Exec(ctx, query, args)
	return err
}

func (r *FlightRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM flights WHERE id = @id`
	args := pgx.NamedArgs{"id": id}
	_, err := r.db.Exec(ctx, query, args)
	return err
}

func (r *FlightRepository) GetStats(ctx context.Context, flightID int) (*entity.FlightStats, error) {
	args := pgx.NamedArgs{"flight_id": flightID}

	var stats entity.FlightStats
	err := r.db.QueryRow(ctx, getStats, args).Scan(
		&stats.TotalSeats,
		&stats.BookedSeats,
		&stats.AvailableSeats,
		&stats.LoadFactor,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *FlightRepository) FindByFlightNumber(ctx context.Context, flightNumber string) (*entity.Flight, error) {
	args := pgx.NamedArgs{
		"flight_number": flightNumber,
	}

	var flight entity.Flight
	err := r.db.QueryRow(ctx, findByFlightNum, args).Scan(
		&flight.ID,
		&flight.FlightNumber,
		&flight.DepartureAirport,
		&flight.ArrivalAirport,
		&flight.DepartureTime,
		&flight.ArrivalTime,
		&flight.Price,
		&flight.AvailableSeats,
		&flight.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &flight, nil
}
