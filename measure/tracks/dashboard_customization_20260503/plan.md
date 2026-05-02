# Dashboard Customization & Themes Plan

## Phase 1: Theme System

- [ ] Create `themes.js` with theme definitions (dark, light, high-contrast)
- [ ] Add CSS custom properties for each theme in style.css
- [ ] Implement `applyTheme(themeName)` function that updates :root variables
- [ ] Create `loadTheme()` to read from localStorage on page load
- [ ] Add theme toggle button to header HTML
- [ ] Wire theme toggle click handler in app.js
- [ ] Add tests for theme application and localStorage persistence

## Phase 2: Layout Preferences

- [ ] Create `preferences.js` for managing all user preferences
- [ ] Add card density CSS classes (compact/comfortable)
- [ ] Implement density toggle in settings UI
- [ ] Add sort order selector (recent/alphabetical/last-pull)
- [ ] Update project card sorting logic in app.js
- [ ] Add section visibility toggles (activity, agents, pull-status)
- [ ] Wire preference changes to localStorage persistence
- [ ] Add tests for preference storage and retrieval

## Phase 3: Integration & Polish

- [ ] Add CSS transitions for theme switching (0.2s ease)
- [ ] Handle flash-of-unstyled-content: apply theme in <head> via inline script
- [ ] Add keyboard shortcuts for quick theme toggle (Ctrl+T)
- [ ] Update service worker to cache theme preference
- [ ] Add integration test for full preference flow
- [ ] Manual smoke test: visual verification across all themes
