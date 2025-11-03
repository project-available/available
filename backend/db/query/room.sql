-- name: CreateRoom :one
INSERT INTO rooms (name, location, image)
VALUES ($1,$2,$3)
RETURNING *;

-- name: ListRooms :many
SELECT * FROM rooms
WHERE is_delete != true
ORDER BY id
LIMIT $1 OFFSET $2;

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