// Package agent holds Lunas' decision layer: a deterministic chase policy
// engine (pure, table-tested), the escalation ladder templates, and the
// provider-agnostic drafting adapter. The LLM writes words; this package
// makes every decision.
package agent

import "time"

// Stage is the escalation ladder (docs/ux-writing.md §11).
type Stage string

const (
	Stage0Reminder    Stage = "stage0_reminder"
	Stage1Polite      Stage = "stage1_polite"
	Stage2Firm        Stage = "stage2_firm"
	Stage3Final       Stage = "stage3_final"
	Stage4Escalation  Stage = "stage4_escalation"
)

var stageOrder = map[Stage]int{
	Stage0Reminder: 0, Stage1Polite: 1, Stage2Firm: 2, Stage3Final: 3, Stage4Escalation: 4,
}

func (s Stage) Ordinal() int {
	if s == "" {
		return -1 // nothing sent yet — stage0 is an escalation from here
	}
	return stageOrder[s]
}

// Tone summarised for plan cards: gentle → polite → firm → final.
func (s Stage) Tone() string {
	switch s {
	case Stage0Reminder:
		return "gentle"
	case Stage1Polite:
		return "polite"
	case Stage2Firm:
		return "firm"
	case Stage3Final:
		return "final"
	default:
		return "escalation"
	}
}

func (s Stage) Label() string {
	switch s {
	case Stage0Reminder:
		return "due-date reminder"
	case Stage1Polite:
		return "polite follow-up"
	case Stage2Firm:
		return "firm follow-up"
	case Stage3Final:
		return "final notice"
	default:
		return "collections dossier"
	}
}

// Reliability buckets from the payer payment-score.
const (
	ReliabilityUnknown    = ""
	ReliabilityOnTime     = "pays_on_time"
	ReliabilityUsually    = "usually_late"
	ReliabilityChronic    = "chronically_late"
)

// Action is what the agent wants to do next for one invoice.
type Action int

const (
	Hold     Action = iota // closed, paused, or waiting for a future date
	Schedule               // nothing due yet; a stage is planned for a future date
	Draft                  // a stage is due now (or overdue) and needs drafting
	Stop                   // settled or written off — chasing must not continue
)

// Input is everything the policy may look at. No clocks, no IO — the caller
// passes today and the payer context; the function stays pure.
type Input struct {
	Today        string // sim clock or real today (ISO date)
	DueOn        string // invoice due date (ISO)
	Status       string // invoice status
	LastSent     Stage // last stage actually sent ("" when none)
	Reliability  string // payer reliability bucket ("" when unknown)
	WarmRelation bool   // relationship note pins the ladder at gentle stages
}

// Decision is the agent's full reasoning, shown to the user and the judges.
type Decision struct {
	Action    Action
	Stage     Stage // the stage this decision is about
	ActOn     string // when to act (for Schedule: the planned date; Draft: today)
	Reasoning string
}

func days(fromISO, toISO string) int {
	from, err1 := time.Parse("2006-01-02", fromISO)
	to, err2 := time.Parse("2006-01-02", toISO)
	if err1 != nil || err2 != nil {
		return 0
	}
	return int(to.Sub(from).Hours() / 24)
}

func addDays(iso string, n int) string {
	t, err := time.Parse("2006-01-02", iso)
	if err != nil {
		return iso
	}
	return t.AddDate(0, 0, n).Format("2006-01-02")
}

// stageFor maps days overdue to the ladder stage the invoice deserves.
func stageFor(daysOverdue int, warm bool) Stage {
	var s Stage
	switch {
	case daysOverdue <= 0:
		s = Stage0Reminder
	case daysOverdue <= 6:
		s = Stage1Polite
	case daysOverdue <= 13:
		s = Stage2Firm
	case daysOverdue <= 20:
		s = Stage3Final
	default:
		s = Stage4Escalation
	}
	// A "keep warm" relationship pins the ladder at gentle stages (0–1);
	// it slows collection but never blocks it silently — reasoning says so.
	if warm && s.Ordinal() > 1 {
		s = Stage1Polite
	}
	return s
}

// actDate returns when the given stage becomes due, relative to due_on.
// Stage 0 lands on the due date — one day early when the payer usually runs
// late, so the reminder arrives before the weekend slips it.
func actDate(stage Stage, dueOn string, reliability string) string {
	switch stage {
	case Stage0Reminder:
		if reliability == ReliabilityUsually || reliability == ReliabilityChronic {
			return addDays(dueOn, -1)
		}
		return dueOn
	case Stage1Polite:
		return addDays(dueOn, 2)
	case Stage2Firm:
		return addDays(dueOn, 7)
	case Stage3Final:
		return addDays(dueOn, 14)
	default:
		return addDays(dueOn, 21)
	}
}

// Decide is the policy engine. Pure: same input, same decision, every time.
func Decide(in Input) Decision {
	switch in.Status {
	case "paid", "written_off":
		return Decision{Action: Stop, Stage: in.LastSent,
			Reasoning: "Invoice is " + in.Status + " — chasing stopped."}
	case "paused":
		return Decision{Action: Hold, Stage: in.LastSent,
			Reasoning: "Chasing paused by you — Lunas holds until resumed."}
	case "draft":
		return Decision{Action: Hold, Stage: in.LastSent, Reasoning: "Draft invoice — not sent to the payer yet."}
	}

	overdue := days(in.DueOn, in.Today)
	target := stageFor(overdue, in.WarmRelation)

	// Never resend a stage already sent; escalate only.
	if target.Ordinal() <= in.LastSent.Ordinal() {
		next := nextAfter(in.LastSent)
		if next == "" || (in.WarmRelation && in.LastSent.Ordinal() >= 1) {
			return Decision{Action: Hold, Stage: in.LastSent, ActOn: "",
				Reasoning: "Every planned chase for this relationship has been sent (" +
					in.LastSent.Label() + "). Warm relationship — Lunas holds at gentle stages until you intervene."}
		}
		return Decision{Action: Schedule, Stage: next, ActOn: actDate(next, in.DueOn, in.Reliability),
			Reasoning: "Holding — next step is the " + next.Label() + " (" + next.Tone() + "), planned for " + actDate(next, in.DueOn, in.Reliability) + "."}
	}

	actOn := actDate(target, in.DueOn, in.Reliability)
	reason := reasoning(target, in.Reliability, overdue, in.WarmRelation)

	if in.Today >= actOn {
		return Decision{Action: Draft, Stage: target, ActOn: in.Today, Reasoning: reason}
	}
	return Decision{Action: Schedule, Stage: target, ActOn: actOn, Reasoning: reason}
}

func nextAfter(s Stage) Stage {
	switch s {
	case Stage0Reminder:
		return Stage1Polite
	case Stage1Polite:
		return Stage2Firm
	case Stage2Firm:
		return Stage3Final
	case Stage3Final:
		return Stage4Escalation
	default:
		return ""
	}
}

func reasoning(target Stage, reliability string, overdue int, warm bool) string {
	r := ""
	switch reliability {
	case ReliabilityOnTime:
		r = "pays on time"
	case ReliabilityUsually:
		r = "usually pays late"
	case ReliabilityChronic:
		r = "is usually well behind on payment"
	default:
		r = "has no settled history yet"
	}
	base := "Scheduling the " + target.Label() + " — this payer " + r + "."
	if overdue > 0 {
		base = "Drafting the " + target.Label() + " (" + target.Tone() + ") — invoice is " + itoa(overdue) + " days overdue and this payer " + r + "."
	}
	if warm {
		base += " Relationship note says keep it warm, so the ladder stays at gentle stages."
	}
	return base
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
