SELECT b.id, b.user_id, b.flight_id, b.passenger_name,
       b.seat_number, b.booking_status, b.booking_date,
       f.flight_number, f.departure_airport, f.arrival_airport,
       f.departure_time, f.arrival_time, f.price
FROM bookings b
         JOIN flights f ON b.flight_id = f.id
WHERE b.id = @id