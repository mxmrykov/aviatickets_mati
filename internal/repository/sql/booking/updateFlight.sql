UPDATE flights
SET available_seats = available_seats + 1
WHERE id = @flight_id