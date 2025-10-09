-- name: CreateBooking :one
INSERT INTO bookings (account_id, room_id, start, "end", phone_booking)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
-- name: GetBookingOfAccount :many
SELECT * FROM bookings WHERE account_id = $1;
-- name: ListBookings :many
SELECT * FROM bookings LIMIT $1 OFFSET $2;
-- name: UpdateBooking :one
UPDATE bookings
SET status = $2
WHERE id = $1
RETURNING *;
