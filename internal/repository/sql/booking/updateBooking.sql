UPDATE bookings
SET booking_status = 'cancelled'
WHERE id = @id