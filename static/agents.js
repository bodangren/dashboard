'use strict';

const agentsListEl = document.getElementById('agents-list');
const addBtn = document.getElementById('add-agent-btn');

function harnessLabel(h) {
  const colors = { opencode: '#39ff14', gemini: '#ff6600', codex: '#00aaff' };
  return '<span class="harness-badge" style="color:' + (colors[h] || '#ccc') + '">' + esc(h) + '</span>';
}

function parseHours(hourStr) {
  if (hourStr === '*') return new Set(Array.from({ length: 24 }, function(_, i) { return i; }));
  const hours = new Set();
  hourStr.split(',').forEach(function(part) {
    if (part.includes('/')) {
      const step = part.split('/')[1];
      for (let i = 0; i < 24; i += parseInt(step)) hours.add(i);
    } else if (part.includes('-')) {
      const bounds = part.split('-').map(Number);
      for (let i = bounds[0]; i <= bounds[1]; i++) hours.add(i);
    } else {
      hours.add(parseInt(part));
    }
  });
  return hours;
}

function parseDays(dow) {
  if (dow === '*') return new Set([0, 1, 2, 3, 4, 5, 6]);
  const days = new Set();
  dow.split(',').forEach(function(part) {
    if (part.includes('-')) {
      const bounds = part.split('-').map(Number);
      for (let i = bounds[0]; i <= bounds[1]; i++) days.add(i);
    } else {
      days.add(parseInt(part));
    }
  });
  return days;
}

function renderTimingVisualization(cron) {
  const parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return '<span class="agent-schedule">' + esc(cron) + '</span>';
  const min = parts[0], hour = parts[1], dow = parts[4];
  const activeHours = parseHours(hour);
  const activeDays = parseDays(dow);
  const minute = min === '*' ? '??' : min.padStart(2, '0');
  const dayLabels = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
  const daysHtml = '<span class="day-labels">' + dayLabels.map(function(d, i) {
    return '<span class="day-label ' + (activeDays.has(i) ? 'active' : 'inactive') + '">' + d + '</span>';
  }).join('') + '</span>';
  const hoursHtml = Array.from({ length: 24 }, function(_, i) {
    return '<span class="hour-block ' + (activeHours.has(i) ? 'active' : 'inactive') + '"></span>';
  }).join('');
  return '<span class="agent-schedule-text">' + esc(scheduleHuman(cron)) + '</span>' +
    '<span class="agent-schedule-visual">' +
      '<span class="sched-legend">' +
        '<span class="sched-legend-item"><span class="day-block active"></span> on</span>' +
        '<span class="sched-legend-item"><span class="day-block inactive"></span> off</span>' +
      '</span>' +
      '<span class="day-labels-row">' + daysHtml + '</span>' +
      '<span class="hours-label">Hours</span>' +
      '<span class="agent-hours">' + hoursHtml + '</span>' +
      '<span class="agent-minute">:' + minute + '</span>' +
    '</span>';
}

function scheduleHuman(cron) {
  const parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return cron;
  const min = parts[0], hour = parts[1], dom = parts[2], mon = parts[3], dow = parts[4];
  if (dom === '*' && mon === '*' && dow === '*') {
    if (min.startsWith('*/')) return 'Every ' + min.slice(2) + ' minutes';
    if (hour.startsWith('*/')) {
      const h = hour.slice(2);
      if (min === '0') return 'Every ' + h + ' hours';
      return 'Every ' + h + 'h at :' + min;
    }
    if (hour !== '*' && min !== '*') return 'Daily at ' + hour + ':' + min.padStart(2, '0');
  }
  return cron;
}

