-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked)
VALUES ($1, $2, $3, false)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token_hash = $1 AND revoked = false;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked = true WHERE user_id = $1;