'use strict';

registerServiceWorker();

const titleEl    = document.getElementById('commit-title');
const metaEl     = document.getElementById('commit-meta');
const diffEl     = document.getElementById('diff-container');

/** Classify a diff line into a CSS class */
function lineClass(line) {
  if (line.startsWith('+') && !line.startsWith('+++')) return 'diff-add';
  if (line.startsWith('-') && !line.startsWith('---')) return 'diff-remove';
  if (line.startsWith('@@'))                            return 'diff-hunk';
  if (line.startsWith('diff ') || line.startsWith('index ') ||
      line.startsWith('---')   || line.startsWith('+++'))   return 'diff-meta';
  return '';
}

/** Render raw unified diff string into colorized HTML */
function renderDiff(raw) {
  if (!raw) return '<span class="diff-hunk">(empty diff)</span>';
  return raw
    .split('\n')
    .map(line => {
      const cls = lineClass(line);
      const escaped = esc(line);
      return cls ? `<span class="${cls}">${escaped}</span>` : escaped;
    })
    .join('\n');
}

async function load() {
  const params = new URLSearchParams(window.location.search);
  const repo = params.get('repo');
  const hash = params.get('hash');

  if (!repo || !hash) {
    diffEl.textContent = 'missing repo or hash parameter';
    return;
  }

  titleEl.textContent = `${hash}`;
  diffEl.innerHTML = '<span class="loading">loading…</span>';

  try {
    const url = `/api/diff?repo=${encodeURIComponent(repo)}&hash=${encodeURIComponent(hash)}`;
    const res = await fetch(url);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();

    if (data.author || data.message) {
      const metaParts = [];
      if (data.author) metaParts.push(`<span>Author: ${esc(data.author)}</span>`);
      if (data.timestamp) metaParts.push(`<span>Date: ${new Date(data.timestamp).toLocaleString()}</span>`);
      if (data.message) metaParts.push(`<span>Message: ${esc(data.message)}</span>`);
      metaEl.innerHTML = metaParts.join('');
    } else {
      metaEl.innerHTML = '';
    }

    titleEl.textContent = data.message || data.hash || hash;
    diffEl.innerHTML = renderDiff(data.diff);
  } catch (err) {
    diffEl.innerHTML = `<span class="error">error: ${esc(err.message)}</span>`;
  }
}

load();
