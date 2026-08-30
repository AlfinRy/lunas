-- name: ListOpenForMatch :many
SELECT i.id AS invoice_id, i.number AS invoice_no, i.client_id, c.name AS client_name,
       i.amount_cents, i.amount_paid_cents, i.due_on
  FROM invoices i
  JOIN clients c ON c.id = i.client_id
 WHERE i.status IN ('scheduled', 'chasing')
 ORDER BY i.due_on ASC;

-- name: ApplyPartialPayment :exec
UPDATE invoices SET amount_paid_cents = amount_paid_cents + ?, agent_state = 'chasing'
 WHERE id = ?;

-- name: CancelPendingDrafts :exec
UPDATE email_drafts SET status = 'skipped' WHERE invoice_id = ? AND status = 'pending';
