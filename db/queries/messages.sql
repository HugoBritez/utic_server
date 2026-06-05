-- name: CreateMessage :one
INSERT INTO messages (
  id, sender_name, sender_phone_number, message_text, created_at
) VALUES (
  ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages;

-- name: GetNewPhoneNumbers :many
SELECT sender_phone_number FROM messages;
