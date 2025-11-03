-- name: CreateCustomField :one
INSERT INTO "custom_fields" ("key")
VALUES ($1)
RETURNING *;

-- name: UpdateCustomFieldShown :exec
UPDATE "custom_fields"
SET "key" = $2, shown = $3
WHERE id = $1;

-- name: ListCustomFields :many
SELECT * FROM "custom_fields"
ORDER BY "key";

-- name: ListRoomCustomFieldValues :many
SELECT * FROM "custom_fields_value"
WHERE room_id = $1;

-- name: ListRoomCustomFieldValuesBatch :many
SELECT * FROM "custom_fields_value"
WHERE room_id = ANY($1::bigint[]);

-- name: UpsertRoomCustomFieldValue :exec
INSERT INTO "custom_fields_value" (room_id, customfield_id, value)
VALUES ($1, $2, $3)
ON CONFLICT (room_id, customfield_id) DO UPDATE SET value = EXCLUDED.value;

