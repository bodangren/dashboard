'use strict';

function getHealthStatus(health) {
  if (!health) return 'unknown';
  const total = health.dirty?.total || 0;
  const stale = health.staleBranches?.count || 0;
  const diverge = (health.divergence?.ahead || 0) + (health.divergence?.behind || 0);

  if (total === 0 && stale === 0 && diverge === 0) return 'healthy';
  if (total <= 3 && stale <= 1 && diverge <= 2) return 'warning';
  return 'critical';
}

function formatHealthTooltip(health) {
  if (!health) return 'Health: unknown';
  const parts = [];
  const dirty = health.dirty || {};
  if (dirty.total > 0) {
    parts.push(`${dirty.total} uncommitted (${dirty.staged} staged, ${dirty.modified} modified, ${dirty.untracked} untracked)`);
  }
  const div = health.divergence || {};
  if (div.ahead > 0) parts.push(`${div.ahead} commits ahead of remote`);
  if (div.behind > 0) parts.push(`${div.behind} commits behind remote`);
  const stale = health.staleBranches || {};
  if (stale.count > 0) {
    const names = stale.branches?.join(', ') || `${stale.count} stale`;
    parts.push(`${stale.count} stale branch${stale.count > 1 ? 'es' : ''}: ${names}`);
  }
  if (parts.length === 0) return 'Health: clean';
  return 'Health: ' + parts.join('; ');
}

function renderHealthBadge(health) {
  const badge = document.createElement('div');
  badge.className = 'health-badge';

  const status = getHealthStatus(health);
  badge.classList.add(`health-${status}`);

  const icon = document.createElement('span');
  icon.className = 'health-icon';
  icon.textContent = status === 'healthy' ? '✓' : status === 'warning' ? '⚠' : status === 'critical' ? '✗' : '?';
  badge.appendChild(icon);

  const label = document.createElement('span');
  label.className = 'health-label';
  label.textContent = status.charAt(0).toUpperCase() + status.slice(1);
  badge.appendChild(label);

  const tooltip = formatHealthTooltip(health);
  badge.title = tooltip;

  return badge;
}

function attachHealthBadge(card, health) {
  const header = card.querySelector('.project-header');
  if (!header) return;

  const existing = header.querySelector('.health-badge');
  if (existing) existing.remove();

  const badge = renderHealthBadge(health);
  header.appendChild(badge);
}