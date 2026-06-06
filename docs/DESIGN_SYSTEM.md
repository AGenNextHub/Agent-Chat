# AGenNext Chat — Design System (draft)

> Draft. The tokens are canonical; components are specified, not yet coded. We refine later.

The platform is headless, so this design system is the **canonical reference any channel or
client UI must follow** — the brand made buildable. It is token-first: `brand/tokens.json`
is the single source of truth, `brand/tokens.css` is its generated export.

## Principles (inherited from the brand)

1. **Inspectable by default** — the UI must make the agent's reasoning visible (the trace),
   never hide logic. Transparency is a UI requirement, not a nicety.
2. **Calm, not flashy** — trust through clarity. Restraint, whitespace, few accents.
3. **Accessible** — WCAG 2.2 AA minimum; keyboard-first; reduced-motion honored.
4. **Open** — open fonts (Inter, JetBrains Mono, both OFL); tokens are portable.
5. **Canonical & undiluted** — use tokens, never raw hex; off-token use is off-brand.

## Tokens

See [`brand/tokens.json`](brand/tokens.json) (canonical) / [`brand/tokens.css`](brand/tokens.css).

| Group | Highlights |
|---|---|
| Color (semantic) | `primary` indigo · `accent` cyan · `danger` rose · `text`/`muted` · `surface`/`bg` |
| Type | Inter (UI), JetBrains Mono (code/traces); scale xs→display (12→36px) |
| Space | 4px base scale (1–12) |
| Radius | sm 6 · md 10 · **lg 16 (the squircle/gate radius)** · pill |
| Elevation | sm/md/lg shadows |
| Motion | 120/200/320ms · easing `cubic-bezier(.2,0,0,1)` · reduced-motion → 0 |

**Contrast guidance (AA):** `primary` (#4338CA) on white passes AA for text. `accent` cyan
(#22D3EE) is **accent-only** — it fails AA as text on white; use it for the agent node,
focus rings, and large/icon elements, never body text.

## Components (spec)

Each component lists anatomy + the states it must implement
(`default · hover · active · focus-visible · disabled`).

- **Button** — `primary` (filled indigo), `secondary` (outline), `ghost`, `danger`. Radius
  `md`; focus-visible = 2px `focusRing`.
- **Composer / Input** — single + multiline; clear focus ring; error uses `danger`.
- **Message bubble** — `user` (surface-raised, right) · `agent` (surface, left, with the
  loop mark). System/escalation notices use a distinct banner (below).
- **Trace viewer** — the inspectability surface. Renders the turn trace as ordered step
  chips with semantic color:
  | Step | Color |
  |---|---|
  | reason | slate |
  | guard `allow` | success · guard `deny` | danger |
  | act | primary |
  | screen/shield `blocked` | danger |
  | observe / persist / exit | slate / accent |
- **Status pills** — admission/exit states: `admitted` (success), `denied`/`blocked`
  (danger), `escalated` (warning).
- **Escalation banner (human-in-the-loop)** — shown when `Result.Escalated`; warning-toned,
  names the stakeholder being engaged. Never blank; the agent always shows an exit.
- **Badge/Tag** — capability names (mono) and sandbox profile (`isolated`/`restricted`/
  `trusted`).

## Patterns

- **Chat surface** — composer pinned bottom; messages scroll; the trace viewer is a
  collapsible side/inline panel (inspectability is one tap away, never hidden).
- **Trust affordances** — every agent answer can reveal its trace; blocked/escalated turns
  state *why*, in plain language (no internal detail bleeds).

## Accessibility checklist

- Contrast AA for text; focus-visible on every interactive element; full keyboard path;
  `prefers-reduced-motion` respected; target size ≥ 24px.

## Open items (refine later)

Full component code, dark-theme contrast pass, motion specs per component, iconography set,
and a token build pipeline (`tokens.json` → CSS/TS/Swift).
