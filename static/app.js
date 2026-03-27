'use strict';

const projectsEl = document.getElementById('projects');
const lastUpdatedEl = document.getElementById('last-updated');

/** Format a UTC ISO timestamp as a relative string, e.g. "2h ago" */
function relativeTime(isoStr) {
  const diffMs = Date.now() - new Date(isoStr).getTime();
  const s = Math.floor(diffMs / 1000);
  if (s < 60)  return `${s}s ago`;
  const m = Math.floor(s / 60);
  if (m < 60)  return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24)  return `${h}h ago`;
  const d = Math.floor(h / 24);
  return `${d}d ago`;
}

/** Absolute timestamp for tooltip */
function absTime(isoStr) {
  return new Date(isoStr).toLocaleString();
}

/** Render a single project card */
function renderProject(project) {
  const card = document.createElement('div');
  card.className = 'project-card';

  const header = document.createElement('div');
  header.className = 'project-header';
  header.innerHTML = `<span class="project-name">${esc(project.name)}</span>`
    + `<span class="project-path">${esc(project.path)}</span>`;
  card.appendChild(header);

  for (const commit of project.commits) {
    const row = document.createElement('a');
    row.className = 'commit-row';
    row.href = `diff.html?repo=${encodeURIComponent(project.path)}&hash=${encodeURIComponent(commit.hash)}`;

    // Hover tooltip: show git notes if present, fall back to body
    const tooltip = commit.notes || commit.body || '';
    if (tooltip) row.title = tooltip;

    row.innerHTML =
      `<span class="commit-hash">${esc(commit.hash)}</span>`
      + `<span class="commit-message">${esc(commit.message)}</span>`
      + `<span class="commit-meta" title="${esc(absTime(commit.timestamp))}">`
      +   `${esc(commit.author)} · ${relativeTime(commit.timestamp)}`
      + `</span>`;

    card.appendChild(row);
  }

  return card;
}

/** Escape HTML special characters */
function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

async function load() {
  projectsEl.innerHTML = '<p class="loading">loading…</p>';
  try {
    const res = await fetch('/api/projects');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const projects = await res.json();

    projectsEl.innerHTML = '';
    if (projects.length === 0) {
      projectsEl.innerHTML = '<p class="loading">no repos found</p>';
      return;
    }
    for (const p of projects) {
      projectsEl.appendChild(renderProject(p));
    }
    lastUpdatedEl.textContent = `updated ${new Date().toLocaleTimeString()}`;
  } catch (err) {
    projectsEl.innerHTML = `<p class="error">error: ${esc(err.message)}</p>`;
  }
}

load();
// Auto-refresh every 15 minutes
setInterval(load, 15 * 60 * 1000);
