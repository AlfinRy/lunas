# Lunas Design System — Theme & Style Guide

Single source of truth for the visual language. Built on the principles of the
`better-colors`, `better-typography`, `better-layout`, and `better-ui` skills.
Every value below was **measured**, not guessed (see §1.5 verification).

Stack context: Vite + React 19 + TanStack Router SPA · Tailwind v4 (`@theme` tokens) · Lucide icons. No SSR — the dashboard is the product; the SPA is embedded into the Go binary at deploy time.

---

## 1. Color

### 1.1 Principles

- **One color, one meaning.** Jade green = the brand *and* the "paid/settled" state — deliberately unified, because getting paid **is** the brand (the "Lunas moment"). Amber = attention/aging. Red = overdue/destructive. Blue = links & information. Never borrow a status color for decoration.
- **Only the primary action is filled.** One colored control per view (the primary button); secondaries are neutral. Status pills are tints, never saturated fills competing with actions.
- **Dark mode is tuned, not reversed.** Chroma and lightness re-derived; every fg/bg pair re-measured.

### 1.2 Semantic tokens — Light (default)

| Token | OKLCH | sRGB fallback | Role |
|---|---|---|---|
| `--bg` | `oklch(0.985 0.003 140)` | `#f9fbf9` | App background |
| `--surface` | `oklch(1 0 0)` | `#ffffff` | Cards, tables, modals |
| `--text-primary` | `oklch(0.24 0.012 155)` | `#1b211d` | Headings, body |
| `--text-secondary` | `oklch(0.47 0.015 155)` | `#545d57` | Supporting text |
| `--text-muted` | `oklch(0.52 0.012 155)` | `#646b66` | Captions, timestamps |
| `--border` | `oklch(0.9 0.008 150)` | `#dadfdb` | Dividers, input borders |
| `--primary-600` | `oklch(0.51 0.115 160)` | `#017a4f` | Primary button bg |
| `--primary-700` | `oklch(0.44 0.09 160)` | `#156141` | Primary hover; text on `primary-100` |
| `--primary-500` | `oklch(0.62 0.135 160)` | `#179e69` | Charts, accents (non-text) |
| `--primary-100` | `oklch(0.955 0.025 160)` | `#e3f6ea` | Paid pill tint, selected states |
| `--success-600` | `oklch(0.52 0.13 150)` | `#1d7d3e` | Success text/icons |
| `--warning-600` | `oklch(0.52 0.11 70)` | `#915c08` | "Due soon" text on tint |
| `--warning-100` | `oklch(0.96 0.03 85)` | `#fbf1dc` | "Due soon" pill tint |
| `--danger-600` | `oklch(0.51 0.17 27)` | `#b32e2a` | Overdue, destructive |
| `--danger-100` | `oklch(0.96 0.018 27)` | `#feeeeb` | Overdue pill tint |
| `--info-600` | `oklch(0.5 0.13 250)` | `#1666aa` | Links |

### 1.3 Semantic tokens — Dark

| Token | OKLCH | sRGB fallback |
|---|---|---|
| `--bg` | `oklch(0.185 0.008 160)` | `#101411` |
| `--surface` | `oklch(0.235 0.01 160)` | `#1a201c` |
| `--text-primary` | `oklch(0.95 0.005 155)` | `#ecefed` |
| `--text-secondary` | `oklch(0.75 0.012 155)` | `#a8b0ab` |
| `--text-muted` | `oklch(0.62 0.012 155)` | `#818883` |
| `--border` | `oklch(0.32 0.012 160)` | `#2e3531` |
| `--primary-400` | `oklch(0.79 0.13 160)` | `#65d49e` |
| `--primary-600` | `oklch(0.51 0.115 160)` | `#017a4f` |
| `--success-400` | `oklch(0.8 0.14 152)` | `#71d790` |
| `--warning-400` | `oklch(0.82 0.13 80)` | `#f0ba59` |
| `--danger-400` | `oklch(0.73 0.15 25)` | `#f87e77` |

Dark-mode status tints reuse `--surface` + a 12% alpha of the 400-level color; text uses the 400-level color (all pairs pass AA, §1.5).

### 1.4 Tailwind v4 wiring

```css
@theme inline {
  --color-bg: oklch(0.985 0.003 140);
  --color-surface: oklch(1 0 0);
  --color-primary: oklch(0.51 0.115 160);
  --color-primary-strong: oklch(0.44 0.09 160);
  --color-primary-soft: oklch(0.955 0.025 160);
  /* ...status tokens as above... */
}
```
Author with semantic utilities (`bg-surface`, `text-primary`); **raw palette values never appear in components.**

### 1.5 Contrast verification (WCAG, measured via OKLCH→sRGB script)

