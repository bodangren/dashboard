'use strict';

const PREFERENCES_KEY = 'dashboard-preferences';

const DEFAULT_PREFERENCES = {
  density: 'comfortable',
  sortOrder: 'recent',
  showActivity: true,
  showAgents: true,
  showPullStatus: true,
  groupByDirectory: false,
};

function getPreferences() {
  try {
    const stored = localStorage.getItem(PREFERENCES_KEY);
    if (stored) {
      const parsed = JSON.parse(stored);
      return { ...DEFAULT_PREFERENCES, ...parsed };
    }
  } catch {
  }
  return { ...DEFAULT_PREFERENCES };
}

function savePreferences(prefs) {
  try {
    localStorage.setItem(PREFERENCES_KEY, JSON.stringify(prefs));
  } catch {
  }
}

function updatePreference(key, value) {
  const prefs = getPreferences();
  prefs[key] = value;
  savePreferences(prefs);
  applyPreferences(prefs);
  return prefs;
}

function applyPreferences(prefs) {
  const root = document.documentElement;
  root.setAttribute('data-density', prefs.density);
  root.setAttribute('data-sort', prefs.sortOrder);

  const activitySection = document.getElementById('activity-section');
  const agentsSection = document.getElementById('agents-section');
  const pullStatusSection = document.getElementById('pull-status-section');

  if (activitySection) {
    activitySection.classList.toggle('hidden', !prefs.showActivity);
  }
  if (agentsSection) {
    agentsSection.classList.toggle('hidden', !prefs.showAgents);
  }
  if (pullStatusSection) {
    pullStatusSection.classList.toggle('hidden', !prefs.showPullStatus);
  }
}

function loadPreferences() {
  const prefs = getPreferences();
  applyPreferences(prefs);
  return prefs;
}