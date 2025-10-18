-- name: CreateAccount :one
INSERT INTO accounts (name, role, email, hashed_password, phone, student_id)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE student_id = $1 AND is_delete != true;

-- name: ListAccounts :many
SELECT id, name, role, email, phone, student_id FROM accounts WHERE is_delete != true LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $2, phone = $3
WHERE id = $1
AND is_delete != true
RETURNING *;

-- name: DeleteAccount :exec
UPDATE accounts
SET is_delete = true
WHERE student_id = $1
AND is_delete != true;