function renderAgentCard(agent) {
  const card = document.createElement('div');
  card.className = 'agent-card' + (agent.enabled ? '' : ' agent-disabled');
  const dir = agent.directory.split('/').pop() || agent.directory;
  const header = agent.section_header || '';
  card.innerHTML =
    (header ? '<div class="agent-section-header">' + esc(header) + '</div>' : '') +
    '<div class="agent-header">' +
      '<span class="agent-project">' + esc(dir) + '</span>' +
      harnessLabel(agent.harness) +
      '<span class="agent-model">' + esc(agent.model) + '</span>' +
      '<span class="agent-schedule" title="' + esc(scheduleHuman(agent.schedule)) + '">' + renderTimingVisualization(agent.schedule) + '</span>' +
      '<span class="agent-status ' + (agent.enabled ? 'status-on' : 'status-off') + '">' + (agent.enabled ? 'ON' : 'OFF') + '</span>' +
      (agent.exit_code ? '<span class="agent-error-badge" title="' + esc(agent.last_error || '') + '">Error (exit ' + agent.exit_code + ')</span>' : '') +
      '<span class="agent-running-badge hidden"><span class="spinner"></span>Running…</span>' +
    '</div>' +
    '<div class="agent-actions">' +
      '<button class="btn-sm btn-run" data-id="' + esc(agent.id) + '">Run Now</button>' +
      '<button class="btn-sm btn-toggle" data-id="' + esc(agent.id) + '">' + (agent.enabled ? 'Disable' : 'Enable') + '</button>' +
      '<button class="btn-sm btn-edit" data-id="' + esc(agent.id) + '">Edit</button>' +
      '<button class="btn-sm btn-delete" data-id="' + esc(agent.id) + '">Delete</button>' +
      '<button class="btn-sm btn-log" data-id="' + esc(agent.id) + '">Logs</button>' +
    '</div>' +
    '<div class="agent-log-container hidden" id="agent-log-' + esc(agent.id) + '"></div>';
  return card;
}

async function loadAgents() {
  agentsListEl.innerHTML = '<p class="loading">loading agents…</p>';
  try {
    // Fetch both projects and agents
    const [projectsRes, agentsRes] = await Promise.all([
      fetch('/api/projects'),
      fetch('/api/agents')
    ]);
    
    if (!projectsRes.ok) throw new Error('HTTP ' + projectsRes.status);
    if (!agentsRes.ok) throw new Error('HTTP ' + agentsRes.status);
    
    const projects = await projectsRes.json();
    const data = await agentsRes.json();
    const agents = data.agents || [];
    
    agentsListEl.innerHTML = '';
    
    if (projects.length === 0) {
      agentsListEl.innerHTML = '<p class="loading">no projects found</p>';
      return;
    }
    
    // Create a map of project path to agents
    const agentsByProject = {};
    agents.forEach(function(a) {
      if (!agentsByProject[a.directory]) agentsByProject[a.directory] = [];
      agentsByProject[a.directory].push(a);
    });
    
    // Show all projects as sections
    projects.forEach(function(project) {
      const group = document.createElement('div');
      group.className = 'agent-group';
      group.innerHTML = '<h2 class="agent-group-title">' + esc(project.name) + '</h2>';
      
      const projectAgents = agentsByProject[project.path] || [];
      if (projectAgents.length === 0) {
        group.innerHTML += '<p class="loading">no agents for this project</p>';
      } else {
        projectAgents.forEach(function(a) { 
          group.appendChild(renderAgentCard(a)); 
        });
      }
      
      agentsListEl.appendChild(group);
    });
  } catch (err) {
    agentsListEl.innerHTML = '<p class="error">error: ' + esc(err.message) + '</p>';
  }
}

function closeInlineForm() {
  const el = document.querySelector('.agent-form-inline');
  if (el) el.remove();
}