| Pair (light) | Ratio | Grade | Pair (dark) | Ratio | Grade |
|---|---|---|---|---|---|
| text-primary / bg | 15.75 | AAA | text-primary / bg | 16.05 | AAA |
| text-secondary / bg | 6.56 | AA | text-secondary / bg | 8.37 | AAA |
| text-muted / bg | 5.26 | AA | text-muted / bg | 5.12 | AA |
| white / primary-600 (btn) | 5.39 | AA | primary-400 / surface | 9.05 | AAA |
| white / primary-700 (hover) | 7.45 | AAA | success-400 / surface | 9.33 | AAA |
| primary-700 / primary-100 | 6.61 | AA | warning-400 / surface | 9.37 | AAA |
| warning-600 / warning-100 | 5.01 | AA | danger-400 / surface | 6.50 | AA |
| danger-600 / danger-100 | 5.56 | AA | | | |
| info-600 / bg (links) | 5.75 | AA | | | |

All palette steps verified in-gamut for sRGB at the listed chroma. Re-run the verification script (`scripts/colorcheck.py`, to be added in build window) whenever a token changes.

---

## 2. Typography

### 2.1 Families

| Role | Family | Weights | Notes |
|---|---|---|---|
| UI text | **Inter Variable** (woff2, self-hosted) | 400 / 500 / 600 | Body, tables, controls. Weights < 400 never used below 18px |
| Display | **Fraunces** (woff2) | 560–640 | Landing hero, section titles, the "Lunas" wordmark only |

Two families, paired for contrast (editorial serif vs neutral sans) — never two near-identical sans. Monospace not used; Inter's `tabular-nums` covers numeric alignment.

### 2.2 Scale (semantic names, not sizes)

| Token | Size / line-height | Usage |
|---|---|---|
| `text-display` | 56px / 1.1, letter-spacing −0.02em | Landing hero only |
| `text-h1` | 30px / 1.15, ls −0.01em | Page titles ("Invoices") |
| `text-h2` | 20px / 1.25 | Section headers, modal titles |
| `text-h3` | 16px / 1.4, weight 600 | Card group headers |
| `text-body` | 15px / 1.55 | Default body, table cells |
| `text-sm` | 14px / 1.5 | Secondary rows, inputs |
| `text-xs` | 13px / 1.5, ls +0.01em | Captions, timestamps, meta |
| `text-label` | 12px / 1.4, ls +0.06em, uppercase | Eyebrow labels ("OUTSTANDING") |

Rules: adjacent heading levels always descend; `text-wrap: balance` on headings; `text-wrap: pretty` on descriptions; inputs render at **16px on mobile** (`text-base sm:text-sm`) to prevent iOS zoom; `antialiased` set once at root.

### 2.3 Numbers & money (critical)

- **`font-variant-numeric: tabular-nums` on every monetary value, date, counter, and DSO figure.** Amounts change and re-sort; proportional digits would jitter the layout.
- Currency format: `$2,400.00` / `Rp 2.400.000` via `Intl.NumberFormat`; amount cells right-aligned, headers aligned to the data (`text-align: end`).
- Big KPI numbers: Inter 600, tabular-nums; count-up animation only on change, ≤ 800 ms, respects reduced motion.

### 2.4 Truncation & wrapping

Client names and subjects: single-line ellipsis + full value in tooltip/detail view. Email draft previews: 2-line `line-clamp`. No fixed-width text containers; test with the longest seeded client name ("Meridian Coffee Company Pty Ltd").

---

## 3. Layout

### 3.1 App shell

- **Left nav (200px):** Dashboard, Invoices, Clients, Agent inbox (badge with pending approvals), Outbox, Settings. Content max-width `1200px`, centered, `24px` inline padding (16px on mobile).
- **Header row:** page title left; demo-time pill + mode toggle (Approve / Full-auto) right. Sticky, floats above content (`backdrop-blur` on `--surface/85%`).
- Reading order = importance: KPIs → attention items (overdue first) → all invoices. Overdue rows are always the top table group.

### 3.2 Spacing scale

`4 · 8 · 12 · 16 · 24 · 32 · 48 · 64` (px). One step per level of subordination.

| Relationship | Gap |
|---|---|
| Within a control group (label→input, icon→label) | 8 |
| Within a card (rows, KPI + caption) | 8–12 |
| Between cards / table groups | 24 |
| Between page sections | 48 |

**Group with space, not lines:** borders only for structure (table row separators, input borders). Never a `<hr>` where `24px` gap communicates the same grouping. Inter-group gap ≥ 2× intra-group gap.

### 3.3 Alignment edges

All content in a card aligns to two shared edges (inline start + end); amounts column forms its own right-edge alignment zone. No stray indents — every edge traces to the 8px grid.

### 3.4 Responsive behavior (content-driven breakpoints)

