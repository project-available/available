-- name: CreateRoom :one
INSERT INTO rooms (name, location, image)
VALUES ($1,$2,$3)
RETURNING *;

-- name: ListRooms :many
SELECT * FROM rooms
WHERE is_delete != true
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListRoomsWithStatus :many
SELECT
    r.id,
    r.location,
    r.name,
    r.image,
    r.is_delete,
    COALESCE(BOOL_OR(
        $1 >= b.start
        AND $1 < b.end
    ), FALSE) AS is_occupied_now,
    MAX(
        CASE
            WHEN $1 >= b.start AND $1 < b.end AND b.status = 'confirmed'
            THEN b.end
        END
    ) AS occupied_until
FROM rooms r
LEFT JOIN bookings b
    ON b.room_id = r.id
    AND $1 >= b.start
    AND $1 < b.end
    AND b.status = 'confirmed'
GROUP BY r.id
ORDER BY r.id
LIMIT $2 OFFSET $3;

-- name: UpdateRoom :one
UPDATE rooms
SET name = $2, location = $3, image = $4
WHERE id = $1
AND is_delete != true
RETURNING *;

-- name: GetRoom :one
SELECT * FROM rooms
WHERE id = $1
AND is_delete != true;

-- name: DeleteRoom :exec
UPDATE rooms
SET is_delete = true
WHERE id = $1
AND is_delete != true
RETURNING *;