agentsListEl.addEventListener('click', async function(e) {
  const btn = e.target.closest('button');
  if (!btn) return;
  const id = btn.dataset.id;

  if (btn.classList.contains('btn-toggle')) {
    btn.disabled = true;
    await fetch('/api/agents/' + encodeURIComponent(id) + '/toggle', { method: 'PATCH' });
    loadAgents();
  } else if (btn.classList.contains('btn-delete')) {
    if (!confirm('Delete this agent from crontab?')) return;
    btn.disabled = true;
    await fetch('/api/agents/' + encodeURIComponent(id), { method: 'DELETE' });
    loadAgents();
  } else if (btn.classList.contains('btn-run')) {
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> Running…';
    const card = btn.closest('.agent-card');
    const runningBadge = card.querySelector('.agent-running-badge');
    if (runningBadge) runningBadge.classList.remove('hidden');
    try {
      await fetch('/api/agents/' + encodeURIComponent(id) + '/trigger', { method: 'POST' });
      setTimeout(function() {
        btn.disabled = false;
        btn.textContent = 'Run Now';
        if (runningBadge) runningBadge.classList.add('hidden');
      }, 2000);
    } catch (err) {
      console.error('trigger error:', err);
      btn.disabled = false;
      btn.textContent = 'Run Now';
      if (runningBadge) runningBadge.classList.add('hidden');
    }
  } else if (btn.classList.contains('btn-log')) {
    const logEl = document.getElementById('agent-log-' + id);
    logEl.classList.toggle('hidden');
    if (!logEl.classList.contains('hidden') && !logEl.dataset.loaded) {
      logEl.innerHTML = '<p class="loading">loading log…</p>';
      try {
        const res = await fetch('/api/agents/' + encodeURIComponent(id) + '/log');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        const logInfo = await res.json();
        logEl.dataset.loaded = '1';
        if (!logInfo.exists) {
          logEl.innerHTML = '<p class="loading">No log file found</p>';
        } else {
          const lines = logInfo.lines || [];
          if (lines.length === 0) {
            logEl.innerHTML = '<p class="loading">Log file is empty</p>';
          } else {
            logEl.innerHTML = '<pre class="log-content">' + esc(lines.join('\n')) + '</pre>';
          }
          if (logInfo.last_run && logInfo.last_run !== '0001-01-01T00:00:00Z') {
            logEl.innerHTML += '<span class="log-last-run">Last run: ' + new Date(logInfo.last_run).toLocaleString() + '</span>';
          }
        }
      } catch (err) {
        logEl.dataset.loaded = '1';
        logEl.innerHTML = '<p class="error">Error loading log: ' + esc(err.message) + '</p>';
      }
    }
  } else if (btn.classList.contains('btn-edit')) {
    try {
      await showEditForm(id, btn);
    } catch (err) {
      console.error('edit error:', err);
      alert('Edit failed: ' + err.message);
    }
  }
});

addBtn.addEventListener('click', function() { showAddForm(); });

let cachedModels = null;
let cachedRepos = null;

async function fetchModels() {
  if (cachedModels) return cachedModels;
  try {
    const res = await fetch('/api/models');
    const data = await res.json();
    cachedModels = data.models || [];
  } catch (e) { cachedModels = []; }
  return cachedModels;
}

async function fetchRepos() {
  if (cachedRepos) return cachedRepos;
  try {
    const res = await fetch('/api/repos');
    const data = await res.json();
    cachedRepos = data.repos || [];
  } catch (e) { cachedRepos = []; }
  return cachedRepos;
}

const DAY_NAMES = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

function parseScheduleParts(cron) {
  const parts = (cron || '').trim().split(/\s+/);
  if (parts.length !== 5) return { minute: '30', hours: [], days: [] };
  const minute = parts[0] === '*' ? '0' : parts[0];
  const hourSet = parseHours(parts[1]);
  const daySet = parseDays(parts[4]);
  return {
    minute: minute,
    hours: Array.from(hourSet).sort(function(a, b) { return a - b; }),
    days: Array.from(daySet).sort(function(a, b) { return a - b; })
  };
}

function buildCronFromParts(minute, hours, days) {
  const hourStr = hours.length === 0 || hours.length === 24 ? '*' : hours.join(',');
  const dowStr = days.length === 0 || days.length === 7 ? '*' : days.join(',');
  return minute + ' ' + hourStr + ' * * ' + dowStr;
}

async function showAddForm() {
  closeInlineForm();
  const repos = await fetchRepos();
  const html = await buildForm('Add Agent', null, repos);
  const wrapper = document.createElement('div');
  wrapper.className = 'agent-form-inline';
  wrapper.innerHTML = html;
  agentsListEl.prepend(wrapper);
  attachProjectSelectHandler(wrapper);
  wrapper.querySelector('input, select')?.focus();
}

