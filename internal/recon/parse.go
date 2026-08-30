// Package recon parses pasted payment notifications and matches them to open
// invoices. Deterministic and table-tested: this is the matcher that tells
// the agent when to stop chasing (PRD F5).
package recon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parsed is what the extractor found in raw text.
type Parsed struct {
	AmountCents int64
	Currency    string
	PaidOn      string // ISO date; "" when absent
	PayerHint   string // client name fragment, "" when absent
	InvoiceHint string // invoice number like INV-0041, "" when absent
}

var (
	reUSD      = regexp.MustCompile(`(?:\$|USD)\s*([0-9][0-9,]*(?:\.[0-9]{1,2})?)`)
	reIDR      = regexp.MustCompile(`Rp\.?\s*([0-9][0-9.]*)(?:,([0-9]{2}))?`)
	reBare     = regexp.MustCompile(`(?i)(?:amount|total|sebesar|jumlah|received|diterima)[:\s]+([0-9][0-9,]*\.[0-9]{2})`)
	reInvoice  = regexp.MustCompile(`(?i)inv[\s-]*([0-9]{2,8})`)
	reISODate  = regexp.MustCompile(`\b(20[0-9]{2})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])\b`)
	reMonthDay = regexp.MustCompile(`(?i)\b(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\.?\s+([0-9]{1,2})(?:st|nd|rd|th)?,?\s+(20[0-9]{2})\b`)
	reDayMonth = regexp.MustCompile(`(?i)\b([0-9]{1,2})\s+(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\.?\s+(20[0-9]{2})\b`)
	reSlash    = regexp.MustCompile(`\b([0-9]{1,2})/([0-9]{1,2})/([0-9]{2,4})\b`)
)

func parseAmount(s string) (int64, string, bool) {
	if m := reIDR.FindStringSubmatch(s); m != nil {
		whole := strings.ReplaceAll(m[1], ".", "")
		cents, _ := strconv.ParseInt(whole, 10, 64)
		cents *= 100
		if m[2] != "" {
			frac, _ := strconv.ParseInt(m[2], 10, 64)
			cents += frac
		}
		return cents, "IDR", true
	}
	if m := reUSD.FindStringSubmatch(s); m != nil {
		return dollars(m[1]), "USD", true
	}
	if m := reBare.FindStringSubmatch(s); m != nil {
		return dollars(m[1]), "USD", true
	}
	return 0, "", false
}

func dollars(s string) int64 {
	s = strings.ReplaceAll(s, ",", "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f*100 + 0.5)
}

var monthNum = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func parseDate(s string, today string) string {
	if m := reISODate.FindStringSubmatch(s); m != nil {
		return m[0]
	}
	if m := reMonthDay.FindStringSubmatch(s); m != nil {
		mo := monthNum[strings.ToLower(m[1])]
		return fmt.Sprintf("%s-%02d-%02d", m[3], mo, atoi(m[2]))
	}
	if m := reDayMonth.FindStringSubmatch(s); m != nil {
		mo := monthNum[strings.ToLower(m[2])]
		return fmt.Sprintf("%s-%02d-%02d", m[3], mo, atoi(m[1]))
	}
	if m := reSlash.FindStringSubmatch(s); m != nil {
		a, b, y := atoi(m[1]), atoi(m[2]), atoi(m[3])
		if len(m[3]) == 2 {
			y += 2000
		}
		// Prefer MM/DD/YYYY; flip when only DD/MM makes sense.
		mo, d := a, b
		if a > 12 && b <= 12 {
			mo, d = b, a
		}
		if mo >= 1 && mo <= 12 {
			return fmt.Sprintf("%04d-%02d-%02d", y, mo, d)
		}
	}
	return today // "received today" is the common case in pastes
}

func atoi(s string) int { n, _ := strconv.Atoi(strings.TrimPrefix(s, "0")); return n }

// normalize lowercases and strips punctuation for fuzzy matching.
func normalize(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// nameMatches: full normalized name appears in text, or ≥2 significant tokens.
func nameMatches(name, textNorm string) bool {
	n := normalize(name)
	if n == "" || strings.Contains(textNorm, n) {
		return n != ""
	}
	tokens := strings.Fields(n)
	hits := 0
	for _, t := range tokens {
		if len(t) > 2 && strings.Contains(textNorm, t) {
			hits++
		}
	}
	return hits >= 2
}

// Parse extracts everything it can; amount is required, the rest may be empty.
func Parse(text, today string) (Parsed, bool) {
	cents, currency, ok := parseAmount(text)
	if !ok || cents <= 0 {
		return Parsed{}, false
	}
	p := Parsed{AmountCents: cents, Currency: currency, PaidOn: parseDate(text, today)}
	if m := reInvoice.FindStringSubmatch(text); m != nil {
		p.InvoiceHint = "INV-" + m[1]
	}
	return p, true
}

// MatchCandidate is one open invoice scored against a parsed payment.
type MatchCandidate struct {
	InvoiceID   int64
	InvoiceNo   string
	ClientID    int64
	ClientName  string
	AmountCents int64
	DueOn       string
	Confidence  string // high | medium | low
}

// OpenInvoice is the matcher's view of an open invoice.
type OpenInvoice struct {
	InvoiceID   int64
	InvoiceNo   string
	ClientID    int64
	ClientName  string
	AmountCents int64
	AmountPaid  int64
	DueOn       string
}

// Match scores open invoices against the parsed payment.
// Invoice-number hits win outright; otherwise amount + payer evidence add up.
func Match(p Parsed, clients []string, open []OpenInvoice, text string) []MatchCandidate {
	textNorm := normalize(text)
	type scored struct {
		c    MatchCandidate
		score int
	}
	var out []scored
	for _, iv := range open {
		remaining := iv.AmountCents - iv.AmountPaid
		score := 0
		confidence := "low"

		exact := p.AmountCents == remaining
		close := remaining > 0 && p.AmountCents > 0 &&
			abs64(p.AmountCents-remaining)*100 <= remaining*2 // within 2%

		if exact {
			score += 50
		} else if close {
			score += 35
		}
		if p.InvoiceHint != "" && strings.EqualFold(p.InvoiceHint, iv.InvoiceNo) {
			score += 100
		}
		if int(iv.ClientID) < len(clients) {
			_ = iv // placeholder, real payer check below
		}
		if nameMatches(iv.ClientName, textNorm) {
			score += 30
		}

		switch {
		case score >= 80:
			confidence = "high"
		case score >= 50:
			confidence = "medium"
		case score >= 30: // payer-name evidence alone still earns a look
			confidence = "low"
		default:
			continue // not a plausible candidate
		}
		out = append(out, scored{MatchCandidate{
			InvoiceID: iv.InvoiceID, InvoiceNo: iv.InvoiceNo, ClientID: iv.ClientID,
			ClientName: iv.ClientName, AmountCents: iv.AmountCents, DueOn: iv.DueOn,
			Confidence: confidence,
		}, score})
	}

	// Sort by score desc, keep top 3.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].score > out[i].score {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if len(out) > 3 {
		out = out[:3]
	}
	result := make([]MatchCandidate, len(out))
	for i, s := range out {
		result[i] = s.c
	}
	return result
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// TodayISO returns real today (helper for handlers that seed the fallback).
func TodayISO() string { return time.Now().UTC().Format("2006-01-02") }
