package entity

import "time"

type Flight struct {
	ID               int       `json:"id"`
	FlightNumber     string    `json:"flight_number"`
	DepartureAirport string    `json:"departure_airport"`
	ArrivalAirport   string    `json:"arrival_airport"`
	DepartureTime    time.Time `json:"departure_time"`
	ArrivalTime      time.Time `json:"arrival_time"`
	Price            float64   `json:"price"`
	AvailableSeats   int       `json:"available_seats"`
	CreatedAt        time.Time `json:"created_at"`
}

type FlightStats struct {
	TotalSeats     int
	BookedSeats    int
	AvailableSeats int
	LoadFactor     float64
}
