-- name: ListOpenInvoicesWithPayer :many
SELECT i.*, c.name AS client_name, c.email AS client_email, c.relationship_note
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 WHERE i.status IN ('scheduled', 'chasing')
 ORDER BY i.due_on ASC;

-- name: ListPendingDraftsWithInvoice :many
SELECT d.*, i.number AS invoice_number, i.amount_cents, i.currency, i.due_on,
       c.name AS client_name, c.email AS client_email
  FROM email_drafts d
  JOIN invoices i ON i.id = d.invoice_id
  JOIN clients c ON c.id = i.client_id
 WHERE d.status = 'pending'
 ORDER BY d.created_at DESC, d.id DESC;

-- name: GetDraftWithInvoice :one
SELECT d.*, i.number AS invoice_number, i.amount_cents, i.currency, i.due_on,
       c.name AS client_name, c.email AS client_email
  FROM email_drafts d
  JOIN invoices i ON i.id = d.invoice_id
  JOIN clients c ON c.id = i.client_id
 WHERE d.id = ?;

-- name: MarkDraftSent :exec
UPDATE email_drafts SET status = 'sent', sent_at = ? WHERE id = ?;

-- name: MarkDraftStatus :exec
UPDATE email_drafts SET status = ? WHERE id = ?;

-- name: SetInvoiceChase :exec
UPDATE invoices
   SET current_stage = ?, next_action_on = ?, agent_state = ?
 WHERE id = ?;

-- name: SetInvoiceAgentState :exec
UPDATE invoices SET agent_state = ? WHERE id = ?;

-- name: CountSentChases :one
SELECT COUNT(*) FROM outbox_emails WHERE invoice_id = ?;

-- name: ListOutbox :many
SELECT * FROM outbox_emails ORDER BY sent_at DESC, id DESC;

-- name: HasPendingDraftForStage :one
SELECT COUNT(*) FROM email_drafts WHERE invoice_id = ? AND stage = ? AND status = 'pending';

-- name: PaymentStats :many
SELECT p.invoice_id, p.amount_cents, p.paid_on, i.client_id, i.issued_on, i.due_on
  FROM payments p JOIN invoices i ON i.id = p.invoice_id
 WHERE i.status = 'paid';
