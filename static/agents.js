'use strict';

// Register Service Worker for PWA
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

const agentsListEl = document.getElementById('agents-list');
const formContainerEl = document.getElementById('agent-form-container');
const addBtn = document.getElementById('add-agent-btn');

function esc(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function harnessLabel(h) {
  const colors = { opencode: '#39ff14', gemini: '#ff6600', codex: '#00aaff' };
  return `<span class="harness-badge" style="color:${colors[h] || '#ccc'}">${esc(h)}</span>`;
}

function parseHours(hourStr) {
  if (hourStr === '*') return Array.from({ length: 24 }, (_, i) => i);
  const hours = new Set();
  hourStr.split(',').forEach(part => {
    if (part.includes('/')) {
      const [, step] = part.split('/');
      for (let i = 0; i < 24; i += parseInt(step)) hours.add(i);
    } else if (part.includes('-')) {
      const [start, end] = part.split('-').map(Number);
      for (let i = start; i <= end; i++) hours.add(i);
    } else {
      hours.add(parseInt(part));
    }
  });
  return hours;
}

function parseDays(dow) {
  if (dow === '*') return new Set([0, 1, 2, 3, 4, 5, 6]);
  const days = new Set();
  dow.split(',').forEach(part => {
    if (part.includes('-')) {
      const [start, end] = part.split('-').map(Number);
      for (let i = start; i <= end; i++) days.add(i);
    } else {
      days.add(parseInt(part));
    }
  });
  return days;
}

function renderTimingVisualization(cron) {
  const parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return `<span class="agent-schedule">${esc(cron)}</span>`;
  const [min, hour, dom, mon, dow] = parts;

  const activeHours = parseHours(hour);
  const activeDays = parseDays(dow);
  const minute = min === '*' ? '??' : min.padStart(2, '0');

  const daysHtml = Array.from({ length: 7 }, (_, i) =>
    `<span class="day-block ${activeDays.has(i) ? 'active' : 'inactive'}"></span>`
  ).join('');

  const hoursHtml = Array.from({ length: 24 }, (_, i) =>
    `<span class="hour-block ${activeHours.has(i) ? 'active' : 'inactive'}"></span>`
  ).join('');

  return `<span class="agent-days">${daysHtml}</span>`
    + `<span class="agent-hours">${hoursHtml}</span>`
    + `<span class="agent-minute">:${minute}</span>`;
}

function scheduleHuman(cron) {
  const parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return cron;
  const [min, hour, dom, mon, dow] = parts;
  if (dom === '*' && mon === '*' && dow === '*') {
    if (min.startsWith('*/')) return `Every ${min.slice(2)} minutes`;
    if (hour.startsWith('*/')) {
      const h = hour.slice(2);
      if (min === '0') return `Every ${h} hours`;
      return `Every ${h}h at :${min}`;
    }
    if (hour !== '*' && min !== '*') return `Daily at ${hour}:${min.padStart(2, '0')}`;
  }
  return cron;
}

function renderAgentCard(agent, index) {
  const card = document.createElement('div');
  card.className = 'agent-card' + (agent.enabled ? '' : ' agent-disabled');

  const dir = agent.directory.split('/').pop() || agent.directory;

  card.innerHTML =
    `<div class="agent-header">`
    +   `<span class="agent-project">${esc(dir)}</span>`
    +   `${harnessLabel(agent.harness)}`
    +   `<span class="agent-model">${esc(agent.model)}</span>`
    +   `<span class="agent-schedule" title="${esc(scheduleHuman(agent.schedule))}">${renderTimingVisualization(agent.schedule)}</span>`
    +   `<span class="agent-status ${agent.enabled ? 'status-on' : 'status-off'}">${agent.enabled ? 'ON' : 'OFF'}</span>`
    + `</div>`
    + `<div class="agent-actions">`
    +   `<button class="btn-sm btn-toggle" data-idx="${index}">${agent.enabled ? '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter"><path d="M10.586 14.957 12 13.543 13.414 14.957a2 2 0 1 0 2.828-2.828L14.828 10.829l1.414-1.415a2 2 0 1 0-2.828-2.828L12 7.586 10.586 6.171a2 2 0 1 0-2.828 2.829L9.171 10.414l-1.414 1.415a2 2 0 1 0 2.828 2.828L10.586 14.957Z"/></svg> Disable' : '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter"><path d="M5 12h14"/><path d="m5 12 7-7 7 7"/><path d="M5 12v7a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-7"/></svg> Enable'}</button>`
    +   `<button class="btn-sm btn-edit" data-idx="${index}"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"/></svg> Edit</button>`
    +   `<button class="btn-sm btn-delete" data-idx="${index}"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg> Delete</button>`
    +   `<button class="btn-sm btn-log" data-idx="${index}"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg> Logs</button>`
    + `</div>`
    + `<div class="agent-log-container hidden" id="agent-log-${index}"></div>`;

  return card;
}

async function loadAgents() {
  agentsListEl.innerHTML = '<p class="loading">loading agents…</p>';
  try {
    const res = await fetch('/api/agents');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();

    agentsListEl.innerHTML = '';
    if (!data.agents || data.agents.length === 0) {
      agentsListEl.innerHTML = '<p class="loading">no agents configured</p>';
      return;
    }

    const grouped = {};
    data.agents.forEach((a, i) => {
      const dir = a.directory.split('/').pop() || a.directory;
      if (!grouped[dir]) grouped[dir] = [];
      grouped[dir].push({ ...a, _idx: i });
    });

    for (const [project, agents] of Object.entries(grouped)) {
      const group = document.createElement('div');
      group.className = 'agent-group';
      group.innerHTML = `<h2 class="agent-group-title">${esc(project)}</h2>`;
      for (const a of agents) {
        group.appendChild(renderAgentCard(a, a._idx));
      }
      agentsListEl.appendChild(group);
    }
  } catch (err) {
    agentsListEl.innerHTML = `<p class="error">error: ${esc(err.message)}</p>`;
  }
}

agentsListEl.addEventListener('click', async (e) => {
  const btn = e.target.closest('button');
  if (!btn) return;
  const idx = btn.dataset.idx;

  if (btn.classList.contains('btn-toggle')) {
    btn.disabled = true;
    await fetch(`/api/agents/${idx}/toggle`, { method: 'PATCH' });
    loadAgents();
  } else if (btn.classList.contains('btn-delete')) {
    if (!confirm('Delete this agent from crontab?')) return;
    btn.disabled = true;
    await fetch(`/api/agents/${idx}`, { method: 'DELETE' });
    loadAgents();
  } else if (btn.classList.contains('btn-log')) {
    const logEl = document.getElementById(`agent-log-${idx}`);
    logEl.classList.toggle('hidden');
    if (!logEl.classList.contains('hidden') && !logEl.dataset.loaded) {
      logEl.innerHTML = '<p class="loading">loading log…</p>';
      const res = await fetch(`/api/agents/${idx}/log`);
      const logInfo = await res.json();
      logEl.dataset.loaded = '1';
      if (!logInfo.exists) {
        logEl.innerHTML = '<p class="loading">No log file found</p>';
      } else {
        logEl.innerHTML = `<pre class="log-content">${esc(logInfo.lines.join('\n'))}</pre>`;
        if (logInfo.last_run && logInfo.last_run !== '0001-01-01T00:00:00Z') {
          logEl.innerHTML += `<span class="log-last-run">Last run: ${new Date(logInfo.last_run).toLocaleString()}</span>`;
        }
      }
    }
  } else if (btn.classList.contains('btn-edit')) {
    showEditForm(idx);
  }
});

addBtn.addEventListener('click', () => showAddForm());

async function showAddForm() {
  const repos = await fetch('/api/projects').then(r => r.json()).catch(() => []);
  formContainerEl.classList.remove('hidden');
  formContainerEl.innerHTML = buildForm('Add Agent', null, repos);
}

async function showEditForm(idx) {
  const [agentsRes, reposRes] = await Promise.all([
    fetch('/api/agents').then(r => r.json()),
    fetch('/api/projects').then(r => r.json()).catch(() => [])
  ]);
  const agent = agentsRes.agents[idx];
  if (!agent) return;
  formContainerEl.classList.remove('hidden');
  formContainerEl.innerHTML = buildForm('Edit Agent', agent, repos, idx);
}

function buildForm(title, agent, repos, editIdx) {
  const dirOptions = repos.map(p =>
    `<option value="${esc(p.path)}" ${agent && agent.directory === p.path ? 'selected' : ''}>${esc(p.name)}</option>`
  ).join('');

  const harnessOptions = ['opencode', 'gemini', 'codex'].map(h =>
    `<option value="${h}" ${agent && agent.harness === h ? 'selected' : ''}>${h}</option>`
  ).join('');

  return `<div class="agent-form">`
    + `<h3>${title}</h3>`
    + `<form id="agent-crud-form" data-edit-idx="${editIdx !== undefined ? editIdx : ''}">`
    +   `<label>Schedule <input name="schedule" value="${agent ? esc(agent.schedule) : '0 */4 * * *'}" placeholder="0 */4 * * *"></label>`
    +   `<label>Directory <select name="directory"><option value="">Select project…</option>${dirOptions}</select></label>`
    +   `<label>Harness <select name="harness">${harnessOptions}</select></label>`
    +   `<label>Model <input name="model" value="${agent ? esc(agent.model) : ''}" placeholder="gpt-4o"></label>`
    +   `<label>Prompt <input name="prompt" value="${agent ? esc(agent.prompt) : ''}" placeholder="tasks.md"></label>`
    +   `<label>Log Path <input name="logPath" value="${agent ? esc(agent.log_path) : ''}" placeholder="/var/log/agent.log"></label>`
    +   `<div class="form-actions">`
    +     `<button type="submit" class="btn">Save</button>`
    +     `<button type="button" class="btn btn-cancel">Cancel</button>`
    +   `</div>`
    + `</form>`
    + `</div>`;
}

formContainerEl.addEventListener('click', (e) => {
  if (e.target.classList.contains('btn-cancel')) {
    formContainerEl.classList.add('hidden');
    formContainerEl.innerHTML = '';
  }
});

formContainerEl.addEventListener('submit', async (e) => {
  e.preventDefault();
  const form = e.target;
  const fd = new FormData(form);
  const body = {
    schedule: fd.get('schedule'),
    directory: fd.get('directory'),
    harness: fd.get('harness'),
    model: fd.get('model'),
    prompt: fd.get('prompt'),
    logPath: fd.get('logPath'),
  };

  const editIdx = form.dataset.editIdx;
  let url = '/api/agents';
  let method = 'POST';
  if (editIdx !== undefined && editIdx !== '') {
    url = `/api/agents/${editIdx}`;
    method = 'PUT';
  }

  const res = await fetch(url, {
    method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    alert('Error: ' + await res.text());
    return;
  }

  formContainerEl.classList.add('hidden');
  formContainerEl.innerHTML = '';
  loadAgents();
});

loadAgents();
