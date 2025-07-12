SELECT id, flight_number, departure_airport, arrival_airport,
       departure_time, arrival_time, price, available_seats, created_at
FROM flights
WHERE departure_airport ILIKE @origin AND arrival_airport ILIKE @destination
ORDER BY departure_time