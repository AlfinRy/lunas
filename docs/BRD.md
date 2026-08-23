# BRD — Business Requirements Document

**Product:** Lunas — the AI collections agent that gets you paid
**Hackathon:** AI Builders Hackathon 2026 (Devpost) — Best SaaS Product, $4,000
**Status:** Planning (pre-build window; all submission artifacts created after Aug 21, 2026)
**Owners:** Team Lunas

---

## 1. Executive Summary

Lunas is an AI agent that autonomously chases unpaid invoices on behalf of freelancers and small agencies. It watches every outstanding invoice, decides **when** to follow up and **how firmly**, drafts the email in the payer's appropriate tone, stops the instant payment is detected, and learns each client's payment behavior over time.

The product turns the most hated, most procrastinated chore of independent work — asking for your own money — into a set-and-forget system with a single, emotionally satisfying outcome: the invoice marked **Lunas** (Indonesian: *paid, settled*).

> One-line pitch for judges: **"Lunas is the AI agent that gets you paid."**

## 2. Problem Statement

- Late and non-payment is endemic among independent workers. Widely-cited industry surveys report that a majority of freelancers experience late payment at least once a year, and a substantial share write off invoices entirely. *(Exact figures to be verified against Freelancers Union / QuickBooks studies before submission — see §10 References.)*
- Chasing money is emotionally taxing: it feels awkward, damages client relationships when done poorly, and is therefore postponed — which directly lengthens payment delays.
- Existing tools (accounting suites, reminder schedulers) send **dumb mail-merge reminders**: fixed schedule, fixed tone, no awareness of payer history, no stopping logic beyond a manual toggle.
- The result is a measurable business pain: **cash-flow gaps, unpredictable DSO (Days Sales Outstanding), and hours of low-value admin per week.**

## 3. Business Objectives

| # | Objective | Measure |
|---|---|---|
| BO-1 | Win **Best SaaS Product** at AI Builders Hackathon 2026 | Judge decision, Sep 2026 |
| BO-2 | Demonstrate a **working product**, not a wrapper | Live demo end-to-end: invoice → autonomous chase → payment → stop |
| BO-3 | Prove quantified impact in-demo | Dashboard shows recovered revenue counter, DSO trend, hours saved |
| BO-4 | Post-hackathon viability | Landing page waitlist + clear SaaS pricing story |

### Why this can win (competitive thesis)

1. **Judging weights fit.** Technical Implementation (25%) + Problem Solving & Impact (25%) = 50%. Lunas is an agentic product with visible reasoning (policy engine + LLM drafting) solving a problem every judge has personally felt (getting paid).
2. **Demo-ability.** The chase loop compresses beautifully into a 3-minute narrative with a built-in time simulator. Emotional hook: the "Lunas" moment — a red overdue invoice flips to green.
3. **"No AI wrappers" rule alignment.** The agent makes real *decisions* (timing, tone, escalation, reconciliation matching), not just text generation. Architecture shows a policy engine + tool-using agent, clearly beyond a ChatGPT wrapper.

## 4. Market & Segments

| Segment | Pain intensity | Willingness to pay | Demo relevance |
|---|---|---|---|
| Freelancers (design, dev, writing) | Very high | $10–20/mo | Primary persona for demo |
| Micro-agencies (2–10 people) | High — payroll depends on it | $30–50/mo | Secondary persona |
| Bookkeepers / fractional CFOs | High — they chase on clients' behalf | $50+/mo | Post-hackathon expansion |

**TAM framing (for pitch):** global freelance platform GDP is estimated in the hundreds of billions USD annually; even 1% leakage from late payment represents billions recoverable. *(Verify before submission.)*

## 5. Value Proposition & Differentiation

| | Reminder tools (QuickBooks, Bonsai) | Virtual assistants | **Lunas** |
|---|---|---|---|
| Decides when to follow up | Fixed schedule | Human decides | **Agent decides per payer history** |
| Adapts tone (gentle → firm) | No | Yes (expensive) | **Yes, automated, 4-stage escalation** |
| Stops when paid | Manual toggle | Manual | **Automatic via reconciliation matching** |
| Learns payer behavior | No | Slowly | **Payer payment-score after 2+ invoices** |
| Cost | $10–30/mo inside suite | $ hundreds/mo | **$19/mo standalone** |

**Core differentiator:** the *decision layer*, not the email generation. Any tool can draft a reminder; Lunas decides **who gets chased, when, how hard, and when to stop** — and shows its reasoning in an auditable timeline.

## 6. Business Model (post-hackathon roadmap)

- **Free:** 5 active invoices, manual approve mode.
- **Pro $19/mo:** unlimited invoices, full-auto chasing, reconciliation, payer scoring.
- **Agency $49/mo:** multi-client workspaces, team members, custom escalation policies.
- **Distribution:** product-led; integrations (Stripe payment links, QuickBooks/Xero import) as acquisition loops.

## 7. Success Metrics (KPIs)

**Hackathon-window metrics**
- End-to-end demo completes in ≤ 3 minutes with zero errors.
- Recovered-revenue counter and DSO chart render live during demo (quantified impact).
- Devpost submission includes: live URL, open-source repo, demo video, docs.

**Product-health metrics (design now, instrument later)**
- % of chased invoices paid within 14 days of first chase (target: beat baseline by 20%).
- Chase-stop precision: % of chases correctly stopped on payment detection (target > 98%).
- Hours saved per user per month (self-reported; target ≥ 3).
- P91 draft-approval latency < 15 s.

## 8. Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| LLM drafts off-brand or too aggressive emails | Medium | High (trust) | Constrained templates + tone ladder; human-in-the-loop default; full-auto opt-in |
| Fake-send mode looks unconvincing to judges | Medium | High | MailPit-style outbox UI showing real SMTP envelope; simulator with realistic timestamps |
| Scope creep kills polish | High | High | MoSCoW discipline (PRD §5); P0 frozen after week 1 |
| "Just a wrapper" perception | Medium | High | Lead with policy engine + reconciliation matcher (deterministic code), LLM only for language |
| Privacy: real invoice data in demo | Low | Medium | Seeded synthetic data only; no real PII |

## 9. Compliance & Ethics Guardrails

- Chase emails must always identify the sender and include opt-out/relationship-preserving language (CAN-SPAM/GDPR-aware defaults).
- Escalation ladder never threatens legal action automatically beyond templated wording the user approved; stage 4 requires explicit user send.
- No dark patterns: pausing chasing is one tap, always visible.

## 10. References (verify before submission)

- Freelancers Union late-payment survey (the "71% of freelancers struggle to get paid" figure) — re-verify current edition.
- QuickBooks small-business invoice payment studies — pick one concrete stat for the pitch.
- Atrato / Chaser / Upflow (competitor pricing anchors).
