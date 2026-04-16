'use strict';

registerServiceWorker();

const projectsEl = document.getElementById('projects');
const lastUpdatedEl = document.getElementById('last-updated');

// relativeTime imported from utils.js

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
  header.innerHTML = `<span class="project-name" title="${esc(project.path)}">${esc(project.name)}</span>`;

  if (project.commits && project.commits.length > 0) {
    const badge = document.createElement('span');
    badge.className = 'commit-age-badge';
    badge.textContent = relativeTime(project.commits[0].timestamp);
    badge.title = absTime(project.commits[0].timestamp);
    header.appendChild(badge);
  } else {
    const badge = document.createElement('span');
    badge.className = 'commit-age-badge';
    badge.textContent = 'no commits yet';
    header.appendChild(badge);
  }

  card.appendChild(header);

  // Pull button
  const pullBtn = document.createElement('button');
  pullBtn.className = 'btn-sm pull-btn';
  pullBtn.textContent = 'Pull';
  pullBtn.addEventListener('click', async () => {
    pullBtn.disabled = true;
    pullBtn.textContent = 'Pulling…';
    try {
      const res = await fetch('/api/pull', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: project.path }),
      });
      const data = await res.json();
      if (!res.ok) {
        pullBtn.textContent = data.error ? data.error.split('\n')[0] : 'Failed';
        pullBtn.title = data.error || '';
        pullBtn.classList.add('pull-error');
      } else {
        pullBtn.textContent = 'Done';
        pullBtn.classList.add('pull-success');
        load();
      }
    } catch (err) {
      pullBtn.textContent = 'Failed';
      pullBtn.classList.add('pull-error');
    } finally {
      setTimeout(() => {
        pullBtn.disabled = false;
        pullBtn.textContent = 'Pull';
        pullBtn.title = '';
        pullBtn.classList.remove('pull-success', 'pull-error');
      }, 5000);
    }
  });
  card.appendChild(pullBtn);

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
