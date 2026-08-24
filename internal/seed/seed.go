// Package seed populates a fresh database with the demo story: two settled
// invoices (pays-on-time payer), one payer on-time streak broken, Meridian
// nine days overdue mid-chase, and a firm draft awaiting approval.
// Dates are relative to seed time so the demo never goes stale.
package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/AlfinRy/lunas/internal/db"
)

func Run(ctx context.Context, q *db.Queries) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	d := func(days int) string { return today.AddDate(0, 0, days).Format("2006-01-02") }
	ts := func(days int) string {
		return today.AddDate(0, 0, days).Format("2006-01-02") + "T09:00:00Z"
	}

	clients := []db.CreateClientParams{
		{Name: "Meridian Coffee Co.", Email: "billing@meridian.co", PaymentTermsDays: 14, RelationshipNote: "Long-term client — keep warm"},
		{Name: "Studio Kawi", Email: "accounts@studiokawi.id", PaymentTermsDays: 14},
		{Name: "Halo Print Co.", Email: "ap@haloprint.co", PaymentTermsDays: 30},
		{Name: "Arkon Digital", Email: "finance@arkondigital.io", PaymentTermsDays: 14, RelationshipNote: "Prefers invoice references in every email"},
		{Name: "Bayu Ventures", Email: "ops@bayuventures.com", PaymentTermsDays: 7},
		{Name: "Nara Interiors", Email: "hello@narainteriors.studio", PaymentTermsDays: 30},
	}
	ids := make(map[string]int64, len(clients))
	for _, c := range clients {
		row, err := q.CreateClient(ctx, c)
		if err != nil {
			return fmt.Errorf("seed client %s: %w", c.Name, err)
		}
		ids[c.Name] = row.ID
	}

	type inv struct {
		client   string
		number   string
		amount   int64
		issued   int
		due      int
		status   string
		agent    string
		stage    sql.NullString
		next     sql.NullString
		activity []db.InsertActivityParams
	}
	mk := func(client, number string, amount int64, issued, due int, status, agent string) inv {
		return inv{client: client, number: number, amount: amount, issued: issued, due: due, status: status, agent: agent}
	}

	invoices := []inv{
		mk("Halo Print Co.", "INV-0038", 225000, -70, -40, "paid", "stopped"),
		mk("Halo Print Co.", "INV-0039", 390000, -40, -10, "paid", "stopped"),
		mk("Meridian Coffee Co.", "INV-0041", 240000, -30, -9, "chasing", "awaiting_approval"),
		mk("Arkon Digital", "INV-0043", 520000, -25, -11, "chasing", "chasing"),
		mk("Studio Kawi", "INV-0040", 115000, -20, 2, "scheduled", "planning"),
		mk("Meridian Coffee Co.", "INV-0042", 86000, -12, 2, "scheduled", "planning"),
		mk("Bayu Ventures", "INV-0045", 98000, -3, 4, "scheduled", "planning"),
		mk("Nara Interiors", "INV-0044", 170000, -8, 22, "scheduled", "idle"),
	}

	// Narrative state for the two mid-chase invoices.
	meridian := &invoices[2]
	meridian.stage = sql.NullString{String: "stage2_firm", Valid: true}
	meridian.next = sql.NullString{String: d(0), Valid: true}
	arkon := &invoices[3]
	arkon.stage = sql.NullString{String: "stage2_firm", Valid: true}
	arkon.next = sql.NullString{String: d(3), Valid: true}

	invIDs := make(map[string]int64, len(invoices))
	for _, iv := range invoices {
		row, err := q.CreateInvoice(ctx, db.CreateInvoiceParams{
			ClientID:    ids[iv.client],
			Number:      iv.number,
			AmountCents: iv.amount,
			Currency:    "USD",
			IssuedOn:    d(iv.issued),
			DueOn:       d(iv.due),
			Status:      iv.status,
			AgentState:  iv.agent,
		})
		if err != nil {
			return fmt.Errorf("seed invoice %s: %w", iv.number, err)
		}
		invIDs[iv.number] = row.ID
		if _, err := q.UpdateInvoiceFields(ctx, db.UpdateInvoiceFieldsParams{
			ID:           row.ID,
			CurrentStage: iv.stage,
			NextActionOn: iv.next,
		}); err != nil {
			return err
		}
		if err := q.InsertActivity(ctx, db.InsertActivityParams{
			InvoiceID: sql.NullInt64{Int64: row.ID, Valid: true},
			Type:      "invoice_created",
			Message:   fmt.Sprintf("Invoice %s added — %s, $%.2f", iv.number, iv.client, float64(iv.amount)/100),
		}); err != nil {
			return err
		}
	}

	// Payments: Halo Print pays on time (twice) → payment score "pays on time".
	payments := []struct {
		invoice string
		amount  int64
		day     int
	}{
		{"INV-0038", 225000, -39},
		{"INV-0039", 390000, -12},
	}
	for _, p := range payments {
		id := invIDs[p.invoice]
		if _, err := q.CreatePayment(ctx, db.CreatePaymentParams{
			InvoiceID:   sql.NullInt64{Int64: id, Valid: true},
			AmountCents: p.amount,
			PaidOn:      d(p.day),
			Source:      "manual",
		}); err != nil {
			return err
		}
		if err := q.ApplyPayment(ctx, db.ApplyPaymentParams{AmountPaidCents: p.amount, Status: "paid", ID: id}); err != nil {
			return err
		}
	}

	// Chase history for Meridian (plan → reminder → polite → firm draft pending).
	meridianID := invIDs["INV-0041"]
	story := []db.InsertActivityParams{
		{InvoiceID: sql.NullInt64{Int64: meridianID, Valid: true}, Type: "plan_made", Message: "Scheduling a soft reminder on the due date — Meridian usually pays 9 days late."},
		{InvoiceID: sql.NullInt64{Int64: meridianID, Valid: true}, Type: "draft_sent", Message: "Sent the due-date reminder to billing@meridian.co."},
		{InvoiceID: sql.NullInt64{Int64: meridianID, Valid: true}, Type: "draft_sent", Message: "Sent a polite follow-up — invoice is 2 days overdue, no reply yet."},
		{InvoiceID: sql.NullInt64{Int64: meridianID, Valid: true}, Type: "draft_ready", Message: "Drafted a firm follow-up for approval — invoice is 9 days overdue."},
	}
	for _, a := range story {
		if err := q.InsertActivity(ctx, a); err != nil {
			return err
		}
	}
	if _, err := q.InsertDraft(ctx, db.InsertDraftParams{
		InvoiceID: meridianID, Stage: "stage2_firm", Status: "pending",
		Subject: "Second notice: invoice INV-0041 is 9 days overdue",
		Body: fmt.Sprintf("Hi Andini,\n\nFollowing up again on invoice INV-0041 — $2,400.00, now 9 days past due. We haven't received a reply to the previous note.\n\nPlease confirm a payment date this week, or settle directly: https://pay.lunas.app/inv/0041\n\nIf there's a dispute or question about the work, flag it now so we can resolve it quickly.\n\nRani Prameswari\nRani Studio\n\n—\nThis is an automated reminder sent by Lunas on behalf of Rani Prameswari."),
	}); err != nil {
		return err
	}

	// Outbox: the sent chases (mock mail client evidence).
	sent := []db.InsertOutboxParams{
		{InvoiceID: meridianID, InvoiceNumber: "INV-0041", ToName: "Meridian Coffee Co.", ToEmail: "billing@meridian.co",
			Subject: "Invoice INV-0041 — due " + d(-9), SentAt: ts(-9),
			Body: "A quick note that invoice INV-0041 for $2,400.00 is due today."},
		{InvoiceID: meridianID, InvoiceNumber: "INV-0041", ToName: "Meridian Coffee Co.", ToEmail: "billing@meridian.co",
			Subject: "Invoice INV-0041 — $2,400.00, now 2 days past due", SentAt: ts(-7),
			Body: "Marking this as unpaid on our side — invoice INV-0041 was due " + d(-9) + "."},
		{InvoiceID: invIDs["INV-0043"], InvoiceNumber: "INV-0043", ToName: "Arkon Digital", ToEmail: "finance@arkondigital.io",
			Subject: "Invoice INV-0043 — $5,200.00, now 3 days past due", SentAt: ts(-8),
			Body: "Marking this as unpaid on our side — invoice INV-0043 was due " + d(-11) + "."},
		{InvoiceID: invIDs["INV-0043"], InvoiceNumber: "INV-0043", ToName: "Arkon Digital", ToEmail: "finance@arkondigital.io",
			Subject: "Second notice: invoice INV-0043 is 7 days overdue", SentAt: ts(-4),
			Body: "Following up again on invoice INV-0043 — $5,200.00, now 7 days past due."},
	}
	for _, o := range sent {
		if err := q.InsertOutbox(ctx, o); err != nil {
			return err
		}
	}

	// Plan cards for the scheduled invoices.
	plans := []db.InsertActivityParams{
		{InvoiceID: sql.NullInt64{Int64: invIDs["INV-0040"], Valid: true}, Type: "plan_made", Message: "Scheduling a soft reminder on the due date — Studio Kawi pays on time."},
		{InvoiceID: sql.NullInt64{Int64: invIDs["INV-0042"], Valid: true}, Type: "plan_made", Message: "Reminder one day before due — Meridian usually pays 9 days late."},
	}
	for _, p := range plans {
		if err := q.InsertActivity(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
