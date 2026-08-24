// Package api implements the generated strict server interface. W1 scope:
// health, dashboard, clients, invoices, activity, settings, demo reset.
// Agent endpoints return 501 until W2/W3 land them.
package api

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"

	oatypes "github.com/oapi-codegen/runtime/types"

	db "github.com/AlfinRy/lunas/internal/db"
	"github.com/AlfinRy/lunas/internal/config"
	"github.com/AlfinRy/lunas/internal/seed"
)

type Handler struct {
	q   *db.Queries
	cfg *config.Config
}

func New(q *db.Queries, cfg *config.Config) *Handler {
	return &Handler{q: q, cfg: cfg}
}

// health ---------------------------------------------------------------------

func (h *Handler) GetHealth(ctx context.Context, _ GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{
		Status: "ok",
		Time:   nowUTC(),
	}, nil
}

// dashboard ------------------------------------------------------------------

func (h *Handler) GetDashboard(ctx context.Context, _ GetDashboardRequestObject) (GetDashboardResponseObject, error) {
	today := h.simToday(ctx)
	rows, err := h.q.ListInvoices(ctx)
	if err != nil {
		return nil, internal(err)
	}

	var outstanding, overdue int64
	counts := struct {
		AwaitingApproval int `json:"awaiting_approval"`
		Chasing          int `json:"chasing"`
		DueSoon          int `json:"due_soon"`
		Overdue          int `json:"overdue"`
	}{}
	for _, r := range rows {
		switch {
		case !isOpen(r.Status):
			continue
		case r.DueOn < today:
			overdue += r.AmountCents - r.AmountPaidCents
			counts.Overdue++
		case daysBetween(today, r.DueOn) <= 3:
			counts.DueSoon++
		}
		outstanding += r.AmountCents - r.AmountPaidCents
		if r.Status == "chasing" {
			counts.Chasing++
		}
	}

	recovered, err := h.q.RecoveredTotal(ctx)
	if err != nil {
		return nil, internal(err)
	}
	pending, err := h.q.CountPendingDrafts(ctx)
	if err != nil {
		return nil, internal(err)
	}
	counts.AwaitingApproval = int(pending)

	dso := h.rollingDSO(ctx, today)

	return GetDashboard200JSONResponse{
		Counts:           counts,
		DsoDays:          dso,
		OutstandingCents: int(outstanding),
		OverdueCents:     int(overdue),
		RecoveredCents:   toInt(recovered),
		SimNow:           apiDate(today),
	}, nil
}

// rollingDSO = average (paid_on − issued_on) over invoices paid in the last 30
// demo-days. Zero when nothing settled in the window.
func (h *Handler) rollingDSO(ctx context.Context, today string) float32 {
	rows, err := h.q.PaymentStats(ctx)
	if err != nil {
		return 0
	}
	var sum float64
	var n int
	for _, r := range rows {
		if daysBetween(r.PaidOn, today) > 30 || daysBetween(r.PaidOn, today) < 0 {
			continue
		}
		sum += float64(daysBetween(r.IssuedOn, r.PaidOn))
		n++
	}
	if n == 0 {
		return 0
	}
	return float32(sum / float64(n))
}

// clients --------------------------------------------------------------------

func (h *Handler) ListClients(ctx context.Context, _ ListClientsRequestObject) (ListClientsResponseObject, error) {
	rows, err := h.q.ListClients(ctx)
	if err != nil {
		return nil, internal(err)
	}
	scores, err := h.clientScores(ctx)
	if err != nil {
		return nil, internal(err)
	}
	out := make([]Client, 0, len(rows))
	for _, c := range rows {
		out = append(out, toClient(c, scores))
	}
	return ListClients200JSONResponse(out), nil
}

func (h *Handler) CreateClient(ctx context.Context, req CreateClientRequestObject) (CreateClientResponseObject, error) {
	fe := fieldErrors{}
	name := strings.TrimSpace(req.Body.Name)
	if name == "" {
		fe["name"] = []string{"Enter the payer's name."}
	}
	if !validEmail(string(req.Body.Email)) {
		fe["email"] = []string{"Enter the payer's email, like billing@meridian.co."}
	}
	if len(fe) > 0 {
		return CreateClient400JSONResponse{BadRequestJSONResponse(errFields(fe))}, nil
	}
	terms := 14
	if req.Body.PaymentTermsDays != nil {
		terms = *req.Body.PaymentTermsDays
	}
	note := ""
	if req.Body.RelationshipNote != nil {
		note = *req.Body.RelationshipNote
	}
	row, err := h.q.CreateClient(ctx, db.CreateClientParams{
		Name: name, Email: string(req.Body.Email), PaymentTermsDays: int64(terms), RelationshipNote: note,
	})
	if err != nil {
		return nil, internal(err)
	}
	scores, _ := h.clientScores(ctx)
	return CreateClient201JSONResponse(toClient(row, scores)), nil
}

