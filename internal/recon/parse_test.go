package recon

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantCents  int64
		wantCur    string
		wantDate   string
		wantInv    string
	}{
		{
			name:      "bank wire notification",
			text:      "Payment received: $2,400.00 from Meridian Coffee Co. on 2026-09-03 — Wire transfer ref INV-0041",
			wantCents: 240000, wantCur: "USD", wantDate: "2026-09-03", wantInv: "INV-0041",
		},
		{
			name:      "US long date, no invoice ref",
			text:      "You received a payment of USD 5,200 from Arkon Digital on Sep 3, 2026",
			wantCents: 520000, wantCur: "USD", wantDate: "2026-09-03",
		},
		{
			name:      "day-month date (Indonesian style English)",
			text:      "Received 1150.00 from Studio Kawi on 3 Sep 2026",
			wantCents: 115000, wantCur: "USD", wantDate: "2026-09-03",
		},
		{
			name:      "slash date MM/DD/YYYY",
			text:      "Payment of $980.00 completed 09/04/2026 — Bayu Ventures",
			wantCents: 98000, wantCur: "USD", wantDate: "2026-09-04",
		},
		{
			name:      "slash date flips to DD/MM when needed",
			text:      "Payment of $980.00 completed 23/08/2026",
			wantCents: 98000, wantCur: "USD", wantDate: "2026-08-23",
		},
		{
			name:      "no date falls back to today",
			text:      "We received your payment of $860. Thank you!",
			wantCents: 86000, wantCur: "USD", wantDate: "2026-08-30",
		},
		{
			name:      "IDR format uses dot thousands",
			text:      "Anda menerima transfer sebesar Rp 2.400.000 dari Studio Kawi pada 2026-09-03",
			wantCents: 240000000, wantCur: "IDR", wantDate: "2026-09-03",
		},
		{
			name:      "bare amount with label keyword",
			text:      "total received: 1,150.00 — thanks",
			wantCents: 115000, wantCur: "USD", wantDate: "2026-08-30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, ok := Parse(tt.text, "2026-08-30")
			if !ok {
				t.Fatalf("Parse(%q) failed, expected success", tt.text)
			}
			if p.AmountCents != tt.wantCents {
				t.Errorf("AmountCents = %d, want %d", p.AmountCents, tt.wantCents)
			}
			if p.Currency != tt.wantCur {
				t.Errorf("Currency = %q, want %q", p.Currency, tt.wantCur)
			}
			if p.PaidOn != tt.wantDate {
				t.Errorf("PaidOn = %q, want %q", p.PaidOn, tt.wantDate)
			}
			if p.InvoiceHint != tt.wantInv {
				t.Errorf("InvoiceHint = %q, want %q", p.InvoiceHint, tt.wantInv)
			}
		})
	}

	t.Run("no amount at all", func(t *testing.T) {
		if _, ok := Parse("Just a friendly note, no numbers here", "2026-08-30"); ok {
			t.Error("Parse should fail when no amount is present")
		}
	})
}

func TestMatch(t *testing.T) {
	open := []OpenInvoice{
		{InvoiceID: 3, InvoiceNo: "INV-0041", ClientID: 1, ClientName: "Meridian Coffee Co.", AmountCents: 240000, DueOn: "2026-08-21"},
		{InvoiceID: 6, InvoiceNo: "INV-0042", ClientID: 1, ClientName: "Meridian Coffee Co.", AmountCents: 86000, DueOn: "2026-09-01"},
		{InvoiceID: 4, InvoiceNo: "INV-0043", ClientID: 4, ClientName: "Arkon Digital", AmountCents: 520000, DueOn: "2026-08-19"},
	}

	t.Run("invoice reference wins with high confidence", func(t *testing.T) {
		p := Parsed{AmountCents: 240000, InvoiceHint: "INV-0041"}
		got := Match(p, nil, open, "received $2,400.00 ref INV-0041")
		if len(got) == 0 || got[0].InvoiceID != 3 || got[0].Confidence != "high" {
			t.Fatalf("got %+v", got)
		}
	})

	t.Run("amount + payer name → high", func(t *testing.T) {
		p := Parsed{AmountCents: 240000}
		got := Match(p, nil, open, "Payment received: $2,400.00 from Meridian Coffee Co.")
		if len(got) == 0 || got[0].InvoiceID != 3 || got[0].Confidence != "high" {
			t.Fatalf("got %+v", got)
		}
	})

	t.Run("amount only, unique → medium", func(t *testing.T) {
		p := Parsed{AmountCents: 520000}
		got := Match(p, nil, open, "some money arrived: $5,200.00")
		if len(got) == 0 || got[0].InvoiceID != 4 || got[0].Confidence != "medium" {
			t.Fatalf("got %+v", got)
		}
	})

	t.Run("partial payment scores as low-medium but listed", func(t *testing.T) {
		p := Parsed{AmountCents: 100000}
		got := Match(p, nil, open, "partial payment $1,000.00 from Arkon Digital")
		if len(got) == 0 || got[0].InvoiceID != 4 {
			t.Fatalf("got %+v", got)
		}
	})

	t.Run("irrelevant payment text yields no candidates", func(t *testing.T) {
		p := Parsed{AmountCents: 500}
		if got := Match(p, nil, open, "someone paid $5.00"); len(got) != 0 {
			t.Fatalf("expected no candidates, got %+v", got)
		}
	})
}
