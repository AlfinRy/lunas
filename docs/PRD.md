# PRD — Product Requirements Document

**Product:** Lunas — the AI collections agent that gets you paid
**Source of truth for scope.** Anything not in this document is a non-goal for the hackathon build.
**Build window:** Aug 21 – Sep 15, 2026 (25 days). Priorities: MoSCoW.

---

## 1. Vision

Freelancers shouldn't have to choose between getting paid and keeping good client relationships. Lunas is an autonomous collections agent: patient, polite, tireless, and firm exactly when it needs to be. The user adds an invoice once; Lunas handles everything until the moment it flips to **Lunas** — paid.

## 2. Personas

| | Rani — Freelance designer | Dimas — 4-person agency ops lead |
|---|---|---|
| Invoices/mo | 5–12 | 30–60 |
| Pain | Feels awkward chasing; postpones; gets paid 3 weeks late | Cash-flow gaps; chases across 12 clients in spreadsheets |
| Mode used | Approve-each-chase (trust building) | Full-auto after week 1 |
| Success | Money arrives without her writing a single awkward email | DSO visible and shrinking |

## 3. North-Star User Journey (this is the demo script)

1. **Onboard** → Rani lands on dashboard with seeded demo data (2 paid, 1 due soon, 2 overdue).
2. **Add invoice** → INV-0042, $2,400, for "Meridian Coffee Co.", due in 3 days. One modal, 5 fields.
3. **Agent plans** → Agent inbox shows a *plan card*: "INV-0042 due in 3 days. Meridian historically pays 9 days late. Schedule: soft reminder on due date, firm follow-up at +7. Tone: gentle — good relationship."
4. **Time passes** (simulator) → Due date hits. Agent drafts the reminder; Rani taps **Approve & send** (or auto-send in full-auto mode).
5. **No reply** (simulator advances) → Agent escalates: firmer email at +7, references original invoice + payment link.
6. **Payment arrives** → Rani pastes the "payment received" bank notification. Matcher links it to INV-0042 → invoice flips to **Lunas** 🎉, chasing stops instantly, recovered-revenue counter ticks up.
7. **Outcome visible** → Dashboard: $2,400 recovered, DSO down 6 days, "3 hours saved this month."

## 4. Feature Requirements (MoSCoW)

### F1 — Aging dashboard `P0`
The home screen. Answers "how much am I owed and what's Lunas doing about it?"

| ID | Requirement | Acceptance criteria |
|---|---|---|
| F1.1 | KPI cards: Outstanding total, Overdue total, DSO, Recovered (lifetime) | Values render from DB in < 500 ms; money uses tabular numerals |
| F1.2 | Invoice table: client, number, amount, due date, status (paid/due soon/overdue), agent state (idle/scheduled/chasing/paused) | Sortable by due date & amount; status pill colors per design tokens |
| F1.3 | Recovered-revenue counter with count-up animation on new settlement | Animation ≤ 800 ms, honors `prefers-reduced-motion` |
| F1.4 | DSO sparkline (30-day window) | Renders from derived data; no third-party chart lib bloat |

### F2 — Invoice & client management `P0`
| ID | Requirement | Acceptance criteria |
|---|---|---|
| F2.1 | Create/edit invoice: client, number, amount, currency, issue date, due date, notes | Validation: due ≥ issue; number unique per client; inline errors per UX-writing doc |
| F2.2 | Client records with email, payment terms, relationship note | Client payment-score displays after 2+ settled invoices |
| F2.3 | Status machine: `draft → scheduled → chasing → paid / written-off` | All transitions logged to activity timeline (F6) |

### F3 — Agent engine (the differentiator) `P0`
| ID | Requirement | Acceptance criteria |
|---|---|---|
| F3.1 | Chase policy engine (deterministic): given invoice age, payer history, relationship note → next action + stage + tone | Pure function; unit-tested; decisions auditable |
| F3.2 | 4-stage escalation ladder (see ux-writing.md): reminder → polite overdue → firm → final notice | Stage templates constrain LLM output; stage 4 always requires human send |
| F3.3 | LLM draft generation: policy engine passes structured context (stage, tone, invoice, payer history) → provider-agnostic adapter → draft | Graceful fallback to deterministic template if no API key or provider error; never blocks UI |
| F3.4 | Plan cards in Agent Inbox: what/when/why for each invoice | Every card shows the *reasoning* ("pays 9 days late on average") — judge-visible intelligence |
| F3.5 | Modes: **Approve each** (default) / **Full-auto** per client or global | Mode switch is one tap; current mode always visible in header |

### F4 — Email sending & outbox `P0` (demo-grade)
| ID | Requirement | Acceptance criteria |
|---|---|---|
| F4.1 | Mock SMTP provider (dev/demo): sends into visual **Outbox** view styled as a real mail client | Shows To/Subject/Body/Status(Sent) with realistic timestamps |
| F4.2 | Resend adapter behind the same interface | Swap via env var; not exercised in demo video |
| F4.3 | "Approve & send" flow from draft → outbox in one tap | Optimistic UI; error toast if provider fails |

### F5 — Payment reconciliation matcher `P0`
The "agent knows when to stop" magic.

