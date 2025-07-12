INSERT INTO bookings (user_id, flight_id, passenger_name, seat_number, booking_status)
VALUES (@user_id, @flight_id, @passenger_name, @seat_number, @booking_status)
RETURNING id, booking_date