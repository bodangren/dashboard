# Dashboard Customization & Themes

## Overview

Add personalization options to the dashboard including theme selection and layout preferences. This addresses a product vision gap where the current fixed UI doesn't accommodate different viewing preferences or accessibility needs.

## Problem Statement

The dashboard currently has a single dark theme with a fixed card layout. Users with different monitor sizes, accessibility needs, or visual preferences have no way to customize their experience. The mobile-first CSS refactor (UX-07) provides a foundation but lacks user-facing controls.

## Solution

Implement a lightweight customization system that persists preferences locally:

### Theme Options
- **Dark Theme**: Current default (dark background, light text)
- **Light Theme**: Inverted colors for bright environments
- **High Contrast**: Enhanced contrast for accessibility

### Layout Options
- **Card Density**: Compact (smaller cards, more visible) vs Comfortable (larger cards, more details)
- **Sort Order**: Recent commits first (default) vs Alphabetical vs Last pull time
- **Show/Hide**: Toggle visibility of sections (activity feed, agent monitoring, pull status)

### Persistence
- Store preferences in localStorage
- Apply theme on page load before render to prevent flash
- No backend changes needed—pure frontend feature

## Acceptance Criteria

- [ ] Theme switcher UI in header/settings area
- [ ] Three themes render correctly without flash of unstyled content
- [ ] Layout density options change card sizing
- [ ] Sort order preference persists across sessions
- [ ] Section visibility toggles work correctly
- [ ] Preferences stored in localStorage and applied on load

## Out of Scope

- Custom color pickers
- Multiple font size options (use browser zoom instead)
- Per-repo customization
- Cloud sync of preferences
