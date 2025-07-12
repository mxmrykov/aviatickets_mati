CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       telegram_id BIGINT UNIQUE NOT NULL,
                       username VARCHAR(100),
                       first_name VARCHAR(100),
                       last_name VARCHAR(100),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE flights (
                         id SERIAL PRIMARY KEY,
                         flight_number VARCHAR(20) UNIQUE NOT NULL,
                         departure_airport VARCHAR(10) NOT NULL,
                         arrival_airport VARCHAR(10) NOT NULL,
                         departure_time TIMESTAMP NOT NULL,
                         arrival_time TIMESTAMP NOT NULL,
                         price DECIMAL(10, 2) NOT NULL,
                         available_seats INT NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bookings (
                          id SERIAL PRIMARY KEY,
                          user_id INT REFERENCES users(id),
                          flight_id INT REFERENCES flights(id),
                          passenger_name VARCHAR(200) NOT NULL,
                          seat_number VARCHAR(10),
                          booking_status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
                          booking_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          UNIQUE (flight_id, seat_number)
);

CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_flight_id ON bookings(flight_id);