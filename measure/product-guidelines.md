# Product Guidelines

## Visual Identity

- **Theme**: Terminal-inspired dark UI
- **Background**: Near-black (`#0d0d0d` or `#111111`)
- **Primary accent**: Green (`#00ff41` or similar terminal green) or amber (`#ffb000`)
- **Text**: Monospace font (e.g., `JetBrains Mono`, `Fira Code`, `monospace` fallback)
- **Borders/separators**: Subtle, low-contrast (e.g., `#333`)
- **Hover states**: Slight glow or brightness increase — no heavy animations

## Typography

- Use monospace throughout for all data (commit hashes, messages, timestamps)
- Keep font sizes compact — this is a dense information display, not a landing page
- Timestamps: relative format preferred (e.g., `2h ago`) with absolute on hover

## Layout

- Card-per-project layout, stacked vertically
- Most recently committed project at the top (dynamic sort)
- Each card shows: project name, latest commits (message, author, timestamp)
- Compact — maximize information density

## Voice & Tone

- Terse and technical
- Labels use short, familiar git vocabulary: `HEAD`, `author`, `hash`, `diff`, `+N -N`
- No marketing language, no tooltips with long explanations
- Error states are direct: `pull failed`, `not a git repo`

## Interactions

- Hover a commit row → show commit note/body in a tooltip or inline expand
- Click a commit row → open diff report view (full-page or modal)
- Diff report uses standard `+` (green) / `-` (red) diff coloring
