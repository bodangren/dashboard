# Spec: Deduplicate Service Worker Registration (JS-02)

## Problem

Both `static/app.js` and `static/diff.js` contain identical service worker registration code (lines 3-14):
```javascript
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
      .then(registration => {
        console.log('SW registered:', registration.scope);
      })
      .catch(error => {
        console.log('SW registration failed:', error);
      });
  });
}
```

## Solution

1. Add `registerServiceWorker()` function to `utils.js` that encapsulates this logic
2. Replace inline code in `app.js` and `diff.js` with a single call to `registerServiceWorker()`
3. All three pages (index.html, agents.html, diff.html) already load utils.js before page-specific scripts, so no HTML changes needed

## Files to Modify

- `static/utils.js` — add `registerServiceWorker()` function
- `static/app.js` — remove lines 3-14, add `registerServiceWorker()` call
- `static/diff.js` — remove lines 3-14, add `registerServiceWorker()` call
- `conductor/tech-debt.md` — mark JS-02 as Resolved