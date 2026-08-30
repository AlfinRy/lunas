package agent

import (
	"strings"
	"testing"
)

// Table-driven: every row is a story the demo tells.
func TestDecide(t *testing.T) {
	tests := []struct {
		name string
		in   Input
		want Action
		wantStage Stage
		wantActOn string
	}{
		{
			name: "fresh invoice, payer on time, before due",
			in:   Input{Today: "2026-09-01", DueOn: "2026-09-05", Status: "scheduled", Reliability: ReliabilityOnTime},
			want: Schedule, wantStage: Stage0Reminder, wantActOn: "2026-09-05",
		},
		{
			name: "fresh invoice, usually-late payer → reminder one day early",
			in:   Input{Today: "2026-09-01", DueOn: "2026-09-05", Status: "scheduled", Reliability: ReliabilityUsually},
			want: Schedule, wantStage: Stage0Reminder, wantActOn: "2026-09-04",
		},
		{
			name: "due today, nothing sent yet → draft the reminder now",
			in:   Input{Today: "2026-09-05", DueOn: "2026-09-05", Status: "scheduled", Reliability: ReliabilityOnTime},
			want: Draft, wantStage: Stage0Reminder, wantActOn: "2026-09-05",
		},
		{
			name: "3 days overdue, nothing sent → polite due now (stage0 window passed)",
			in:   Input{Today: "2026-09-08", DueOn: "2026-09-05", Status: "chasing", Reliability: ""},
			want: Draft, wantStage: Stage1Polite, wantActOn: "2026-09-08",
		},
		{
			name: "9 days overdue, reminder+polite sent → firm due (the Meridian story)",
			in:   Input{Today: "2026-09-14", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage1Polite, Reliability: ReliabilityUsually},
			want: Draft, wantStage: Stage2Firm, wantActOn: "2026-09-14",
		},
		{
			name: "10 days overdue but firm already sent → wait for final window",
			in:   Input{Today: "2026-09-15", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage2Firm},
			want: Schedule, wantStage: Stage3Final, wantActOn: "2026-09-19",
		},
		{
			name: "15 days overdue, firm sent → final notice now",
			in:   Input{Today: "2026-09-20", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage2Firm},
			want: Draft, wantStage: Stage3Final, wantActOn: "2026-09-20",
		},
		{
			name: "25 days overdue, final sent → escalation dossier",
			in:   Input{Today: "2026-09-30", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage3Final},
			want: Draft, wantStage: Stage4Escalation, wantActOn: "2026-09-30",
		},
		{
			name: "warm relationship pins the ladder at polite",
			in:   Input{Today: "2026-09-20", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage1Polite, WarmRelation: true},
			want: Hold, wantStage: Stage1Polite,
		},
		{
			name: "paid → stop immediately",
			in:   Input{Today: "2026-09-20", DueOn: "2026-09-05", Status: "paid", LastSent: Stage1Polite},
			want: Stop, wantStage: Stage1Polite,
		},
		{
			name: "paused → hold",
			in:   Input{Today: "2026-09-08", DueOn: "2026-09-05", Status: "paused"},
			want: Hold, wantStage: "",
		},
		{
			name: "stage0 sent, still overdue ≤6, no more stages until +7",
			in:   Input{Today: "2026-09-06", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage0Reminder},
			want: Schedule, wantStage: Stage1Polite, wantActOn: "2026-09-07",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decide(tt.in)
			if got.Action != tt.want || got.Stage != tt.wantStage {
				t.Errorf("Decide() = %+v, want action=%v stage=%v", got, tt.want, tt.wantStage)
			}
			if tt.wantActOn != "" && got.ActOn != tt.wantActOn {
				t.Errorf("Decide().ActOn = %q, want %q", got.ActOn, tt.wantActOn)
			}
			if strings.TrimSpace(got.Reasoning) == "" {
				t.Errorf("Decide().Reasoning is empty — judges must always see the why")
			}
		})
	}
}

func TestStageForWarmCap(t *testing.T) {
	if got := stageFor(30, true); got != Stage1Polite {
		t.Errorf("stageFor(30, warm) = %v, want stage1_polite", got)
	}
	if got := stageFor(30, false); got != Stage4Escalation {
		t.Errorf("stageFor(30, cold) = %v, want stage4_escalation", got)
	}
}

func TestFullAutoRule(t *testing.T) {
	// Stage 3 (final notice) always needs a human, even in full-auto — the
	// policy itself never decides to auto-send it; that guard lives in tick,
	// but the stage math that feeds it must be stable.
	in := Input{Today: "2026-09-20", DueOn: "2026-09-05", Status: "chasing", LastSent: Stage2Firm}
	if d := Decide(in); d.Stage != Stage3Final {
		t.Fatalf("final stage math unstable: %+v", d)
	}
}
