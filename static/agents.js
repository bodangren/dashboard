'use strict';

const agentsListEl = document.getElementById('agents-list');
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
  return '<span class="harness-badge" style="color:' + (colors[h] || '#ccc') + '">' + esc(h) + '</span>';
}

function parseHours(hourStr) {
  if (hourStr === '*') return new Set(Array.from({ length: 24 }, function(_, i) { return i; }));
  var hours = new Set();
  hourStr.split(',').forEach(function(part) {
    if (part.includes('/')) {
      var step = part.split('/')[1];
      for (var i = 0; i < 24; i += parseInt(step)) hours.add(i);
    } else if (part.includes('-')) {
      var bounds = part.split('-').map(Number);
      for (var i = bounds[0]; i <= bounds[1]; i++) hours.add(i);
    } else {
      hours.add(parseInt(part));
    }
  });
  return hours;
}

function parseDays(dow) {
  if (dow === '*') return new Set([0, 1, 2, 3, 4, 5, 6]);
  var days = new Set();
  dow.split(',').forEach(function(part) {
    if (part.includes('-')) {
      var bounds = part.split('-').map(Number);
      for (var i = bounds[0]; i <= bounds[1]; i++) days.add(i);
    } else {
      days.add(parseInt(part));
    }
  });
  return days;
}

function renderTimingVisualization(cron) {
  var parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return '<span class="agent-schedule">' + esc(cron) + '</span>';
  var min = parts[0], hour = parts[1], dow = parts[4];
  var activeHours = parseHours(hour);
  var activeDays = parseDays(dow);
  var minute = min === '*' ? '??' : min.padStart(2, '0');
  var daysHtml = Array.from({ length: 7 }, function(_, i) {
    return '<span class="day-block ' + (activeDays.has(i) ? 'active' : 'inactive') + '"></span>';
  }).join('');
  var hoursHtml = Array.from({ length: 24 }, function(_, i) {
    return '<span class="hour-block ' + (activeHours.has(i) ? 'active' : 'inactive') + '"></span>';
  }).join('');
  return '<span class="agent-days">' + daysHtml + '</span>' +
    '<span class="agent-hours">' + hoursHtml + '</span>' +
    '<span class="agent-minute">:' + minute + '</span>';
}

function scheduleHuman(cron) {
  var parts = cron.trim().split(/\s+/);
  if (parts.length !== 5) return cron;
  var min = parts[0], hour = parts[1], dom = parts[2], mon = parts[3], dow = parts[4];
  if (dom === '*' && mon === '*' && dow === '*') {
    if (min.startsWith('*/')) return 'Every ' + min.slice(2) + ' minutes';
    if (hour.startsWith('*/')) {
      var h = hour.slice(2);
      if (min === '0') return 'Every ' + h + ' hours';
      return 'Every ' + h + 'h at :' + min;
    }
    if (hour !== '*' && min !== '*') return 'Daily at ' + hour + ':' + min.padStart(2, '0');
  }
  return cron;
}

function renderAgentCard(agent, index) {
  var card = document.createElement('div');
  card.className = 'agent-card' + (agent.enabled ? '' : ' agent-disabled');
  var dir = agent.directory.split('/').pop() || agent.directory;
  var header = agent.section_header || '';
  card.innerHTML =
    (header ? '<div class="agent-section-header">' + esc(header) + '</div>' : '') +
    '<div class="agent-header">' +
      '<span class="agent-project">' + esc(dir) + '</span>' +
      harnessLabel(agent.harness) +
      '<span class="agent-model">' + esc(agent.model) + '</span>' +
      '<span class="agent-schedule" title="' + esc(scheduleHuman(agent.schedule)) + '">' + renderTimingVisualization(agent.schedule) + '</span>' +
      '<span class="agent-status ' + (agent.enabled ? 'status-on' : 'status-off') + '">' + (agent.enabled ? 'ON' : 'OFF') + '</span>' +
    '</div>' +
    '<div class="agent-actions">' +
      '<button class="btn-sm btn-toggle" data-idx="' + index + '">' + (agent.enabled ? 'Disable' : 'Enable') + '</button>' +
      '<button class="btn-sm btn-edit" data-idx="' + index + '">Edit</button>' +
      '<button class="btn-sm btn-delete" data-idx="' + index + '">Delete</button>' +
      '<button class="btn-sm btn-log" data-idx="' + index + '">Logs</button>' +
    '</div>' +
    '<div class="agent-log-container hidden" id="agent-log-' + index + '"></div>';
  return card;
}

