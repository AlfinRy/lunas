package api

import (
	"time"

	db "github.com/AlfinRy/lunas/internal/db"
	oatypes "github.com/oapi-codegen/runtime/types"
)

func isoDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func apiDate(s string) oatypes.Date {
	return oatypes.Date{Time: isoDate(s)}
}

// toClient converts a DB row into the API shape, computing the payment score
// from settled history when enough evidence exists (PRD F8).
func toClient(c db.Client, scores map[int64]*scoreCalc) Client {
	out := Client{
		Id:              c.ID,
		Name:            c.Name,
		Email:           oatypes.Email(c.Email),
		PaymentTermsDays: int(c.PaymentTermsDays),
	}
	if c.RelationshipNote != "" {
		s := c.RelationshipNote
		out.RelationshipNote = &s
	}
	if sc, ok := scores[c.ID]; ok && sc.count >= 2 {
		out.PaymentScore = &PaymentScore{
			AvgDaysLate:  float32(sc.avgLate()),
			Reliability:  sc.reliability(),
			SettledCount: sc.count,
		}
	}
	return out
}

type scoreCalc struct {
	sumLate float64
	count   int
}

func (s *scoreCalc) add(daysLate float64) { s.sumLate += daysLate; s.count++ }

func (s *scoreCalc) avgLate() float64 {
	if s.count == 0 {
		return 0
	}
	// One decimal, half-up, for display stability.
	return float64(int(s.sumLate/float64(s.count)*10+0.5)) / 10
}

func (s *scoreCalc) reliability() Reliability {
	switch avg := s.avgLate(); {
	case avg <= 1:
		return PaysOnTime
	case avg <= 7:
		return UsuallyLate
	default:
		return ChronicallyLate
	}
}

func toInvoice(i db.ListInvoicesRow, today string) Invoice {
	paid := int(i.AmountPaidCents)
	days := daysBetween(i.DueOn, today)
	if !isOpen(i.Status) || days < 0 {
		days = 0
	}
	out := Invoice{
		Id:              i.ID,
		ClientId:        i.ClientID,
		ClientName:      i.ClientName,
		Number:          i.Number,
		AmountCents:     int(i.AmountCents),
		AmountPaidCents: &paid,
		Currency:        i.Currency,
		IssuedOn:        apiDate(i.IssuedOn),
		DueOn:           apiDate(i.DueOn),
		Status:          InvoiceStatus(i.Status),
		AgentState:      AgentState(i.AgentState),
		DaysOverdue:     &days,
		CreatedAt:       parseTS(i.CreatedAt),
	}
	if i.Notes != "" {
		n := i.Notes
		out.Notes = &n
	}
	if i.CurrentStage.Valid {
		s := ChaseStage(i.CurrentStage.String)
		out.CurrentStage = &s
	}
	if i.NextActionOn.Valid {
		d := apiDate(i.NextActionOn.String)
		out.NextActionOn = &d
	}
	return out
}

// invoiceFromFull adapts the single-row query result to the shared converter.
func invoiceFromFull(i db.GetInvoiceRow, today string) Invoice {
	row := db.ListInvoicesRow{
		ID: i.ID, ClientID: i.ClientID, Number: i.Number, AmountCents: i.AmountCents,
		AmountPaidCents: i.AmountPaidCents, Currency: i.Currency, IssuedOn: i.IssuedOn,
		DueOn: i.DueOn, Status: i.Status, Notes: i.Notes, AgentState: i.AgentState,
		CurrentStage: i.CurrentStage, NextActionOn: i.NextActionOn, CreatedAt: i.CreatedAt,
		ClientName: i.ClientName,
	}
	return toInvoice(row, today)
}

func toActivity(a db.Activity) Activity {
	out := Activity{
		Id:        a.ID,
		Type:      ActivityType(a.Type),
		Message:   a.Message,
		CreatedAt: parseTS(a.CreatedAt),
	}
	if a.InvoiceID.Valid {
		id := int(a.InvoiceID.Int64)
		out.InvoiceId = &id
	}
	return out
}

func isOpen(status string) bool {
	switch status {
	case "scheduled", "chasing", "paused", "draft":
		return true
	}
	return false
}
