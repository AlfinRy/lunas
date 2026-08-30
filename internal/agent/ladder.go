package agent

import (
	"fmt"
	"strings"
)

// ChaseFacts is the structured context handed to any drafting provider —
// template or LLM. Same facts either way: the decision layer owns what gets
// said; the provider only shapes the words.
type ChaseFacts struct {
	Stage        Stage
	PayerName    string
	PayerEmail   string
	InvoiceNo    string
	AmountCents  int64
	Currency     string
	DueOn        string
	DaysOverdue  int
	PaymentLink  string
	SenderName   string
	SenderEmail  string
	ChasesSent   int
}

// DraftEmail is a rendered chase.
type DraftEmail struct {
	Subject string
	Body    string
}

func firstName(full string) string {
	if i := strings.IndexByte(full, ' '); i > 0 {
		return full[:i]
	}
	return full
}

func money(cents int64, currency string) string {
	if currency == "USD" {
		return "$" + FormatAmount(cents)
	}
	return currency + " " + FormatAmount(cents)
}

// FormatAmount renders 123456 → "1,234.56" without locale magic.
func FormatAmount(cents int64) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	whole := cents / 100
	frac := cents % 100
	s := fmt.Sprintf("%d", whole)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := strings.Join(parts, ",") + fmt.Sprintf(".%02d", frac)
	if neg {
		return "-" + out
	}
	return out
}

// footer keeps the agent honest (BRD §9): automated, on behalf of the sender.
func footer(f ChaseFacts) string {
	return fmt.Sprintf("\n\n—\nThis is an automated reminder sent by Lunas on behalf of %s.", f.SenderName)
}

// Template renders the canonical ladder copy (docs/ux-writing.md §11).
// It is both the no-LLM fallback and the LLM's grounding example.
func Template(f ChaseFacts) DraftEmail {
	who := firstName(f.PayerName)
	amt := money(f.AmountCents, f.Currency)
	var d DraftEmail
	switch f.Stage {
	case Stage0Reminder:
		d.Subject = fmt.Sprintf("Invoice %s — due %s", f.InvoiceNo, f.DueOn)
		d.Body = fmt.Sprintf("Hi %s,\n\nA quick note that invoice %s for %s is due %s.\n\nPay here: %s\n\nIf it's already scheduled, feel free to ignore this note.\n\n%s%s",
			who, bold(f.InvoiceNo), bold(amt), f.DueOn, f.PaymentLink, f.SenderName, footer(f))

	case Stage1Polite:
		d.Subject = fmt.Sprintf("Invoice %s — %s, now %d days past due", f.InvoiceNo, amt, f.DaysOverdue)
		d.Body = fmt.Sprintf("Hi %s,\n\nMarking this as unpaid on our side — invoice %s for %s was due %s.\n\nIf anything's blocking on your end (approval, PO, paperwork), reply here and %s will sort it out today. Otherwise, the fastest way to settle is this link: %s\n\nThanks,\n%s%s",
			who, bold(f.InvoiceNo), bold(amt), f.DueOn, f.SenderName, f.PaymentLink, f.SenderName, footer(f))

	case Stage2Firm:
		d.Subject = fmt.Sprintf("Second notice: invoice %s is %d days overdue", f.InvoiceNo, f.DaysOverdue)
		d.Body = fmt.Sprintf("Hi %s,\n\nFollowing up again on invoice %s — %s, now %d days past due. We haven't received a reply to the previous note.\n\nPlease confirm a payment date this week, or settle directly: %s\n\nIf there's a dispute or question about the work, flag it now so we can resolve it quickly.\n\n%s%s",
			who, bold(f.InvoiceNo), bold(amt), f.DaysOverdue, f.PaymentLink, f.SenderName, footer(f))

	case Stage3Final:
		d.Subject = fmt.Sprintf("Final notice — invoice %s, %d days overdue", f.InvoiceNo, f.DaysOverdue)
		d.Body = fmt.Sprintf("Hi %s,\n\nThis is a final notice for invoice %s (%s, due %s). To date we've sent %d reminders without payment or an agreed payment date.\n\nPlease settle within 5 business days: %s\n\nIf payment isn't received or a plan agreed by %s, this account moves to formal collections and late fees may apply.\n\nWe'd rather resolve it directly. Reply with a payment date and this is closed.\n\n%s%s",
			who, bold(f.InvoiceNo), bold(amt), f.DueOn, f.ChasesSent, f.PaymentLink, addDays(f.DueOn, 21), f.SenderName, footer(f))

	case Stage4Escalation:
		// Not an email: the dossier line the UI renders as a packet summary.
		d.Subject = fmt.Sprintf("Collections dossier — invoice %s (%s)", f.InvoiceNo, amt)
		d.Body = fmt.Sprintf("Everything a collections process needs, in one file:\n\n• Invoice %s for %s, due %s (%d days overdue)\n• %d chases delivered — full timeline attached\n• Amount breakdown and payer contact: %s <%s>\n• Draft referral letter included\n\nReview and send from here when you're ready.",
			f.InvoiceNo, amt, f.DueOn, f.DaysOverdue, f.ChasesSent, f.PayerName, f.PayerEmail)
	}
	return d
}

func bold(s string) string { return s }