async function loadAgents() {
  agentsListEl.innerHTML = '<p class="loading">loading agents…</p>';
  try {
    // Fetch both projects and agents
    var [projectsRes, agentsRes] = await Promise.all([
      fetch('/api/projects'),
      fetch('/api/agents')
    ]);
    
    if (!projectsRes.ok) throw new Error('HTTP ' + projectsRes.status);
    if (!agentsRes.ok) throw new Error('HTTP ' + agentsRes.status);
    
    var projects = await projectsRes.json();
    var data = await agentsRes.json();
    var agents = data.agents || [];
    
    agentsListEl.innerHTML = '';
    
    if (projects.length === 0) {
      agentsListEl.innerHTML = '<p class="loading">no projects found</p>';
      return;
    }
    
    // Create a map of project path to agents
    var agentsByProject = {};
    agents.forEach(function(a, i) {
      if (!agentsByProject[a.directory]) agentsByProject[a.directory] = [];
      agentsByProject[a.directory].push(Object.assign({}, a, { _idx: i }));
    });
    
    // Show all projects as sections
    projects.forEach(function(project) {
      var group = document.createElement('div');
      group.className = 'agent-group';
      group.innerHTML = '<h2 class="agent-group-title">' + esc(project.name) + '</h2>';
      
      var projectAgents = agentsByProject[project.path] || [];
      if (projectAgents.length === 0) {
        group.innerHTML += '<p class="loading">no agents for this project</p>';
      } else {
        projectAgents.forEach(function(a) { 
          group.appendChild(renderAgentCard(a, a._idx)); 
        });
      }
      
      agentsListEl.appendChild(group);
    });
  } catch (err) {
    agentsListEl.innerHTML = '<p class="error">error: ' + esc(err.message) + '</p>';
  }
}

function closeInlineForm() {
  var el = document.querySelector('.agent-form-inline');
  if (el) el.remove();
}

