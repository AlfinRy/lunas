package agent

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	db "github.com/AlfinRy/lunas/internal/db"
)

// Engine orchestrates the agent loop over the store. Decisions come from the
// pure policy; words come from the provider; everything lands in the store
// with an auditable activity trail.
type Engine struct {
	Q        *db.Queries
	Provider DraftProvider
}

func nowISO() string { return time.Now().UTC().Format(time.RFC3339) }

// reliability buckets per client from settled history (same thresholds as the
// API's payment score — pinned by internal/api tests).
func reliabilityBuckets(ctx context.Context, q *db.Queries) (map[int64]string, error) {
	rows, err := q.PaymentStats(ctx)
	if err != nil {
		return nil, err
	}
	type acc struct {
		sum   float64
		count int
	}
	accs := map[int64]*acc{}
	for _, r := range rows {
		a := accs[r.ClientID]
		if a == nil {
			a = &acc{}
			accs[r.ClientID] = a
		}
		a.sum += float64(days(r.DueOn, r.PaidOn))
		a.count++
	}
	out := map[int64]string{}
	for id, a := range accs {
		if a.count < 2 {
			continue // need two settled invoices before we judge (PRD F8)
		}
		avg := a.sum / float64(a.count)
		switch {
		case avg <= 1:
			out[id] = ReliabilityOnTime
		case avg <= 15:
			out[id] = ReliabilityUsually
		default:
			out[id] = ReliabilityChronic
		}
	}
	return out, nil
}

func isWarm(note string) bool {
	return strings.Contains(strings.ToLower(note), "warm")
}

// Run executes one agent pass for every open invoice, as of `today`.
// Returns a human summary of what it did (for the activity log / response).
func (e *Engine) Run(ctx context.Context, today, mode string) (string, error) {
	rows, err := e.Q.ListOpenInvoicesWithPayer(ctx)
	if err != nil {
		return "", err
	}
	buckets, err := reliabilityBuckets(ctx, e.Q)
	if err != nil {
		return "", err
	}

	drafted, sent, planned := 0, 0, 0
	for _, r := range rows {
		in := Input{
			Today:        today,
			DueOn:        r.DueOn,
			Status:       r.Status,
			LastSent:     Stage(""),
			Reliability:  buckets[r.ClientID],
			WarmRelation: isWarm(r.RelationshipNote),
		}
		if r.CurrentStage.Valid {
			in.LastSent = Stage(r.CurrentStage.String)
		}
		d := Decide(in)

		switch d.Action {
		case Draft:
			already, err := e.Q.HasPendingDraftForStage(ctx, db.HasPendingDraftForStageParams{
				InvoiceID: r.ID, Stage: string(d.Stage),
			})
			if err == nil && already > 0 {
				_ = e.Q.SetInvoiceAgentState(ctx, db.SetInvoiceAgentStateParams{AgentState: "awaiting_approval", ID: r.ID})
				continue
			}
			if err := e.draftAndMaybeSend(ctx, r, d.Stage, mode, d.Reasoning); err != nil {
				return "", err
			}
			if mode == "full_auto" && d.Stage.Ordinal() < 3 {
				sent++
			} else {
				drafted++
			}

		case Schedule:
			next := sql.NullString{}
			if r.NextActionOn.Valid {
				next = r.NextActionOn
			}
			stageChanged := !r.CurrentStage.Valid || Stage(r.CurrentStage.String).Ordinal() < d.Stage.Ordinal()-0
			if next.String != d.ActOn || (!r.NextActionOn.Valid) || stageChanged && r.AgentState != "awaiting_approval" {
				if err := e.Q.SetInvoiceChase(ctx, db.SetInvoiceChaseParams{
					CurrentStage:  r.CurrentStage,
					NextActionOn:  sql.NullString{String: d.ActOn, Valid: true},
					AgentState:    "planning",
					ID:            r.ID,
				}); err != nil {
					return "", err
				}
				if r.NextActionOn.String != d.ActOn {
					e.log(ctx, &r.ID, "plan_made", d.Reasoning)
				}
				planned++
			}

		case Hold:
			if r.AgentState == "awaiting_approval" || r.AgentState == "planning" || r.AgentState == "idle" {
				_ = e.Q.SetInvoiceAgentState(ctx, db.SetInvoiceAgentStateParams{AgentState: "idle", ID: r.ID})
			}

		case Stop:
			_ = e.Q.SetInvoiceAgentState(ctx, db.SetInvoiceAgentStateParams{AgentState: "stopped", ID: r.ID})
		}
	}

	switch {
	case drafted+sent+planned == 0:
		return "Nothing needs attention — every open invoice is on schedule.", nil
	default:
		parts := []string{}
		if sent > 0 {
			parts = append(parts, fmt.Sprintf("%d sent", sent))
		}
		if drafted > 0 {
			parts = append(parts, fmt.Sprintf("%d awaiting approval", drafted))
		}
		if planned > 0 {
			parts = append(parts, fmt.Sprintf("%d planned", planned))
		}
		return fmt.Sprintf("Agent run — %s.", strings.Join(parts, ", ")), nil
	}
}

