'use strict';

const feedEl = document.getElementById('activity-feed');
const loadingEl = document.getElementById('activity-loading');
const errorEl = document.getElementById('activity-error');
const emptyEl = document.getElementById('activity-empty');

const FILTER_KEY = 'activity_filters';
const CURSOR_KEY = 'activity_last_seen';

const typeFilters = {
  commit: true,
  agent: true,
  pull: true,
};

let ws = null;
let wsConnected = false;
let reconnectTimeout = null;
let lastSeenEventID = null;

function loadFilterState() {
  try {
    const saved = localStorage.getItem(FILTER_KEY);
    if (saved) {
      const parsed = JSON.parse(saved);
      Object.assign(typeFilters, parsed);
    }
  } catch (_) {}
}

function saveFilterState() {
  try {
    localStorage.setItem(FILTER_KEY, JSON.stringify(typeFilters));
  } catch (_) {}
}

function syncFilterButtons() {
  document.querySelectorAll('.filter-btn').forEach(btn => {
    const type = btn.dataset.type;
    btn.classList.toggle('active', typeFilters[type]);
  });
}

function eventIcon(type) {
  if (type === 'commit') {
    return '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square"><line x1="6" y1="3" x2="6" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/></svg>';
  }
  if (type === 'agent') {
    return '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square"><rect width="8" height="8" x="2" y="2" rx="1"/><path d="M14 2c1.1 0 2 .9 2 2v4c0 1.1-.9 2-2 2"/><path d="M20 2c1.1 0 2 .9 2 2v4c0 1.1-.9 2-2 2"/><path d="M10 22H5a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h5"/><path d="M16 22h5a2 2 0 0 0 2-2v-7a2 2 0 0 0-2-2h-5"/></svg>';
  }
  if (type === 'pull') {
    return '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square"><path d="M18 18H6"/><path d="M18 9a9 9 0 0 1-9 9"/><path d="M6 15a9 9 0 0 0 9-9"/></svg>';
  }
  return '';
}

function typeColor(type) {
  if (type === 'commit') return 'var(--green)';
  if (type === 'agent') return 'var(--amber)';
  if (type === 'pull') return 'var(--text-dim)';
  return 'var(--text)';
}

function repoName(repoPath) {
  if (!repoPath) return '';
  const parts = repoPath.split('/');
  return parts[parts.length - 1] || repoPath;
}

function renderEvent(event) {
  const card = document.createElement('div');
  card.className = 'event-card';
  card.dataset.id = event.id;

  const meta = event.metadata ? JSON.parse(event.metadata) : {};

  const header = document.createElement('div');
  header.className = 'event-header';

  const icon = document.createElement('span');
  icon.className = 'event-icon';
  icon.style.color = typeColor(event.type);
  icon.innerHTML = eventIcon(event.type);

  const repo = document.createElement('span');
  repo.className = 'event-repo';
  repo.textContent = repoName(event.repo);

  const time = document.createElement('span');
  time.className = 'event-time';
  time.textContent = relativeTime(event.timestamp);
  time.title = new Date(event.timestamp).toLocaleString();

  header.appendChild(icon);
  header.appendChild(repo);
  header.appendChild(time);

  const body = document.createElement('div');
  body.className = 'event-body';

  const msg = document.createElement('div');
  msg.className = 'event-message';
  msg.textContent = event.message || '(no message)';

  let subtext = '';
  if (event.type === 'commit' && meta.author) {
    subtext = meta.author;
  } else if (event.type === 'agent') {
    subtext = meta.agent_name || meta.agent_id || '';
    if (meta.status) {
      subtext += ' · ' + meta.status;
      if (meta.exit_code !== undefined && meta.exit_code !== 0) {
        subtext += ' (exit ' + meta.exit_code + ')';
      }
    }
    if (meta.error) {
      subtext += ' · ' + meta.error.split('\n')[0];
    }
  } else if (event.type === 'pull') {
    subtext = meta.success ? 'success' : ('error: ' + (meta.error || 'unknown'));
  }

  const sub = document.createElement('div');
  sub.className = 'event-sub';
  sub.textContent = subtext;

  body.appendChild(msg);
  if (subtext) body.appendChild(sub);

  card.appendChild(header);
  card.appendChild(body);

  return card;
}

let events = [];
let loadingMore = false;
let seenIds = new Set();

function applyFilters() {
  feedEl.innerHTML = '';
  let visible = 0;
  for (const ev of events) {
    if (typeFilters[ev.type]) {
      feedEl.appendChild(renderEvent(ev));
      visible++;
    }
  }
  emptyEl.classList.toggle('hidden', visible > 0);
}

async function load(since) {
  loadingEl.classList.remove('hidden');
  errorEl.classList.add('hidden');

  const params = new URLSearchParams();
  params.set('limit', '50');
  const types = Object.keys(typeFilters).filter(k => typeFilters[k]).join(',');
  if (types) params.set('types', types);
  if (since) params.set('since', since);

  try {
    const res = await fetch('/api/activity?' + params.toString());
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const data = await res.json();
    const newEvents = data.events || [];

    if (!since) {
      events = newEvents;
    } else {
      for (const ev of newEvents) {
        if (!seenIds.has(ev.id)) {
          events.unshift(ev);
          seenIds.add(ev.id);
        }
      }
    }

    for (const ev of events) {
      seenIds.add(ev.id);
      if (!lastSeenEventID || ev.id > lastSeenEventID) {
        lastSeenEventID = ev.id;
      }
    }
    saveLastSeenCursor();

    applyFilters();
  } catch (err) {
    errorEl.textContent = 'error loading activity: ' + err.message;
    errorEl.classList.remove('hidden');
  } finally {
    loadingEl.classList.add('hidden');
  }
}

function saveLastSeenCursor() {
  try {
    localStorage.setItem(CURSOR_KEY, lastSeenEventID || '');
  } catch (_) {}
}

function loadLastSeenCursor() {
  try {
    const saved = localStorage.getItem(CURSOR_KEY);
    if (saved) lastSeenEventID = saved;
  } catch (_) {}
}

function prependEvent(ev) {
  if (seenIds.has(ev.id)) return;
  seenIds.add(ev.id);
  if (!lastSeenEventID || ev.id > lastSeenEventID) {
    lastSeenEventID = ev.id;
  }
  saveLastSeenCursor();

  events.unshift(ev);

  if (typeFilters[ev.type]) {
    const card = renderEvent(ev);
    if (feedEl.firstChild) {
      feedEl.insertBefore(card, feedEl.firstChild);
    } else {
      feedEl.appendChild(card);
    }
    emptyEl.classList.add('hidden');
  }
}

function connectWS() {
  if (ws) {
    ws.close();
  }

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(protocol + '//' + location.host + '/ws/activity');

  ws.onopen = () => {
    wsConnected = true;
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    }
  };

  ws.onmessage = (msg) => {
    try {
      const ev = JSON.parse(msg.data);
      prependEvent(ev);
    } catch (_) {}
  };

  ws.onclose = () => {
    wsConnected = false;
    ws = null;
    reconnectTimeout = setTimeout(connectWS, 2000);
  };

  ws.onerror = () => {
    ws.close();
  };
}

function init() {
  loadFilterState();
  syncFilterButtons();
  loadLastSeenCursor();

  document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const type = btn.dataset.type;
      typeFilters[type] = !typeFilters[type];
      btn.classList.toggle('active', typeFilters[type]);
      saveFilterState();
      applyFilters();
    });
  });

  load();
  connectWS();
}

init();