func (h *Handler) GetClient(ctx context.Context, req GetClientRequestObject) (GetClientResponseObject, error) {
	row, err := h.q.GetClient(ctx, req.Id)
	if errors.Is(err, errNoRows) {
		return GetClient404JSONResponse{NotFoundJSONResponse(errMessage("Payer not found."))}, nil
	}
	if err != nil {
		return nil, internal(err)
	}
	scores, _ := h.clientScores(ctx)
	return GetClient200JSONResponse(toClient(row, scores)), nil
}

func (h *Handler) UpdateClient(ctx context.Context, req UpdateClientRequestObject) (UpdateClientResponseObject, error) {
	if _, err := h.q.GetClient(ctx, req.Id); errors.Is(err, errNoRows) {
		return UpdateClient404JSONResponse{NotFoundJSONResponse(errMessage("Payer not found."))}, nil
	} else if err != nil {
		return nil, internal(err)
	}
	fe := fieldErrors{}
	if req.Body.Email != nil && !validEmail(string(*req.Body.Email)) {
		fe["email"] = []string{"Enter the payer's email, like billing@meridian.co."}
	}
	if req.Body.Name != nil && strings.TrimSpace(*req.Body.Name) == "" {
		fe["name"] = []string{"Enter the payer's name."}
	}
	if len(fe) > 0 {
		return UpdateClient400JSONResponse{BadRequestJSONResponse(errFields(fe))}, nil
	}
	row, err := h.q.UpdateClient(ctx, db.UpdateClientParams{
		ID:                req.Id,
		Name:              strPtrOrNil(req.Body.Name),
		Email:             emailPtrToStrPtr(req.Body.Email),
		PaymentTermsDays:  intPtrOrNil(req.Body.PaymentTermsDays),
		RelationshipNote:  strPtrOrNil(req.Body.RelationshipNote),
	})
	if err != nil {
		return nil, internal(err)
	}
	scores, _ := h.clientScores(ctx)
	return UpdateClient200JSONResponse(toClient(row, scores)), nil
}

// invoices -------------------------------------------------------------------

