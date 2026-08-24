-- name: ListClients :many
SELECT * FROM clients ORDER BY name;

-- name: GetClient :one
SELECT * FROM clients WHERE id = ?;

-- name: CreateClient :one
INSERT INTO clients (name, email, payment_terms_days, relationship_note)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateClient :one
UPDATE clients
   SET name = COALESCE(sqlc.narg('name'), name),
       email = COALESCE(sqlc.narg('email'), email),
       payment_terms_days = COALESCE(sqlc.narg('payment_terms_days'), payment_terms_days),
       relationship_note = COALESCE(sqlc.narg('relationship_note'), relationship_note)
 WHERE id = sqlc.arg('id')
RETURNING *;