| ID | Requirement | Acceptance criteria |
|---|---|---|
| F5.1 | Paste payment text (bank notification / "paid" email) → parser extracts amount, date, payer name/ID | Regex + fuzzy rules; works on 5 seeded formats (see test fixtures) |
| F5.2 | Matcher links payment → invoice(s): exact amount, amount + payer fuzzy match, partial payment → keeps chasing remainder | Confidence shown (High/Medium); low-confidence asks user to confirm |
| F5.3 | On match: status → paid, chase stops, counter updates, celebration state | Invoice row animates to green once (no loop) |

### F6 — Activity timeline `P1`
Per-invoice & global feed: every agent decision, draft, send, payment, pause — timestamped, human-readable ("Lunas sent a firm follow-up — Meridian is 9 days overdue").

### F7 — Time simulator `P0` (demo-critical)
| ID | Requirement | Acceptance criteria |
|---|---|---|
| F7.1 | "Advance time" control: +1d / +7d / jump-to-date | Recomputes all schedules, due statuses, and triggers agent actions as of simulated date |
| F7.2 | Simulated clock visible in header ("Demo time: Sep 2") | Never confusable with real time; clearly labeled |

### F8 — Payer payment-score `P1`
After ≥ 2 settled invoices: avg days late + reliability badge (Pays on time / Usually late / Chronically late). Feeds policy engine (F3.1) and plan-card reasoning.

### F9 — Settings `P1`
Sender identity (name/email), default payment terms, global mode toggle, escalation timing overrides, danger zone (reset demo data).

### F10 — Landing page `P1`
Above-the-fold: product name, one-line promise, hero showing the "Lunas moment" (invoice flipping to paid), waitlist CTA. Serif display type per design-system.

## 5. Non-Goals (explicitly out of scope)

- Real OAuth/accounting-suite integrations (QuickBooks/Xero import), real bank APIs
- Multi-user auth, teams, roles (single-workspace demo with seeded owner)
- Multi-currency FX conversion (store + display currency, no conversion)
- Mobile app (responsive web only)
- Legal-document generation beyond the stage-4 template

## 6. Non-Functional Requirements

| Area | Requirement |
|---|---|
| Performance | Dashboard first paint < 1.5 s on seeded data; interactions < 100 ms perceived |
| Accessibility | WCAG 2.1 AA for contrast & focus; full keyboard operability for core loop (add invoice → approve → reconcile) |
| Resilience | No API key configured → template mode banner, product fully usable (judges can run repo without keys) |
| Privacy | No real PII; seeded synthetic clients; LLM calls contain synthetic data only |
| Testability | Policy engine (F3.1) and matcher (F5.1) unit-tested; core loop E2E-tested (Playwright) |
| Stack | **Go 1.24 API** (Echo, sqlc, modernc SQLite) + **Vite SPA** (React 19, TanStack Router, TanStack Query, Tailwind v4, `motion` v12) · **[oa-design](https://github.com/OpenLabs-so/oa-design) design language** (MIT, vendored — jade-ink adaptation, see docs/design-system.md) · OpenAPI contract with dual codegen (openapi-typescript FE / oapi-codegen BE) · single-binary deploy via `go:embed` · provider-agnostic LLM adapter |

## 7. Data Model (summary)

```
Client      id, name, email, payment_terms, relationship_note, payment_score?
Invoice     id, client_id, number, amount (int cents), currency, issued_at,
            due_at, status, notes, created_at
ChasePolicy id, invoice_id, mode (approve|auto), current_stage (0-4),
            next_action_at, tone_override?
EmailDraft  id, invoice_id, stage, subject, body, status (draft|approved|sent|skipped),
            created_at, sent_at
Payment     id, invoice_id?, amount, paid_at, source (paste|csv|manual),
            confidence, raw_text
Activity    id, invoice_id?, type, message, created_at
Settings    singleton: sender_name, sender_email, default_terms, global_mode, sim_clock
```

## 8. Judging-Criteria Mapping (how the build scores)

| Criterion (weight) | Where we earn it |
|---|---|
| Technical 25% | Policy engine (pure, tested), reconciliation matcher, provider-agnostic agent adapter, clean data model |
| Problem/Impact 25% | Quantified dashboard (recovered $, DSO), universal pain, BRD metrics |
| Innovation 20% | Decision-layer agent (not text gen), payer-aware timing, chase-stop precision |
| UX/Design 15% | This design system; the "Lunas moment"; empty/error states done right |
| Presentation 15% | Demo script §3 + 3-min video storyboard (docs/demo-video.md, to be created in build window) |

## 9. 25-Day Build Plan (Aug 21 – Sep 15)

| Week | Dates | Deliverables |
|---|---|---|
| W1 Foundation | Aug 21–27 | Repo + monorepo layout, OpenAPI contract & dual codegen, Go scaffold (Echo/sqlc) + Vite scaffold, tokens/theme, data model + seed, invoice CRUD (F2), dashboard skeleton (F1) |
| W2 Agent core | Aug 28–Sep 3 | Policy engine + ladder (F3.1–3.3), drafts + outbox (F4), time simulator (F7), plan cards (F3.4) |
| W3 Closure loop | Sep 4–10 | Reconciliation matcher (F5), payment-score (F8), timeline (F6), settings (F9), landing (F10) |
| W4 Polish & submit | Sep 11–15 | E2E tests, a11y pass, demo data narrative, video script + record, Devpost submission **(buffer: 2 days)** |

Milestone gates: end of W2 = complete chase loop demoable; end of W3 = feature freeze; Sep 13 = code freeze, submission draft complete.
