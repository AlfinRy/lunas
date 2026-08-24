# Lunas Design System v2 — the OA Design language, wearing Lunas colors

**Base:** [oa-design](https://github.com/OpenLabs-so/oa-design) (MIT) — the shipped
design language of Open Analytics, vendored complete at `docs/design/oa-design/`
(skill: 12 component recipes, 7 guides, token CSS, type-checked sources).

**The look in one sentence:** white squircle surfaces resting on a quiet grey stage,
drawn in a single green-cast ink, moved by a single spring family, with one jade
accent spent in exactly one place per view — *the color of getting paid.*

v1 (OKLCH scale + concentric radius + CSS transitions) is superseded. What survives
from v1: the **measured contrast pairs** (§5) — our jade/ink values carry over as
text tints. Everything else follows OA construction.

---

## 1. Identity mapping — OA construction, Lunas ink

OA's rule: *derive, don't pick* — every neutral is `--ink` mixed at a stated
percentage. We keep that construction and pour Lunas' identity into the inputs:

| OA token | Open Analytics | **Lunas value** | Why |
|---|---|---|---|
| `--ink` | `#292929` | `#1b211d` (oklch 0.24 0.012 155) | Faint green cast; every derived neutral (borders, washes, inputs) inherits it |
| `--primary` | `#305dde` | `#017a4f` jade | White-on-primary = 5.39:1 ✓. The accent appears on the primary button and the primary chart line and **nowhere else** |
| `--primary` (dark) | `#296FF0` | `#179e69` | One notch brighter, per OA's dark rule |
| `--ring` | `#3ba6f1` | jade at 70% | Focus ring |
| `--chart-1` | `#296FF0` | `#179e69` | DSO sparkline line |
| `--background` | `#f6f6f6` | `oklch(0.985 0.003 140)` #f9fbf9 | The grey-green stage between plates |
| `--card`, `--popover` | white | white | unchanged |
| `--secondary` | `#e9e9e9` | `color-mix(ink 8%, #fff)` | flat grey, no border |
| semantic text tints | red/emerald/amber 700-row (light), 400-row (dark) | **our measured set** — success `#1d7d3e` / warning `#915c08` / destructive `#b32e2a`; dark: `#71d790` / `#f0ba59` / `#f87e77` | §5 verified |

Install `docs/design/oa-design/_root.css` as the base block (with the substitutions
above), then mirror under Tailwind v4 `@theme inline` → `--color-*`. Components use
semantic utilities only; raw values never appear in components.

**Status rendering (OA rule 7/8, replaces v1 tinted pills):** semantic colors are
*text tints and small dots, never fills*. A status is a neutral-surface pill with a
6px dot + tinted text: `● Paid`, `● Overdue`, `● Chasing`. Calmer tables, and the
jade accent keeps its full power for the moments that matter.

## 2. Typography — Inter Tight, ceiling at 500

- **Inter Tight** (variable, woff2, self-hosted in repo) 300–500. **No semibold,
  no bold anywhere.** Hierarchy = size, color (`foreground/80` vs `muted-foreground`),
  and spacing. KPI numbers: Inter Tight 500, `tabular-nums`, large size.
- **Geist Mono** 400–500 for invoice numbers (`INV-0042`), emails, amounts where
  tabular isn't enough — `code, kbd, samp, pre` at base layer.
- Headings `font-medium tracking-tight`; `text-wrap: balance` on marketing headlines.
- **`tabular-nums` on every monetary value, date, counter, DSO figure** (unchanged
  from v1 — non-negotiable).
- Inputs render at 16px on mobile (`text-base sm:text-sm`).
- Scale (semantic names, from OA): `text-display` 48–56 landing only · `text-h1` 28 ·
  `text-h2` 20 · `text-h3` 16 (500) · `text-body` 15 · `text-sm` 14 · `text-xs` 13 ·
  `text-label` 12 +0.06em uppercase eyebrows.

Fraunces (v1 display serif) is **dropped** — the calm, dense, native OA voice is the
better identity for a collections product.

## 3. Surfaces — squircles and pills

- **Squircle cards** (continuous-curvature corners) for every surface: KPI cards,
  table plates, panels, the outbox message card. Radius scale from `--radius: 10px`
  (`calc()` multipliers, see `_tokens.md`) — never a bespoke radius.
- **Pills for actions**: every clickable that is not a card is a pill
  (`rounded-full`). A `rounded-lg` rectangle on a control is the smell of a foreign
  component.
- **Two layers, and the gap is the page**: white plates on the grey-green stage; the
  page background is the only divider. No `<hr>`, no horizontal rules.
- Shadows: only where OA puts them (popovers, modals). Cards sit on the stage by
  contrast of layer, not by shadow.
- Images/desk content: 1px pure-black/white outline at 10% opacity (v1 rule kept).

## 4. Motion — one spring family, nothing over 200ms

`motion` (v12) is a frontend dependency. **The seven springs, verbatim** — do not
invent an eighth:

| Spring | Value | Lunas usage |
|---|---|---|
| PANEL | 550/38 | dropdowns, mode switcher, invoice row menus |
| LAYOUT | 550/40 | tab-bar traveling highlight, measured-height dialogs |
| POP | 400/26 | add-invoice modal, reconciliation dialog entrance |
| POP_EXIT | 380/28 | dialog exits |
| BANNER | 400/30 | demo-time floating pill, template-mode banner |
| FLICK | 900/50 | **the Lunas moment** — status dot/icon swap, chevrons |
| CHART | 300/28 | DSO sparkline tooltip/crosshair |

- Micro fades 0.10–0.18s easeOut; **nothing in app chrome tweens past 0.2s**.
- **Chrome never waits:** layout and titles render instantly; only data swaps from a
  pixel-matched skeleton, arriving **by blur, not by pop** (skeleton recipe).
- Signature moves per `_motion.md`: measured-height choreography (multi-step
  reconciliation), the sliding highlight (invoice status filters), the label mask.
- **The Lunas moment (v2):** payment matched → status dot + text tint swap on FLICK
  (danger→jade), row wash fades, recovered counter ticks (tabular-nums, LAYOUT),
  success toast pulses once per OA toast recipe. Total ≤ 1.2s, honors
  `prefers-reduced-motion` (quality floor rule 10).
- Enter staggers only on first dashboard load (~80ms beats, reveal recipe), once.

## 5. Contrast verification (carried from v1, still enforced)

All pairs re-measured WCAG via the OKLCH→sRGB script:

| Pair (light) | Ratio | | Pair (dark) | Ratio |
|---|---|---|---|---|
| ink `#1b211d` / stage | 15.7:1 AAA | | text / bg | 16.0:1 AAA |
| `#545d57` secondary text | 6.6:1 AA | | secondary | 8.4:1 AAA |
| `#646b66` muted text | 5.3:1 AA | | muted | 5.1:1 AA |
| white / jade `#017a4f` | 5.4:1 AA | | jade-dark `#179e69` text | 9.1:1 AAA |
| success `#1d7d3e` text | 5.0:1 AA | | `#71d790` | 9.3:1 AAA |
| warning `#915c08` text | 5.0:1 AA | | `#f0ba59` | 9.4:1 AAA |
| destructive `#b32e2a` text | 5.5:1 AA | | `#f87e77` | 6.5:1 AA |

Re-run `scripts/colorcheck.py` whenever a token changes.

## 6. Layout (delta from v1)

Unchanged: app shell (sidebar 200px, content 1200px, sticky header), spacing scale
4–64, 2× grouping rule, alignment edges, importance ordering (overdue first),
content-driven breakpoints (1280 / 1024 icon-rail / 768 stacked), progressive
disclosure with peek, logical properties, 16px mobile inputs.

OA additions: settings screens follow `_layout.md` plate anatomy; landing page per
`_landing.md` (hero, reveals in 80ms beats, header-morph glass pill).

## 7. Component mapping — Lunas features → OA recipes

| Lunas surface | OA recipe |
|---|---|
| KPI cards, invoice table plate, client cards, outbox message | `01-squircle-card` |
| Primary/secondary/destructive actions, honest loading | `02-button` |
| Mode switcher (Approve / Full-auto), filters, selects | `03-dropdown` + `04-tab-bar` |
| Add/edit invoice | `05-modal` |
| Reconciliation confirm (match → link → settled) | `06-multi-step-dialog` |
| Dashboard & table loads | `07-skeleton` (pixel-matched, blur arrival) |
| Template-fallback banner ("AI provider unreachable…") | `08-notice-strip` |
| Demo-time pill + time controls | `09-floating-pill` |
| Lunas moment, chase sent, demo reset | `10-toast` |
| Landing header | `11-header-morph` |
| Landing sections | `12-reveal` |

## 8. Quality floor (OA rule 10, enforced in review)

Focus rings visible in both themes (`role` of `--ring`), `role="status"` on live
regions/toasts, `aria-hidden` on decorative glyphs, `prefers-reduced-motion`
respected, no horizontal page scroll, text selectable by default, full keyboard
operability of the core loop (add → approve → reconcile).

## 9. Attribution

Design language: [oa-design](https://github.com/OpenLabs-so/oa-design) by
Voprex Labs / Open Analytics, MIT — vendored at `docs/design/oa-design/` with
Lunas token substitutions documented in §1. Credited in README and the demo video.