| Range | Layout |
|---|---|
| ≥ 1280 | Full: nav + 4 KPI cards in one row + full table |
| 1024–1279 | KPIs 2×2; table keeps amount + status columns |
| 768–1023 | Nav collapses to icon rail (56px); table drops "issued" column |
| < 768 | Nav becomes top bar with menu; invoices render as stacked cards; **inputs stay 16px** |

Collapse late (keep the expanded layout as long as it genuinely fits). Core demo target is desktop 1440px; mobile must remain fully operable for the add→approve→reconcile loop.

### 3.5 Progressive disclosure

Draft detail (full email body) lives behind a chevron; the row shows subject + 2-line preview. The next collapsed section always peeks 16–32px past the fold — never a dead-flat scroll edge.

### 3.6 Growth & i18n readiness

No fixed widths on text containers; labels wrap. Buttons sized by content (`min-w` only where the demo needs a stable tap target). Logical properties everywhere (`ps-*`, `pe-*`, `text-start`) so an RTL mirror doesn't break alignment.

---

## 4. Surfaces, depth & components

### 4.1 Radius — concentric

| Element | Radius |
|---|---|
| Cards, modals | 16px |
| Buttons, inputs, pills | 10px |
| Nested chip inside a pill row | 6px (10 − padding) |

`outer = inner + padding` — verified on every nested pair (the classic "off" feeling comes from breaking this).

### 4.2 Elevation

```css
--shadow-card: 0 1px 2px oklch(0 0 0 / 0.05), 0 4px 12px oklch(0 0 0 / 0.06);
--shadow-pop:  0 2px 4px oklch(0 0 0 / 0.06), 0 12px 32px oklch(0 0 0 / 0.12);
```
Layered, transparent shadows for depth (cards, dropdowns, modals). **Borders for structure only**: table separators, input borders, focus rings, selected-state outline. Dark mode: shadows become near-invisible → depth carried by `--surface` vs `--bg` lightness gap + 1px borders.

### 4.3 Buttons

| Variant | Style |
|---|---|
| Primary | Fill `--primary-600`, white text, hover `--primary-700`, press `scale(0.96)` |
| Secondary | `--surface` fill, `--border` border, text `--text-primary` |
| Destructive | Fill `--danger-600` — only in confirmation dialogs |
| Ghost | Text-only, `--text-secondary`, hover bg `--bg` |

One primary per view. All labels verb-first (see ux-writing.md). Icon buttons carry `aria-label` + tooltip.

### 4.4 Status pills

Tinted bg (100-level) + 600-level text + 1px same-hue border: `Paid`, `Due soon`, `Overdue`, `Chasing`, `Scheduled`, `Paused`. Pill = `text-xs`, uppercase off (sentence case reads calmer), `white-space: nowrap`. Never filled-saturated — they must not compete with the primary action.

### 4.5 Icons (Lucide)

- Default: outline, stroke **1.5px** next to 400/500 text; **2px** next to 600 semibold labels; sizes 16 (inline), 20 (buttons/nav), 24 (empty states).
- One library only; `currentColor` for every state (hover/active/disabled via CSS `color` + `opacity`); outline default, filled variant marks the active nav item.
- Icon-next-to-text is optically inset (icon visually centered against cap-height, not bounding box).

### 4.6 Motion

| Interaction | Spec |
|---|---|
| State hovers | `color, background-color, border-color` 150ms `ease-out` — never `transition: all` |
| Button press | `scale(0.96)` via `active:`, transform 100ms |
| Dashboard first load | Cards/rows stagger-enter ~100ms apart, `translateY(4px)→0` + opacity, `ease-out`; **once only** |
| New draft appears | Same enter values; icon swap animates `scale 0.25→1`, `opacity 0→1`, `blur 4px→0` |
| Exits (row dismissed) | Subtle: `translateY(-4px)` + fade 120ms — softer than enter |
| "Lunas moment" | Invoice row: danger→success pill crossfade + check icon scale-in; confetti **rejected** — restraint |
| Reduced motion | `prefers-reduced-motion: reduce` → all transforms/fades become instant; color change remains (static cue always accompanies motion) |

No staggered motion on high-frequency interactions (typing, table sorting). No `will-change` except on the count-up counter if first-frame stutter is observed.

### 4.7 Focus

`focus-visible` ring: `2px` `--primary-600` offset 2px, visible in both modes — keyboard operability is part of the demo story (core loop fully keyboard-driven).

### 4.8 Empty & loading states

Every empty state (per ux-writing.md): one-line orientation + muted explanation + one verb-first action. Table loading: 3 skeleton rows shimmering 1.2s max, replaced progressively.

---

## 5. The "Lunas moment" (signature state)

When reconciliation matches a payment: row background crossfades `danger-100 → primary-100`, status pill flips `Overdue → Paid`, check icon scales in, recovered counter ticks up with tabular digits, and a one-time toast reads from the ux-writing inventory. Total elapsed ≤ 1.2s. This is the emotional peak of the demo — restraint everywhere else makes it land.