agentsListEl.addEventListener('click', async function(e) {
  var btn = e.target.closest('button');
  if (!btn) return;
  var idx = btn.dataset.idx;

  if (btn.classList.contains('btn-toggle')) {
    btn.disabled = true;
    await fetch('/api/agents/' + idx + '/toggle', { method: 'PATCH' });
    loadAgents();
  } else if (btn.classList.contains('btn-delete')) {
    if (!confirm('Delete this agent from crontab?')) return;
    btn.disabled = true;
    await fetch('/api/agents/' + idx, { method: 'DELETE' });
    loadAgents();
  } else if (btn.classList.contains('btn-log')) {
    var logEl = document.getElementById('agent-log-' + idx);
    logEl.classList.toggle('hidden');
    if (!logEl.classList.contains('hidden') && !logEl.dataset.loaded) {
      logEl.innerHTML = '<p class="loading">loading log…</p>';
      try {
        var res = await fetch('/api/agents/' + idx + '/log');
        if (!res.ok) throw new Error('HTTP ' + res.status);
        var logInfo = await res.json();
        logEl.dataset.loaded = '1';
        if (!logInfo.exists) {
          logEl.innerHTML = '<p class="loading">No log file found</p>';
        } else {
          var lines = logInfo.lines || [];
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
      await showEditForm(idx, btn);
    } catch (err) {
      console.error('edit error:', err);
      alert('Edit failed: ' + err.message);
    }
  }
});

addBtn.addEventListener('click', function() { showAddForm(); });

var cachedModels = null;
var cachedRepos = null;

async function fetchModels() {
  if (cachedModels) return cachedModels;
  try {
    var res = await fetch('/api/models');
    var data = await res.json();
    cachedModels = data.models || [];
  } catch (e) { cachedModels = []; }
  return cachedModels;
}

async function fetchRepos() {
  if (cachedRepos) return cachedRepos;
  try {
    var res = await fetch('/api/projects');
    cachedRepos = await res.json();
  } catch (e) { cachedRepos = []; }
  return cachedRepos;
}

var DAY_NAMES = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

function parseScheduleParts(cron) {
  var parts = (cron || '').trim().split(/\s+/);
  if (parts.length !== 5) return { minute: '30', hours: [], days: [] };
  var minute = parts[0] === '*' ? '0' : parts[0];
  var hourSet = parseHours(parts[1]);
  var daySet = parseDays(parts[4]);
  return {
    minute: minute,
    hours: Array.from(hourSet).sort(function(a, b) { return a - b; }),
    days: Array.from(daySet).sort(function(a, b) { return a - b; })
  };
}

function buildCronFromParts(minute, hours, days) {
  var hourStr = hours.length === 0 || hours.length === 24 ? '*' : hours.join(',');
  var dowStr = days.length === 0 || days.length === 7 ? '*' : days.join(',');
  return minute + ' ' + hourStr + ' * * ' + dowStr;
}

async function showAddForm() {
  closeInlineForm();
  var repos = await fetchRepos();
  var html = await buildForm('Add Agent', null, repos);
  var wrapper = document.createElement('div');
  wrapper.className = 'agent-form-inline';
  wrapper.innerHTML = html;
  agentsListEl.prepend(wrapper);
  attachProjectSelectHandler(wrapper);
  wrapper.querySelector('input, select')?.focus();
}

async function showEditForm(idx, btn) {
  closeInlineForm();
  var agentsRes = await fetch('/api/agents').then(function(r) { return r.json(); });
  var agent = agentsRes.agents[idx];
  if (!agent) { alert('Agent not found at index ' + idx); return; }
  var repos = await fetchRepos();
  var html = await buildForm('Edit Agent', agent, repos, idx);
  var card = btn.closest('.agent-card');
  var wrapper = document.createElement('div');
  wrapper.className = 'agent-form-inline';
  wrapper.innerHTML = html;
  if (card) { card.after(wrapper); } else { agentsListEl.prepend(wrapper); }
  attachProjectSelectHandler(wrapper);
  wrapper.querySelector('input, select')?.focus();
  wrapper.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

function attachProjectSelectHandler(wrapper) {
  var select = wrapper.querySelector('#project-select');
  var sectionHeaderField = wrapper.querySelector('#section-header-field');
  if (select && sectionHeaderField) {
    select.addEventListener('change', function() {
      var selected = select.options[select.selectedIndex];
      if (selected && selected.dataset.name) {
        sectionHeaderField.value = selected.dataset.name;
      }
    });
  }
}

async function buildForm(title, agent, repos, editIdx) {
  var models = await fetchModels();
  var sched = parseScheduleParts(agent ? agent.schedule : '');

  var dirOptions = repos.map(function(p) {
    var sel = agent && agent.directory === p.path ? ' selected' : '';
    return '<option value="' + esc(p.path) + '" data-name="' + esc(p.name) + '"' + sel + '>' + esc(p.name) + '</option>';
  }).join('');

  var currentModel = agent ? (agent.model || '') : '';

  var binaryPath = agent ? (agent.binary_path || '') : '';
  var defaultBinary = '/home/daniel-bo/.nvm/versions/node/v24.4.0/bin/opencode';
  var prompt = agent ? (agent.prompt || '') : '';
  var logPath = agent ? (agent.log_path || '') : '';
  var sectionHeader = agent ? (agent.section_header || '') : '';

  var minuteVal = agent ? sched.minute : '30';
  var selectedHours = agent ? sched.hours : [3, 7, 11, 15, 23];
  var selectedDays = agent ? sched.days : [0, 1, 2, 3, 4, 5, 6];

  var hourCheckboxes = '';
  for (var h = 0; h < 24; h++) {
    var checked = selectedDays.length === 0 || selectedHours.indexOf(h) !== -1 ? ' checked' : '';
    hourCheckboxes += '<label class="hour-cb"><input type="checkbox" name="hours" value="' + h + '"' + checked + '>' + String(h).padStart(2, '0') + '</label>';
  }

  var dayCheckboxes = '';
  for (var d = 0; d < 7; d++) {
    var chk = selectedDays.length === 0 || selectedDays.indexOf(d) !== -1 ? ' checked' : '';
    dayCheckboxes += '<label class="day-cb"><input type="checkbox" name="days" value="' + d + '"' + chk + '>' + DAY_NAMES[d] + '</label>';
  }

  var modelsJSON = esc(JSON.stringify(models));

  return '<div class="agent-form">' +
    '<h3>' + title + '</h3>' +
    '<form class="agent-crud-form" data-edit-idx="' + (editIdx !== undefined ? editIdx : '') + '">' +
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
      '<input type="hidden" name="binary_path" value="' + esc(binaryPath || defaultBinary) + '">' +
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
  var models = JSON.parse(picker.dataset.models);
  var query = (picker.querySelector('.model-search').value || '').toLowerCase();
  var dropdown = picker.querySelector('.model-dropdown');

  var filtered = query
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
  var picker = e.target.closest('.model-picker');
  if (picker) renderModelDropdown(picker);
});

document.addEventListener('input', function(e) {
  if (!e.target.classList.contains('model-search')) return;
  var picker = e.target.closest('.model-picker');
  if (picker) renderModelDropdown(picker);
});

document.addEventListener('click', function(e) {
  var option = e.target.closest('.model-option');
  if (option && option.dataset.value) {
    var picker = option.closest('.model-picker');
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
  var form = e.target;
  var fd = new FormData(form);

  var hours = fd.getAll('hours').map(Number);
  var days = fd.getAll('days').map(Number);
  var schedule = buildCronFromParts(fd.get('minute') || '30', hours, days);

  var body = {
    schedule: schedule,
    directory: fd.get('directory'),
    harness: 'opencode',
    binary_path: fd.get('binary_path'),
    model: fd.get('model'),
    prompt: fd.get('prompt'),
    log_path: fd.get('logPath'),
    section_header: fd.get('section_header')
  };

  var editIdx = form.dataset.editIdx;
  var url = '/api/agents';
  var method = 'POST';
  if (editIdx !== undefined && editIdx !== '') {
    url = '/api/agents/' + editIdx;
    method = 'PUT';
  }

  var res = await fetch(url, {
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
