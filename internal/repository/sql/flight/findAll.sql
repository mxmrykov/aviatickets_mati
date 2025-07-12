SELECT id, flight_number, departure_airport, arrival_airport,
       departure_time, arrival_time, price, available_seats, created_at
FROM flights
ORDER BY departure_time