async function showEditForm(id, btn) {
  closeInlineForm();
  const agentsRes = await fetch('/api/agents').then(function(r) { return r.json(); });
  const agent = agentsRes.agents.find(function(a) { return a.id === id; });
  if (!agent) { alert('Agent not found with id ' + id); return; }
  const repos = await fetchRepos();
  const html = await buildForm('Edit Agent', agent, repos, id);
  const card = btn.closest('.agent-card');
  const wrapper = document.createElement('div');
  wrapper.className = 'agent-form-inline';
  wrapper.innerHTML = html;
  if (card) { card.after(wrapper); } else { agentsListEl.prepend(wrapper); }
  attachProjectSelectHandler(wrapper);
  wrapper.querySelector('input, select')?.focus();
  wrapper.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

function attachProjectSelectHandler(wrapper) {
  const select = wrapper.querySelector('#project-select');
  const sectionHeaderField = wrapper.querySelector('#section-header-field');
  if (select && sectionHeaderField) {
    select.addEventListener('change', function() {
      const selected = select.options[select.selectedIndex];
      if (selected && selected.dataset.name) {
        sectionHeaderField.value = selected.dataset.name;
      }
    });
  }
}

async function buildForm(title, agent, repos, editId) {
  const models = await fetchModels();
  const sched = parseScheduleParts(agent ? agent.schedule : '');

  const dirOptions = repos.map(function(p) {
    const sel = agent && agent.directory === p.path ? ' selected' : '';
    return '<option value="' + esc(p.path) + '" data-name="' + esc(p.name) + '"' + sel + '>' + esc(p.name) + '</option>';
  }).join('');

  const currentModel = agent ? (agent.model || '') : '';

  const binaryPath = agent ? (agent.binary_path || '') : '';
  const prompt = agent ? (agent.prompt || '') : '';
  const logPath = agent ? (agent.log_path || '') : '';
  const sectionHeader = agent ? (agent.section_header || '') : '';

  const minuteVal = agent ? sched.minute : '30';
  const selectedHours = agent ? sched.hours : [3, 7, 11, 15, 23];
  const selectedDays = agent ? sched.days : [0, 1, 2, 3, 4, 5, 6];

  let hourCheckboxes = '';
  for (let h = 0; h < 24; h++) {
    const checked = selectedDays.length === 0 || selectedHours.indexOf(h) !== -1 ? ' checked' : '';
    hourCheckboxes += '<label class="hour-cb"><input type="checkbox" name="hours" value="' + h + '"' + checked + '>' + String(h).padStart(2, '0') + '</label>';
  }

  let dayCheckboxes = '';
  for (let d = 0; d < 7; d++) {
    const chk = selectedDays.length === 0 || selectedDays.indexOf(d) !== -1 ? ' checked' : '';
    dayCheckboxes += '<label class="day-cb"><input type="checkbox" name="days" value="' + d + '"' + chk + '>' + DAY_NAMES[d] + '</label>';
  }

  const modelsJSON = esc(JSON.stringify(models));

  return '<div class="agent-form">' +
    '<h3>' + title + '</h3>' +
    '<form class="agent-crud-form" data-edit-id="' + (editId !== undefined ? esc(editId) : '') + '">' +
      '<input type="hidden" name="section_header" id="section-header-field" value="' + esc(sectionHeader) + '">' +
      '<fieldset class="sched-fieldset"><legend>Schedule</legend>' +
        '<div class="sched-row"><label>Minute <input name="minute" type="number" min="0" max="59" value="' + esc(minuteVal) + '"></label></div>' +
        '<div class="sched-row"><span class="sched-label">Hours</span><div class="hour-grid">' + hourCheckboxes + '</div></div>' +
        '<div class="sched-row"><span class="sched-label">Days</span><div class="day-grid">' + dayCheckboxes + '</div></div>' +
      '</fieldset>' +
      '<label>Project <select name="directory" id="project-select"><option value="">Select project…</option>' + dirOptions + '</select></label>' +
      '<label>Model</label>' +
      '<div class="model-picker" data-models="' + modelsJSON + '" data-current="' + esc(currentModel) + '">' +
        '<input type="hidden" name="model" value="' + esc(currentModel) + '">' +
        '<input type="text" class="model-search" placeholder="Search models…" value="' + esc(currentModel) + '" autocomplete="off">' +
        '<div class="model-dropdown"></div>' +
      '</div>' +
      '<label>Prompt <input name="prompt" value="' + esc(prompt) + '" placeholder="conductor/autonomous_prompt.md"></label>' +
      '<label>Log Path <input name="logPath" value="' + esc(logPath) + '" placeholder="/path/to/output.log"></label>' +
      '<input type="hidden" name="binary_path" value="' + esc(binaryPath) + '">' +
      '<div class="form-actions">' +
        '<button type="submit" class="btn">Save</button>' +
        '<button type="button" class="btn btn-cancel">Cancel</button>' +
      '</div>' +
    '</form>' +
    '</div>';
}

document.addEventListener('click', function(e) {
  if (e.target.classList.contains('btn-cancel')) closeInlineForm();
});

// ── Model picker ──

function renderModelDropdown(picker) {
  const models = JSON.parse(picker.dataset.models);
  const query = (picker.querySelector('.model-search').value || '').toLowerCase();
  const dropdown = picker.querySelector('.model-dropdown');

  const filtered = query
    ? models.filter(function(m) { return m.toLowerCase().includes(query); })
    : models;

  if (filtered.length === 0) {
    dropdown.innerHTML = '<div class="model-option model-empty">No matches</div>';
  } else {
    dropdown.innerHTML = filtered.slice(0, 50).map(function(m) {
      return '<div class="model-option" data-value="' + esc(m) + '">' + esc(m) + '</div>';
    }).join('');
    if (filtered.length > 50) {
      dropdown.innerHTML += '<div class="model-option model-empty">+' + (filtered.length - 50) + ' more…</div>';
    }
  }
  dropdown.classList.add('open');
}

document.addEventListener('focusin', function(e) {
  if (!e.target.classList.contains('model-search')) return;
  const picker = e.target.closest('.model-picker');
  if (picker) renderModelDropdown(picker);
});

document.addEventListener('input', function(e) {
  if (!e.target.classList.contains('model-search')) return;
  const picker = e.target.closest('.model-picker');
  if (picker) renderModelDropdown(picker);
});

document.addEventListener('click', function(e) {
  const option = e.target.closest('.model-option');
  if (option && option.dataset.value) {
    const picker = option.closest('.model-picker');
    picker.querySelector('input[name="model"]').value = option.dataset.value;
    picker.querySelector('.model-search').value = option.dataset.value;
    picker.querySelector('.model-dropdown').classList.remove('open');
    e.stopPropagation();
    return;
  }

  // Close dropdowns when clicking outside
  if (!e.target.closest('.model-picker')) {
    document.querySelectorAll('.model-dropdown.open').forEach(function(d) {
      d.classList.remove('open');
    });
  }
});

document.addEventListener('submit', async function(e) {
  if (!e.target.classList.contains('agent-crud-form')) return;
  e.preventDefault();
  const form = e.target;
  const fd = new FormData(form);

  const hours = fd.getAll('hours').map(Number);
  const days = fd.getAll('days').map(Number);
  const schedule = buildCronFromParts(fd.get('minute') || '30', hours, days);

  const body = {
    schedule: schedule,
    directory: fd.get('directory'),
    harness: 'opencode',
    binary_path: fd.get('binary_path'),
    model: fd.get('model'),
    prompt: fd.get('prompt'),
    log_path: fd.get('logPath'),
    section_header: fd.get('section_header')
  };

  const editId = form.dataset.editId;
  const url = '/api/agents';
  const method = 'POST';
  if (editId !== undefined && editId !== '') {
    url = '/api/agents/' + encodeURIComponent(editId);
    method = 'PUT';
  }

  const res = await fetch(url, {
    method: method,
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });

  if (!res.ok) {
    alert('Error: ' + await res.text());
    return;
  }

  closeInlineForm();
  loadAgents();
});

loadAgents();
