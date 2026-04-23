# Track: Fix JS-01 — Replace `var` with `const`/`let` in agents.js

## Context

The JavaScript styleguide forbids `var` (rule: "Use `const` by default, `let` if reassignment is needed"). Currently agents.js has ~30 instances of `var` that must be converted.

## Scope

- `static/agents.js` only
- No other files in this track

## Approach

1. Read agents.js line by line
2. For each `var` declaration, determine if reassigned later:
   - If never reassigned → `const`
   - If reassigned → `let`
3. Update only the declaration, not the usage

## Verification

- Manual: load agents page in browser, verify no console errors
- The JavaScript is loaded via `<script src="agents.js">` in a page that already works

## Track Completion

Mark JS-01 as Resolved in tech-debt.md