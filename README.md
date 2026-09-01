# Lunas — Autonomous Invoice-Chasing Agent

> "Lunas" = Bahasa Indonesia for **paid / settled**.
> Lunas is the AI collections agent that gets your invoices to *lunas* — politely, persistently, automatically.

Dibangun untuk **AI Builders Hackathon 2026** (Devpost) — kategori Best SaaS Product.
Submission deadline: **15 September 2026, 23:00 EDT**.

## The Problem
Freelancers & small agencies lose **50+ hours a year** chasing unpaid invoices, and poor follow-up is the #1 reason invoices get paid late. Existing tools are dumb mail-merge reminders — wrong tone, wrong timing, awkward for client relationships.

## The Product
An AI agent that autonomously manages accounts receivable:

1. **Aging dashboard** — every outstanding invoice, DSO metric, recovered-revenue counter
2. **Agent brain** — per client & invoice, the agent decides *when* to follow up, *what tone* to use (gentle → firm → final notice), and drafts the email referencing full context
3. **Human-in-the-loop or full-auto** — approve drafts with one tap, or let it run
4. **Payment reconciliation** — upload a bank CSV / paste a payment email → auto-matched, invoice closed, chasing stops instantly
5. **Payer scoring** — learns each client's payment behavior, predicts who pays late, schedules accordingly
6. **Time simulator** — fast-forward days for the demo video

## Stack (planned)

**Backend — Go 1.24:** Echo · sqlc + modernc.org/sqlite (CGO-free) · OpenAPI → oapi-codegen server stubs · provider-agnostic LLM adapter with deterministic template fallback · testify (table-driven tests for policy engine & payment matcher)

**Frontend — Vite + React 19:** TanStack Router (SPA) · TanStack Query · Tailwind v4 with OKLCH `@theme` tokens (see docs/design-system.md) · Lucide icons · Playwright E2E for the core loop

**Deploy:** `vite build` output embedded into the Go binary (`go:embed`) → one artifact, one URL, no CORS, no cold starts. Local dev: `air` + Vite proxy. Repo quickstart: `docker compose up`.

## Documentation

| Doc | What it contains |
|---|---|
| `docs/hackathon-research.md` | Why this hackathon was picked; judging criteria |
| `docs/BRD.md` | Business case, market, differentiators, risks, KPIs |
| `docs/PRD.md` | Scope: features (MoSCoW), acceptance criteria, data model, 25-day plan |
| `docs/design-system.md` | Theme: OKLCH color tokens (verified contrast), type scale, layout, motion |
| `docs/ux-writing.md` | Voice, microcopy inventory, and the 4-stage chase-email ladder |

## Quickstart

Requires Go ≥ 1.25 and [Bun](https://bun.sh). No database setup — SQLite is a file.

**Dev (hot reload):**

```sh
# Terminal 1 — API on :8080 (auto-seeds demo data on first run)
go run ./cmd/lunas

# Terminal 2 — SPA on :5173 (proxies /api)
cd web && bun install && bun run dev
```

**Single binary (production):**

```sh
cd web && bun run build
cd .. && go build -o lunas ./cmd/lunas && ./lunas   # SPA embedded, one port
```

**From your phone:** run the binary, then open `http://<laptop-ip>:8080` from any
device on the same network (allow the port through the firewall once).

**Reset the demo:** delete `lunas.db*` and restart — or Settings → Danger zone
in the app.

### Environment

| Variable | Default | Purpose |
|---|---|---|
| `LUNAS_PORT` | `8080` | Server port |
| `LUNAS_DB` | `lunas.db` | SQLite file path |
| `LLM_API_KEY` | *(empty)* | Any OpenAI-compatible key; empty = template mode |
| `LLM_BASE_URL` | `https://api.openai.com/v1` | Swap providers freely |
| `LLM_MODEL` | `gpt-4o-mini` | Model used for drafting |

No key needed: decisions (policy engine, matcher) are deterministic Go; the LLM
only polishes email copy and falls back to the canonical templates on any error.

## Credits

- Design language: [oa-design](https://github.com/OpenLabs-so/oa-design) by Voprex Labs / Open Analytics — MIT

## Status

🚧 **Build in progress** — hackathon window Aug 21 – Sep 15, 2026 (day 3/25).
Follow the plan in `docs/PRD.md` §9 (W1: foundation → W4: polish & submit).
