-- name: CreateAccount :one
INSERT INTO accounts (name, role, email, password, phone, student_id)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1 AND is_delete != true;

-- name: ListAccounts :many
SELECT * FROM accounts WHERE is_delete != true LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET name = $2, role = $3, email = $4, password = $5, phone = $6, student_id = $7
WHERE id = $1
AND is_delete != true
RETURNING *;

-- name: DeleteAccount :one
UPDATE accounts
SET is_delete = true
WHERE id = $1
AND is_delete != true
RETURNING *;