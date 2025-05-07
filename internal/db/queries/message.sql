-- name: InsertMessage :one
INSERT INTO messages (sender_id, receiver_id, content, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: MarkAsDelivered :exec
UPDATE messages
SET delivered = TRUE
WHERE id = $1;

-- name: GetUndeliveredMessages :many
SELECT * FROM messages
WHERE receiver_id = $1 AND delivered = FALSE;
