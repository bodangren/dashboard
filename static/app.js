'use strict';

registerServiceWorker();

(function() {
  const stored = localStorage.getItem('dashboard-theme');
  const theme = stored && THEMES[stored] ? stored : 'dark';
  applyTheme(theme);
  updateThemeIcon(theme);

  if (typeof NotificationService !== 'undefined') {
    const pref = NotificationService.getPreference();
    if (pref === 'denied' && localStorage.getItem('dashboard-notification-enabled') === null) {
      NotificationService.setEnabled(true);
    }
  }
})();

function updateThemeIcon(theme) {
  const iconDark = document.getElementById('theme-icon-dark');
  const iconLight = document.getElementById('theme-icon-light');
  const iconContrast = document.getElementById('theme-icon-contrast');
  if (!iconDark) return;
  iconDark.classList.add('hidden');
  iconLight.classList.add('hidden');
  iconContrast.classList.add('hidden');
  if (theme === 'light') {
    iconLight.classList.remove('hidden');
  } else if (theme === 'highContrast') {
    iconContrast.classList.remove('hidden');
  } else {
    iconDark.classList.remove('hidden');
  }
}

const themeToggle = document.getElementById('theme-toggle');
if (themeToggle) {
  themeToggle.addEventListener('click', function() {
    const newTheme = cycleTheme();
    updateThemeIcon(newTheme);
  });
}

document.addEventListener('keydown', function(e) {
  if (e.ctrlKey && e.key === 't') {
    e.preventDefault();
    const newTheme = cycleTheme();
    updateThemeIcon(newTheme);
  }
});

const settingsToggle = document.getElementById('settings-toggle');
const settingsPanel = document.getElementById('settings-panel');
const densitySelect = document.getElementById('density-select');
const sortSelect = document.getElementById('sort-select');

if (settingsToggle) {
  settingsToggle.addEventListener('click', function() {
    settingsPanel.classList.toggle('hidden');
  });
}

if (densitySelect) {
  densitySelect.value = getPreferences().density;
  densitySelect.addEventListener('change', function() {
    updatePreference('density', densitySelect.value);
  });
}

if (sortSelect) {
  sortSelect.value = getPreferences().sortOrder;
  sortSelect.addEventListener('change', function() {
    updatePreference('sortOrder', sortSelect.value);
    renderProjects();
  });
}

const groupDirSelect = document.getElementById('group-directory');
if (groupDirSelect) {
  groupDirSelect.checked = getPreferences().groupByDirectory;
  groupDirSelect.addEventListener('change', function() {
    updatePreference('groupByDirectory', groupDirSelect.checked);
    renderProjects();
  });
}

(function() {
  const notifEnabled = document.getElementById('notifications-enabled');
  const healthAlerts = document.getElementById('health-alerts-enabled');
  const agentAlerts = document.getElementById('agent-alerts-enabled');
  const aiAlerts = document.getElementById('ai-alerts-enabled');
  const quietStart = document.getElementById('quiet-start');
  const quietEnd = document.getElementById('quiet-end');
  const permBtn = document.getElementById('notification-permission-btn');

  if (notifEnabled && typeof NotificationService !== 'undefined') {
    notifEnabled.checked = NotificationService.isEnabled();
    notifEnabled.addEventListener('change', function() {
      NotificationService.setEnabled(notifEnabled.checked);
    });
  }

  if (healthAlerts) {
    healthAlerts.checked = NotificationService.isEnabled();
  }

  if (agentAlerts && typeof NotificationService !== 'undefined') {
    agentAlerts.checked = NotificationService.isAgentErrorsEnabled();
    agentAlerts.addEventListener('change', function() {
      NotificationService.setAgentErrorsEnabled(agentAlerts.checked);
    });
  }

  if (aiAlerts && typeof NotificationService !== 'undefined') {
    aiAlerts.checked = NotificationService.isAIInsightsEnabled();
    aiAlerts.addEventListener('change', function() {
      NotificationService.setAIInsightsEnabled(aiAlerts.checked);
    });
  }

  if (quietStart && quietEnd && typeof NotificationService !== 'undefined') {
    const qh = NotificationService.getQuietHours();
    if (qh.start) quietStart.value = qh.start;
    if (qh.end) quietEnd.value = qh.end;

    quietStart.addEventListener('change', function() {
      NotificationService.setQuietHours(quietStart.value, quietEnd.value);
    });
    quietEnd.addEventListener('change', function() {
      NotificationService.setQuietHours(quietStart.value, quietEnd.value);
    });
  }

  if (permBtn && typeof NotificationService !== 'undefined') {
    permBtn.addEventListener('click', function() {
      NotificationService.requestPermission();
    });
  }
})();

