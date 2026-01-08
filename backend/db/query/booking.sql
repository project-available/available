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

-- name: GetBookingsOnDate :many
SELECT * FROM bookings
WHERE room_id = $1 AND status = 'confirmed' AND start >= $2 AND start < $3;

-- name: UpdateBooking :one
UPDATE bookings
SET status = $2
WHERE id = $1
RETURNING *;

-- name: CheckBookingOverlap :one
SELECT COUNT(*) as overlap_count
FROM bookings
WHERE room_id = $1
  AND status IN ('pending', 'confirmed')
  AND NOT (
    "end" <= $2           -- existing booking ends before our start
    OR "start" >= $3      -- existing booking starts after our end
  );
