-- name: GetSettings :one
SELECT * FROM settings WHERE id = 1;

-- name: InitSettings :exec
INSERT OR IGNORE INTO settings (id) VALUES (1);

-- name: UpdateSettings :one
UPDATE settings
   SET sender_name        = COALESCE(sqlc.narg('sender_name'), sender_name),
       sender_email       = COALESCE(sqlc.narg('sender_email'), sender_email),
       default_terms_days = COALESCE(sqlc.narg('default_terms_days'), default_terms_days),
       global_mode        = COALESCE(sqlc.narg('global_mode'), global_mode),
       sim_now            = sqlc.narg('sim_now')
 WHERE id = 1
RETURNING *;

-- name: RecoveredTotal :one
SELECT COALESCE(SUM(p.amount_cents), 0) AS total FROM payments p
  JOIN invoices i ON i.id = p.invoice_id
 WHERE i.status = 'paid';

-- name: ListActivityForInvoice :many
SELECT * FROM activities WHERE invoice_id = ? ORDER BY created_at DESC, id DESC;

-- name: InsertActivity :exec
INSERT INTO activities (invoice_id, type, message) VALUES (?, ?, ?);

-- name: CountClients :one
SELECT COUNT(*) FROM clients;

-- name: CountPendingDrafts :one
SELECT COUNT(*) FROM email_drafts WHERE status = 'pending';


-- name: ResetDemoData :exec
DELETE FROM outbox_emails; DELETE FROM email_drafts; DELETE FROM payments; DELETE FROM activities; DELETE FROM invoices; DELETE FROM clients;
DELETE FROM sqlite_sequence WHERE name IN ('outbox_emails','email_drafts','payments','activities','invoices','clients');

-- name: CreatePayment :one
INSERT INTO payments (invoice_id, amount_cents, paid_on, source, confidence, raw_text)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ApplyPayment :exec
UPDATE invoices
   SET amount_paid_cents = amount_paid_cents + ?,
       status = ?,
       agent_state = 'stopped',
       current_stage = NULL,
       next_action_on = NULL
 WHERE id = ?;

-- name: InsertDraft :one
INSERT INTO email_drafts (invoice_id, stage, subject, body, status)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: InsertOutbox :exec
INSERT INTO outbox_emails (invoice_id, invoice_number, to_name, to_email, subject, body, sent_at)
VALUES (?, ?, ?, ?, ?, ?, ?);
