# Lunas UX Writing & Voice Guide

Interface copy, microcopy, and the chase-email ladder. Built on the `better-writing`
principles: clear beats clever, one voice with flexible tone, errors are instructions,
empty states point forward.

Language: English (submission-facing). Currency examples seeded in USD.

---

## 1. Voice

**One voice:** calm, competent, quietly on your side — a great account manager who is
polite to your clients and relentless about your money.

| Context | Tone | Example |
|---|---|---|
| Success, "Lunas moment", onboarding | Warm, light | "Invoice settled. Chasing stopped." |
| Routine (tables, settings, nav) | Neutral, minimal | "Due in 3 days" |
| Errors, destructive actions | Calm, plain, zero playfulness | "Unable to send draft. Check the sender email in Settings, then try again." |
| Data loss / danger zone | Serious, explicit | "Reset demo data. This deletes all invoices, drafts, and payment history." |

Address the reader as **you**. Never "the user". No "we" in errors ("We're having trouble…" → "Unable to…").

## 2. Terminology (use exactly these words)

| Term | Meaning | Never |
|---|---|---|
| **Invoice** | the bill being chased | "bill", "document" |
| **Chase** | one agent follow-up action | "reminder" alone (reminder is stage 0's name) |
| **Payer** | the client who owes | "customer" (we serve the user, not the payer) |
| **Settled / Lunas** | invoice fully paid | "closed", "complete" |
| **Recovered** | money in after a chase | "won", "collected" |
| **Agent** | the Lunas engine | "bot", "AI assistant" |

## 3. Capitalization & format policy

- **Sentence case everywhere**: buttons, headings, nav, toasts, pills. ("Save draft", not "Save Draft").
- Numbers via full templated strings with pluralization — never string concatenation around variables.
- Digits tabular (`tabular-nums`); dates `Sep 2`; relative time "in 3 days", "9 days overdue".

## 4. Buttons (verb-first, consequence-repeating)

| Button | Context |
|---|---|
| **Add invoice** | Primary, dashboard/empty state |
| **Approve & send** | Draft review — names both steps |
| **Skip this draft** | Draft review (secondary) |
| **Edit draft** | Draft review |
| **Mark as settled** | Invoice detail (manual reconciliation confirm) |
| **Pause chasing** | Invoice row/detail — toggle, labels the result |
| **Resume chasing** | The paused state |
| **Link payment** | Reconciliation confirm (medium confidence) |
| **Send final notice** | Stage 4 — repeats the consequence instead of "Yes" |
| **Cancel** | All secondary dialogs |
| **Reset demo data** | Danger zone — explicit, no "OK" |

Flow vocabulary: **Get started → Continue → Done** (one set, never alternated synonyms).

## 5. Errors = instructions, next to the field

| Bad | Good (ship this) |
|---|---|
| "Invalid date" | "Pick a due date after the issue date." |
| "Amount required" | "Enter an amount, like 2400.00." |
| "Client email invalid" | "Enter the payer's email, like billing@meridian.co." |
| "Oops! Something went wrong." | "Unable to send draft. Check the sender email in Settings, then try again." |
| "LLM error" | "Couldn't reach the AI provider. Using the standard template instead." *(product keeps working — resilience is the message)* |

Hints shown **before** the mistake (placeholder/format hints), phrased positively ("Use only letters…"). **Report all field problems at once** — fixing them one at a time is a needlessly slow form (add-invoice modal validates every field on submit).

## 6. Empty states (orientation + one action)

- **Dashboard, no invoices:** heading "No invoices yet" · body "Add an invoice and Lunas will start planning the first chase." · button "Add invoice".
- **Agent inbox, all clear:** "Nothing needs your approval" · "Lunas is watching 5 invoices. You'll see drafts here the moment one needs chasing."
- **Outbox:** "No emails sent yet" · "Approved chases land here with their delivery status."
- **Search/filter no results:** "No invoices match 'meridian'." · action "Clear filters".
- **Reconciliation, no match:** "No invoice matches this payment" · "Check the amount, or link it manually." · buttons "Link manually" / "Dismiss".

## 7. Placeholders are examples, never labels

Every field keeps a visible label; placeholders show format only: `INV-0042`, `billing@meridian.co`, `2400.00`, `DD MMM YYYY`.

## 8. Settings — label the ON state

- ✅ "Send chases automatically" (not "Don't ask before sending")
- ✅ "Copy me on every chase"
- ✅ "Escalate to final notice after 14 days overdue"
- Links describe destination: "Edit payment terms", never "Click here" / bare "Learn more".

## 8.5 Never claim what the system hasn't confirmed

A pasted bank notification is *input*, not confirmation. The reconciliation flow says "Confirming payment…" until the matcher links it and the user confirms (or confidence is High); only then: "INV-0042 settled — $2,400 recovered." Waiting states say what is awaited and how long it usually takes — only durations we've measured ("This usually takes a few seconds." for the LLM draft).

## 9. Toasts

| Event | Toast |
|---|---|
| Draft approved | "Chase sent to Meridian Coffee Co." |
| Payment matched | "INV-0042 settled — $2,400 recovered." |
| Chase paused | "Chasing paused for INV-0042." |
| Provider fallback | "AI provider unreachable — sent the standard template." |
| Demo data reset | "Demo data reset." |

## 10. Agent status lines (timeline & inbox)

Status lines are how judges *see* the agent think — always show the reasoning:

- "Scheduling a soft reminder for Sep 2 — Meridian usually pays 9 days late."
- "Drafting a firm follow-up — invoice is 7 days overdue, no reply to the first chase."
- "Holding — payer replied promising payment Friday. Re-checking Sep 5."
- "Chase stopped — payment of $2,400 matched this morning."
- "Skipping weekend — chases send Tuesday–Thursday mornings, payer's local time."

## 11. The Chase Ladder (product copy, 4 stages)

The escalation ladder is the product's soul. Each stage constrains the LLM draft; these
are the canonical templates (also the no-API fallback).

**Rules across all stages:** state the facts (invoice number, amount, days overdue) · one
clear next step (the payment link) · never guilt-trip, never threaten before stage 3 ·
stage 4 is never auto-sent · sign as the user's billing tool, transparent that it's automated:
*"This is an automated reminder sent by Lunas on behalf of {sender_name}."*

### Stage 0 — Reminder (due date, or 1 day before if payer is "usually late")
> **Subject:** Invoice {number} — due {date}
>
> Hi {payer_first_name},
>
> A quick note that invoice **{number}** for **{amount}** is due **{date}**.
>
> Pay here: {payment_link}
>
> If it's already scheduled, feel free to ignore this note.
>
> {sender_name}

### Stage 1 — Polite overdue (+1–3 days)
> **Subject:** Invoice {number} — {amount}, now {days} days past due
>
> Hi {payer_first_name},
>
> Marking this as unpaid on our side — invoice **{number}** for **{amount}** was due **{date}**.
>
> If anything's blocking on your end (approval, PO, paperwork), reply here and {sender_name} will sort it out today. Otherwise, the fastest way to settle is this link: {payment_link}
>
> Thanks,
> {sender_name}

### Stage 2 — Firm follow-up (+7 days, no reply)
> **Subject:** Second notice: invoice {number} is {days} days overdue
>
> Hi {payer_first_name},
>
> Following up again on invoice **{number}** — **{amount}**, now **{days} days past due**. We haven't received a reply to the previous note.
>
> Please confirm a payment date this week, or settle directly: {payment_link}
>
> If there's a dispute or question about the work, flag it now so we can resolve it quickly.
>
> {sender_name}

### Stage 3 — Final notice (+14 days) — *requires user send*
> **Subject:** Final notice — invoice {number}, {days} days overdue
>
> Hi {payer_first_name},
>
> This is a final notice for invoice **{number}** (**{amount}**, due **{date}**). To date we've sent {n_chases} reminders without payment or an agreed payment date.
>
> Please settle within **5 business days**: {payment_link}
>
> If payment isn't received or a plan agreed by {deadline}, this account moves to formal collections, and {late_fee_clause}.
>
> We'd rather resolve it directly. Reply with a payment date and this is closed.
>
> {sender_name}

### Stage 4 — Escalation packet (post-final, +21 days) — *user-sent only*
Not an email: Lunas compiles the **collections dossier** — full chase timeline, delivered-email proof, amount breakdown, and a draft referral letter. Copy: "Everything your collections process needs, in one file."

**Tone escalation summary** for plan cards: `gentle → polite → firm → final`. Relationship
note "long-term client, keep warm" pins the ladder at stages 0–1 and adds a human check-in
line; it never blocks collection silently.

## 12. Landing page copy (draft)

- **Hero:** (Fraunces display) "Get paid. Without the awkward emails." · sub: "Lunas is the AI collections agent that chases your invoices — politely, persistently, and only until they're paid."
- **CTA:** "Try the live demo" (secondary: "Read how it works")
- **Proof strip:** "Invoice settled. Chasing stopped." — the product's favorite sentence.
