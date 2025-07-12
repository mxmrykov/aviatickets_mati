SELECT EXISTS(
    SELECT 1 FROM bookings
    WHERE flight_id = @flight_id AND seat_number = @seat_number
)