(function() {
  loadPreferences();
})();

const projectsEl = document.getElementById('projects');
const searchInput = document.getElementById('search-input');
const searchBtn = document.getElementById('search-btn');
const filterToggle = document.getElementById('filter-toggle');
const filterPanel = document.getElementById('filter-panel');
const filterRepo = document.getElementById('filter-repo');
const filterAuthor = document.getElementById('filter-author');
const filterDate = document.getElementById('filter-date');
const searchResults = document.getElementById('search-results');

let projects = [];
let searchTimeout = null;
let activeTagFilter = null;

function renderTagChips(card, repoPath) {
  const tags = TagManager.getTagsForRepo(repoPath);
  if (tags.length === 0) return;

  const tagRow = document.createElement('div');
  tagRow.className = 'tag-chips-row';

  for (const tag of tags) {
    const chip = document.createElement('span');
    chip.className = 'tag-chip';
    chip.textContent = tag;
    chip.dataset.tag = tag;
    chip.dataset.repo = repoPath;
    chip.addEventListener('click', function() {
      toggleTagFilter(tag);
    });
    tagRow.appendChild(chip);
  }

  card.appendChild(tagRow);
}

function toggleTagFilter(tag) {
  if (activeTagFilter === tag) {
    activeTagFilter = null;
  } else {
    activeTagFilter = tag;
  }
  syncTagFilterBar();
  renderProjects();
}

function syncTagFilterBar() {
  const bar = document.getElementById('tag-filter-bar');
  const chips = document.getElementById('tag-filter-chips');
  const clearBtn = document.getElementById('tag-filter-clear');

  if (!bar) return;

  if (activeTagFilter) {
    bar.classList.remove('hidden');
    chips.innerHTML = '<span class="tag-filter-active">' + esc(activeTagFilter) + '</span>';
    clearBtn.classList.remove('hidden');
  } else {
    bar.classList.add('hidden');
    chips.innerHTML = '';
    clearBtn.classList.add('hidden');
  }
}

function clearTagFilter() {
  activeTagFilter = null;
  syncTagFilterBar();
  renderProjects();
}

function initTagFilterBar() {
  const clearBtn = document.getElementById('tag-filter-clear');
  if (clearBtn) {
    clearBtn.addEventListener('click', clearTagFilter);
  }
}

