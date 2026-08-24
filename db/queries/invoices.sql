-- name: ListInvoices :many
SELECT i.*, c.name AS client_name
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 ORDER BY i.due_on ASC;

-- name: ListInvoicesByStatus :many
SELECT i.*, c.name AS client_name
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 WHERE i.status = ?
 ORDER BY i.due_on ASC;

-- name: ListInvoicesByClient :many
SELECT i.*, c.name AS client_name
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 WHERE i.client_id = ?
 ORDER BY i.due_on ASC;

-- name: GetInvoice :one
SELECT i.*, c.name AS client_name
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 WHERE i.id = ?;

-- name: CreateInvoice :one
INSERT INTO invoices (client_id, number, amount_cents, currency, issued_on, due_on, notes, status, agent_state)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateInvoiceFields :one
UPDATE invoices
   SET number        = COALESCE(sqlc.narg('number'), number),
       amount_cents  = COALESCE(sqlc.narg('amount_cents'), amount_cents),
       due_on        = COALESCE(sqlc.narg('due_on'), due_on),
       notes         = COALESCE(sqlc.narg('notes'), notes),
       status        = COALESCE(sqlc.narg('status'), status),
       agent_state   = COALESCE(sqlc.narg('agent_state'), agent_state),
       current_stage = COALESCE(sqlc.narg('current_stage'), current_stage),
       next_action_on = sqlc.narg('next_action_on'),
       amount_paid_cents = COALESCE(sqlc.narg('amount_paid_cents'), amount_paid_cents)
 WHERE id = sqlc.arg('id')
RETURNING *;

-- name: CountInvoiceNumberForClient :one
SELECT COUNT(*) FROM invoices WHERE client_id = ? AND number = ? AND id != sqlc.narg('exclude_id');
