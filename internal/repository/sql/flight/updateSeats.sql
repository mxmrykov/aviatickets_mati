UPDATE flights
SET available_seats = available_seats - @seats
WHERE id = @flight_id AND available_seats >= @seats