// draftAndMaybeSend creates the draft; in full-auto (stages 0–2) it also
// sends. Stage 3 (final notice) and stage 4 (dossier) always wait for a human.
func (e *Engine) draftAndMaybeSend(ctx context.Context, r db.ListOpenInvoicesWithPayerRow, stage Stage, mode, reasoning string) error {
	settings, err := e.Q.GetSettings(ctx)
	if err != nil {
		return err
	}
	chases, _ := e.Q.CountSentChases(ctx, r.ID)
	facts := ChaseFacts{
		Stage:       stage,
		PayerName:   r.ClientName,
		PayerEmail:  r.ClientEmail,
		InvoiceNo:   r.Number,
		AmountCents: r.AmountCents,
		Currency:    r.Currency,
		DueOn:       r.DueOn,
		DaysOverdue: days(r.DueOn, nowToday(ctx, e)),
		PaymentLink: fmt.Sprintf("https://pay.lunas.app/i/%s", strings.ToLower(strings.TrimPrefix(r.Number, "INV-"))),
		SenderName:  settings.SenderName,
		SenderEmail: settings.SenderEmail,
		ChasesSent:  int(chases),
	}
	email, err := e.Provider.Draft(ctx, facts)
	if err != nil {
		return err
	}

	draft, err := e.Q.InsertDraft(ctx, db.InsertDraftParams{
		InvoiceID: r.ID, Stage: string(stage), Subject: email.Subject, Body: email.Body, Status: "pending",
	})
	if err != nil {
		return err
	}

	auto := mode == "full_auto" && stage.Ordinal() < 3
	if auto {
		return e.SendDraft(ctx, draft.ID)
	}
	if err := e.Q.SetInvoiceAgentState(ctx, db.SetInvoiceAgentStateParams{AgentState: "awaiting_approval", ID: r.ID}); err != nil {
		return err
	}
	e.log(ctx, &r.ID, "draft_ready", fmt.Sprintf("Drafted the %s for approval — %s", stage.Label(), strings.ToLower(reasoning[0:1])+reasoning[1:]))
	return nil
}

// SendDraft approves and sends one draft into the outbox, advancing the
// ladder and scheduling the next step.
func (e *Engine) SendDraft(ctx context.Context, draftID int64) error {
	d, err := e.Q.GetDraftWithInvoice(ctx, draftID)
	if err != nil {
		return err
	}
	now := nowISO()
	if err := e.Q.MarkDraftSent(ctx, db.MarkDraftSentParams{SentAt: sql.NullString{String: now, Valid: true}, ID: draftID}); err != nil {
		return err
	}
	if err := e.Q.InsertOutbox(ctx, db.InsertOutboxParams{
		InvoiceID: d.InvoiceID, InvoiceNumber: d.InvoiceNumber, ToName: d.ClientName,
		ToEmail: d.ClientEmail, Subject: d.Subject, Body: d.Body,
		SentAt: now,
	}); err != nil {
		return err
	}

	today := e.today(ctx)
	stage := Stage(d.Stage)
	next := nextAfter(stage)
	var nextOn sql.NullString
	state := "planning"
	if next == "" {
		state = "stopped"
	} else {
		nextOn = sql.NullString{String: actDate(next, d.DueOn, ""), Valid: true}
		if nextOn.String <= today {
			nextOn = sql.NullString{String: addDays(today, 1), Valid: true}
		}
	}
	if err := e.Q.SetInvoiceChase(ctx, db.SetInvoiceChaseParams{
		CurrentStage: sql.NullString{String: string(stage), Valid: true},
		NextActionOn: nextOn, AgentState: state, ID: d.InvoiceID,
	}); err != nil {
		return err
	}
	e.log(ctx, &d.InvoiceID, "draft_sent",
		fmt.Sprintf("Sent the %s (%s) to %s.", stage.Label(), stage.Tone(), d.ClientEmail))
	return nil
}

func (e *Engine) log(ctx context.Context, invoiceID *int64, typ, msg string) {
	var id sql.NullInt64
	if invoiceID != nil {
		id = sql.NullInt64{Int64: *invoiceID, Valid: true}
	}
	_ = e.Q.InsertActivity(ctx, db.InsertActivityParams{InvoiceID: id, Type: typ, Message: msg})
}

func nowToday(_ context.Context, e *Engine) string { return e.today(nil) }

// today returns the effective clock: settings.sim_now when set, else real today.
func (e *Engine) today(_ context.Context) string {
	if s, err := e.Q.GetSettings(context.Background()); err == nil && s.SimNow.Valid {
		return s.SimNow.String
	}
	return time.Now().UTC().Format("2006-01-02")
}
