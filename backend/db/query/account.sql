-- name: create_account :one
INSERT INTO accounts (name, role, email, password, phone, student_id)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: get_account :one
SELECT * FROM accounts WHERE id = $1 AND is_delete = false;

-- name: list_accounts :many
SELECT * FROM accounts WHERE is_delete = false LIMIT $1 OFFSET $2;

-- name: update_account :one
UPDATE accounts
SET name = $2, role = $3, email = $4, password = $5, phone = $6
WHERE id = $1
AND is_delete = false
RETURNING *;

-- name: delete_account :one
UPDATE accounts
SET is_delete = true
WHERE id = $1
AND is_delete = false
RETURNING *;