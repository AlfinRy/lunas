package api

import "testing"

func TestDaysBetween(t *testing.T) {
	cases := []struct {
		from, to string
		want     int
	}{
		{"2026-08-01", "2026-08-01", 0},
		{"2026-08-01", "2026-08-10", 9},
		{"2026-08-10", "2026-08-01", -9},
		{"2026-07-25", "2026-08-15", 21}, // Meridian seed case
		{"garbage", "2026-08-01", 0},
	}
	for _, c := range cases {
		if got := daysBetween(c.from, c.to); got != c.want {
			t.Errorf("daysBetween(%q, %q) = %d, want %d", c.from, c.to, got, c.want)
		}
	}
}

func TestValidEmail(t *testing.T) {
	valid := []string{"billing@meridian.co", "a@b.id", "x.y+tag@studio.kawi.jp"}
	invalid := []string{"", "nope", "nope@", "@nope.com", "a b@c.com"}
	for _, s := range valid {
		if !validEmail(s) {
			t.Errorf("validEmail(%q) = false, want true", s)
		}
	}
	for _, s := range invalid {
		if validEmail(s) {
			t.Errorf("validEmail(%q) = true, want false", s)
		}
	}
}

func TestValidCurrency(t *testing.T) {
	for _, s := range []string{"USD", "IDR", "EUR"} {
		if !validCurrency(s) {
			t.Errorf("validCurrency(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"", "usd", "US", "USDD", "U1D"} {
		if validCurrency(s) {
			t.Errorf("validCurrency(%q) = true, want false", s)
		}
	}
}

// scoreCalc drives the payer payment-score (PRD F8) and the agent's reasoning
// lines, so its thresholds are pinned by tests.
func TestScoreReliability(t *testing.T) {
	cases := []struct {
		late  []float64 // days late per settled invoice
		want  Reliability
		avg   float64
	}{
		{[]float64{-1, 1}, PaysOnTime, 0},    // early then late → on time
		{[]float64{0, 2}, PaysOnTime, 1},       // avg exactly 1 → on time (small-noise tolerance)
		{[]float64{8, 10}, UsuallyLate, 9},    // Meridian's story: 9 days avg
		{[]float64{20, 30}, ChronicallyLate, 25},
		{[]float64{9}, UsuallyLate, 9},       // Meridian single-invoice story
	}
	for _, c := range cases {
		s := &scoreCalc{}
		for _, d := range c.late {
			s.add(d)
		}
		if got := s.reliability(); got != c.want {
			t.Errorf("reliability(%v) = %v, want %v", c.late, got, c.want)
		}
		if got := s.avgLate(); got != c.avg {
			t.Errorf("avgLate(%v) = %v, want %v", c.late, got, c.avg)
		}
	}
}