func (h *Handler) ListInvoices(ctx context.Context, req ListInvoicesRequestObject) (ListInvoicesResponseObject, error) {
	today := h.simToday(ctx)
	rows, err := h.q.ListInvoices(ctx)
	if err != nil {
		return nil, internal(err)
	}
	if req.Params.Status != nil {
		filtered := rows[:0]
		for _, r := range rows {
			if r.Status == string(*req.Params.Status) {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}
	if req.Params.ClientId != nil {
		filtered := rows[:0]
		for _, r := range rows {
			if r.ClientID == *req.Params.ClientId {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}
	// Overdue first (due_on already sorts past-before-future), then by due date.
	sort.SliceStable(rows, func(i, j int) bool {
		io, jo := isOpen(rows[i].Status) && rows[i].DueOn < today, isOpen(rows[j].Status) && rows[j].DueOn < today
		if io != jo {
			return io
		}
		return rows[i].DueOn < rows[j].DueOn
	})
	out := make([]Invoice, 0, len(rows))
	for _, r := range rows {
		out = append(out, toInvoice(r, today))
	}
	return ListInvoices200JSONResponse(out), nil
}

func (h *Handler) CreateInvoice(ctx context.Context, req CreateInvoiceRequestObject) (CreateInvoiceResponseObject, error) {
	b := req.Body
	fe := fieldErrors{}
	if strings.TrimSpace(b.Number) == "" {
		fe["number"] = []string{"Enter an invoice number, like INV-0042."}
	}
	if _, err := h.q.GetClient(ctx, b.ClientId); errors.Is(err, errNoRows) {
		fe["client_id"] = []string{"Pick the payer for this invoice."}
	} else if err != nil {
		return nil, internal(err)
	}
	if b.AmountCents < 1 {
		fe["amount_cents"] = []string{"Enter an amount greater than zero."}
	}
	if !validCurrency(b.Currency) {
		fe["currency"] = []string{"Use a 3-letter currency code, like USD."}
	}
	if !validDate(dstr(b.IssuedOn)) || !validDate(dstr(b.DueOn)) {
		fe["due_on"] = []string{"Pick valid issue and due dates."}
	} else if dstr(b.DueOn) < dstr(b.IssuedOn) {
		fe["due_on"] = []string{"Pick a due date after the issue date."}
	}
	if len(fe) > 0 {
		return CreateInvoice400JSONResponse{BadRequestJSONResponse(errFields(fe))}, nil
	}

	dup, err := h.q.CountInvoiceNumberForClient(ctx, db.CountInvoiceNumberForClientParams{
		ClientID: b.ClientId, Number: b.Number,
	})
	if err != nil {
		return nil, internal(err)
	}
	if dup > 0 {
		return CreateInvoice409JSONResponse{ConflictJSONResponse(errMessage("This payer already has an invoice with that number."))}, nil
	}

	status := "scheduled"
	if dstr(b.DueOn) < h.simToday(ctx) {
		status = "chasing"
	}
	row, err := h.q.CreateInvoice(ctx, db.CreateInvoiceParams{
		ClientID: b.ClientId, Number: b.Number, AmountCents: int64(b.AmountCents), Currency: b.Currency,
		IssuedOn: dstr(b.IssuedOn), DueOn: dstr(b.DueOn), Status: status,
		AgentState: "planning", Notes: strVal(b.Notes),
	})
	if err != nil {
		return nil, internal(err)
	}
	_ = h.q.InsertActivity(ctx, db.InsertActivityParams{
		InvoiceID: nullInt64(row.ID), Type: "invoice_created",
		Message: "Invoice " + row.Number + " added — the agent starts planning its chase.",
	})
	full, err := h.q.GetInvoice(ctx, row.ID)
	if err != nil {
		return nil, internal(err)
	}
	return CreateInvoice201JSONResponse(invoiceFromFull(full, h.simToday(ctx))), nil
}

func (h *Handler) GetInvoice(ctx context.Context, req GetInvoiceRequestObject) (GetInvoiceResponseObject, error) {
	row, err := h.q.GetInvoice(ctx, req.Id)
	if errors.Is(err, errNoRows) {
		return GetInvoice404JSONResponse{NotFoundJSONResponse(errMessage("Invoice not found."))}, nil
	}
	if err != nil {
		return nil, internal(err)
	}
	return GetInvoice200JSONResponse(invoiceFromFull(row, h.simToday(ctx))), nil
}

func (h *Handler) UpdateInvoice(ctx context.Context, req UpdateInvoiceRequestObject) (UpdateInvoiceResponseObject, error) {
	cur, err := h.q.GetInvoice(ctx, req.Id)
	if errors.Is(err, errNoRows) {
		return UpdateInvoice404JSONResponse{NotFoundJSONResponse(errMessage("Invoice not found."))}, nil
	}
	if err != nil {
		return nil, internal(err)
	}
	b := req.Body
	fe := fieldErrors{}
	if b.DueOn != nil && dstr(*b.DueOn) < cur.IssuedOn {
		fe["due_on"] = []string{"Pick a due date after the issue date."}
	}
	if b.AmountCents != nil && *b.AmountCents < 1 {
		fe["amount_cents"] = []string{"Enter an amount greater than zero."}
	}
	if b.Number != nil {
		dup, err := h.q.CountInvoiceNumberForClient(ctx, db.CountInvoiceNumberForClientParams{
			ClientID: cur.ClientID, Number: *b.Number, ExcludeID: nullInt64(cur.ID),
		})
		if err == nil && dup > 0 {
			fe["number"] = []string{"This payer already has an invoice with that number."}
		}
	}
	if len(fe) > 0 {
		return UpdateInvoice400JSONResponse{BadRequestJSONResponse(errFields(fe))}, nil
	}

	// Status transitions with side effects (pause/resume/write-off).
	agentState := cur.AgentState
	if b.Status != nil {
		switch *b.Status {
		case "paused":
			agentState = "holding"
		case "chasing", "scheduled":
			if cur.Status == "paused" {
				agentState = "planning"
			}
		case "written_off":
			agentState = "stopped"
		case "paid":
			agentState = "stopped"
		}
		_ = h.q.InsertActivity(ctx, db.InsertActivityParams{
			InvoiceID: nullInt64(cur.ID), Type: "note",
			Message: "Status changed from " + cur.Status + " to " + string(*b.Status) + ".",
		})
	}

	row, err := h.q.UpdateInvoiceFields(ctx, db.UpdateInvoiceFieldsParams{
		ID:              req.Id,
		Number:          strPtrOrNil(b.Number),
		AmountCents:     int64PtrOrNil(b.AmountCents),
		DueOn:           datePtrOrNil(b.DueOn),
		Notes:           strPtrOrNil(b.Notes),
		Status:          statusPtrOrNil(b.Status),
		AgentState:      sql.NullString{String: agentState, Valid: true},
		CurrentStage:    cur.CurrentStage,
		NextActionOn:    cur.NextActionOn,
		AmountPaidCents: sql.NullInt64{Int64: cur.AmountPaidCents, Valid: true},
	})
	if err != nil {
		return nil, internal(err)
	}
	full, err := h.q.GetInvoice(ctx, row.ID)
	if err != nil {
		return nil, internal(err)
	}
	return UpdateInvoice200JSONResponse(invoiceFromFull(full, h.simToday(ctx))), nil
}

func (h *Handler) ListInvoiceActivity(ctx context.Context, req ListInvoiceActivityRequestObject) (ListInvoiceActivityResponseObject, error) {
	if _, err := h.q.GetInvoice(ctx, req.Id); errors.Is(err, errNoRows) {
		return ListInvoiceActivity404JSONResponse{NotFoundJSONResponse(errMessage("Invoice not found."))}, nil
	} else if err != nil {
		return nil, internal(err)
	}
	rows, err := h.q.ListActivityForInvoice(ctx, nullInt64(req.Id))
	if err != nil {
		return nil, internal(err)
	}
	out := make([]Activity, 0, len(rows))
	for _, a := range rows {
		out = append(out, toActivity(a))
	}
	return ListInvoiceActivity200JSONResponse(out), nil
}

// settings -------------------------------------------------------------------

func (h *Handler) toSettings(s db.Setting) Settings {
	return Settings{
		SenderName:       s.SenderName,
		SenderEmail:      oatypes.Email(s.SenderEmail),
		DefaultTermsDays: int(s.DefaultTermsDays),
		GlobalMode:       SettingsGlobalMode(s.GlobalMode),
		TemplateMode:     h.cfg.TemplateMode(),
		SimNow:           strDatePtrOrNull(s.SimNow),
	}
}

func (h *Handler) GetSettings(ctx context.Context, _ GetSettingsRequestObject) (GetSettingsResponseObject, error) {
	s, err := h.q.GetSettings(ctx)
	if err != nil {
		return nil, internal(err)
	}
	return GetSettings200JSONResponse(h.toSettings(s)), nil
}

func (h *Handler) UpdateSettings(ctx context.Context, req UpdateSettingsRequestObject) (UpdateSettingsResponseObject, error) {
	cur, err := h.q.GetSettings(ctx)
	if err != nil {
		return nil, internal(err)
	}
	fe := fieldErrors{}
	if req.Body.SenderEmail != nil && !validEmail(string(*req.Body.SenderEmail)) {
		fe["sender_email"] = []string{"Enter your billing email, like billing@ranistudio.co."}
	}
	if len(fe) > 0 {
		return UpdateSettings400JSONResponse{BadRequestJSONResponse(errFields(fe))}, nil
	}
	s, err := h.q.UpdateSettings(ctx, db.UpdateSettingsParams{
		SenderName:       strPtrOrNil(req.Body.SenderName),
		SenderEmail:      emailPtrToStrPtr(req.Body.SenderEmail),
		DefaultTermsDays: intPtrOrNil(req.Body.DefaultTermsDays),
		GlobalMode:       modePtr(req.Body.GlobalMode),
		SimNow:           cur.SimNow,
	})
	if err != nil {
		return nil, internal(err)
	}
	return UpdateSettings200JSONResponse(h.toSettings(s)), nil
}

// demo reset -----------------------------------------------------------------

func (h *Handler) ResetDemo(ctx context.Context, _ ResetDemoRequestObject) (ResetDemoResponseObject, error) {
	if err := h.q.ResetDemoData(ctx); err != nil {
		return nil, internal(err)
	}
	if err := h.q.InitSettings(ctx); err != nil {
		return nil, internal(err)
	}
	if err := seed.Run(ctx, h.q); err != nil {
		return nil, internal(err)
	}
	return ResetDemo200JSONResponse{Ok: true}, nil
}

// W2/W3 stubs ----------------------------------------------------------------

func (h *Handler) GetAgentInbox(context.Context, GetAgentInboxRequestObject) (GetAgentInboxResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) ApproveDraft(context.Context, ApproveDraftRequestObject) (ApproveDraftResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) SkipDraft(context.Context, SkipDraftRequestObject) (SkipDraftResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) ListOutbox(context.Context, ListOutboxRequestObject) (ListOutboxResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) ParsePayment(context.Context, ParsePaymentRequestObject) (ParsePaymentResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) ReconcilePayment(context.Context, ReconcilePaymentRequestObject) (ReconcilePaymentResponseObject, error) {
	return nil, errNotYet()
}
func (h *Handler) SimulateAdvance(context.Context, SimulateAdvanceRequestObject) (SimulateAdvanceResponseObject, error) {
	return nil, errNotYet()
}
