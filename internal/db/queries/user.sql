-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, fname, lname, email, mobile_number, password_hash, is_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUserFirstName :one
UPDATE users
SET 
    fname = $1
WHERE id = $2
RETURNING id;

-- name: UpdateUserLastName :one
UPDATE users
SET 
    lname = $1
WHERE id = $2
RETURNING id;

-- name: UpdateUserMobileNumber :one
UPDATE users
SET 
    mobile_number = $1
WHERE id = $2
RETURNING id;

-- name: UpdateUserRole :one
UPDATE users
SET 
    is_admin = $1
WHERE id = $2
RETURNING id;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id;

-- name: GetUsers :many
SELECT *
FROM users
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;