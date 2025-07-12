SELECT
    f.available_seats + COUNT(b.id) as total_seats,
    COUNT(b.id) as booked_seats,
    f.available_seats,
    COUNT(b.id)::float / (f.available_seats + COUNT(b.id)) as load_factor
FROM flights f
         LEFT JOIN bookings b ON f.id = b.flight_id AND b.booking_status = 'confirmed'
WHERE f.id = @flight_id
GROUP BY f.id, f.available_seats