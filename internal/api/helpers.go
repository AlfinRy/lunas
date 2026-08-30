package api

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	oatypes "github.com/oapi-codegen/runtime/types"
)

var errNoRows = sql.ErrNoRows

func nowUTC() time.Time { return time.Now().UTC() }

// API error constructors -----------------------------------------------------

func internal(err error) error {
	return &echo.HTTPError{Code: 500, Message: "Unable to load data. Try again.", Internal: err}
}

func errNotYet() error {
	return &echo.HTTPError{Code: 501, Message: "This part of the agent ships in week 2."}
}

func errMessage(msg string) Error {
	return Error{Message: msg}
}

func errFields(fe fieldErrors) Error {
	f := fe.orEmpty()
	return Error{Message: "Check the highlighted fields.", Fields: &f}
}

// validation helpers ---------------------------------------------------------

type fieldErrors map[string][]string

func (fe fieldErrors) orEmpty() map[string][]string {
	if fe == nil {
		return map[string][]string{}
	}
	return fe
}

func validEmail(s string) bool {
	if strings.ContainsAny(s, " \t\r\n") {
		return false
	}
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-1 && strings.Contains(s[at:], ".")
}

func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func validCurrency(s string) bool {
	if len(s) != 3 {
		return false
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

// time helpers ---------------------------------------------------------------

func parseTS(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Time{}
}

func daysBetween(fromISO, toISO string) int {
	from, err1 := time.Parse("2006-01-02", fromISO)
	to, err2 := time.Parse("2006-01-02", toISO)
	if err1 != nil || err2 != nil {
		return 0
	}
	return int(to.Sub(from).Hours() / 24)
}

// pointer helpers (API optional fields → sqlc nullable params) ----------------

func strPtrOrNil(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func intPtrOrNil(p *int) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*p), Valid: true}
}

func int64PtrOrNil(p *int) sql.NullInt64 { return intPtrOrNil(p) }

func datePtrOrNil(p *oatypes.Date) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: p.Time.Format("2006-01-02"), Valid: true}
}

func strDatePtrOrNull(p sql.NullString) *oatypes.Date {
	if !p.Valid {
		return nil
	}
	d := oatypes.Date{Time: isoDate(p.String)}
	return &d
}

func statusPtrOrNil(p *InvoiceStatus) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*p), Valid: true}
}

func modePtr(p *SettingsUpdateGlobalMode) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*p), Valid: true}
}

func strVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func nullInt64(v int64) sql.NullInt64 {
	return sql.NullInt64{Int64: v, Valid: true}
}

// score computation ----------------------------------------------------------

func (h *Handler) clientScores(ctx context.Context) (map[int64]*scoreCalc, error) {
	rows, err := h.q.PaymentStats(ctx)
	if err != nil {
		return nil, err
	}
	scores := map[int64]*scoreCalc{}
	for _, r := range rows {
		sc, ok := scores[r.ClientID]
		if !ok {
			sc = &scoreCalc{}
			scores[r.ClientID] = sc
		}
		sc.add(float64(daysBetween(r.DueOn, r.PaidOn)))
	}
	return scores, nil
}

var _ = errors.Is

// simToday returns the effective "now" date: the demo clock when set, real today otherwise.
func (h *Handler) simToday(ctx context.Context) string {
	if s, err := h.q.GetSettings(ctx); err == nil && s.SimNow.Valid {
		return s.SimNow.String
	}
	return time.Now().UTC().Format("2006-01-02")
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int64:
		return int(x)
	case int:
		return x
	case float64:
		return int(x)
	}
	return 0
}

func emailStrPtr(e oatypes.Email) *string {
	s := string(e)
	return &s
}

// dstr renders an API date as an ISO date string.
func dstr(d oatypes.Date) string { return d.Time.Format("2006-01-02") }

func emailPtrToStrPtr(e *oatypes.Email) sql.NullString {
	if e == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*e), Valid: true}
}