function escapeRegex(str) {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

function highlightMatch(text, query) {
  if (!query) return esc(text);
  const regex = new RegExp('(' + escapeRegex(query) + ')', 'gi');
  return esc(text).replace(regex, '<span class="highlight">$1</span>');
}

function renderSearchResult(result, query) {
  const item = document.createElement('div');
  item.className = 'search-result-item';

  const repo = document.createElement('div');
  repo.className = 'search-result-repo';
  repo.textContent = result.repo_path || result.repoPath || '';
  item.appendChild(repo);

  const hash = document.createElement('span');
  hash.className = 'search-result-hash';
  hash.textContent = result.hash.substring(0, 7);
  item.appendChild(hash);

  const message = document.createElement('div');
  message.className = 'search-result-message';
  message.innerHTML = highlightMatch(result.message || '', query);
  item.appendChild(message);

  const meta = document.createElement('div');
  meta.className = 'search-result-meta';
  meta.textContent = `${result.author} · ${relativeTime(result.timestamp)}`;
  item.appendChild(meta);

  item.addEventListener('click', () => {
    const repoPath = encodeURIComponent(result.repo_path || result.repoPath);
    const hash = encodeURIComponent(result.hash);
    window.location.href = `diff.html?repo=${repoPath}&hash=${hash}`;
  });

  return item;
}

async function performSearch(query) {
  if (!query.trim()) {
    searchResults.classList.add('hidden');
    searchResults.innerHTML = '';
    return;
  }

  const params = new URLSearchParams({ q: query });
  const repoPath = filterRepo.value;
  const author = filterAuthor.value.trim();
  const dateFrom = filterDate.value;

  if (repoPath) params.append('repo', repoPath);
  if (author) params.append('author', author);
  if (dateFrom) params.append('date', dateFrom);

  try {
    const res = await fetch(`/api/search?${params}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const data = await res.json();
    const results = data.results || [];

    searchResults.innerHTML = '';
    if (results.length === 0) {
      searchResults.innerHTML = '<p class="loading">no results found</p>';
    } else {
      for (const r of results) {
        searchResults.appendChild(renderSearchResult(r, query));
      }
    }
    searchResults.classList.remove('hidden');
  } catch (err) {
    searchResults.innerHTML = `<p class="error">search error: ${esc(err.message)}</p>`;
    searchResults.classList.remove('hidden');
  }
}

function scheduleSearch() {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    performSearch(searchInput.value);
  }, 300);
}

function toggleFilterPanel() {
  filterPanel.classList.toggle('hidden');
  filterToggle.classList.toggle('active');
}

if (searchInput) {
  searchInput.addEventListener('input', scheduleSearch);
  searchInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      if (searchTimeout) clearTimeout(searchTimeout);
      performSearch(searchInput.value);
    }
  });
}

if (searchBtn) {
  searchBtn.addEventListener('click', () => {
    if (searchTimeout) clearTimeout(searchTimeout);
    performSearch(searchInput.value);
  });
}

if (filterToggle) {
  filterToggle.addEventListener('click', toggleFilterPanel);
}

// relativeTime imported from utils.js

/** Absolute timestamp for tooltip */
function absTime(isoStr) {
  return new Date(isoStr).toLocaleString();
}

/** Render a single project card */
/** Sort projects by current sort order preference */
function sortProjects(projects) {
  const prefs = getPreferences();
  const sorted = [...projects];
  if (prefs.sortOrder === 'alphabetical') {
    sorted.sort((a, b) => a.name.localeCompare(b.name));
  } else {
    sorted.sort((a, b) => {
      const aTime = a.commits && a.commits[0] ? new Date(a.commits[0].timestamp) : new Date(0);
      const bTime = b.commits && b.commits[0] ? new Date(b.commits[0].timestamp) : new Date(0);
      return bTime - aTime;
    });
  }
  return sorted;
}

/** Render all projects (call after preference changes) */
function renderProjects() {
  const prefs = getPreferences();
  const sorted = sortProjects(projects);
  projectsEl.innerHTML = '';

  if (prefs.groupByDirectory) {
    renderProjectsGrouped(sorted);
  } else {
    for (const p of sorted) {
      if (activeTagFilter && !TagManager.getTagsForRepo(p.path).includes(activeTagFilter)) {
        continue;
      }
      projectsEl.appendChild(renderProject(p));
    }
  }
}

function getParentDirectory(repoPath) {
  const parts = repoPath.split('/');
  if (parts.length > 1) {
    return parts.slice(0, -1).join('/');
  }
  return '';
}

function renderProjectsGrouped(sorted) {
  const groups = {};
  for (const p of sorted) {
    if (activeTagFilter && !TagManager.getTagsForRepo(p.path).includes(activeTagFilter)) {
      continue;
    }
    const dir = getParentDirectory(p.path);
    if (!groups[dir]) groups[dir] = [];
    groups[dir].push(p);
  }

  const dirs = Object.keys(groups).sort();
  for (const dir of dirs) {
    const groupHeader = document.createElement('div');
    groupHeader.className = 'directory-group-header';

    const toggle = document.createElement('button');
    toggle.className = 'directory-group-toggle';
    toggle.textContent = dir || 'Root';
    toggle.addEventListener('click', function() {
      toggle.classList.toggle('collapsed');
      contentEl.classList.toggle('hidden');
    });

    const contentEl = document.createElement('div');
    contentEl.className = 'directory-group-content';

    for (const p of groups[dir]) {
      contentEl.appendChild(renderProject(p));
    }

    groupHeader.appendChild(toggle);
    groupHeader.appendChild(contentEl);
    projectsEl.appendChild(groupHeader);
  }
}

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

  renderTagChips(card, project.path);

  if (project.health) {
    attachHealthBadge(card, project.health);
  }

  attachBranchSelector(card, project);
  attachStashToggle(card, project);

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
    const loadedProjects = await res.json();
    projects = loadedProjects;

    if (typeof NotificationService !== 'undefined') {
      NotificationService.checkHealthAndNotify(projects);
    }

    filterRepo.innerHTML = '<option value="">All repos</option>';
    for (const p of projects) {
      const opt = document.createElement('option');
      opt.value = p.path;
      opt.textContent = p.name;
      filterRepo.appendChild(opt);
    }

    projectsEl.innerHTML = '';
    if (projects.length === 0) {
      projectsEl.innerHTML = '<p class="loading">no repos found</p>';
      return;
    }
    renderProjects();
  } catch (err) {
    projectsEl.innerHTML = `<p class="error">error: ${esc(err.message)}</p>`;
  }
}

initTagFilterBar();
load();
// Auto-refresh every 15 minutes
setInterval(load, 15 * 60 * 1000);
