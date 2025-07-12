package entity

import "time"

type Booking struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	FlightID      int       `json:"flight_id"`
	PassengerName string    `json:"passenger_name"`
	SeatNumber    string    `json:"seat_number"`
	BookingStatus string    `json:"booking_status"`
	BookingDate   time.Time `json:"booking_date"`
	Flight        Flight    `json:"flight"`
}
